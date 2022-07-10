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

func TestFilterGroup_FillValue(t *testing.T) {
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
					NewFilter("Id", PredicateIs, WithFilterValue(1)),
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
					NewFilter("Id", PredicateBetween, WithFilterValue([]any{1, 3})),
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
					NewFilter("Id", PredicateIs, WithFilterValue(1)),
					NewFilter("Id", PredicateBetween, WithFilterValue([]any{1, 3})),
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
						NewFilter("Id", PredicateIs, WithFilterValue(1)),
					}, LogicOperatorAnd),
					NewFilterGroupWithFilters([]*Filter{
						NewFilter("Name", PredicateContains, WithFilterValue("Lily")),
						NewFilter("Id", PredicateBetween, WithFilterValue([]any{1, 3})),
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
			err := tt.group.FillValue(tt.args.values)
			if (err != nil) != tt.wantErr {
				t.Errorf("FillValue() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.group != nil {
				fmt.Println("actual =", tt.group.String())
			}
			if tt.wantGroup != nil {
				fmt.Println("expect =", tt.wantGroup.String())
			}
			if !tt.wantErr && !reflect.DeepEqual(tt.group, tt.wantGroup) {
				t.Errorf("FillValue() actual = %v, expect = %v", tt.group, tt.wantGroup)
			}
		})
	}
}
