package main

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
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

var blockProducer *BlockProducer

func Init(producerAccount *Account, dataChan chan interface{}, blockRepository BlockRepository) error {
	blockProducer = &BlockProducer{
		producerAccount: producerAccount,
		blockRepository: blockRepository,
		dataChan:        dataChan,
		isStart:         false,
	}
	lastBlock, _, err := blockProducer.blockRepository.GetLastBlock()
	if err != nil {
		return err
	}
	blockProducer.lastBlock = lastBlock

	return nil
}

func Run() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	done := make(chan bool, 1)

	go func() {
		<-sigs
		close(done)
		return
	}()

	if blockProducer == nil {
		panic("Block Producer is not init yet")
	}

	for {
		select {
		case <-done:
			return
		case data := <-blockProducer.dataChan:
			blockNumber := uint64(0)
			previousLastBlock := common.Hash{}
			if blockProducer.lastBlock != nil {
				blockNumber = blockProducer.lastBlock.BlockData.BlockNumber + 1
				previousLastBlock = blockProducer.lastBlock.Hash
			}
			newBlockData := NewBlockData(
				previousLastBlock,
				blockNumber,
				rand.Uint64(),
				data,
				blockProducer.producerAccount.Address,
			)
			newBlock, err := CreateNewBlock(newBlockData, blockProducer.producerAccount.PrivateKey)
			if err != nil {
				fmt.Printf("Block Creating Error: blockNumber=%d - data=%+v - error=%+v", blockNumber, data, err)
			}
			if err := blockProducer.blockRepository.StoreBlock(newBlock); err != nil {
				fmt.Printf("Block Store Error: blockNumber=%d - data=%+v - error=%+v", blockNumber, data, err)

			}
		default:
			time.Sleep(100 * time.Millisecond)
		}
	}
}
