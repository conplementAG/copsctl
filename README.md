# copsctl

## Introduction

copsctl - the conplement AG Kubernetes developer tooling

## Requirements
- [kubectl >= 1.29.9](https://kubernetes.io/docs/tasks/tools/install-kubectl/)

## Installation

Simply download a binary for your target system from [Releases](https://github.com/conplementAG/copsctl/releases), set it to you PATH, and you are ready to go.

## Contributing

For contributing to the project, and for development instructions, please check [CONTRIBUTING.md](CONTRIBUTING.md)

## Getting Started

### Connect to a cluster
`copsctl connect -e <environment-tag>`

(*Environment-tag determines the name of the cluster.*)

### Create a kubernetes namespace

`copsctl namespace create -n <namespace-name> -u John.Smith@conplement.de`

(*Namespace-name specifies the name of the kubernetes namespace.*)

### Show help

`copsctl --help`

### Show environment or cluster information

`copsctl info cluster`
`copsctl info environment`

### [Manage build agent pool](internal/corebuild/readme.md)
