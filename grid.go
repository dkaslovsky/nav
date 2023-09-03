package main

import (
	"math"
	"strings"
)

const separator = "    " // Separator between columns.

var separatorLen = len(separator)

func grid(dispNames []*displayName, width int, height int) ([][]string, gridLayout) {
	// Target number of columns to use 1/3 of the available height.
	tgtColumns := len(dispNames) / (height / 3)
	if tgtColumns < 1 {
		tgtColumns = 1
	}

	layout := newGridLayout(dispNames, tgtColumns, width, height)
	names := make([][]string, layout.columns)

	for col := 0; col < layout.columns; col++ {
		colNames := make([]string, layout.rows)
		for row := 0; row < layout.rows; row++ {
			idx := row + (col * layout.rows)
			if idx < len(dispNames) {
				n := dispNames[idx]
				colNames[row] = n.String() + strings.Repeat(" ", layout.maxColumnLen[col]-n.Len())
			}
		}
		names[col] = colNames
	}
	return names, layout
}

// gridLayout defines the shape and constraints of the display grid.
type gridLayout struct {
	rows         int
	columns      int
	maxColumnLen []int
}

// newGridLayout constructs a gridLayout for given display names from a target number of columns.
func newGridLayout(dispNames []*displayName, tgtColumns int, width int, height int) gridLayout {
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
				idx := row + (col * layout.rows)
				if idx < len(dispNames) {
					curLen := dispNames[idx].Len()
					rowLen += (curLen + separatorLen)
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
