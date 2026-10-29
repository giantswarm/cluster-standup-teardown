package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"time"

	. "github.com/onsi/gomega" // nolint:staticcheck

	fluxhelmv2 "github.com/fluxcd/helm-controller/api/v2"
	"github.com/giantswarm/clustertest/v5"
	"github.com/giantswarm/clustertest/v5/pkg/application"
	"github.com/giantswarm/clustertest/v5/pkg/client"
	"github.com/giantswarm/clustertest/v5/pkg/helmrelease"
	"github.com/giantswarm/clustertest/v5/pkg/organization"
	"github.com/giantswarm/clustertest/v5/pkg/utils"
	"github.com/giantswarm/clustertest/v5/pkg/wait"
	"github.com/spf13/cobra"
	apitypes "k8s.io/apimachinery/pkg/types"
	cr "sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/giantswarm/cluster-standup-teardown/v6/cmd/standup/types"
	cb "github.com/giantswarm/cluster-standup-teardown/v6/pkg/clusterbuilder"
	"github.com/giantswarm/cluster-standup-teardown/v6/pkg/clusterbuilder/providers/capz"
	"github.com/giantswarm/cluster-standup-teardown/v6/pkg/standup"
	"github.com/giantswarm/cluster-standup-teardown/v6/pkg/values"
)

var (
	standupCmd = &cobra.Command{
		Use:     "standup",
		Long:    "Standup create a test workload cluster in a standard, reproducible way.\nA valid Management Cluster kubeconfig must be available and set to the `E2E_KUBECONFIG` environment variable.",
		Example: "standup --provider aws --context capa",
		Args:    cobra.NoArgs,
		RunE:    run,
	}

	provider        string
	kubeContext     string
	clusterValues   string
	releaseVersion  string
	releaseCommit   string
	clusterVersion  string
	outputDirectory string

	controlPlaneNodes int
	workerNodes       int

	waitForApps bool

	// Functions to run after cluster creation to confirm it is up and ready to use
	clusterReadyFns []func(wcClient *client.Client) = []func(wcClient *client.Client){
		func(wcClient *client.Client) {
			// Control plane nodes are joined one after another, so the timeout scales with their number. The first
			// node is already registered once the WC API is reachable and only needs a few minutes to become ready,
			// each further one needs a full machine provisioning, which can take up to 10 minutes on slower providers
			// (e.g. Cloud Director).
			_ = wait.For(
				wait.AreNumNodesReady(context.Background(), wcClient, controlPlaneNodes, &cr.MatchingLabels{"node-role.kubernetes.io/control-plane": ""}),
				wait.WithTimeout(time.Duration(controlPlaneNodes)*10*time.Minute),
				wait.WithInterval(15*time.Second),
			)
		},
		func(wcClient *client.Client) {
			_ = wait.For(
				wait.AreNumNodesReady(context.Background(), wcClient, workerNodes, client.DoesNotHaveLabels{"node-role.kubernetes.io/control-plane"}),
				wait.WithTimeout(20*time.Minute),
				wait.WithInterval(15*time.Second),
			)
		},
	}
)

func init() {
	standupCmd.Flags().StringVar(&provider, "provider", "", "The provider (required)")
	standupCmd.Flags().StringVar(&kubeContext, "context", "", "The kubernetes context to use (required)")

	standupCmd.Flags().StringVar(&clusterValues, "cluster-values", "", "The path to the cluster app values")
	standupCmd.Flags().IntVar(&controlPlaneNodes, "control-plane-nodes", 1, "The number of control plane nodes to wait for being ready")
	standupCmd.Flags().IntVar(&workerNodes, "worker-nodes", 1, "The number of worker nodes to wait for being ready")
	standupCmd.Flags().StringVar(&outputDirectory, "output", "./", "The directory to store the results.json and kubeconfig in")

	standupCmd.Flags().StringVar(&releaseVersion, "release", application.ReleaseLatest, "The version of the Release to use")
	standupCmd.Flags().StringVar(&releaseCommit, "release-commit", "", "The git commit to get the Release version from (defaults to main default if unset)")
	standupCmd.Flags().StringVar(&clusterVersion, "cluster-version", "latest", "The version of the cluster app to install")

	standupCmd.Flags().BoolVar(&waitForApps, "wait-for-apps-ready", false, "Wait until all default apps are installed")

	_ = standupCmd.MarkFlagRequired("provider")
	_ = standupCmd.MarkFlagRequired("context")
}

func main() {
	if err := standupCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) error {
	// Required to be able to use our module with Gomega assertions without Ginkgo
	RegisterFailHandler(func(message string, callerSkip ...int) {
		panic(message)
	})

	cmd.SilenceUsage = true

	ctx := context.Background()

	framework, err := clustertest.New(kubeContext)
	if err != nil {
		return err
	}

	provider := application.Provider(provider)
	clusterName := utils.GenerateRandomName("t")
	orgName := utils.GenerateRandomName("t")

	clusterValuesOverrides := []string{values.MustLoadValuesFile(clusterValues)}

	var cluster *application.Cluster
	clusterBuilder, err := cb.GetClusterBuilderForContext(kubeContext)
	if err != nil {
		fmt.Printf("Failed to automatically get cluster builder based on context, falling back to building out a cluster...\n")
		cluster = application.NewClusterApp(clusterName, provider).
			WithRelease(application.ReleasePair{
				Version: releaseVersion,
				Commit:  releaseCommit,
			}).
			WithAppVersions(clusterVersion).WithOrg(organization.New(orgName)).
			WithAppValuesFile(path.Clean(clusterValues), &application.TemplateValues{
				ClusterName:  clusterName,
				Organization: orgName,
			})
	} else {
		cluster = clusterBuilder.NewClusterApp(clusterName, orgName, clusterValuesOverrides).
			WithRelease(application.ReleasePair{
				Version: releaseVersion,
				Commit:  releaseCommit,
			}).
			WithAppVersions(clusterVersion)
	}

	if provider == application.ProviderEKS || provider == capz.ProviderAKS {
		// As EKS and AKS have a managed control plane (no control-plane nodes are
		// visible in the workload cluster) we only check for worker nodes being ready
		clusterReadyFns = []func(wcClient *client.Client){func(wcClient *client.Client) {
			_ = wait.For(
				wait.AreNumNodesReady(context.Background(), wcClient, workerNodes, client.DoesNotHaveLabels{"node-role.kubernetes.io/control-plane"}),
				wait.WithTimeout(20*time.Minute),
				wait.WithInterval(15*time.Second),
			)
		},
		}
	}

	if waitForApps {
		// Additionally wait for all HelmReleases of the cluster to be ready once the nodes are ready
		clusterReadyFns = append(clusterReadyFns, func(wcClient *client.Client) {
			// wait.For calls the condition below every interval until it returns true or an error, or the timeout is
			// reached. The result is ignored, so standup continues even if the HelmReleases never become ready.
			_ = wait.For(
				func() (bool, error) {
					// List the HelmReleases created by the cluster chart. They already exist by the time the nodes are
					// ready, as they are created alongside the Cluster resource in the same Helm install.
					helmReleases := &fluxhelmv2.HelmReleaseList{}
					err := framework.MC().List(context.Background(), helmReleases, cr.InNamespace(cluster.GetNamespace()), cr.MatchingLabels{"giantswarm.io/cluster": cluster.Name})
					if err != nil {
						// Returning an error stops the wait
						return false, err
					}

					// Collect the name and namespace of each HelmRelease, as AreAllReady expects
					namespacedNames := []apitypes.NamespacedName{}
					for _, helmRelease := range helmReleases.Items {
						namespacedNames = append(namespacedNames, apitypes.NamespacedName{Name: helmRelease.Name, Namespace: helmRelease.Namespace})
					}

					// AreAllReady builds a check that logs the status of each HelmRelease and returns an error if any
					// isn't ready yet. We call it right away and treat that error as "not done" rather than as a
					// failure, so wait.For keeps polling.
					return helmrelease.AreAllReady(context.Background(), framework.MC(), namespacedNames)() == nil, nil
				},
				wait.WithTimeout(20*time.Minute),
				wait.WithInterval(15*time.Second),
			)
		})
	}

	// Create the results file with the details we have already incase the cluster creation fails
	result := types.StandupResult{
		Provider:       string(provider),
		ClusterName:    clusterName,
		OrgName:        orgName,
		Namespace:      cluster.GetNamespace(),
		ClusterVersion: cluster.ClusterApp.Version,
		KubeconfigPath: "",
	}

	resultsFile, err := os.Create(path.Join(outputDirectory, "results.json")) // nolint:gosec
	if err != nil {
		return err
	}

	resultBytes, err := json.Marshal(result)
	if err != nil {
		return err
	}
	_, err = resultsFile.Write(resultBytes)
	if err != nil {
		return err
	}

	resultsFile.Close() // nolint:errcheck,gosec

	fmt.Printf("Standing up cluster...\n\nProvider:\t\t%s\nCluster Name:\t\t%s\nOrg Name:\t\t%s\nResults Directory:\t%s\n\n", provider, clusterName, orgName, outputDirectory)

	cluster, err = standup.New(framework, false, clusterReadyFns...).Standup(cluster)
	if err != nil {
		return err
	}

	// Save the kubeconfig for the WC
	kubeconfigFile, err := os.Create(path.Join(outputDirectory, "kubeconfig")) // nolint:gosec
	if err != nil {
		return err
	}
	defer kubeconfigFile.Close() // nolint:errcheck

	kubeconfig, err := framework.MC().GetClusterKubeConfig(ctx, cluster.Name, cluster.GetNamespace())
	if err != nil {
		return err
	}
	_, err = kubeconfigFile.Write([]byte(kubeconfig))
	if err != nil {
		return err
	}

	// Update the results file with the kubeconfig path
	result.KubeconfigPath = kubeconfigFile.Name()

	resultsFile, err = os.Create(path.Join(outputDirectory, "results.json")) // nolint:gosec
	if err != nil {
		return err
	}
	defer resultsFile.Close() // nolint:errcheck

	resultBytes, err = json.Marshal(result)
	if err != nil {
		return err
	}
	_, err = resultsFile.Write(resultBytes)

	return err
}
