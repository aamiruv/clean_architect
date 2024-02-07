this is my clean architecture template for golang project.

### you need config.json file to run project. for example:
```go run cmd/web/main.go -config=/path/to/config_file.json```

### run tests:
```go test ./...```

#### [docker repository url](https://hub.docker.com/r/92276992/clean_architect)

#### you could run docker image easily via command for example:
```docker run --rm -p 8070:8070 -p 8071:8071  -it  --name my_container clean_architect:1.0```