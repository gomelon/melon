package engine

var Engines map[string]Engine = map[string]Engine{}

type Engine interface {
	Dialect() string
	Escape(str string) string
	BuildColumn(str string) string
	BuildContains(str string) string
	BuildStartsWith(str string) string
	BuildEndsWith(str string) string
	BuildLimit(offset, limit string) string
}
