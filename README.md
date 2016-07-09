
# go-keymaker

Distributed generation of RSA and ECDSA private keys.  Used when a process needs
a large number of keys generated securely, without overloading the local host.
A pilot process contacts one or more key generation drones using a simple RPC
protocol.  The pilot uses buffered channels and asynchronous RPC to reduce
latency.  Channels are created and kept full by background goroutines.  Private
keys are PEM-encoded PKCS#8.

Primary goal is to keep the system lightweight and simple.  A YAML configuration
file defines the topology of the system but dynamic discovery of drones using
Consul (ZK, Etcd) can be added if desired.
