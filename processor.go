package mongodb

import (
	"github.com/goravel/framework/contracts/database/driver"
)

var _ driver.Processor = &Processor{}

// Processor is a minimal implementation to satisfy Goravel's driver interface
// MongoDB operations don't use SQL-based column/index processing
type Processor struct {
}

func NewProcessor() *Processor {
	return &Processor{}
}

func (r Processor) ProcessColumns(dbColumns []driver.DBColumn) []driver.Column {
	// MongoDB doesn't use SQL-based column definitions
	return nil
}

func (r Processor) ProcessForeignKeys(dbForeignKeys []driver.DBForeignKey) []driver.ForeignKey {
	// MongoDB doesn't use foreign keys
	return nil
}

func (r Processor) ProcessIndexes(dbIndexes []driver.DBIndex) []driver.Index {
	// MongoDB index management is handled through native MongoDB operations
	return nil
}

func (r Processor) ProcessTypes(types []driver.Type) []driver.Type {
	// MongoDB doesn't use SQL types
	return types
}
