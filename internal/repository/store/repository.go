package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	gen_store "github.com/oapi-codegen-multiple-packages-example/internal/gen/store"
)

type Repository interface {
	CreateOrder(ctx context.Context, order gen_store.Order) (*gen_store.Order, error)
	GetOrderByID(ctx context.Context, orderID int64) (*gen_store.Order, error)
	DeleteOrder(ctx context.Context, orderID int64) error

	GetInventory(ctx context.Context) (map[string]int32, error)
}

type mysqlRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &mysqlRepository{db: db}
}

func (r *mysqlRepository) CreateOrder(ctx context.Context, order gen_store.Order) (*gen_store.Order, error) {
	query := `INSERT INTO orders (pet_id, quantity, ship_date, status, complete) 
              VALUES (?, ?, ?, ?, ?)`

	var (
		petID    *int64
		quantity *int32
		shipDate *time.Time
		status   *string
		complete *bool
	)

	if order.PetId != nil {
		petID = order.PetId
	}

	if order.Quantity != nil {
		quantity = order.Quantity
	}

	if order.ShipDate != nil {
		shipDate = order.ShipDate
	}

	if order.Status != nil {
		orderStatus := string(*order.Status)
		status = &orderStatus
	}

	if order.Complete != nil {
		complete = order.Complete
	}

	result, err := r.db.ExecContext(ctx, query, petID, quantity, shipDate, status, complete)
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert ID: %w", err)
	}

	order.Id = &id
	return &order, nil
}

func (r *mysqlRepository) GetOrderByID(ctx context.Context, orderID int64) (*gen_store.Order, error) {
	query := `SELECT id, pet_id, quantity, ship_date, status, complete FROM orders WHERE id = ?`

	var (
		id       int64
		petID    sql.NullInt64
		quantity sql.NullInt32
		shipDate sql.NullTime
		status   sql.NullString
		complete sql.NullBool
	)

	err := r.db.QueryRowContext(ctx, query, orderID).Scan(&id, &petID, &quantity, &shipDate, &status, &complete)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	order := gen_store.Order{
		Id: &id,
	}

	if petID.Valid {
		order.PetId = &petID.Int64
	}

	if quantity.Valid {
		order.Quantity = &quantity.Int32
	}

	if shipDate.Valid {
		order.ShipDate = &shipDate.Time
	}

	if status.Valid {
		orderStatus := gen_store.OrderStatus(status.String)
		order.Status = &orderStatus
	}

	if complete.Valid {
		order.Complete = &complete.Bool
	}

	return &order, nil
}

func (r *mysqlRepository) DeleteOrder(ctx context.Context, orderID int64) error {
	query := `DELETE FROM orders WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query, orderID)
	if err != nil {
		return fmt.Errorf("failed to delete order: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("order not found")
	}

	return nil
}

func (r *mysqlRepository) GetInventory(ctx context.Context) (map[string]int32, error) {
	query := `SELECT status, COUNT(*) as count FROM orders GROUP BY status`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get inventory: %w", err)
	}
	defer rows.Close()

	inventory := make(map[string]int32)

	for rows.Next() {
		var status string
		var count int32

		if err := rows.Scan(&status, &count); err != nil {
			return nil, fmt.Errorf("failed to scan inventory row: %w", err)
		}

		inventory[status] = count
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating inventory rows: %w", err)
	}

	return inventory, nil
}
