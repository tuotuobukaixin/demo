// Copyright 2014 beego Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package orm

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// oracle operators.
var oracleOperators = map[string]string{
	"exact":     "= :1",
	"iexact":    "LIKE :1",
	"contains":  "LIKE BINARY :1",
	"icontains": "LIKE :1",
	// "regex":       "REGEXP BINARY :1",
	// "iregex":      "REGEXP :1",
	"gt":          "> :1",
	"gte":         ">= :1",
	"lt":          "< :1",
	"lte":         "<= :1",
	"startswith":  "LIKE BINARY :1",
	"endswith":    "LIKE BINARY :1",
	"istartswith": "LIKE :1",
	"iendswith":   "LIKE :1",
}

// oracle column field types.
var oracleTypes = map[string]string{
	"auto":            "NOT NULL PRIMARY KEY",
	"pk":              "NOT NULL PRIMARY KEY",
	"bool":            "char(1)",
	"string":          "varchar2(%d)",
	"string-text":     "clob",
	"time.Time-date":  "date",
	"time.Time":       "timestamp",
	"int8":            "number(2)",
	"int16":           "number(4)",
	"int32":           "number(9)",
	"int64":           "number(18)",
	"uint8":           "number(2)",
	"uint16":          "number(4)",
	"uint32":          "number(9)",
	"uint64":          "number(18)",
	"float32":         "numeric(%d, %d)",
	"float64":         "numeric(%d, %d)",
	"float64-decimal": "numeric(%d, %d)",
}

// oracle dbBaser
type dbBaseOracle struct {
	dbBase
}

var _ dbBaser = new(dbBaseOracle)

// create oracle dbBaser.
func newdbBaseOracle() dbBaser {
	b := new(dbBaseOracle)
	b.ins = b
	return b
}

// get oracle operator.
func (d *dbBaseOracle) OperatorSql(operator string) string {
	return oracleOperators[operator]
}

// get oracle table field types.
func (d *dbBaseOracle) DbTypes() map[string]string {
	return oracleTypes
}

// oracle quote is blank.
func (d *dbBaseOracle) TableQuote() string {
	return ""
}

// show table sql for oracle.
func (d *dbBaseOracle) ShowTablesQuery() string {
	return "SELECT table_name FROM user_tab_comments"
}

// show columns sql of table for oracle.
func (d *dbBaseOracle) ShowColumnsQuery(table string) string {
	return fmt.Sprintf("select column_name, data_type, nullable from user_tab_columns where table_name = '%s'", table)
}

// execute sql to check index exist.
func (d *dbBaseOracle) IndexExists(db dbQuerier, table string, name string) bool {
	sql := fmt.Sprintf("select count(*) from user_ind_columns t,user_indexes i where "+
		"t.index_name = i.index_name and t.table_name = i.table_name and t.table_name = %s "+
		"and i.INDEX_NAME = %s", strings.ToUpper(table), name)
	row := db.QueryRow(sql)
	var cnt int
	if err := row.Scan(&cnt); err != nil {
		fmt.Println("select index error:", err)
	}
	return cnt > 0
}

/*func (d *dbBaseOracle) SuitForKeyWords(fi *fieldInfo) {
	if value, ok := oracleKeyWords[strings.ToLower(fi.column)]; ok {
		fi.column = value
	}
}*/

// oracle value placeholder is :n.
// replace default ? to :n.
func (d *dbBaseOracle) ReplaceMarks(query *string) {
	q := *query
	num := 0
	for _, c := range q {
		if c == '?' {
			num += 1
		}
	}
	if num == 0 {
		return
	}
	data := make([]byte, 0, len(q)+num)
	num = 1
	for i := 0; i < len(q); i++ {
		c := q[i]
		if c == '?' {
			data = append(data, ':')
			data = append(data, []byte(strconv.Itoa(num))...)
			num += 1
		} else {
			data = append(data, c)
		}
	}
	*query = string(data)
}

func (d *dbBaseOracle) GetLimitSql(query string, offset, limit int64) (sql string) {
	if limit <= 0 {
		limit = int64(DefaultRowsLimit)
	}

	if strings.Contains(query, "WHERE") {
		if strings.Contains(query, "ORDER BY") {
			seps := strings.Split(query, "ORDER BY")
			sql = fmt.Sprintf("%s AND rownum BETWEEN %d AND %d ORDER BY%s", seps[0], offset, limit, seps[1])
		} else {
			sql = fmt.Sprintf("%s AND rownum BETWEEN %d AND %d", query, offset, limit)
		}
	} else {
		if strings.Contains(query, "ORDER BY") {
			seps := strings.Split(query, "ORDER BY")
			sql = fmt.Sprintf("%s WHERE rownum BETWEEN %d AND %d ORDER BY%s", seps[0], offset, limit, seps[1])
		} else {
			sql = fmt.Sprintf("%s WHERE rownum BETWEEN %d AND %d", query, offset, limit)
		}
	}

	return
}

func (d *dbBaseOracle) SuitForNull(values []interface{}) {
	for i, value := range values {
		if value == `` {
			switch reflect.TypeOf(value).Kind() {
			case reflect.String:
				values[i] = " "
			}
		}
	}
}
