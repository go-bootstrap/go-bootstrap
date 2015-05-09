#!/bin/bash

apt-get update
apt-get install -y software-properties-common

# Install Go, godeb will install the latest version of Go.
curl https://godeb.s3.amazonaws.com/godeb-amd64.tar.gz | tar zx -C /usr/local/bin
GOPATH=/go godeb install

# Setup Go
export GOPATH=/go
rm -rf $GOPATH/pkg/linux_amd64
echo 'GOPATH=/go' > /etc/profile.d/go.sh
echo 'PATH=$GOPATH/bin:$PATH' >> /etc/profile.d/go.sh

# Place ENV variables in /home/vagrant/.bashrc
if ! grep -Fxq "# Go Evironment Variables" /home/vagrant/.bashrc ; then
    echo -e "\n# Go Evironment Variables" >> /home/vagrant/.bashrc
    echo -e ". /etc/profile.d/go.sh" >> /home/vagrant/.bashrc
fi

GOPATH=/go go get github.com/tools/godep

cd /go/src/$GO_BOOTSTRAP_REPO_NAME/$GO_BOOTSTRAP_REPO_USER/$GO_BOOTSTRAP_PROJECT_NAME
GOPATH=/go go get ./...
