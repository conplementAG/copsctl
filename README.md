# copsctl

## Introduction

copsctl - the Conplement AG Kubernetes developer tooling

[![Build Status](https://cpgithub.visualstudio.com/GitHubPipelines/_apis/build/status/conplementAG.copsctl?branchName=master)](https://cpgithub.visualstudio.com/GitHubPipelines/_build/latest?definitionId=9&branchName=master)

## Contributing

For contributing to the project, and fo development instructions, please check [CONTRIBUTING.md](CONTRIBUTING.md)

## Getting Started

### Build the tool

```bash

# Clone the project into your %GOPATH%/src/github.com/conplementAG
git clone https://github.com/conplementAG/copsctl.git

# restore dependencies (dep package manager for Golang required)
cd $GOPATH/src/github.com/conplementAG
dep ensure

# Build the binary
cd cmd/copsctl
go build .

# run the executable, e.g. ./copsctl or set this directory to your path so that copsctl works as well

```

### Basic Commands

#### Connect to a cluster
`copsctrl connect -e <environment-tag>`

(*Environment-tag determines the name of the cluster.*)

#### Create a kubernetes namespace

`copsctl namespace create -n <namespace-name>`

(*Namespace-name specifies the name of the kubernetes namespace.*)

#### Show help

`copsctl --help`