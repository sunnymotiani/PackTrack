package users

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"golang.org/x/crypto/bcrypt"
)

func (us *UserService) ValidateByIDPassword(email string, password string) error {
	query := sq.Select("id", "password_hash", "account_status").
		From(TableUsers).Where(sq.Eq{"email": email})
	sql, args, err := query.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return fmt.Errorf("err validating user %w", err)
	}
	var user User
	var passwordHash string
	row := us.DB.QueryRow(sql, args...)
	err = row.Scan(&user.ID, &passwordHash, &user.AccountStatus)
	if err != nil {
		return fmt.Errorf("err validating user %w", err)
	}
	if user.AccountStatus == false {
		return fmt.Errorf("err user account status invalid %w", err)
	}
	// Compare the hashed password
	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
	if err != nil {
		return fmt.Errorf("invalid credentials")
	}
	return nil
}
