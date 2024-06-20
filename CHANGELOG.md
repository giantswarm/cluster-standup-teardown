# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.7.1] - 2024-06-13

### Fixed

- Upgraded `clustertest` to latest v1.1.0

### Changed

- Added `scaleDownUnneededTime` of 15m to `cluster-autoscaler` config for CAPA clusters to help avoid scale-down occurring during our scaling tests.

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

[Unreleased]: https://github.com/giantswarm/cluster-standup-teardown/compare/v1.7.1...HEAD
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
