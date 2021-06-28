# goldden-go 介绍
````
goldden-go 是一个golang的工具包，也是一个golang程序的基础服务。
````
## 作为工具包使用
````
所有的工具都在 pkg/utils 下面，包括：
auth
base_dir
captcha
config
crypto
gin_middleware
....
````
## 基础服务使用
````
1、你可以直接引用这个库，然后使用cmd/server.go 下面的init_server方法来启动服务
2、也可以直接build  main.go 然后最为一个服务启动
````
