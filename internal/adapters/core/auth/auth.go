package auth

import "github.com/Renewdxin/selfMade/internal/ports/core/auth"

type AccountAdapter struct{}

func NewAdapter() AccountAdapter {
	return AccountAdapter{}
}

func (accountAdapter AccountAdapter) CreateAccount(id, password string) (auth.Account, error) {
	return auth.Account{ID: id, Password: password}, nil
}

func (accountAdapter AccountAdapter) TableName() string {
	return "accountInfo"
}
