# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantipush c Versioning](https://semver.org/spec/v2.0.0.html).

## [3.2.0] - 2020.3.3

### Added
- TLSConfig option tweaked

### Fixed
- linter errors

### Changed
- validator and translator inited via a sync.Once

## [3.1.0] - 2020.2.21

### Added
- Option type to init server
- Error Service Unavailable
- Dump method, which dump exposed routes 
- isReady channel pattern
- More unit test

### Changed
- InitServer method signature
- Log routes groups name instead of full routes
- Validation translator init once

### Fixed
- golangci linter

## [3.0.1] - 2020.1.21

### Fixed
- Wrong custom context implementation ... 

## [3.0.0] - 2020.1.20

### Added
- go 1.13 support
- IErrorHandler can now wrap errors
- JSONAccepted method
- IContext setter and getter for native context.Context

### Changed
- IContext Set* method call can be chained

## [2.5.0] - 2019.12.17

### Added
- Translation of validation error

### Changed
- `psjon` became `pretty`
- test for the route package

### Fixed
- logging middleware was using the wrong logger

Initial release

v1.0 and v1.1 where broken

### Added
### Fixed
### Changed
### Removed
