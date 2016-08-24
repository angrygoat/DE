#!/usr/bin/env bash

docker build --no-cache -t visr-docker-registry.irss.unc.edu:443/config_de:2.5 -f DockerfileConfigs .
docker push visr-docker-registry.irss.unc.edu:443/config_de:2.5
