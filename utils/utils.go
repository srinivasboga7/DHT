package utils

import (
	"fmt"
	"log"
	"math/rand"

	crypto "github.com/libp2p/go-libp2p-crypto"
	peer "github.com/libp2p/go-libp2p-peer"
	"github.com/multiformats/go-multiaddr"
)

// NullValidator is a validator that does no valiadtion
type NullValidator struct{}

// Validate always returns success
func (nv NullValidator) Validate(key string, value []byte) error {
	return nil
}

// Select always selects the first record
func (nv NullValidator) Select(key string, values [][]byte) (int, error) {
	strs := make([]string, len(values))
	for i := 0; i < len(values); i++ {
		strs[i] = string(values[i])
	}
	return 0, nil
}

// MakePeer converts an addtess in a string format to get peerID and multiaddr
func MakePeer(dest string) (peer.ID, multiaddr.Multiaddr) {
	ipfsAddr, err := multiaddr.NewMultiaddr(dest)
	if err != nil {
		log.Fatal(err)
	}
	peerIDStr, err := ipfsAddr.ValueForProtocol(multiaddr.P_P2P)
	if err != nil {
		log.Fatal(err)
	}
	peerID, err := peer.IDB58Decode(peerIDStr)
	if err != nil {
		log.Fatal(err)
	}
	targetPeerAddr, err := multiaddr.NewMultiaddr(fmt.Sprintf("/p2p/%s", peer.IDB58Encode(peerID)))
	if err != nil {
		log.Fatal(err)
	}
	targetAddr := ipfsAddr.Decapsulate(targetPeerAddr)
	return peerID, targetAddr
}

// GenerateMultiAddr creates a multiaddr from IP and port
func GenerateMultiAddr(port string, ip string) (multiaddr.Multiaddr, error) {
	return multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/%s/tcp/%s", ip, port))
}

// GeneratePrivateKey - creates a private key with the given seed
func GeneratePrivateKey(seed int64) crypto.PrivKey {
	randBytes := rand.New(rand.NewSource(seed))
	prvKey, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, randBytes)

	if err != nil {
		log.Fatalf("Could not generate Private Key: %v", err)
	}

	return prvKey
}
