# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Changed

- CAPA: use lower heartbeat timeout to allow spot instances to terminate more quickly

## [1.27.3] - 2024-11-13

### Changed

- Update CAPV values to deploy on new Neoedge provider

## [1.27.2] - 2024-11-07

## Fixed

- Updated `clustertest` to fix release name when having provider with a hyphen.

## [1.27.1] - 2024-11-07

## Fixed

- Updated `clustertest` to fix Cloud Director.

## [1.27.0] - 2024-11-07

## Changed

- Updated `clustertest` to include latest supported Release providers

## [1.26.1] - 2024-10-31

### Changed

- Drop deprecated values from `cluster-vsphere` tests

## [1.26.0] - 2024-10-21

### Changed

- Updated `clustertest` to support `cluster-cloud-director` as a unified app

## [1.25.8] - 2024-10-15

### Fixed

- Updated `clustertest` to include latest supported Release providers

## [1.25.7] - 2024-10-11

- Update CAPVCD values for provider version v0.61.0

## [1.25.6] - 2024-10-08

### Changed

- Change number of reserved IPs for LB in vSphere pool to 1.

## [1.25.5] - 2024-10-08

### Fixed

- Update `clustertest` with Provider fix when loading existing workload cluster

## [1.25.4] - 2024-10-08

### Changed

- Update CAPVCD values with proxy vars

## [1.25.3] - 2024-09-25

### Changed

- Update CAPVCD values with `name` value

## [1.25.2] - 2024-09-24

### Changed

- Update capvcd values to not include node classes.

## [1.25.1] - 2024-09-19

### Fixed

- Update `clustertest` with GitHub latest release fix

## [1.25.0] - 2024-09-16

### Changed

- Updated Go version to v1.23.1

## [1.24.0] - 2024-09-06

### Changed

- Upgraded `clustertest` to v1.23.0 and transient deps to latest

## [1.23.1] - 2024-09-02

### Fixed

- Updated clustertest with fix for version prefix on releases

## [1.23.0] - 2024-08-30

### Changed

- Generate base labels from `clustertest`

## [1.22.0] - 2024-08-27

### Added

- Updated `clustertest` with support for unified cluster-vsphere app.

## [1.21.0] - 2024-08-21

### Added

- Support for private CAPZ cluster

## [1.20.0] - 2024-08-19

### Changed

- Updated `clustertest` to v1.19.0 to make use of Teleport kubeconfig if available

## [1.19.0] - 2024-08-15

### Fixed

- Replace `containerdVolumeSizeGB` and `kubeletVolumeSizeGB` with `libVolumeSizeGB`.

### Changed

- Updated all modules to latest (including support for Kubernetes v1.31)

## [1.18.1] - 2024-08-08

### Fixed

- Upgraded `clustertest` to include fix for correctly handling Releases version prefixes.

## [1.18.0] - 2024-08-06

### Changed

- Updated `clustertest` to latest v1.17.0

## [1.17.2] - 2024-07-26

### Fixed

- Bump clustertest with latest releases SDK to correctly get latest release

## [1.17.1] - 2024-07-23

### Fixed

- Bump releases SDK to actually handle Azure

## [1.17.0] - 2024-07-23

### Changed

- Bumped clustertest to 0.16.0 with Releases support for CAPZ

## [1.16.0] - 2024-07-22

### Added

- Added a new `--wait-for-apps-ready` flag to the standup CLI that will wait until all default apps are installed

## [1.15.0] - 2024-07-22

### Changed

- Update CAPV values with `name` value

## [1.14.0] - 2024-07-11

### Added

- Added support for specifying the Release version when using the `standup` CLI

## [1.13.0] - 2024-07-09

### Changed

- Updated `capv` cluster values following removal of nodeClasses.

## [1.12.1] - 2024-07-05

### Fixed

- Use `ShouldSkipUpgrade()` to check if upgrade is valid or not

## [1.12.0] - 2024-07-04

### Fixed

- Use `E2E_WC_KEEP` from `clustertest` and added a small note when using `E2E_WC_NAME` but not `E2E_WC_KEEP`.

## [1.11.0] - 2024-07-02

### Changed

- Updated clustertest to v1.10.0 and changed the `IsUpgrade` support in Standup to handle setting the Release version used by the cluster before upgrade.

## [1.10.0] - 2024-06-25

### Changed

- Updated all depenedncies to latest version

## [1.9.0] - 2024-06-24

### Changed

- Updated all depenedncies to latest version

## [1.8.0] - 2024-06-20

### Changed

- Added `scaleDownUnneededTime` of 15m to `cluster-autoscaler` config for CAPA clusters to help avoid scale-down occurring during our scaling tests.
- Updated `clustertest` to latest with additional logging

## [1.7.1] - 2024-06-13

### Fixed

- Upgraded `clustertest` to latest v1.1.0

## [1.7.0] - 2024-06-12

### Added

- Add CAPA China cluster builder.

## [1.6.0] - 2024-06-10

### Changed

- Update `clustertest` to v1.0.0 to support Releases with cluster Apps

## [1.5.0] - 2024-06-07

### Changed

- Update `cluster-cloud-director` values for refactored chart.

## [1.4.0] - 2024-05-20

### Added

- If relevant env vars are found populate the Cluster values with labels containing the Tekton run names

### Changed

- Added a `values.BuildBaseValues` function to handle generic cluster values that are specific to test environments and apply over all providers. The prometheus volume size label has been moved into this function.

## [1.3.0] - 2024-05-16

### Changed

- Update `cluster-vsphere` values for refactored chart.

## [1.2.0] - 2024-05-14

### Changed

- Reduce prom volume size in test clusters

### Added

- Added support for the `E2E_WC_KEEP` environment variable that prevent cluster deletion during the teardown phase
- Add support for unified cluster-aws app.

## [1.1.0] - 2024-05-12

### Changed

- Update clustertest to v0.20.0

## [1.0.3] - 2024-05-07

### Added

- Added fields to ClusterBuilder structs to allow consumers to use custom kubeconfig contexts.

## [1.0.2] - 2024-04-29

### Added

- Each ClusterBuilder now includes a function to say what KubeContext it expects / supports
- Added a function to get a ClusterBuilder based on a given KubeContext

## [1.0.1] - 2024-04-26

### Changed

- Implemented value file merging / overlaying. Values provided to the clusterbuilder will be merged ontop of the default ones included in this module.
- `clusterbuilder` updated to take in a slice of values overrides that are layered ontop of the default values
- Updated `LoadOrBuildCluster` and `standup` to work with the `clusterbuilder` refactoring.

## [1.0.0] - 2024-04-26

### Added

- ClusterBuilder for CAPA, CAPV, CAPVCD, CAPZ and EKS along with default values for each
- standup and teardown modules for use in other projects
- `standup` and `teardown` CLIs
- Dockerfile containing the two CLIs

[Unreleased]: https://github.com/giantswarm/cluster-standup-teardown/compare/v1.27.3...HEAD
[1.27.3]: https://github.com/giantswarm/cluster-standup-teardown/compare/v1.27.2...v1.27.3
[1.27.2]: https://github.com/giantswarm/cluster-standup-teardown/compare/v1.27.1...v1.27.2
[1.27.1]: https://github.com/giantswarm/cluster-standup-teardown/compare/v1.27.0...v1.27.1
[1.27.0]: https://github.com/giantswarm/cluster-standup-teardown/compare/v1.26.1...v1.27.0
[1.26.1]: https://github.com/giantswarm/cluster-standup-teardown/compare/v1.26.0...v1.26.1
[1.26.0]: https://github.com/giantswarm/cluster-standup-teardown/compare/v1.25.8...v1.26.0
[1.25.8]: https://github.com/giantswarm/cluster-standup-teardown/compare/v1.25.7...v1.25.8
[1.25.7]: https://github.com/giantswarm/cluster-standup-teardown/compare/v1.25.6...v1.25.7
[1.25.6]: https://github.com/giantswarm/cluster-standup-teardown/compare/v1.25.5...v1.25.6
[1.25.5]: https://github.com/giantswarm/cluster-standup-teardown/compare/v1.25.4...v1.25.5
[1.25.4]: https://github.com/giantswarm/cluster-standup-teardown/compare/v1.25.3...v1.25.4
[1.25.3]: https://github.com/giantswarm/cluster-standup-teardown/compare/v1.25.2...v1.25.3
[1.25.2]: https://github.com/giantswarm/cluster-standup-teardown/compare/v1.25.1...v1.25.2
[1.25.1]: https://github.com/giantswarm/cluster-standup-teardown/compare/v1.25.0...v1.25.1
[1.25.0]: https://github.com/giantswarm/cluster-standup-teardown/compare/v1.24.0...v1.25.0
[1.24.0]: https://github.com/giantswarm/cluster-standup-teardown/compare/v1.23.1...v1.24.0
[1.23.1]: https://github.com/giantswarm/cluster-standup-teardown/compare/v1.23.0...v1.23.1
[1.23.0]: https://github.com/giantswarm/cluster-standup-teardown/compare/v1.22.0...v1.23.0
[1.22.0]: https://github.com/giantswarm/cluster-standup-teardown/compare/v1.21.0...v1.22.0
[1.21.0]: https://github.com/giantswarm/cluster-standup-teardown/compare/v1.20.0...v1.21.0
[1.20.0]: https://github.com/giantswarm/cluster-standup-teardown/compare/v1.19.0...v1.20.0
[1.19.0]: https://github.com/giantswarm/cluster-standup-teardown/compare/v1.18.1...v1.19.0
[1.18.1]: https://github.com/giantswarm/cluster-standup-teardown/compare/v1.18.0...v1.18.1
[1.18.0]: https://github.com/giantswarm/cluster-standup-teardown/compare/v1.17.2...v1.18.0
[1.17.2]: https://github.com/giantswarm/cluster-standup-teardown/compare/v1.17.1...v1.17.2
[1.17.1]: https://github.com/giantswarm/cluster-standup-teardown/compare/v1.17.0...v1.17.1
[1.17.0]: https://github.com/giantswarm/cluster-standup-teardown/compare/v1.16.0...v1.17.0
[1.16.0]: https://github.com/giantswarm/cluster-standup-teardown/compare/v1.15.0...v1.16.0
[1.15.0]: https://github.com/giantswarm/cluster-standup-teardown/compare/v1.14.0...v1.15.0
[1.14.0]: https://github.com/giantswarm/cluster-standup-teardown/compare/v1.13.0...v1.14.0
[1.13.0]: https://github.com/giantswarm/cluster-standup-teardown/compare/v1.12.1...v1.13.0
[1.12.1]: https://github.com/giantswarm/cluster-standup-teardown/compare/v1.12.0...v1.12.1
[1.12.0]: https://github.com/giantswarm/cluster-standup-teardown/compare/v1.11.0...v1.12.0
[1.11.0]: https://github.com/giantswarm/cluster-standup-teardown/compare/v1.10.0...v1.11.0
[1.10.0]: https://github.com/giantswarm/cluster-standup-teardown/compare/v1.9.0...v1.10.0
[1.9.0]: https://github.com/giantswarm/cluster-standup-teardown/compare/v1.8.0...v1.9.0
[1.8.0]: https://github.com/giantswarm/cluster-standup-teardown/compare/v1.7.1...v1.8.0
[1.7.1]: https://github.com/giantswarm/cluster-standup-teardown/compare/v1.7.0...v1.7.1
[1.7.0]: https://github.com/giantswarm/cluster-standup-teardown/compare/v1.6.0...v1.7.0
[1.6.0]: https://github.com/giantswarm/cluster-standup-teardown/compare/v1.5.0...v1.6.0
[1.5.0]: https://github.com/giantswarm/cluster-standup-teardown/compare/v1.4.0...v1.5.0
[1.4.0]: https://github.com/giantswarm/cluster-standup-teardown/compare/v1.3.0...v1.4.0
[1.3.0]: https://github.com/giantswarm/cluster-standup-teardown/compare/v1.2.0...v1.3.0
[1.2.0]: https://github.com/giantswarm/cluster-standup-teardown/compare/v1.1.0...v1.2.0
[1.1.0]: https://github.com/giantswarm/cluster-standup-teardown/compare/v1.0.3...v1.1.0
[1.0.3]: https://github.com/giantswarm/cluster-standup-teardown/compare/v1.0.2...v1.0.3
[1.0.2]: https://github.com/giantswarm/cluster-standup-teardown/compare/v1.0.1...v1.0.2
[1.0.1]: https://github.com/giantswarm/cluster-standup-teardown/compare/v1.0.0...v1.0.1
[1.0.0]: https://github.com/giantswarm/cluster-standup-teardown/releases/tag/v1.0.0
