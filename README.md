## Docker 运行


```shell
docker run --restart=always --name monitor_logger -p 9081:80 -v /Logs:/Logs -d sqzxcv/monitor_logger:latest
```
**-v /:/** 限定容器能访问数组的文件目录, 可以指定全部映射, 也可以只映射日志目录所在位置


```shell
docker rm -f monitor_logger && docker run --restart=always --name monitor_logger -p 9081:80 -v /Logs:/Logs -d sqzxcv/monitor_logger:latest
```
强制更新容器