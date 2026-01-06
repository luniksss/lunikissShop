package repository

import (
	"context"
	"database/sql"
	"fmt"
	"lunikissShop/pkg/domain/model"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (ur UserRepository) ListAllUsers(ctx context.Context) ([]model.User, error) {
	query := `
        SELECT 
            id, name, surname, email, role, phone
        FROM user
        ORDER BY id
    `

	rows, err := ur.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var us model.User
		err := rows.Scan(
			&us.ID, &us.Name, &us.Surname, &us.Email, &us.Role, &us.Phone,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, us)
	}

	return users, nil
}

func (ur UserRepository) GetUserByID(ctx context.Context, userID string) (model.User, error) {
	query := `
        SELECT 
            id, name, surname, email, role, phone
        FROM user
        WHERE id=?
    `
	rows, err := ur.db.QueryContext(ctx, query, userID)
	if err != nil {
		return model.User{}, err
	}
	defer rows.Close()

	var oi model.User
	for rows.Next() {
		err := rows.Scan(&oi.ID, &oi.Name, &oi.Surname, &oi.Email, &oi.Role, &oi.Phone)
		if err != nil {
			return model.User{}, err
		}
	}

	return oi, nil
}

func (ur UserRepository) GetUserByEmail(ctx context.Context, userEmail string) (model.User, error) {
	query := `
        SELECT 
            id, name, surname, email, role, phone
        FROM user
        WHERE email=?
    `
	rows, err := ur.db.QueryContext(ctx, query, userEmail)
	if err != nil {
		return model.User{}, err
	}
	defer rows.Close()

	var oi model.User
	for rows.Next() {
		err := rows.Scan(&oi.ID, &oi.Name, &oi.Surname, &oi.Email, &oi.Role, &oi.Phone)
		if err != nil {
			return model.User{}, err
		}
	}

	return oi, nil
}

func (ur UserRepository) GetUserPassword(ctx context.Context, userID string) (string, error) {
	query := `
        SELECT password
        FROM user
        WHERE id=?
    `
	var password string
	err := ur.db.QueryRowContext(ctx, query, userID).Scan(&password)
	if err != nil {
		return "", err
	}
	return password, nil
}

func (ur UserRepository) AddUser(ctx context.Context, userInfo *model.UserInfo) error {
	query := `INSERT INTO user (name, surname, email, password, role, phone) VALUES (?, ?, ?, ?, ?, ?)`
	_, err := ur.db.ExecContext(ctx, query, userInfo.Name, userInfo.Surname, userInfo.Email, userInfo.Password, userInfo.Role, userInfo.Phone)
	if err != nil {
		return err
	}
	return nil
}

func (ur UserRepository) UpdateUser(ctx context.Context, newUserInfo *model.UserInfo) error {
	query := `UPDATE user 
		SET name = ?, surname = ?, email = ?, role = ?, phone = ?
		WHERE id = ?
	`
	_, err := ur.db.ExecContext(ctx, query, newUserInfo.Name, newUserInfo.Surname, newUserInfo.Email, newUserInfo.Role, newUserInfo.Phone, newUserInfo.ID)
	if err != nil {
		return err
	}

	return nil
}

func (ur UserRepository) UpdateUserRole(ctx context.Context, userID string, newUserRole string) error {
	query := `UPDATE user 
		SET role = ?
		WHERE id = ?
	`
	_, err := ur.db.ExecContext(ctx, query, newUserRole, userID)
	if err != nil {
		return err
	}

	return nil
}

func (ur UserRepository) UpdatePassword(ctx context.Context, userID string, newPassword string) error {
	query := `UPDATE user 
		SET password = ?
		WHERE id = ?
	`
	_, err := ur.db.ExecContext(ctx, query, newPassword, userID)
	if err != nil {
		return err
	}

	return nil
}

func (ur UserRepository) DeleteUser(ctx context.Context, userID string) error {
	query := `DELETE FROM user WHERE id = ?`
	_, err := ur.db.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}
