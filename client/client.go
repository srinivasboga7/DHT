package client

import (
	"DHT/utils"
	"context"
	"time"

	"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	dhtopts "github.com/libp2p/go-libp2p-kad-dht/opts"
)

var (
	ctx context.Context
)

// NewDHTclient initializes and connects to destPeer DHT node
func NewDHTclient(destPeer string) (*dht.IpfsDHT, error) {
	ctx = context.Background()

	host, err := libp2p.New(ctx)

	if err != nil {
		return nil, err
	}

	destID, destAddr := utils.MakePeer(destPeer)
	host.Peerstore().AddAddr(destID, destAddr, 24*time.Hour)
	kademliaDHT, err := dht.New(ctx, host, dhtopts.Client(true), dhtopts.Validator(utils.NullValidator{}))

	if err != nil {
		return nil, err
	}

	kademliaDHT.Update(ctx, destID)
	return kademliaDHT, nil

}

// probably include a discovery protocol to discover the nodes in the network for
// client to connect and query

// PutValue inserts value into DHT
func PutValue(kad *dht.IpfsDHT, key string, value []byte) error {
	err := kad.PutValue(ctx, key, value)
	return err
}

// GetValue retrieves value corresponding to a key
func GetValue(kad *dht.IpfsDHT, key string) ([]byte, error) {
	val, err := kad.GetValue(ctx, key)
	return val, err
}
