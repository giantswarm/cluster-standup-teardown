package clusterbuilder

import (
	"fmt"
	"os"
	"strings"

	. "github.com/onsi/gomega" // nolint:staticcheck

	"github.com/giantswarm/clustertest/v3"
	"github.com/giantswarm/clustertest/v3/pkg/application"
	"github.com/giantswarm/clustertest/v3/pkg/env"
	"github.com/giantswarm/clustertest/v3/pkg/logger"

	"github.com/giantswarm/cluster-standup-teardown/v4/pkg/clusterbuilder/providers/capa"
	"github.com/giantswarm/cluster-standup-teardown/v4/pkg/clusterbuilder/providers/capv"
	"github.com/giantswarm/cluster-standup-teardown/v4/pkg/clusterbuilder/providers/capvcd"
	"github.com/giantswarm/cluster-standup-teardown/v4/pkg/clusterbuilder/providers/capz"
	"github.com/giantswarm/cluster-standup-teardown/v4/pkg/values"
)

// ClusterBuilder is an interface that provides a function for building provider-specific Cluster apps
type ClusterBuilder interface {
	NewClusterApp(clusterName string, orgName string, clusterValuesOverrides []string) *application.Cluster
	KubeContext() string
}

// LoadOrBuildCluster attempts to load a pre-built workload cluster if the appropriate env vars are set and if not will build a new Cluster
// For now, when building a cluster it is assumed that the values for the cluster can be found at:
// ./test_data/cluster_values.yaml
func LoadOrBuildCluster(framework *clustertest.Framework, clusterBuilder ClusterBuilder) *application.Cluster {
	// If env vars are set, load pre-built WC
	cluster, err := framework.LoadCluster()
	Expect(err).NotTo(HaveOccurred())
	if cluster != nil {
		logger.Log("Using existing cluster %s/%s", cluster.Name, cluster.GetNamespace())
		return cluster
	}

	cluster = clusterBuilder.NewClusterApp(
		"", "",
		[]string{values.MustLoadValuesFile("./test_data/cluster_values.yaml")},
	)

	// Apply app overrides from E2E_OVERRIDE_VERSIONS for Release CR apps
	cluster = ApplyAppOverridesFromEnv(cluster)

	return cluster
}

// ApplyAppOverridesFromEnv reads the E2E_OVERRIDE_VERSIONS environment variable and applies
// overrides for apps that are part of the Release CR.
//
// The E2E_OVERRIDE_VERSIONS env var is a comma-separated list of app=version pairs.
// For example: "cluster-aws=1.2.3,karpenter=2.0.0,aws-ebs-csi-driver=4.1.0"
//
// This function enables testing specific versions of bundled apps (like karpenter,
// aws-ebs-csi-driver, etc.) that are defined in the Release CR.
//
// Note: The cluster app (e.g., cluster-aws) is also processed here, but WithAppOverride()
// will silently ignore it since it's not a "default app" in the Release CR (it's a component).
// The cluster app version is handled separately via WithAppVersions() in the application package.
func ApplyAppOverridesFromEnv(cluster *application.Cluster) *application.Cluster {
	overrides := os.Getenv(env.OverrideVersions)
	if overrides == "" {
		return cluster
	}

	for _, pair := range strings.Split(overrides, ",") {
		parts := strings.Split(pair, "=")
		if len(parts) != 2 {
			continue
		}
		appName := strings.TrimSpace(parts[0])
		version := strings.TrimSpace(parts[1])

		if appName == "" || version == "" {
			continue
		}

		logger.Log("Attempting Release app override from E2E_OVERRIDE_VERSIONS: %s=%s", appName, version)
		app := *application.New(appName, appName).WithVersion(version)
		cluster = cluster.WithAppOverride(app)
	}

	return cluster
}

// GetClusterBuilderForContext returns a suitable ClusterBuilder instance that supports the provided KubeContext
func GetClusterBuilderForContext(context string) (ClusterBuilder, error) {
	knownBuilders := []ClusterBuilder{
		&capa.ClusterBuilder{}, &capa.ManagedClusterBuilder{}, &capa.PrivateClusterBuilder{},
		&capv.ClusterBuilder{}, &capz.PrivateClusterBuilder{},
		&capvcd.ClusterBuilder{},
		&capz.ClusterBuilder{},
	}

	for _, builder := range knownBuilders {
		if strings.EqualFold(builder.KubeContext(), strings.ToLower(context)) {
			return builder, nil
		}
	}

	return nil, fmt.Errorf("unable to find matching ClusterBuilder")
}
