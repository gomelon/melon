package query

import (
	"fmt"
	"strings"
)

type Query struct {
	//TODO 要添加projection
	table               Table
	subject             *Subject
	subjectModifier     *SubjectModifier
	subjectModifierArgs map[SubjectModifierArg]any
	filterGroup         *FilterGroup
	sorts               []*Sort
	pager               Pager
}

func New(subject *Subject, opts ...Option) *Query {
	q := &Query{
		subject: subject,
	}
	for _, opt := range opts {
		opt(q)
	}
	return q
}

func (q *Query) With(opts ...Option) *Query {
	newQuery := &Query{
		table:               q.table,
		subject:             q.subject,
		subjectModifier:     q.subjectModifier,
		subjectModifierArgs: q.subjectModifierArgs,
		filterGroup:         q.filterGroup,
		sorts:               q.sorts,
		pager:               q.pager,
	}
	for _, opt := range opts {
		opt(newQuery)
	}
	return newQuery
}

func (q *Query) Table() Table {
	return q.table
}

func (q *Query) Subject() *Subject {
	return q.subject
}

func (q *Query) SubjectModifier() *SubjectModifier {
	return q.subjectModifier
}

func (q *Query) SubjectModifierArgs() map[SubjectModifierArg]any {
	return q.subjectModifierArgs
}

func (q *Query) FilterGroup() *FilterGroup {
	return q.filterGroup
}

func (q *Query) Sorts() []*Sort {
	return q.sorts
}

func (q *Query) Pager() Pager {
	return q.pager
}

type Option func(q *Query)

func WithTable(table Table) Option {
	return func(q *Query) {
		q.table = table
	}
}

func WithSubjectModifier(modifier *SubjectModifier) Option {
	return func(q *Query) {
		q.subjectModifier = modifier
	}
}

func WithSubjectModifierArgs(subjectModifierArgs map[SubjectModifierArg]any) Option {
	return func(q *Query) {
		q.subjectModifierArgs = subjectModifierArgs
	}
}

func WithFilterGroup(filterGroup *FilterGroup) Option {
	return func(q *Query) {
		q.filterGroup = filterGroup
	}
}

func WithSorts(sorts []*Sort) Option {
	return func(q *Query) {
		q.sorts = sorts
	}
}

func WithPager(pager Pager) Option {
	return func(q *Query) {
		q.pager = pager
	}
}

func (q Query) String() string {
	builder := strings.Builder{}
	builder.Grow(256)
	builder.WriteString(q.subject.String())

	if q.subjectModifier != nil {
		builder.WriteRune(' ')
		builder.WriteString(q.subjectModifier.String())
	}

	if len(q.subjectModifierArgs) > 0 {
		builder.WriteRune(' ')
		builder.WriteString(fmt.Sprintf("%v", q.subjectModifierArgs))
	}

	if q.filterGroup != nil {
		builder.WriteString(" WHERE ")
		builder.WriteString(q.filterGroup.String())
	}

	if len(q.sorts) > 0 {
		builder.WriteString(" Order By ")
		for i, sort := range q.sorts {
			if i > 0 {
				builder.WriteString(", ")
			}
			builder.WriteString(sort.String())
		}
	}

	if q.pager != nil {
		builder.WriteRune(' ')
		builder.WriteString(q.pager.String())
	}
	return builder.String()
}
