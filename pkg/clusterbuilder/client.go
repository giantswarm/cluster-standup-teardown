package clusterbuilder

import (
	"fmt"
	"os"
	"strings"

	. "github.com/onsi/gomega" // nolint:staticcheck

	"github.com/giantswarm/clustertest/v4"
	"github.com/giantswarm/clustertest/v4/pkg/application"
	"github.com/giantswarm/clustertest/v4/pkg/env"
	"github.com/giantswarm/clustertest/v4/pkg/logger"

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
// The E2E_OVERRIDE_VERSIONS env var is a comma-separated list of app=version[:catalog] pairs.
// The catalog is optional; if not specified, the default catalog from clustertest is used.
//
// Examples:
//   - "cluster-aws=1.2.3" - uses default catalog
//   - "karpenter=2.0.0:giantswarm" - uses 'giantswarm' catalog
//   - "cluster-aws=1.2.3,karpenter=2.0.0:giantswarm,aws-ebs-csi-driver=4.1.0:default"
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

	// Save the cluster app state before processing overrides.
	// WithAppOverride() calls IsDefaultApp() which calls GetRelease() which
	// calls c.ClusterApp.Build(). Build() has side effects: it resolves the
	// cluster app version from E2E_OVERRIDE_VERSIONS and potentially changes
	// the catalog to "<catalog>-test" when the version has a SHA suffix.
	//
	// This side effect is problematic because:
	// 1. It prematurely resolves the cluster app version and catalog.
	// 2. The catalog change is one-way (WithVersion only adds "-test", never
	//    removes it), so if the version is later changed to a stable one
	//    (e.g., via WithAppVersions("latest") in upgrade tests), the catalog
	//    stays as "cluster-test" causing Helm to look for the stable version
	//    in the test catalog where it doesn't exist.
	//
	// We save and restore the state to undo these side effects.
	savedVersion := cluster.ClusterApp.Version
	savedCatalog := cluster.ClusterApp.Catalog

	for _, pair := range strings.Split(overrides, ",") {
		parts := strings.Split(pair, "=")
		if len(parts) != 2 {
			continue
		}
		appName := strings.TrimSpace(parts[0])
		versionAndCatalog := strings.TrimSpace(parts[1])

		if appName == "" || versionAndCatalog == "" {
			continue
		}

		// Parse version and optional catalog (format: version or version:catalog)
		version, catalog := parseVersionAndCatalog(versionAndCatalog)

		if version == "" {
			continue
		}

		app := application.New(appName, appName).WithVersion(version)
		if catalog != "" {
			app = app.WithCatalog(catalog)
			logger.Log("Attempting Release app override from E2E_OVERRIDE_VERSIONS: %s=%s (catalog: %s)", appName, version, catalog)
		} else {
			logger.Log("Attempting Release app override from E2E_OVERRIDE_VERSIONS: %s=%s", appName, version)
		}
		cluster = cluster.WithAppOverride(*app)
	}

	// Restore the cluster app state to undo the Build() side effects.
	cluster.ClusterApp.Version = savedVersion
	cluster.ClusterApp.Catalog = savedCatalog

	return cluster
}

// parseVersionAndCatalog splits a version string that may contain an optional catalog suffix.
// Format: "version" or "version:catalog"
// Returns the version and catalog (empty string if no catalog specified).
func parseVersionAndCatalog(versionAndCatalog string) (version, catalog string) {
	// Find the last colon to split version and catalog
	// We use LastIndex because version strings can contain colons in edge cases,
	// but catalog names should be simple identifiers
	lastColonIdx := strings.LastIndex(versionAndCatalog, ":")
	if lastColonIdx == -1 {
		// No catalog specified
		return versionAndCatalog, ""
	}

	version = strings.TrimSpace(versionAndCatalog[:lastColonIdx])
	catalog = strings.TrimSpace(versionAndCatalog[lastColonIdx+1:])

	// If catalog is empty after the colon, treat as no catalog
	if catalog == "" {
		return versionAndCatalog, ""
	}

	return version, catalog
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
