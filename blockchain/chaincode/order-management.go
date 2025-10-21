package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// OrderContract 订单管理智能合约
type OrderContract struct {
	contractapi.Contract
}

// Order 订单结构
type Order struct {
	ID                string    `json:"id"`
	OrderNumber       string    `json:"orderNumber"`
	UserID            string    `json:"userId"`
	Status            string    `json:"status"`
	Amount            float64   `json:"amount"`
	Currency          string    `json:"currency"`
	PaymentMethod     string    `json:"paymentMethod"`
	PaymentID         string    `json:"paymentId"`
	PaymentStatus     string    `json:"paymentStatus"`
	Items             string    `json:"items"` // JSON string
	ShippingAddress   string    `json:"shippingAddress"`
	BillingAddress    string    `json:"billingAddress"`
	TrackingNumber    string    `json:"trackingNumber"`
	ShippingMethod    string    `json:"shippingMethod"`
	EstimatedDelivery time.Time `json:"estimatedDelivery"`
	ActualDelivery    time.Time `json:"actualDelivery"`
	ChainEventHash    string    `json:"chainEventHash"`
	RefundHash        string    `json:"refundHash"`
	RefundReason      string    `json:"refundReason"`
	Metadata          string    `json:"metadata"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
}

// OrderEvent 订单事件结构
type OrderEvent struct {
	ID          string    `json:"id"`
	OrderID     string    `json:"orderId"`
	EventType   string    `json:"eventType"`
	PayloadHash string    `json:"payloadHash"`
	Timestamp   time.Time `json:"timestamp"`
	Description string    `json:"description"`
	Metadata    string    `json:"metadata"`
}

// RefundPolicy 退款政策结构
type RefundPolicy struct {
	ID                string    `json:"id"`
	OrderID           string    `json:"orderId"`
	RefundWindow      int       `json:"refundWindow"` // 天数
	RefundPercentage  float64   `json:"refundPercentage"`
	Conditions        string    `json:"conditions"`
	IsActive          bool      `json:"isActive"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
}

// CreateOrder 创建订单
func (oc *OrderContract) CreateOrder(ctx contractapi.TransactionContextInterface, orderData string) error {
	var order Order
	err := json.Unmarshal([]byte(orderData), &order)
	if err != nil {
		return fmt.Errorf("failed to unmarshal order data: %v", err)
	}

	// 验证必需字段
	if order.OrderNumber == "" || order.UserID == "" || order.Amount <= 0 {
		return fmt.Errorf("missing required fields")
	}

	// 设置默认值
	order.Status = "pending"
	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()

	// 生成唯一ID
	order.ID = fmt.Sprintf("order_%s_%d", order.OrderNumber, order.CreatedAt.Unix())

	// 存储订单
	orderJSON, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("failed to marshal order: %v", err)
	}

	err = ctx.GetStub().PutState(order.ID, orderJSON)
	if err != nil {
		return fmt.Errorf("failed to put order: %v", err)
	}

	// 创建复合键用于查询
	compositeKey, err := ctx.GetStub().CreateCompositeKey("order", []string{order.UserID, order.Status, order.ID})
	if err != nil {
		return fmt.Errorf("failed to create composite key: %v", err)
	}

	err = ctx.GetStub().PutState(compositeKey, orderJSON)
	if err != nil {
		return fmt.Errorf("failed to put composite key: %v", err)
	}

	// 创建订单事件
	event := OrderEvent{
		ID:          fmt.Sprintf("event_%s_%d", order.ID, time.Now().Unix()),
		OrderID:     order.ID,
		EventType:   "CREATED",
		PayloadHash: fmt.Sprintf("hash_%s", order.ID),
		Timestamp:   time.Now(),
		Description: "Order created",
		Metadata:    "",
	}

	eventJSON, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %v", err)
	}

	err = ctx.GetStub().PutState(event.ID, eventJSON)
	if err != nil {
		return fmt.Errorf("failed to put event: %v", err)
	}

	// 创建退款政策
	refundPolicy := RefundPolicy{
		ID:               fmt.Sprintf("refund_%s", order.ID),
		OrderID:          order.ID,
		RefundWindow:     90, // 90天退款窗口
		RefundPercentage: 100.0,
		Conditions:       "Product must be unopened and in original condition",
		IsActive:         true,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	refundPolicyJSON, err := json.Marshal(refundPolicy)
	if err != nil {
		return fmt.Errorf("failed to marshal refund policy: %v", err)
	}

	err = ctx.GetStub().PutState(refundPolicy.ID, refundPolicyJSON)
	if err != nil {
		return fmt.Errorf("failed to put refund policy: %v", err)
	}

	// 发出事件
	eventData := map[string]interface{}{
		"type":        "OrderCreated",
		"orderId":     order.ID,
		"orderNumber": order.OrderNumber,
		"userId":      order.UserID,
		"amount":      order.Amount,
		"timestamp":   order.CreatedAt,
	}

	eventDataJSON, _ := json.Marshal(eventData)
	ctx.GetStub().SetEvent("OrderCreated", eventDataJSON)

	return nil
}

// GetOrder 获取订单
func (oc *OrderContract) GetOrder(ctx contractapi.TransactionContextInterface, orderID string) (*Order, error) {
	orderJSON, err := ctx.GetStub().GetState(orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to read order: %v", err)
	}

	if orderJSON == nil {
		return nil, fmt.Errorf("order not found: %s", orderID)
	}

	var order Order
	err = json.Unmarshal(orderJSON, &order)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal order: %v", err)
	}

	return &order, nil
}

// UpdateOrderStatus 更新订单状态
func (oc *OrderContract) UpdateOrderStatus(ctx contractapi.TransactionContextInterface, orderID string, status string, description string) error {
	order, err := oc.GetOrder(ctx, orderID)
	if err != nil {
		return err
	}

	// 验证状态转换
	validTransitions := map[string][]string{
		"pending":   {"paid", "cancelled"},
		"paid":      {"shipped", "cancelled"},
		"shipped":   {"delivered", "cancelled"},
		"delivered": {"refunded"},
		"cancelled": {},
		"refunded":  {},
	}

	if !isValidTransition(order.Status, status, validTransitions) {
		return fmt.Errorf("invalid status transition from %s to %s", order.Status, status)
	}

	// 更新订单状态
	order.Status = status
	order.UpdatedAt = time.Now()

	orderJSON, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("failed to marshal order: %v", err)
	}

	err = ctx.GetStub().PutState(orderID, orderJSON)
	if err != nil {
		return fmt.Errorf("failed to put order: %v", err)
	}

	// 创建订单事件
	event := OrderEvent{
		ID:          fmt.Sprintf("event_%s_%d", orderID, time.Now().Unix()),
		OrderID:     orderID,
		EventType:   status,
		PayloadHash: fmt.Sprintf("hash_%s_%s", orderID, status),
		Timestamp:   time.Now(),
		Description: description,
		Metadata:    "",
	}

	eventJSON, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %v", err)
	}

	err = ctx.GetStub().PutState(event.ID, eventJSON)
	if err != nil {
		return fmt.Errorf("failed to put event: %v", err)
	}

	// 发出事件
	eventData := map[string]interface{}{
		"type":        "OrderStatusUpdated",
		"orderId":     orderID,
		"oldStatus":   order.Status,
		"newStatus":   status,
		"description": description,
		"timestamp":   time.Now(),
	}

	eventDataJSON, _ := json.Marshal(eventData)
	ctx.GetStub().SetEvent("OrderStatusUpdated", eventDataJSON)

	return nil
}

// ProcessRefund 处理退款
func (oc *OrderContract) ProcessRefund(ctx contractapi.TransactionContextInterface, orderID string, refundReason string) error {
	order, err := oc.GetOrder(ctx, orderID)
	if err != nil {
		return err
	}

	// 检查订单状态
	if order.Status != "delivered" {
		return fmt.Errorf("order must be delivered to process refund")
	}

	// 获取退款政策
	refundPolicyID := fmt.Sprintf("refund_%s", orderID)
	refundPolicyJSON, err := ctx.GetStub().GetState(refundPolicyID)
	if err != nil {
		return fmt.Errorf("failed to read refund policy: %v", err)
	}

	if refundPolicyJSON == nil {
		return fmt.Errorf("refund policy not found")
	}

	var refundPolicy RefundPolicy
	err = json.Unmarshal(refundPolicyJSON, &refundPolicy)
	if err != nil {
		return fmt.Errorf("failed to unmarshal refund policy: %v", err)
	}

	// 检查退款窗口
	refundDeadline := order.ActualDelivery.AddDate(0, 0, refundPolicy.RefundWindow)
	if time.Now().After(refundDeadline) {
		return fmt.Errorf("refund window has expired")
	}

	// 更新订单状态
	order.Status = "refunded"
	order.RefundReason = refundReason
	order.RefundHash = fmt.Sprintf("refund_hash_%s_%d", orderID, time.Now().Unix())
	order.UpdatedAt = time.Now()

	orderJSON, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("failed to marshal order: %v", err)
	}

	err = ctx.GetStub().PutState(orderID, orderJSON)
	if err != nil {
		return fmt.Errorf("failed to put order: %v", err)
	}

	// 创建退款事件
	event := OrderEvent{
		ID:          fmt.Sprintf("event_%s_%d", orderID, time.Now().Unix()),
		OrderID:     orderID,
		EventType:   "REFUND_PROCESSED",
		PayloadHash: order.RefundHash,
		Timestamp:   time.Now(),
		Description: fmt.Sprintf("Refund processed: %s", refundReason),
		Metadata:    "",
	}

	eventJSON, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %v", err)
	}

	err = ctx.GetStub().PutState(event.ID, eventJSON)
	if err != nil {
		return fmt.Errorf("failed to put event: %v", err)
	}

	// 发出退款事件
	eventData := map[string]interface{}{
		"type":         "RefundProcessed",
		"orderId":      orderID,
		"refundHash":   order.RefundHash,
		"refundReason": refundReason,
		"amount":       order.Amount * refundPolicy.RefundPercentage / 100,
		"timestamp":    time.Now(),
	}

	eventDataJSON, _ := json.Marshal(eventData)
	ctx.GetStub().SetEvent("RefundProcessed", eventDataJSON)

	return nil
}

// GetOrderEvents 获取订单事件
func (oc *OrderContract) GetOrderEvents(ctx contractapi.TransactionContextInterface, orderID string) ([]*OrderEvent, error) {
	iterator, err := ctx.GetStub().GetStateByPartialCompositeKey("event", []string{orderID})
	if err != nil {
		return nil, fmt.Errorf("failed to get order events: %v", err)
	}
	defer iterator.Close()

	var events []*OrderEvent
	for iterator.HasNext() {
		response, err := iterator.Next()
		if err != nil {
			return nil, fmt.Errorf("failed to get next event: %v", err)
		}

		var event OrderEvent
		err = json.Unmarshal(response.Value, &event)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal event: %v", err)
		}

		events = append(events, &event)
	}

	return events, nil
}

// GetOrdersByUser 根据用户ID获取订单
func (oc *OrderContract) GetOrdersByUser(ctx contractapi.TransactionContextInterface, userID string) ([]*Order, error) {
	iterator, err := ctx.GetStub().GetStateByPartialCompositeKey("order", []string{userID})
	if err != nil {
		return nil, fmt.Errorf("failed to get orders by user: %v", err)
	}
	defer iterator.Close()

	var orders []*Order
	for iterator.HasNext() {
		response, err := iterator.Next()
		if err != nil {
			return nil, fmt.Errorf("failed to get next order: %v", err)
		}

		var order Order
		err = json.Unmarshal(response.Value, &order)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal order: %v", err)
		}

		orders = append(orders, &order)
	}

	return orders, nil
}

// isValidTransition 检查状态转换是否有效
func isValidTransition(currentStatus, newStatus string, validTransitions map[string][]string) bool {
	allowedTransitions, exists := validTransitions[currentStatus]
	if !exists {
		return false
	}

	for _, allowed := range allowedTransitions {
		if allowed == newStatus {
			return true
		}
	}

	return false
}

func main() {
	orderContract, err := contractapi.NewChaincode(&OrderContract{})
	if err != nil {
		fmt.Printf("Error creating order contract: %v", err)
		return
	}

	if err := orderContract.Start(); err != nil {
		fmt.Printf("Error starting order contract: %v", err)
	}
}