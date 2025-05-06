package repository

import (
	"database/sql"

	"github.com/TechwizsonORG/product-service/entity"
	"github.com/TechwizsonORG/product-service/util"
)

type SizeRepository struct {
	db *sql.DB
}

func NewSizeRepository(db *sql.DB) *SizeRepository {
	return &SizeRepository{
		db: db,
	}
}
func (s *SizeRepository) AddSize(entity *entity.Size) error {
	query := `
		INSERT INTO "size" (id, name, created_at, updated_at)
		VALUES (?, ?, ?, ?)
	`
	ConvertTemplate(&query)
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	current := util.GetCurrentUtcTime(7)
	_, err = stmt.Exec(entity.Id, entity.Name, current, current)
	return err
}
func (s *SizeRepository) GetSizes() ([]entity.Size, error) {
	query := `
		SELECT
			s.id,
			s.Name
		FROM "size" s
	`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := []entity.Size{}
	for rows.Next() {
		var size entity.Size
		err = rows.Scan(&size.Id, &size.Name)
		if err != nil {
			return nil, err
		}
		result = append(result, size)
	}
	return result, nil
}
