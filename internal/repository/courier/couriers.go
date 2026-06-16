package courier

import (
	modelCourier "avito/internal/model/courier"
	"context"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Env struct {
	Host string
	Port string
	User string
	Pswd string
	Db   string
}
type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(dbPool *pgxpool.Pool) *Repository {
	return &Repository{dbPool}
}

func (db *Repository) GetById(ctx context.Context, id int) (modelCourier.Courier, error) {
	query := sq.Select(
		"id", 
		"name", 
		"phone", 
		"status", 
		"created_at", 
		"updated_at").
		From("couriers").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar)
	sql, args, err := query.ToSql()
	if err != nil {
		return modelCourier.Courier{}, fmt.Errorf("ToSql: %w", err)
	}
	courier := modelCourier.Courier{}
	err = db.pool.QueryRow(ctx, sql, args...).Scan(
		&courier.ID,
		&courier.Name,
		&courier.Phone,
		&courier.Status,
		&courier.CreatedAt,
		&courier.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return modelCourier.Courier{}, modelCourier.ErrIdNotFound
		}
		return modelCourier.Courier{}, fmt.Errorf("QueryRow: %w", err)
	}
	return courier, nil
}

func (db *Repository) Create(ctx context.Context, courier modelCourier.Courier) (int, error) {
	query := sq.Insert("couriers").
		Columns("name", "phone", "status").
		Values(courier.Name, courier.Phone, courier.Status).
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar)
	sql, args, err := query.ToSql()
	if err != nil {
		return -1, fmt.Errorf("ToSql: %w", err)
	}
	id := 0
	err = db.pool.QueryRow(ctx, sql, args...).Scan(&id)
	if err != nil {
		if isUniqueViolation(err) {
			return -1, modelCourier.ErrPhoneExists
		}
		return -1, fmt.Errorf("Exec: %w", err)
	}
	return id, nil
}

func (db *Repository) GetAll(ctx context.Context) ([]modelCourier.Courier, error) {
	query := sq.Select(
		"id",
		"name",
		"phone",
		"status",
		"created_at",
		"updated_at",
	).
		From("couriers")
	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("ToSql: %w", err)
	}
	rows, err := db.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("Query: %w", err)
	}
	defer rows.Close()
	couriers := make([]modelCourier.Courier, 0)
	for rows.Next() {
		var courier modelCourier.Courier
		if err := rows.Scan(
			&courier.ID,
			&courier.Name,
			&courier.Phone,
			&courier.Status,
			&courier.CreatedAt,
			&courier.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("Scan: %w", err)
		}
		couriers = append(couriers, courier)
	}
	return couriers, nil
}

func (db *Repository) Update(ctx context.Context, courier modelCourier.Courier) error {
	old, err := db.GetById(ctx, courier.ID)
	if err != nil {
		return fmt.Errorf("GetById: %w", err)
	}
	if old.Name == courier.Name && old.Phone == courier.Phone && old.Status == courier.Status {
		return nil
	}
	query := getUpdateQuery(courier)
	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("ToSql: %w", err)
	}
	result, err := db.pool.Exec(ctx, sql, args...)
	if err != nil {
		if isUniqueViolation(err) {
			return modelCourier.ErrPhoneExists
		}
		return fmt.Errorf("Exec: %w", err)
	}
	if result.RowsAffected() == 0 {
		return modelCourier.ErrIdNotFound
	}
	return nil
}

func (r *Repository) Delete(ctx context.Context, id int) error {
	query := sq.Delete("couriers").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar)
	sql, args, err := query.ToSql()
	res, err := r.pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("Exec: %w", err)
	}
	if res.RowsAffected() == 0 {
		return modelCourier.ErrIdNotFound
	}
	return nil
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation
}

func getUpdateQuery(courier modelCourier.Courier) sq.UpdateBuilder {
	query := sq.Update("couriers")
	if courier.Name != "" {
		query = query.Set("name", courier.Name)
	}
	if courier.Phone != "" {
		query = query.Set("phone", courier.Phone)
	}
	if courier.Status != "" {
		query = query.Set("status", courier.Status)
	}
	query = query.Set("updated_at", sq.Expr("NOW()")).
		Where(sq.Eq{"id": courier.ID}).
		PlaceholderFormat(sq.Dollar)
	return query
}
