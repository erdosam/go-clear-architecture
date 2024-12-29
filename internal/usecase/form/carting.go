package form

import (
	"github.com/arendi-project/ba-version-2/internal/entity"
	"github.com/arendi-project/ba-version-2/internal/usecase"
	"github.com/arendi-project/ba-version-2/pkg/logger"
)

type CartingUseCase struct {
	dao usecase.CartingDAO
	log logger.Interface
}

func NewCarting(l logger.Interface, d usecase.CartingDAO) *CartingUseCase {
	return &CartingUseCase{
		dao: d,
		log: l,
	}
}

func (us *CartingUseCase) GetItems(c entity.Cart) ([]entity.CartItem, error) {
	items, err := us.dao.GetItemsByCart(c)
	if err != nil {
		us.log.Error(err)
		return nil, err
	}
	return items, nil
}

func (us *CartingUseCase) AddItemToCart(i entity.CartItem, c entity.Cart) error {
	return nil
}
