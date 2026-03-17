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
make test
make lint
make ci
make clean
```

`make lint` builds `.bin/gclint` from `.custom-gclint.yml` and runs it against the repository.

## Current Rules

- `errshort`: forbids `:=` from reusing an existing `err` in the same scope
- `namedreturn`: forbids explicit mirror returns after named return values have already been assigned
- `chinesekey`: forbids Chinese JSON tag keys and Chinese string keys in map literals
- `layerdep`: forbids lower-level packages from importing configured higher-level packages
- `varreuse`: heuristically forbids reusing a semantic variable as a container for a different business object

## Layer Dependency Configuration

`layerdep` is configuration-driven.

Example in `.golangci.yml`:

```yaml
linters:
  settings:
    custom:
      style:
        type: module
        description: Enforce custom Go style rules with style.
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
- [ ] Meaningless spaces, blank lines, and dirty formatting are forbidden
- [ ] Different logic blocks must be clearly separated and documented when needed
- [ ] Comments must stay clear, logical, and use English when needed
- [x] Once `err` already exists in the same scope, short declarations must not reuse it
- [x] Reusing a semantic variable for another business object is checked with a heuristic semantic-token analyzer
- [ ] Large duplicated code blocks should be removed
- [x] Low-level packages must not import higher-level packages when configured through `settings.dependency_rules`
- [x] Chinese JSON tag keys and Chinese string keys in map literals are forbidden
- [x] After named returns are assigned, explicit mirror returns are forbidden and bare `return` must be used

## Current Scope

- `chinesekey` currently checks explicit `json` struct tags and explicit string keys in map literals
- `layerdep` only enforces the package prefixes listed in `settings.dependency_rules`
- `varreuse` is intentionally heuristic and focuses on descriptive variable names plus stable semantic tokens from assignment sources
