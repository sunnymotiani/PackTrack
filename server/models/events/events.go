package events

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/sunnymotiani/PackTrack/server/models/users"
)

const (
	TableEvents           = "events"
	TableEventMemberships = "event_memberships"
)

type Event struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	OwnerID     *string   `json:"owner_id,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

type EventMembership struct {
	ID      string `json:"id"`
	UserID  string `json:"user_id"`
	EventID string `json:"event_id"`
	Role    string `json:"role"` // owner, admin, member, viewer
	Name    string `json:"name"`
	Email   string `json:"email"`
}

type EventService struct {
	DB *sql.DB
}

func (es *EventService) CreateEvent(ctx context.Context, name, description, ownerID string) (*Event, error) {
	eventID := uuid.NewString()
	now := time.Now()

	// Start transaction
	tx, err := es.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Insert into events
	insertEvent := sq.Insert(TableEvents).
		Columns("id", "name", "description", "owner_id", "created_at").
		Values(eventID, name, description, ownerID, now).
		PlaceholderFormat(sq.Dollar)

	sql1, args1, err := insertEvent.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build event insert: %w", err)
	}
	if _, err := tx.ExecContext(ctx, sql1, args1...); err != nil {
		return nil, fmt.Errorf("insert event: %w", err)
	}

	// Insert into event_memberships
	membershipID := uuid.NewString()
	insertMembership := sq.Insert(TableEventMemberships).
		Columns("id", "user_id", "event_id", "role").
		Values(membershipID, ownerID, eventID, "owner").
		PlaceholderFormat(sq.Dollar)

	sql2, args2, err := insertMembership.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build membership insert: %w", err)
	}
	if _, err := tx.ExecContext(ctx, sql2, args2...); err != nil {
		return nil, fmt.Errorf("insert event membership: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &Event{
		ID:          eventID,
		Name:        name,
		Description: description,
		OwnerID:     &ownerID,
		CreatedAt:   now,
	}, nil
}
func (es *EventService) AddMember(ctx context.Context, eventID, userID, role string) error {
	// Validate role
	validRoles := map[string]bool{
		"owner":  true,
		"admin":  true,
		"member": true,
		"viewer": true,
	}
	if !validRoles[role] {
		return fmt.Errorf("invalid role: %s", role)
	}

	query := sq.Insert(TableEventMemberships).
		Columns("id", "user_id", "event_id", "role").
		Values(uuid.NewString(), userID, eventID, role).
		PlaceholderFormat(sq.Dollar)

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("error building AddMember query: %w", err)
	}

	_, err = es.DB.ExecContext(ctx, sqlStr, args...)
	if err != nil {
		return fmt.Errorf("error executing AddMember query: %w", err)
	}

	return nil
}
func (es *EventService) GetEventMembers(ctx context.Context, eventID string) ([]EventMembership, error) {
	query := sq.
		Select("em.user_id", "em.role", "u.name", "u.email").
		From(TableEventMemberships + " em").
		Join(users.TableUsers + " u ON em.user_id = u.id").
		Where(sq.Eq{"em.event_id": eventID}).
		PlaceholderFormat(sq.Dollar)

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("error building GetEventMembers query: %w", err)
	}

	rows, err := es.DB.QueryContext(ctx, sqlStr, args...)
	if err != nil {
		return nil, fmt.Errorf("error executing GetEventMembers query: %w", err)
	}
	defer rows.Close()

	var members []EventMembership
	for rows.Next() {
		var m EventMembership
		if err := rows.Scan(&m.UserID, &m.Role, &m.Name, &m.Email); err != nil {
			return nil, fmt.Errorf("error scanning member row: %w", err)
		}
		members = append(members, m)
	}
	return members, nil
}
func (es *EventService) UpdateMemberRole(ctx context.Context, eventID, userID, newRole string) error {
	query := sq.
		Update(TableEventMemberships).
		Set("role", newRole).
		Where(sq.Eq{"event_id": eventID, "user_id": userID}).
		PlaceholderFormat(sq.Dollar)

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("error building UpdateMemberRole query: %w", err)
	}

	_, err = es.DB.ExecContext(ctx, sqlStr, args...)
	if err != nil {
		return fmt.Errorf("error executing UpdateMemberRole query: %w", err)
	}

	return nil
}
func (es *EventService) RemoveMember(ctx context.Context, eventID, userID string) error {
	query := sq.
		Delete(TableEventMemberships).
		Where(sq.Eq{"event_id": eventID, "user_id": userID}).
		PlaceholderFormat(sq.Dollar)

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("error building RemoveMember query: %w", err)
	}

	_, err = es.DB.ExecContext(ctx, sqlStr, args...)
	if err != nil {
		return fmt.Errorf("error executing RemoveMember query: %w", err)
	}

	return nil
}
func (es *EventService) GetEventsForUser(ctx context.Context, userID string) ([]*Event, error) {
	query := sq.
		Select("e.id", "e.name", "e.description", "e.owner_id", "e.created_at").
		From(fmt.Sprintf("%s AS e", TableEvents)).
		Join(fmt.Sprintf("%s AS em ON em.event_id = e.id", TableEventMemberships)).
		Where(sq.Eq{"em.user_id": userID}).
		PlaceholderFormat(sq.Dollar)

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("error building GetEventsForUser query: %w", err)
	}

	rows, err := es.DB.QueryContext(ctx, sqlStr, args...)
	if err != nil {
		return nil, fmt.Errorf("error executing GetEventsForUser query: %w", err)
	}
	defer rows.Close()

	var events []*Event
	for rows.Next() {
		var e Event
		if err := rows.Scan(&e.ID, &e.Name, &e.Description, &e.OwnerID, &e.CreatedAt); err != nil {
			return nil, fmt.Errorf("error scanning event row: %w", err)
		}
		events = append(events, &e)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row error in GetEventsForUser: %w", err)
	}

	return events, nil
}
