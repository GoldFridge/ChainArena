package handlers

import (
	"encoding/json"
	"net/http"

	"GuildVault/internal/db"
	"GuildVault/internal/models"
)

type CreateOrganizationRequest struct {
	CreatorWallet string `json:"creator_wallet"`
	ContractAddr  string `json:"contract_addr"`
}

func CreateTournamentHandler(w http.ResponseWriter, r *http.Request) {
	var req CreateOrganizationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	org := models.Tournament{
		ContractAddr:  req.ContractAddr,
		CreatorWallet: req.CreatorWallet,
	}

	if err := db.DB.Create(&org).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(org)
}

func JoinTournamentHandler(w http.ResponseWriter, r *http.Request) {
	type JoinRequest struct {
		UserWallet   string `json:"user_wallet"`
		ContractAddr string `json:"contract_addr"`
	}

	var req JoinRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var org models.Tournament
	if err := db.DB.Where("contract_addr = ?", req.ContractAddr).First(&org).Error; err != nil {
		http.Error(w, "organization not found", http.StatusNotFound)
		return
	}

	var user models.User
	if err := db.DB.Where("wallet_addr = ?", req.UserWallet).First(&user).Error; err != nil {
		user = models.User{WalletAddr: req.UserWallet}
		db.DB.Create(&user)
	}

	if err := db.DB.Model(&org).Association("Members").Append(&user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"status": "joined successfully",
	})
}
