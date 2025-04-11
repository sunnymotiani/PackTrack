package templates

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
)

type TemplatesService struct {
	DB *sql.DB
}

type Template struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type TemplateItem struct {
	ID         string `json:"id"`
	TemplateID string `json:"template_id"`
	Category   string `json:"category"`
	Name       string `json:"name"`
	Quantity   int    `json:"quantity"`
}

const TableTemplates = "templates"
const TableTemplateItems = "template_items"

func (s *TemplatesService) CreateTemplate(ctx context.Context, name string) (*Template, error) {
	query := sq.
		Insert(TableTemplates).
		Columns("name").
		Values(name).
		Suffix("RETURNING id, name, created_at").
		PlaceholderFormat(sq.Dollar)

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("create template sql error: %w", err)
	}

	t := &Template{}
	err = s.DB.QueryRowContext(ctx, sqlStr, args...).
		Scan(&t.ID, &t.Name, &t.CreatedAt)
	return t, err
}

func (s *TemplatesService) AddItemToTemplate(ctx context.Context, templateID, category, name string, quantity int) (*TemplateItem, error) {
	query := sq.
		Insert(TableTemplateItems).
		Columns("template_id", "category", "name", "quantity").
		Values(templateID, category, name, quantity).
		Suffix("RETURNING id, template_id, category, name, quantity").
		PlaceholderFormat(sq.Dollar)

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("add item sql error: %w", err)
	}

	item := &TemplateItem{}
	err = s.DB.QueryRowContext(ctx, sqlStr, args...).
		Scan(&item.ID, &item.TemplateID, &item.Category, &item.Name, &item.Quantity)
	return item, err
}

func (s *TemplatesService) GetTemplateByID(ctx context.Context, id string) (*Template, []TemplateItem, error) {
	query := sq.
		Select("id", "name", "created_at").
		From(TableTemplates).
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar)

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return nil, nil, fmt.Errorf("get template sql error: %w", err)
	}

	t := &Template{}
	err = s.DB.QueryRowContext(ctx, sqlStr, args...).
		Scan(&t.ID, &t.Name, &t.CreatedAt)
	if err != nil {
		return nil, nil, err
	}

	// Now fetch items
	itemsQuery := sq.
		Select("id", "template_id", "category", "name", "quantity").
		From(TableTemplateItems).
		Where(sq.Eq{"template_id": id}).
		PlaceholderFormat(sq.Dollar)

	sqlStr, args, err = itemsQuery.ToSql()
	if err != nil {
		return nil, nil, fmt.Errorf("get template items sql error: %w", err)
	}

	rows, err := s.DB.QueryContext(ctx, sqlStr, args...)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	var items []TemplateItem
	for rows.Next() {
		var i TemplateItem
		if err := rows.Scan(&i.ID, &i.TemplateID, &i.Category, &i.Name, &i.Quantity); err != nil {
			return nil, nil, err
		}
		items = append(items, i)
	}

	return t, items, nil
}

func (s *TemplatesService) DeleteTemplate(ctx context.Context, id string) error {
	query := sq.
		Delete(TableTemplates).
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar)

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = s.DB.ExecContext(ctx, sqlStr, args...)
	return err
}

func (s *TemplatesService) DeleteItem(ctx context.Context, itemID string) error {
	query := sq.
		Delete(TableTemplateItems).
		Where(sq.Eq{"id": itemID}).
		PlaceholderFormat(sq.Dollar)

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = s.DB.ExecContext(ctx, sqlStr, args...)
	return err
}
