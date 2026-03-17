# gclint

English | [中文](README.zh-CN.md)

Practical Go style, shipped as `gclint`.

`gclint` builds a custom `golangci-lint` binary named `gclint`.

The repository currently contains one module plugin package: `style`.

## Installation

Releases use the version format `YYYY.MM.DD-SHORT_HASH_ID`, for example `2026.03.17-deadbee`.

The release workflow derives its build matrix from the official `golangci-lint` release assets for the pinned version in `.custom-gclint.yml`.

### Install the latest release with `install.sh`

`install.sh` detects the current platform automatically and downloads the matching asset.

```bash
curl -fsSL https://raw.githubusercontent.com/liasica/gclint/master/install.sh | sh
```

Install to a custom directory:

```bash
curl -fsSL https://raw.githubusercontent.com/liasica/gclint/master/install.sh | GCLINT_INSTALL_DIR="$HOME/.local/bin" sh
```

Install a specific release:

```bash
curl -fsSL https://raw.githubusercontent.com/liasica/gclint/master/install.sh | GCLINT_VERSION="2026.03.17-deadbee" sh
```

Override auto-detection when needed:

```bash
curl -fsSL https://raw.githubusercontent.com/liasica/gclint/master/install.sh | \
  GCLINT_OS=linux GCLINT_ARCH=armv7 GCLINT_INSTALL_DIR="$HOME/.local/bin" sh
```

### Download a release asset manually

Release assets are published as `gclint_<version>_<os>_<arch>.tar.gz` for Unix-like targets and `gclint_<version>_<os>_<arch>.zip` for Windows.

Current release targets:

- `darwin`: `amd64`, `arm64`
- `freebsd`: `386`, `amd64`, `arm64`, `armv6`, `armv7`
- `illumos`: `amd64`
- `linux`: `386`, `amd64`, `arm64`, `armv6`, `armv7`, `loong64`, `mips64`, `mips64le`, `ppc64le`, `riscv64`, `s390x`
- `netbsd`: `386`, `amd64`, `arm64`, `armv6`, `armv7`
- `windows`: `386`, `amd64`, `arm64`

### Build from source

```bash
make install-lint
make build-lint
install -m 0755 ./.bin/gclint /usr/local/bin/gclint
```

## Quick Start

```bash
make install-lint
make format-check
make test
make lint
make ci
make clean
```

`make format-check` validates `gofmt`. `make lint` builds `.bin/gclint` from `.custom-gclint.yml` and runs it against the repository.

## Current Rules

Custom analyzers:

- `redeclare`: forbids short declarations from reusing an existing variable name in the current function, including inner-block shadowing and `err`
- `namedreturn`: forbids explicit returns after named return values have already been assigned
- `chinesekey`: forbids Chinese JSON tag keys, Chinese string keys in persistent maps, and Chinese keys inside raw JSON string constants
- `layerdep`: forbids lower-level packages from importing configured higher-level packages
- `varreuse`: heuristically forbids reusing a semantic variable as a container for a different business object

Official linters enabled by default:

- `dupl`: detects large duplicated code fragments
- `dupword`: detects duplicated words in code and comments
- `godot`: normalizes sentence-ending punctuation in comments
- `misspell`: detects common English misspellings
- `whitespace`: forbids meaningless leading or trailing blank lines in blocks
- `wsl_v5`: enforces empty-line structure between logical blocks

## Layer Dependency Configuration

`layerdep` is configuration-driven.

Example in `.golangci.yml`:

```yaml
linters:
  settings:
    custom:
      style:
        type: module
        description: Enforce custom Go style rules with the gclint style plugin.
        settings:
          dependency_rules:
            - source: github.com/example/project/internal/repository
              forbidden:
                - github.com/example/project/internal/service
                - github.com/example/project/internal/handler
```

Replace those package prefixes with your real architecture layers.

## Rule Coverage

The checklist below maps the current repository status to the Go style document used for this project.

- [ ] Names must match real business semantics and avoid arbitrary abbreviations
- [ ] Singular and plural naming must stay accurate
- [ ] Different business objects must not share the same name
- [x] Same-scope short declarations must not reuse an existing variable
- [x] Meaningless spaces, blank lines, and dirty formatting are checked with `gofmt`, `whitespace`, and `wsl_v5`
- [ ] Different logic blocks must be clearly separated and documented when needed
- [ ] Comments must stay clear, logical, and use English when needed
- [x] Once `err` already exists in the same scope, short declarations must not reuse it
- [x] Reusing a semantic variable for another business object is checked with a heuristic semantic-token analyzer
- [x] Large duplicated code blocks are checked with `dupl`
- [x] Low-level packages must not import higher-level packages when configured through `settings.dependency_rules`
- [x] Chinese JSON tag keys, persistent map keys, and raw JSON string keys are forbidden
- [x] After named returns are assigned, explicit returns are forbidden and bare `return` must be used

## Current Scope

- The repository includes analyzer fixtures under `style/testdata/src`, built from recommended and discouraged examples in the Go style document.
- `chinesekey` checks explicit `json` struct tags, string-keyed map literals, string-keyed map index assignments, and raw JSON string constants
- `layerdep` only enforces the package prefixes listed in `settings.dependency_rules`
- `namedreturn` is intentionally strict: once a named return value has been assigned, later explicit returns are rejected
- `redeclare` covers same-block reuse and inner-block shadowing such as `if err := ...` inside the same function
- `varreuse` is intentionally heuristic and focuses on descriptive variable names plus stable semantic tokens from assignment sources
- Business-semantic naming, singular/plural correctness, and comment quality still require code review judgment
