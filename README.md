# cluster-standup-teardown

<a href="https://godoc.org/github.com/giantswarm/cluster-standup-teardown"><img src="https://godoc.org/github.com/giantswarm/cluster-standup-teardown?status.svg"></a>

A helper module for use in Giant Swarm E2E test frameworks to handle the creation and deletion of workload clusters to perform tests against.

> Note: This should not be used as a general purpose tool to create Giant Swarm clusters and is solely designed for use in our test frameworks.

## Installation

```shell
go get github.com/giantswarm/cluster-standup-teardown
```

### Environment Variables

#### E2E_KUBECONFIG

When using this module to standup a workload cluster it is expected that the `E2E_KUBECONFIG` environment variable is set and pointing to a valid kubeconfig with expected contexts defined.

#### E2E_OVERRIDE_VERSIONS

The `E2E_OVERRIDE_VERSIONS` environment variable allows you to override app versions when creating test clusters. This is useful for testing specific versions of apps.

**Format:** Comma-separated list of `app-name=version` pairs.

**Example:**
```bash
export E2E_OVERRIDE_VERSIONS="cluster-aws=7.2.5-164a75740365c5c21ca8aed69ebeb05f75c07fd8,karpenter=2.0.0,aws-ebs-csi-driver=4.1.0"
```

This will:
- Use the specified version for `cluster-aws` (the main cluster app)
- Override `karpenter` and `aws-ebs-csi-driver` versions in the Release CR

**Note:** When using a version with a commit SHA suffix (e.g., `7.2.5-164a75740365c5c21ca8aed69ebeb05f75c07fd8`), the app catalog will automatically be changed to `<catalog>-test` (e.g., `cluster` → `cluster-test`).

### KubeConfig Contexts

Each [ClusterBuilder](./pkg/clusterbuilder/) in this module has a specific KubeContext that it supports and that is expected to exist in the provided KubeConfig when running.

Example kubeconfig:

```yaml
apiVersion: v1
kind: Config
contexts:
- context:
    cluster: glippy
    user: glippy-admin
  name: capz
- context:
    cluster: grizzly
    user: grizzly-admin
  name: capa
- context:
    cluster: gcapeverde
    user: gcapeverde-admin
  name: capv
- context:
    cluster: gerbil
    user: gerbil-admin
  name: capvcd
clusters:
- cluster:
    certificate-authority-data: [REDACTED]
    server: https://[REDACTED]:6443
  name: glippy
- cluster:
    certificate-authority-data: [REDACTED]
    server: https://[REDACTED]:6443
  name: grizzly
- cluster:
    certificate-authority-data: [REDACTED]
    server: https://[REDACTED]:6443
  name: gcapeverde
- cluster:
    certificate-authority-data: [REDACTED]
    server: https://[REDACTED]:6443
  name: gerbil
current-context: grizzly
preferences: {}
users:
- name: glippy-admin
  user:
    client-certificate-data: [REDACTED]
    client-key-data: [REDACTED]
- name: grizzly-admin
  user:
    client-certificate-data: [REDACTED]
    client-key-data: [REDACTED]
- name: gcapeverde-admin
  user:
    client-certificate-data: [REDACTED]
    client-key-data: [REDACTED]
- name: gerbil-admin
  user:
    client-certificate-data: [REDACTED]
    client-key-data: [REDACTED]
```

## API Documentation

Documentation can be found at: [pkg.go.dev/github.com/giantswarm/cluster-standup-teardown](https://pkg.go.dev/github.com/giantswarm/cluster-standup-teardown).

## CLI Documentation

### standup

Create a workload cluster using the same known-good configuration as the E2E test suites.

Once the workload cluster is ready two files will be produced:

* a `kubeconfig` to use to access the cluster
* a `results.json` that contains details about the cluster created and can be used by `teardown` to cleanup the cluster when done

#### Install

```shell
go install github.com/giantswarm/cluster-standup-teardown/cmd/standup@latest
```

If this fails with the error `The go.mod file for the module providing named packages contains one or more replace directives.` then
you can install the binary directly from this repo:

```
cd cmd/standup
go install ./
```

#### Usage

```
$ standup --help

Standup create a test workload cluster in a standard, reproducible way.
A valid Management Cluster kubeconfig must be available and set to the `E2E_KUBECONFIG` environment variable.

Usage:
  standup [flags]

Examples:
standup --provider aws --context capa

Flags:
      --cluster-values string         The path to the cluster app values
      --cluster-version string        The version of the cluster app to install (default "latest")
      --context string                The kubernetes context to use (required)
      --control-plane-nodes int       The number of control plane nodes to wait for being ready (default 1)
  -h, --help                          help for standup
      --output string                 The directory to store the results.json and kubeconfig in (default "./")
      --provider string               The provider (required)
      --release string                The version of the Release to use (default "latest")
      --release-commit string         The git commit to get the Release version from (defaults to main default if unset)
      --wait-for-apps-ready           Wait until all default apps are installed
      --worker-nodes int              The number of worker nodes to wait for being ready (default 1)
```

##### Example

```
export E2E_KUBECONFIG=/path/to/kubeconfig
standup --provider aws --context capa
```

### teardown

Cleans up a workload cluster previously created by `standup`. Makes use of the `kubeconfig` and `results.json` produced by `standup`.

#### Preventing cluster deletion

It is possible to prevent cluster deletion during the teardown stage by setting the `E2E_WC_KEEP` environment variable to anything other than `false`. If this env var is found the teardown function will not actually perform any actions against the cluster and will instead just log out to the logger that deletion has been skipped and the user must then maunally clean up the resources in the MC.

#### Install

```shell
go install github.com/giantswarm/cluster-standup-teardown/cmd/teardown@latest
```

#### Usage

```
$ teardown --help

Teardown completely removes a previously created test cluster.
Can take in the results.json produced by `standup` to quickly clean up a test cluster.
A valid Management Cluster kubeconfig must be available and set to the `E2E_KUBECONFIG` environment variable.

Usage:
  teardown [flags]

Examples:
teardown --context capa --standup-directory ./

Flags:
      --cleanup-standup-results    Remove the results and kubeconfig generated by 'standup' (default true)
      --cluster-name string        The name of the cluster to tear down
      --context string             The kubernetes context to use (required)
  -h, --help                       help for teardown
      --org-name string            The org the cluster belongs to
      --provider string            The provider
      --standup-directory string   The directory containing the results from 'standup'
```

##### Example

```
teardown --provider aws --context capa --standup-directory ./
```
