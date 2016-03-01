#!/usr/bin/env bash

docker build -t diceunc/config_anon_files -f DockerfileConfigs .
docker build -t diceunc/config_apps -f DockerfileConfigs  .
docker build -t diceunc/config_iplant_data_apps -f DockerfileConfigs .
docker build -t diceunc/config_clockwork -f DockerfileConfigs .
docker build -t diceunc/config_clm -f DockerfileConfigs .
docker build -t diceunc/config_data_info -f DockerfileConfigs .
docker build -t diceunc/config_dewey -f DockerfileConfigs .
docker build -t diceunc/config_de_ui -f DockerfileConfigs .
docker build -t diceunc/config_iplant_data_de_ui -f DockerfileConfigs .
docker build -t diceunc/config_de_ui_nginx -f DockerfileConfigs .
docker build -t diceunc/config_iplant_data_de_ui_nginx -f DockerfileConfigs .
docker build -t diceunc/config_exim_sender -f DockerfileConfigs .
docker build -t diceunc/config_info_typer -f DockerfileConfigs .
docker build -t diceunc/config_iplant_email -f DockerfileConfigs .
docker build -t diceunc/config_jex_events -f DockerfileConfigs .
docker build -t diceunc/config_kifshare -f DockerfileConfigs .
docker build -t diceunc/config_metadata -f DockerfileConfigs .
docker build -t diceunc/config_monkey -f DockerfileConfigs .
docker build -t diceunc/config_notification_agent -f DockerfileConfigs .
docker build -t diceunc/config_saved_searches -f DockerfileConfigs .
docker build -t diceunc/config_terrain -f DockerfileConfigs .
docker build -t diceunc/config_iplant_data_terrain -f DockerfileConfigs .
docker build -t diceunc/config_tree_urls -f DockerfileConfigs .
docker build -t diceunc/config_user_preferences -f DockerfileConfigs .
docker build -t diceunc/config_user_sessions -f DockerfileConfigs .


docker push diceunc/config_anon_files
sleep 2
docker push diceunc/config_apps
sleep 2
docker push diceiunc/config_iplant_data_apps
sleep 2
docker push diceunc/config_data_info
sleep 2
docker push diceunc/config_dewey
sleep 2
docker push diceunc/config_de_ui
sleep 2
docker push diceunc/config_iplant_data_de_ui
sleep 2
docker push diceunc/config_de_ui_nginx
sleep 2
docker push diceunc/config_iplant_data_de_ui_nginx
sleep 2
docker push diceunc/config_exim_sender
sleep 2
docker push diceunc/config_info_typer
sleep 2
docker push diceunc/config_iplant_email
sleep 2
docker push diceunc/config_jex_events
sleep 2
docker push diceunc/config_kifshare
sleep 2
docker push diceunc/config_metadata
sleep 2
docker push diceunc/config_monkey
sleep 2
docker push diceunc/config_notification_agent
sleep 2
docker push diceunc/config_saved_searches
sleep 2
docker push diceunc/config_terrain
sleep 2
docker push diceunc/config_iplant_data_terrain
sleep 2
docker push diceunc/config_tree_urls
sleep 2
docker push diceunc/config_user_preferences
sleep 2
docker push diceunc/config_user_sessions
