package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"time"

	. "github.com/onsi/gomega"

	"github.com/giantswarm/clustertest"
	"github.com/giantswarm/clustertest/pkg/application"
	"github.com/giantswarm/clustertest/pkg/client"
	"github.com/giantswarm/clustertest/pkg/organization"
	"github.com/giantswarm/clustertest/pkg/utils"
	"github.com/giantswarm/clustertest/pkg/wait"
	"github.com/spf13/cobra"
	cr "sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/giantswarm/cluster-standup-teardown/cmd/standup/types"
	"github.com/giantswarm/cluster-standup-teardown/pkg/clusterbuilder/providers/capa"
	"github.com/giantswarm/cluster-standup-teardown/pkg/clusterbuilder/providers/capv"
	"github.com/giantswarm/cluster-standup-teardown/pkg/clusterbuilder/providers/capvcd"
	"github.com/giantswarm/cluster-standup-teardown/pkg/clusterbuilder/providers/capz"
	"github.com/giantswarm/cluster-standup-teardown/pkg/clusterbuilder/providers/eks"
	"github.com/giantswarm/cluster-standup-teardown/pkg/standup"
)

var (
	standupCmd = &cobra.Command{
		Use:     "standup",
		Long:    "Standup create a test workload cluster in a standard, reproducible way.\nA valid Management Cluster kubeconfig must be available and set to the `E2E_KUBECONFIG` environment variable.",
		Example: "standup --provider aws --context capa --cluster-values ./cluster_values.yaml --default-apps-values ./default-apps_values.yaml",
		Args:    cobra.NoArgs,
		RunE:    run,
	}

	provider          string
	kubeContext       string
	clusterValues     string
	defaultAppValues  string
	clusterVersion    string
	defaultAppVersion string
	outputDirectory   string

	controlPlaneNodes int
	workerNodes       int

	// Functions to run after cluster creation to confirm it is up and ready to use
	clusterReadyFns []func(wcClient *client.Client) = []func(wcClient *client.Client){
		func(wcClient *client.Client) {
			wait.For(
				wait.AreNumNodesReady(context.Background(), wcClient, controlPlaneNodes, &cr.MatchingLabels{"node-role.kubernetes.io/control-plane": ""}),
				wait.WithTimeout(20*time.Minute),
				wait.WithInterval(15*time.Second),
			)
		},
		func(wcClient *client.Client) {
			wait.For(
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
	standupCmd.Flags().StringVar(&clusterVersion, "cluster-version", "latest", "The version of the cluster app to install")
	standupCmd.Flags().StringVar(&defaultAppVersion, "default-apps-version", "latest", "The version of the default-apps app to install")

	standupCmd.MarkFlagRequired("provider")
	standupCmd.MarkFlagRequired("context")
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

	fmt.Printf("Standing up cluster...\n\nProvider:\t\t%s\nCluster Name:\t\t%s\nOrg Name:\t\t%s\nResults Directory:\t%s\n\n", provider, clusterName, orgName, outputDirectory)

	var cluster *application.Cluster
	switch provider {
	case application.ProviderVSphere:
		clusterBuilder := capv.ClusterBuilder{}
		cluster = clusterBuilder.NewClusterApp(clusterName, orgName, clusterValues, defaultAppValues).
			WithAppVersions(clusterVersion, defaultAppVersion)
	case application.ProviderCloudDirector:
		clusterBuilder := capvcd.ClusterBuilder{}
		cluster = clusterBuilder.NewClusterApp(clusterName, orgName, clusterValues, defaultAppValues).
			WithAppVersions(clusterVersion, defaultAppVersion)
	case application.ProviderAWS:
		clusterBuilder := capa.ClusterBuilder{}
		cluster = clusterBuilder.NewClusterApp(clusterName, orgName, clusterValues, defaultAppValues).
			WithAppVersions(clusterVersion, defaultAppVersion)
	case application.ProviderEKS:
		clusterBuilder := eks.ClusterBuilder{}
		cluster = clusterBuilder.NewClusterApp(clusterName, orgName, clusterValues, defaultAppValues).
			WithAppVersions(clusterVersion, defaultAppVersion)
		// As EKS has no control plane we only check for worker nodes being ready
		clusterReadyFns = []func(wcClient *client.Client){
			func(wcClient *client.Client) {
				wait.For(
					wait.AreNumNodesReady(context.Background(), wcClient, workerNodes, &cr.MatchingLabels{"node-role.kubernetes.io/worker": ""}),
					wait.WithTimeout(20*time.Minute),
					wait.WithInterval(15*time.Second),
				)
			},
		}
	case application.ProviderAzure:
		clusterBuilder := capz.ClusterBuilder{}
		cluster = clusterBuilder.NewClusterApp(clusterName, orgName, clusterValues, defaultAppValues).
			WithAppVersions(clusterVersion, defaultAppVersion)
	default:
		cluster = application.NewClusterApp(clusterName, provider).
			WithAppVersions(clusterVersion, defaultAppVersion).
			WithOrg(organization.New(orgName)).
			WithAppValuesFile(path.Clean(clusterValues), path.Clean(defaultAppValues), &application.TemplateValues{
				ClusterName:  clusterName,
				Organization: orgName,
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

	resultsFile.Close()

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
