#!/bin/bash

set -e
set -x

ln -svf /vagrant/.vagrant-skel/bashrc ~/.bashrc

source ~/.bashrc

go get -x github.com/kr/godep
