package usecase

import (
	"github.com/arendi-project/ba-version-2/internal/entity"
	repo "github.com/arendi-project/ba-version-2/internal/usecase/dao"
	"github.com/arendi-project/ba-version-2/pkg/logger"
)

type orderUseCase struct {
	cartingDAO repo.CartingDAO
	orderDAO   repo.OrderDAO
	log        logger.Interface
}

func NewOrderUseCase(l logger.Interface, od repo.OrderDAO, cd repo.CartingDAO) Order {
	return &orderUseCase{
		cartingDAO: cd,
		orderDAO:   od,
		log:        l,
	}
}

func (o orderUseCase) CreateMultipleItemsOrder([]entity.CartItem) (entity.Order, error) {
	return entity.Order{}, nil
}

func (o orderUseCase) GetOrderById(id string) (entity.Order, error) {
	o.log.Info(id)
	return entity.Order{}, nil
}
