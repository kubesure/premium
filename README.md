# premium

apt-get install procps


kubectl create configmap --from-file=slave.conf=./slave.conf --from-file=master.conf=./master.conf --from-file=sentinel.conf=./sentinel.conf --from-file=init.sh=./init.sh --from-file=sentinel.sh=./sentinel.sh redis-config

k exec redis-premium-master-2 -c sentinel -- redis-cli -p 26379 sentinel get-master-addr-by-name redis-premium-master


sudo apt-get install jq

jq -n '{"code": "1A","sumInsured": "100000","dateOfBirth": "1990-06-07"}' | curl -H "Content-Type: application/json" -d@- http://172.17.0.8:8000/api/v1/healths/premiums | jq .

curl -i http://172.17.0.8:8000/api/v1/healths/premiums/loads

curl -i http://172.17.0.8:8000/api/v1/healths/premiums/unloads
