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

ハンドラ内に貼ったブレイクポイントをデバッグできる。
![スクリーンショット 2020-04-18 12 35 52](https://user-images.githubusercontent.com/2452581/79627399-d64f5680-8172-11ea-8d56-96c01efd0882.png)

注意：main関数をデバッグしたい…TODO


