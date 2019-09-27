package main

import (
	"fmt"
	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/llgcode/draw2d/draw2dkit"
	"github.com/markfarnan/go-canvas/canvas"
	"image/color"
	"syscall/js"
)

var cells [][]Cell

var done chan struct{}

var cvs *canvas.Canvas2d
var width = 40
var height = 20
var cellSize = 20.0


func main() {

	cvs, _ = canvas.NewCanvas2d(true)
	js.Global().Set("click", js.FuncOf(handleClick))
	fmt.Println(js.
		Global().
		Get("document").
		Get("body").
		Get("children").
		Index(0).
		Call("addEventListener", "click", js.Global().Get("click")))

	for i := 0; i < width; i++ {
		var column []Cell
		for j := 0; j < height; j++ {
			var t CellType = Liquid
			var v = 0.0
			if j == height - 1 {
				t = Solid
			} else if false {
				t = Liquid
				v = 1.0
			}
			column = append(column, NewCell(i, j, t, v))
		}
		cells = append(cells, column)
	}
	cells[3][height-2].cellType = Solid
	cells[3][height-3].cellType = Solid
	cells[7][height-2].cellType = Solid
	cells[7][height-3].cellType = Solid
	cells[3][height-4].cellType = Solid
	cells[3][height-5].cellType = Solid
	cells[7][height-4].cellType = Solid
	cells[7][height-5].cellType = Solid

	cells[5][height-2].cellType = Solid
	cells[5][height-3].cellType = Solid
	cells[5][height-4].cellType = Solid
	cells[5][height-5].cellType = Solid

	cells[8][height-3].cellType = Solid
	cells[9][height-3].cellType = Solid
	cvs.Start(60, Render)
	<-done
}

func handleClick(this js.Value, args []js.Value) interface{} {
	mouseEvent := args[0]
	btn := mouseEvent.Get("button")
	x := mouseEvent.Get("clientX").Int()/int(cellSize)
	y := mouseEvent.Get("clientY").Int()/int(cellSize)
	if cells[x][y].cellType == Void || cells[x][y].cellType == Liquid {
		cells[x][y].cellType = Solid
	} else {
		cells[x][y].cellType = Liquid
		cells[x][y].volume = maxVolume + maxPressure
	}
	fmt.Println(x, y, btn)
	return nil
}

func drawCell(gc *draw2dimg.GraphicContext, c Cell) {
	col := c.GetColor()
	gc.SetFillColor(col)
	gc.SetStrokeColor(col)
	step := cellSize / maxVolume
	liquidVolume := 0.0
	if c.volume > 0 && c.cellType == Liquid {
		liquidVolume = cellSize - step * c.GetVolume()
	}
	gc.BeginPath()
	draw2dkit.Rectangle(
		gc,
		float64(c.x) * cellSize,
		float64(c.y) * cellSize + liquidVolume,
		float64(c.x + 1) * cellSize,
		float64(c.y + 1) * cellSize,
		)
	gc.FillStroke()
	gc.Close()
}

func Render(gc *draw2dimg.GraphicContext) bool {
	simulate(cells)
	simulate(cells)
	simulate(cells)
	gc.SetFillColor(color.RGBA{0xff, 0xff, 0xff, 0xff})
	gc.Clear()
	for _, row := range cells {
		for _, cell := range row {
			drawCell(gc, cell)
		}
	}
	return true
}

