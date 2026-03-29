package domain

import "time"

type User struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u *User) Fields() []interface{} {
	return []interface{}{
		&u.ID,
		&u.Name,
		&u.Email,
		&u.CreatedAt,
		&u.UpdatedAt,
	}
}

type UserTable struct{}

func (UserTable) TableName() string {
	return "users"
}

func (UserTable) Columns() []string {
	return []string{
		"users.id",
		"users.name",
		"users.email",
		"users.created_at",
		"users.updated_at",
	}
}
