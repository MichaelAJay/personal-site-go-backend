package triangle

import (
	"bytes"
	"image/color"
	"math"

	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	_ "gonum.org/v1/plot/vg/vgsvg"
)

const (
	s = 68
)

type Triangle struct {
	V1, V2, V3 vg.Point
}

func GenerateSierpinskiSVG(iterations int, color color.Color) (string, error) {
	h := s * math.Sqrt(3) / 2

	// Create canvas
	width, height := vg.Length(s), vg.Length(h)
	canvasWriterTo, err := draw.NewFormattedCanvas(width, height, "svg")
	if err != nil {
		panic(err)
	}

	// Create a draw.Canvas from vg.CanvasWriterTo
	canvas := draw.NewCanvas(canvasWriterTo, width, height)

	// Define vertices
	v1 := vg.Point{X: 0, Y: 0}
	v2 := vg.Point{X: vg.Length(s), Y: 0}
	v3 := vg.Point{X: vg.Length(s / 2), Y: vg.Length(h)}
	borderTriangle := Triangle{V1: v1, V2: v2, V3: v3}

	// Draw the border
	drawTriangle(&canvas, borderTriangle, color)

	// Iterate
	source := []Triangle{borderTriangle}
	for i := 1; i <= iterations; i++ {
		source = iteration(&canvas, source, color)
	}

	var svgBuffer bytes.Buffer

	_, writeToErr := canvasWriterTo.WriteTo(&svgBuffer)
	if writeToErr != nil {
		return "", writeToErr
	}

	return svgBuffer.String(), nil
}

func drawTriangle(c *draw.Canvas, t Triangle, col color.Color) {
	lines := []vg.Point{t.V1, t.V2, t.V3, t.V1}
	c.StrokeLines(draw.LineStyle{Color: col, Width: vg.Length(0.2)}, lines)
}

func iteration(c *draw.Canvas, source []Triangle, col color.Color) []Triangle {
	var output []Triangle

	for _, triangle := range source {
		nextSourceTriangles := generateSubTriangles(c, triangle, col)
		output = append(output, nextSourceTriangles...)
	}

	return output
}

func generateSubTriangles(c *draw.Canvas, t Triangle, col color.Color) []Triangle {
	midpoint_v1_v2 := calculateMidpoint(t.V1, t.V2)
	midpoint_v2_v3 := calculateMidpoint(t.V2, t.V3)
	midpoint_v3_v1 := calculateMidpoint(t.V3, t.V1)

	vertices := []vg.Point{midpoint_v1_v2, midpoint_v2_v3, midpoint_v3_v1}
	c.FillPolygon(col, vertices)

	t1 := Triangle{V1: t.V1, V2: midpoint_v1_v2, V3: midpoint_v3_v1}
	t2 := Triangle{V1: t.V2, V2: midpoint_v1_v2, V3: midpoint_v2_v3}
	t3 := Triangle{V1: t.V3, V2: midpoint_v2_v3, V3: midpoint_v3_v1}

	return []Triangle{t1, t2, t3}
}

func calculateMidpoint(p1, p2 vg.Point) vg.Point {
	return vg.Point{
		X: (p1.X + p2.X) / 2,
		Y: (p1.Y + p2.Y) / 2,
	}
}
