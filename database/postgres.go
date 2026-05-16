package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
	"github.com/cr1m1/expense-tracker-service/models"
)

type DB struct {
	conn *sql.DB
}

func Connect(dsn string) (*DB, error) {
	conn, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	db := &DB{conn: conn}
	if err := db.createTable(); err != nil {
		return nil, err
	}

	return db, nil
}

func (db *DB) createTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS expenses (
		id SERIAL PRIMARY KEY,
		amount DECIMAL(10, 2) NOT NULL,
		category VARCHAR(100) NOT NULL,
		description TEXT,
		date TIMESTAMP NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)
	`
	_, err := db.conn.Exec(query)
	return err
}

func (db *DB) CreateExpense(expense *models.Expense) (int, error) {
	var id int
	query := `INSERT INTO expenses (amount, category, description, date, created_at)
	          VALUES ($1, $2, $3, $4, $5) RETURNING id`
	err := db.conn.QueryRow(query, expense.Amount, expense.Category, expense.Description, expense.Date, time.Now()).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to create expense: %w", err)
	}
	return id, nil
}

func (db *DB) GetAllExpenses() ([]models.Expense, error) {
	query := `SELECT id, amount, category, description, date, created_at FROM expenses ORDER BY date DESC`
	rows, err := db.conn.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch expenses: %w", err)
	}
	defer rows.Close()

	var expenses []models.Expense
	for rows.Next() {
		var e models.Expense
		if err := rows.Scan(&e.ID, &e.Amount, &e.Category, &e.Description, &e.Date, &e.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan expense: %w", err)
		}
		expenses = append(expenses, e)
	}
	return expenses, nil
}

func (db *DB) GetExpenseByID(id int) (*models.Expense, error) {
	query := `SELECT id, amount, category, description, date, created_at FROM expenses WHERE id = $1`
	var e models.Expense
	err := db.conn.QueryRow(query, id).Scan(&e.ID, &e.Amount, &e.Category, &e.Description, &e.Date, &e.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to fetch expense: %w", err)
	}
	return &e, nil
}

func (db *DB) UpdateExpense(id int, expense *models.Expense) error {
	query := `UPDATE expenses SET amount = $1, category = $2, description = $3, date = $4 WHERE id = $5`
	result, err := db.conn.Exec(query, expense.Amount, expense.Category, expense.Description, expense.Date, id)
	if err != nil {
		return fmt.Errorf("failed to update expense: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check affected rows: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("expense not found")
	}
	return nil
}

func (db *DB) DeleteExpense(id int) error {
	query := `DELETE FROM expenses WHERE id = $1`
	result, err := db.conn.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete expense: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check affected rows: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("expense not found")
	}
	return nil
}

func (db *DB) Close() error {
	return db.conn.Close()
}
