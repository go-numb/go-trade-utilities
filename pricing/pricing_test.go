package pricing

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestIsVolatirity(t *testing.T) {
	// 10ms配信を仮定し、100秒の値動きを作る
	count := 10000
	pa := New(100, count, 0.01)
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)

	for i := 0; i < count; i++ {
		// ランダムウォークな価格変化を起こす乱数
		n := r.Int()

		// 疑似板
		// - midから0.1pipsより狭く上下に板
		// - 0.2pips(銭)以下のスプレッド
		if i == 0 {
			// 基準となる初期価格を作る
			pa.Prices = append(pa.Prices, Price{
				Ltp:       125.05,
				Mid:       125.05,
				Ask:       125.05 + r.Float64()*0.001,
				Bid:       125.05 - r.Float64()*0.001,
				Timestamp: time.Now().UnixMilli(),
			})
		} else {
			// 過去の価格を踏襲し疑似的にプライスアクションを作る
			priceaction := r.Float64() * 0.001
			ltp := pa.Prices[i-1].Ltp
			if n%3 == 0 {
				ltp += priceaction
			} else if n%5 == 0 {
				ltp -= priceaction
			}

			pa.Prices = append(pa.Prices, Price{
				Ltp:       ltp,
				Mid:       ltp,
				Ask:       ltp + r.Float64()*0.001,
				Bid:       ltp - r.Float64()*0.001,
				Timestamp: time.Now().UnixMilli(),
			})
		}

		// 擬似的に積極的な値動きを加える
		// 300 の倍数で1pipsのPump ≒ 3秒に一度
		// 500 の倍数で1pipsのDump ≒ 5秒に一度
		// 1500 の倍数で2pipsのPump&Dump ≒ 15秒に一度

		if n%1500 == 0 {
			pa.Prices[i].Ltp += 0.02
			pa.Prices[i].Ask += 0.02
		} else if n%500 == 0 {
			pa.Prices[i].Ltp -= 0.02
			pa.Prices[i].Bid -= 0.02
		} else if n%300 == 0 {
			pa.Prices[i].Ask += 0.02
			pa.Prices[i].Ltp += 0.02
		}

		// 配信頻度
		// 値の到着など通信の時間を見込む
		time.Sleep(10 * time.Millisecond)

		if i < 1 {
			continue
		}

		// 逐次値が変わる状態を検知
		// 判定を挟むとする
		fmt.Printf("%d: get price: %.4f to %.4f, ", i, pa.Prices[i-1].Ltp, pa.Prices[i].Ltp)
		isPump, isVolatility := pa.Volatility()
		if isVolatility {
			if isPump == 1 {
				fmt.Printf("FLAG is %v, PUMP↑ DO SOMETHING!!", isVolatility)
			} else if isPump == -1 {
				fmt.Printf("FLAG is %v, DUMP↓ DO SOMETHING!!", isVolatility)
			}
		}

		fmt.Println("THROUGH")
	}

}
