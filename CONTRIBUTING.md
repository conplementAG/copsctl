# Contributing Guidelines

## How to get the source-code

`go get github.com/conplementAG/copsctl` will install the package into your `$GOPATH`.

## How to build the project

1. Resolve all dependencies

`cd $GOPATH/src/github.com/conplementAG/copsctl`
`dep ensure -v`

2. Build the tool

`cd $GOPATH/src/github.com/conplementAG/copsctl/cmd/copsctl`
`go build .`

## How to create a release

Release creation is partially a manual process. You need to tag the master branch first, and then you can either:

- start the release in Azure DevOps
- or perform the release manually using [GoReleaser](https://goreleaser.com/). For reference on the command, check the Azure DevOps release task.