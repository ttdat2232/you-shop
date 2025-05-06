package repository

import (
	"database/sql"

	"github.com/TechwizsonORG/auth-service/entity"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type RoleRepository struct {
	db     *sql.DB
	logger zerolog.Logger
}

func NewRoleRepository(db *sql.DB, logger zerolog.Logger) *RoleRepository {
	logger = logger.
		With().
		Str("Infrastructure", "Role Repository").
		Logger()
	return &RoleRepository{
		db:     db,
		logger: logger,
	}
}

func (r *RoleRepository) GetRolesByUserId(userId uuid.UUID) []entity.Role {
	query := `
		SELECT r.id, r.name FROM "role" r 
		JOIN user_role ur ON r.id = ur.role_id
		WHERE ur.user_id = $1
	`
	rows, queryErr := r.db.Query(query, userId)
	var result = []entity.Role{}
	if queryErr != nil {
		r.logger.Error().Err(queryErr).Msg("")
		return result
	}
	defer rows.Close()
	for rows.Next() {
		var role entity.Role
		scanErr := rows.Scan(&role.Id, &role.Name)
		if scanErr != nil {
			return result
		}
		result = append(result, role)
	}
	return result
}
