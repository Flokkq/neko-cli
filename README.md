
<div align="center">
<h1>Neko Cli</h1>
<img alt="Neko-Cli Logo" height="500" src="neko-cli-logo.png" width="500"/>

<br />
</div>


  [![GitHub release](https://img.shields.io/github/v/release/nekoman-hq/neko-cli?style=flat-square)](https://github.com/nekoman-hq/neko-cli/releases)
  [![Go Report Card](https://goreportcard.com/badge/github.com/nekoman-hq/neko-cli)](https://goreportcard.com/report/github.com/nekoman-hq/neko-cli)
  [![Issues](https://img.shields.io/github/issues/nekoman-hq/neko-cli?style=flat-square)](https://github.com/nekoman-hq/neko-cli/issues)
  [![Last Commit](https://img.shields.io/github/last-commit/nekoman-hq/neko-cli?style=flat-square)](https://github.com/nekoman-hq/neko-cli/commits)
  [![Contributors](https://img.shields.io/github/contributors/nekoman-hq/neko-cli?style=flat-square)](https://github.com/nekoman-hq/neko-cli/graphs/contributors)

---

Neko is a **universal CLI tool** for orchestrating release workflows across frontend and backend projects.

It helps developers and release engineers automate tasks like:

- 🚀 Initialize new release workflows
- 🔄 Automatically update version numbers
- ✅ Validate release readiness

---

## ⚡ Getting Started

**Requirements**

- A GitHub Personal Access Token named `GITHUB_TOKEN`
- CLI tool of your chosen release system (e.g., [goreleaser](https://goreleaser.com/install/))

**Global Flags**

| Flag | Description |
|------|-------------|
| `-h` | Show help |
| `-v` | Verbose output |

---

## 🛠 Commands

### `neko init`
Initialize Neko in the current project with the underlying release system.

**Supported Systems:** `goreleaser`, `release-it`, `jreleaser`

### `neko release`
Run the release process using the detected or configured tool.

**Args / Flags:**
- `patch` : increment by 0.0.1
- `minor` : increment by 0.1.0
- `major` : increment by 1.0.0

### `neko version`
Show or set the current version of the repo.

**Args / Flags:**
- `--set=<version>` : set specific version (in progress)

### `neko validate`
Validate or show the Neko configuration.

**Args / Flags:**
- `--config-show` : display current configuration

### `neko history`
Show release/tag history

### `neko status` *(in progress)*
Display current release status (checks include git clean state, branch, version file, changelog status)

### `neko check-release` *(in progress)*
Validate project readiness for release (pre-flight checks)

### `neko start` *(in progress)*
Start the whole project with one command