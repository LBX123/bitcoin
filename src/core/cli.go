package core

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
)

type CLI struct {
	Bc *Blockchain
}

func (cli *CLI) printUsage() {
	fmt.Println("Usage")
	fmt.Println("addblock -data BLOCK_DATA - add a block to the blockchain")
	fmt.Println("printchain - print all the blocks of the blockchain")
}

// 验证参数个数
func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}
}

func (cli *CLI) addBlock(data string) {
	cli.Bc.AddBlock(data)
	fmt.Println("Success!")
}

func (cli *CLI) printChain() {
	bci := cli.Bc.Iterator()
	for {
		block := bci.Next() //对区块进行迭代
		fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		pow := NewProofOfWork(block)                                //生成工作量证明
		fmt.Printf("Pow: %s\n", strconv.FormatBool(pow.Vaildate())) //验证工作量
		fmt.Println()

		if len(block.PrevBlockHash) == 0 { //迭代完毕
			break
		}
	}
}

func (cli *CLI) Run() {
	cli.validateArgs()
	//解析命令行参数
	addBlockCmd := flag.NewFlagSet("addblock", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	addBlockData := addBlockCmd.String("data", "", "Block data")

	switch os.Args[1] {
	case "addblock":
		//解析子命令
		err := addBlockCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "printchain":
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		cli.printUsage()
		os.Exit(1)
	}
	if addBlockCmd.Parsed() {
		if *addBlockData == "" {
			addBlockCmd.Usage()
			os.Exit(1)
		}
		//添加区块
		cli.addBlock(*addBlockData)
	}
	if printChainCmd.Parsed() {
		cli.printChain()
	}
}
