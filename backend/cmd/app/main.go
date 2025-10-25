package main

import (
	"GuildVault/internal/handlers"
	"GuildVault/internal/services"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"log"
	"net/http"
)

func main() {

	priv, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}
	privBytes := crypto.FromECDSA(priv)
	fmt.Printf("PRIVATE_KEY=0x%x\n", privBytes)
	pub := priv.Public().(*ecdsa.PublicKey)
	fmt.Printf("ADDRESS: %s\n", crypto.PubkeyToAddress(*pub).Hex())
	//	db.InitDB()
	services.TournamentDeploy()
	http.HandleFunc("/tournament/create", handlers.CreateTournamentHandler)
	http.HandleFunc("/tournament/join", handlers.JoinTournamentHandler)

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
	services.TestDeploy()
}
