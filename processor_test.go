package mongodb

import (
	"testing"

	"github.com/goravel/framework/contracts/database/driver"
	"github.com/stretchr/testify/suite"
)

type ProcessorTestSuite struct {
	suite.Suite
	processor *Processor
}

func TestProcessorTestSuite(t *testing.T) {
	suite.Run(t, new(ProcessorTestSuite))
}

func (s *ProcessorTestSuite) SetupTest() {
	s.processor = NewProcessor()
}

func (s *ProcessorTestSuite) TestProcessColumns() {
	// MongoDB doesn't use SQL-based column definitions
	dbColumns := []driver.DBColumn{
		{Name: "id", Type: "string"},
		{Name: "name", Type: "string"},
	}

	result := s.processor.ProcessColumns(dbColumns)
	s.Nil(result)
}

func (s *ProcessorTestSuite) TestProcessForeignKeys() {
	// MongoDB doesn't use foreign keys
	dbForeignKeys := []driver.DBForeignKey{
		{Name: "fk_user_role"},
	}

	result := s.processor.ProcessForeignKeys(dbForeignKeys)
	s.Nil(result)
}

func (s *ProcessorTestSuite) TestProcessIndexes() {
	// MongoDB index management is handled through native MongoDB operations
	dbIndexes := []driver.DBIndex{
		{Name: "idx_name"},
	}

	result := s.processor.ProcessIndexes(dbIndexes)
	s.Nil(result)
}

func (s *ProcessorTestSuite) TestProcessTypes() {
	// MongoDB doesn't use SQL types, but this method returns input as-is
	types := []driver.Type{
		{Name: "string"},
		{Name: "integer"},
	}

	result := s.processor.ProcessTypes(types)
	s.Equal(types, result)
}
