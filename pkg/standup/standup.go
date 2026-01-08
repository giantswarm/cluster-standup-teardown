package standup

import (
	"context"
	"os"
	"strings"
	"time"

	. "github.com/onsi/gomega" // nolint:staticcheck

	"github.com/giantswarm/clustertest/v3"
	"github.com/giantswarm/clustertest/v3/pkg/application"
	clustertestclient "github.com/giantswarm/clustertest/v3/pkg/client"
	"github.com/giantswarm/clustertest/v3/pkg/env"
	"github.com/giantswarm/clustertest/v3/pkg/logger"
	"github.com/giantswarm/clustertest/v3/pkg/utils"
	"github.com/giantswarm/clustertest/v3/pkg/wait"
	cr "sigs.k8s.io/controller-runtime/pkg/client"
)

// Client is the client responsible for handling standing up a given cluster
type Client struct {
	Framework    *clustertest.Framework
	IsUpgade     bool
	ClusterReady []func(client *clustertestclient.Client)
}

// New builds a new standup client
func New(framework *clustertest.Framework, isUpgrade bool, clusterReadyChecks ...func(client *clustertestclient.Client)) *Client {
	return &Client{
		Framework:    framework,
		IsUpgade:     isUpgrade,
		ClusterReady: clusterReadyChecks,
	}
}

// Standup takes in a Cluster app and applies it to the Management Cluster.
// After applying it checks for the cluster being ready and usable.
func (c *Client) Standup(cluster *application.Cluster) (*application.Cluster, error) {
	if c.IsUpgade {
		Expect(utils.ShouldSkipUpgrade()).To(BeFalse())

		releaseVersion := application.ReleaseLatest
		if os.Getenv(env.ReleasePreUpgradeVersion) != "" {
			releaseVersion = os.Getenv(env.ReleasePreUpgradeVersion)
		}

		cluster = cluster.
			WithAppVersions("latest").
			WithRelease(application.ReleasePair{
				Version: releaseVersion,
				Commit:  "",
			})
	}
	logger.Log("Workload cluster name: %s", cluster.Name)
	logger.Log("Organisation name: %s", cluster.Organization.Name)

	// In certain cases, when connecting over the VPN, it is possible that the tunnel
	// isn't ready and can take a short while to become usable. This attempts to wait
	// for the connection to be usable before starting the tests.
	Eventually(func() error {
		logger.Log("Checking connection to MC is available.")
		logger.Log("MC API Endpoint: '%s'", c.Framework.MC().GetAPIServerEndpoint())
		mcName := c.Framework.MC().GetClusterName()
		// Strip the Teleport prefix that is automatically added to the cluster name so we just log out the installation name
		mcName = strings.TrimPrefix(mcName, "teleport.giantswarm.io-")
		logger.Log("MC name: '%s'", mcName)
		return c.Framework.MC().CheckConnection()
	}).
		WithTimeout(5 * time.Minute).
		WithPolling(5 * time.Second).
		Should(Succeed())

	ctx := context.Background()
	applyCtx, cancelApplyCtx := context.WithTimeout(ctx, 25*time.Minute)
	defer cancelApplyCtx()

	client, err := c.Framework.ApplyCluster(applyCtx, cluster)
	Expect(err).NotTo(HaveOccurred())

	if len(c.ClusterReady) > 0 {
		// Use provided functions to check if cluster is ready.
		// This is mainly used for managed clusters such as EKS that need to check for worker nodes rather than control plane nodes.
		for _, fn := range c.ClusterReady {
			fn(client)
		}
	} else {
		// If no custom check functions are provided we default to checking for a single control plane node being ready
		Eventually(
			wait.AreNumNodesReady(ctx, client, 1, &cr.MatchingLabels{"node-role.kubernetes.io/control-plane": ""}),
			20*time.Minute, 15*time.Second,
		).Should(BeTrue())
	}

	return cluster, nil
}
