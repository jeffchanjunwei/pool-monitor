# Installation instructions
- local compile
step1 go build .
step2 ./pool-monitor --etcdEndPoints https://10.142.21.224:2379 --caCert /etc/cni/net.d/calico-tls/etcd-ca --key /etc/cni/net.d/calico-tls/etcd-key --cert /etc/cni/net.d/calico-tls/etcd-cert
