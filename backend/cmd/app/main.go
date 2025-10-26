package main

import (
	"GuildVault/internal/handlers"
	"GuildVault/internal/services"
	"log"
	"net/http"
)

func main() {
	services.TournamentDeploy("smartcontracts/TournamentInstance.json")
	services.TournamentDeploy("smartcontracts/TournamentFactory.json")
	services.TournamentDeploy("smartcontracts/ItemVault.json")
	services.TournamentDeploy("smartcontracts/WrappedItemToken.json")
	services.TournamentDeploy("smartcontracts/InventoryManager.json")

	http.HandleFunc("/tournament/create", handlers.CreateTournamentHandler)
	http.HandleFunc("/tournament/join", handlers.JoinTournamentHandler)
	http.HandleFunc("/login", handlers.LoginHandler)

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
	services.TestDeploy()
}
