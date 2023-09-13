package main

import (
	"fmt"
	"math"
	"strings"
)

const columnSeparatorLen = 4

var columnSeparator = strings.Repeat(" ", columnSeparatorLen)

type lenStringer interface {
	fmt.Stringer
	Len() int
}

func gridSingleColumn(dispNames []lenStringer, width int, _ int) ([][]string, gridLayout) {
	tgtColumns := 1

	layout := newGridLayout(dispNames, tgtColumns, width)
	names := grid(dispNames, layout)
	return names, layout
}

func gridMultiColumn(dispNames []lenStringer, width int, height int) ([][]string, gridLayout) {
	// Target number of columns to use 1/3 of the available height.
	tgtColumns := len(dispNames) / (height / 3)
	if tgtColumns < 1 {
		tgtColumns = 1
	}

	layout := newGridLayout(dispNames, tgtColumns, width)
	names := grid(dispNames, layout)
	return names, layout
}

func grid(dispNames []lenStringer, layout gridLayout) [][]string {
	names := make([][]string, layout.columns)
	for col := 0; col < layout.columns; col++ {
		colNames := make([]string, layout.rows)
		for row := 0; row < layout.rows; row++ {
			idx := index(col, row, layout.rows)
			if idx < len(dispNames) {
				n := dispNames[idx]
				colNames[row] = n.String() + strings.Repeat(" ", layout.maxColumnLen[col]-n.Len())
			}
		}
		names[col] = colNames
	}
	return names
}

// gridLayout defines the shape and constraints of the display grid.
type gridLayout struct {
	rows         int
	columns      int
	maxColumnLen []int
}

// newGridLayout constructs a gridLayout for given display names from a target number of columns.
func newGridLayout(dispNames []lenStringer, tgtColumns int, width int) gridLayout {
	layout := gridLayout{}

tgtLoop:
	// Evaluate if the display names will fit given the target columns and associated layout and
	// continue to decrease the target columns until a fit is found.
	for tgt := tgtColumns; tgt >= 1; tgt-- {
		layout.columns = tgt
		layout.rows = int(math.Ceil(float64(len(dispNames)) / float64(tgt)))
		layout.maxColumnLen = make([]int, tgt)

		for row := 0; row < layout.rows; row++ {
			rowLen := 0
			for col := 0; col < tgt; col++ {
				idx := index(col, row, layout.rows)
				if idx < len(dispNames) {
					curLen := dispNames[idx].Len()
					rowLen += (curLen + columnSeparatorLen)
					// Reject the tgt if it results in any row length greater than the width.
					if rowLen > width && tgt > 1 {
						continue tgtLoop
					}
					if curLen > layout.maxColumnLen[col] {
						layout.maxColumnLen[col] = curLen
					}
				}
			}
		}

		// The layout has not been rejected so break the loop and return.
		break tgtLoop
	}

	return layout
}

func gridRowMajorFixedLayout(items []lenStringer, columns int, rows int) [][]string {
	rowMajorIndex := func(c int, r int, columns int) int {
		return c + (r * columns)
	}

	maxColumnLen := make([]int, columns)
	for col := 0; col < columns; col++ {
		for row := 0; row < rows; row++ {
			idx := rowMajorIndex(col, row, columns)
			if idx < len(items) {
				curLen := items[idx].Len()
				if curLen > maxColumnLen[col] {
					maxColumnLen[col] = curLen
				}
			}
		}
	}

	names := make([][]string, rows)
	for row := 0; row < rows; row++ {
		rowNames := make([]string, columns)
		for col := 0; col < columns; col++ {
			idx := rowMajorIndex(col, row, columns)
			if idx < len(items) {
				n := items[idx]
				rowNames[col] = n.String() + strings.Repeat(" ", maxColumnLen[col]-n.Len())
			}
		}
		names[row] = rowNames
	}
	return names
}
