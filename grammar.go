package mongodb

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/goravel/framework/contracts/database/driver"
	"github.com/goravel/framework/contracts/log"
	"gorm.io/gorm/clause"
)

var _ driver.Grammar = &Grammar{}

// Grammar is a minimal implementation to satisfy Goravel's driver interface
// MongoDB operations bypass this grammar and use native MongoDB queries
type Grammar struct {
	log log.Log
}

func NewGrammar(log log.Log) *Grammar {
	return &Grammar{
		log: log,
	}
}

// Minimal implementations - MongoDB operations don't use SQL grammar
func (g *Grammar) CompileAdd(blueprint driver.Blueprint, command *driver.Command) string {
	return ""
}

func (g *Grammar) CompileChange(blueprint driver.Blueprint, command *driver.Command) []string {
	return nil
}

func (g *Grammar) CompileColumns(schema, table string) (string, error) {
	return "", nil
}

func (g *Grammar) CompileDefault(blueprint driver.Blueprint, command *driver.Command) string {
	return ""
}

func (g *Grammar) CompileComment(blueprint driver.Blueprint, command *driver.Command) string {
	return ""
}

func (g *Grammar) CompileCreate(blueprint driver.Blueprint) string {
	return ""
}

func (g *Grammar) CompileDisableWriteableSchema() string {
	return ""
}

func (g *Grammar) CompileDrop(blueprint driver.Blueprint) string {
	return ""
}

func (g *Grammar) CompileDropAllDomains(domains []string) string {
	return ""
}

func (g *Grammar) CompileDropAllTables(database string, tables []driver.Table) []string {
	return nil
}

func (g *Grammar) CompileDropAllTypes(database string, types []driver.Type) []string {
	return nil
}

func (g *Grammar) CompileDropAllViews(database string, views []driver.View) []string {
	return nil
}

func (g *Grammar) CompileDropColumn(blueprint driver.Blueprint, command *driver.Command) []string {
	return nil
}

func (g *Grammar) CompileDropForeign(blueprint driver.Blueprint, command *driver.Command) string {
	return ""
}

func (g *Grammar) CompileDropFullText(blueprint driver.Blueprint, command *driver.Command) string {
	return ""
}

func (g *Grammar) CompileDropIfExists(blueprint driver.Blueprint) string {
	return ""
}

func (g *Grammar) CompileDropIndex(blueprint driver.Blueprint, command *driver.Command) string {
	return ""
}

func (g *Grammar) CompileDropPrimary(blueprint driver.Blueprint, command *driver.Command) string {
	return ""
}

func (g *Grammar) CompileDropUnique(blueprint driver.Blueprint, command *driver.Command) string {
	return ""
}

func (g *Grammar) CompileEnableWriteableSchema() string {
	return ""
}

func (g *Grammar) CompileForeign(blueprint driver.Blueprint, command *driver.Command) string {
	return ""
}

func (g *Grammar) CompileForeignKeys(database, table string) string {
	return ""
}

func (g *Grammar) CompileFullText(blueprint driver.Blueprint, command *driver.Command) string {
	return ""
}

func (g *Grammar) CompileIndex(blueprint driver.Blueprint, command *driver.Command) string {
	return ""
}

func (g *Grammar) CompileIndexes(database, table string) (string, error) {
	return "", nil
}

func (g *Grammar) CompileJsonColumnsUpdate(values map[string]any) (map[string]any, error) {
	return values, nil
}

func (g *Grammar) CompileJsonContains(column string, value any, isNot bool) (string, []any, error) {
	return "", nil, nil
}

func (g *Grammar) CompileJsonContainsKey(column string, isNot bool) string {
	return ""
}

func (g *Grammar) CompileJsonLength(column string) string {
	return ""
}

func (g *Grammar) CompileJsonSelector(column string) string {
	return ""
}

func (g *Grammar) CompileJsonValues(args ...any) []any {
	return args
}

func (g *Grammar) CompileLockForUpdate(builder sq.SelectBuilder, conditions *driver.Conditions) sq.SelectBuilder {
	return builder
}

func (g *Grammar) CompileLockForUpdateForGorm() clause.Expression {
	return nil
}

func (g *Grammar) CompilePlaceholderFormat() driver.PlaceholderFormat {
	return nil
}

func (g *Grammar) CompilePrimary(blueprint driver.Blueprint, command *driver.Command) string {
	return ""
}

func (g *Grammar) CompileInRandomOrder(builder sq.SelectBuilder, conditions *driver.Conditions) sq.SelectBuilder {
	return builder
}

func (g *Grammar) CompilePrune(database string) string {
	return ""
}

func (g *Grammar) CompileRandomOrderForGorm() string {
	return ""
}

func (g *Grammar) CompileRename(blueprint driver.Blueprint, command *driver.Command) string {
	return ""
}

func (g *Grammar) CompileRenameColumn(blueprint driver.Blueprint, command *driver.Command, columns []driver.Column) (string, error) {
	return "", nil
}

func (g *Grammar) CompileRenameIndex(blueprint driver.Blueprint, command *driver.Command, indexes []driver.Index) []string {
	return nil
}

func (g *Grammar) CompileSharedLock(builder sq.SelectBuilder, conditions *driver.Conditions) sq.SelectBuilder {
	return builder
}

func (g *Grammar) CompileSharedLockForGorm() clause.Expression {
	return nil
}

func (g *Grammar) CompileTables(database string) string {
	return ""
}

func (g *Grammar) CompileTableComment(blueprint driver.Blueprint, command *driver.Command) string {
	return ""
}

func (g *Grammar) CompileTypes() string {
	return ""
}

func (g *Grammar) CompileUnique(blueprint driver.Blueprint, command *driver.Command) string {
	return ""
}

func (g *Grammar) CompileVersion() string {
	return ""
}

func (g *Grammar) CompileViews(database string) string {
	return ""
}

func (g *Grammar) GetAttributeCommands() []string {
	return nil
}

func (g *Grammar) GetModifiers() []func(driver.Blueprint, driver.ColumnDefinition) string {
	return nil
}

// Type methods - not used in MongoDB but required for interface
func (g *Grammar) TypeBigInteger(column driver.ColumnDefinition) string   { return "" }
func (g *Grammar) TypeBoolean(column driver.ColumnDefinition) string     { return "" }
func (g *Grammar) TypeChar(column driver.ColumnDefinition) string        { return "" }
func (g *Grammar) TypeDate(column driver.ColumnDefinition) string        { return "" }
func (g *Grammar) TypeDateTime(column driver.ColumnDefinition) string    { return "" }
func (g *Grammar) TypeDateTimeTz(column driver.ColumnDefinition) string  { return "" }
func (g *Grammar) TypeDecimal(column driver.ColumnDefinition) string     { return "" }
func (g *Grammar) TypeDouble(column driver.ColumnDefinition) string      { return "" }
func (g *Grammar) TypeEnum(column driver.ColumnDefinition) string        { return "" }
func (g *Grammar) TypeFloat(column driver.ColumnDefinition) string       { return "" }
func (g *Grammar) TypeInteger(column driver.ColumnDefinition) string     { return "" }
func (g *Grammar) TypeJson(column driver.ColumnDefinition) string        { return "" }
func (g *Grammar) TypeJsonb(column driver.ColumnDefinition) string       { return "" }
func (g *Grammar) TypeLongText(column driver.ColumnDefinition) string    { return "" }
func (g *Grammar) TypeMediumInteger(column driver.ColumnDefinition) string { return "" }
func (g *Grammar) TypeMediumText(column driver.ColumnDefinition) string  { return "" }
func (g *Grammar) TypeSmallInteger(column driver.ColumnDefinition) string { return "" }
func (g *Grammar) TypeString(column driver.ColumnDefinition) string      { return "" }
func (g *Grammar) TypeText(column driver.ColumnDefinition) string        { return "" }
func (g *Grammar) TypeTime(column driver.ColumnDefinition) string        { return "" }
func (g *Grammar) TypeTimeTz(column driver.ColumnDefinition) string      { return "" }
func (g *Grammar) TypeTimestamp(column driver.ColumnDefinition) string   { return "" }
func (g *Grammar) TypeTimestampTz(column driver.ColumnDefinition) string { return "" }
func (g *Grammar) TypeTinyInteger(column driver.ColumnDefinition) string { return "" }
func (g *Grammar) TypeTinyText(column driver.ColumnDefinition) string    { return "" }
func (g *Grammar) TypeUuid(column driver.ColumnDefinition) string        { return "" }