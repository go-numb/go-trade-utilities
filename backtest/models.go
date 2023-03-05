package backtest

import (
	"math"

	"gonum.org/v1/gonum/stat"
)

const (
	BUY  = "BUY"
	SELL = "SELL"

	RANDOMETHRESHOLD = 1.64
)

type Controller struct {
	Trades []Trade

	Expected float64
	PnLs     []float64
}

func New() *Controller {
	return &Controller{
		Trades: make([]Trade, 0),
	}
}

/*
Expect
ここまで下記が満たされる場合
期待値 = 純損失/取引回数

以下、
期待値 = (勝率*平均利益)/(敗率*平均損失)
*/
func (t *Controller) Expect() (pnl float64, pnls []float64) {
	pnl, pnls = t.PNL()
	t.Expected = pnl / float64(t.Count())
	t.PnLs = pnls
	return t.Expected, t.PnLs
}

// SharpRatio リスク・リターン指標
func (t *Controller) SharpRatio(rfr float64) float64 {
	pnls := t.PnLs
	if pnls == nil {
		_, pnls = t.PNL()
	}

	mean, std := stat.MeanStdDev(pnls, nil)

	return (mean - rfr) / std
}

/*
ProfitFactor
リターンの大きさを表す指標
*/
func (t *Controller) ProfitFactor() float64 {
	var (
		p float64
		l float64
	)
	for i := 0; i < len(t.Trades); i++ {
		if !t.Trades[i].win() {
			l += t.Trades[i].pnL()
			continue
		}
		p += t.Trades[i].pnL()
	}

	return float64(p) / -float64(l)
}

/*
PRR
取引数が多いほど高く評価するPF
*/
func (t *Controller) PRR() float64 {
	var (
		profitN, lossN = 0, 0
		p              float64
		l              float64
	)
	for i := 0; i < len(t.Trades); i++ {
		if !t.Trades[i].win() {
			lossN++
			l += t.Trades[i].pnL()
			continue
		}
		profitN++
		p += t.Trades[i].pnL()
	}

	a := (float64(profitN) - math.Sqrt(float64(profitN))) * p
	b := (float64(lossN) + math.Sqrt(float64(lossN))) * l

	return a / -b
}

// RandomPF ランダムトレードでの収束ProfitFactor
func (t *Controller) RandomPF() float64 {
	var n = float64(t.Count())
	a := n + RANDOMETHRESHOLD*math.Sqrt(n)
	b := n - RANDOMETHRESHOLD*math.Sqrt(n)

	return a / b
}

// PNL 純損益
func (t *Controller) PNL() (pnl float64, pnls []float64) {
	pnls = make([]float64, t.Count())
	for i := 0; i < len(t.Trades); i++ {
		pnl += t.Trades[i].pnL()
		pnls[i] = t.Trades[i].pnl
		t.Trades[i].win()
	}

	return pnl, pnls
}

// Profit
func (t *Controller) Profit() float64 {
	var f float64
	for i := 0; i < len(t.Trades); i++ {
		if !t.Trades[i].win() {
			continue
		}
		f += t.Trades[i].pnL()
	}

	return f
}

// Loss
func (t *Controller) Loss() float64 {
	var f float64
	for i := 0; i < len(t.Trades); i++ {
		if t.Trades[i].win() {
			continue
		}
		f += t.Trades[i].pnL()
	}

	return f
}

// Count 取引回数
func (t *Controller) Count() int {
	return len(t.Trades)
}

// WinR 勝率
func (t *Controller) WinR() float64 {
	var wins int
	for i := 0; i < t.Count(); i++ {
		if t.Trades[i].win() {
			wins++
		}
	}

	return float64(wins) / float64(t.Count())
}
