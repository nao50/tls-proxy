# tcp proxy

## Server
```bash
$ go run 0_tcp_server/main.go
```

## Proxy
```bash
$ go run 1_tcp_proxy/main.go
```

## Client
```bash
$ go run tcp_client.go -remoteAddress localhost:8080
```