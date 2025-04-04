# Contributing Guidelines

## How to build the tool

```bash
# Clone the project 
git clone https://github.com/conplementAG/copsctl.git

# Compile to executable
go build .
```

## How to create a release

Release creation is fully automated with gitHub workflows. A combination of release-please and goreleaser actions.
Just create a pull request to the main branch with [Conventional Commit messages](https://www.conventionalcommits.org/).
After completion of this pull request a release-please pull request is automatically created, which could be reviewed. Once approved, the release process starts. After completion, the newest release could found on the [Releases page](https://github.com/conplementAG/copsctl/releases).
