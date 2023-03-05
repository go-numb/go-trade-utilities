package orders

import (
	"sync"
	"time"
)

type Controller struct {
	// Count is number of orders
	// Use for contract ratio
	Count          int
	TotalOrderSize float64

	Orders sync.Map
}

type Order struct {
	ID  int
	OID string

	Symbol string
	Type   string

	Side  string
	Price float64
	Size  float64

	IsContracted bool
	Commission   float64

	CreatedAt time.Time
	UpdatedAt time.Time
}

func New() *Controller {
	return &Controller{
		Orders: sync.Map{},
	}
}

func (p *Controller) Length() (lenght int) {
	p.Orders.Range(func(key, value any) bool {
		lenght++
		return true
	})
	return lenght
}

func (p *Controller) Set(k, v any) {
	p.Count++
	p.TotalOrderSize += v.(Order).Size
	p.Orders.Store(k, v)
}

func (p *Controller) Load(k any) Order {
	v, isThere := p.Orders.Load(k)
	if !isThere {
		return Order{}
	}
	return v.(Order)
}

func (p *Controller) Update(k, updateV any) {
	v, isThere := p.Orders.LoadOrStore(k, updateV)
	if !isThere {
		p.Count++
		p.TotalOrderSize += v.(Order).Size
		return
	}

	p.TotalOrderSize -= v.(Order).Size
	p.TotalOrderSize += updateV.(Order).Size

	// Replace
	p.Orders.Store(k, updateV)
}

func (p *Controller) Delete(k any) {
	p.Orders.Delete(k)
}
