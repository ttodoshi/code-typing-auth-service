# Code Typing

## Auth Service on Golang

Auth service written on Golang with MongoDB and Gin framework

Swagger URL: [http://localhost:8090/swagger-ui/index.html](http://localhost:8090/swagger-ui/index.html)

### Run locally

```shell
make run
```

### Build docker container

```shell
docker buildx build . -t ghcr.io/ttodoshi/code-typing-auth-service:latest
```
