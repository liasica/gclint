GO ?= go
GOFMT ?= gofmt
GOLANGCI_LINT ?= golangci-lint
CUSTOM_GCLINT_CONFIG ?= ./.custom-gclint.yml
# golangci-lint custom only auto-loads .custom-gcl.yml.
CUSTOM_GCLINT_COMPAT_CONFIG ?= ./.custom-gcl.yml
GOLANGCI_LINT_VERSION ?= $(shell sed -n 's/^version:[[:space:]]*//p' $(CUSTOM_GCLINT_CONFIG) | head -n 1)
CUSTOM_GCLINT ?= ./.bin/gclint
RELEASE_DIR ?= ./.release
TARGET_OS ?= $(shell $(GO) env GOOS)
TARGET_ARCH ?= $(shell $(GO) env GOARCH)
TARGET_GOARM ?=
TARGET_GOMIPS ?= hardfloat
VERSION ?= dev
GO_SOURCE_FILES ?= $(shell git ls-files '*.go')

.PHONY: install-lint format-check test verify-config build-lint lint ci package-release clean print-golangci-lint-version

install-lint:
	$(GO) install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)

format-check:
	@unformatted_files="$$( $(GOFMT) -l $(GO_SOURCE_FILES) )"; \
	if [ -n "$$unformatted_files" ]; then \
		printf '%s\n' "$$unformatted_files"; \
		exit 1; \
	fi

test:
	$(GO) test ./...

verify-config:
	$(GOLANGCI_LINT) config verify -c .golangci.yml

build-lint:
	@set -eu; \
	source_config="$(CUSTOM_GCLINT_CONFIG)"; \
	compat_config="$(CUSTOM_GCLINT_COMPAT_CONFIG)"; \
	trap 'rm -f "$$compat_config"' EXIT; \
	cp "$$source_config" "$$compat_config"; \
	$(GOLANGCI_LINT) custom

lint: verify-config build-lint
	$(CUSTOM_GCLINT) run ./...

ci: format-check test lint

package-release:
	@set -eu; \
	source_config="$(CUSTOM_GCLINT_CONFIG)"; \
	compat_config="$(CUSTOM_GCLINT_COMPAT_CONFIG)"; \
	trap 'rm -f "$$compat_config"' EXIT; \
	cp "$$source_config" "$$compat_config"; \
	release_root="$$(mkdir -p "$(RELEASE_DIR)" && cd "$(RELEASE_DIR)" && pwd)"; \
	asset_arch="$(TARGET_ARCH)"; \
	goarm="$(TARGET_GOARM)"; \
	archive_format="tar.gz"; \
	package_binary_name="gclint"; \
	if [ "$(TARGET_ARCH)" = "arm" ]; then \
		if [ -z "$$goarm" ]; then \
			echo "TARGET_GOARM is required when TARGET_ARCH=arm" >&2; \
			exit 1; \
		fi; \
		asset_arch="armv$$goarm"; \
	fi; \
	if [ "$(TARGET_OS)" = "windows" ]; then \
		archive_format="zip"; \
		package_binary_name="gclint.exe"; \
		command -v zip >/dev/null 2>&1 || { echo "zip is required to package Windows artifacts" >&2; exit 1; }; \
	fi; \
	build_dir="$$release_root/build/$(TARGET_OS)_$$asset_arch"; \
	package_dir="$$release_root/package/$(TARGET_OS)_$$asset_arch"; \
	rm -rf "$$build_dir" "$$package_dir"; \
	mkdir -p "$$build_dir" "$$package_dir"; \
	GOOS="$(TARGET_OS)" GOARCH="$(TARGET_ARCH)" GOARM="$$goarm" GOMIPS="$(TARGET_GOMIPS)" CGO_ENABLED=0 "$(GOLANGCI_LINT)" custom --name gclint --destination "$$build_dir"; \
	cp "$$build_dir/gclint" "$$package_dir/$$package_binary_name"; \
	archive_path="$$release_root/gclint_$(VERSION)_$(TARGET_OS)_$$asset_arch.$$archive_format"; \
	rm -f "$$archive_path"; \
	if [ "$$archive_format" = "zip" ]; then \
		(cd "$$package_dir" && zip -q "$$archive_path" "$$package_binary_name"); \
	else \
		tar -C "$$package_dir" -czf "$$archive_path" "$$package_binary_name"; \
	fi

clean:
	rm -rf ./.bin ./.release

print-golangci-lint-version:
	@printf '%s\n' "$(GOLANGCI_LINT_VERSION)"
