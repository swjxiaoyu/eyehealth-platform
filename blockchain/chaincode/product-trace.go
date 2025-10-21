package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// ProductTraceContract 产品溯源智能合约
type ProductTraceContract struct {
	contractapi.Contract
}

// TraceRecord 溯源记录结构
type TraceRecord struct {
	ID            string    `json:"id"`
	ProductID     string    `json:"productId"`
	Stage         string    `json:"stage"`
	DocumentHash  string    `json:"documentHash"`
	Issuer        string    `json:"issuer"`
	IssuerName    string    `json:"issuerName"`
	Location      string    `json:"location"`
	Coordinates   string    `json:"coordinates"`
	Timestamp     time.Time `json:"timestamp"`
	Temperature   float64   `json:"temperature"`
	Humidity      float64   `json:"humidity"`
	Metadata      string    `json:"metadata"`
	CertificateURL string   `json:"certificateUrl"`
	CertificateHash string  `json:"certificateHash"`
	IsVerified    bool      `json:"isVerified"`
	VerificationMethod string `json:"verificationMethod"`
}

// Product 产品结构
type Product struct {
	ID              string    `json:"id"`
	SKU             string    `json:"sku"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	ManufacturerID  string    `json:"manufacturerId"`
	ManufacturerName string   `json:"manufacturerName"`
	QRCode          string    `json:"qrCode"`
	Barcode         string    `json:"barcode"`
	Price           float64   `json:"price"`
	Currency        string    `json:"currency"`
	Category        string    `json:"category"`
	Subcategory     string    `json:"subcategory"`
	Specifications  string    `json:"specifications"`
	Ingredients     string    `json:"ingredients"`
	ExpiryDate      time.Time `json:"expiryDate"`
	BatchNumber     string    `json:"batchNumber"`
	SerialNumber    string    `json:"serialNumber"`
	IsActive        bool      `json:"isActive"`
	IsVerified      bool      `json:"isVerified"`
	Metadata        string    `json:"metadata"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

// CreateTrace 创建溯源记录
func (pc *ProductTraceContract) CreateTrace(ctx contractapi.TransactionContextInterface, traceData string) error {
	var trace TraceRecord
	err := json.Unmarshal([]byte(traceData), &trace)
	if err != nil {
		return fmt.Errorf("failed to unmarshal trace data: %v", err)
	}

	// 验证必需字段
	if trace.ProductID == "" || trace.Stage == "" || trace.DocumentHash == "" || trace.Issuer == "" {
		return fmt.Errorf("missing required fields")
	}

	// 设置时间戳
	trace.Timestamp = time.Now()

	// 生成唯一ID
	trace.ID = fmt.Sprintf("trace_%s_%d", trace.ProductID, trace.Timestamp.Unix())

	// 存储溯源记录
	traceJSON, err := json.Marshal(trace)
	if err != nil {
		return fmt.Errorf("failed to marshal trace: %v", err)
	}

	err = ctx.GetStub().PutState(trace.ID, traceJSON)
	if err != nil {
		return fmt.Errorf("failed to put trace: %v", err)
	}

	// 创建复合键用于查询
	compositeKey, err := ctx.GetStub().CreateCompositeKey("trace", []string{trace.ProductID, trace.Stage, trace.ID})
	if err != nil {
		return fmt.Errorf("failed to create composite key: %v", err)
	}

	err = ctx.GetStub().PutState(compositeKey, traceJSON)
	if err != nil {
		return fmt.Errorf("failed to put composite key: %v", err)
	}

	// 发出事件
	event := map[string]interface{}{
		"type":      "TraceCreated",
		"traceId":   trace.ID,
		"productId": trace.ProductID,
		"stage":     trace.Stage,
		"issuer":    trace.Issuer,
		"timestamp": trace.Timestamp,
	}

	eventJSON, _ := json.Marshal(event)
	ctx.GetStub().SetEvent("TraceCreated", eventJSON)

	return nil
}

// GetTrace 获取溯源记录
func (pc *ProductTraceContract) GetTrace(ctx contractapi.TransactionContextInterface, traceID string) (*TraceRecord, error) {
	traceJSON, err := ctx.GetStub().GetState(traceID)
	if err != nil {
		return nil, fmt.Errorf("failed to read trace: %v", err)
	}

	if traceJSON == nil {
		return nil, fmt.Errorf("trace not found: %s", traceID)
	}

	var trace TraceRecord
	err = json.Unmarshal(traceJSON, &trace)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal trace: %v", err)
	}

	return &trace, nil
}

// GetTracesByProduct 根据产品ID获取所有溯源记录
func (pc *ProductTraceContract) GetTracesByProduct(ctx contractapi.TransactionContextInterface, productID string) ([]*TraceRecord, error) {
	iterator, err := ctx.GetStub().GetStateByPartialCompositeKey("trace", []string{productID})
	if err != nil {
		return nil, fmt.Errorf("failed to get traces by product: %v", err)
	}
	defer iterator.Close()

	var traces []*TraceRecord
	for iterator.HasNext() {
		response, err := iterator.Next()
		if err != nil {
			return nil, fmt.Errorf("failed to get next trace: %v", err)
		}

		var trace TraceRecord
		err = json.Unmarshal(response.Value, &trace)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal trace: %v", err)
		}

		traces = append(traces, &trace)
	}

	return traces, nil
}

// VerifyTrace 验证溯源记录
func (pc *ProductTraceContract) VerifyTrace(ctx contractapi.TransactionContextInterface, traceID string, verificationMethod string) error {
	trace, err := pc.GetTrace(ctx, traceID)
	if err != nil {
		return err
	}

	// 更新验证状态
	trace.IsVerified = true
	trace.VerificationMethod = verificationMethod

	traceJSON, err := json.Marshal(trace)
	if err != nil {
		return fmt.Errorf("failed to marshal trace: %v", err)
	}

	err = ctx.GetStub().PutState(traceID, traceJSON)
	if err != nil {
		return fmt.Errorf("failed to put trace: %v", err)
	}

	// 发出验证事件
	event := map[string]interface{}{
		"type":               "TraceVerified",
		"traceId":            traceID,
		"productId":          trace.ProductID,
		"verificationMethod": verificationMethod,
		"timestamp":          time.Now(),
	}

	eventJSON, _ := json.Marshal(event)
	ctx.GetStub().SetEvent("TraceVerified", eventJSON)

	return nil
}

// CreateProduct 创建产品
func (pc *ProductTraceContract) CreateProduct(ctx contractapi.TransactionContextInterface, productData string) error {
	var product Product
	err := json.Unmarshal([]byte(productData), &product)
	if err != nil {
		return fmt.Errorf("failed to unmarshal product data: %v", err)
	}

	// 验证必需字段
	if product.SKU == "" || product.Name == "" || product.ManufacturerID == "" {
		return fmt.Errorf("missing required fields")
	}

	// 设置时间戳
	product.CreatedAt = time.Now()
	product.UpdatedAt = time.Now()
	product.IsActive = true
	product.IsVerified = false

	// 生成唯一ID
	product.ID = fmt.Sprintf("product_%s_%d", product.SKU, product.CreatedAt.Unix())

	// 存储产品
	productJSON, err := json.Marshal(product)
	if err != nil {
		return fmt.Errorf("failed to marshal product: %v", err)
	}

	err = ctx.GetStub().PutState(product.ID, productJSON)
	if err != nil {
		return fmt.Errorf("failed to put product: %v", err)
	}

	// 创建复合键用于查询
	compositeKey, err := ctx.GetStub().CreateCompositeKey("product", []string{product.SKU, product.ID})
	if err != nil {
		return fmt.Errorf("failed to create composite key: %v", err)
	}

	err = ctx.GetStub().PutState(compositeKey, productJSON)
	if err != nil {
		return fmt.Errorf("failed to put composite key: %v", err)
	}

	// 发出事件
	event := map[string]interface{}{
		"type":      "ProductCreated",
		"productId": product.ID,
		"sku":       product.SKU,
		"name":      product.Name,
		"timestamp": product.CreatedAt,
	}

	eventJSON, _ := json.Marshal(event)
	ctx.GetStub().SetEvent("ProductCreated", eventJSON)

	return nil
}

// GetProduct 获取产品信息
func (pc *ProductTraceContract) GetProduct(ctx contractapi.TransactionContextInterface, productID string) (*Product, error) {
	productJSON, err := ctx.GetStub().GetState(productID)
	if err != nil {
		return nil, fmt.Errorf("failed to read product: %v", err)
	}

	if productJSON == nil {
		return nil, fmt.Errorf("product not found: %s", productID)
	}

	var product Product
	err = json.Unmarshal(productJSON, &product)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal product: %v", err)
	}

	return &product, nil
}

// GetProductBySKU 根据SKU获取产品
func (pc *ProductTraceContract) GetProductBySKU(ctx contractapi.TransactionContextInterface, sku string) (*Product, error) {
	iterator, err := ctx.GetStub().GetStateByPartialCompositeKey("product", []string{sku})
	if err != nil {
		return nil, fmt.Errorf("failed to get product by SKU: %v", err)
	}
	defer iterator.Close()

	if iterator.HasNext() {
		response, err := iterator.Next()
		if err != nil {
			return nil, fmt.Errorf("failed to get next product: %v", err)
		}

		var product Product
		err = json.Unmarshal(response.Value, &product)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal product: %v", err)
		}

		return &product, nil
	}

	return nil, fmt.Errorf("product not found with SKU: %s", sku)
}

// VerifyProduct 验证产品
func (pc *ProductTraceContract) VerifyProduct(ctx contractapi.TransactionContextInterface, productID string) error {
	product, err := pc.GetProduct(ctx, productID)
	if err != nil {
		return err
	}

	// 更新验证状态
	product.IsVerified = true
	product.UpdatedAt = time.Now()

	productJSON, err := json.Marshal(product)
	if err != nil {
		return fmt.Errorf("failed to marshal product: %v", err)
	}

	err = ctx.GetStub().PutState(productID, productJSON)
	if err != nil {
		return fmt.Errorf("failed to put product: %v", err)
	}

	// 发出验证事件
	event := map[string]interface{}{
		"type":      "ProductVerified",
		"productId": productID,
		"sku":       product.SKU,
		"timestamp": time.Now(),
	}

	eventJSON, _ := json.Marshal(event)
	ctx.GetStub().SetEvent("ProductVerified", eventJSON)

	return nil
}

func main() {
	productTraceContract, err := contractapi.NewChaincode(&ProductTraceContract{})
	if err != nil {
		fmt.Printf("Error creating product trace contract: %v", err)
		return
	}

	if err := productTraceContract.Start(); err != nil {
		fmt.Printf("Error starting product trace contract: %v", err)
	}
}