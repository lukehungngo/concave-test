package main

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"math/rand"
	"time"
)

type Account struct {
	PrivateKey *ecdsa.PrivateKey
	Address    common.Address
}
type BlockProducer struct {
	lastBlock       *Block
	blockRepository BlockRepository
	producerAccount *Account
	isStart         bool
	dataChan        chan interface{}
}

func Init(producerAccount *Account, dataChan chan interface{}, blockRepository BlockRepository) (*BlockProducer, error) {
	blockProducer := &BlockProducer{
		producerAccount: producerAccount,
		blockRepository: blockRepository,
		dataChan:        dataChan,
		isStart:         false,
	}
	lastBlock, _, err := blockProducer.blockRepository.GetLastBlock()
	if err != nil {
		return nil, err
	}
	blockProducer.lastBlock = lastBlock

	return blockProducer, nil
}

func (bp *BlockProducer) Run(done chan bool) {
	if bp.isStart {
		fmt.Println("Block Producer already started")
	}

	if bp == nil {
		panic("Block Producer is not init yet")
	}
	bp.isStart = true
	defer func() {
		bp.isStart = false
	}()
	fmt.Println("Block Producer is started")
	for {
		select {
		case <-done:
			close(bp.dataChan)
			fmt.Println("Block Producer is stopped")
			return
		case data := <-bp.dataChan:
			blockNumber := uint64(0)
			previousLastBlock := common.Hash{}
			if bp.lastBlock != nil {
				blockNumber = bp.lastBlock.BlockData.BlockNumber + 1
				previousLastBlock = bp.lastBlock.Hash
			}
			newBlockData := NewBlockData(
				previousLastBlock,
				blockNumber,
				uint64(rand.Uint32()),
				data.(string),
				bp.producerAccount.Address,
			)
			newBlock, err := CreateNewBlock(newBlockData, bp.producerAccount.PrivateKey)
			if err != nil {
				fmt.Printf("Block Creating Error: blockNumber=%d - data=%+v - error=%+v\n", blockNumber, data, err)
				continue
			}
			if err := bp.blockRepository.StoreBlock(newBlock); err != nil {
				fmt.Printf("Block Store Error: blockNumber=%d - data=%+v - error=%+v\n", blockNumber, data, err)
				continue
			}
			bp.lastBlock = newBlock
			fmt.Printf("BLOCK CREATED | number=%d - hash=%+v\n", newBlock.BlockData.BlockNumber, newBlock.Hash)
		default:
			time.Sleep(100 * time.Millisecond)
		}
	}
}
