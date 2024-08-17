package csvdb

import (
	"database/sql"
	"errors"
	"fmt"
)

func (cf csvFile) ColumnStats() ([]ColumnStats, error) {
	stmt := fmt.Sprintf("SUMMARIZE %s", cf.tableName)

	rows, err := cf.database.Query(stmt)
	if err != nil {
		return nil, errors.New("error: query not executed")
	}
	defer rows.Close()

	var columnStats []ColumnStats

	var meanValue sql.NullFloat64
	var standardDev sql.NullFloat64
	var quantile25 sql.NullFloat64
	var quantile50 sql.NullFloat64
	var quantile75 sql.NullFloat64
	var nullPercentage any

	for rows.Next() {
		var row ColumnStats
		if err := rows.Scan(&row.Name, &row.Type, &row.Minimum,
			&row.Maximum, &row.Unique, &meanValue, &standardDev,
			&quantile25, &quantile50, &quantile75, &row.Count, &nullPercentage); err != nil {
			return columnStats, err
		}
		if meanValue.Valid {
			row.Mean = round(meanValue.Float64, 2)
		}
		if standardDev.Valid {
			row.StandardDev = round(standardDev.Float64, 2)
		}
		if quantile25.Valid {
			row.Q25 = round(quantile25.Float64, 2)
		}
		if quantile50.Valid {
			row.Q50 = round(quantile50.Float64, 2)
		}
		if quantile75.Valid {
			row.Q75 = round(quantile75.Float64, 2)
		}
		columnStats = append(columnStats, row)
	}
	return columnStats, nil
}

func (cf csvFile) getStat(stmt string) (float64, error) {
	var result float64

	row := cf.database.QueryRow(stmt)
	if err := row.Scan(&result); err != nil {
		if err == sql.ErrNoRows {
			return 0.0, errors.New("error: no results")
		}
		return 0.0, errors.New("error: not a numerical column")
	}
	return round(result,2), nil
}

func (cf csvFile) Mean(column string) (float64, error) {

	stmt := fmt.Sprintf("SELECT avg(%s) AS RESULT FROM %s", column, cf.tableName)
	return cf.getStat(stmt)
}

func (cf csvFile) Median(column string) (float64, error) {

	stmt := fmt.Sprintf("SELECT quantile(%s, 0.5) AS RESULT FROM %s", column, cf.tableName)
	return cf.getStat(stmt)
}

func (cf csvFile) Sum(column string) (float64, error) {
	stmt := fmt.Sprintf("SELECT sum(%s) AS RESULT FROM %s", column, cf.tableName)
	return cf.getStat(stmt)
}

func (cf csvFile) Min(column string) (float64, error) {
	stmt := fmt.Sprintf("SELECT min(%s) AS RESULT FROM %s", column, cf.tableName)
	return cf.getStat(stmt)
}

func (cf csvFile) Max(column string) (float64, error) {
	stmt := fmt.Sprintf("SELECT max(%s) AS RESULT FROM %s", column, cf.tableName)
	return cf.getStat(stmt)
}

func (cf csvFile) NullCount(column string) (int64, error) {
	stmt := fmt.Sprintf("SELECT COUNT(*) - COUNT(%s) AS RESULT FROM %s", column, cf.tableName)
	result, err := cf.getStat(stmt)
	return int64(result), err
}