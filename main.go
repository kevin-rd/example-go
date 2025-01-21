package main

import (
	"context"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	client, err := ethclient.Dial("http://testnet-rpc.mechain.tech")
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	block, err := client.BlockByNumber(context.TODO(), big.NewInt(372379))
	if err != nil {
		log.Fatalf("Failed to retrieve balance: %v", err)
	}
	log.Println(block.Transactions().Len())

	for i, transaction := range block.Transactions() {
		txHash := transaction.Hash()
		receipt, err := client.TransactionReceipt(context.TODO(), txHash)
		if err != nil {
			log.Fatalf("Failed to get tx Receipt: %v", err)
		}
		log.Println(i, receipt.TransactionIndex, receipt.TxHash.Hex(), receipt.BlockNumber, receipt.CumulativeGasUsed, receipt.Status)
	}
}
