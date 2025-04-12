package items

import (
	"context"
	"database/sql"
	"fmt"
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

const TableItems = "items"
const TableItemStatusHistory = "item_status_history"

func (is *ItemsService) AddItem(item *Item) error {
	item.ID = uuid.NewString()
	query := sq.Insert(TableItems).
		Columns("id", "category_id", "name", "quantity", "assigned_to", "status", "notes").
		Values(item.ID, item.CategoryID, item.Name, item.Quantity, item.AssignedTo, item.Status, item.Notes).
		PlaceholderFormat(sq.Dollar)

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("err generating add item row : %w", err)

	}

	_, err = is.DB.Exec(sqlStr, args...)
	return fmt.Errorf("err scanning category name row : %w", err)

}

func (is *ItemsService) GetItemByCategory(ctx context.Context, catID string) (*[]Item, error) {
	var items []Item
	query := sq.Select("id", "name", "quantity", "assigned_to", "status", "notes", "created_at").
		From(TableItems).Where(sq.Eq{"category_id": catID})
	sql, args, err := query.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, fmt.Errorf("err generating item by cat row : %w", err)
	}
	rows, err := is.DB.QueryContext(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("err item by cat row : %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var itm Item
		err := rows.Scan(&itm.ID, &itm.Name, &itm.Quantity, &itm.AssignedTo, &itm.Status, &itm.Notes, &itm.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("err scanning row for get item by id : %w", err)
		}
		itm.CategoryID = catID
		items = append(items, itm)
	}
	return &items, nil
}
func (is *ItemsService) UpdateItemStatus(itemID, newStatus, userID string) error {
	var oldStatus string

	// Get old status using squirrel
	selectQuery := sq.Select("status").
		From(TableItems).
		Where(sq.Eq{"id": itemID}).
		PlaceholderFormat(sq.Dollar)

	selectSQL, selectArgs, err := selectQuery.ToSql()
	if err != nil {
		return fmt.Errorf("error building select item status query: %w", err)
	}

	err = is.DB.QueryRow(selectSQL, selectArgs...).Scan(&oldStatus)
	if err != nil {
		return fmt.Errorf("error fetching old status: %w", err)
	}

	tx, err := is.DB.Begin()
	if err != nil {
		return fmt.Errorf("error beginning tx: %w", err)
	}
	defer tx.Rollback()

	// Update item status
	updateQuery := sq.Update(TableItems).
		Set("status", newStatus).
		Where(sq.Eq{"id": itemID}).
		PlaceholderFormat(sq.Dollar)

	updateSQL, updateArgs, err := updateQuery.ToSql()
	if err != nil {
		return fmt.Errorf("error building update item status query: %w", err)
	}

	_, err = tx.Exec(updateSQL, updateArgs...)
	if err != nil {
		return fmt.Errorf("error updating item status: %w", err)
	}

	// Insert status change history
	historyID := uuid.NewString()
	insertQuery := sq.Insert(TableItemStatusHistory).
		Columns("id", "item_id", "user_id", "old_status", "new_status").
		Values(historyID, itemID, userID, oldStatus, newStatus).
		PlaceholderFormat(sq.Dollar)

	insertSQL, insertArgs, err := insertQuery.ToSql()
	if err != nil {
		return fmt.Errorf("error building insert history query: %w", err)
	}

	_, err = tx.Exec(insertSQL, insertArgs...)
	if err != nil {
		return fmt.Errorf("error inserting item status history: %w", err)
	}
	return tx.Commit()
}

func (is *ItemsService) AssignItem(itemID, userID string) error {
	query := sq.Update(TableItems).Set("assigned_to", userID).Where(sq.Eq{"id": itemID})
	sql, args, err := query.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return fmt.Errorf("error assigning item to user query generation: %w", err)
	}
	_, err = is.DB.Exec(sql, args...)
	return err
}

func (is *ItemsService) UnassignItem(itemID string) error {
	query := sq.Update(TableItems).Set("assigned_to", nil).Where(sq.Eq{"id": itemID})
	sql, args, err := query.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return fmt.Errorf("error unassigning item  generation: %w", err)
	}
	_, err = is.DB.Exec(sql, args...)
	return err
}

func (is *ItemsService) EditItem(ctx context.Context, itemID string, updates map[string]interface{}) error {
	if len(updates) == 0 {
		return fmt.Errorf("no fields to update")
	}

	query := sq.Update(TableItems).
		SetMap(updates).
		Where(sq.Eq{"id": itemID}).
		PlaceholderFormat(sq.Dollar)

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("error building update query: %w", err)
	}

	_, err = is.DB.ExecContext(ctx, sqlStr, args...)
	if err != nil {
		return fmt.Errorf("error executing item update: %w", err)
	}

	return nil
}

func (is *ItemsService) DeleteItem(id string) error {
	query := sq.Delete(TableItems).Where(sq.Eq{"id": id})
	sql, args, err := query.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return fmt.Errorf("error generating delete item query: %w", err)
	}
	_, err = is.DB.Exec(sql, args...)
	return err
}
