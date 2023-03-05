package controllers

import (
	"github.com/go-numb/go-trade-utilities/orders"
	"github.com/go-numb/go-trade-utilities/positions"
	"github.com/go-numb/go-trade-utilities/pricing"
)

type CentralController struct {
	// PriceController 指定長の価格情報を保存
	PriceController *pricing.Controller
	// OrderController 注文待機建玉管理
	OrderController *orders.Controller
	// PositionController 約定済み建玉管理
	PositionController *positions.Controller

	// HistoricalController 注文及び約定履歴を保存
	HistoricalController *orders.HistoricalController
}

func New(
	targetMillisecBefore int,
	maxPricesLength int,
	targetActionPrice float64,
) *CentralController {
	return &CentralController{
		PriceController:    pricing.New(targetMillisecBefore, maxPricesLength, targetActionPrice),
		OrderController:    orders.New(),
		PositionController: positions.New(),
	}
}

func (p *CentralController) F() {
	p.OrderController.Orders.Store("key", "value")
}
