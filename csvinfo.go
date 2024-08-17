package csvdb

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
)

func (cf csvFile) ColumnNames() ([]string, error) {

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

func (cf csvFile) ColumnInfo() ([]ColumnInfo, error) {
	stmt := fmt.Sprintf("DESCRIBE %s", cf.tableName)

	rows, err := cf.database.Query(stmt)
	if err != nil {
		return nil, errors.New("error: query not executed")
	}
	defer rows.Close()

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
