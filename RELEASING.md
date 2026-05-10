# Releasing samfile

Each `samfile` release is a three-step ritual. You handle steps 1 and 3; GitHub Actions handles step 2.

1. You tag the release locally and push the tag.
2. `release.yml` runs the tests, creates a **draft** GitHub Release, and attaches cross-compiled binaries.
3. You add prose release notes to the draft and publish it.

The `scripts/release.sh` helper takes care of step 1 (and the post-tag README refresh).

## Prerequisites

- Local clone on `main`, fully up to date with `origin/main`.
- Working tree clean (no modified or untracked files).
- All work intended for this release has been merged into `main` via PR (CI is the gate).
- A GPG key configured for signed tags (`git config user.signingkey <KEY-ID>`).

## Step 1 — tag the release

```
./scripts/release.sh X.Y.Z
```

For an alpha release, use the `X.Y.ZalphaN` form (e.g. `3.1.0alpha1`). The script will refuse to proceed if any of these hold:

- The version string isn't `X.Y.Z` or `X.Y.ZalphaN`.
- A `vX.Y.Z` tag already exists on `origin`.
- The working tree isn't clean.
- Local `HEAD` isn't at `origin/main` (non-alpha only).

What it does on success:

- Creates a signed annotated tag `vX.Y.Z` with message `Release X.Y.Z`.
- Pushes the tag to `origin`.
- Runs `scripts/refresh_readme.sh` to regenerate the `samfile --help` block in `README.md` with the new version banner.
- Commits the README change as `Refreshed README with samfile output from vX.Y.Z`.
- Pushes that commit to `origin/main` (non-alpha only).

## Step 2 — GitHub Actions

The tag push triggers `.github/workflows/release.yml`. Watch it at  
[github.com/petemoore/samfile/actions/workflows/release.yml](https://github.com/petemoore/samfile/actions/workflows/release.yml).

The workflow:

- Re-runs `go vet` and `go test -race`.
- Creates a **draft** GitHub Release titled `samfile X.Y.Z` if one doesn't already exist.
- Cross-builds the six platform binaries in parallel:
  - `samfile-darwin-amd64`, `samfile-darwin-arm64`
  - `samfile-linux-amd64`, `samfile-linux-arm64`
  - `samfile-windows-amd64.exe`, `samfile-windows-arm64.exe`
- Uploads each to the draft via `gh release upload --clobber`.

Total runtime is ~2 minutes.

## Step 3 — write notes and publish

Visit [github.com/petemoore/samfile/releases](https://github.com/petemoore/samfile/releases) and edit the `vX.Y.Z` draft. Write the release notes (prior-art examples: [v2.1.0](https://github.com/petemoore/samfile/releases/tag/v2.1.0), [v3.0.0](https://github.com/petemoore/samfile/releases/tag/v3.0.0)) and click **Publish**.

A typical samfile release-notes structure:

- One-paragraph framing of what's new.
- Bug-fix / feature sections, each with a short rationale and citations to the Tech Manual / ROM disasm where relevant.
- Breaking changes, if any, plus migration guidance.
- Roadmap / known limitations.

After publishing:

- `git tag -v vX.Y.Z` prints `Good signature from Pete Moore`.
- The release page lists all six assets.
- The `Latest` badge on the releases list page points at `vX.Y.Z`.
- `go install github.com/petemoore/samfile/vN/cmd/samfile@vX.Y.Z` succeeds and reports the right version banner.
