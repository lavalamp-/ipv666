# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

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

[0.0.2]: https://github.com/lavalamp-/ipv666/compare/v0.1...v0.2