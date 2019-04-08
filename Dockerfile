FROM golang:latest

RUN apt-get update && \
    apt-get install lsb-release -y

# Azure CLI
RUN echo "deb [arch=amd64] https://packages.microsoft.com/repos/azure-cli/ $(lsb_release -cs) main" > /etc/apt/sources.list.d/azure-cli.list
RUN curl -L https://packages.microsoft.com/keys/microsoft.asc | apt-key add -
RUN apt-get install apt-transport-https
RUN apt-get update && apt-get install azure-cli

# kubectl
RUN curl -s https://packages.cloud.google.com/apt/doc/apt-key.gpg | apt-key add -
RUN touch /etc/apt/sources.list.d/kubernetes.list 
RUN echo "deb http://apt.kubernetes.io/ kubernetes-xenial main" | tee -a /etc/apt/sources.list.d/kubernetes.list
RUN apt-get update && apt-get install -y kubectl

RUN go version

# dep package manager
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
 
ADD . /go/src/github.com/conplementAG/copsctl
WORKDIR /go/src/github.com/conplementAG/copsctl

RUN dep ensure

WORKDIR /go/src/github.com/conplementAG/copsctl/cmd/copsctl
RUN go build -o copsctl .

WORKDIR /go/src/github.com/conplementAG/copsctl
RUN go test ./... --cover

CMD [ "/bin/bash" ]