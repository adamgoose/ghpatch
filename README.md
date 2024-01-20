# ghpatch

ghpatch is a simple Golang utility for patching files on GitHub.

## man

```
ghpatch is a tool for patching GitHub files

Usage:
  ghpatch {path} [flags] -- {cmd}

Flags:
  -e, --author-email string   [GHPATCH_COMMIT_AUTHOR_EMAIL] Commit Author email
  -n, --author-name string    [GHPATCH_COMMIT_AUTHOR_NAME] Commit Author name
  -b, --branch string         [GHPATCH_GITHUB_BRANCH] GitHub branch
      --dry-run               [GHPATCH_DRY_RUN] Dry run
  -h, --help                  help for ghpatch
  -m, --message string        [GHPATCH_COMMIT_MESSAGE] Commit message
  -o, --owner string          [GHPATCH_GITHUB_OWNER] GitHub owner
  -r, --repo string           [GHPATCH_GITHUB_REPO] GitHub repo
  -t, --token string          [GHPATCH_GITHUB_TOKEN] GitHub token

```
