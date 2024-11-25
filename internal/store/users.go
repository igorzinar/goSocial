package store

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  password  `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	IsActive  bool      `json:"is_active"`
}

var (
	ErrDuplicateEmail    = errors.New("a user with that email already exists")
	ErrDuplicateUsername = errors.New("a user with that username already exists")
)

type password struct {
	text *string
	hash []byte
}

func (p *password) Set(text string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	p.text = &text
	p.hash = hash
	return nil
}

type UserStore struct {
	db *sql.DB
}

func (s *UserStore) Create(ctx context.Context, tx *sql.Tx, user *User) error {
	query := `INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING id, created_at`
	ctx, cancel := context.WithTimeout(ctx, TimeoutDuration)
	defer cancel()

	err := tx.QueryRowContext(ctx, query, user.Username, user.Email, user.Password.hash).Scan(&user.ID, &user.CreatedAt)

	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		case err.Error() == `pq: duplicate key value violates unique constraint "users_username_key"`:
			return ErrDuplicateUsername
		default:
			return err
		}
	}

	return nil
}

func (s *UserStore) Delete(ctx context.Context, id int64) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {

		if err := s.delete(ctx, tx, id); err != nill {
			return err
		}
		if err := s.deleteUserInvitation(ctx, tx, id); err != nil {
			return err
		}

		return nil
	})

}

func (s *UserStore) GetByID(ctx context.Context, id int64) (*User, error) {
	query := `SELECT id, username, created_at, email FROM users WHERE id = $1`
	var user User
	err := s.db.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.Username, &user.CreatedAt, &user.Email)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}

func (s *UserStore) CreateAndInvite(ctx context.Context, user *User, token string, invitationExp time.Duration) error {
	// create transaction wrapper
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		// create the user
		if err := s.Create(ctx, tx, user); err != nil {
			return err
		}
		// create user invite
		if err := s.createUserInvitation(ctx, tx, token, invitationExp, user.ID); err != nil {
			return err
		}
		return nil
	})

}

func (s *UserStore) Activate(ctx context.Context, token string) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		// find user and check if not expired token
		user, err := s.getUserFromInvitation(ctx, tx, token)
		if err != nil {
			return err
		}
		// update the user
		user.IsActive = true
		if err := s.update(ctx, tx, user); err != nil {
			return err
		}
		// clean invitations

		if err := s.deleteUserInvitation(ctx, tx, user.ID); err != nil {
			return err
		}

		return nil

	})
}
func (s *UserStore) createUserInvitation(ctx context.Context, tx *sql.Tx, token string, invitationExp time.Duration, userID int64) error {
	query := `

INSERT INTO user_invitations (token, user_id, expiry) VALUES ($1, $2, $3)`

	ctx, cancel := context.WithTimeout(ctx, TimeoutDuration)
	defer cancel()
	_, err := tx.ExecContext(ctx, query, token, userID, time.Now().Add(invitationExp))
	if err != nil {
		return err
	}
	return nil
}

func (s *UserStore) getUserFromInvitation(ctx context.Context, tx *sql.Tx, token string) (*User, error) {
	query := `
SELECT u.id, u.username, u.email, u.created_at, u.is_active 
FROM users u 
JOIN user_invitations ui ON ui.user_id = u.id
WHERE ui.token = $1 AND  ui.expiry > $2
`

	hash := sha256.Sum256([]byte(token))
	hashToken := hex.EncodeToString(hash[:])

	ctx, cancel := context.WithTimeout(ctx, TimeoutDuration)
	defer cancel()
	user := &User{}
	err := tx.QueryRowContext(ctx, query, hashToken, time.Now()).Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt, &user.IsActive)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}
	return user, nil
}

func (s *UserStore) update(ctx context.Context, tx *sql.Tx, user *User) error {
	query := `UPDATE users SET username = $1, email = $2, is_active = $3 WHERE id = $4`

	ctx, cancel := context.WithTimeout(ctx, TimeoutDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, user.Username, user.Email, user.IsActive, user.ID)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserStore) deleteUserInvitation(ctx context.Context, tx *sql.Tx, id int64) error {
	query := `DELETE FROM user_invitations WHERE user_id = $1`
	ctx, cancel := context.WithTimeout(ctx, TimeoutDuration)
	defer cancel()
	_, err := tx.ExecContext(ctx, query, id)
	if err != nil {
		return err

	}
	return nil
}

func (s *UserStore) delete(ctx context.Context, tx *sql.Tx, id int64) error {
	query := `DELETE FROM users WHERE id = $1`
	ctx, cancel := context.WithTimeout(ctx, TimeoutDuration)
	defer cancel()
	_, err := tx.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}
