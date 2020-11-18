## Develop

Clone repo in go/src/github.com/kubesure #todo #defect 18 to convert to go modules  

### start redis. conf folder contains configuration

```
redis-server conf/master.conf
redis-server conf/slave.conf
redis-server conf/sentinel.conf --sentinel
```

### run premium calc service. install Go 1.15 and above 
export redissvc=localhost
go run health.go

## load premium matrix in redis

```
curl -i http://<hostname>:8000/api/v1/healths/premiums/loads
```

### calculate premium
```

curl -i -X POST http://<hostname>:8000/api/v1/healths/premiums \
-H "Content-Type: application/json" \
-d '{"code": "1A","sumInsured": "100000", "dateOfBirth": "1990-06-07"}' | jq .
```

## unload premium matrix. reload excel based matrix

```
curl -i http://<hostname>:8000/api/v1/healths/premiums/unloads
```
