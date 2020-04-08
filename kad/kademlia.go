package main

import (
	"DHT/utils"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/libp2p/go-libp2p"
	host "github.com/libp2p/go-libp2p-host"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	dhtopts "github.com/libp2p/go-libp2p-kad-dht/opts"
	peerstore "github.com/libp2p/go-libp2p-peerstore"
	multiaddr "github.com/multiformats/go-multiaddr"
)

func createHost(ctx context.Context, hostAddr multiaddr.Multiaddr) (*dht.IpfsDHT, host.Host) {

	// In the options add the privatekey
	host, err := libp2p.New(ctx,
		libp2p.ListenAddrs([]multiaddr.Multiaddr{hostAddr}...),
	)

	if err != nil {
		log.Fatal(err)
	}

	// add the DHT options
	kad, err := dht.New(ctx, host, dhtopts.Validator(utils.NullValidator{}))
	if err != nil {
		log.Fatal(err)
	}
	return kad, host
}

func addPeers(ctx context.Context, peersList []string, h host.Host, kad *dht.IpfsDHT) {

	for _, addr := range peersList {
		peerID, peerAddr := utils.MakePeer(addr)
		h.Peerstore().AddAddr(peerID, peerAddr, peerstore.PermanentAddrTTL)
		kad.Update(ctx, peerID)
	}

}

func main() {
	ctx := context.Background()
	port := os.Args[1]
	addr, err := utils.GenerateMultiAddr(port, "127.0.0.1")
	if err != nil {
		log.Fatal(err)
	}

	kademliaDHT, host := createHost(ctx, addr)

	log.Println("host address ", fmt.Sprintf("%s/p2p/%s", addr, host.ID().Pretty()))

	var peerList []string

	if len(os.Args) < 3 {
		log.Println("First node in the network")
	} else {
		peer := os.Args[2]
		peerList = append(peerList, peer)
		addPeers(ctx, peerList, host, kademliaDHT)
	}

	select {}
}
