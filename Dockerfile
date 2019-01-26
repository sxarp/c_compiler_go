FROM golang:1.11

# install dep
RUN DEP_VERSION=0.5.0; curl -L -s "https://github.com/golang/dep/releases/download/v${DEP_VERSION}/dep-linux-amd64" -o /bin/dep; chmod +x /bin/dep

RUN apt-get update

RUN apt-get install less
