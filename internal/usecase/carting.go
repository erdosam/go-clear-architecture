package usecase

import (
	"github.com/arendi-project/ba-version-2/internal/entity"
	repo "github.com/arendi-project/ba-version-2/internal/usecase/dao"
	"github.com/arendi-project/ba-version-2/pkg/logger"
)

var _ Carting = &cartingUseCase{}

type cartingUseCase struct {
	dao repo.CartingDAO
	log logger.Interface
}

func NewCartingUseCase(l logger.Interface, d repo.CartingDAO) Carting {
	return &cartingUseCase{
		dao: d,
		log: l,
	}
}

func (us *cartingUseCase) GetItems(c entity.Cart) ([]entity.CartItem, error) {
	items, err := us.dao.FindItemsByCart(c)
	if err != nil {
		us.log.Error(err)
		return nil, err
	}
	return items, nil
}

func (us *cartingUseCase) GetItem(c entity.Cart, id string) (entity.CartItem, error) {
	item, err := us.dao.FindOneItem("id", id, "user_id", c.UserId)
	if err != nil {
		us.log.Error(err)
		return entity.CartItem{}, err
	}
	return item, nil
}

func (us *cartingUseCase) AddItemToCart(i entity.CartItem, c entity.Cart) error {
	return nil
}

func (us *cartingUseCase) GetCart(u entity.User) (entity.Cart, error) {
	return entity.Cart{UserId: u.Id}, nil
}
