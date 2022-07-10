package query

import (
	"fmt"
	"strings"
)

type FilterGroup struct {
	groups        []*FilterGroup
	filters       []*Filter
	logicOperator LogicOperator
}

func NewFilterGroup(groups []*FilterGroup, logicOperator LogicOperator) *FilterGroup {
	return &FilterGroup{groups: groups, logicOperator: logicOperator}
}

func NewFilterGroupWithFilters(filters []*Filter, logicOperator LogicOperator) *FilterGroup {
	return &FilterGroup{filters: filters, logicOperator: logicOperator}
}

func (fg *FilterGroup) NumValue() (num int) {
	for _, group := range fg.groups {
		num += group.NumValue()
	}
	for _, filter := range fg.filters {
		num += filter.NumValue()
	}
	return
}

func (fg *FilterGroup) FillValue(values []any) error {
	numValue := fg.NumValue()
	if numValue != len(values) {
		return fmt.Errorf("expected %d values, but actual %d values", numValue, len(values))
	}

	fg.fillValue(values)
	return nil
}

func (fg *FilterGroup) IsEmpty() bool {
	return len(fg.groups) == 0 && len(fg.filters) == 0
}

func (fg *FilterGroup) fillValue(values []any) {
	remainingValues := values
	if fg.groups != nil {
		for _, group := range fg.groups {
			numValue := group.NumValue()
			group.fillValue(remainingValues[:numValue])
			if numValue >= len(remainingValues) {
				continue
			}
			remainingValues = remainingValues[numValue:]
		}
	}
	if fg.filters != nil {
		for _, filter := range fg.filters {
			numValue := filter.NumValue()
			filter.fillValue(remainingValues[:numValue])
			if numValue >= len(remainingValues) {
				continue
			}
			remainingValues = remainingValues[numValue:]
		}
	}
}

func (fg FilterGroup) String() string {

	builder := strings.Builder{}
	builder.Grow(256)
	isMultiple := len(fg.groups) > 1 || len(fg.filters) > 1
	if isMultiple {
		builder.WriteRune('(')
	}
	for i, group := range fg.groups {
		if i > 0 {
			builder.WriteRune(' ')
			builder.WriteString(group.logicOperator.String())
			builder.WriteRune(' ')
		}
		builder.WriteString(group.String())
	}
	for i, filter := range fg.filters {
		if i > 0 {
			builder.WriteRune(' ')
			builder.WriteString(fg.logicOperator.String())
			builder.WriteRune(' ')
		}
		builder.WriteRune('(')
		builder.WriteString(filter.String())
		builder.WriteRune(')')
	}
	if isMultiple {
		builder.WriteRune(')')
	}
	return builder.String()
}

type Filter struct {
	fieldName string
	predicate *Predicate
	modifier  *FilterModifier
	value     any
}

func NewFilter(fieldName string, predicate *Predicate, opts ...FilterOption) *Filter {
	filter := &Filter{fieldName: fieldName, predicate: predicate}
	for _, opt := range opts {
		opt(filter)
	}
	return filter
}

func (f *Filter) NumValue() (num int) {
	return f.predicate.numArgs
}

func (f *Filter) FillValue(values []any) error {
	numValue := f.NumValue()
	if numValue != len(values) {
		return fmt.Errorf("expected %d values, but actual %d values", numValue, len(values))
	}

	f.fillValue(values)
	return nil
}

func (f *Filter) fillValue(values []any) {
	switch f.predicate.numArgs {
	case 0:
	case 1:
		f.value = values[0]
	default:
		f.value = values
	}
}

func (f *Filter) FieldName() string {
	return f.fieldName
}

func (f *Filter) Predicate() *Predicate {
	return f.predicate
}

func (f *Filter) FilterModifier() *FilterModifier {
	return f.modifier
}

func (f Filter) String() string {
	builder := strings.Builder{}
	builder.Grow(256)
	builder.WriteString(f.fieldName)
	builder.WriteRune(' ')
	builder.WriteString(f.predicate.String())
	if f.modifier != nil {
		builder.WriteRune('(')
		builder.WriteString(f.modifier.String())
		builder.WriteRune(')')
	}
	builder.WriteRune(' ')
	builder.WriteString(fmt.Sprintf("%#v", f.value))
	return builder.String()
}

type FilterOption func(filter *Filter)

func WithFilterValue(value any) FilterOption {
	return func(filter *Filter) {
		filter.value = value
	}
}
func WithFilterModifier(modifier *FilterModifier) FilterOption {
	return func(filter *Filter) {
		filter.modifier = modifier
	}
}

type Predicate struct {
	keywords []string
	numArgs  int
}

func (p Predicate) String() string {
	return p.keywords[0]
}

func (p *Predicate) Keywords() []string {
	return p.keywords
}

func (p *Predicate) NumArgs() int {
	return p.numArgs
}

var (
	PredicateContains   = &Predicate{keywords: []string{"Contains"}, numArgs: 1}
	PredicateStartsWith = &Predicate{keywords: []string{"StartsWith"}, numArgs: 1}
	PredicateEndsWith   = &Predicate{keywords: []string{"EndsWith"}, numArgs: 1}
	PredicateIsNull     = &Predicate{keywords: []string{"IsNull"}, numArgs: 0}
	PredicateIsNotNull  = &Predicate{keywords: []string{"IsNotNull"}, numArgs: 0}
	//PredicateIsEmpty is collection empty,
	//If you want to express that the string is not empty, you can use 'Is',
	//and then use an empty string as a parameter, the same below
	PredicateIsEmpty = &Predicate{keywords: []string{"IsEmpty"}, numArgs: 0}
	//PredicateIsNotEmpty is collection not empty
	PredicateIsNotEmpty = &Predicate{keywords: []string{"IsNotEmpty"}, numArgs: 0}
	PredicateIsFalse    = &Predicate{keywords: []string{"IsFalse"}, numArgs: 0}
	PredicateIsTrue     = &Predicate{keywords: []string{"IsTrue"}, numArgs: 0}
	PredicateMatches    = &Predicate{keywords: []string{"Matches"}, numArgs: 1}
	//PredicateBetween The BETWEEN operator is inclusive: begin and end values are included.
	PredicateBetween = &Predicate{keywords: []string{"Between"}, numArgs: 2}
	PredicateNotIn   = &Predicate{keywords: []string{"NotIn"}, numArgs: 1}
	PredicateIn      = &Predicate{keywords: []string{"In"}, numArgs: 1}
	PredicateGT      = &Predicate{keywords: []string{"GT"}, numArgs: 1}
	PredicateLT      = &Predicate{keywords: []string{"LT"}, numArgs: 1}
	PredicateGTE     = &Predicate{keywords: []string{"GTE"}, numArgs: 1}
	PredicateLTE     = &Predicate{keywords: []string{"LTE"}, numArgs: 1}
	PredicateIsNot   = &Predicate{keywords: []string{"IsNot", "NotEquals", "NE"}, numArgs: 1}
	PredicateIs      = &Predicate{keywords: []string{"Equals", "Is", "EQ", ""}, numArgs: 1}
)

var Predicates = []*Predicate{
	PredicateContains, PredicateStartsWith, PredicateEndsWith, PredicateIsNull, PredicateIsNotNull,
	PredicateIsEmpty, PredicateIsNotEmpty, PredicateIsFalse, PredicateIsTrue, PredicateMatches,
	PredicateBetween, PredicateNotIn, PredicateIn, PredicateGT, PredicateLT,
	PredicateGTE, PredicateLTE, PredicateIsNot, PredicateIs,
}

type FilterModifier struct {
	keywords []string
	global   bool //wait support
}

func (p FilterModifier) String() string {
	return p.keywords[0]
}

func (p *FilterModifier) Keywords() []string {
	return p.keywords
}

func (p *FilterModifier) Global() bool {
	return p.global
}

var (
	FilterModifierIgnoreCase    = &FilterModifier{keywords: []string{"IgnoreCase", "IC"}, global: false}
	FilterModifierAllIgnoreCase = &FilterModifier{keywords: []string{"AllIgnoreCase", "AllIC"}, global: true}
)

var FilterModifiers = []*FilterModifier{
	FilterModifierIgnoreCase, FilterModifierAllIgnoreCase,
}

type LogicOperator string

func (l LogicOperator) String() string {
	return string(l)
}

const (
	LogicOperatorAnd LogicOperator = "And"
	LogicOperatorOr  LogicOperator = "Or"
)
