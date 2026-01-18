package csvdb

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
)

func (cf CsvFile) ColumnNames() ([]string, error) {

	stmt := fmt.Sprintf("SELECT * FROM %s LIMIT 0", cf.tableName)

	rows, err := cf.database.Query(stmt)
	if err != nil {
		return nil, errors.New("error: query not executed")
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		log.Fatal(err)
	}
	return columns, nil
}

func (cf CsvFile) ColumnInfo() ([]ColumnInfo, error) {
	stmt := fmt.Sprintf("DESCRIBE %s", cf.tableName)

	rows, err := cf.database.Query(stmt)
	if err != nil {
		return nil, errors.New("error: query not executed")
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			panic(err)
		}
	}(rows)

	var columnInfo []ColumnInfo
	var nullable string
	var primaryKey sql.NullString
	var defaultValue sql.NullString
	var extraInfo sql.NullString
	for rows.Next() {
		var row ColumnInfo
		if err := rows.Scan(&row.Name, &row.Type, &nullable, &primaryKey, &defaultValue, &extraInfo); err != nil {
			return columnInfo, err
		}
		columnInfo = append(columnInfo, row)
	}
	return columnInfo, nil
}

func (cf CsvFile) UniqueCounts(column string) ([]KeyValue[int64], error) {
	stmt := fmt.Sprintf("SELECT %s, COUNT(*) AS RESULT FROM %s GROUP BY %s ORDER BY RESULT DESC",
		column,
		cf.tableName,
		column)

	rows, err := cf.database.Query(stmt)
	if err != nil {
		return nil, errors.New("error: query not executed")
	}
	defer rows.Close()

	uniqueCounts := []KeyValue[int64]{}
	for rows.Next() {
		var keyValue KeyValue[int64]
		if err := rows.Scan(&keyValue.Key, &keyValue.Value); err != nil {
			return uniqueCounts, err
		}
		uniqueCounts = append(uniqueCounts, keyValue)
	}
	return uniqueCounts, nil
}
