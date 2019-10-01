package main

import (
	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/llgcode/draw2d/draw2dkit"
	"github.com/markfarnan/go-canvas/canvas"
	"image/color"
	"syscall/js"
)

var cells [][]Cell

var done chan struct{}

var cvs *canvas.Canvas2d
const width = 40
const height = 20
const cellSize = 20.0
var pressed = false


func main() {

	cvs, _ = canvas.NewCanvas2d(true)
	cvs_obj := js.
		Global().
		Get("document").
		Call("getElementsByTagName", "canvas").
		Index(0)
	cvs_obj.Call("addEventListener", "mousemove", js.FuncOf(handleMouseMove))
	cvs_obj.Call("addEventListener", "mousedown", js.FuncOf(handleMouseDown))
	cvs_obj.Call("addEventListener", "mouseup"  , js.FuncOf(handleMouseUp  ))

	for i := 0; i < width; i++ {
		var column []Cell
		for j := 0; j < height; j++ {
			var t CellType = Liquid
			var v = 0.0
			if j == height - 1 || j == 0 || i == 0 || i == width - 1 {
				t = Solid
			} else if false {
				t = Liquid
				v = 1.0
			}
			column = append(column, NewCell(i, j, t, v))
		}
		cells = append(cells, column)
	}

	cvs.Start(60, Render)
	<-done
}


func handleMouseDown(this js.Value, args []js.Value) interface{} {
	pressed = true
	return nil
}

func handleMouseUp(this js.Value, args []js.Value) interface{} {
	pressed = false
	return nil
}

func handleMouseMove(this js.Value, args []js.Value) interface{} {
	mouseEvent := args[0]
	js.Global().Get("console").Call("log", mouseEvent)
	shift := mouseEvent.Get("shiftKey").Bool()
	ctrl := mouseEvent.Get("ctrlKey").Bool()
	if pressed {
		x := int((mouseEvent.Get("offsetX").Float()) / cellSize)
		y := int((mouseEvent.Get("offsetY").Float()) / cellSize)
		if x > width - 1 || y > height - 1 || y < 1 || x < 1 {
			return nil
		}
		if !ctrl && !shift {
			cells[x][y].cellType = Liquid
			cells[x][y].volume = maxVolume + maxPressure
		} else if ctrl && !shift {
			cells[x][y].cellType = Solid
		} else if shift && !ctrl {
			cells[x][y].cellType = Void
			cells[x][y].volume = 0
		}
	}
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

