package main

import (
	"DHT/utils"
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

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

	if len(peersList) == 0 {
		return
	}

	for _, addr := range peersList {
		peerID, peerAddr := utils.MakePeer(addr)
		h.Peerstore().AddAddr(peerID, peerAddr, peerstore.PermanentAddrTTL)
		kad.Update(ctx, peerID)
	}

}

func main() {
	ctx := context.Background()
	port := os.Args[1]

	// contact discovery server
	conn, err := net.Dial("tcp", os.Args[2])
	if err != nil {
		log.Fatal("Failed to query discovery server", err)
	}

	ipaddr := conn.LocalAddr().String()
	ipaddr = ipaddr[:strings.IndexByte(ipaddr, ':')]
	addr, err := utils.GenerateMultiAddr(port, ipaddr)

	kad, host := createHost(ctx, addr)
	hostAddr := fmt.Sprintf("%s/p2p/%s", addr, host.ID().Pretty())
	log.Println(hostAddr)

	buf := []byte{0x01}
	payload := []byte(hostAddr)
	var l uint32
	l = uint32(len(payload))

	b := new(bytes.Buffer)
	binary.Write(b, binary.LittleEndian, l)
	buf = append(buf, b.Bytes()...)
	buf = append(buf, payload...)

	conn.Write(buf)

	resp := make([]byte, 1024)
	Len, _ := conn.Read(resp)

	// decoding the list of peers
	var peerAddr []string
	json.Unmarshal(resp[:Len], &peerAddr)

	log.Println(peerAddr)
	// connecting with peers
	addPeers(ctx, peerAddr, host, kad)

	select {}
}
