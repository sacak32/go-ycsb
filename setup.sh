#!/bin/bash

curl -OL https://golang.org/dl/go1.21.4.linux-amd64.tar.gz
sha256sum go1.21.4.linux-amd64.tar.gz
sudo tar -C /usr/local -xvf go1.21.4.linux-amd64.tar.gz
rm go1.21.4.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.profile
source ../.profile
