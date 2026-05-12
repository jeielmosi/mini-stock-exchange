package main

import (
	"context"
	"fmt"
	order_heaps "mini-stock-exchange/internal/domain-service/match-service/order-heap"
	"mini-stock-exchange/internal/entity"
	"mini-stock-exchange/internal/repository"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

func lessQt(lhs *order_heaps.QueryTrigger, rhs *order_heaps.QueryTrigger) bool {
	if lhs == nil {
		return true
	}
	if rhs == nil {
		return false
	}

	if lhs.Price.Equal(rhs.Price) {
		return lhs.CreatedAt.Before(rhs.CreatedAt)
	}
	return lhs.Price.LessThan(rhs.Price)
}

func lessOrder(lhs *entity.Order, rhs *entity.Order) bool {
	if lhs == nil {
		return true
	}
	if rhs == nil {
		return false
	}

	return lessQt(order_heaps.NewQueryTrigger(*lhs), order_heaps.NewQueryTrigger(*rhs))
}

func main() {
	symbol := "AAPL"
	now := time.Now().UTC()
	ctx := context.Background()
	db, cleanup, err := repository.SetupTestDB(ctx)
	defer cleanup()
	repo, err := repository.NewOrderRepository(db)
	if err != nil {
		panic(err)
	}
	defer repo.Stop()
	brokerRepo, err := repository.NewBrokerRepository(db)
	if err != nil {
		panic(err)
	}

	b0 := uuid.New()
	b1 := uuid.New()
	b2 := uuid.New()
	err = brokerRepo.Insert(entity.Broker{ID: b0, Name: "broker0"})
	if err != nil {
		panic(err)
	}
	err = brokerRepo.Insert(entity.Broker{ID: b1, Name: "broker1"})
	if err != nil {
		panic(err)
	}
	err = brokerRepo.Insert(entity.Broker{ID: b2, Name: "broker2"})
	if err != nil {
		panic(err)
	}

	ah := order_heaps.NewAskHeap(symbol, 2, repo)

	// Order in repo that should be picked up by fill
	o0 := entity.Order{
		ID:                uuid.New(),
		BrokerID:          b0,
		OwnerDoc:          "doc0",
		Symbol:            symbol,
		Price:             decimal.NewFromInt(5),
		CreatedAt:         now,
		Type:              entity.Ask,
		ValidUntil:        now.Add(24 * time.Hour),
		Status:            entity.Pending,
		Quantity:          100,
		RemainingQuantity: 100,
	}
	ah.Push(o0)
	err = repo.Insert(o0)
	if err != nil {
		panic(err)
	}

	o1 := entity.Order{
		ID:                uuid.New(),
		BrokerID:          b1,
		OwnerDoc:          "doc1",
		Symbol:            symbol,
		Price:             decimal.NewFromInt(10),
		CreatedAt:         now,
		Type:              entity.Ask,
		ValidUntil:        now.Add(24 * time.Hour),
		Status:            entity.Pending,
		Quantity:          100,
		RemainingQuantity: 100,
	}
	ah.Push(o1)
	err = repo.Insert(o1)
	if err != nil {
		panic(err)
	}

	o2 := entity.Order{
		ID:                uuid.New(),
		BrokerID:          b2,
		OwnerDoc:          "doc2",
		Symbol:            symbol,
		Price:             decimal.NewFromInt(20),
		CreatedAt:         now,
		Type:              entity.Ask,
		ValidUntil:        now.Add(24 * time.Hour),
		Status:            entity.Pending,
		Quantity:          100,
		RemainingQuantity: 100,
	}
	ah.Push(o2)
	err = repo.Insert(o2)
	if err != nil {
		panic(err)
	}

	// Pop should trigger fill because qt(10) < top(20)
	match := entity.Order{
		ID:                uuid.New(),
		Symbol:            symbol,
		Price:             decimal.NewFromInt(200),
		Quantity:          100,
		RemainingQuantity: 100,
		Type:              entity.Bid,
		OwnerDoc:          "match",
		BrokerID:          uuid.New(),
	}

	res, err := ah.Pop(match)
	if err != nil {
		panic(err)
	}
	rm := []uuid.UUID{}
	for _, m := range res.Matches {
		rm = append(rm, m.ID)
	}
	err = repo.Expire(rm)
	if err != nil {
		panic(err)
	}
	fmt.Println("matches: ", res.Matches)
	fmt.Println("expired: ", res.Expired)
	fmt.Println(res.Matches[0].ID == o0.ID, "match: ", res.Matches[0].OwnerDoc, "order: ", o0.OwnerDoc)
	fmt.Println()

	res, err = ah.Pop(match)
	if err != nil {
		panic(err)
	}
	rm = []uuid.UUID{}
	for _, m := range res.Matches {
		rm = append(rm, m.ID)
	}
	err = repo.Expire(rm)
	if err != nil {
		panic(err)
	}
	fmt.Println("matches: ", res.Matches)
	fmt.Println("expired: ", res.Expired)
	fmt.Println(res.Matches[0].ID == o1.ID, "match: ", res.Matches[0].OwnerDoc, "order: ", o1.OwnerDoc)
	fmt.Println()

	res, err = ah.Pop(match)
	if err != nil {
		panic(err)
	}
	rm = []uuid.UUID{}
	for _, m := range res.Matches {
		rm = append(rm, m.ID)
	}
	err = repo.Expire(rm)
	if err != nil {
		panic(err)
	}
	fmt.Println("matches: ", res.Matches)
	fmt.Println("expired: ", res.Expired)
	fmt.Println(res.Matches[0].ID == o2.ID, "match: ", res.Matches[0].OwnerDoc, "order: ", o2.OwnerDoc)
	fmt.Println()
}
