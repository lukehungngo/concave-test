package main

import (
	"crypto/ecdsa"
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type Block struct {
	BlockData *BlockData  `json:"block_data"`
	Hash      common.Hash `json:"hash"`
	Signature []byte      `json:"signature"`
}

type BlockData struct {
	PreviousBlockHash common.Hash    `json:"previous_block_hash"`
	BlockNumber       uint64         `json:"block"`
	Nonce             uint64         `json:"nonce"`
	Data              interface{}    `json:"data"`
	ProducerAddress   common.Address `json:"producer_address"`
}

func NewBlockData(previousBlockHash common.Hash, blockNumber uint64, nonce uint64, data interface{}, producerAddress common.Address) *BlockData {
	return &BlockData{PreviousBlockHash: previousBlockHash, BlockNumber: blockNumber, Nonce: nonce, Data: data, ProducerAddress: producerAddress}
}

func CreateNewBlock(blockData *BlockData, privateKey *ecdsa.PrivateKey) (*Block, error) {
	res, err := json.Marshal(blockData)
	if err != nil {
		return nil, err
	}
	digestHash := crypto.Keccak256Hash(res)
	sig, err := crypto.Sign(digestHash[:], privateKey)
	if err != nil {
		return nil, err
	}
	return &Block{
		BlockData: blockData,
		Hash:      digestHash,
		Signature: sig,
	}, nil
}
