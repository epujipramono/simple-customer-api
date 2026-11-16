package repository

import (
	"context"
	"customer-api/internal/model"
	"database/sql"
	"errors"
	"fmt"
)

type CustomerAPIRepository interface {
	GetAllCustomers(ctx context.Context) ([]model.GetCustomer, error)
	GetCustomerByID(ctx context.Context, id int64) (*model.GetCustomer, error)
	AddNewCustomer(ctx context.Context, p *model.Customer) error
	UpdateCustomerByID(ctx context.Context, p *model.Customer) error
	DeleteCustomerByID(ctx context.Context, id int64) error
}

type mysqlCustomerAPIRepository struct {
	db *sql.DB
}

func NewCustomerAPIRepository(db *sql.DB) CustomerAPIRepository {
	return &mysqlCustomerAPIRepository{db: db}
}

func (r *mysqlCustomerAPIRepository) GetAllCustomers(ctx context.Context) ([]model.GetCustomer, error) {
	query := `
    	SELECT customers.id AS ID, customers.name AS Name, customers.email AS Email, customers.phone AS Phone
    	FROM customers
	`
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]model.GetCustomer, 0)
	for rows.Next() {
		var customer model.GetCustomer
		if err := rows.Scan(&customer.ID, &customer.Name, &customer.Email, &customer.Phone); err != nil {
			return nil, err
		}
		result = append(result, customer)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, rows.Err()
}

func (r *mysqlCustomerAPIRepository) GetCustomerByID(ctx context.Context, id int64) (*model.GetCustomer, error) {
	query := fmt.Sprintf(`
    	SELECT customers.id AS ID, customers.name AS Name, customers.email AS Email, customers.phone AS Phone
    	FROM customers
		WHERE id = '%d'
	`, id)
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var customer model.GetCustomer
	err = stmt.QueryRowContext(ctx).Scan(&customer.ID, &customer.Name, &customer.Email, &customer.Phone)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("Customer not found")
	}
	if err != nil {
		return nil, err
	}
	return &customer, nil
}

func (r *mysqlCustomerAPIRepository) AddNewCustomer(ctx context.Context, p *model.Customer) error {
	err := r.checkEmail(p.Email)
	if err != nil {
		return err
	}

	query := "INSERT INTO customers(name, email, phone, created_at) VALUES (?, ?, ?, NOW())"
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, p.Name, p.Email, p.Phone)
	if err != nil {
		return err
	}
	p.ID, err = result.LastInsertId()
	return err
}

func (r *mysqlCustomerAPIRepository) UpdateCustomerByID(ctx context.Context, p *model.Customer) error {

	var oldEmail string
	queryOldEmail := fmt.Sprintf(`SELECT customers.email FROM customers where customers.id = '%d'`, p.ID)
	err := r.db.QueryRow(queryOldEmail).Scan(&oldEmail)
	if err != nil {
		return err
	}
	if oldEmail != p.Email {
		err := r.checkEmail(p.Email)
		if err != nil {
			return err
		}
	}

	query := "UPDATE customers SET customers.name = ?, customers.email = ?, customers.phone = ?, customers.updated_at = NOW() WHERE customers.id = ?"
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, p.Name, p.Email, p.Phone, p.ID)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return errors.New("Customer not found")
	}
	return nil
}

func (r *mysqlCustomerAPIRepository) DeleteCustomerByID(ctx context.Context, id int64) error {
	query := "DELETE FROM customers WHERE customers.id = ?"
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, id)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return errors.New("Customer not found")
	}
	return nil
}

func (r *mysqlCustomerAPIRepository) checkEmail(email string) error {
	var count int
	queryCheckEmail := fmt.Sprintf(`SELECT COUNT(*) FROM customers where customers.email = '%s'`, email)
	err := r.db.QueryRow(queryCheckEmail).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("Duplicate email")
	}
	return nil
}
