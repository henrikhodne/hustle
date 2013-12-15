#!/bin/bash

set -e
set -x

apt-get update -yq
apt-get install -yq curl python-software-properties

apt-add-repository -y ppa:nginx/development

apt-get install -yq build-essential git mercurial screen byobu nginx

cp -v /vagrant/.vagrant-skel/etc-default-docker /etc/default/docker

if ! which docker ; then
  curl -s https://get.docker.io | sh
fi

if ! go env ; then
  curl -s -L https://go.googlecode.com/files/go1.2.linux-amd64.tar.gz | \
    tar xzf - -C /usr/local
  ln -svf /usr/local/go/bin/* /usr/local/bin/
fi

mkdir -p /gopath
chown vagrant:vagrant /gopath

if ! redis-server -v ; then
  REDIS_VERSION=2.8.3
  mkdir -p /var/tmp/redis-$REDIS_VERSION
  pushd /var/tmp/redis-$REDIS_VERSION
  curl -s -L http://download.redis.io/releases/redis-${REDIS_VERSION}.tar.gz | \
	tar xzf - --strip-components=1
  make
  cp -v src/redis-server /usr/local/bin/
  cp -v src/redis-cli /usr/local/bin/
fi

if ! which shoreman ; then
  curl -s -L -o /usr/local/bin/shoreman https://raw.github.com/hecticjeff/shoreman/master/shoreman.sh
  chmod +x /usr/local/bin/shoreman
fi

service nginx stop
ln -svf /vagrant/.vagrant-skel/nginx.conf /etc/nginx/nginx.conf
rm -vf /etc/nginx/sites-enabled/hustle
service nginx start

stop redis-server || true
cp -v  /vagrant/.vagrant-skel/etc-init-redis-server.conf /etc/init/redis-server.conf
start redis-server

su - vagrant -c /vagrant/.vagrant-provision-as-vagrant.sh
