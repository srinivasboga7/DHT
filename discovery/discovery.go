package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"log"
	"math/rand"
	"net"
	"sync"
	"time"
)

type discvServer struct {
	mux      sync.Mutex
	DHTPeers []string
	maxPeers int
}

const (
	pingMsg = 0x01
)

// func (srv *discvServer) removePeer(peerAddr string) {
// 	srv.mux.Lock()
// 	delete(srv.DHTPeers, peerAddr)
// 	srv.mux.Unlock()
// }

func (srv *discvServer) addPeer(peerMultiAddr string) {
	srv.mux.Lock()
	srv.DHTPeers = append(srv.DHTPeers, peerMultiAddr)
	srv.mux.Unlock()
}

// func (srv *discvServer) pingPeer(peerAddr string, peerMultiAddr string) {
// 	conn, err := net.Dial("tcp", peerAddr)

// 	if err != nil {
// 		log.Println(err)
// 		return
// 	}

// 	for {
// 		time.Sleep(15 * time.Second)
// 		buf := []byte{pingMsg}
// 		_, err := conn.Write(buf)
// 		if err != nil {
// 			break
// 		}
// 	}
// 	srv.removePeer(peerAddr)
// 	return
// }

func (srv *discvServer) handleConnection(conn net.Conn) {

	// reading the message type
	buf1 := make([]byte, 1)
	conn.Read(buf1)

	if buf1[0] == 0x01 {
		// reading the length of the message
		buf2 := make([]byte, 4)
		conn.Read(buf2)
		r := bytes.NewReader(buf2)
		var l uint32
		binary.Read(r, binary.LittleEndian, &l)
		// reading the content of the message
		buf3 := make([]byte, l)
		conn.Read(buf3)
		multiaddr := string(buf3)

		var peerAddrs []string

		// selecting random peers
		srv.mux.Lock()
		if len(srv.DHTPeers) == 1 {
			peerAddrs = append(peerAddrs, srv.DHTPeers[0])
		} else if len(srv.DHTPeers) > 1 {
			randPerm := rand.Perm(len(srv.DHTPeers))
			peerAddrs = append(peerAddrs, srv.DHTPeers[randPerm[0]])
			peerAddrs = append(peerAddrs, srv.DHTPeers[randPerm[1]])
		}
		srv.mux.Unlock()
		s, _ := json.Marshal(peerAddrs)
		// replying with peers
		conn.Write(s)
		log.Println(multiaddr)
		srv.addPeer(multiaddr)

	} else if buf1[0] == 0x02 {
		var peerAddr string
		// selecting random peer
		log.Println("Client Request for DHT node address")
		srv.mux.Lock()
		if len(srv.DHTPeers) == 1 {
			peerAddr = srv.DHTPeers[0]
		} else if len(srv.DHTPeers) > 1 {
			randPerm := rand.Perm(len(srv.DHTPeers))
			peerAddr = srv.DHTPeers[randPerm[0]]
		}
		srv.mux.Unlock()

		s := []byte(peerAddr)
		// replying with peer
		conn.Write(s)
	}

}

func listenForNodes() {

	listener, err := net.Listen("tcp", ":8000")

	if err != nil {
		log.Fatal(err)
	}

	var srv discvServer

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go srv.handleConnection(conn)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	listenForNodes()
}
