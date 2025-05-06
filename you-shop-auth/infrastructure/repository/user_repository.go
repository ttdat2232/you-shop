package repository

import (
	"database/sql"
	"errors"

	"github.com/TechwizsonORG/auth-service/entity"
	"github.com/TechwizsonORG/auth-service/util"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
)

type UserRepository struct {
	db  *sql.DB
	log zerolog.Logger
}

func NewUserRepository(db *sql.DB, logger zerolog.Logger) *UserRepository {
	logger = logger.
		With().
		Str("Infrastructure", "User Repository").
		Logger()
	return &UserRepository{
		db:  db,
		log: logger,
	}
}

func (u UserRepository) GetUserByEmailOrUsername(email string, username string) (*entity.User, error) {
	query := `
		SELECT 
			u.id,
			u.created_at,
			u.updated_at,
			u.username,
			u.email,
			u.password_hash,
			u.avatar_url,
			u.is_active,
			u.phone_number
		FROM users u
		WHERE (u.email = $1 OR u.username = $2) AND u.is_active = true
	`
	var user entity.User
	error := u.db.QueryRow(query, email, username).Scan(
		&user.Id,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.AvatarUrl,
		&user.IsActive,
		&user.PhoneNumber,
	)

	if error != nil {
		if error.Error() != sql.ErrNoRows.Error() {
			return nil, error
		}
		return nil, nil
	}
	return &user, nil
}

func (u UserRepository) GetUserByID(id uuid.UUID) (*entity.User, error) {
	query := `
		SELECT 
			u.id,
			u.created_at,
			u.updated_at,
			u.username,
			u.email,
			u.password_hash,
			u.avatar_url,
			u.is_active,
			u.phone_number
		FROM users u
		WHERE u.id = $1 AND u.is_active = true
	`
	var user entity.User
	error := u.db.QueryRow(query, id.String()).Scan(
		&user.Id,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.AvatarUrl,
		&user.IsActive,
		&user.PhoneNumber,
	)

	if error != nil {
		return nil, error
	}
	return &user, nil

}

func (u UserRepository) AddUser(user entity.User) (entity.User, error) {
	query := `
		INSERT INTO users (id, username, normalized_username, password_hash, email, normalized_email, avatar_url, is_active, phone_number)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id
	`
	smtp, err := u.db.Prepare(query)
	if err != nil {
		return entity.User{}, err
	}
	defer smtp.Close()

	result, err := smtp.Exec(user.Id, user.Username, user.NormalizedEmail, user.PasswordHash, user.Email, user.NormalizedEmail, user.AvatarUrl, user.IsActive, user.PhoneNumber)
	if err != nil {
		return entity.User{}, err
	}
	if rows, err := result.RowsAffected(); rows == 0 || err != nil {
		return entity.User{}, errors.New("user haven't been created")
	}
	return user, nil
}

func (u UserRepository) UpdateUser(user entity.User) (entity.User, error) {
	query := `
		UPDATE user
		SET
			avatar_url = $1,
			is_active = $2,
			phone_number = $3,
			updated_at = $4
		WHERE id = $5
		RETURNING updated_at
	`
	smtp, err := u.db.Prepare(query)

	if err != nil {
		return entity.User{}, err
	}
	defer smtp.Close()

	err = smtp.
		QueryRow(user.AvatarUrl, user.IsActive, user.PhoneNumber, util.GetCurrentUtcTime(7), user.Id.String()).
		Scan(&user.UpdatedAt)
	if err != nil {
		return entity.User{}, err
	}
	return user, nil
}
