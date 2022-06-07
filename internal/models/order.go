package models

const (
	Created       string = "created"
	Paid                 = "paid"
	Shipping             = "shipping"
	OnIssuePoint         = "on_issue_point"
	ReadyForIssue        = "ready_for_issue"
	Issued               = "issued"
)

type Order struct {
	Id     int64
	Status string
}

type OrderRetrieve struct {
	Order *Order
	Error error
}

type OrderItem struct {
	Id      int64
	OrderId int64
	Volume  float64
}

type OrderItemRetrieve struct {
	OrderItem *OrderItem
	Error     error
}

type OrderItemsRetrieve struct {
	OrderItems []OrderItem
	Error      error
}
