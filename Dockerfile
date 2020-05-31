FROM golang:1.13.4

RUN apt-get update && \
    apt-get install lsb-release -y

RUN apt-get install apt-transport-https

# kubectl
RUN curl -s https://packages.cloud.google.com/apt/doc/apt-key.gpg | apt-key add -
RUN touch /etc/apt/sources.list.d/kubernetes.list 
RUN echo "deb http://apt.kubernetes.io/ kubernetes-xenial main" | tee -a /etc/apt/sources.list.d/kubernetes.list
RUN apt-get update && apt-get install -y kubectl

RUN go version

ADD . /go/src/github.com/conplementAG/copsctl

# Trigger resource embedding
WORKDIR /go/src/github.com/conplementAG/copsctl/cmd/copsctl
RUN go get -u github.com/mjibson/esc
RUN go generate

# simple build
RUN go build -o copsctl .

# run the tests
WORKDIR /go/src/github.com/conplementAG/copsctl
RUN go test ./... --cover

# complex build with all platforms, optionally create Release with latest tag in GitHub as well
WORKDIR /go/src/github.com/conplementAG/copsctl/cmd/copsctl
ARG GITHUB_TOKEN
RUN if [ "x$GITHUB_TOKEN" = "x" ] ; then curl -sL http://git.io/goreleaser | VERSION=v0.123.3 bash -s -- release --skip-validate --rm-dist --skip-publish --snapshot ; else curl -sL http://git.io/goreleaser | VERSION=v0.123.3 bash -s -- release --skip-validate --rm-dist ; fi

CMD [ "/bin/bash" ]