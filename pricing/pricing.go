package pricing

// 仕様:
// 様々な取引所データに対応する
// ボラティリティの検知
// 発火条件を指定

type Controller struct {
	// ターゲットとする過去データをミリ秒指定
	TargetMillisecBefore int
	TargetActionPrice    float64

	// Flag
	flag int

	MaxLength int
	Prices    []Price
}

func New(
	targetMillisecBefore,
	maxPricesLength int,
	targetActionPrice float64,
) *Controller {
	return &Controller{
		TargetMillisecBefore: targetMillisecBefore,
		TargetActionPrice:    targetActionPrice,

		flag: 0,

		MaxLength: maxPricesLength,
		Prices:    []Price{},
	}
}

type Price struct {
	// 共通項
	Timestamp int64

	// 約定データのみで計算
	Ltp float64

	// BestAsk/Bidで計算
	Mid    float64
	Ask    float64
	Bid    float64
	Volume float64
}

func (p *Controller) Volatility() (isPump int, isVolatility bool) {
	// 保持データが少なければFalse
	l := len(p.Prices)
	if len(p.Prices) < p.TargetMillisecBefore {
		return p._changeFlag(0)
	}

	// 最新価格と基準Pump&Dump
	ltp := p.Prices[l-1].Ltp
	// 過去の時間
	preview := p.Prices[l-1].Timestamp - int64(p.TargetMillisecBefore)

	for i := 1; i < l; i++ {
		// 振り返る時間が指定ミリ秒より過去の場合はFalse
		if p.Prices[l-1-i].Timestamp < preview {
			return p._changeFlag(0)
		}

		if ltp > p.Prices[l-1-i].Ltp+p.TargetActionPrice { // Pump
			return p._changeFlag(1)
		} else if ltp < p.Prices[l-1-i].Ltp-p.TargetActionPrice { // Dump
			return p._changeFlag(-1)
		}
	}

	return p._changeFlag(0)
}

func (p *Controller) Set(price Price) {
	p.Prices = append(p.Prices, price)
	if len(p.Prices) > p.MaxLength {
		p.Prices = p.Prices[len(p.Prices)-p.MaxLength:]
	}
}

func (p *Controller) IsFlag() bool {
	return p.flag != 0
}

func (p *Controller) Flag() int {
	return p.flag
}

func (p *Controller) _changeFlag(newFlag int) (int, bool) {
	if p.flag != newFlag {
		p.flag = newFlag
		if p.IsFlag() {
			return p.flag, true
		}
	}
	return p.flag, false
}
