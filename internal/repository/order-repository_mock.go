package repository

import "context"

func NewMockOrderRepository() (OrderRepository, func()) {
	ctx := context.Background()
	sql, cleanup := SetupTestDB(ctx)
	return NewOrderRepository(sql), cleanup
}
