#!/bin/bash
cp ../../dbscripts/initscript.sql docker/initscripts/
cp ../../swagger/api-spec.yml docker/
cd docker
docker-compose up -d
