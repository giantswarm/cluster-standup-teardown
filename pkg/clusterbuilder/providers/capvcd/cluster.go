package capvcd

import (
	_ "embed"

	"github.com/giantswarm/cluster-standup-teardown/pkg/values"

	applicationv1alpha1 "github.com/giantswarm/apiextensions-application/api/v1alpha1"
	"github.com/giantswarm/clustertest/pkg/application"
	"github.com/giantswarm/clustertest/pkg/organization"
	"github.com/giantswarm/clustertest/pkg/utils"
)

const (
	RegCredSecretName      = "container-registries-configuration"
	RegCredSecretNamespace = "default"
)

var (
	//go:embed values/cluster_values.yaml
	baseClusterValues string
	//go:embed values/default-apps_values.yaml
	baseDefaultAppsValues string
)

// ClusterBuilder is the CAPVCD ClusterBuilder
type ClusterBuilder struct{}

// NewClusterApp builds a new CAPVCD cluster App
func (c *ClusterBuilder) NewClusterApp(clusterName string, orgName string, clusterValuesOverrides []string, defaultAppsValuesOverrides []string) *application.Cluster {
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
			values.MustMergeValues(append([]string{baseDefaultAppsValues}, defaultAppsValuesOverrides...)...),
			&application.TemplateValues{
				ClusterName:  clusterName,
				Organization: orgName,
			},
		).
		WithExtraConfigs([]applicationv1alpha1.AppExtraConfig{
			{
				Kind:      "secret",
				Name:      RegCredSecretName,
				Namespace: RegCredSecretNamespace,
				Priority:  25,
			},
		})
}
