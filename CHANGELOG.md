# CHANGELOG

## v0.3.0

New features
- Add file string expr for command (#36, #37)
- Accept multiple specs from command line (#38, #39)

Breaking chagnes
- Rename `tests[].stdin.type` to `tests[].stdin.format`

## v0.2.0

New features
- Receive JSON format spec from stdin (#8)
- Receive YAML format spec from stdin (#9)
- Add `--format json` for JSON format report (#11)
- Add stream matcher `eqJSON` (#12)
- Add `--strict` for strict spec parsing (#26)
- Add `$.spexec` for version declaration (#31)

Bug fixes
- Fix output of failure message (#6)
- Make error when JSON spec has extra token (#13)
- Output occured error at exit (#16, #17)
- Fix error message for unknown stream matcher (#25)

## v0.1.0

Initial release

- Add status matchers: `eq`, `success`
- Add stream matchers: `any`, `be_empty`, `contain`, `eq`, `not`, `satisfy`
