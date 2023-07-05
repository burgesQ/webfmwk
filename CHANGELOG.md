# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantipush c Versioning](https://semver.org/spec/v2.0.0.html).

## [5.3.0] (Wed Jul 5 13:3136 2023)

### Added
- support for http2
- test are now shuffled
- pprof handler toggler

### Changed
- ShutAndWait renamed to ShutdownAndWait
- tls.GetTLSCfg now accept an tls.Config and an optional http2 param 
- tls.NewListener(IConfig) changed to tls.NewListner(tls.Config)
- Shutdown now return an error 
- WorkerLauncher is now an interface
- WorkerLauncher.Run now only take a function as unique parameter
- RegisterValidatorRule and RegisterValidatorAlias now return an error
- Shutdown now return an error
- InitServer now return an error (which should be handled as a fatal)
- Server.Shutdown now call to WaitAndStop
- WaitAndStop renamed to WaitForStop

### Fixed 
- poll ping over EOF tls connection 

### Removed 
- internal call to the logger (not in a server/context), affected functions now return an error
- WorkerLauncher call to logger
- Shutdown call to logger

## [5.2.5] (Fri Apr 7 08:41:20 2023)

### Added
- log.LogPrefix inteface
- handler.Logging example

### Changed
- handler.Logging object now prefix all request w/ the rid
- updated go toolchain from 1.15 to 1.19

### Fixed
- fixed example

## [5.2.[2,3,4]] (Fri Sep 30 11:38:53 2022, Wed Oct 19 16:35:20 2022, Mon Dec 12 14:30:42 2022)

### Added
- websocket support
- Socket.io handler support

### Changed
- TLS config implementation

## [5.2.1] (Thu Sep 29 10:08:39 2022)

### Added
- tls/cmd package, for interpolation with spf13' cobra n viper
- gitlab ci file

## [5.1.1]

### Added
- paramtrize TLS level
- tls get it own pacakge
- new GetLevel method to the IConfig interface

### Remove
- webfmwk.tls doens't exist anymore - please use webfmwk/v5/tls.(I)Config
- webfmwk.LoadTLSListener -> webfmwk/v5/tls.LoadListener
- webfmwk.GetTLSLCfg -> webfmwk/v5/tls.GetCfg

## [5.1.0] - 2020-19-09

### Added
- Support for mTLS (#43)
- LoadTLSListener and GetTLSCfg method

### Changed
- linting configuration

### Removed
- internal server methods

## [5.0.7] - 2022-06-09

### Added
- more unit test
- support for sighup
- new logger to reduce external deps
- swagger handler
- fields validation error now use json fields name
- max body size parameter

### Changed
- tls check for ca pool
- fasthttp updated

### Fixed

### Removed
- old logger
## [5.0.0] - 2021-03-05

### Added
- `GetFastContext() fasthttp.RequestCtx` to the `Context` interface
- `DecodeAndValidateQP` which is a successive call to `DecodQP` then `Validate`
- go-test github action
- logging.Handler prepend a request_id

### Changed
- `valaya/fasthttp` in favor of `net/http`
- `segmentio/encoding` in favor of `encoding/json`
- `fasthttp/router` in favor of `gorilla/mux`
- handler now live in `/webfmwk/v5/handler/{logger,recover,redoc,security}`
- updated the test to reflect the changes
- `GetQuery` signature return some `*fasthttp.Args`
- use testify for assertion

### Fixed

### Removed
- `handler.RequestID` (handled by the logger one)
- `webfmwktest` is obsolete
- dependency to burgesQ/gommon/log and burgesQ/gommon/pretty
- go-build github action

## [4.2.4] - 2020-15-12

### Changed
- redoc handler
- error comment for go-swagger compat

### Fixed
- datarace on the logger
- time.After not always gc'ed

## [4.2.1] - 2020-17-10

### Added
- NewForbidden error method
- godoc target to makefile

### Fixed
- context max output was 2014 instead of 2048

### Changed
- Wrapped support for external doc handlers
- Moved example to the `example` sub directory

## [4.2.0] - 2020-11-10

### Added
- Runner method and Address struct
- Log source IP in logger handler and 404/405

### Fixed
- pollPingEndpoint didn't follow server context
- linting

### Changed
- update golangci linting config

## [4.1.9] - 2020-10-9

### Added
- Option type to start server
- Address as param of Run
- Custom handler for 404 and 405
- Status code to response payload when possible
- Expose the validator (v10)

### Changed
- UseHandler -> Use, applyOption-> UseOption(s)
- Sleep a bit in CheckIsUp
- Gommon assertion
- Better context handling
- Validator v9 -> v10
- AnonymousError -> Error
- Trunkat logged payload to 1kb

### Fixed
- Don't poll PingEndpoint in case of tls

### Removed

## [4.1.0] - 2020-4-27

### Added
- new doc example runner
- recover hanlder
- controller return error
- ErrorHandled returned generate an API error response

### Removed
- panic/recover to error pattern
- custom context setter, please use a hanlder

### Changed
- privatizate interface implementation

### Fixed
- logger was fetched to early, fetch it via a sync.Once

## [4.0.3] - 2020.4.7

### Added
- IContext.XMLBlob method

## [4.0.2] - 2020.4.7

### Added
- new middlware implemented, the handlers use an IContext
- add the webfmwktest package which warp the httptest one

### Changed
- pass more server field in private, please use the Options object to setup the server

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
