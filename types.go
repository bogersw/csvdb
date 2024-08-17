package csvdb

import (
	"database/sql"
)

type csvFile struct {
	Options   Options
	FileName  string
	database  *sql.DB
	tableName string
}

type Options struct {
	DateFormat       string
	DecimalSeparator string
	Header           bool
	IgnoreErrors     bool
	RowsToRead       int64
	Separator        string
	SampleSize       int64
}

type ColumnInfo struct {
	Name string
	Type string
}

type ColumnStats struct {
	Name        string
	Type        string
	Minimum     any
	Maximum     any
	Unique      int64
	Mean        float64
	StandardDev float64
	Q25         float64
	Q50         float64
	Q75         float64
	Count       int64
}

// Note: important date format specifiers
// %d zero-padded day of the month
// %m zero-padded month of the year
// %y zero-padded year without century
// %Y year with century
