package store

import (
	"github.com/VIPowERuS/nsu_postman/internal/app/model"
)

// UserRepository ...
type UserRepository struct {
	store *Store
}

// FindByEmail ...
func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	u := &model.User{}
	if err := r.store.db.QueryRow("SELECT id, email, encrypted_password, access FROM users WHERE email = $1",
		email).Scan(&u.ID, &u.Email, &u.EncryptedPassword, &u.Access); err != nil {
		return nil, err
	}
	return u, nil
}
