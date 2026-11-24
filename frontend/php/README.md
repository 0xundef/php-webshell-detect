# Vendored PHP frontend

This directory is a fork of [z7zmey/php-parser](https://github.com/z7zmey/php-parser)
(MIT, see [LICENSE](LICENSE)), adapted to parse PHP 5/7 into the AST consumed by
`internal/core/ir`.

- `version/scanner` and parts of `version/php5`/`php7` are generated (ragel) and
  are checked in; they are **not** linted or vetted by this repository.
- Local changes are limited to the module import path and small
  analyser-driven adjustments. The parser covers PHP 5 and PHP 7; the analyser
  configures the PHP 7.4 grammar.
- Do not edit generated files by hand. Upstream is unmaintained; patch only what
  the analyser needs and add a test for every grammar change.
