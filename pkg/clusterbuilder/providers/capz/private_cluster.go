package capz

import (
	_ "embed"

	"github.com/giantswarm/clustertest/v3/pkg/application"
	"github.com/giantswarm/clustertest/v3/pkg/organization"
	"github.com/giantswarm/clustertest/v3/pkg/utils"

	"github.com/giantswarm/cluster-standup-teardown/v3/pkg/values"
)

var (
	//go:embed values/private-cluster_values.yaml
	basePrivateClusterValues string
)

// PrivateClusterBuilder is the private CAPZ ClusterBuilder
type PrivateClusterBuilder struct {
	CustomKubeContext string
}

// NewClusterApp builds a new private CAPZ cluster App
func (c *PrivateClusterBuilder) NewClusterApp(clusterName string, orgName string, clusterValuesOverrides []string) *application.Cluster {
	if clusterName == "" {
		clusterName = utils.GenerateRandomName("t")
	}
	if orgName == "" {
		orgName = utils.GenerateRandomName("t")
	}

	return application.NewClusterApp(clusterName, application.ProviderAzure).
		WithOrg(organization.New(orgName)).
		WithAppValues(
			values.MustMergeValues(append([]string{basePrivateClusterValues}, clusterValuesOverrides...)...),
			&application.TemplateValues{
				ClusterName:  clusterName,
				Organization: orgName,
			},
		)
}

// KubeContext returns the known KubeConfig context that this builder expects
func (c *PrivateClusterBuilder) KubeContext() string {
	if c.CustomKubeContext != "" {
		return c.CustomKubeContext
	}
	return "capz-private"
}
