package dao

import (
	"github.com/arendi-project/ba-version-2/internal/entity"
	"github.com/arendi-project/ba-version-2/pkg/logger"
	"github.com/arendi-project/ba-version-2/pkg/postgres"
)

var _ TrashCategoryDAO = &categoryDAO{}

type categoryDAO struct {
	*postgres.Postgres
	log logger.Interface
}

func NewTrashCategoryDAO(l logger.Interface, pg *postgres.Postgres) TrashCategoryDAO {
	return &categoryDAO{pg, l}
}

func (dao *categoryDAO) FindCategories(arg ...interface{}) ([]entity.TrashCategory, error) {
	var rows []entity.TrashCategory
	query, args := buildFindQuery(`
		SELECT c.id,
			   c.name,
			   c.parent_category_id,
			   c."group",
			   c.status,
			   cd.description,
			   cd.image_url_md as "image.image_url_md",
			   cd.image_url_sm as "image.image_url_sm",
			   cd.image_url_lg as "image.image_url_lg"
		FROM public.trash_category c
				 JOIN public.trash_category_detail cd ON cd.category_id = c.id
		`, arg...)
	q := dao.Rebind(query)
	if err := dao.Select(&rows, q, args...); err != nil {
		dao.log.Debug(err)
		if err.Error() == errorSqlEmptyResult {
			return []entity.TrashCategory{}, err
		}
		panic(err.(any))
	}
	return rows, nil
}
