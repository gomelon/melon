package engine

import (
	"fmt"
	"github.com/huandu/xstrings"
	"strings"
)

func UseMySQL() {
	mySQL := NewMySQL()
	Engines[mySQL.Dialect()] = mySQL
}

type MySQL struct {
}

func NewMySQL() *MySQL {
	return &MySQL{}
}

func (m *MySQL) Dialect() string {
	return "mysql"
}

func (m *MySQL) Escape(str string) string {
	if strings.HasPrefix(str, "`") {
		return str
	}
	return fmt.Sprintf("`%s`", str)
}

func (m *MySQL) BuildColumn(str string) string {
	return xstrings.ToSnakeCase(str)
}

func (m *MySQL) BuildContains(str string) string {
	return fmt.Sprintf("LIKE CONCAT('%%',%s,'%%')", str)
}

func (m *MySQL) BuildStartsWith(str string) string {
	return fmt.Sprintf("LIKE CONCAT(%s,'%%')", str)
}

func (m *MySQL) BuildEndsWith(str string) string {
	return fmt.Sprintf("LIKE CONCAT('%%',%s)", str)
}

func (m *MySQL) BuildLimit(offset, limit string) string {
	return fmt.Sprintf("LIMIT %s, %s", offset, limit)
}
