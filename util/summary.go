package util

import (
	"bytes"
	"fmt"
)

//------------------------------------------------------------
// Summary displays data as a formatted table
//------------------------------------------------------------

type Summary struct {
	Header        string
	Columns       []string
	ColumnsData   map[string]*ColumnData
	ColumnPadding int
    LineCount     int
	TotalWidth    int
}

type ColumnData struct {
	Lines    []string
	Width int
}

// Creates new summary table with name and column names.
func NewSummary(header string, columns ...string) (sum *Summary) {

	sum = &Summary{
		Header:        header,
		Columns:       columns,
		ColumnsData:   map[string]*ColumnData{},
		ColumnPadding: 8,
	}

	for _, col := range columns {
		sum.ColumnsData[col] = &ColumnData{}
	}

	return sum
}

// Adds empty line of values.
func (s *Summary) AddBlankLine() {
    s.AddLine(make([]string, len(s.Columns))...)
}

// Adds table line of values.
func (s *Summary) AddLine(vals ...string) {

	if len(vals) != len(s.Columns) {
		panic(fmt.Sprintf("Column count must match: %v (created) != %v (adding)", len(s.Columns), len(vals)))
	}

    lineWidth := 0

	for i, col := range s.Columns {
		cdata := s.ColumnsData[col]
		val := vals[i]
		cdata.Lines = append(cdata.Lines, val)

		if len(val) > cdata.Width {
			cdata.Width = len(val)
        }

        lineWidth += cdata.Width + s.ColumnPadding
	}

    if lineWidth > s.TotalWidth {
        s.TotalWidth = lineWidth
    }

    s.LineCount++
}

// String.
func (s *Summary) String() string {

	var buf bytes.Buffer

	buf.WriteString(s.Header)
    buf.WriteByte('\n')

    // Separator line
    s.bufPad(&buf, s.TotalWidth, '-')
    buf.WriteByte('\n')

    // Column names
    for _, col := range s.Columns {
        buf.WriteString(col)
        cdata := s.ColumnsData[col]
        s.bufPad(&buf, cdata.Width - len(col) + s.ColumnPadding, ' ')
    }
    buf.WriteByte('\n')

    // Separator line
    s.bufPad(&buf, s.TotalWidth, '-')
    buf.WriteByte('\n')

    // Each line
    for l := 0; l < s.LineCount; l++ {

        // Each column
        for _, col := range s.Columns {
            cdata := s.ColumnsData[col]
            val := cdata.Lines[l]

            buf.WriteString(val)
            s.bufPad(&buf, cdata.Width - len(val) + s.ColumnPadding, ' ')
        }
        buf.WriteByte('\n')
    }

    // Separator line
    s.bufPad(&buf, s.TotalWidth, '-')
    buf.WriteByte('\n')

	return buf.String()
}

// Writes N copies of B into buffer.
func (s *Summary) bufPad(buf *bytes.Buffer, n int, b byte) {
    for i := 0; i < n; i++ {
        buf.WriteByte(b)
    }
}
