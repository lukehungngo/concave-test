package main

import (
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"strconv"
)

type Block struct {
	Hash      common.Hash `json:"hash"`
	Signature []byte      `json:"signature"`
	BlockData *BlockData  `json:"block_data"`
}

type BlockData struct {
	PreviousBlockHash common.Hash    `json:"previous_block_hash"`
	BlockNumber       uint64         `json:"block"`
	Nonce             uint64         `json:"nonce"`
	ProducerAddress   common.Address `json:"producer_address"`
	Data              string         `json:"data"`
}

func (b BlockData) toBytes() ([]byte, error) {
	var res []byte
	res = append(res, b.PreviousBlockHash[:]...)
	blockNumberBytes := common.FromHex(strconv.FormatUint(b.BlockNumber, 16))
	res = append(res, common.LeftPadBytes(blockNumberBytes, 32)...)
	nonceBytes := common.FromHex(strconv.FormatUint(b.Nonce, 16))
	res = append(res, common.LeftPadBytes(nonceBytes, 32)...)
	res = append(res, common.LeftPadBytes(b.ProducerAddress.Bytes(), 32)...)
	res = append(res, []byte(b.Data)...)
	return res, nil
}

func NewBlockData(previousBlockHash common.Hash, blockNumber uint64, nonce uint64, data string, producerAddress common.Address) *BlockData {
	return &BlockData{PreviousBlockHash: previousBlockHash, BlockNumber: blockNumber, Nonce: nonce, Data: data, ProducerAddress: producerAddress}
}

func (b Block) toBytes() ([]byte, error) {
	var res []byte
	res = append(res, b.Hash[:]...)
	//signatureBytes := make([]byte, 32*3)
	//copy(signatureBytes, b.Signature)
	//res = append(res, signatureBytes...)
	res = append(res, b.Signature...)
	blockDataBytes, err := b.BlockData.toBytes()
	if err != nil {
		return nil, err
	}
	res = append(res, blockDataBytes...)
	return res, nil
}

func CreateNewBlock(blockData *BlockData, privateKey *ecdsa.PrivateKey) (*Block, error) {
	res, err := blockData.toBytes()
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
