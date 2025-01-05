package webapi

import (
	"github.com/arendi-project/ba-version-2/internal/entity"
)

type junoService struct {
}

func NewJunoService() *junoService {
	return &junoService{}
}

func (j junoService) GetAccount() (account *entity.UserAccount, err error) {
	return nil, nil
}
