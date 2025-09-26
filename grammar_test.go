package mongodb

import (
	"testing"

	contractsdriver "github.com/goravel/framework/contracts/database/driver"
	mocksdriver "github.com/goravel/framework/mocks/database/driver"
	mockslog "github.com/goravel/framework/mocks/log"
	"github.com/stretchr/testify/suite"
)

type GrammarSuite struct {
	suite.Suite
	grammar *Grammar
	mockLog *mockslog.Log
}

func TestGrammarSuite(t *testing.T) {
	suite.Run(t, &GrammarSuite{})
}

func (s *GrammarSuite) SetupTest() {
	s.mockLog = mockslog.NewLog(s.T())
	s.grammar = NewGrammar(s.mockLog)
}

func (s *GrammarSuite) TestCompileForeign() {
	// MongoDB doesn't have foreign keys like SQL databases
	mockBlueprint := mocksdriver.NewBlueprint(s.T())
	command := &contractsdriver.Command{}
	result := s.grammar.CompileForeign(mockBlueprint, command)
	s.Equal("", result)
}

func (s *GrammarSuite) TestCompileAdd() {
	// MongoDB doesn't use SQL ALTER TABLE commands
	mockBlueprint := mocksdriver.NewBlueprint(s.T())
	mockColumn := mocksdriver.NewColumnDefinition(s.T())

	result := s.grammar.CompileAdd(mockBlueprint, &contractsdriver.Command{
		Column: mockColumn,
	})

	s.Equal("", result)
}

func (s *GrammarSuite) TestCompileCreate() {
	// MongoDB doesn't use SQL CREATE TABLE commands
	// Collections are created automatically when documents are inserted
	mockBlueprint := mocksdriver.NewBlueprint(s.T())

	result := s.grammar.CompileCreate(mockBlueprint)
	s.Equal("", result)
}

func (s *GrammarSuite) TestCompileDropColumn() {
	// MongoDB doesn't use SQL ALTER TABLE commands
	mockBlueprint := mocksdriver.NewBlueprint(s.T())

	result := s.grammar.CompileDropColumn(mockBlueprint, &contractsdriver.Command{
		Name:    "name",
		Columns: []string{"id", "name"},
	})

	s.Nil(result)
}

func (s *GrammarSuite) TestCompileDropIfExists() {
	// MongoDB doesn't use SQL DROP TABLE commands
	mockBlueprint := mocksdriver.NewBlueprint(s.T())

	result := s.grammar.CompileDropIfExists(mockBlueprint)
	s.Equal("", result)
}

func (s *GrammarSuite) TestCompileIndex() {
	// MongoDB handles indexes differently than SQL
	mockBlueprint := mocksdriver.NewBlueprint(s.T())
	command := &contractsdriver.Command{
		Index:   "users",
		Columns: []string{"role_id", "permission_id"},
	}

	result := s.grammar.CompileIndex(mockBlueprint, command)
	s.Equal("", result)
}

func (s *GrammarSuite) TestCompileJsonColumnsUpdate() {
	// MongoDB handles JSON naturally, no special compilation needed
	values := map[string]any{
		"data": map[string]any{
			"name": "test",
		},
	}

	result, err := s.grammar.CompileJsonColumnsUpdate(values)
	s.NoError(err)
	s.Equal(values, result)
}

func (s *GrammarSuite) TestCompileJsonContains() {
	// MongoDB handles JSON queries differently
	result, values, err := s.grammar.CompileJsonContains("data->details", "value1", false)
	s.NoError(err)
	s.Equal("", result)
	s.Nil(values)
}

func (s *GrammarSuite) TestCompileJsonContainsKey() {
	// MongoDB handles JSON queries differently
	result := s.grammar.CompileJsonContainsKey("data->details", false)
	s.Equal("", result)
}

func (s *GrammarSuite) TestCompileJsonLength() {
	// MongoDB handles JSON queries differently
	result := s.grammar.CompileJsonLength("data->details")
	s.Equal("", result)
}

func (s *GrammarSuite) TestCompileRenameColumn() {
	// MongoDB doesn't rename columns in the SQL sense
	mockBlueprint := mocksdriver.NewBlueprint(s.T())
	mockColumn := mocksdriver.NewColumnDefinition(s.T())

	result, err := s.grammar.CompileRenameColumn(mockBlueprint, &contractsdriver.Command{
		Column: mockColumn,
		From:   "before",
		To:     "after",
	}, nil)

	s.NoError(err)
	s.Equal("", result)
}

func (s *GrammarSuite) TestCompileRenameIndex() {
	// MongoDB doesn't use SQL index commands
	mockBlueprint := mocksdriver.NewBlueprint(s.T())
	command := &contractsdriver.Command{
		From: "users",
		To:   "admins",
	}

	result := s.grammar.CompileRenameIndex(mockBlueprint, command, nil)
	s.Nil(result)
}

func (s *GrammarSuite) TestTypeBoolean() {
	mockColumn := mocksdriver.NewColumnDefinition(s.T())
	result := s.grammar.TypeBoolean(mockColumn)
	s.Equal("", result)
}

func (s *GrammarSuite) TestTypeEnum() {
	mockColumn := mocksdriver.NewColumnDefinition(s.T())
	result := s.grammar.TypeEnum(mockColumn)
	s.Equal("", result)
}

func (s *GrammarSuite) TestTypeUuid() {
	mockColumn := mocksdriver.NewColumnDefinition(s.T())
	result := s.grammar.TypeUuid(mockColumn)
	s.Equal("", result)
}