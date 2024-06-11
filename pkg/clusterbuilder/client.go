package clusterbuilder

import (
	"fmt"
	"strings"

	. "github.com/onsi/gomega"

	"github.com/giantswarm/clustertest"
	"github.com/giantswarm/clustertest/pkg/application"
	"github.com/giantswarm/clustertest/pkg/logger"

	"github.com/giantswarm/cluster-standup-teardown/pkg/clusterbuilder/providers/capa"
	"github.com/giantswarm/cluster-standup-teardown/pkg/clusterbuilder/providers/capv"
	"github.com/giantswarm/cluster-standup-teardown/pkg/clusterbuilder/providers/capvcd"
	"github.com/giantswarm/cluster-standup-teardown/pkg/clusterbuilder/providers/capz"
	"github.com/giantswarm/cluster-standup-teardown/pkg/values"
)

// ClusterBuilder is an interface that provides a function for building provider-specific Cluster apps
type ClusterBuilder interface {
	NewClusterApp(clusterName string, orgName string, clusterValuesOverrides []string, defaultAppsValuesOverrides []string) *application.Cluster
	KubeContext() string
}

// LoadOrBuildCluster attempts to load a pre-built workload cluster if the appropriate env vars are set and if not will build a new Cluster
// For now, when building a cluster it is assumed that the values for the cluster and default-apps can be found at:
// ./test_data/cluster_values.yaml and ./test_data/default-apps_values.yaml
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
		[]string{values.MustLoadValuesFile("./test_data/default-apps_values.yaml")},
	)

	return cluster
}

// GetClusterBuilderForContext returns a suitable ClusterBuilder instance that supports the provided KubeContext
func GetClusterBuilderForContext(context string) (ClusterBuilder, error) {
	knownBuilders := []ClusterBuilder{
		&capa.ClusterBuilder{}, &capa.ManagedClusterBuilder{}, &capa.PrivateClusterBuilder{},
		&capv.ClusterBuilder{},
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
