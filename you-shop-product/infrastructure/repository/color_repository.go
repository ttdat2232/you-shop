package repository

import (
	"database/sql"

	"github.com/TechwizsonORG/product-service/entity"
	"github.com/TechwizsonORG/product-service/util"
	"github.com/google/uuid"
)

type ColorRepository struct {
	db *sql.DB
}

func NewColorRepository(db *sql.DB) *ColorRepository {
	return &ColorRepository{
		db: db,
	}
}

func (c *ColorRepository) AddColor(color *entity.Color) error {
	query := `
		INSERT INTO color (id, name, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
	`
	stmt, err := c.db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(color.Id, color.Name, util.GetCurrentUtcTime(7), util.GetCurrentUtcTime(7))
	if err != nil {
		return err
	}
	return nil
}
func (c *ColorRepository) GetColors() ([]entity.Color, error) {
	query := `
		SELECT
			c.id,
			c.name
		FROM color c
	`
	rows, err := c.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := []entity.Color{}
	for rows.Next() {
		var color entity.Color
		err = rows.Scan(&color.Id, &color.Name)
		if err != nil {
			return nil, err
		}
		result = append(result, color)
	}
	return result, nil
}
func (c *ColorRepository) GetBydId(id uuid.UUID) (*entity.Color, error) {

	query := `
		SELECT
			c.id,
			c.name
		FROM color c
		WHERE c.id = $1
	`
	row := c.db.QueryRow(query, id)
	var color entity.Color
	err := row.Scan(&color.Id, &color.Name)
	if err != nil {
		return nil, err
	}
	return &color, nil
}
