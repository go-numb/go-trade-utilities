package plot

import (
	"strings"

	"github.com/go-gota/gota/dataframe"
	"github.com/mattn/go-pairplot"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/vg"
)

type Seaborn struct {
	Data dataframe.DataFrame
}

func NewSeaborn(s interface{}) *Seaborn {
	var df dataframe.DataFrame
	switch v := s.(type) {
	case []string:
		df = dataframe.ReadCSV(strings.NewReader(strings.Join(v, "\n")))
	default:
		df = dataframe.LoadStructs(s)
	}

	return &Seaborn{
		Data: df,
	}
}

// SubPlot plot like a seaborn, too much consuming
func (p *Seaborn) SubPlot(w, h int, filename string) error {
	plt := plot.New()
	pp, err := pairplot.NewPairPlotDataFrame(p.Data)
	if err != nil {
		return err
	}
	// pp.SetHue("Name")
	plt.HideAxes()
	plt.Add(pp)
	plt.Save(vg.Length(w)*vg.Inch, vg.Length(h)*vg.Inch, filename)

	return nil
}
