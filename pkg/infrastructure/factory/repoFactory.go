package factory

import (
	"database/sql"
	"lunikissShop/pkg/infrastructure/mysql/repository"
)

type RepositoryFactory struct {
	db *sql.DB
}

func NewRepositoryFactory(db *sql.DB) *RepositoryFactory {
	return &RepositoryFactory{db: db}
}

func (f *RepositoryFactory) NewProductRepository() *repository.ProductRepository {
	return repository.NewProductRepository(f.db)
}

func (f *RepositoryFactory) NewSalesOutletRepository() *repository.SalesOutletRepository {
	return repository.NewSalesOutletRepository(f.db)
}

func (f *RepositoryFactory) NewOrderRepository() *repository.OrderRepository {
	return repository.NewOrderRepository(f.db)
}
