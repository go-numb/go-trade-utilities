package plot

import (
	"fmt"
	"image/color"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

type Scatter struct {
	Title string
	// Value
	Length         int
	xLabel, yLabel string
	xValue, yValue []float64
}

func NewScatter(title, xLabel, yLabel string, xValues, yValues []float64) *Scatter {
	if len(xValues) != len(yValues) {
		return &Scatter{}
	}

	return &Scatter{
		Title:  title,
		Length: len(xValues),
		xLabel: xLabel,
		yLabel: yLabel,
		xValue: xValues,
		yValue: yValues,
	}
}

func (p *Scatter) Save(output string) error {
	if p.Length < 1 {
		return fmt.Errorf("struct has not length")
	}

	plt := plot.New()

	plt.Y.Color = color.RGBA{R: 189, G: 189, B: 189, A: 255}
	plt.Y.Tick.Color = color.RGBA{R: 189, G: 189, B: 189, A: 255}
	plt.X.Color = color.RGBA{R: 189, G: 189, B: 189, A: 255}
	plt.X.Tick.Color = color.RGBA{R: 189, G: 189, B: 189, A: 255}

	plt.Title.Text = p.Title
	plt.Y.Label.Text = "data_y"
	if p.yLabel != "" {
		plt.Y.Label.Text = p.yLabel
	}
	plt.X.Label.Text = "data_x"
	if p.xLabel != "" {
		plt.X.Label.Text = p.xLabel
	}

	points := make(plotter.XYs, p.Length)
	for i := range points {
		points[i].X = p.xValue[i]
		points[i].Y = p.yValue[i]
	}

	// if err = plotutil.AddScatters(plt, points); err != nil {
	// 	return err
	// }

	l, err := plotter.NewScatter(points)
	if err != nil {
		return err
	}
	// #82b1ff blue accent-1
	l.Color = color.RGBA{R: 130, G: 177, B: 255, A: 255}
	plt.Add(l)

	if err := plt.Save(4*vg.Inch, 3*vg.Inch, output); err != nil {
		return err
	}

	if err := plt.Save(4*vg.Inch, 4*vg.Inch, output); err != nil {
		return err
	}

	return nil
}
