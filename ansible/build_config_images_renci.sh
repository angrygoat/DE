#!/usr/bin/env bash

docker build -t angrygoat.irss.unc.edu/config_happygoat -f DockerfileConfigs .
docker push angrygoat.irss.unc.edu/config_happygoat
