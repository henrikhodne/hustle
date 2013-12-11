#!/bin/bash

set -e
set -x

apt-get update -yq
apt-get install -yq build-essential curl git mercurial

cp -v /vagrant/.vagrant-skel/etc-default-docker /etc/default/docker

if ! which docker ; then
  curl -s https://get.docker.io | sh
fi

if ! which go ; then
  curl -s -L https://go.googlecode.com/files/go1.2.linux-amd64.tar.gz | \
    tar xzf - -C /usr/local
  ln -svf /usr/local/go/bin/* /usr/local/bin/
fi

mkdir -p /gopath
chown vagrant:vagrant /gopath

su - vagrant -c /vagrant/.vagrant-provision-as-vagrant.sh
