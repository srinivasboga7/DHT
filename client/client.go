package client

import (
	"DHT/utils"
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"math/rand"
	"strconv"
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

// TestFunc tests the DHT network by inserting and retrieving random data
func TestFunc(kad *dht.IpfsDHT) error {
	randNums := rand.Perm(100)

	for i, n := range randNums {
		b := new(bytes.Buffer)
		binary.Write(b, binary.LittleEndian, n)
		val := b.Bytes()
		var key string
		key = strconv.Itoa(i)
		PutValue(kad, key, val)
	}

	for j, m := range randNums {
		key := strconv.Itoa(j)
		val, err := GetValue(kad, key)

		if err != nil {
			return err
		}

		var p int
		b := bytes.NewReader(val)
		binary.Read(b, binary.LittleEndian, &p)

		if m != p {
			return errors.New("mismatch values")
		}

	}

	return nil
}
