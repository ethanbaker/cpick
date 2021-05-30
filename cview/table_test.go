package cview

import (
	"fmt"
	"testing"
)

var tableTestCases = generateTableTestCases()

type tableTestCase struct {
	rows         int
	columns      int
	fixedRows    int
	fixedColumns int
}

func (c *tableTestCase) String() string {
	return fmt.Sprintf("Rows=%d/Cols=%d/FixedRows=%d/FixedCols=%d", c.rows, c.columns, c.fixedRows, c.fixedColumns)
}

func TestTable(t *testing.T) {
	t.Parallel()

	for _, c := range tableTestCases {
		c := c // Capture

		t.Run(c.String(), func(t *testing.T) {
			t.Parallel()

			table := tc(c)

			app, err := newTestApp(table)
			if err != nil {
				t.Errorf("failed to initialize Application: %s", err)
			}

			for row := 0; row < c.rows; row++ {
				for column := 0; column < c.columns; column++ {
					contents := table.GetCell(row, column).GetText()
					expected := fmt.Sprintf("%d,%d", column, row)
					if contents != expected {
						t.Errorf("failed to either get or set TableCell text: expected %s, got %s", expected, contents)
					}
				}
			}

			table.Draw(app.screen)

			table.Clear()
		})
	}
}

func BenchmarkTableDraw(b *testing.B) {
	for _, c := range tableTestCases {
		c := c // Capture

		b.Run(c.String(), func(b *testing.B) {
			table := tc(c)

			app, err := newTestApp(table)
			if err != nil {
				b.Errorf("failed to initialize Application: %s", err)
			}

			table.Draw(app.screen)

			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				table.Draw(app.screen)
			}
		})
	}
}

func generateTableTestCases() []*tableTestCase {
	var cases []*tableTestCase
	for i := 1; i < 3; i++ {
		rows := i * 5
		for i := 1; i < 3; i++ {
			columns := i * 7
			for fixedRows := 0; fixedRows < 3; fixedRows++ {
				for fixedColumns := 0; fixedColumns < 3; fixedColumns++ {
					cases = append(cases, &tableTestCase{rows, columns, fixedRows, fixedColumns})
				}
			}
		}
	}
	return cases
}

func tc(c *tableTestCase) *Table {
	table := NewTable()

	for row := 0; row < c.rows; row++ {
		for column := 0; column < c.columns; column++ {
			table.SetCellSimple(row, column, fmt.Sprintf("%d,%d", column, row))
		}
	}

	table.SetFixed(c.fixedRows, c.fixedColumns)

	return table
}
