package users

import (
	"database/sql"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Email         string    `json:"email"`
	CreatedAt     time.Time `json:"created_at"`
	AccountStatus bool      `json:"account_status"`
}

type UserService struct {
	DB          *sql.DB
	RedisClient *redis.Client
}

const TableUsers = "users"

func (us *UserService) CreateUser(name string, email string, password string) (*User, error) {
	var user User
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("create user %w", err)
	}
	PasswordHash := string(hashedBytes)
	query := sq.Insert(TableUsers).Columns("name", "email", "password_hash").
		Values(name, email, PasswordHash).Suffix("RETURNING id").PlaceholderFormat(sq.Dollar)
	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("create user generating query %w", err)
	}
	row := us.DB.QueryRow(sql, args...)
	err = row.Scan(&user.ID)
	if err != nil {
		return nil, err
	}
	user.Email = email
	user.Name = name
	return &user, nil
}
func (us *UserService) GetUserByID(id string) (*User, error) {
	query := sq.Select("name", "email", "created_at", "account_status").
		From(TableUsers).Where(sq.Eq{"id": id})
	sql, args, err := query.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, fmt.Errorf("get user by ID query %w", err)
	}
	user := &User{}
	row := us.DB.QueryRow(sql, args...)
	err = row.Scan(&user.Name, &user.Email, &user.CreatedAt, &user.AccountStatus)
	return user, fmt.Errorf("get user by ID scanning row %w", err)
}
func (us *UserService) GetUserByEmail(email string) (*User, error) {
	query := sq.Select("name", "id", "created_at", "account_status").
		From(TableUsers).Where(sq.Eq{"email": email})
	sql, args, err := query.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, fmt.Errorf("get user by ID query %w", err)
	}
	user := &User{}
	row := us.DB.QueryRow(sql, args...)
	err = row.Scan(&user.Name, &user.ID, &user.CreatedAt, &user.AccountStatus)
	return user, fmt.Errorf("get user by ID scanning row %w", err)
}
func (us *UserService) UpdateUser(new User) error {
	values := sq.Eq{}
	values["name"] = new.Name
	values["email"] = new.Email
	query := sq.Update(TableUsers).SetMap(values).PlaceholderFormat(sq.Dollar)
	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("err updating user %w", err)
	}
	_, err = us.DB.Exec(sql, args...)
	if err != nil {
		return fmt.Errorf("error updating user query: %w", err)
	}
	return nil
}
func (us *UserService) UpdateUserPassword(id string, password string) error {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("create user %w", err)
	}
	PasswordHash := string(hashedBytes)
	query := sq.Update(TableUsers).Set("password_hash", PasswordHash).Where(sq.Eq{"id": id}).PlaceholderFormat(sq.Dollar)
	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("error updating uper password query: %w", err)
	}
	_, err = us.DB.Exec(sql, args...)
	if err != nil {
		return fmt.Errorf("error updating user password: %w", err)
	}
	return nil
}
