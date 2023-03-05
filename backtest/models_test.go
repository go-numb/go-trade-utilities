package backtest

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

func TestExpect(t *testing.T) {
	start := time.Now()
	defer fmt.Printf("%.3fs\n", time.Since(start).Seconds())

	side := BUY
	basef := 10000.0
	add := basef * 0.1

	s := rand.NewSource(time.Now().UnixMilli())
	r := rand.New(s)

	for i := 0; i < 10; i++ {
		data := New()
		// 1000Tradeを擬似的に作る
		count := 1000
		for i := 0; i < count; i++ {
			a := basef + r.Float64()*100
			if i%2 == 0 {
				side = SELL
				if i%19 == 0 {
					a += add
				}
			} else {
				side = BUY
				a -= add
			}
			data.Trades = append(data.Trades, Trade{
				Unique:     fmt.Sprintf("%d", i),
				Side:       side,
				EntryPrice: a,
				ExitPrice:  basef + r.Float64()*100,
			})
		}

		f, pnls := data.Expect()
		fmt.Printf("expected: %f\n", f)
		fmt.Printf("勝率: %.3f%%\n", data.WinR()*100)
		fmt.Printf("%.1f, %.1f\n", data.Profit(), data.Loss())
		fmt.Printf("pf: %.2f, PRR: %.2f, random: %.2f\n", data.ProfitFactor(), data.PRR(), data.RandomPF())
		fmt.Printf("Sharp ratio: %f\n", data.SharpRatio(0))
		fmt.Println("")

		plt := plot.New()

		plt.Title.Text = "pnl addition"
		plt.Y.Label.Text = "PnL"
		plt.X.Label.Text = "count"

		points := make(plotter.XYs, len(pnls))
		var add float64
		for i := range points {
			add += pnls[i]
			points[i].X = float64(i)
			points[i].Y = add
		}

		l, _ := plotter.NewLine(points)
		plt.Add(l)

		plt.Save(4*vg.Inch, 3*vg.Inch, "./line.png")

		time.Sleep(3 * time.Second)
	}
}
