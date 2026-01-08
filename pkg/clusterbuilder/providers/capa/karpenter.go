package capa

import (
	_ "embed"

	"github.com/giantswarm/cluster-standup-teardown/v2/pkg/values"

	"github.com/giantswarm/clustertest/v3/pkg/application"
	"github.com/giantswarm/clustertest/v3/pkg/organization"
	"github.com/giantswarm/clustertest/v3/pkg/utils"
)

var (
	//go:embed values/karpenter-cluster_values.yaml
	baseKarpenterClusterValues string
)

// KarpenterBuilder is the CAPA KarpenterBuilder
type KarpenterBuilder struct {
	CustomKubeContext string
}

// NewClusterApp builds a new CAPA cluster App
func (c *KarpenterBuilder) NewClusterApp(clusterName string, orgName string, clusterValuesOverrides []string, defaultAppsValuesOverrides []string) *application.Cluster {
	if clusterName == "" {
		clusterName = utils.GenerateRandomName("t")
	}
	if orgName == "" {
		orgName = utils.GenerateRandomName("t")
	}

	return application.NewClusterApp(clusterName, application.ProviderAWS).
		WithOrg(organization.New(orgName)).
		WithAppValues(
			values.MustMergeValues(append([]string{baseKarpenterClusterValues}, clusterValuesOverrides...)...),
			"",
			&application.TemplateValues{
				ClusterName:  clusterName,
				Organization: orgName,
			},
		)
}

// KubeContext returns the known KubeConfig context that this builder expects
func (c *KarpenterBuilder) KubeContext() string {

	if c.CustomKubeContext != "" {
		return c.CustomKubeContext
	}

	return "capa-karpenter"
}
