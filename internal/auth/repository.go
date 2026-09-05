package auth

import (
	"database/sql"
	"time"
)

type AuthRepository interface {
	CreateUser(username string, passwordHash string, role string) (int64, error)
	FindByUsername(username string) (*User, error)
	FindByID(id int64) (*User, error)
	SaveRefreshToken(userID int64, token string, expires time.Time) error
	GetRefreshToken(token string) (*RefreshToken, error)
	RevokeRefreshToken(token string) error
	RevokeAllRefreshTokens(userID int64) error
}

type PostgresAuthRepository struct {
	DB *sql.DB
}

func NewRepository(db *sql.DB) AuthRepository {
	return &PostgresAuthRepository{DB: db}
}

func (r *PostgresAuthRepository) CreateUser(username string, passwordHash string, role string) (int64, error) {
	var id int64
	query := `INSERT INTO users (username, password_hash,role) VALUES ($1, $2,$3) RETURNING id;`
	err := r.DB.QueryRow(query, username, passwordHash, role).Scan(&id)
	return id, err
}

func (r *PostgresAuthRepository) FindByUsername(username string) (*User, error) {
	query := `SELECT id, username, password_hash, role, created_at FROM users WHERE username = $1`

	user := &User{}
	err := r.DB.QueryRow(query, username).Scan(
		&user.ID, &user.Username, &user.PasswordHash, &user.Role, &user.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return user, nil
}

func (r *PostgresAuthRepository) FindByID(id int64) (*User, error) {
	query := `SELECT id, username, password_hash, role, created_at 
			  FROM users 
			  WHERE id = $1`

	user := &User{}

	err := r.DB.QueryRow(query, id).Scan(
		&user.ID,
		&user.Username,
		&user.PasswordHash,
		&user.Role,
		&user.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return user, err
}

func (r *PostgresAuthRepository) SaveRefreshToken(userID int64, token string, expires time.Time) error {
	query := `INSERT INTO refresh_tokens (user_id, token, expires_at) VALUES ($1, $2, $3);`

	_, err := r.DB.Exec(query, userID, token, expires)
	return err
}

func (r *PostgresAuthRepository) GetRefreshToken(token string) (*RefreshToken, error) {
	query := `SELECT id, user_id, token, expires_at, created_at, revoked
            FROM refresh_tokens
            WHERE token = $1;
            `

	rt := &RefreshToken{}
	err := r.DB.QueryRow(query, token).Scan(
		&rt.ID, &rt.UserID, &rt.Token, &rt.ExpiresAt, &rt.CreatedAt, &rt.Revoked,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	return rt, err
}

func (r *PostgresAuthRepository) RevokeRefreshToken(token string) error {
	query := `UPDATE refresh_tokens SET revoked = true WHERE token = $1;`
	_, err := r.DB.Exec(query, token)
	return err
}

func (r *PostgresAuthRepository) RevokeAllRefreshTokens(userID int64) error {
	query := `UPDATE refresh_tokens SET revoked = true WHERE user_id = $1;`
	_, err := r.DB.Exec(query, userID)
	return err
}
