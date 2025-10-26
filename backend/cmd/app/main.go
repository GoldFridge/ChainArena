package main

import (
	"GuildVault/internal/handlers"
	"GuildVault/internal/services"
	"log"
	"net/http"
)

func main() {
	services.TournamentDeploy()
	http.HandleFunc("/tournament/create", handlers.CreateTournamentHandler)
	http.HandleFunc("/tournament/join", handlers.JoinTournamentHandler)
	http.HandleFunc("/login", handlers.LoginHandler)

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
	services.TestDeploy()
}
