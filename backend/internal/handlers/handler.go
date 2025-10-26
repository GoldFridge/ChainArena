package handlers

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"strings"
	"time"

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

var jwtSecret = []byte("super_secret_key")

var nonces = map[string]string{}

type NonceResponse struct {
	Nonce string `json:"nonce"`
}

type LoginRequest struct {
	Address   string `json:"address"`
	Signature string `json:"signature"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func nonceHandler(w http.ResponseWriter, r *http.Request) {
	addr := r.URL.Query().Get("address")
	if addr == "" {
		http.Error(w, "missing address", http.StatusBadRequest)
		return
	}

	nonce := fmt.Sprintf("Login with wallet %s at %d", addr, time.Now().UnixNano())
	nonces[strings.ToLower(addr)] = nonce

	json.NewEncoder(w).Encode(NonceResponse{Nonce: nonce})
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	nonce, ok := nonces[strings.ToLower(req.Address)]
	if !ok {
		http.Error(w, "nonce not found", http.StatusUnauthorized)
		return
	}

	sig := strings.TrimPrefix(req.Signature, "0x")

	sigBytes, err := hex.DecodeString(sig)
	if err != nil {
		http.Error(w, "invalid signature format", http.StatusBadRequest)
		return
	}

	msg := accounts.TextHash([]byte(nonce))

	if sigBytes[64] == 27 || sigBytes[64] == 28 {
		sigBytes[64] -= 27
	}

	pubKey, err := crypto.SigToPub(msg, sigBytes)
	if err != nil {
		http.Error(w, "signature recovery failed", http.StatusUnauthorized)
		return
	}

	recoveredAddr := crypto.PubkeyToAddress(*pubKey)

	if !strings.EqualFold(recoveredAddr.Hex(), req.Address) {
		http.Error(w, "invalid signature", http.StatusUnauthorized)
		return
	}

	// Создаем JWT токен
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"address": req.Address,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenStr, err := token.SignedString(jwtSecret)
	if err != nil {
		http.Error(w, "token creation failed", http.StatusInternalServerError)
		return
	}

	delete(nonces, strings.ToLower(req.Address))

	json.NewEncoder(w).Encode(LoginResponse{Token: tokenStr})
}
