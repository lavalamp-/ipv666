# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.4.0] - 2019-05-27
### Added
- Opt-in functionality for uploading discovered addresses to [our web site](https://ipv6.exposed/)

### Changed
- Reduced the number of hosts attempted for host-based fan-out

## [0.3.0] - 2019-03-18
### Added
- Utility for generating IPv6 addresses from a cluster model
- New predictive clustering model based off of [6gen](https://zakird.com/papers/imc17-6gen.pdf)
- Fanning out from initial discovered addresses

### Changed
- Command line syntax for a few of the subcommands

### Removed
- A bunch of unused configuration values

## [0.2.0] - 2019-01-20
### Added
- Utility for converting IPv6 files to different formats
- Dependency management with Golang modules
- ASCII art
- Changelog

### Changed
- Logging now uses log levels
- Removed dependence on Zmapv6
- Structure of repository now supports install from `go get`
- Integrated Cobra for command line invocation
- Integrated Viper for configuration management
- Packaged up assets using `packr`
- State files written to disk under a dot file in user's home directory

### Removed
- S3 integration

## 0.1.0 - 2018-11-26
### Added
- Initial release
- Reliance on Zmapv6
- Dependent on file structure in present working directory

[0.4.0]: https://github.com/lavalamp-/ipv666/compare/77f2a59...ad0302a
[0.3.0]: https://github.com/lavalamp-/ipv666/compare/f86fe91...77f2a59
[0.2.0]: https://github.com/lavalamp-/ipv666/compare/20b731c...f86fe91