package positions

import (
	"sync"
	"time"
)

type Controller struct {
	// Count is number of contract
	// Use for contract ratio
	Count             int
	TotalContractSize float64

	Positions sync.Map
}

type Position struct {
	ID  int
	OID string

	Symbol string
	Type   string

	Side  string
	Price float64
	Size  float64

	Commission float64

	CreatedAt time.Time
	UpdatedAt time.Time
}

func New() *Controller {
	return &Controller{
		Positions: sync.Map{},
	}
}

func (p *Controller) Size() (size float64) {
	p.Positions.Range(func(key, value any) bool {
		size += value.(Position).Size
		return true
	})
	return size
}

func (p *Controller) TotalCommission() (fee float64) {
	p.Positions.Range(func(key, value any) bool {
		fee += value.(Position).Commission
		return true
	})
	return fee
}

func (p *Controller) Set(k, v any) {
	p.Count++
	p.TotalContractSize += v.(Position).Size
	p.Positions.Store(k, v)
}

func (p *Controller) Load(k any) Position {
	v, isThere := p.Positions.Load(k)
	if !isThere {
		return Position{}
	}
	return v.(Position)
}

func (p *Controller) Update(k, updateV any) {
	v, isThere := p.Positions.LoadOrStore(k, updateV)
	if !isThere {
		p.Count++
		p.TotalContractSize += v.(Position).Size
		return
	}

	p.TotalContractSize -= v.(Position).Size
	p.TotalContractSize += updateV.(Position).Size

	// Replace
	p.Positions.Store(k, updateV)
}

func (p *Controller) Delete(k any) {
	p.Positions.Delete(k)
}
