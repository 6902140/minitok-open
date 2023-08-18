#!/bin/bash

# 删除所有的容器
sudo docker ps -a -q | xargs sudo docker rm

# 删除名为 "minitok_tiktoklite" 的镜像
sudo docker rmi minitok_minitok

# 执行 docker-compose up
sudo docker-compose up
