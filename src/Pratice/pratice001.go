package Pratice

import (
	"fmt"
	"strconv"
)

type myInt struct {
	value int
}
baidufunc (a myInt) Less(b int) (ret bool) {
	return a.value < b
}
 
func (a *myInt) Add(b int) {
	a.value += b
}

func (a myInt) String() (str string) {
	return "show value= " + strconv.Itoa(a.value)
}

type LessAdder interface {
	Less(b int) (ret bool)
	Add(b int)
	String() (str string)
}

const (
	White = iota
	Black
	Blue
	Red
	Yellow
)

type Color byte

func (c Color) ColorToString() (ret string) {
	nameList := []string{"white", "black", "blue", "red", "yellow"}
	return nameList[c]
}

type Box struct {
	Width, Height, Depth float64
	Color                Color
}

func (box Box) Cal() float64 {
	return box.Width * box.Height * box.Depth
}

func (box *Box) SetColor(newColor Color) {
	box.Color = newColor
}

type BoxList []Box

func (bList BoxList) GetBiggestColor() (color Color) {

	var a myInt = myInt{value: 1}
	var b LessAdder = &a
	b.Add(2)
	fmt.Println(b)

	biggestValue := 0.00
	var biggestColor Color = 0
	for _, v := range bList {
		if thisCal := v.Cal(); thisCal > biggestValue {
			biggestValue = thisCal
			biggestColor = v.Color
		}
	}
	return biggestColor
}

func TestCanSee() {

}
