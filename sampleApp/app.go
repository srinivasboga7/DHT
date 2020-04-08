package main

import (
	"DHT/client"
	"log"
	"net"
	"os"
)

func main() {
	// connect to discovery server to get an existing active node address

	conn, err := net.Dial("tcp", os.Args[1])
	if err != nil {
		log.Fatal("Failed to query discovery server", err)
	}

	buf := []byte{0x02}
	conn.Write(buf)

	resp := make([]byte, 128)
	Len, _ := conn.Read(resp)

	addr := string(resp[:Len])

	log.Println(addr)

	kad, cancel, err := client.NewDHTclient(addr)
	defer cancel()

	if err != nil {
		log.Fatal("Failed to intialize DHTclient", err)
	}

	err = client.TestFunc(kad)

	if err != nil {
		log.Println(err)
	} else {
		log.Println("DHT tested successfully")
	}
}
