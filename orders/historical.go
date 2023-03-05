package orders

import (
	"strings"
	"sync"

	"github.com/go-numb/go-trade-utilities/types"
)

type HistoricalController struct {
	// Count is number of orders
	// Use for contract ratio
	Count          int
	TotalOrderSize float64

	Orders sync.Map
}

func NewHistorical() *HistoricalController {
	return &HistoricalController{
		Orders: sync.Map{},
	}
}

func (p *HistoricalController) Length() (buy, sell int) {
	p.Orders.Range(func(key, value any) bool {
		v := value.(Order)
		if strings.ToUpper(v.Side) == types.BUY {
			buy++
		} else if strings.ToUpper(v.Side) == types.SELL {
			sell++
		}
		return true
	})
	return buy, sell
}

func (p *HistoricalController) Size() (buy, sell float64) {
	p.Orders.Range(func(key, value any) bool {
		v := value.(Order)
		if strings.ToUpper(v.Side) == types.BUY {
			buy += v.Size
		} else if strings.ToUpper(v.Side) == types.SELL {
			sell += v.Size
		}
		return true
	})
	return buy, sell
}

func (p *HistoricalController) TotalCommission() (buy, sell float64) {
	p.Orders.Range(func(key, value any) bool {
		v := value.(Order)
		if strings.ToUpper(v.Side) == types.BUY {
			buy += v.Commission
		} else if strings.ToUpper(v.Side) == types.SELL {
			sell += v.Commission
		}
		return true
	})
	return buy, sell
}

func (p *HistoricalController) Set(k, v any) {
	p.Count++
	p.Orders.Store(k, v)
}

func (p *HistoricalController) Load(k any) Order {
	v, isThere := p.Orders.Load(k)
	if !isThere {
		return Order{}
	}
	return v.(Order)
}

func (p *HistoricalController) Update(k, updateV any) {
	v, isThere := p.Orders.LoadOrStore(k, updateV)
	if !isThere {
		p.Count++
		return
	}

	p.TotalOrderSize -= v.(Order).Size
	p.TotalOrderSize += updateV.(Order).Size

	// Replace
	p.Orders.Store(k, updateV)
}

func (p *HistoricalController) Delete(k any) {
	p.Orders.Delete(k)
}
