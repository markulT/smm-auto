## Default

Setup required docker containers


```bash
cd ./docker
docker-compose -p smm-system up
```
**!Setting the name of the project is necessarily**

And then simply run the main.go

```bash
go run main.go
```

## Only swagger
This version only includes swagger docs

To run it , you have to simply run \cmd\swagger\main.go without docker containers

```bash
go run .\cmd\swagger\main.go
```

### How to access

1. To access docs you can either go to html version (/swagger/index.html) or by using .json or .yaml files