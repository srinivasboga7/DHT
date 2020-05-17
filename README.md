# DHT
A Distributed Hash Table forked from libp2p/go-libp2p-kad-dht.


# Usage

Build discovery server  
cd /discovery
go build discovery.go

Run discovery server
./discovery

Build kademlia DHT
cd /kad
go build kademlia.go

Run DHT node
./kademlia <port> <discovery address> <logging>
