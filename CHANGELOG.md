# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- Add `CiliumEniModeClusterBuilder`

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

[Unreleased]: https://github.com/giantswarm/cluster-standup-teardown/compare/v1.0.2...HEAD
[1.0.2]: https://github.com/giantswarm/cluster-standup-teardown/compare/v1.0.1...v1.0.2
[1.0.1]: https://github.com/giantswarm/cluster-standup-teardown/compare/v1.0.0...v1.0.1
[1.0.0]: https://github.com/giantswarm/cluster-standup-teardown/releases/tag/v1.0.0
