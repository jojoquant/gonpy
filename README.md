
## Ubuntu / wsl2 预准备
```
apt update -y && apt upgrade -y
apt install docker.io -y && apt install docker-compose -y
```

端口占用查询：
```
lsof -i:****
```
停止使用这个端口的程序
```
kill pid
```

## Docker 相关
docker-compose.yml [version 3 编写说明在此](https://docs.docker.com/compose/compose-file/compose-file-v3/)

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

grafana 首次登录账号：admin， 密码：admin

目前考虑使用 grafana 是因为支持echarts plugin，尽管influxdb自带了dashboard，但是简单扫了一下，好像没有更强大的可视化支持，需要后续再研究一下

[mongodb 和 mongo-express 配置说明在此](https://hub.docker.com/_/mongo)


-----------------------------------------------------
## gRPC 相关

### Install Go plugins for the protocol compiler:

Linux，使用apt或apt-get，例如：
```
$ apt install -y protobuf-compiler
$ protoc --version  # Ensure compiler version is 3+
```

MacOS，使用Homebrew：
```
$ brew install protobuf
$ protoc --version  # Ensure compiler version is 3+
```

Windows, 从 (github.com/google/protobuf/releases) 手动下载与您的操作系统和计算机体系结构 ( protoc-<version>-<os><arch>.zip)对应的 zip 文件, 下载后解压然后将bin目录添加到系统环境变量中。

protocol buffer 编译器需要一个插件来生成 Go 代码。通过运行以下命令使用 Go 1.16 或更高版本安装它：
```
go install google.golang.org/protobuf/cmd/protoc-gen-go
```

### 代码生成
在 .../gonpy/trader/grpc/proto> 下执行下面命令：
```
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    hello.proto
```


### Install Go plugins for the protocol compiler:

With Go module support (Go 1.11+), simply add the following import

```
import "google.golang.org/grpc"
```
to your code, and then go [build|run|test] will automatically fetch the necessary dependencies.

Otherwise, to install the grpc-go package, run the following command:

```
$ go get -u google.golang.org/grpc
```
