package core

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
)

var (
	maxNonce = math.MaxInt64
)

// 挖矿的难度值
const targetBits = 20

// 工作量证明结构体
type ProofOfWork struct {
	block  *Block
	target *big.Int
}

// 新工作量证明
func NewProofOfWork(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits)) //左移
	pow := &ProofOfWork{b, target}
	return pow
}

// 准备数据
func (pow *ProofOfWork) prepareData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.block.PrevBlockHash,
			pow.block.Data,
			IntToHex(pow.block.Timestamp),
			IntToHex(int64(targetBits)), //目标值
			IntToHex(int64(nonce)),      //计数器
		}, []byte{})
	return data
}

// 工作量证明计算
func (pow *ProofOfWork) Run() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte
	nonce := 0

	fmt.Printf("Mining the block containing \"%s\"\n", pow.block.Data)
	for nonce < maxNonce {
		data := pow.prepareData(nonce) //准备数据

		hash = sha256.Sum256(data) //计算hash
		fmt.Printf("\r%x", hash)
		hashInt.SetBytes(hash[:]) //转成整数形式
		//比较hash和目标值 小于目标值就返回
		if hashInt.Cmp(pow.target) == -1 {
			break
		} else {
			nonce++
		}
	}
	fmt.Print("\n\n")
	return nonce, hash[:]
}

// 验证工作量证明
func (pow *ProofOfWork) Vaildate() bool {
	var hashInt big.Int
	data := pow.prepareData(pow.block.Nonce) //计数器
	hash := sha256.Sum256(data)              // 再次计算hash值
	hashInt.SetBytes(hash[:])
	isValid := hashInt.Cmp(pow.target) == -1 //验证工作量
	return isValid
}
