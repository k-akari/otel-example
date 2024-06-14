# otel-example

## How to Run
### 1. Run Containers
```sh
$ docker compose up -d
```

### 2. Open Jaeger Dashboard in browser
`http://localhost:16686`

### 3. Open Prometheus Dashboard in browser
`http://localhost:9090`

### 4. Call Endpoints
```sh
# OK
$ curl -i -XGET localhost:9080/ -d '{"value": "hello", "sleep_time_ms": 10, "error_code_returned": 0}'

# INVALID ARGUMENT
$ curl -i -XGET localhost:9080/ -d '{"value": "hello", "sleep_time_ms": 10, "error_code_returned": 3}'
```
