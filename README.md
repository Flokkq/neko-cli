# Neko - CLI

Neko is a universal CLI tool for orchestrating release workflows across frontend and backend projects.

Neko is a command-line tool designed to streamline release management
for the Nekoman project. It helps developers and release engineers automate
common tasks such as version bumping, changelog generation, and deployment.

With neko-cli, you can:
- Initialize a new release management
- Automatically update version numbers
- Generate changelogs from git commits
- Validate release readiness

---
## Hot To Use

**Global Flags**

`-h` Help displays 

`-v` Verbose Output

## Commands

### `neko init`
Initialize Neko in the current project. 

### `neko release`
Run the release process using the detected or configured tool.  
**Args / Flags:**
- `--dry-run` : simulate the release
- `--push` : push changes and tags (Default)

### `neko config`
Show or edit the Neko configuration.  
**Args / Flags:**
- `--show` : display current configuration
- `--edit` : open configuration in editor

### `neko version`
Show current version of this repo.  
**Args / Flags:**
- `--set=<version>` : set specific version

### `neko history`
Show release/tag history.  

### `neko status`
Display current release status.  
**Checks include:** git clean state, branch, version file, changelog status

### `neko release-check`
Validate whether the project is ready for release (pre-flight checks).





