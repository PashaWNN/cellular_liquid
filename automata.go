package main

import (
	"image/color"
	"math"
)

type CellType int

const (
	Void = iota
	Liquid
	Solid
)


const maxSpeed = 2.0
const maxVolume = 1.0
const maxPressure = 0.15
const minFlow = 0.001
const flowCoef = 0.6
const minVolume = 0.0001  // Treat cells with vol < minVol as void

const colorSlope = (50 - 0) / (1.0 - 0.0)

type Cell struct {
	x, y int
	cellType CellType
	volume float64
	deltaVolume float64
}

func (c* Cell) GetColor() color.RGBA {
	if c.cellType == Liquid {
		if c.volume <= 0 {
			return color.RGBA{0xff, 0xff, 0xff, 0xff}
		}
		col := uint8(colorSlope * c.GetPressure())
		return color.RGBA{
			0xcc - col,
			0xe6 - col,
			255 - col,
			0xff}
	} else if c.cellType == Solid {
		return color.RGBA{0x00, 0x00, 0x33, 0xff}
	}
	return color.RGBA{0xff, 0xff, 0xff, 0xff}
}

func (c* Cell) GetPressure() float64 {
	pressure := c.volume - maxVolume
	if pressure < 0 {
		pressure = 0
	}
	return pressure
}

func (c* Cell) GetVolume() float64 {
	volume := c.volume
	if volume > maxVolume {
		volume = maxVolume
	}
	return volume
}

func (c* Cell) Tick() {
	c.volume += c.deltaVolume
	c.deltaVolume = 0
}

func NewCell(x, y int, cellType CellType, volume float64) Cell {
	return Cell{x, y, cellType, volume, 0}
}


func simulate(cells [][]Cell) {

	for x := 1; x < len(cells) - 1; x++ {
		for y := 1; y < len(cells[x]) - 1; y++ {
			if cells[x][y].cellType == Solid {
				continue
			}

			var this   *Cell = &cells[x  ][y  ]
			var left   *Cell = &cells[x-1][y  ]
			var right  *Cell = &cells[x+1][y  ]
			var up     *Cell = &cells[x  ][y-1]
			var bottom *Cell = &cells[x  ][y+1]

			if FallDown(this, bottom) {
				if Spread(this, left, right) {
					Decompress(this, up)
				}
			}
		}
	}
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {

			if cells[x][y].cellType != Solid {

				if cells[x][y].volume > minVolume {
					cells[x][y].cellType = Liquid
				} else {
					cells[x][y].cellType = Void
				}
			}

			cells[x][y].Tick()
		}
	}
}


func FallDown(c1* Cell, c2* Cell) bool {
	if c2.cellType != Solid {
		remaining_mass := c1.volume - c1.deltaVolume
		flow := getStableState(remaining_mass + c2.volume) - c2.volume
		flow *= 2.0
		if flow > minFlow { flow *= flowCoef }
		flow = constrain(flow, 0, math.Min(maxSpeed, remaining_mass))
		c1.deltaVolume -= flow
		c2.deltaVolume += flow
	}
	return c1.volume - c1.deltaVolume > 0
}


func Spread(c, left, right *Cell) bool {
	if left.cellType != Solid {
		flow := (c.volume - left.volume) / 1
		if flow > minFlow { flow *= flowCoef }
		flow = constrain(flow, 0, c.volume - c.deltaVolume)

		c.volume -= flow
		left.volume += flow
	}
	if c.volume - c.deltaVolume <= 0 { return false }
	if right.cellType != Solid {
		flow := (c.volume - right.volume) / 1
		if flow > minFlow { flow *= flowCoef }
		flow = constrain(flow, 0, c.volume - c.deltaVolume)

		c.volume -= flow
		right.volume += flow
	}
	return c.volume - c.deltaVolume > 0
}


func Decompress(c, up *Cell) bool {

	if c.cellType != Solid {
		remaining_mass := c.volume - c.deltaVolume
		flow := remaining_mass - getStableState(remaining_mass + up.volume)
		flow /= 2.0
		if flow > minFlow { flow *= flowCoef }
		flow = constrain(flow, 0, math.Min(maxSpeed, remaining_mass))

		c.deltaVolume -= flow
		up.deltaVolume += flow
	}

	return c.volume - c.deltaVolume > 0
}


func getStableState(totalMass float64) float64 {
	if totalMass <= 1 {
		return 1
	} else if totalMass < 2 * maxVolume + maxPressure {
		return (maxVolume * maxVolume + totalMass * maxPressure) / maxVolume - maxPressure
	} else {
		return (totalMass + maxPressure) / 2
	}
}


func constrain(x, a, b float64) float64 {
	result := x
	if x > b {
		result = b
	}
	if result < a {
		result = a
	}
	return result
}