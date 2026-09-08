# 项目介绍



# 部署步骤

> 本项目基于docker服务运行，使用前需要自行配置好docker环境

分别进入目录`deploy`下的三个子目录执行下面的命令

### `deploy/api/`目录下执行

`goinspector/backend:1.0` 我的go使用的是官方包仓库，所以使用 `--build-arg`指定系统代理，**按需修改.**

```shell
docker build --network=host --build-arg http_proxy=http://127.0.0.1:7890 --build-arg https_proxy=http://127.0.0.1:7890 -f Dockerfile -t goinspector/backend:1.0 ../../backend/
```



### `deploy/frontend/`目录下执行

`goinpector/frontend:1.0`

```shell
cp *.conf ../../ # 将定制的nginx配置文件移动到镜像构建上下文的根目录上，为构建镜像做准备。
docker build -f Dockerfile -t goinspector/frontend:1.0 ../../
```



### `deploy/mysql/`目录下执行

`goinpector/mysql:1.0`

```shell
docker build -f Dockerfile -t goinspector/mysql:1.0 .
```



## 启动容器

到目录`go-inspector/`下

```shell
docker compose up -d
```



浏览器访问`localhost:8080`端口验证服务是否启动成功。



## 服务使用方法