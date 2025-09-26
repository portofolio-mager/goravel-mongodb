// MongoDB dialector for GORM integration

package mongodb

import (
	"context"
	"database/sql"
	"strconv"

	"gorm.io/gorm"
	"gorm.io/gorm/callbacks"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type Dialector struct {
	DSN  string
	Conn gorm.ConnPool
}

// dummyConnPool is a minimal implementation of gorm.ConnPool for MongoDB
// MongoDB operations will use the native MongoDB driver directly
type dummyConnPool struct{}

func (d *dummyConnPool) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	return nil, gorm.ErrNotImplemented
}

func (d *dummyConnPool) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return nil, gorm.ErrNotImplemented
}

func (d *dummyConnPool) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return nil, gorm.ErrNotImplemented
}

func (d *dummyConnPool) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return nil
}

// Open opens a GORM dialector from a data source name.
func Open(dsn string) gorm.Dialector {
	return &Dialector{DSN: dsn}
}

// Open opens a GORM dialector from a database handle.
func OpenDB(db gorm.ConnPool) gorm.Dialector {
	return &Dialector{Conn: db}
}

func (dialector Dialector) Name() string {
	return "mongodb"
}

func (dialector Dialector) Initialize(db *gorm.DB) (err error) {
	if dialector.Conn != nil {
		db.ConnPool = dialector.Conn
	} else {
		// For MongoDB, we don't use SQL connection pools like other GORM dialectors
		// We'll handle connections through the native MongoDB driver
		// Set a dummy connection to satisfy GORM requirements
		db.ConnPool = &dummyConnPool{}
	}

	callbacks.RegisterDefaultCallbacks(db, &callbacks.Config{
		CreateClauses:        []string{"INSERT", "VALUES", "ON CONFLICT", "RETURNING"},
		UpdateClauses:        []string{"UPDATE", "SET", "FROM", "WHERE", "RETURNING"},
		DeleteClauses:        []string{"DELETE", "FROM", "WHERE", "RETURNING"},
		LastInsertIDReversed: true,
	})

	for k, v := range dialector.ClauseBuilders() {
		db.ClauseBuilders[k] = v
	}
	return
}

func (dialector Dialector) ClauseBuilders() map[string]clause.ClauseBuilder {
	return map[string]clause.ClauseBuilder{
		"INSERT": func(c clause.Clause, builder clause.Builder) {
			if insert, ok := c.Expression.(clause.Insert); ok {
				if stmt, ok := builder.(*gorm.Statement); ok {
					_, _ = stmt.WriteString("INSERT ")
					if insert.Modifier != "" {
						_, _ = stmt.WriteString(insert.Modifier)
						_ = stmt.WriteByte(' ')
					}

					_, _ = stmt.WriteString("INTO ")
					if insert.Table.Name == "" {
						stmt.WriteQuoted(stmt.Table)
					} else {
						stmt.WriteQuoted(insert.Table)
					}
					return
				}
			}

			c.Build(builder)
		},
		"LIMIT": func(c clause.Clause, builder clause.Builder) {
			if limit, ok := c.Expression.(clause.Limit); ok {
				var lmt = -1
				if limit.Limit != nil && *limit.Limit >= 0 {
					lmt = *limit.Limit
				}
				if lmt >= 0 || limit.Offset > 0 {
					_, _ = builder.WriteString("LIMIT ")
					_, _ = builder.WriteString(strconv.Itoa(lmt))
				}
				if limit.Offset > 0 {
					_, _ = builder.WriteString(" OFFSET ")
					_, _ = builder.WriteString(strconv.Itoa(limit.Offset))
				}
			}
		},
		"FOR": func(c clause.Clause, builder clause.Builder) {
			if _, ok := c.Expression.(clause.Locking); ok {
				// SQLite3 does not support row-level locking.
				return
			}
			c.Build(builder)
		},
	}
}

func (dialector Dialector) DefaultValueOf(field *schema.Field) clause.Expression {
	if field.AutoIncrement {
		return clause.Expr{SQL: "NULL"}
	}

	// doesn't work, will raise error
	return clause.Expr{SQL: "DEFAULT"}
}

func (dialector Dialector) Migrator(db *gorm.DB) gorm.Migrator {
	return nil
}

func (dialector Dialector) BindVarTo(writer clause.Writer, stmt *gorm.Statement, v interface{}) {
	_ = writer.WriteByte('?')
}

func (dialector Dialector) QuoteTo(writer clause.Writer, str string) {
	var (
		underQuoted, selfQuoted bool
		continuousBacktick      int8
		shiftDelimiter          int8
	)

	for _, v := range []byte(str) {
		switch v {
		case '`':
			continuousBacktick++
			if continuousBacktick == 2 {
				_, _ = writer.WriteString("``")
				continuousBacktick = 0
			}
		case '.':
			if continuousBacktick > 0 || !selfQuoted {
				shiftDelimiter = 0
				underQuoted = false
				continuousBacktick = 0
				_, _ = writer.WriteString("`")
			}
			_ = writer.WriteByte(v)
			continue
		default:
			if shiftDelimiter-continuousBacktick <= 0 && !underQuoted {
				_, _ = writer.WriteString("`")
				underQuoted = true
				if selfQuoted = continuousBacktick > 0; selfQuoted {
					continuousBacktick -= 1
				}
			}

			for ; continuousBacktick > 0; continuousBacktick -= 1 {
				_, _ = writer.WriteString("``")
			}

			_ = writer.WriteByte(v)
		}
		shiftDelimiter++
	}

	if continuousBacktick > 0 && !selfQuoted {
		_, _ = writer.WriteString("``")
	}
	_, _ = writer.WriteString("`")
}

func (dialector Dialector) Explain(sql string, vars ...interface{}) string {
	return logger.ExplainSQL(sql, nil, `'`, vars...)
}

func (dialector Dialector) DataTypeOf(field *schema.Field) string {
	switch field.DataType {
	case schema.Bool:
		return "numeric"
	case schema.Int, schema.Uint:
		if field.AutoIncrement {
			// doesn't check `PrimaryKey`, to keep backward compatibility
			// https://sqlite.org/autoinc.html
			return "integer PRIMARY KEY AUTOINCREMENT"
		} else {
			return "integer"
		}
	case schema.Float:
		return "real"
	case schema.String:
		return "text"
	case schema.Time:
		// Distinguish between schema.Time and tag time
		if val, ok := field.TagSettings["TYPE"]; ok {
			return val
		} else {
			return "datetime"
		}
	case schema.Bytes:
		return "blob"
	}

	return string(field.DataType)
}

func (dialectopr Dialector) SavePoint(tx *gorm.DB, name string) error {
	tx.Exec("SAVEPOINT " + name)
	return nil
}

func (dialectopr Dialector) RollbackTo(tx *gorm.DB, name string) error {
	tx.Exec("ROLLBACK TO SAVEPOINT " + name)
	return nil
}
