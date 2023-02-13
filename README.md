
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
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/xiaoyuzhou

//-ldflags="-s -w" 
// -s：忽略符号表和调试信息。
// -w：忽略DWARFv3调试信息，使用该选项后将无法使用gdb进行调试。


### 启动命令：
nohup ./bin/xiaoyuzhou > ./log/run.log 2>&1 &