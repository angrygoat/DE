#!/usr/bin/env bash

docker build -t dfc-de-db.edc.renci.org/config_de -f DockerfileConfigs .
docker push dfc-de-db.edc.renci.org/config_de
