# go-syslog-server 使用说明
这是一个Go语言写的迷你的SYSLOG服务器

- 监听在*:514[UDP]端口,即所有可用接口的514接口
- 将接收到的日志写入 `./logs/<SRC-IP>/YYYY-MM-DD.log` 文件
- 同时会将日志打印到屏幕上cd