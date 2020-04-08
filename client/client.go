package client

import (
	"DHT/utils"
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	dhtopts "github.com/libp2p/go-libp2p-kad-dht/opts"
)

var (
	ctx    context.Context
	cancel context.CancelFunc
)

// NewDHTclient initializes and connects to destPeer DHT node
func NewDHTclient(destPeer string) (*dht.IpfsDHT, context.CancelFunc, error) {
	ctx, cancel = context.WithCancel(context.Background())

	host, err := libp2p.New(ctx, libp2p.Identity(utils.GeneratePrivateKey(time.Now().Unix())))

	if err != nil {
		return nil, cancel, err
	}

	destID, destAddr := utils.MakePeer(destPeer)
	host.Peerstore().AddAddr(destID, destAddr, 24*time.Hour)
	kademliaDHT, err := dht.New(ctx, host, dhtopts.Client(true), dhtopts.Validator(utils.NullValidator{}))

	if err != nil {
		return nil, cancel, err
	}

	kademliaDHT.Update(ctx, destID)

	err = kademliaDHT.Ping(ctx, destID)
	if err != nil {
		log.Println(err)
	}

	return kademliaDHT, cancel, nil
}

// probably include a discovery protocol to discover the nodes in the network for
// client to connect and query

// PutValue inserts value into DHT
func PutValue(kad *dht.IpfsDHT, key string, value []byte) error {
	err := kad.PutValue(ctx, key, value, dht.Quorum(1))
	return err
}

// GetValue retrieves value corresponding to a key
func GetValue(kad *dht.IpfsDHT, key string) ([]byte, error) {
	val, err := kad.GetValue(ctx, key, dht.Quorum(1))
	return val, err
}

// TestFunc tests the DHT network by inserting and retrieving random data
func TestFunc(kad *dht.IpfsDHT) error {
	randNums := rand.Perm(1)

	for i, n := range randNums {
		un := uint32(n)
		b := new(bytes.Buffer)
		binary.Write(b, binary.LittleEndian, un)
		val := b.Bytes()
		log.Println(val)
		var key string
		key = strconv.Itoa(i)
		err := PutValue(kad, key, val)
		if err != nil {
			return err
		}
	}

	for j, m := range randNums {
		key := strconv.Itoa(j)
		val, err := GetValue(kad, key)
		if err != nil {
			return err
		}
		var p uint32
		b := bytes.NewReader(val)
		binary.Read(b, binary.LittleEndian, &p)
		if uint32(m) != p {
			return errors.New("mismatch values")
		}
	}

	return nil
}
