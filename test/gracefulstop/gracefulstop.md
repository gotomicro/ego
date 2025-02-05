构造一个接口需要响应 10s
发送信号量 kill -s SIGTERM PID

curl http://127.0.0.1:9001/hello
kill -s SIGTERM 53065