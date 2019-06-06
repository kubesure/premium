## Install Dev

```
redis-server ./master-dev.conf
redis-server ./slave-dev.conf
redis-server ./sentinel-dev.conf --sentinel
```

## Install k8s
```
kubectl apply -f health/config/premium.yaml

kubectl exec redis-premium-master-2 -c sentinel -- redis-cli -p 26379 sentinel get-master-addr-by-name redis-premium-master

sudo apt-get install jq

kubectl get po -o wide
```

## load premium matrix in redis

```
curl -i http://<ip of premiumcalc pod>:8000/api/v1/healths/premiums/loads
```

## calculate premium
```
jq -n '{"code": "1A","sumInsured": "100000","dateOfBirth": "1990-06-07"}' | \
curl -H "Content-Type: application/json" -d@- http://<ip of premiumcalc pod>:8000/api/v1/healths/premiums | jq .

 curl -i -X POST http://<ip of premiumcalc pod>:8000/api/v1/healths/premiums -H "Content-Type: application/json" \
 -d '{"code": "1A","sumInsured": "100000", "dateOfBirth": "1990-06-07"}' 

 ingress

 curl -i -X POST http://<ingress lb ip>/api/v1/healths/premiums -H "Content-Type: application/json" -d '{"code": "1A","sumInsured": "100000", "dateOfBirth": "1990-06-07"}'
 ```

## unload premium matrix

```
curl -i http://<ip of premiumcalc pod>:8000/api/v1/healths/premiums/unloads
```