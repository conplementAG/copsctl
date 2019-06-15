# copsctl

## Introduction

copsctl - the Conplement AG Kubernetes developer tooling

[![Build Status](https://cpgithub.visualstudio.com/GitHubPipelines/_apis/build/status/conplementAG.copsctl?branchName=master)](https://cpgithub.visualstudio.com/GitHubPipelines/_build/latest?definitionId=9&branchName=master)

## Installation

Simply download a binary for your target system from [Releases](https://github.com/conplementAG/copsctl/releases), set it to you PATH, and you are ready to go.

## Contributing

For contributing to the project, and fo development instructions, please check [CONTRIBUTING.md](CONTRIBUTING.md)

## Getting Started

#### Connect to a cluster
`copsctrl connect -e <environment-tag>`

(*Environment-tag determines the name of the cluster.*)

#### Create a kubernetes namespace

`copsctl namespace create -n <namespace-name> -u John.Smith@conplement.de`

(*Namespace-name specifies the name of the kubernetes namespace.*)

#### Show help

`copsctl --help`