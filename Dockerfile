FROM golang:1.24.6-bullseye

RUN go version

ADD . /go/src/github.com/conplementAG/copsctl

# simple build with vet
WORKDIR /go/src/github.com/conplementAG/copsctl/cmd/copsctl
RUN go build -o copsctl && \
    go vet ./...

# run the tests
WORKDIR /go/src/github.com/conplementAG/copsctl
RUN go test ./... --cover

# simple cli run
WORKDIR /go/src/github.com/conplementAG/copsctl/cmd/copsctl
RUN ./copsctl --version

CMD [ "/bin/bash" ]