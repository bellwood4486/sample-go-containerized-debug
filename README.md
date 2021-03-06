# sample-go-containerized-debug

## 0. Code Base

### How to run
```
❯ cd 0_base
❯ docker-compose up -d --build
❯ curl http://localhost:8080
hello world
```

## 1. Debugging with Delve

### How to run
```
❯ cd 1_dlv
❯ docker-compose up -d --build
```

The port 8080 is for the application. 2345 is for Delve.
```
❯ docker-compose logs -f
Attaching to 1_dlv_app_1
app_1  | API server listening at: [::]:2345
app_1  | Launching server at ":8080" ...
```

The application returns "hello world".
```
❯ curl http://localhost:8080
hello world
```

In the container, there are 2 processes. PID:1 is a process of Delve. PID:12 is of the application launched by Delve.
```
❯ docker-compose exec app ash
/go/src/github.com/bellwood4486/sample-go-containerized-debug/1_dlv # ps
PID   USER     TIME  COMMAND
    1 root      0:00 /go/bin/dlv exec /bin/sample-go-server-debug --listen=:2345 --headless --api-version=2 --continue --accept-multiclient
   12 root      0:00 /bin/sample-go-server-debug
```

In the tree view, you can see that the application is being launched by Delve.
```
/go/src/github.com/bellwood4486/sample-go-containerized-debug/1_dlv # pstree -p
dlv(1)---sample-go-serve(12)
```

### How to debug on GoLand

Add a configuration of "Go Remote".

![スクリーンショット 2020-04-18 12 27 27](https://user-images.githubusercontent.com/2452581/79627378-9b4d2300-8172-11ea-93a7-d095e95b086d.png)

Specify the number defined in the Dockerfile as the port. (2345 in this case)

![スクリーンショット 2020-04-18 12 29 14](https://user-images.githubusercontent.com/2452581/79627388-aef88980-8172-11ea-9a4a-7aaf5b72e3ab.png)

When the debugger attaches, the following message is displayed. The cause has not been investigated.
```
❯ docker-compose logs -f
Attaching to 1_dlv_app_1
... snip ...
app_1  | 2020-04-18T07:03:47Z error layer=rpc writing response:write tcp 127.0.0.1:2345->127.0.0.1:48800: use of closed network connection
```

Stops at the break point.

![スクリーンショット 2020-04-18 12 35 52](https://user-images.githubusercontent.com/2452581/79627399-d64f5680-8172-11ea-8d56-96c01efd0882.png)

This sample can't debug the main function. This is because the main function is finished before debugger attaches to it. To debug it, change the option of Delve.
see: https://github.com/derekparker/delve/blob/master/Documentation/faq.md#how-do-i-use-delve-with-docker

## 2. Debugging a hot-reload app with Delve

### How to run
```
❯ cd 2_dlv_realize
❯ docker-compose up -d --build
❯ curl http://localhost:8080
hello world
```

In the following logs, Realize has completed the build in 7.766 seconds.
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

In the container, there are 3 process. They are Realize, Delve and the sample app.
```
❯ docker-compose exec app ash
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

Modify the code in main.go to confirm the hot-reload. ("hello reload" -> "foo bar")

main.go
```go
//... snip ...
		_, _ = fmt.Fprintf(w, "foo bar")
//... snip ...
```

You can see from the log that it has been rebuilt.
```
❯ docker-compose logs -f
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

The result of curl also changes.
```
❯ curl http://localhost:8080
foo bar
```

The PID of Delve and the app has changed.
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
