package capz

import (
	_ "embed"

	"github.com/giantswarm/cluster-standup-teardown/v3/pkg/values"

	"github.com/giantswarm/clustertest/v3/pkg/application"
	"github.com/giantswarm/clustertest/v3/pkg/organization"
	"github.com/giantswarm/clustertest/v3/pkg/utils"
)

var (
	//go:embed values/cluster_values.yaml
	baseClusterValues string
	//go:embed values/default-apps_values.yaml
	baseDefaultAppsValues string
)

// ClusterBuilder is the CAPZ ClusterBuilder
type ClusterBuilder struct {
	CustomKubeContext string
}

// NewClusterApp builds a new CAPZ cluster App
func (c *ClusterBuilder) NewClusterApp(clusterName string, orgName string, clusterValuesOverrides []string, defaultAppsValuesOverrides []string) *application.Cluster {
	if clusterName == "" {
		clusterName = utils.GenerateRandomName("t")
	}
	if orgName == "" {
		orgName = utils.GenerateRandomName("t")
	}

	return application.NewClusterApp(clusterName, application.ProviderAzure).
		WithOrg(organization.New(orgName)).
		WithAppValues(
			values.MustMergeValues(append([]string{baseClusterValues}, clusterValuesOverrides...)...),
			values.MustMergeValues(append([]string{baseDefaultAppsValues}, defaultAppsValuesOverrides...)...),
			&application.TemplateValues{
				ClusterName:  clusterName,
				Organization: orgName,
			},
		)
}

// KubeContext returns the known KubeConfig context that this builder expects
func (c *ClusterBuilder) KubeContext() string {
	if c.CustomKubeContext != "" {
		return c.CustomKubeContext
	}
	return "capz"
}
