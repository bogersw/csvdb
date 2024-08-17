# csvdb

Go package to read csv files into a duckdb database for further analysis. In development!

Installation: `go get github.com/bogersw/csvdb`.

# Types in package

```
type Options struct {
	DateFormat       string
	DecimalSeparator string
	Header           bool
	IgnoreErrors     bool
	RowsToRead       int64
	Separator        string
	SampleSize       int64
}
```

```
type csvFile struct {
	Options   Options
	FileName  string
	database  *sql.DB
	tableName string
}
```

```
type ColumnInfo struct {
	Name string
	Type string
}
```

```
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
```

# Usage

```
import (
    "fmt"

    "github.com/bogersw/csvdb"
)

options := Options{
    DecimalSeparator: ",",
    Header:           true,
    Separator:        ";",
}

csv := csvdb.NewCsv(options)

if err := csv.Read(); err != nil {
    fmt.Println(err)
    return
}

csv.SetFileName(".", "data", "file.csv")
...
```

The returned `csv` is of type `csvFile` (struct). The following methods are associated with this type:

**Settings**

-   `SetFileName(parts ...string) error`

**Information**

-   `ColumnInfo() ([]ColumnInfo, error)`: returns slice of structs with basic column info.
-   `ColumnNames() ([]string, error)`: returns string slice with column names.
-   `ColumnStats() ([]ColumnStats, error)`: returns slice of structs with more extended column info.

**Statistics**

-   `Max(column string) (float64, error)`: returns maximum value of specified column.
-   `Mean(column string) (float64, error)`: returns mean value of specified column.
-   `Median(column string) (float64, error)`: returns median value of specified column.
-   `Min(column string) (float64, error)`: returns minimum value of specified column.
-   `NullCount(column string) (int64, error)`: returns null count of specified column.
-   `Sum(column string) (float64, error)`: returns sum of specified column.
