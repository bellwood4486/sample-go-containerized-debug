# sample-go-containerized-debug

## 0. Code Base

### How to run
```
❯ cd 0_base
❯ docker-compose up -d --build
❯ curl http://localhost:8080
hello world
```

## 1. Debug with Delve

### How to run
```
❯ cd 1_dlv
❯ docker-compose up -d --build
```

2345と8080をリッスンしているのがわかります。
```
❯ dc logs -f
Attaching to 1_dlv_app_1
app_1  | API server listening at: [::]:2345
app_1  | Launching server at ":8080" ...
```

curlで結果が返ります。
```
❯ curl http://localhost:8080
hello world
```

In the container, there are 2 processes. One is Delve's process. Another is an application process run via Delve.
```
❯ docker-compose exec app ash
/go/src/github.com/bellwood4486/sample-go-containerized-debug/1_dlv # ps
PID   USER     TIME  COMMAND
    1 root      0:00 /go/bin/dlv exec /bin/sample-go-server-debug --listen=:2345 --headless --api-version=2 --continue --accept-multiclient
   12 root      0:00 /bin/sample-go-server-debug
```

ツリー表示してみると、dlvから呼ばれていることがわかる。
```
/go/src/github.com/bellwood4486/sample-go-containerized-debug/1_dlv # pstree -p
dlv(1)---sample-go-serve(12)
```

### How to debug on GoLand

`Go Remote`設定を追加する

![スクリーンショット 2020-04-18 12 27 27](https://user-images.githubusercontent.com/2452581/79627378-9b4d2300-8172-11ea-93a7-d095e95b086d.png)

ポートはDockerfileで定義した`2345`を指定。

![スクリーンショット 2020-04-18 12 29 14](https://user-images.githubusercontent.com/2452581/79627388-aef88980-8172-11ea-9a4a-7aaf5b72e3ab.png)


デバッガからアタッチすると次のメッセージが表示されるが、今のところ無視してる。
```
❯ dc logs -f
Attaching to 1_dlv_app_1
... snip ...
app_1  | 2020-04-18T07:03:47Z error layer=rpc writing response:write tcp 127.0.0.1:2345->127.0.0.1:48800: use of closed network connection
```

ハンドラ内に貼ったブレイクポイントをデバッグできる。

![スクリーンショット 2020-04-18 12 35 52](https://user-images.githubusercontent.com/2452581/79627399-d64f5680-8172-11ea-8d56-96c01efd0882.png)

注意：main関数をデバッグしたい…TODO

## 2. Delve + Realize

### How to run
```
❯ cd 2_dlv_realize
❯ docker-compose up -d --build
❯ curl http://localhost:8080
hello world
```

realizeの出力を確認できる。ビルドが7.766 sで完了し、2345と8080ポートをリッスンしているのがわかる。
```
❯ docker-compose logs -f
Attaching to 2_dlv_realize_app_1
app_1  | [05:58:33][APP-DEBUG] : Watching 1 file/s 1 folder/s
app_1  | [05:58:33][APP-DEBUG] : Command "pkill -INT /bin/sample-go-server-debug"
app_1  | [05:58:33][APP-DEBUG] : Build started
app_1  | [05:58:41][APP-DEBUG] : Build completed in 7.766 s
app_1  | [05:58:41][APP-DEBUG] : Running..
app_1  | [05:58:41][APP-DEBUG] : API server listening at: [::]:2345
app_1  | [05:58:41][APP-DEBUG] : Launching server at ":8080" ...
```

3つのプロセスが起動しているのがわかる。また、realizeからdelveを経由して、アプリケーションが起動されているのがわかる。
```
❯ dc exec app ash
/go/src/github.com/bellwood4486/sample-go-containerized-debug/2_dlv_realize # ps
PID   USER     TIME  COMMAND
    1 root      0:00 /go/bin/realize start
 1000 root      0:00 /go/bin/dlv exec /bin/sample-go-server-debug --listen :2345 --headless --api-version 2 --continue --accept-multiclient
 1007 root      0:00 /bin/sample-go-server-debug
 1018 root      0:00 ash
 1023 root      0:00 ps
/go/src/github.com/bellwood4486/sample-go-containerized-debug/2_dlv_realize # pstree -p
realize(1)---dlv(1000)---sample-go-serve(1007)
```

コードを変更してみる。
main.go
```go
//... snip ...
		_, _ = fmt.Fprintf(w, "foo bar")
//... snip ...
```

変更が検知されて、再ビルドされているのがわかる。
```
❯ dc logs -f
Attaching to 2_dlv_realize_app_1
... snip ...
app_1  | [06:32:05][APP-DEBUG] : GO changed /go/src/github.com/bellwood4486/sample-go-containerized-debug/2_dlv_realize/app/main.go
app_1  | [06:32:05][APP-DEBUG] : Command "pkill -INT /bin/sample-go-server-debug"
app_1  | [06:32:05][APP-DEBUG] : Build started
app_1  | [06:32:06][APP-DEBUG] : Build completed in 0.894 s
app_1  | [06:32:06][APP-DEBUG] : Running..
app_1  | [06:32:06][APP-DEBUG] : API server listening at: [::]:2345
app_1  | [06:32:06][APP-DEBUG] : Launching server at ":8080" ...
```

curlの結果も変わる
```
❯ curl http://localhost:8080
foo bar
```

再度プロセスを見てみると、dlvとsample-go-server-debugのPIDが変わっている。
```
/go/src/github.com/bellwood4486/sample-go-containerized-debug/2_dlv_realize # ps
PID   USER     TIME  COMMAND
    1 root      0:00 /go/bin/realize start
 1018 root      0:00 ash
 1091 root      0:00 /go/bin/dlv exec /bin/sample-go-server-debug --listen :2345 --headless --api-version 2 --continue --accept-multiclient
 1100 root      0:00 /bin/sample-go-server-debug
 1109 root      0:00 ps
/go/src/github.com/bellwood4486/sample-go-containerized-debug/2_dlv_realize # pstree -p
realize(1)---dlv(1091)---sample-go-serve(1100)
```
