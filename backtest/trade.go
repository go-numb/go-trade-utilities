package backtest

import "strings"

type Trade struct {
	Unique     string
	Side       string
	EntryPrice float64
	ExitPrice  float64

	IsWin bool
	pnl   float64
}

func (t *Trade) isBUY() bool {
	return strings.ToUpper(t.Side) == BUY
}

func (t *Trade) pnL() float64 {
	if t.isBUY() {
		t.pnl = t.ExitPrice - t.EntryPrice
		return t.pnl
	}
	t.pnl = t.EntryPrice - t.ExitPrice
	return t.pnl
}

func (t *Trade) win() bool {
	if t.pnl <= 0 { // 0 is Loose
		return t.IsWin
	}
	t.IsWin = true
	return t.IsWin
}
