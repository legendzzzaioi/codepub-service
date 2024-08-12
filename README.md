# codepub-service
This backend project provides a solution for code publish using the Gin framework.

### Usage

Build Docker Image

```
docker build -t legendzzzaioi/codepub-service:v1 .
```

Run Docker Container

```
# mariadb
# docker run -d \
#   --restart always \
#   -p 3306:3306 \
#   -v /data/mysql:/var/lib/mysql \
#   --name mariadb \
#   --env MARIADB_ROOT_PASSWORD=root \
#   -d mariadb:latest

# modify config.yaml, then
docker run -d \
  --restart always \
  -p 8000:8000 \
  --name codepub-service \
  -d legendzzzaioi/codepub-service:v1
```

Usage with Kubernetes

```
# mariadb
# kubectl -n xxx apply -f mariadb.yaml

# modify codepub-service.yaml
kubectl -n xxx apply -f codepub-service.yaml
```
