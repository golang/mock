# Changelog
All notable changes to this project will be documented in this file.

This project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## 0.2.0 (06 Jul 2023)
### Added
- `Controller.Satisfied` that lets you check whether all expected calls
  bound to a Controller have been satisfied.
- `NewController` now takes optional `ControllerOption` parameter.
- `WithOverridableExpectations` is a `ControllerOption` that configures
  Controller to override existing expectations upon a new EXPECT().
- `-typed` flag for generating type-safe methods in the generated mock.

## 0.1.0 (29 Jun 2023)

This is a minor version that mirrors the original golang/mock
project that this project originates from.

Any users on golang/mock project should be able to migrate to
this project as-is, and expect exact same set of features (apart
from supported Go versions. See [README](README.md#supported-go-versions)
for more details.
