Adb Backup
===================

通过ADB备份Android短信存储在数据库中，并使用webhook通知

工作原理
-----------------

使用`adb shell content query --uri content://sms/`定时读取手机短信

如果配置了通知地址，则每次读取到新短信时，会向通知地址发送POST请求

构建/运行
-----------------

```go
go mod download
// 构建二进制
// go build
// 直接运行
// go run
// 使用 .env 文件运行
go run -tags=dev .
```

Docker
-----------------

docker 镜像 `tablebird/adb-backup`

示例：[docker-compose.yml](./docker-compose.yml)

环境变量以及默认值
-----------------

- `ADB_PORT` adb server端口（默认 5037）
- `DB_HOST` postgres 数据库host （默认 postgres.lan）
- `DB_PORT` 数据库端口号 （默认 5432）
- `DB_USER` 数据库 （默认 backup）
- `DB_PASS` 数据库密码 （默认 backup）
- `DB_NAME` 书库库名称 （默认 backup）
- `DB_SSLMODE` 是否启动SSL （默认 disable）
- `DEBUG_LOG` 是否启动debug日志 （默认 false）
- `READ_INTERVAL` 读取消息得间隔 （默认 5s）
- `WAIT_DEVICE_INTERVAL` 等待Android设备连接的检测间隔 （默认 10）
- `NOTIFY_WEBHOOK_URL` 开始监听后新的短信的POST通知地址 （默认 空 不通知）
