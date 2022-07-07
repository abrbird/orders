package sql_repository

import (
	"github.com/abrbird/orders/internal/repository"
	"github.com/jackc/pgx/v4/pgxpool"
)

type SQLRepository struct {
	dbConnectionPool *pgxpool.Pool
	orderRepository  *SQLOrderRepository
}

func New(dbConnPool *pgxpool.Pool) *SQLRepository {
	repo := &SQLRepository{
		dbConnectionPool: dbConnPool,
	}
	repo.orderRepository = &SQLOrderRepository{store: repo}
	return repo
}

func (s *SQLRepository) Order() repository.OrderRepository {
	return s.orderRepository
}
