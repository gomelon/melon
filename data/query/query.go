package query

import (
	"fmt"
	"strings"
)

type Query struct {
	//TODO 要添加projection及pageable
	//要添加table的赋值
	table               Table
	subject             *Subject
	subjectModifier     *SubjectModifier
	subjectModifierArgs map[SubjectModifierArg]any
	filterGroup         *FilterGroup
	sorts               []*Sort
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
	return builder.String()
}
