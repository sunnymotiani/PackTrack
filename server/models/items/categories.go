package items

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
)

//CATEGORIES IMPLIMENTATION

const TableCategories = "categories"

func (is *ItemsService) CreateCategory(ctx context.Context, eventID string, name string) (*Category, error) {
	var category Category
	query := sq.Insert(TableCategories).Columns("event_id", "name").Values(eventID, name).
		Suffix("RETURNING id")
	sql, args, err := query.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, fmt.Errorf("create category sql error: %w", err)
	}
	row := is.DB.QueryRowContext(ctx, sql, args...)
	err = row.Scan(&category.ID)
	category.Name = name
	category.EventID = eventID
	return &category, nil
}

func (is *ItemsService) GetCategories(ctx context.Context, eventID string) (*[]Category, error) {
	var categories []Category
	query := sq.Select("id", "name").From(TableCategories).Where(sq.Eq{"event_id": eventID})
	sql, args, err := query.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, fmt.Errorf("get categories sql error: %w", err)
	}
	rows, err := is.DB.QueryContext(ctx, sql, args)
	if err != nil {
		return nil, fmt.Errorf("get categories sql query error: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var cat Category
		cat.EventID = eventID
		err := rows.Scan(&cat.ID, &cat.Name)
		if err != nil {
			return nil, fmt.Errorf("get categories err scanning rows : %w", err)
		}
		categories = append(categories, cat)
	}
	return &categories, nil
}

func (is *ItemsService) DeleteCategory(ctx context.Context, id string) error {
	query := sq.Delete(TableCategories).Where(sq.Eq{"id": id})
	sql, args, err := query.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return fmt.Errorf("err generating delete category row : %w", err)
	}
	_, err = is.DB.ExecContext(ctx, sql, args...)
	return err
}

func (is *ItemsService) UpdateCategoryName(ctx context.Context, id string, name string) error {
	query := sq.Update(TableCategories).Set("name", name).Where(sq.Eq{"id": id})
	sql, args, err := query.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return fmt.Errorf("err generating update category name row : %w", err)
	}
	_, err = is.DB.ExecContext(ctx, sql, args...)
	return err
}
