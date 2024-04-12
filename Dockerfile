FROM golang:1.22.2-bullseye

RUN apt-get update && \
    apt-get install lsb-release -y

RUN apt-get install apt-transport-https

RUN go version

ADD . /go/src/github.com/conplementAG/copsctl

# simple build
WORKDIR /go/src/github.com/conplementAG/copsctl/cmd/copsctl
RUN go build -o copsctl .

# run the tests
WORKDIR /go/src/github.com/conplementAG/copsctl
RUN go test ./... --cover

# complex build with all platforms, optionally create Release with latest tag in GitHub as well
WORKDIR /go/src/github.com/conplementAG/copsctl/cmd/copsctl
ARG GITHUB_TOKEN
ARG GO_RELEASER_VERSION=v1.25.1
RUN if [ "x$GITHUB_TOKEN" = "x" ] ; then curl -sL https://git.io/goreleaser | VERSION=${GO_RELEASER_VERSION} bash -s -- release --skip-validate --rm-dist --skip-publish --snapshot ; else curl -sL https://git.io/goreleaser | VERSION=${GO_RELEASER_VERSION} bash -s -- release --skip-validate --rm-dist ; fi

CMD [ "/bin/bash" ]