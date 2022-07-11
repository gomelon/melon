package query

import (
	"fmt"
	"reflect"
	"testing"
)

func TestFilterGroup_NumValue(t *testing.T) {
	tests := []struct {
		name    string
		group   *FilterGroup
		wantNum int
	}{
		{
			name: "one filter with zero value",
			group: NewFilterGroupWithFilters(
				[]*Filter{
					NewFilter("Name", PredicateIsEmpty),
				},
				LogicOperatorAnd,
			),
			wantNum: 0,
		},
		{
			name: "one filter with one value",
			group: NewFilterGroupWithFilters(
				[]*Filter{
					NewFilter("Id", PredicateIs),
				},
				LogicOperatorAnd,
			),
			wantNum: 1,
		},
		{
			name: "one filter with two value",
			group: NewFilterGroupWithFilters(
				[]*Filter{
					NewFilter("Id", PredicateBetween),
				},
				LogicOperatorAnd,
			),
			wantNum: 2,
		},
		{
			name: "mix filter",
			group: NewFilterGroupWithFilters(
				[]*Filter{
					NewFilter("Name", PredicateIsEmpty),
					NewFilter("Id", PredicateIs),
					NewFilter("Id", PredicateBetween),
				},
				LogicOperatorAnd,
			),
			wantNum: 3,
		},
		{
			name: "filter groups with ",
			group: NewFilterGroup(
				[]*FilterGroup{
					NewFilterGroupWithFilters([]*Filter{
						NewFilter("Name", PredicateIsEmpty),
						NewFilter("Id", PredicateIs),
					}, LogicOperatorAnd),
					NewFilterGroupWithFilters([]*Filter{
						NewFilter("Name", PredicateContains),
						NewFilter("Id", PredicateBetween),
					}, LogicOperatorAnd),
				},
				LogicOperatorOr,
			),
			wantNum: 4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotNum := tt.group.NumValue(); gotNum != tt.wantNum {
				t.Errorf("NumValue() = %v, want %v", gotNum, tt.wantNum)
			}
		})
	}
}

func TestFilterGroup_FillValues(t *testing.T) {
	type args struct {
		values []any
	}
	tests := []struct {
		name      string
		group     *FilterGroup
		args      args
		wantGroup *FilterGroup
		wantErr   bool
	}{
		{
			name: "not enough values",
			group: NewFilterGroupWithFilters(
				[]*Filter{
					NewFilter("Name", PredicateIs),
				},
				LogicOperatorAnd,
			),
			args: args{
				values: nil,
			},
			wantErr: true,
		},
		{
			name: "too many values",
			group: NewFilterGroupWithFilters(
				[]*Filter{
					NewFilter("Name", PredicateIs),
				},
				LogicOperatorAnd,
			),
			args: args{
				values: []any{1, 1},
			},
			wantErr: true,
		},
		{
			name: "one filter with zero value",
			group: NewFilterGroupWithFilters(
				[]*Filter{
					NewFilter("Name", PredicateIsEmpty),
				},
				LogicOperatorAnd,
			),
			args: args{
				values: nil,
			},
			wantGroup: NewFilterGroupWithFilters(
				[]*Filter{
					NewFilter("Name", PredicateIsEmpty),
				},
				LogicOperatorAnd,
			),
			wantErr: false,
		},
		{
			name: "one filter with one value",
			group: NewFilterGroupWithFilters(
				[]*Filter{
					NewFilter("Id", PredicateIs),
				},
				LogicOperatorAnd,
			),
			wantGroup: NewFilterGroupWithFilters(
				[]*Filter{
					NewFilter("Id", PredicateIs, WithFilterValues(1)),
				},
				LogicOperatorAnd,
			),
			args: args{
				values: []any{1},
			},
			wantErr: false,
		},
		{
			name: "one filter with two value",
			group: NewFilterGroupWithFilters(
				[]*Filter{
					NewFilter("Id", PredicateBetween),
				},
				LogicOperatorAnd,
			),
			wantGroup: NewFilterGroupWithFilters(
				[]*Filter{
					NewFilter("Id", PredicateBetween, WithFilterValues(1, 3)),
				},
				LogicOperatorAnd,
			),
			args: args{
				values: []any{1, 3},
			},
			wantErr: false,
		},
		{
			name: "mix filter",
			group: NewFilterGroupWithFilters(
				[]*Filter{
					NewFilter("Name", PredicateIsEmpty),
					NewFilter("Id", PredicateIs),
					NewFilter("Id", PredicateBetween),
				},
				LogicOperatorAnd,
			),
			wantGroup: NewFilterGroupWithFilters(
				[]*Filter{
					NewFilter("Name", PredicateIsEmpty),
					NewFilter("Id", PredicateIs, WithFilterValues(1)),
					NewFilter("Id", PredicateBetween, WithFilterValues(1, 3)),
				},
				LogicOperatorAnd,
			),
			args: args{
				values: []any{1, 1, 3},
			},
			wantErr: false,
		},
		{
			name: "filter groups with ",
			group: NewFilterGroup(
				[]*FilterGroup{
					NewFilterGroupWithFilters([]*Filter{
						NewFilter("Name", PredicateIsEmpty),
						NewFilter("Id", PredicateIs),
					}, LogicOperatorAnd),
					NewFilterGroupWithFilters([]*Filter{
						NewFilter("Name", PredicateContains),
						NewFilter("Id", PredicateBetween),
					}, LogicOperatorAnd),
				},
				LogicOperatorOr,
			),
			wantGroup: NewFilterGroup(
				[]*FilterGroup{
					NewFilterGroupWithFilters([]*Filter{
						NewFilter("Name", PredicateIsEmpty),
						NewFilter("Id", PredicateIs, WithFilterValues(1)),
					}, LogicOperatorAnd),
					NewFilterGroupWithFilters([]*Filter{
						NewFilter("Name", PredicateContains, WithFilterValues("Lily")),
						NewFilter("Id", PredicateBetween, WithFilterValues(1, 3)),
					}, LogicOperatorAnd),
				},
				LogicOperatorOr,
			),
			args: args{
				values: []any{1, "Lily", 1, 3},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.group.FillValues(tt.args.values)
			if (err != nil) != tt.wantErr {
				t.Errorf("FillValues() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.group != nil {
				fmt.Println("actual =", tt.group.String())
			}
			if tt.wantGroup != nil {
				fmt.Println("expect =", tt.wantGroup.String())
			}
			if !tt.wantErr && !reflect.DeepEqual(tt.group, tt.wantGroup) {
				t.Errorf("FillValues() actual = %v, expect = %v", tt.group, tt.wantGroup)
			}
		})
	}
}

func TestFilterGroup_FillNamedArgs(t *testing.T) {
	type args struct {
		namedArgs []string
	}
	tests := []struct {
		name      string
		group     *FilterGroup
		args      args
		wantGroup *FilterGroup
		wantErr   bool
	}{
		{
			name: "not enough values",
			group: NewFilterGroupWithFilters(
				[]*Filter{
					NewFilter("Name", PredicateIs),
				},
				LogicOperatorAnd,
			),
			args: args{
				namedArgs: nil,
			},
			wantErr: true,
		},
		{
			name: "too many values",
			group: NewFilterGroupWithFilters(
				[]*Filter{
					NewFilter("Name", PredicateIs),
				},
				LogicOperatorAnd,
			),
			args: args{
				namedArgs: []string{"name1", "name2"},
			},
			wantErr: true,
		},
		{
			name: "one filter with zero value",
			group: NewFilterGroupWithFilters(
				[]*Filter{
					NewFilter("Name", PredicateIsEmpty),
				},
				LogicOperatorAnd,
			),
			args: args{
				namedArgs: nil,
			},
			wantGroup: NewFilterGroupWithFilters(
				[]*Filter{
					NewFilter("Name", PredicateIsEmpty),
				},
				LogicOperatorAnd,
			),
			wantErr: false,
		},
		{
			name: "one filter with one value",
			group: NewFilterGroupWithFilters(
				[]*Filter{
					NewFilter("Id", PredicateIs),
				},
				LogicOperatorAnd,
			),
			wantGroup: NewFilterGroupWithFilters(
				[]*Filter{
					NewFilter("Id", PredicateIs, WithFilterNamedArgs("id")),
				},
				LogicOperatorAnd,
			),
			args: args{
				namedArgs: []string{"id"},
			},
			wantErr: false,
		},
		{
			name: "one filter with two value",
			group: NewFilterGroupWithFilters(
				[]*Filter{
					NewFilter("Id", PredicateBetween),
				},
				LogicOperatorAnd,
			),
			wantGroup: NewFilterGroupWithFilters(
				[]*Filter{
					NewFilter("Id", PredicateBetween, WithFilterNamedArgs("min_id", "max_id")),
				},
				LogicOperatorAnd,
			),
			args: args{
				namedArgs: []string{"min_id", "max_id"},
			},
			wantErr: false,
		},
		{
			name: "mix filter",
			group: NewFilterGroupWithFilters(
				[]*Filter{
					NewFilter("Name", PredicateIsEmpty),
					NewFilter("Id", PredicateIs),
					NewFilter("Id", PredicateBetween),
				},
				LogicOperatorAnd,
			),
			wantGroup: NewFilterGroupWithFilters(
				[]*Filter{
					NewFilter("Name", PredicateIsEmpty),
					NewFilter("Id", PredicateIs, WithFilterNamedArgs("id")),
					NewFilter("Id", PredicateBetween, WithFilterNamedArgs("min_id", "max_id")),
				},
				LogicOperatorAnd,
			),
			args: args{
				namedArgs: []string{"id", "min_id", "max_id"},
			},
			wantErr: false,
		},
		{
			name: "filter groups with ",
			group: NewFilterGroup(
				[]*FilterGroup{
					NewFilterGroupWithFilters([]*Filter{
						NewFilter("Name", PredicateIsEmpty),
						NewFilter("Id", PredicateIs),
					}, LogicOperatorAnd),
					NewFilterGroupWithFilters([]*Filter{
						NewFilter("Name", PredicateContains),
						NewFilter("Id", PredicateBetween),
					}, LogicOperatorAnd),
				},
				LogicOperatorOr,
			),
			wantGroup: NewFilterGroup(
				[]*FilterGroup{
					NewFilterGroupWithFilters([]*Filter{
						NewFilter("Name", PredicateIsEmpty),
						NewFilter("Id", PredicateIs, WithFilterNamedArgs("id")),
					}, LogicOperatorAnd),
					NewFilterGroupWithFilters([]*Filter{
						NewFilter("Name", PredicateContains, WithFilterNamedArgs("name")),
						NewFilter("Id", PredicateBetween, WithFilterNamedArgs("min_id", "max_id")),
					}, LogicOperatorAnd),
				},
				LogicOperatorOr,
			),
			args: args{
				namedArgs: []string{"id", "name", "min_id", "max_id"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.group.FillNamedArgs(tt.args.namedArgs)
			if (err != nil) != tt.wantErr {
				t.Errorf("FillValues() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.group != nil {
				fmt.Println("actual =", tt.group.String())
			}
			if tt.wantGroup != nil {
				fmt.Println("expect =", tt.wantGroup.String())
			}
			if !tt.wantErr && !reflect.DeepEqual(tt.group, tt.wantGroup) {
				t.Errorf("FillValues() actual = %v, expect = %v", tt.group, tt.wantGroup)
			}
		})
	}
}
