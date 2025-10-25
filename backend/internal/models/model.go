package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID            uuid.UUID    `gorm:"type:uuid;primaryKey" json:"id"`
	WalletAddr    string       `gorm:"uniqueIndex;not null" json:"wallet_addr"`
	Organizations []Tournament `gorm:"many2many:user_organizations;" json:"organizations,omitempty"`
}

type Tournament struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	ContractAddr  string    `gorm:"uniqueIndex;not null" json:"contract_addr"`
	CreatorWallet string    `gorm:"not null" json:"creator_wallet"`
	Members       []User    `gorm:"many2many:user_organizations;" json:"members,omitempty"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return
}
