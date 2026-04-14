package capmox

import (
	_ "embed"

	"github.com/giantswarm/cluster-standup-teardown/v5/pkg/values"

	"github.com/giantswarm/clustertest/v4/pkg/application"
	"github.com/giantswarm/clustertest/v4/pkg/organization"
	"github.com/giantswarm/clustertest/v4/pkg/utils"
)

var (
	//go:embed values/cluster_values.yaml
	baseClusterValues string
)

// ClusterBuilder is the CAPMOX ClusterBuilder
type ClusterBuilder struct {
	CustomKubeContext string
}

// NewClusterApp builds a new CAPMOX cluster App
func (c *ClusterBuilder) NewClusterApp(clusterName string, orgName string, clusterValuesOverrides []string) *application.Cluster {
	if clusterName == "" {
		clusterName = utils.GenerateRandomName("t")
	}
	if orgName == "" {
		orgName = utils.GenerateRandomName("t")
	}

	return application.NewClusterApp(clusterName, application.ProviderProxmox).
		WithOrg(organization.New(orgName)).
		WithAppValues(
			values.MustMergeValues(append([]string{baseClusterValues}, clusterValuesOverrides...)...),
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
	return "capmox"
}
