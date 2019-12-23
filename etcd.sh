export PUBLIC_IP=$(ipconfig getifaddr en0)
docker run -d -p 8001:8001 -p 5001:5001 quay.io/coreos/etcd:v0.4.6 -peer-addr ${PUBLIC_IP}:8001 -addr ${PUBLIC_IP}:5001 -name etcd-node1
curl -L $PUBLIC_IP:5001/v2/stats/leader