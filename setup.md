## Develop

Clone repo in go/src/github.com/kubesure #todo #defect 18 to convert to go modules  

### start redis. conf folder contains configuration 

# update redis data path in master.conf
```
redis-server conf/master.conf 
redis-server conf/slave.conf
redis-server conf/sentinel.conf --sentinel
```

### run premium calc service. install Go 1.15 and above 
export redissvc=localhost
go run health.go

### load premium matrix in redis

```
curl -i http://<hostname>:8000/api/v1/healths/premiums/loads
```

### calculate premium
```

curl -X POST http://<hostname>:8000/api/v1/healths/premiums \
-H "Content-Type: application/json" \
-d '{"code": "1A","sumInsured": "100000", "dateOfBirth": "1990-06-07"}' | jq .
```

### unload premium matrix. reload excel based matrix

```
curl -i http://<hostname>:8000/api/v1/healths/premiums/unloads
```

### deploy and run in Kind. Create Kind with Nginx Ingress Controller (https://kind.sigs.k8s.io/docs/user/ingress/#using-ingress)

```
git clone https://github.com/kubesure/helm-charts.git
cd helm-charts/premium
helm install premium .

kubectl run curl --image=radial/busyboxplus --restart=Never -it -- sh

## loads api not to be executed as its done during container startup

curl -i -X POST http://<ip of premiumcalc pod>:8000/api/v1/healths/premiums \
-H "Content-Type: application/json" \
-d '{"code": "1A","sumInsured": "100000", "dateOfBirth": "1990-06-07"}' 

cd helm-charts/ingress-egress 
helm install kubesure .
kubectl get ingress

curl -i -X POST http://localhost/api/v1/healths/premiums \
-H "Content-Type: application/json" \
-d '{"code": "1A","sumInsured": "100000", "dateOfBirth": "1990-06-07"}'
```