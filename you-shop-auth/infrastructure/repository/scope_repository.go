package repository

import (
	"database/sql"

	"github.com/TechwizsonORG/auth-service/entity"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type ScopeRepository struct {
	db     *sql.DB
	logger zerolog.Logger
}

func NewScopeRepository(db *sql.DB, logger zerolog.Logger) *ScopeRepository {
	logger = logger.
		With().
		Str("Infrastructure", "Scope Repository").
		Logger()
	return &ScopeRepository{
		db:     db,
		logger: logger,
	}
}

func (s *ScopeRepository) GetScopesByUserId(userId uuid.UUID) []entity.Scope {
	query := `
		SELECT s.id, s.name FROM "scope" s 
		JOIN user_scope us ON s.id = us.scope_id
		WHERE us.user_id = $1
	`
	rows, queryErr := s.db.Query(query, userId)
	var result = []entity.Scope{}
	if queryErr != nil {
		s.logger.Error().Err(queryErr).Msg("")
		return result
	}
	defer rows.Close()
	for rows.Next() {
		var scope entity.Scope
		scanErr := rows.Scan(&scope.Id, &scope.Name)
		if scanErr != nil {
			return result
		}
		result = append(result, scope)
	}
	return result
}
