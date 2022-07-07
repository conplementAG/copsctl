# Contributing Guidelines

### How to build the tool

```bash
# Clone the project 
git clone https://github.com/conplementAG/copsctl.git

# Compile to executable
go build .
```

## How to create a release

Release creation is partially a manual process. You need to tag the master branch first, and then you can either:

- start the release in Azure DevOps
- or perform the release manually using [GoReleaser](https://goreleaser.com/). For reference on the command, check the Dockerfile / Azure DevOps task.

Release will be created for the new tag, including all changes since the previous tag in the changelog.