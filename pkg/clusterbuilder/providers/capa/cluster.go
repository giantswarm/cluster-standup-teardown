package capa

import (
	_ "embed"
	"fmt"
	"math/rand"

	"github.com/giantswarm/cluster-standup-teardown/pkg/values"

	"github.com/giantswarm/clustertest/pkg/application"
	"github.com/giantswarm/clustertest/pkg/organization"
	"github.com/giantswarm/clustertest/pkg/utils"
)

var (
	//go:embed values/cluster_values.yaml
	baseClusterValues string
	//go:embed values/default-apps_values.yaml
	baseDefaultAppsValues string
)

// ClusterBuilder is the CAPA ClusterBuilder
type ClusterBuilder struct{}

// NewClusterApp builds a new CAPA cluster App
func (c *ClusterBuilder) NewClusterApp(clusterName string, orgName string, clusterValuesFile string, defaultAppsValuesFile string) *application.Cluster {
	if clusterName == "" {
		clusterName = utils.GenerateRandomName("t")
	}
	if orgName == "" {
		orgName = utils.GenerateRandomName("t")
	}

	return application.NewClusterApp(clusterName, application.ProviderAWS).
		WithOrg(organization.New(orgName)).
		WithAppValues(
			values.MustOverlayValues(baseClusterValues, clusterValuesFile),
			values.MustOverlayValues(baseDefaultAppsValues, defaultAppsValuesFile),
			&application.TemplateValues{
				ClusterName:  clusterName,
				Organization: orgName,
			},
		)
}

var (
	//go:embed values/private-cluster_values.yaml
	basePrivateClusterValues string
	//go:embed values/private-default-apps_values.yaml
	basePrivateDefaultAppsValues string
)

// PrivateClusterBuilder is the private CAPA ClusterBuilder
type PrivateClusterBuilder struct{}

// NewClusterApp builds a new private CAPA cluster App
func (c *PrivateClusterBuilder) NewClusterApp(clusterName string, orgName string, clusterValuesFile string, defaultAppsValuesFile string) *application.Cluster {
	if clusterName == "" {
		clusterName = utils.GenerateRandomName("t")
	}
	if orgName == "" {
		orgName = utils.GenerateRandomName("t")
	}

	// WC CIDRs have to not overlap and be in the 10.225. - 10.255. range, so
	// we select a random number in that range and set it as the second octet.
	randomOctet := rand.Intn(30) + 225
	cidrOctet := fmt.Sprintf("%d", randomOctet)
	templateValues := &application.TemplateValues{
		ClusterName:  clusterName,
		Organization: orgName,
		ExtraValues: map[string]string{
			"CIDRSecondOctet": cidrOctet,
		},
	}

	return application.NewClusterApp(clusterName, application.ProviderAWS).
		WithOrg(organization.New(orgName)).
		WithAppValues(
			values.MustOverlayValues(basePrivateClusterValues, clusterValuesFile),
			values.MustOverlayValues(basePrivateDefaultAppsValues, defaultAppsValuesFile),
			templateValues,
		)
}
