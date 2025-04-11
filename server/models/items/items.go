package items

import (
	"database/sql"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

type ItemsService struct {
	DB *sql.DB
}

type Category struct {
	ID      string `json:"id" db:"id"`
	EventID string `json:"event_id" db:"event_id"`
	Name    string `json:"name" db:"name"`
}

type Item struct {
	ID         string    `json:"id" db:"id"`
	CategoryID string    `json:"category_id" db:"category_id"`
	Name       string    `json:"name" db:"name"`
	Quantity   int       `json:"quantity" db:"quantity"`
	AssignedTo *string   `json:"assigned_to,omitempty" db:"assigned_to"`
	Status     string    `json:"status" db:"status"`
	Notes      *string   `json:"notes,omitempty" db:"notes"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}
type ItemStatusHistory struct {
	ID        string    `json:"id" db:"id"`
	ItemID    string    `json:"item_id" db:"item_id"`
	UserID    *string   `json:"user_id,omitempty" db:"user_id"`
	OldStatus string    `json:"old_status" db:"old_status"`
	NewStatus string    `json:"new_status" db:"new_status"`
	ChangedAt time.Time `json:"changed_at" db:"changed_at"`
}

func (s *ItemsService) AddItem(item *Item) error {
	item.ID = uuid.NewString()
	query := sq.Insert("items").
		Columns("id", "category_id", "name", "quantity", "assigned_to", "status", "notes").
		Values(item.ID, item.CategoryID, item.Name, item.Quantity, item.AssignedTo, item.Status, item.Notes).
		PlaceholderFormat(sq.Dollar)

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = s.DB.Exec(sqlStr, args...)
	return err
}
