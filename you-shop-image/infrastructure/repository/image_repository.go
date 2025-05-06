package repository

import (
	"database/sql"

	"github.com/TechwizsonORG/image-service/entity"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/rs/zerolog"
)

type ImageRepository struct {
	logger zerolog.Logger
	db     *sql.DB
}

func NewImageRepository(logger zerolog.Logger, db *sql.DB) *ImageRepository {
	logger = logger.
		With().
		Str("infrastructure", "repository").
		Logger()
	return &ImageRepository{
		logger: logger,
		db:     db,
	}
}

func (i *ImageRepository) GetById(id uuid.UUID) (*entity.Image, error) {
	query := `
		SELECT
			i.image_url,
			i.content_type,
			i.alt
		FROM image i
		WHERE i.is_public = true
			AND i.id = $1
	`
	var contentTye sql.NullString
	image := &entity.Image{}
	err := i.db.QueryRow(query, id).Scan(&image.ImageUrl, &contentTye, &image.Alt)
	image.ContentType = contentTye.String
	if err != nil {
		return nil, err
	}

	return image, nil
}
func (i *ImageRepository) GetByOwnerIds(ownerIds []uuid.UUID) ([]entity.Image, error) {
	query := `
		SELECT
			i.id,
			COALESCE(i.image_url, 'https://i.pinimg.com/736x/5f/62/78/5f6278d64c70ab1bc7b2495117e0c23e.jpg'),
			i.content_type,
			i.owner_id,
			i.alt,
			owner_id.id
		FROM
			unnest($1::uuid[]) AS owner_id(id)
		LEFT JOIN image i ON i.owner_id = owner_id.id
	`
	result := []entity.Image{}
	rows, err := i.db.Query(query, pq.Array(ownerIds))
	if err != nil {
		i.logger.Error().Err(err).Msg("Error occurred")
		return result, nil
	}
	var contentType sql.NullString
	defer rows.Close()
	for rows.Next() {
		image := &entity.Image{}
		err = rows.Scan(&image.Id, &image.ImageUrl, &contentType, &image.OwnerId, &image.Alt, &image.OwnerId)
		if err != nil {
			i.logger.Error().Err(err).Msg("Error occurred")
			return []entity.Image{}, nil
		}
		image.ContentType = contentType.String
		result = append(result, *image)
	}
	return result, nil
}
func (i *ImageRepository) AddImage(image entity.Image) (int, error) {
	images := []entity.Image{}
	images = append(images, image)
	return i.AddImages(images)
}

func (i *ImageRepository) AddImages(images []entity.Image) (int, error) {

	tx, err := i.db.Begin()
	if err != nil {
		return 0, err
	}
	query := `
		INSERT INTO image (id, image_url, filename, content_type, owner_id, size, width, height, is_public, alt, type)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	for _, image := range images {
		stmt, err := tx.Prepare(query)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
		_, err = stmt.Exec(image.Id, image.ImageUrl, image.Filename, image.ContentType, image.OwnerId, image.Size, image.Width, image.Height, image.IsPublic, image.Alt, image.Type)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
	}
	err = tx.Commit()
	if err != nil {
		return 0, err
	}
	return len(images), nil
}

func (i *ImageRepository) DeleteImage(id uuid.UUID) (int, error) {
	ids := []uuid.UUID{}
	ids = append(ids, id)
	return i.DeleteImageByIds(ids)
}

// TODO: Change this function to hard delete. Soft delete will implemented by service layer
func (i *ImageRepository) DeleteImageByIds(ids []uuid.UUID) (int, error) {
	query := `
		UPDATE image
		SET is_public = false
		WHERE id = ANY($1)
	`
	stmt, err := i.db.Prepare(query)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(pq.Array(ids))
	if err != nil {
		return 0, err
	}
	numberOfRows, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return int(numberOfRows), nil
}

func (i *ImageRepository) UpdateImage(image entity.Image) (int, error) {
	return 0, nil
}

func (i *ImageRepository) GetByType(imageType entity.ImageType) ([]entity.Image, error) {
	query := `
		SELECT
			i.id,
			i.image_url,
			i.content_type,
			i.alt
		FROM IMAGE i
		WHERE i.is_public = true
			AND i."type" = $1
	`
	result := []entity.Image{}
	rows, err := i.db.Query(query, imageType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		image := entity.Image{}
		var contentType sql.NullString
		scanErr := rows.Scan(&image.Id, &image.ImageUrl, &contentType, &image.Alt)
		if scanErr != nil {
			return nil, scanErr
		}
		result = append(result, image)
	}
	return result, nil
}
