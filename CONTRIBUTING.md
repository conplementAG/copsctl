# Contributing Guidelines

### How to build the tool

This project uses the Go modules, so make sure you use Golang > 1.13 for everything to work smoothly. 

```bash
# Clone the project 
git clone https://github.com/conplementAG/copsctl.git

# Embedd resources

cd cmd/copsctl
go get -u github.com/mjibson/esc
go generate .

# Compile to executable
go build .
```

*Additional Info:*

The  snippet above will put the `esc`-executable into your `$GOPATH/bin` directory, so it is available in the `go generate` phase.
This is required because `esc` will search for `yaml-files` and embedd those into the final binary, so the executable can run idependant from any working directory.

## How to create a release

Release creation is partially a manual process. You need to tag the master branch first, and then you can either:

- start the release in Azure DevOps
- or perform the release manually using [GoReleaser](https://goreleaser.com/). For reference on the command, check the Dockerfile / Azure DevOps task.

Release will be created for the new tag, including all changes since the previous tag in the changelog.