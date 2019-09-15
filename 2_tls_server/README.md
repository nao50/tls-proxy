# tls proxy

## Generate cert
```bash
$ cd 2_tls_server/certificate
$ go run ca.go
$ go run server.go
$ go run client.go
```

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