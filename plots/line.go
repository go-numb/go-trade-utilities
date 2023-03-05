package plot

import (
	"fmt"
	"image/color"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

type Outputer interface {
	Save(output string) error
}

type Line struct {
	Title string
	// Value
	Length         int
	xLabel, yLabel string
	Value          []float64
}

func NewLine(title, xLabel, yLabel string, values []float64) *Line {
	return &Line{
		Title:  title,
		Length: len(values),
		xLabel: xLabel,
		yLabel: yLabel,
		Value:  values,
	}
}

func (p *Line) Save(output string) error {
	if p.Length < 1 {
		return fmt.Errorf("struct has not length")
	}

	plt := plot.New()

	plt.Y.Color = color.RGBA{R: 189, G: 189, B: 189, A: 255}
	plt.Y.Tick.Color = color.RGBA{R: 189, G: 189, B: 189, A: 255}
	plt.X.Color = color.RGBA{R: 189, G: 189, B: 189, A: 255}
	plt.X.Tick.Color = color.RGBA{R: 189, G: 189, B: 189, A: 255}

	plt.Title.Text = p.Title
	plt.Y.Label.Text = "data"
	if p.yLabel != "" {
		plt.Y.Label.Text = p.yLabel
	}
	plt.X.Label.Text = "time"
	if p.xLabel != "" {
		plt.X.Label.Text = p.xLabel
	}

	points := make(plotter.XYs, p.Length)
	for i := range points {
		points[i].X = float64(i)
		points[i].Y = p.Value[i]
	}

	// Color自動設定: R:255
	// if err = plotutil.AddLines(plt, points); err != nil {
	// 	return err
	// }

	l, err := plotter.NewLine(points)
	if err != nil {
		return err
	}
	// #82b1ff blue accent-1
	l.Color = color.RGBA{R: 130, G: 177, B: 255, A: 255}
	plt.Add(l)

	if err := plt.Save(4*vg.Inch, 3*vg.Inch, output); err != nil {
		return err
	}

	return nil
}
