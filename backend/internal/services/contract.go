package services

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/zksync-sdk/zksync2-go/accounts"
	"github.com/zksync-sdk/zksync2-go/clients"
	"log"
	"os"
)

func TestDeploy() {
	ZkSyncEraProvider := "https://sepolia.era.zksync.dev"
	//ZkSyncEraWSProvider := "ws://testnet.era.zksync.dev:3051"

	client, err := clients.Dial(ZkSyncEraProvider)
	if err != nil {
		log.Panic(err)
	}
	defer client.Close()

	chainID, err := client.ChainID(context.Background())
	if err != nil {
		log.Panic(err)
	}
	fmt.Println("Chain ID: ", chainID)

}

func TournamentDeploy() {
	var (
		PrivateKey        = os.Getenv("0x6D4112a188c4FeddA51C945EDF3ACb7980A692bC")
		ZkSyncEraProvider = "https://sepolia.era.zksync.dev"
	)

	// Connect to ZKsync network
	client, err := clients.Dial(ZkSyncEraProvider)
	if err != nil {
		log.Panic(err)
	}
	defer client.Close()

	// Create wallet
	wallet, err := accounts.NewWallet(common.Hex2Bytes(PrivateKey), client, nil)
	if err != nil {
		log.Panic(err)
	}

	// Read smart contract bytecode
	bytecode, err := os.ReadFile("tournament_sol_Tournament.bin")
	if err != nil {
		log.Panic(err)
	}

	//Deploy smart contract
	hash, err := wallet.DeployWithCreate(nil, accounts.CreateTransaction{Bytecode: bytecode})
	if err != nil {
		panic(err)
	}
	fmt.Println("Transaction: ", hash)

	// Wait unit transaction is finalized
	receipt, err := client.WaitMined(context.Background(), hash)
	if err != nil {
		log.Panic(err)
	}

	contractAddress := receipt.ContractAddress
	fmt.Println("Smart contract address", contractAddress.String())
}
