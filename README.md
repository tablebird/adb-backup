# Adb Backup

[中文说明](./README_zh.md)

Backup Android SMS via ADB and store in database with webhook notification

## How it works

Uses `adb shell content query --uri content://sms/` to periodically read SMS messages from the phone.

If a notification address is configured, a POST request will be sent to the notification address every time a new SMS is read

## Build/Run

```go
cp example.env .env

go mod download
// build
// go build
// run
// go run
// use .env file
go run -tags=dev .
```

## Docker

docker images `tablebird/adb-backup`

example: [docker-compose.yml](./docker-compose.yml)

## Environment Variables and Default Values

- `ADB_HOST`: adb server host (default: localhost), Prioritize connecting to adb services for this host address
- `ADB_PORT`: adb server port (default: 5037)
- `DB_HOST`: PostgreSQL database host (default: postgres.lan)
- `DB_PORT`: Database port number (default: 5432)
- `DB_USER`: Database user (default: backup)
- `DB_PASS`: Database password (default: backup)
- `DB_NAME`: Database name (default: backup)
- `DB_SSLMODE`: Whether to enable SSL (default: disable)
- `DEBUG_LOG`: Whether to enable debug logging (default: false)
- `READ_INTERVAL`: Message reading interval (default: 5s)
- `WAIT_DEVICE_INTERVAL`: Interval for checking Android device connection (default: 10)
- `NOTIFY_WEBHOOK_URL`: POST notification URL for new SMS after starting monitoring (default: empty, no notification)
- `NOTIFY_STATUS_WEBHOOK_URL`: POST notification URL for device status after starting monitoring (default: empty, no notification)
- `WEB_ADDRESS`: Web service address (default: all/0.0.0.0)
- `WEB_PORT`: Web service port (default: 8080)
- `ADMIN_NAME`: Initial administrator username (default: admin)
- `ADMIN_PASS`: Initial administrator password (default: admin)
- `LDAP_HOST`: LDAP server host (default: empty, no LDAP authentication)
- `LDAP_PORT`: LDAP server port (default: 389)
- `LDAP_ENABLE_TLS` : Whether to enable TLS (default: false)
- `LDAP_BIND_DN`: LDAP bind DN (default: empty, no LDAP authentication)
- `LDAP_BIND_PASS`: LDAP bind password (default: empty, no LDAP authentication)
- `LDAP_BASE_DN`: LDAP base DN (default: empty)
- `LDAP_USER_FILTER`: LDAP user filter (default: empty)
