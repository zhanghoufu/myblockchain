package main

import (
	"fmt"
)

func main() {
	fmt.Println("block chain")
	private, public := GetKeyPair()
	fmt.Printf("key pair \nprivatekey:%x\npublickey:%x\n", private, public)
	wal := NewWallet()
	wal.GetAddress()
}
