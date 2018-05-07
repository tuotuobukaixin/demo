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
	"os"
	"strings"
)

type dbIndex struct {
	Table string
	Name  string
	Sql   string
}

// create database drop sql.
func getDbDropSql(al *alias) (sqls []string) {
	if len(modelCache.cache) == 0 {
		fmt.Println("no Model found, need register your model")
		os.Exit(2)
	}

	Q := al.DbBaser.TableQuote()

	for _, mi := range modelCache.allOrdered() {
		sqls = append(sqls, fmt.Sprintf(`DROP TABLE IF EXISTS %s%s%s`, Q, mi.table, Q))
	}
	return sqls
}

// get database column type string.
func getColumnTyp(al *alias, fi *fieldInfo) (col string) {
	T := al.DbBaser.DbTypes()
	fieldType := fi.fieldType

checkColumn:
	switch fieldType {
	case TypeBooleanField:
		col = T["bool"]
	case TypeCharField:
		col = fmt.Sprintf(T["string"], fi.size)
	case TypeTextField:
		col = T["string-text"]
	case TypeDateField:
		col = T["time.Time-date"]
	case TypeDateTimeField:
		col = T["time.Time"]
	case TypeBitField:
		col = T["int8"]
	case TypeSmallIntegerField:
		col = T["int16"]
	case TypeIntegerField:
		col = T["int32"]
	case TypeBigIntegerField:
		if al.Driver == DR_Sqlite {
			fieldType = TypeIntegerField
			goto checkColumn
		}
		col = T["int64"]
	case TypePositiveBitField:
		col = T["uint8"]
	case TypePositiveSmallIntegerField:
		col = T["uint16"]
	case TypePositiveIntegerField:
		col = T["uint32"]
	case TypePositiveBigIntegerField:
		col = T["uint64"]
	case TypeFloatField:
		col = T["float64"]
	case TypeDecimalField:
		s := T["float64-decimal"]
		if strings.Index(s, "%d") == -1 {
			col = s
		} else {
			col = fmt.Sprintf(s, fi.digits, fi.decimals)
		}
	case RelForeignKey, RelOneToOne:
		fieldType = fi.relModelInfo.fields.pk.fieldType
		goto checkColumn
	}

	return
}

// create alter sql string.
func getColumnAddQuery(al *alias, fi *fieldInfo) string {
	Q := al.DbBaser.TableQuote()
	typ := getColumnTyp(al, fi)

	if fi.null == false {
		typ += " " + "NOT NULL"
	}

	var result string
	if al.Driver == DR_Oracle {
		result = fmt.Sprintf("ALTER TABLE %s%s%s ADD COLUMN %s%s%s %s %s",
			Q, fi.mi.table, Q,
			Q, fi.column, Q,
			getColumnDefault(al.Driver, fi), typ,
		)
	} else {
		result = fmt.Sprintf("ALTER TABLE %s%s%s ADD COLUMN %s%s%s %s %s",
			Q, fi.mi.table, Q,
			Q, fi.column, Q,
			typ, getColumnDefault(al.Driver, fi),
		)
	}

	return result
}

// create database creation string.
func getDbCreateSql(al *alias) (sqls []string, tableIndexes map[string][]dbIndex) {
	if len(modelCache.cache) == 0 {
		fmt.Println("no Model found, need register your model")
		os.Exit(2)
	}

	var oracleAuto map[string]string
	Q := al.DbBaser.TableQuote()
	T := al.DbBaser.DbTypes()
	sep := fmt.Sprintf("%s, %s", Q, Q)

	tableIndexes = make(map[string][]dbIndex)

	for _, mi := range modelCache.allOrdered() {
		sql := fmt.Sprintf("-- %s\n", strings.Repeat("-", 50))
		sql += fmt.Sprintf("--  Table Structure for `%s`\n", mi.fullName)
		sql += fmt.Sprintf("-- %s\n", strings.Repeat("-", 50))

		switch al.Driver {
		case DR_Oracle:
			oracleAuto = make(map[string]string)
			sql += fmt.Sprintf("DECLARE cnt NUMBER;\n")
			sql += fmt.Sprintf("BEGIN\n")
			sql += fmt.Sprintf("SELECT count(*) INTO cnt FROM all_tables WHERE table_name = '%s';\n", strings.ToUpper(mi.table))
			sql += fmt.Sprintf("IF cnt = 0 THEN\n")
			sql += fmt.Sprintf("EXECUTE IMMEDIATE 'CREATE TABLE %s%s%s (\n", Q, mi.table, Q)
		default:
			sql += fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s%s%s (\n", Q, mi.table, Q)
		}

		columns := make([]string, 0, len(mi.fields.fieldsDB))

		sqlIndexes := [][]string{}

		for _, fi := range mi.fields.fieldsDB {

			column := fmt.Sprintf("    %s%s%s ", Q, fi.column, Q)
			col := getColumnTyp(al, fi)

			if fi.auto {
				switch al.Driver {
				case DR_Sqlite, DR_Postgres:
					column += T["auto"]
				case DR_Oracle:
					column += col + " " + T["auto"]
					oracleAuto[mi.table] = fi.column
				default:
					column += col + " " + T["auto"]
				}
			} else if fi.pk {
				column += col + " " + T["pk"]
			} else {
				column += col

				if al.Driver == DR_Oracle {
					column += getColumnDefault(al.Driver, fi)
					if fi.null == false {
						column += " " + "NOT NULL"
					}
				} else {
					if fi.null == false {
						column += " " + "NOT NULL"
					}
					column += getColumnDefault(al.Driver, fi)
				}

				if fi.unique {
					column += " " + "UNIQUE"
				}

				if fi.index {
					sqlIndexes = append(sqlIndexes, []string{fi.column})
				}
			}

			if strings.Index(column, "%COL%") != -1 {
				column = strings.Replace(column, "%COL%", fi.column, -1)
			}

			columns = append(columns, column)
		}

		if mi.model != nil {
			allnames := getTableUnique(mi.addrField)
			if !mi.manual && len(mi.uniques) > 0 {
				allnames = append(allnames, mi.uniques)
			}
			for _, names := range allnames {
				cols := make([]string, 0, len(names))
				for _, name := range names {
					if fi, ok := mi.fields.GetByAny(name); ok && fi.dbcol {
						cols = append(cols, fi.column)
					} else {
						panic(fmt.Errorf("cannot found column `%s` when parse UNIQUE in `%s.TableUnique`", name, mi.fullName))
					}
				}
				column := fmt.Sprintf("    UNIQUE (%s%s%s)", Q, strings.Join(cols, sep), Q)
				columns = append(columns, column)
			}
		}

		sql += strings.Join(columns, ",\n")
		if al.Driver == DR_Oracle {
			sql += "\n)';\n"
			//just for oracle auto_increment
			for key, value := range oracleAuto {
				sql += "EXECUTE IMMEDIATE '\n"
				sql += fmt.Sprintf("CREATE SEQUENCE %s_%s\n", key, value)
				sql += "';\nEXECUTE IMMEDIATE '\n"
				sql += fmt.Sprintf("CREATE OR REPLACE TRIGGER %s_%s\n", key, value)
				sql += fmt.Sprintf("BEFORE INSERT ON %s\n", key)
				sql += fmt.Sprintf("FOR EACH ROW\n")
				sql += fmt.Sprintf("BEGIN\n")
				sql += fmt.Sprintf("SELECT %s_%s.nextval\n", key, value)
				sql += fmt.Sprintf("INTO :new.%s\n", value)
				sql += fmt.Sprintf("FROM dual;\n")
				sql += fmt.Sprintf("END;\n';\n")
			}
			sql += "END IF;\n"
			sql += "END;\n"
		} else {
			sql += "\n)"
		}

		if al.Driver == DR_MySQL {
			var engine string
			if mi.model != nil {
				engine = getTableEngine(mi.addrField)
			}
			if engine == "" {
				engine = al.Engine
			}
			sql += " ENGINE=" + engine
		}

		if al.Driver != DR_Oracle {
			sql += ";"
		}

		sqls = append(sqls, sql)

		if mi.model != nil {
			for _, names := range getTableIndex(mi.addrField) {
				cols := make([]string, 0, len(names))
				for _, name := range names {
					if fi, ok := mi.fields.GetByAny(name); ok && fi.dbcol {
						cols = append(cols, fi.column)
					} else {
						panic(fmt.Errorf("cannot found column `%s` when parse INDEX in `%s.TableIndex`", name, mi.fullName))
					}
				}
				sqlIndexes = append(sqlIndexes, cols)
			}
		}

		for _, names := range sqlIndexes {
			name := mi.table + "_" + strings.Join(names, "_")
			cols := strings.Join(names, sep)
			var sql string
			if al.Driver == DR_Oracle {
				if len(name) > 30 {
					name = name[len(name)-30 : len(name)]
					if name[0] == '_' {
						name = name[1:]
					}
				}
				sql += "DECLARE cnt NUMBER;\n"
				sql += "BEGIN\n"
				sql += fmt.Sprintf("SELECT count(*) INTO cnt FROM user_indexes WHERE index_name = '%s';\n", strings.ToUpper(name))
				sql += "IF cnt = 0 THEN\n"
				sql += "EXECUTE IMMEDIATE '\n"
				sql += fmt.Sprintf("CREATE INDEX %s%s%s ON %s%s%s (%s%s%s)\n", Q, name, Q, Q, mi.table, Q, Q, cols, Q)
				sql += "';\n"
				sql += "END IF;\n"
				sql += "END;"
			} else {
				sql = fmt.Sprintf("CREATE INDEX %s%s%s ON %s%s%s (%s%s%s);", Q, name, Q, Q, mi.table, Q, Q, cols, Q)
			}

			index := dbIndex{}
			index.Table = mi.table
			index.Name = name
			index.Sql = sql

			tableIndexes[mi.table] = append(tableIndexes[mi.table], index)
		}

	}

	return
}

// Get string value for the attribute "DEFAULT" for the CREATE, ALTER commands
func getColumnDefault(driver DriverType, fi *fieldInfo) string {
	var (
		v, t, d string
	)

	// Skip default attribute if field is in relations
	if fi.rel || fi.reverse {
		return v
	}

	if driver == DR_Oracle {
		t = " DEFAULT '''||'%s'||''' "
	} else {
		t = " DEFAULT '%s' "
	}

	// These defaults will be useful if there no config value orm:"default" and NOT NULL is on
	switch fi.fieldType {
	case TypeDateField, TypeDateTimeField:
		return v

	case TypeBooleanField:
		if driver == DR_Oracle {
			d = "F"
		} else {
			d = "0"
		}

	case TypeBitField, TypeSmallIntegerField, TypeIntegerField,
		TypeBigIntegerField, TypePositiveBitField, TypePositiveSmallIntegerField,
		TypePositiveIntegerField, TypePositiveBigIntegerField, TypeFloatField,
		TypeDecimalField:
		d = "0"
	}

	if fi.colDefault {
		if !fi.initial.Exist() {
			v = fmt.Sprintf(t, "")
		} else {
			if fi.fieldType == TypeBooleanField {
				if driver == DR_Oracle {
					switch strings.ToLower(fi.initial.String()) {
					case "true", "'true'":
						fi.initial.Set("T")
					case "false", "'false'":
						fi.initial.Set("F")
					}
				} else {
					switch strings.ToLower(fi.initial.String()) {
					case "true", "'true'":
						fi.initial.Set("1")
					case "false", "'false'":
						fi.initial.Set("0")
					}
				}
			}

			if strings.Contains(fi.initial.String(), "'") {
				fi.initial.Set(strings.Replace(fi.initial.String(), "'", "", -1))
			}

			v = fmt.Sprintf(t, fi.initial.String())
		}
	} else {
		if !fi.null {
			v = fmt.Sprintf(t, d)
		}
	}

	return v
}
