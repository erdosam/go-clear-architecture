package usecase

import (
	"github.com/arendi-project/ba-version-2/internal/entity"
	repo "github.com/arendi-project/ba-version-2/internal/usecase/dao"
	"github.com/arendi-project/ba-version-2/pkg/logger"
	"github.com/go-playground/validator/v10"
)

var _ Category = &categoryUsecase{}

type categoryUsecase struct {
	log logger.Interface
	dao repo.TrashCategoryDAO
	val *validator.Validate
}

func NewCategoryUsecase(l logger.Interface, d repo.TrashCategoryDAO, v *validator.Validate) Category {
	return &categoryUsecase{l, d, v}
}

func (uc *categoryUsecase) GetAvailableCategories() ([]entity.TrashCategory, error) {
	categories, err := uc.dao.FindCategories("status", "enabled")
	if err != nil {
		return nil, err
	}
	//TODO must return TrashCategoryResponse
	return categories, nil
}
