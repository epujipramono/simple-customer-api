package usecase

import (
	"context"
	"customer-api/internal/model"
	"customer-api/internal/repository"
)

type CustomerAPIUsecase struct {
	repo repository.CustomerAPIRepository
}

func NewCustomerAPIUsecase(repo repository.CustomerAPIRepository) *CustomerAPIUsecase {
	return &CustomerAPIUsecase{repo: repo}
}

func (u *CustomerAPIUsecase) GetAllCustomers(ctx context.Context) ([]model.GetCustomer, error) {
	return u.repo.GetAllCustomers(ctx)
}

func (u *CustomerAPIUsecase) GetCustomerByID(ctx context.Context, id int64) (*model.GetCustomer, error) {
	return u.repo.GetCustomerByID(ctx, id)
}

func (u *CustomerAPIUsecase) AddNewCustomer(ctx context.Context, p *model.Customer) error {
	return u.repo.AddNewCustomer(ctx, p)
}

func (u *CustomerAPIUsecase) UpdateCustomerByID(ctx context.Context, p *model.Customer) error {
	return u.repo.UpdateCustomerByID(ctx, p)
}

func (u *CustomerAPIUsecase) DeleteCustomerByID(ctx context.Context, id int64) error {
	return u.repo.DeleteCustomerByID(ctx, id)
}
