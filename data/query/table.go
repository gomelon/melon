package query

type Table interface {
	Schema() string
	Name() string
}

type simpleTable struct {
	schema string
	name   string
}

type TableOption func(table *simpleTable)

func WithTableSchema(schema string) TableOption {
	return func(table *simpleTable) {
		table.schema = schema
	}
}

func NewTable(name string, opts ...TableOption) Table {
	table := &simpleTable{name: name}
	for _, opt := range opts {
		opt(table)
	}
	return table
}

func (s *simpleTable) Schema() string {
	return s.schema
}

func (s *simpleTable) Name() string {
	return s.name
}
