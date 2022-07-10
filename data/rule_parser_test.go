package data

import (
	"fmt"
	"github.com/gomelon/melon/data/query"
	"reflect"
	"testing"
)

func TestRuleParser_Parse_Find(t *testing.T) {
	tests := []struct {
		name       string
		methodName string
		fieldNames []string
		wantQuery  *query.Query
		wantErr    bool
	}{
		{
			name:       "simple find",
			methodName: "FindById",
			wantQuery: query.New(
				query.SubjectFind,
				query.WithFilterGroup(
					query.NewFilterGroupWithFilters([]*query.Filter{
						query.NewFilter("Id", query.PredicateIs),
					}, query.LogicOperatorAnd),
				),
			),
			wantErr: false,
		},
		{
			name:       "simple find and",
			methodName: "FindByIdAndName",
			wantQuery: query.New(
				query.SubjectFind,
				query.WithFilterGroup(
					query.NewFilterGroupWithFilters([]*query.Filter{
						query.NewFilter("Id", query.PredicateIs),
						query.NewFilter("Name", query.PredicateIs),
					}, query.LogicOperatorAnd)),
			),
			wantErr: false,
		},
		{
			name:       "simple find or",
			methodName: "FindByIdOrName",
			wantQuery: query.New(
				query.SubjectFind,
				query.WithFilterGroup(
					query.NewFilterGroup([]*query.FilterGroup{
						query.NewFilterGroupWithFilters(
							[]*query.Filter{
								query.NewFilter("Id", query.PredicateIs),
							},
							query.LogicOperatorAnd),
						query.NewFilterGroupWithFilters(
							[]*query.Filter{
								query.NewFilter("Name", query.PredicateIs),
							},
							query.LogicOperatorAnd),
					}, query.LogicOperatorOr),
				),
			),
			wantErr: false,
		},
		{
			name:       "filter find and order by",
			methodName: "FindByIdAndNameOrderByFirstnameAscLastnameDesc",
			wantQuery: query.New(
				query.SubjectFind,
				query.WithFilterGroup(
					query.NewFilterGroupWithFilters([]*query.Filter{
						query.NewFilter("Id", query.PredicateIs),
						query.NewFilter("Name", query.PredicateIs),
					}, query.LogicOperatorAnd),
				),
				query.WithSorts(
					[]*query.Sort{
						query.NewSort("Firstname", query.DirectionAsc),
						query.NewSort("Lastname", query.DirectionDesc),
					},
				),
			),
			wantErr: false,
		},
		{
			name:       "filter with predicate and order by",
			methodName: "FindByIdIsAndNameContainsOrderByFirstnameAscLastnameDesc",
			wantQuery: query.New(
				query.SubjectFind,
				query.WithFilterGroup(query.NewFilterGroupWithFilters([]*query.Filter{
					query.NewFilter("Id", query.PredicateIs),
					query.NewFilter("Name", query.PredicateContains),
				}, query.LogicOperatorAnd),
				),
				query.WithSorts([]*query.Sort{
					query.NewSort("Firstname", query.DirectionAsc),
					query.NewSort("Lastname", query.DirectionDesc),
				},
				),
			),
			wantErr: false,
		},
		{
			name:       "filter with predicate and order by",
			methodName: "FindByIdIsAndNameContainsOrderByFirstnameAscLastnameDesc",
			wantQuery: query.New(
				query.SubjectFind,
				query.WithFilterGroup(
					query.NewFilterGroupWithFilters([]*query.Filter{
						query.NewFilter("Id", query.PredicateIs),
						query.NewFilter("Name", query.PredicateContains),
					}, query.LogicOperatorAnd),
				),
				query.WithSorts(
					[]*query.Sort{
						query.NewSort("Firstname", query.DirectionAsc),
						query.NewSort("Lastname", query.DirectionDesc),
					},
				),
			),
			wantErr: false,
		},
		{
			name:       "filter with predicate and or order by",
			methodName: "FindByIdIsAndNameContainsOrAgeGTEOrderByFirstnameAscLastnameDesc",
			wantQuery: query.New(
				query.SubjectFind,
				query.WithFilterGroup(
					query.NewFilterGroup([]*query.FilterGroup{
						query.NewFilterGroupWithFilters([]*query.Filter{
							query.NewFilter("Id", query.PredicateIs),
							query.NewFilter("Name", query.PredicateContains),
						}, query.LogicOperatorAnd),
						query.NewFilterGroupWithFilters([]*query.Filter{
							query.NewFilter("Age", query.PredicateGTE),
						}, query.LogicOperatorAnd),
					}, query.LogicOperatorOr),
				),
				query.WithSorts(
					[]*query.Sort{
						query.NewSort("Firstname", query.DirectionAsc),
						query.NewSort("Lastname", query.DirectionDesc),
					},
				),
			),
			wantErr: false,
		},
		{
			name:       "find distinct filter and order by",
			methodName: "FindDistinctByIdAndNameOrderByFirstname",
			wantQuery: query.New(
				query.SubjectFind,
				query.WithSubjectModifier(query.SubjectModifierDistinct),
				query.WithFilterGroup(
					query.NewFilterGroupWithFilters([]*query.Filter{
						query.NewFilter("Id", query.PredicateIs),
						query.NewFilter("Name", query.PredicateIs),
					}, query.LogicOperatorAnd),
				),
				query.WithSorts(
					[]*query.Sort{
						query.NewSort("Firstname", query.DirectionAsc),
					},
				),
			),
			wantErr: false,
		},
		{
			name:       "find top filter and order by",
			methodName: "FindTop10ByIdAndNameOrderByFirstname",
			wantQuery: query.New(
				query.SubjectFind,
				query.WithSubjectModifier(query.SubjectModifierTop),
				query.WithSubjectModifierArgs(map[query.SubjectModifierArg]any{query.SubjectModifierArgLimit: 10}),
				query.WithFilterGroup(
					query.NewFilterGroupWithFilters([]*query.Filter{
						query.NewFilter("Id", query.PredicateIs),
						query.NewFilter("Name", query.PredicateIs),
					}, query.LogicOperatorAnd),
				),
				query.WithSorts(
					[]*query.Sort{
						query.NewSort("Firstname", query.DirectionAsc),
					},
				),
			),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRuleParser()
			gotQuery, err := r.Parse(tt.methodName)
			if gotQuery != nil {
				fmt.Println("actual =", gotQuery.String())
			}
			if tt.wantQuery != nil {
				fmt.Println("expect =", tt.wantQuery.String())
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotQuery, tt.wantQuery) {
				t.Errorf("Parse() \nactual = %+v, \nexpect = %+v", gotQuery, tt.wantQuery)
			}
		})
	}
}

func TestRuleParser_Parse_Count(t *testing.T) {
	tests := []struct {
		name       string
		methodName string
		fieldNames []string
		wantQuery  *query.Query
		wantErr    bool
	}{
		{
			name:       "simple count",
			methodName: "CountByName",
			wantQuery: query.New(
				query.SubjectCount,
				query.WithFilterGroup(
					query.NewFilterGroupWithFilters([]*query.Filter{
						query.NewFilter("Name", query.PredicateIs),
					}, query.LogicOperatorAnd),
				),
			),
			wantErr: false,
		},
		{
			name:       "simple count and",
			methodName: "CountByIdAndName",
			wantQuery: query.New(
				query.SubjectCount,
				query.WithFilterGroup(
					query.NewFilterGroupWithFilters([]*query.Filter{
						query.NewFilter("Id", query.PredicateIs),
						query.NewFilter("Name", query.PredicateIs),
					}, query.LogicOperatorAnd),
				),
			),
			wantErr: false,
		},
		{
			name:       "simple count or",
			methodName: "CountByIdOrName",
			wantQuery: query.New(
				query.SubjectCount,
				query.WithFilterGroup(
					query.NewFilterGroup([]*query.FilterGroup{
						query.NewFilterGroupWithFilters(
							[]*query.Filter{
								query.NewFilter("Id", query.PredicateIs),
							},
							query.LogicOperatorAnd),
						query.NewFilterGroupWithFilters(
							[]*query.Filter{
								query.NewFilter("Name", query.PredicateIs),
							},
							query.LogicOperatorAnd),
					}, query.LogicOperatorOr),
				),
			),
			wantErr: false,
		},
		{
			name:       "filter count and order by",
			methodName: "CountByIdOrderByFirstname",
			wantQuery:  nil,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRuleParser()
			gotQuery, err := r.Parse(tt.methodName)
			if gotQuery != nil {
				fmt.Println("actual =", gotQuery.String())
			}
			if tt.wantQuery != nil {
				fmt.Println("expect =", tt.wantQuery.String())
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotQuery, tt.wantQuery) {
				t.Errorf("Parse() actual = %v, expect = %v", gotQuery, tt.wantQuery)
			}
		})
	}
}

func TestRuleParser_splitByOrKeyword(t *testing.T) {
	type args struct {
		str string
	}
	var tests = []struct {
		name string
		args args
		want []string
	}{
		{
			name: "simple none or",
			args: args{
				str: "Id",
			},
			want: []string{"Id"},
		},
		{
			name: "simple one or",
			args: args{
				str: "IdOrName",
			},
			want: []string{"Id", "Name"},
		},
		{
			name: "simple two or",
			args: args{
				str: "IdOrNameOrAge",
			},
			want: []string{"Id", "Name", "Age"},
		},
		{
			name: "start with or",
			args: args{
				str: "OrNameOrAge",
			},
			want: []string{"OrName", "Age"},
		},
		{
			name: "end with or",
			args: args{
				str: "IdOr",
			},
			want: []string{"IdOr"},
		},
		{
			name: "two or and one is ending",
			args: args{
				str: "IdOrNameOr",
			},
			want: []string{"Id", "NameOr"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRuleParser()
			if got := r.splitByOrKeyword(tt.args.str); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("splitByOrKeyword() = %v, want %v", got, tt.want)
			}
		})
	}
}
