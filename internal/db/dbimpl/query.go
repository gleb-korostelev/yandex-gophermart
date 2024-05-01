package dbimpl

import (
	"context"
	"time"

	"github.com/gleb-korostelev/gophermart.git/internal/config"
	"github.com/gleb-korostelev/gophermart.git/internal/db"
	"github.com/gleb-korostelev/gophermart.git/internal/models"
)

func GetUserCred(db db.DB, ctx context.Context, login string) (string, error) {
	var password string
	var isDeleted bool
	sql := `SELECT password FROM user_data WHERE login = $1`
	err := db.QueryRow(ctx, sql, login).Scan(&password)
	if err != nil {
		return "", err
	}
	if isDeleted {
		return "", config.ErrGone
	}
	return password, nil
}

func GetOrders(db db.DB, ctx context.Context, login string) ([]models.OrdersData, error) {
	sql := `
	SELECT order_id, status, accrual, uploaded_at
	FROM orders
	WHERE login=$1
	ORDER BY uploaded_at ASC
	`
	rows, err := db.Query(context.Background(), sql, login)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []models.OrdersData
	for rows.Next() {
		var order models.OrdersData
		var accrual *float64
		var uploadedAt time.Time
		if err := rows.Scan(&order.Number, &order.Status, &accrual, &uploadedAt); err != nil {
			return orders, err
		}
		if accrual != nil {
			order.Accrual = *accrual
		}
		order.UploadedAt = uploadedAt.Format(time.RFC3339)
		orders = append(orders, order)
	}
	return orders, nil
}

func Balance(db db.DB, ctx context.Context, login string) (models.BalanceData, error) {
	sql := `
	SELECT current, withdrawn
	FROM balances
	WHERE login=$1
	`
	var balance models.BalanceData
	err := db.QueryRow(ctx, sql, login).Scan(&balance.Current, &balance.Withdrawn)
	if err != nil {
		return models.BalanceData{}, err
	}
	return balance, nil
}

func GetWithdrawals(db db.DB, ctx context.Context, login string) ([]models.Withdraws, error) {
	sql := `
	SELECT order_id, sum, processed_at
	FROM withdrawals
	WHERE login=$1
	`
	rows, err := db.Query(ctx, sql, login)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var withdraws []models.Withdraws
	for rows.Next() {
		var withdraw models.Withdraws
		var processedAt time.Time
		if err := rows.Scan(&withdraw.Order, &withdraw.Sum, &processedAt); err != nil {
			return nil, err
		}
		withdraw.ProcessedAt = processedAt.Format(time.RFC3339)
		withdraws = append(withdraws, withdraw)
	}
	return withdraws, nil
}

func GetOrderByNumber(db db.DB, ctx context.Context, orderInfo models.OrderResponse, resultChan chan<- models.OrderResponse) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		var response models.OrderResponse
		sql := `UPDATE orders SET status = $1, accrual = $2 WHERE order_number = $3`
		_, err := db.Exec(ctx, sql, orderInfo.Status, orderInfo.Accrual, orderInfo.Order)
		if err != nil {
			return err
		}
		select {
		case resultChan <- response:
		case <-ctx.Done():
			return ctx.Err()
		}

		return nil
	}
}
