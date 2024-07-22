package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"time"

	. "github.com/onsi/gomega"

	"github.com/giantswarm/apiextensions-application/api/v1alpha1"
	"github.com/giantswarm/clustertest"
	"github.com/giantswarm/clustertest/pkg/application"
	"github.com/giantswarm/clustertest/pkg/client"
	"github.com/giantswarm/clustertest/pkg/organization"
	"github.com/giantswarm/clustertest/pkg/utils"
	"github.com/giantswarm/clustertest/pkg/wait"
	"github.com/spf13/cobra"
	apitypes "k8s.io/apimachinery/pkg/types"
	cr "sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/giantswarm/cluster-standup-teardown/cmd/standup/types"
	cb "github.com/giantswarm/cluster-standup-teardown/pkg/clusterbuilder"
	"github.com/giantswarm/cluster-standup-teardown/pkg/standup"
	"github.com/giantswarm/cluster-standup-teardown/pkg/values"
)

var (
	standupCmd = &cobra.Command{
		Use:     "standup",
		Long:    "Standup create a test workload cluster in a standard, reproducible way.\nA valid Management Cluster kubeconfig must be available and set to the `E2E_KUBECONFIG` environment variable.",
		Example: "standup --provider aws --context capa",
		Args:    cobra.NoArgs,
		RunE:    run,
	}

	provider          string
	kubeContext       string
	clusterValues     string
	defaultAppValues  string
	releaseVersion    string
	releaseCommit     string
	clusterVersion    string
	defaultAppVersion string
	outputDirectory   string

	controlPlaneNodes int
	workerNodes       int

	waitForApps bool

	// Functions to run after cluster creation to confirm it is up and ready to use
	clusterReadyFns []func(wcClient *client.Client) = []func(wcClient *client.Client){
		func(wcClient *client.Client) {
			_ = wait.For(
				wait.AreNumNodesReady(context.Background(), wcClient, controlPlaneNodes, &cr.MatchingLabels{"node-role.kubernetes.io/control-plane": ""}),
				wait.WithTimeout(20*time.Minute),
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
	standupCmd.Flags().StringVar(&defaultAppValues, "default-apps-values", "", "The path to the default-apps app values")
	standupCmd.Flags().IntVar(&controlPlaneNodes, "control-plane-nodes", 1, "The number of control plane nodes to wait for being ready")
	standupCmd.Flags().IntVar(&workerNodes, "worker-nodes", 1, "The number of worker nodes to wait for being ready")
	standupCmd.Flags().StringVar(&outputDirectory, "output", "./", "The directory to store the results.json and kubeconfig in")

	standupCmd.Flags().StringVar(&releaseVersion, "release", application.ReleaseLatest, "The version of the Release to use")
	standupCmd.Flags().StringVar(&releaseCommit, "release-commit", "", "The git commit to get the Release version from (defaults to main default if unset)")
	standupCmd.Flags().StringVar(&clusterVersion, "cluster-version", "latest", "The version of the cluster app to install")
	standupCmd.Flags().StringVar(&defaultAppVersion, "default-apps-version", "latest", "The version of the default-apps app to install")

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
	defaultAppValuesOverrides := []string{values.MustLoadValuesFile(defaultAppValues)}

	var cluster *application.Cluster
	clusterBuilder, err := cb.GetClusterBuilderForContext(kubeContext)
	if err != nil {
		fmt.Printf("Failed to automatically get cluster builder based on context, falling back to building out a cluster...\n")
		cluster = application.NewClusterApp(clusterName, provider).
			WithRelease(application.ReleasePair{
				Version: releaseVersion,
				Commit:  releaseCommit,
			}).
			WithAppVersions(clusterVersion, defaultAppVersion).WithOrg(organization.New(orgName)).
			WithAppValuesFile(path.Clean(clusterValues), path.Clean(defaultAppValues), &application.TemplateValues{
				ClusterName:  clusterName,
				Organization: orgName,
			})
	} else {
		cluster = clusterBuilder.NewClusterApp(clusterName, orgName, clusterValuesOverrides, defaultAppValuesOverrides).
			WithRelease(application.ReleasePair{
				Version: releaseVersion,
				Commit:  releaseCommit,
			}).
			WithAppVersions(clusterVersion, defaultAppVersion)
	}

	if provider == application.ProviderEKS {
		// As EKS has no control plane we only check for worker nodes being ready
		clusterReadyFns = []func(wcClient *client.Client){func(wcClient *client.Client) {
			_ = wait.For(
				wait.AreNumNodesReady(context.Background(), wcClient, workerNodes, &cr.MatchingLabels{"node-role.kubernetes.io/worker": ""}),
				wait.WithTimeout(20*time.Minute),
				wait.WithInterval(15*time.Second),
			)
		},
		}
	}

	if waitForApps {
		clusterReadyFns = append(clusterReadyFns, func(wcClient *client.Client) {
			fmt.Printf("Waiting for all Apps to be ready...\n")
			_ = wait.For(func() (bool, error) {
				appList := &v1alpha1.AppList{}
				err := framework.MC().List(context.Background(), appList, cr.InNamespace(cluster.GetNamespace()), cr.MatchingLabels{"giantswarm.io/cluster": cluster.Name})
				if err != nil {
					return false, err
				}

				appNamespacedNames := []apitypes.NamespacedName{}
				for _, app := range appList.Items {
					appNamespacedNames = append(appNamespacedNames, apitypes.NamespacedName{Name: app.Name, Namespace: app.Namespace})
				}

				return wait.IsAllAppDeployed(context.Background(), framework.MC(), appNamespacedNames)()
			},
				wait.WithTimeout(20*time.Minute),
				wait.WithInterval(15*time.Second))
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

	resultsFile, err := os.Create(path.Join(outputDirectory, "results.json"))
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

	resultsFile.Close()

	fmt.Printf("Standing up cluster...\n\nProvider:\t\t%s\nCluster Name:\t\t%s\nOrg Name:\t\t%s\nResults Directory:\t%s\n\n", provider, clusterName, orgName, outputDirectory)

	cluster, err = standup.New(framework, false, clusterReadyFns...).Standup(cluster)
	if err != nil {
		return err
	}

	// Save the kubeconfig for the WC
	kubeconfigFile, err := os.Create(path.Join(outputDirectory, "kubeconfig"))
	if err != nil {
		return err
	}
	defer kubeconfigFile.Close()

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

	resultsFile, err = os.Create(path.Join(outputDirectory, "results.json"))
	if err != nil {
		return err
	}
	defer resultsFile.Close()

	resultBytes, err = json.Marshal(result)
	if err != nil {
		return err
	}
	_, err = resultsFile.Write(resultBytes)

	return err
}
