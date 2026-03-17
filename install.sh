#!/bin/sh

set -eu

REPOSITORY="${GCLINT_REPOSITORY:-${GCL_REPOSITORY:-liasica/gclint}}"
INSTALL_DIRECTORY="${GCLINT_INSTALL_DIR:-${GCL_INSTALL_DIR:-}}"
REQUESTED_VERSION="${GCLINT_VERSION:-${GCL_VERSION:-}}"
REQUESTED_OS="${GCLINT_OS:-${GCL_OS:-}}"
REQUESTED_ARCH="${GCLINT_ARCH:-${GCL_ARCH:-}}"
BINARY_NAME="gclint"

log() {
	printf '%s\n' "$*"
}

fail() {
	printf '%s\n' "$*" >&2
	exit 1
}

download_to_stdout() {
	url="$1"

	if command -v curl >/dev/null 2>&1; then
		if [ -n "${GITHUB_TOKEN:-}" ]; then
			curl -fsSL \
				-H "Authorization: Bearer ${GITHUB_TOKEN}" \
				-H "Accept: application/vnd.github+json" \
				"$url"
			return
		fi

		curl -fsSL -H "Accept: application/vnd.github+json" "$url"
		return
	fi

	if command -v wget >/dev/null 2>&1; then
		if [ -n "${GITHUB_TOKEN:-}" ]; then
			wget -qO- \
				--header="Authorization: Bearer ${GITHUB_TOKEN}" \
				--header="Accept: application/vnd.github+json" \
				"$url"
			return
		fi

		wget -qO- "$url"
		return
	fi

	fail "curl or wget is required"
}

download_to_file() {
	url="$1"
	destination="$2"

	if command -v curl >/dev/null 2>&1; then
		if [ -n "${GITHUB_TOKEN:-}" ]; then
			curl -fsSL \
				-H "Authorization: Bearer ${GITHUB_TOKEN}" \
				-H "Accept: application/vnd.github+json" \
				-o "$destination" \
				"$url"
			return
		fi

		curl -fsSL -o "$destination" "$url"
		return
	fi

	if command -v wget >/dev/null 2>&1; then
		if [ -n "${GITHUB_TOKEN:-}" ]; then
			wget -qO "$destination" \
				--header="Authorization: Bearer ${GITHUB_TOKEN}" \
				--header="Accept: application/vnd.github+json" \
				"$url"
			return
		fi

		wget -qO "$destination" "$url"
		return
	fi

	fail "curl or wget is required"
}

detect_os() {
	if [ -n "$REQUESTED_OS" ]; then
		case "$REQUESTED_OS" in
			linux|darwin|windows|freebsd|netbsd|illumos)
				printf '%s\n' "$REQUESTED_OS"
				return
				;;
			*)
				fail "unsupported operating system override: ${REQUESTED_OS}"
				;;
		esac
	fi

	case "$(uname -s)" in
		Linux)
			printf '%s\n' "linux"
			;;
		Darwin)
			printf '%s\n' "darwin"
			;;
		FreeBSD)
			printf '%s\n' "freebsd"
			;;
		NetBSD)
			printf '%s\n' "netbsd"
			;;
		SunOS)
			printf '%s\n' "illumos"
			;;
		CYGWIN*|MINGW*|MSYS*)
			printf '%s\n' "windows"
			;;
		*)
			fail "unsupported operating system: $(uname -s)"
			;;
	esac
}

detect_arch() {
	if [ -n "$REQUESTED_ARCH" ]; then
		case "$REQUESTED_ARCH" in
			386|amd64|arm64|armv6|armv7|loong64|mips64|mips64le|ppc64le|riscv64|s390x)
				printf '%s\n' "$REQUESTED_ARCH"
				return
				;;
			*)
				fail "unsupported architecture override: ${REQUESTED_ARCH}"
				;;
		esac
	fi

	case "$(uname -m)" in
		x86_64|amd64)
			printf '%s\n' "amd64"
			;;
		i386|i686|x86)
			printf '%s\n' "386"
			;;
		i86pc)
			printf '%s\n' "amd64"
			;;
		arm64|aarch64)
			printf '%s\n' "arm64"
			;;
		armv6|armv6l)
			printf '%s\n' "armv6"
			;;
		armv7|armv7l)
			printf '%s\n' "armv7"
			;;
		loong64|loongarch64)
			printf '%s\n' "loong64"
			;;
		mips64)
			printf '%s\n' "mips64"
			;;
		mips64el|mips64le)
			printf '%s\n' "mips64le"
			;;
		ppc64le)
			printf '%s\n' "ppc64le"
			;;
		riscv64)
			printf '%s\n' "riscv64"
			;;
		s390x)
			printf '%s\n' "s390x"
			;;
		*)
			fail "unsupported architecture: $(uname -m)"
			;;
	esac
}

resolve_install_directory() {
	operating_system="$1"

	if [ -n "$INSTALL_DIRECTORY" ]; then
		printf '%s\n' "$INSTALL_DIRECTORY"
		return
	fi

	if [ "$operating_system" = "windows" ]; then
		printf '%s\n' "${HOME}/bin"
		return
	fi

	if [ -w /usr/local/bin ]; then
		printf '%s\n' "/usr/local/bin"
		return
	fi

	printf '%s\n' "${HOME}/.local/bin"
}

resolve_archive_extension() {
	operating_system="$1"

	case "$operating_system" in
		windows)
			printf '%s\n' "zip"
			;;
		*)
			printf '%s\n' "tar.gz"
			;;
	esac
}

resolve_packaged_binary_name() {
	operating_system="$1"

	case "$operating_system" in
		windows)
			printf '%s\n' "${BINARY_NAME}.exe"
			;;
		*)
			printf '%s\n' "${BINARY_NAME}"
			;;
	esac
}

resolve_version() {
	if [ -n "$REQUESTED_VERSION" ]; then
		printf '%s\n' "$REQUESTED_VERSION"
		return
	fi

	release_json="$(download_to_stdout "https://api.github.com/repos/${REPOSITORY}/releases/latest")"
	release_version="$(printf '%s' "$release_json" | tr -d '\n' | sed -n 's/.*"tag_name"[[:space:]]*:[[:space:]]*"\([^"]*\)".*/\1/p')"

	if [ -z "$release_version" ]; then
		fail "failed to resolve latest release version from GitHub"
	fi

	printf '%s\n' "$release_version"
}

verify_checksum() {
	checksum_file="$1"

	if command -v sha256sum >/dev/null 2>&1; then
		sha256sum -c "$checksum_file"
		return
	fi

	if command -v shasum >/dev/null 2>&1; then
		shasum -a 256 -c "$checksum_file"
		return
	fi

	fail "sha256sum or shasum is required to verify the download"
}

extract_archive() {
	archive_path="$1"
	destination="$2"

	case "$archive_path" in
		*.tar.gz)
			tar -xzf "$archive_path" -C "$destination"
			;;
		*.zip)
			if command -v unzip >/dev/null 2>&1; then
				unzip -q "$archive_path" -d "$destination"
				return
			fi

			if command -v tar >/dev/null 2>&1; then
				tar -xf "$archive_path" -C "$destination" >/dev/null 2>&1 && return
			fi

			fail "unzip is required to extract zip archives"
			;;
		*)
			fail "unsupported archive format: ${archive_path}"
			;;
	esac
}

install_binary() {
	source_path="$1"
	destination_path="$2"

	if command -v install >/dev/null 2>&1; then
		install -m 0755 "$source_path" "$destination_path"
		return
	fi

	cp "$source_path" "$destination_path"
	chmod 0755 "$destination_path"
}

main() {
	operating_system="$(detect_os)"
	architecture="$(detect_arch)"
	version="$(resolve_version)"
	target_directory="$(resolve_install_directory "$operating_system")"
	archive_extension="$(resolve_archive_extension "$operating_system")"
	packaged_binary_name="$(resolve_packaged_binary_name "$operating_system")"
	asset_name="${BINARY_NAME}_${version}_${operating_system}_${architecture}.${archive_extension}"
	download_url="https://github.com/${REPOSITORY}/releases/download/${version}/${asset_name}"
	checksum_url="https://github.com/${REPOSITORY}/releases/download/${version}/checksums.txt"
	temporary_directory="$(mktemp -d)"

	trap 'find "$temporary_directory" -type f -delete 2>/dev/null || true; find "$temporary_directory" -depth -type d -empty -delete 2>/dev/null || true' EXIT INT TERM

	log "Installing ${BINARY_NAME} ${version} for ${operating_system}/${architecture}"

	mkdir -p "$target_directory"

	download_to_file "$download_url" "${temporary_directory}/${asset_name}"
	download_to_file "$checksum_url" "${temporary_directory}/checksums.txt"

	grep " ${asset_name}\$" "${temporary_directory}/checksums.txt" > "${temporary_directory}/${asset_name}.sha256" || fail "missing checksum for ${asset_name}"
	(
		cd "$temporary_directory"
		verify_checksum "${asset_name}.sha256"
	)

	extract_archive "${temporary_directory}/${asset_name}" "$temporary_directory"

	if [ ! -f "${temporary_directory}/${packaged_binary_name}" ]; then
		fail "binary ${packaged_binary_name} not found in ${asset_name}"
	fi

	install_binary "${temporary_directory}/${packaged_binary_name}" "${target_directory}/${packaged_binary_name}"

	log "Installed ${packaged_binary_name} to ${target_directory}/${packaged_binary_name}"
}

main "$@"
