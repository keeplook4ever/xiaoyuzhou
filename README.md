
## How to run

### Required

- Mysql
- Redis

### Ready

Create a **blog database** and execute SQL from docs/sql/blog.sql

### Conf

You should modify `conf/app.ini`

```
[database]
Type = mysql
User = root
Password =
Host = 127.0.0.1:3306
Name = blog


[redis]
Host = 127.0.0.1:6379
Password =
MaxIdle = 30
MaxActive = 30
IdleTimeout = 200
...
```

### Run
```
$ cd $GOPATH/src/go-gin-example

$ go run main.go 
```


## 发布
### Build命令
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/xiaoyuzhou



### 启动命令：
nohup ./bin/xiaoyuzhou > ./log/run.log 2>&1 &