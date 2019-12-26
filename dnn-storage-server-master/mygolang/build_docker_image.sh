#!/bin/sh -eu
docker login
docker build -t macchiang/mygolang:1.8.3 .
docker push macchiang/mygolang:1.8.3