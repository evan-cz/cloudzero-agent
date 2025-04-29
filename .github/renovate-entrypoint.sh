#!/bin/bash

apt update

apt install -y golang-go npm

runuser -u ubuntu renovate
