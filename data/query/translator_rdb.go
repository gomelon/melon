package query

import (
	"context"
	"fmt"
	"github.com/gomelon/melon/data/engine"
	"strings"
)

type RDBTranslator struct {
	engin engine.Engine
}

func NewRDBTranslator(engin engine.Engine) *RDBTranslator {
	return &RDBTranslator{engin: engin}
}

func (t *RDBTranslator) Translate(ctx context.Context, query *Query) (result string, err error) {
	switch query.Subject() {
	case SubjectFind:
		result, err = t.TranslateFind(ctx, query)
	case SubjectCount:
		result, err = t.TranslateCount(ctx, query)
	case SubjectExists:
		result, err = t.TranslateExists(ctx, query)
	case SubjectDelete:
		result, err = t.TranslateDelete(ctx, query)
	default:
		err = fmt.Errorf("translate query fail: unsupported subject [%s]", query.subject.String())
	}
	return
}

func (t *RDBTranslator) TranslateTable(ctx context.Context, table Table) (string, error) {
	if table == nil {
		return t.engin.Escape("$$_table_$$"), nil
	}
	if len(table.Schema()) == 0 {
		return t.engin.Escape(table.Name()), nil
	}
	return t.engin.Escape(table.Schema()) + "." + t.engin.Escape(table.Name()), nil
}

func (t *RDBTranslator) TranslateFind(ctx context.Context, query *Query) (result string, err error) {
	var subjectStr string
	switch query.subjectModifier {
	case SubjectModifierDistinct:
		subjectStr = "SELECT DISTINCT * FROM "
	default:
		subjectStr = "SELECT * FROM "
	}
	tableStr, err := t.TranslateTable(ctx, query.Table())
	if err != nil {
		return
	}

	whereStr, err := t.TranslateFilterGroup(ctx, query.FilterGroup())
	if err != nil {
		return
	}

	sortsStr, err := t.TranslateSorts(ctx, query.Sorts())
	if err != nil {
		return
	}

	pagerStr, err := t.TranslatePager(ctx, query.Pager())
	if err != nil {
		return
	}

	builder := &strings.Builder{}
	t.build(builder, subjectStr, tableStr, whereStr, sortsStr, pagerStr)
	result = builder.String()
	return
}

func (t *RDBTranslator) TranslateCount(ctx context.Context, query *Query) (result string, err error) {
	var subjectStr string
	switch query.subjectModifier {
	case SubjectModifierDistinct:
		subjectStr = "SELECT COUNT(DISTINCT *) FROM "
	default:
		subjectStr = "SELECT COUNT(*) FROM "
	}
	tableStr, err := t.TranslateTable(ctx, query.Table())
	if err != nil {
		return
	}

	whereStr, err := t.TranslateFilterGroup(ctx, query.FilterGroup())
	if err != nil {
		return
	}

	builder := &strings.Builder{}
	t.build(builder, subjectStr, tableStr, whereStr, "", "")
	result = builder.String()
	return
}

func (t *RDBTranslator) TranslateExists(ctx context.Context, query *Query) (result string, err error) {
	subjectStr := "SELECT 1 FROM "
	tableStr, err := t.TranslateTable(ctx, query.Table())
	if err != nil {
		return
	}

	whereStr, err := t.TranslateFilterGroup(ctx, query.FilterGroup())
	if err != nil {
		return
	}

	builder := &strings.Builder{}
	t.build(builder, subjectStr, tableStr, whereStr, "", t.engin.BuildLimit("0", "1"))
	result = builder.String()
	return
}

func (t *RDBTranslator) TranslateDelete(ctx context.Context, query *Query) (result string, err error) {
	subjectStr := "DELETE FROM "
	tableStr, err := t.TranslateTable(ctx, query.Table())
	if err != nil {
		return
	}

	whereStr, err := t.TranslateFilterGroup(ctx, query.FilterGroup())
	if err != nil {
		return
	}

	sortsStr, err := t.TranslateSorts(ctx, query.Sorts())
	if err != nil {
		return
	}

	pagerStr, err := t.TranslatePager(ctx, query.Pager())
	if err != nil {
		return
	}

	builder := &strings.Builder{}
	t.build(builder, subjectStr, tableStr, whereStr, sortsStr, pagerStr)
	result = builder.String()
	return
}

func (t *RDBTranslator) TranslateFilterGroup(ctx context.Context, fg *FilterGroup) (result string, err error) {
	operator, err := t.TranslateLogicOperator(ctx, fg.logicOperator)
	if err != nil {
		return
	}
	if fg == nil || fg.IsEmpty() {
		return
	}
	builder := strings.Builder{}
	builder.Grow(256)

	isMultiple := len(fg.groups) > 1 || len(fg.filters) > 1
	if isMultiple {
		builder.WriteRune('(')
	}

	var elementResult string
	for i, group := range fg.groups {
		if i > 0 {
			builder.WriteRune(' ')
			builder.WriteString(operator)
			builder.WriteRune(' ')
		}
		elementResult, err = t.TranslateFilterGroup(ctx, group)
		if err != nil {
			return
		}
		builder.WriteString(elementResult)
	}
	for i, filter := range fg.filters {
		if i > 0 {
			builder.WriteRune(' ')
			builder.WriteString(operator)
			builder.WriteRune(' ')
		}
		elementResult, err = t.TranslateFilter(ctx, filter)
		if err != nil {
			return
		}
		builder.WriteString(elementResult)
	}

	if isMultiple {
		builder.WriteRune(')')
	}

	result = builder.String()
	return
}

func (t *RDBTranslator) TranslateFilter(ctx context.Context, f *Filter) (result string, err error) {
	e := t.engin
	column := e.Escape(e.BuildColumn(f.FieldName()))
	switch f.Predicate() {
	case PredicateIs:
		result = fmt.Sprintf("(%s = %s)", column, t.namedArgOrValue(f, 0))
	case PredicateIsNot:
		result = fmt.Sprintf("(%s != %s)", column, t.namedArgOrValue(f, 0))
	case PredicateGT:
		result = fmt.Sprintf("(%s > %s)", column, t.namedArgOrValue(f, 0))
	case PredicateLT:
		result = fmt.Sprintf("(%s < %s)", column, t.namedArgOrValue(f, 0))
	case PredicateGTE:
		result = fmt.Sprintf("(%s >= %s)", column, t.namedArgOrValue(f, 0))
	case PredicateLTE:
		result = fmt.Sprintf("(%s <= %s)", column, t.namedArgOrValue(f, 0))
	case PredicateBetween:
		result = fmt.Sprintf("(%s >= %s AND %s <= %s)",
			column, column, t.namedArgOrValue(f, 0), t.namedArgOrValue(f, 1))
	case PredicateIn:
		result = fmt.Sprintf("(%s in (%s))", column, t.namedArgOrValue(f, 0))
	case PredicateNotIn:
		result = fmt.Sprintf("(%s NOT IN (%s))", column, t.namedArgOrValue(f, 0))
	case PredicateContains:
		result = fmt.Sprintf("(%s %s)", column, e.BuildContains(t.namedArgOrValue(f, 0)))
	case PredicateStartsWith:
		result = fmt.Sprintf("(%s %s)", column, e.BuildStartsWith(t.namedArgOrValue(f, 0)))
	case PredicateEndsWith:
		result = fmt.Sprintf("(%s %s)", column, e.BuildEndsWith(t.namedArgOrValue(f, 0)))
	case PredicateIsNull:
		result = fmt.Sprintf("(%s IS NULL)", column)
	case PredicateIsNotNull:
		result = fmt.Sprintf("(%s IS NOT NULL)", column)
	case PredicateIsEmpty:
		//TODO wait implement
		panic("wait implement")
	case PredicateIsNotEmpty:
		//TODO wait implement
		panic("wait implement")
	case PredicateIsFalse:
		result = fmt.Sprintf("(!%s)", column)
	case PredicateIsTrue:
		result = fmt.Sprintf("(%s)", column)
	case PredicateMatches:
		//TODO wait implement
		panic("wait implement")
	default:
		err = fmt.Errorf("translate query fail: unsupoorted predicate [%s]", f.Predicate().String())
	}
	return
}

func (t *RDBTranslator) TranslateLogicOperator(ctx context.Context, operator LogicOperator) (result string, err error) {
	switch operator {
	case LogicOperatorAnd:
		result = "AND"
	case LogicOperatorOr:
		result = "OR"
	default:
		err = fmt.Errorf("translate query fail: unsupoorted logic operator [%s]", operator.String())
	}
	return
}

func (t *RDBTranslator) TranslateSorts(ctx context.Context, sorts []*Sort) (result string, err error) {
	if len(sorts) == 0 {
		return
	}
	builder := strings.Builder{}
	builder.Grow(64)
	builder.WriteString("ORDER BY ")
	var sortStr string
	for i, sort := range sorts {
		if i > 0 {
			builder.WriteString(", ")
		}
		sortStr, err = t.TranslateSort(ctx, sort)
		if err != nil {
			return
		}
		builder.WriteString(sortStr)
	}
	result = builder.String()
	return
}

func (t *RDBTranslator) TranslateSort(ctx context.Context, sort *Sort) (result string, err error) {
	column := t.engin.Escape(t.engin.BuildColumn(sort.FieldName()))
	return fmt.Sprintf("%s %s", column, strings.ToUpper(string(sort.Direction()))), nil
}

func (t *RDBTranslator) TranslatePager(ctx context.Context, pager Pager) (result string, err error) {
	if pager == nil {
		return
	}
	return t.engin.BuildLimit("?", "?"), nil
}

func (t *RDBTranslator) build(builder *strings.Builder, subjectStr string, tableStr string,
	whereStr string, sortsStr string, pagerStr string) {

	builder.Grow(len(subjectStr) + len(tableStr) + len(whereStr) + len(sortsStr) + len(pagerStr) + 8)
	builder.WriteString(subjectStr)
	builder.WriteString(tableStr)
	if len(whereStr) > 0 {
		builder.WriteString(" WHERE ")
		builder.WriteString(whereStr)
	}
	if len(sortsStr) > 0 {
		builder.WriteRune(' ')
		builder.WriteString(sortsStr)
	}
	if len(pagerStr) > 0 {
		builder.WriteRune(' ')
		builder.WriteString(pagerStr)
	}
}

func (t *RDBTranslator) namedArgOrValue(f *Filter, index int) string {
	if f.NamedArgs() != nil {
		return ":" + f.NamedArgs()[index]
	}
	return "?"
}
