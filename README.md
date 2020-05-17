# DHT
A Distributed Hash Table forked from libp2p/go-libp2p-kad-dht.


## Pre-requisites
go version 1.13. \
Dependencies will be automatically downloaded by go.mod, no need to setup dependencies.

## Usage

### Discovery Node
discovery dir has the code for discovery server.\
Run discovery node using : go run discovery.go.\
Discovery Node runs on the port 8000.\
Note : Run the discovery server before any of the DHT nodes.\

### DHT Node
kad dir has the code for the DHT node.\
Building : go build kademlia.go .\
Running : ./kademlia <port> <discovery Node Address> <logging> . \
For example : ./kademlia 6666 127.0.0.1:8000 enableLogging , this command runs the node on port 6666 with logging enabled.
