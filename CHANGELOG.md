# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Changed
- Updated minimum Go version to 1.22

## [0.1.0] - 2024-12-27

### Added
- Initial release
- Jinja2/Twig-inspired template syntax for Go
- Lexer → Parser → Compiler → Runtime architecture
- Template composition (layouts, includes, blocks)
- Filters for data transformation
- Control flow (if/else, range/for loops)
- Variable expressions with dot notation
- Special loop variables (@first, @last, @index)
- Template inheritance
- Comprehensive test coverage (80.7%)
- Performance optimization and caching

[Unreleased]: https://github.com/toutaio/toutago-fith-renderer/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/toutaio/toutago-fith-renderer/releases/tag/v0.1.0
