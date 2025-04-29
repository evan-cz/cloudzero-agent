#!/bin/bash

apt update

apt install -y golang-go

runuser -u ubuntu renovate
