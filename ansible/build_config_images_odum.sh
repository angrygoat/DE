#!/usr/bin/env bash

docker build -t docker.irss.unc.edu/config_happygoat -f DockerfileConfigs .
docker push docker.irss.unc.edu/config_happygoat
