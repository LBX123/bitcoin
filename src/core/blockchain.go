package core

import (
	"fmt"
	"github.com/boltdb/bolt"
	"log"
)

const dbFile = "blockchain.db"
const blockBucket = "blocks"

type Blockchain struct {
	tip []byte
	Db  *bolt.DB
}
type BlockchainIterator struct {
	currentHash []byte
	Db          *bolt.DB
}

func NewGenesisBlock() *Block {
	return NewBlock("Genesis Block", []byte{})
}

// 区块链迭代器 从尾到头
func (bc *Blockchain) Iterator() *BlockchainIterator {
	bci := &BlockchainIterator{bc.tip, bc.Db}
	return bci
}
func (i *BlockchainIterator) Next() *Block {
	var block *Block
	//查找当前末尾区块
	err := i.Db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockBucket))
		encodeBlock := b.Get(i.currentHash)
		//反序列化到内存中
		block = DeserializeBlock(encodeBlock)
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	//更细末尾区块hash
	i.currentHash = block.PrevBlockHash

	return block
}

// 增加区块
func (bc *Blockchain) AddBlock(data string) {
	var lashHash []byte
	//只读区块
	err := bc.Db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockBucket))
		lashHash = b.Get([]byte("1"))

		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	newBlock := NewBlock(data, lashHash) //挖出一个新区块
	err = bc.Db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockBucket))
		err := b.Put(newBlock.Hash, newBlock.Serialize())
		if err != nil {
			log.Panic(err)
		}
		err = b.Put([]byte("1"), newBlock.Hash) //更新最后一个区块的hash
		if err != nil {
			log.Panic(err)
		}
		bc.tip = newBlock.Hash

		return nil
	})
}

func NewBlockchain() *Blockchain {
	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil) //打开BoltDB文件
	if err != nil {
		log.Panic(err)
	}
	//读写事务
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockBucket)) //读取bucket
		if b == nil {
			fmt.Println("No existing blockchain found. Creating a new one...")
			genesis := NewGenesisBlock()                   //创世区块
			b, err := tx.CreateBucket([]byte(blockBucket)) //创建bucket
			if err != nil {
				log.Panic(err)
			}
			err = b.Put(genesis.Hash, genesis.Serialize()) //持久化
			if err != nil {
				log.Panic(err)
			}
			err = b.Put([]byte("1"), genesis.Hash)
			if err != nil {
				log.Panic(err)
			}
			tip = genesis.Hash //指向最后一个块的hash
		} else { //如果存在就直接获取
			tip = b.Get([]byte("1"))
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	//创建区块链
	bc := Blockchain{tip, db}
	return &bc
}
