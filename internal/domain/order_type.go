package domain

type OrderType string

const (
	Bid OrderType = "BID"
	Ask OrderType = "ASK"
)

type OrderStatus string

const (
	Pending OrderStatus = "PENDING"
	Partial OrderStatus = "PARTIAL"
	Filled  OrderStatus = "FILLED"
	Expired OrderStatus = "EXPIRED"
)
