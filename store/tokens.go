package store

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mikolysz/enably/model"
)

// PostgresTokenStore 		lets you store and retrieve tokens.
type PostgresTokenStore struct {
	DB *pgxpool.Pool
}

// AddToken 		inserts a token into the database.
// The returned token will have the "id" field filled in with the ID of the new token.
func (s PostgresTokenStore) AddToken(c context.Context, t model.Token) (model.Token, error) {
	query := "INSERT INTO session_tokens(email_address, token) VALUES($1, $2) RETURNING id"
	row := s.DB.QueryRow(c, query, t.Email, t.Token)
	if err := row.Scan(&t.ID); err != nil {
		return model.Token{}, fmt.Errorf("error when inserting token: %s", err)
	}
	return t, nil
}

// GetByTokenContents 		returns the token with the given contents.
// returns model.ErrInvalidToken if the token does not exist.
func (s PostgresTokenStore) GetByTokenContents(c context.Context, token string) (model.Token, error) {

	query := "SELECT id, email_address, token, created_at FROM session_tokens WHERE token = $1"

	var t model.Token
	row := s.DB.QueryRow(c, query, token)
	err := row.Scan(&t.ID, &t.Email, &t.Token, &t.CreatedAt)

	if err == pgx.ErrNoRows {
		return model.Token{}, model.ErrInvalidToken
	}

	if err != nil {
		return model.Token{}, fmt.Errorf("error when scanning token: %s", err)
	}

	return t, nil
}
