package csvdb

import (
	"database/sql"
	"errors"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/bogersw/csvdb/internal/db"
)

var duckDb *sql.DB = db.DuckDb

func NewCsv(options Options) *csvFile {
	return &csvFile{
		Options:   options,
		FileName:  "",
		database:  duckDb,
		tableName: "",
	}
}

func (c csvFile) Read() error {
	// Convert options to valid DuckDb options
	options := []string{}
	if c.Options.DateFormat != "" {
		options = append(options, fmt.Sprintf("dateformat='%s'", c.Options.DateFormat))
	}
	if c.Options.DecimalSeparator != "" {
		options = append(options, fmt.Sprintf("decimal_separator='%s'", c.Options.DecimalSeparator))
	}
	if c.Options.Header {
		options = append(options, "header=true")
	}
	if c.Options.IgnoreErrors {
		options = append(options, "ignore_errors=true")
	}
	if c.Options.Separator != "" {
		options = append(options, fmt.Sprintf("delim='%s'", c.Options.Separator))
	}
	if c.Options.SampleSize > 0 {
		options = append(options, fmt.Sprintf("sample_size=%d", c.Options.SampleSize))
	}
	stmt := fmt.Sprintf("CREATE TABLE %s AS FROM read_csv('%s', %s)",
		c.tableName,
		c.FileName,
		strings.Join(options, ","))
	if c.Options.RowsToRead > 0 {
		stmt += fmt.Sprintf(" LIMIT %d", c.Options.RowsToRead)
	}
	_, err := c.database.Exec(stmt)
	if err != nil {
		if strings.Contains(err.Error(), "Invalid unicode") {
			return errors.New("error: encoding not UTF-8")
		} else {
			return errors.New("error: file could not be read")
		}
	}
	return nil
}

func (c *csvFile) SetFileName(parts ...string) error {
	// Build the path and check if file exists
	fileName := filepath.Join(parts...)
	_, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		return errors.New("error: file does not exist")
	}
	c.FileName = fileName
	// Get filename without extensions as the table name
	c.tableName = strings.Split(filepath.Base(fileName), ".")[0]
	return nil
}

// func (cf csvFile) query(stmt string) ([][]interface{}, error) {

// 	rows, err := cf.database.Query(stmt)
// 	if err != nil {
// 		return nil, errors.New("error: query execution failed")
// 	}
// 	defer rows.Close()

// 	// Get the column names
// 	columns, err := rows.Columns()
// 	if err != nil {
// 		return nil, errors.New("error: column names could not be determined")
// 	}

// 	// Create a slice to store the results
// 	var results [][]interface{}

// 	// Iterate over the rows and append the values to the results slice
// 	for rows.Next() {
// 		values := make([]interface{}, len(columns))
// 		valuePtrs := make([]interface{}, len(columns))
// 		for i := range columns {
// 			valuePtrs[i] = &values[i]
// 		}
// 		err := rows.Scan(valuePtrs...)
// 		if err != nil {
// 			return nil, errors.New("error: rows could not be read")
// 		}
// 		results = append(results, values)
// 	}
// 	return results, nil
// }

func round(val float64, precision int) float64 {
	p := math.Pow10(precision)
	return math.Round(val*p) / p
}
