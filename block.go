package main

import (
	"bytes"
	"crypto/sha256"
	"strconv"
	"time"
)

type Block struct {
	Timestamp int64
	PreHash   []byte
	Hash      []byte
	Data      []byte
}

type BlockChain struct {
	blocks []*Block
}

func NewBlock(data string, preblockchainhash []byte) *Block {
	var block = &Block{
		Timestamp: time.Now().Unix(),
		Data:      []byte(data),
		Hash:      []byte{},
		PreHash:   preblockchainhash}
	block.SetHash()
	return block

}

func (b *Block) SetHash() {
	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
	headers := bytes.Join([][]byte{b.Data, b.PreHash, timestamp}, []byte{})
	hash := sha256.Sum256(headers)
	b.Hash = hash[:]
}
