#!/bin/bash

TAG=0.0.10
docker build -t sqzxcv/monitor_logger:$TAG .
docker push sqzxcv/monitor_logger:$TAG

docker tag sqzxcv/monitor_logger:$TAG sqzxcv/monitor_logger
docker push sqzxcv/monitor_logger
#docker build -t sqzxcv/monitor_logger:latest .
#docker push sqzxcv/monitor_logger:latest
