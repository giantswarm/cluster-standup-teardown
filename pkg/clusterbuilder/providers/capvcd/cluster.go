package capvcd

import (
	_ "embed"

	"github.com/giantswarm/cluster-standup-teardown/v3/pkg/values"

	"github.com/giantswarm/clustertest/v3/pkg/application"
	"github.com/giantswarm/clustertest/v3/pkg/organization"
	"github.com/giantswarm/clustertest/v3/pkg/utils"
)

const (
	RegCredSecretName      = "container-registries-configuration"
	RegCredSecretNamespace = "default"
)

var (
	//go:embed values/cluster_values.yaml
	baseClusterValues string
)

// ClusterBuilder is the CAPVCD ClusterBuilder
type ClusterBuilder struct {
	CustomKubeContext string
}

// NewClusterApp builds a new CAPVCD cluster App
func (c *ClusterBuilder) NewClusterApp(clusterName string, orgName string, clusterValuesOverrides []string) *application.Cluster {
	if clusterName == "" {
		clusterName = utils.GenerateRandomName("t")
	}
	if orgName == "" {
		orgName = utils.GenerateRandomName("t")
	}

	return application.NewClusterApp(clusterName, application.ProviderCloudDirector).
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
	return "capvcd"
}
