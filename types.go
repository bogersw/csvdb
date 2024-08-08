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
	RowsToRead       int
	Separator        string
}

// Note: important date format specifiers
// %d zero-padded day of the month
// %m zero-padded month of the year
// %y zero-padded year without century
// %Y year with century
