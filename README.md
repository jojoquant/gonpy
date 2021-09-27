
## Ubuntu / wsl2 预准备
```
apt update -y && apt upgrade -y
apt install docker.io -y && apt install docker-compose -y
```

## Docker 相关
启动容器:
> docker-compose up -d

更新镜像:
> docker-compose pull

dockerfile更新后, 重新build：
> docker-compose up -d --build

单独启动服务
> docker-compose up grafana

如果发现报错，对'/xxxx/xxxx/grafana/plugins'没有权限创建目录，那么就赋予权限： chmod 777 /xxx/xxx/grafana


[grafana 容器相关配置在此](https://grafana.com/docs/grafana/latest/installation/docker/)

目前考虑使用 grafana 是因为支持echarts plugin，尽管influxdb自带了dashboard，但是简单扫了一下，好像没有更强大的可视化支持，需要后续再研究一下

[mongodb 和 mongo-express 配置说明在此](https://hub.docker.com/_/mongo)