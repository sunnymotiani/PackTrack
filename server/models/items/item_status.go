package items

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

func (is *ItemsService) LogStatusChange(ctx context.Context, itemID, userID, oldStatus, newStatus string) error {
	id := uuid.NewString()

	query := sq.Insert(TableItemStatusHistory).
		Columns("id", "item_id", "user_id", "old_status", "new_status", "changed_at").
		Values(id, itemID, userID, oldStatus, newStatus, time.Now()).
		PlaceholderFormat(sq.Dollar)

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("error building log insert query: %w", err)
	}

	_, err = is.DB.ExecContext(ctx, sqlStr, args...)
	if err != nil {
		return fmt.Errorf("error inserting item status log: %w", err)
	}

	return nil
}

func (is *ItemsService) GetItemHistory(ctx context.Context, itemID string) ([]ItemStatusHistory, error) {
	query := sq.Select("id", "item_id", "user_id", "old_status", "new_status", "changed_at").
		From(TableItemStatusHistory).
		Where(sq.Eq{"item_id": itemID}).
		OrderBy("changed_at ASC").
		PlaceholderFormat(sq.Dollar)

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("error building select history query: %w", err)
	}

	rows, err := is.DB.QueryContext(ctx, sqlStr, args...)
	if err != nil {
		return nil, fmt.Errorf("error querying item history: %w", err)
	}
	defer rows.Close()

	var history []ItemStatusHistory

	for rows.Next() {
		var record ItemStatusHistory
		var userID sql.NullString

		if err := rows.Scan(
			&record.ID,
			&record.ItemID,
			&userID,
			&record.OldStatus,
			&record.NewStatus,
			&record.ChangedAt,
		); err != nil {
			return nil, fmt.Errorf("error scanning history row: %w", err)
		}

		if userID.Valid {
			record.UserID = &userID.String
		}

		history = append(history, record)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return history, nil
}
