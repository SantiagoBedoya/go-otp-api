package mysql

import (
	"database/sql"
	"errors"

	"github.com/SantiagoBedoya/otp-api/internal/model"
	"github.com/SantiagoBedoya/otp-api/internal/service"
	"github.com/go-sql-driver/mysql"
)

const (
	GetUsersQuery       = "SELECT id, first_name, last_name, email FROM users"
	GetUserByEmailQuery = "SELECT id, email, password FROM users WHERE email = ?"
	GetUserByIDQuery    = "SELECT email, secret_2fa FROM users WHERE id = ?"
	SaveUserQuery       = "INSERT INTO users (first_name, last_name, email, password) VALUES (?, ?, ?, ?)"
	SaveUserSecretQuery = "UPDATE users SET secret_2fa = ? WHERE id = ?"
	SetValidSecretQuery = "UPDATE users SET 2fa_valid = 1 WHERE id = ?"
)

type userRepo struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) service.UserRepository {
	return &userRepo{
		db: db,
	}
}

func (r *userRepo) SaveSecret(userID, secret string) error {
	stmt, err := r.db.Prepare(SaveUserSecretQuery)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(secret, userID)
	if err != nil {
		return err
	}
	return nil
}

func (r *userRepo) SetValidSecret(userID string) error {
	stmt, err := r.db.Prepare(SetValidSecretQuery)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(userID)
	if err != nil {
		return err
	}
	return nil
}

func (r *userRepo) Save(u *model.User) error {
	stmt, err := r.db.Prepare(SaveUserQuery)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(u.FirstName, u.LastName, u.Email, u.Password)
	if err != nil {
		me, ok := err.(*mysql.MySQLError)
		if !ok {
			return err
		}
		if me.Number == 1062 {
			return service.ErrEmailInUse
		}
	}
	lastID, err := result.LastInsertId()
	if err != nil {
		return err
	}
	u.ID = lastID
	return nil
}

func (r *userRepo) GetAll() ([]model.User, error) {
	stmt, err := r.db.Prepare(GetUsersQuery)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	users := make([]model.User, 0)
	for rows.Next() {
		var user model.User
		err = rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *userRepo) GetByEmail(email string) (*model.User, error) {
	stmt, err := r.db.Prepare(GetUserByEmailQuery)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var user model.User
	err = stmt.QueryRow(email).Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, service.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) GetByID(userID string) (*model.User, error) {
	stmt, err := r.db.Prepare(GetUserByIDQuery)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var user model.User
	err = stmt.QueryRow(userID).Scan(&user.Email, &user.Secret2FA)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, service.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}
