package main

import (
	"encoding/binary"
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

type BlockRepository interface {
	StoreBlock(block *Block) error
	GetLastBlock() (*Block, bool, error)
	GetBlocks(limit, offset int) ([]*Block, error)
	GetBlockByNumber(blockNumber uint64) (*Block, bool, error)
}

const (
	// Define our schema here
	LAST_BLOCK_KEY      = "LAST_BLOCK"
	BLOCK_NUMBER_PREFIX = "BLOCK_NUMBER_"
	BLOCK_HASH_PREFIX   = "BLOCK_HASH_"
)

type LevelDbBlockRepository struct {
	db *leveldb.DB
}

func NewLevelDbBlockRepository(db *leveldb.DB) *LevelDbBlockRepository {
	return &LevelDbBlockRepository{db: db}
}

func (l LevelDbBlockRepository) StoreBlock(block *Block) error {
	batch := new(leveldb.Batch)

	lastBlockKey := []byte(LAST_BLOCK_KEY)
	lastBlockValue := block.Hash[:]
	batch.Put(lastBlockKey, lastBlockValue)

	blockNumberKey := append([]byte(BLOCK_NUMBER_PREFIX), convertUint64ToBytes(block.BlockData.BlockNumber)...)
	blockNumberValue := block.Hash[:]
	batch.Put(blockNumberKey, blockNumberValue)

	blockHashKey := append([]byte(BLOCK_HASH_PREFIX), block.Hash[:]...)
	blockHashValue, err := json.Marshal(block)
	if err != nil {
		return err
	}
	batch.Put(blockHashKey, blockHashValue)

	return l.db.Write(batch, nil)

}

func (l LevelDbBlockRepository) GetLastBlock() (*Block, bool, error) {
	lastBlockKey := []byte(LAST_BLOCK_KEY)
	if has, err := l.db.Has(lastBlockKey, nil); err != nil {
		return nil, false, err
	} else if !has {
		return nil, false, nil
	}
	blockHash, err := l.db.Get(lastBlockKey, nil)
	if err != nil {
		return nil, false, err
	}
	blockHashKey := append([]byte(BLOCK_HASH_PREFIX), blockHash[:]...)
	blockData, err := l.db.Get(blockHashKey, nil)
	block := &Block{}
	if err := json.Unmarshal(blockData, block); err != nil {
		return nil, false, err
	}

	return block, true, err
}

func (l LevelDbBlockRepository) GetBlocks(limit, offset int) ([]*Block, error) {
	blocks := []*Block{}
	prefix := []byte(BLOCK_NUMBER_PREFIX)
	iter := l.db.NewIterator(util.BytesPrefix(prefix), nil)
	defer iter.Release()

	count := 0
	for iter.Next() {
		// Skip elements until the offset is reached
		if count < offset {
			count++
			continue
		}

		// Process the key-value pair
		value := copyByteArray(iter.Value())
		blockHash := common.BytesToHash(value)
		blockHashKey := append([]byte(BLOCK_HASH_PREFIX), blockHash[:]...)
		blockData, err := l.db.Get(blockHashKey, nil)
		if err != nil {
			return nil, err
		}
		block := &Block{}
		if err := json.Unmarshal(blockData, block); err != nil {
			return nil, err
		}
		blocks = append(blocks, block)
		// Break the loop once the limit is reached
		count++
		if count >= offset+limit {
			break
		}
	}

	if err := iter.Error(); err != nil {
		return nil, err
	}
	return blocks, nil
}

func (l LevelDbBlockRepository) GetBlockByNumber(blockNumber uint64) (*Block, bool, error) {
	blockNumberKey := append([]byte(BLOCK_NUMBER_PREFIX), convertUint64ToBytes(blockNumber)...)
	if has, err := l.db.Has(blockNumberKey, nil); err != nil {
		return nil, false, err
	} else if !has {
		return nil, false, nil
	}

	blockHashValue, err := l.db.Get(blockNumberKey, nil)
	if err != nil {
		return nil, false, err

	}
	blockHashKey := append([]byte(BLOCK_HASH_PREFIX), blockHashValue[:]...)
	blockData, err := l.db.Get(blockHashKey, nil)
	if err != nil {
		return nil, false, err
	}
	block := &Block{}
	if err := json.Unmarshal(blockData, block); err != nil {
		return nil, false, err
	}

	return block, true, nil
}

func convertUint64ToBytes(number uint64) []byte {
	bytes := make([]byte, 8)

	binary.BigEndian.PutUint64(bytes, number)

	return bytes
}

func copyByteArray(src []byte) []byte {
	dst := make([]byte, len(src))
	copy(dst, src)
	return dst
}
