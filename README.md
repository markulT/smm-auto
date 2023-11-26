### Default

Setup required docker containers

```bash
cd ./docker
docker-compose up
```
And then simply run the main.go
```bash
go run main.go
```

### Only swagger
This version only includes swagger docs

To run it , you have to simply run \cmd\swagger\main.go without docker containers

```bash
go run .\cmd\swagger\main.go
```