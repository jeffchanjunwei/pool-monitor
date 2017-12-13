# Introduction
pool-monitor is used to update the config file of cni plugin once the element in etcd changed.

# Installation instructions
- local compile
1. go build .
2. pool-monitor --etcdEndPoints https://xxxx:xxxx --caCert xxxx --key xxxx --cert xxxx

# Versions
- calicoctl v1.1.0
- calico v2.6
