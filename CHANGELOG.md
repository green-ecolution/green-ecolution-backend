# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [v1.0.0] - 2024-10-22

### Added

- Initial release
- Create database scheme
- Implement sqlc to generate type-safe code from SQL (#35)
- Implement repository pattern for database access for postgres 
- Added seed for demo data (#110)
- Versioning in API routes (#34)
- Implement api endpoints for trees, tree cluster, sensor, etc. (#107)
- Implement handler and service logic for trees, tree cluster, sensor, etc. (#107)
- Implement user authentication (#37)
- Implement base structure for user management (#38)
- Implement logging using slogger (#26)
- Use Linter to analyze Code 
- Implement import of trees over csv file (#43)
- Calculate tree cluster center point based on region (#86)
- Get region by coordinates (#60)
- Build, test and deploy pipeline

### Fixed
- Use pgxpool for concurrency-safe connection pool for pgx (#94)

[Unreleased]: https://github.com/green-ecolution/green-ecolution-backend/compare/v1.0.0...HEAD
[v1.0.0]: https://github.com/green-ecolution/green-ecolution-backend/compare/dfdebe...v1.0.0

