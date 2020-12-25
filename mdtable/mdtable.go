// Package mdtable helps construct markdown tables.
package mdtable

import (
	"encoding/csv"
	"fmt"
	"strings"
	"unicode/utf8"
)

// Alignment specifies a column alignment.
type Alignment string

// Valid Alignment values.
const (
	DefaultAlignment Alignment = ""
	Left             Alignment = ":--"
	Right            Alignment = "--:"
	Center           Alignment = ":-:"
)

// Builder helps construct a Markdown table.
type Builder struct {
	header     []string
	alignments []Alignment
	rows       [][]string
}

// SetAlignment sets the alignment of the columns.
func (b *Builder) SetAlignment(colAlignments []Alignment) *Builder {
	b.alignments = colAlignments
	return b
}

// SetHeader sets the header cell contents of the table.
func (b *Builder) SetHeader(cells []string) *Builder {
	b.header = cells
	return b
}

// AddRow adds a given markdown literals as contents of a row.
func (b *Builder) AddRow(row []string) *Builder {
	b.rows = append(b.rows, row)
	return b
}

func (b *Builder) Build() string {
	colLengths := maxColLengths(b.allRows())
	formatRow := func(row []string) string {
		formattedCells := []string{}
		for i := 0; i < len(colLengths); i++ {
			contents := "  "
			if i < len(row) {
				contents = " " + row[i] + " "
			}
			formattedCells = append(formattedCells, padRight(contents, colLengths[i]+2))
		}
		return fmt.Sprintf("|%s|", strings.Join(formattedCells, "|"))
	}
	lines := []string{
		formatRow(b.header),
		b.renderAlignmentRow(colLengths),
	}
	for _, row := range b.rows {
		lines = append(lines, formatRow(row))
	}

	return strings.Join(lines, "\n")
}

// BuildCSV returns the table as a CSV.
func (b *Builder) BuildCSV() string {
	sb := &strings.Builder{}
	w := csv.NewWriter(sb)

	w.Write(b.header)
	w.WriteAll(b.rows)

	w.Flush()
	return sb.String()
}

func (b *Builder) allRows() [][]string {
	out := [][]string{b.header}
	out = append(out, b.rows...)
	return out
}

func (b *Builder) renderAlignmentRow(colLengths []int) string {
	cells := []string{}
	for i := 0; i < len(colLengths); i++ {
		align := DefaultAlignment
		if i < len(b.alignments) {
			align = b.alignments[i]
		}
		dashes := repeatString("-", colLengths[i])
		cell := ""
		switch align {
		case DefaultAlignment:
			cell = fmt.Sprintf("-%s-", dashes)
		case Left:
			cell = fmt.Sprintf(":%s-", dashes)
		case Right:
			cell = fmt.Sprintf("-%s:", dashes)
		case Center:
			cell = fmt.Sprintf(":%s:", dashes)
		}
		cells = append(cells, cell)
	}
	return fmt.Sprintf("|%s|", strings.Join(cells, "|"))
}

func repeatString(str string, times int) string {
	out := ""
	for i := 0; i < times; i++ {
		out += str
	}
	return out
}

func padRight(str string, length int) string {
	size := utf8.RuneCountInString(str)
	if size >= length {
		return str
	}
	return str + repeatString(" ", length-size)
}

func maxColLengths(rows [][]string) []int {
	lengths := make(map[int]int)
	colCount := 0
	for _, row := range rows {
		colCount = maxInt(len(row), colCount)
		for col, contents := range row {
			size := utf8.RuneCountInString(contents)
			if size > lengths[col] {
				lengths[col] = maxInt(lengths[col], size)
			}
		}
	}
	var out []int
	for i := 0; i < colCount; i++ {
		out = append(out, lengths[i])
	}
	return out
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
