package data

import (
	"fmt"
	"github.com/gomelon/melon/data/query"
	"regexp"
	"strconv"
	"strings"
)

const (
	keywordBy      = "By"
	keywordOrderBy = "OrderBy"
)

//RuleParser
//Format:     $Subject [$SubjectModifier] [@Filter] [$Sort]
//
//Subject:
//    Find Query Get Search:General query method returning typically the repository type,slice or struct
//    Count                : Count projection returning a numeric result.
//    Exists               : Exists projection, returning typically a boolean result.
//    Delete Remove        : Delete query method returning either no result (void) or the delete count.
//
//SubjectModifier:
//    Distinct:    Use a distinct query to return only unique results.
//    Top<Number>: Limit the query results to the first <number> of results.
//
//Filter:  By$Field$Predicate[$FilterModifier][And|Or $Field$Predicate[$FilterModifier]]
//    Remark: If you want to support nested fields later, use _ to separate the nesting
//Predicate:
//    Is, Equals, (or no keyword)
//    Contains: for string contains substring or collection contains an element
//    StartsWith EndsWith: for string
//    Between: The BETWEEN operator is inclusive: begin and end values are included.
//    GT LT GTE LTE: for comparable
//    IsNull IsNotNull:
//    IsEmpty IsNotEmpty: for string or collection is empty
//    IsFalse IsTrue: for bool
//    In NotIn:
//    Matches: match the regex
//    wait for support: Exists, ContainsAny(for array), ContainsAll(for array)
//FilterModifier:
//    IgnoreCase:    Used with a predicate keyword for case-insensitive comparison.
//    AllIgnoreCase: Ignore case for all suitable properties. Used somewhere in the query method predicate.
//
//Sort:       Specify a static sorting order followed by the field path and direction
//            Format: OrderBy$Field[$Direction], EX: OrderByFirstnameAscLastnameDesc
//$Direction:
//    Desc Asc: default is Desc
type RuleParser struct {
}

func NewRuleParser() *RuleParser {
	return &RuleParser{}
}

func (r *RuleParser) Parse(method string) (q *query.Query, err error) {
	var nextIndex int
	remaining := method
	subject, nextIndex, err := r.parseSubject(remaining)
	if err != nil {
		return
	}

	remaining = remaining[nextIndex:]
	subjectModifier, subjectModifierArgs, nextIndex, err := r.parseSubjectModifier(subject, remaining)
	if err != nil {
		return
	}

	remaining = remaining[nextIndex:]
	filterGroup, nextIndex, err := r.parseFilters(remaining)
	if err != nil {
		return
	}

	remaining = remaining[nextIndex:]

	var sorts []*query.Sort
	if subject.Sortable() {
		sorts, nextIndex, err = r.parseSort(remaining)
		if err != nil {
			return
		}
	}

	if nextIndex < len(remaining) {
		err = fmt.Errorf("method rule parse fail: can not parse [%s]", remaining)
		return
	}
	q = query.New(subject, query.WithSubjectModifier(subjectModifier),
		query.WithSubjectModifierArgs(subjectModifierArgs), query.WithFilterGroup(filterGroup),
		query.WithSorts(sorts),
	)
	return
}

func (r *RuleParser) parseSubject(str string) (s *query.Subject, nextIndex int, err error) {
	keywords := make([]string, 0, len(query.Subjects)*2)
	for _, subject := range query.Subjects {
		for _, keyword := range subject.Keywords() {
			if strings.HasPrefix(str, keyword) {
				s = subject
				nextIndex = len(keyword)
				return
			}
			keywords = append(keywords, keyword)
		}
	}

	err = fmt.Errorf("method rule parse fail: [%s] can not find subject, method must starts with [%s]",
		str, strings.Join(keywords, "|"))
	return
}

func (r *RuleParser) parseSubjectModifier(subject *query.Subject, str string) (
	modifier *query.SubjectModifier, args map[query.SubjectModifierArg]any, nextIndex int, err error) {

	for _, m := range query.SubjectModifiers {
		if !m.Subjects()[subject] {
			continue
		}
		for _, keyword := range m.Keywords() {
			if strings.HasPrefix(str, keyword) {
				modifier = m
				nextIndex = len(keyword)
				break
			}
		}
		if modifier != nil {
			break
		}
	}

	if query.SubjectModifierTop == modifier {
		topN := -1
		topLen := 0
		for i := nextIndex + 1; i < len(str); i++ {
			s := str[nextIndex:i]
			n, err := strconv.Atoi(s)
			if err != nil {
				break
			}
			topLen++
			topN = n
		}
		if topN <= 0 {
			err = fmt.Errorf("method rule parse fail: [%s] top n is invalid, n must great than 0", str)
			return
		}
		nextIndex += topLen
		args = map[query.SubjectModifierArg]any{query.SubjectModifierArgLimit: topN}
	}
	return
}

func (r *RuleParser) parseFilters(str string) (group *query.FilterGroup, nextIndex int, err error) {
	if !strings.HasPrefix(str, keywordBy) {
		return
	}
	orderByIndex := strings.Index(str, keywordOrderBy)
	var filtersStr string
	if orderByIndex > 0 {
		nextIndex = orderByIndex
		filtersStr = str[len(keywordBy):orderByIndex]
	} else {
		nextIndex = len(str)
		filtersStr = str[len(keywordBy):]
	}

	orParts := r.splitByOrKeyword(filtersStr)
	if len(orParts) == 1 {
		group = r.parseAndFilterGroup(orParts[0])
		return
	}
	groups := make([]*query.FilterGroup, 0, len(orParts))
	for _, part := range orParts {
		groups = append(groups, r.parseAndFilterGroup(part))
	}
	group = query.NewFilterGroup(groups, query.LogicOperatorOr)
	return
}

func (r *RuleParser) parseSort(str string) (sorts []*query.Sort, nextIndex int, err error) {
	if !strings.HasPrefix(str, keywordOrderBy) {
		return
	}

	orderByStr := str[len(keywordOrderBy):]
	nextIndex = len(str)
	parts := make([]string, 0, 2)
	for _, part1 := range r.splitByKeyword(orderByStr, string(query.DirectionAsc)) {
		for _, part2 := range r.splitByKeyword(part1, string(query.DirectionDesc)) {
			parts = append(parts, part2)
		}
	}
	sorts = make([]*query.Sort, 0, len(parts))
	for _, part := range parts {
		if strings.HasSuffix(part, string(query.DirectionDesc)) {
			sorts = append(sorts, query.NewSort(part[:len(part)-len(query.DirectionDesc)], query.DirectionDesc))
		} else if strings.HasSuffix(part, string(query.DirectionAsc)) {
			sorts = append(sorts, query.NewSort(part[:len(part)-len(query.DirectionAsc)], query.DirectionAsc))
		} else {
			sorts = append(sorts, query.NewSort(part, query.DirectionAsc))
		}
	}
	return
}

func (r *RuleParser) parseAndFilterGroup(str string) *query.FilterGroup {
	andParts := r.splitByAndKeyword(str)
	filters := make([]*query.Filter, 0, len(andParts))
	for _, part := range andParts {
		filters = append(filters, r.parseFilter(part))
	}
	return query.NewFilterGroupWithFilters(filters, query.LogicOperatorAnd)
}

func (r *RuleParser) parseFilter(str string) *query.Filter {
	remainingStr := str
	var modifier *query.FilterModifier
	for _, m := range query.FilterModifiers {
		for _, keyword := range m.Keywords() {
			if strings.HasSuffix(remainingStr, keyword) {
				modifier = m
				remainingStr = remainingStr[:len(remainingStr)-len(keyword)]
				break
			}
		}
		if modifier != nil {
			break
		}
	}
	var predicate *query.Predicate
	for _, p := range query.Predicates {
		for _, keyword := range p.Keywords() {
			if strings.HasSuffix(remainingStr, keyword) {
				predicate = p
				remainingStr = remainingStr[:len(remainingStr)-len(keyword)]
				break
			}
		}
		if predicate != nil {
			break
		}
	}
	return query.NewFilter(remainingStr, predicate, query.WithFilterModifier(modifier))
}

func (r *RuleParser) splitByOrKeyword(str string) []string {
	return r.splitByKeyword(str, string(query.LogicOperatorOr))
}

func (r *RuleParser) splitByAndKeyword(str string) []string {
	return r.splitByKeyword(str, string(query.LogicOperatorAnd))
}

func (r *RuleParser) splitByKeyword(str, keyword string) []string {
	compile := regexp.MustCompile(keyword + "[A-Z]+")
	rangeIndexes := compile.FindAllStringIndex(str, 10)
	if len(rangeIndexes) == 0 {
		return []string{str}
	}
	keywordLen := len(keyword)
	parts := make([]string, 0, len(rangeIndexes))
	var lastIndex int
	for _, r := range rangeIndexes {
		if r[0] == 0 {
			continue
		}
		parts = append(parts, str[lastIndex:r[0]])
		lastIndex = r[0] + keywordLen
	}
	if lastIndex > 0 {
		parts = append(parts, str[lastIndex:])
	}
	return parts
}
