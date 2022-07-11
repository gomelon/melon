package query

import (
	"context"
	"github.com/gomelon/melon/data/engine"
	"testing"
)

func TestRDBTranslator_TranslateFind(t1 *testing.T) {
	tests := []struct {
		name        string
		query       *Query
		engines     map[string]engine.Engine
		wantResults map[string]string
		wantErr     bool
	}{
		{
			name: "find one filter one sort pager",
			query: New(
				SubjectFind,
				WithTable(NewTable("user")),
				WithFilterGroup(
					NewFilterGroupWithFilters(
						[]*Filter{
							NewFilter("Id", PredicateIs),
						},
						LogicOperatorAnd),
				),
				WithSorts(
					[]*Sort{
						NewSort("Firstname", DirectionAsc),
					},
				),
				WithPager(NewPageRequest(1, 10, false)),
			),
			engines: map[string]engine.Engine{"MySQL": engine.NewMySQL()},
			wantResults: map[string]string{
				"MySQL": "SELECT * FROM `user` WHERE (`id` = ?) ORDER BY `firstname` ASC LIMIT ?, ?",
			},
			wantErr: false,
		},
		{
			name: "find multiple filter one sort pager",
			query: New(
				SubjectFind,
				WithTable(NewTable("user")),
				WithFilterGroup(
					NewFilterGroupWithFilters(
						[]*Filter{
							NewFilter("Id", PredicateIs),
							NewFilter("Name", PredicateIs),
						},
						LogicOperatorAnd),
				),
				WithSorts(
					[]*Sort{
						NewSort("Firstname", DirectionAsc),
					},
				),
				WithPager(NewPageRequest(1, 10, false)),
			),
			engines: map[string]engine.Engine{"MySQL": engine.NewMySQL()},
			wantResults: map[string]string{
				"MySQL": "SELECT * FROM `user` WHERE ((`id` = ?) AND (`name` = ?)) ORDER BY `firstname` ASC LIMIT ?, ?",
			},
			wantErr: false,
		},
		{
			name: "find multiple filter group multiple sort pager",
			query: New(
				SubjectFind,
				WithTable(NewTable("user")),
				WithFilterGroup(
					NewFilterGroup(
						[]*FilterGroup{
							NewFilterGroupWithFilters(
								[]*Filter{
									NewFilter("Id", PredicateIs),
									NewFilter("Name", PredicateContains),
								},
								LogicOperatorAnd),
							NewFilterGroupWithFilters(
								[]*Filter{
									NewFilter("Age", PredicateGTE),
								},
								LogicOperatorAnd),
						},
						LogicOperatorOr),
				),
				WithSorts(
					[]*Sort{
						NewSort("Firstname", DirectionAsc),
						NewSort("Lastname", DirectionDesc),
					},
				),
				WithPager(NewPageRequest(1, 10, false)),
			),
			engines: map[string]engine.Engine{"MySQL": engine.NewMySQL()},
			wantResults: map[string]string{
				"MySQL": "SELECT * FROM `user` " +
					"WHERE (((`id` = ?) AND (`name` LIKE CONCAT('%',?,'%'))) OR (`age` >= ?)) " +
					"ORDER BY `firstname` ASC, `lastname` DESC LIMIT ?, ?",
			},
			wantErr: false,
		},
		{
			name: "find distinct one filter one sort pager",
			query: New(
				SubjectFind,
				WithSubjectModifier(SubjectModifierDistinct),
				WithTable(NewTable("user")),
				WithFilterGroup(
					NewFilterGroupWithFilters(
						[]*Filter{
							NewFilter("Id", PredicateIs),
						},
						LogicOperatorAnd),
				),
				WithSorts(
					[]*Sort{
						NewSort("Firstname", DirectionAsc),
					},
				),
				WithPager(NewPageRequest(1, 10, false)),
			),
			engines: map[string]engine.Engine{"MySQL": engine.NewMySQL()},
			wantResults: map[string]string{
				"MySQL": "SELECT DISTINCT * FROM `user` WHERE (`id` = ?) ORDER BY `firstname` ASC LIMIT ?, ?",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		for dialect, dbEngine := range tt.engines {
			t1.Run(tt.name, func(t1 *testing.T) {
				translator := NewRDBTranslator(dbEngine)
				gotResult, err := translator.Translate(context.Background(), tt.query)
				if (err != nil) != tt.wantErr {
					t1.Errorf("TranslateFind() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				wantResult := tt.wantResults[dialect]
				if gotResult != wantResult {
					t1.Errorf("TranslateFind() \nactual = %v, \nexpect = %v", gotResult, wantResult)
				}
			})
		}

	}
}

func TestRDBTranslator_TranslateCount(t1 *testing.T) {
	tests := []struct {
		name        string
		query       *Query
		engines     map[string]engine.Engine
		wantResults map[string]string
		wantErr     bool
	}{
		{
			name: "count one filter one",
			query: New(
				SubjectCount,
				WithTable(NewTable("user")),
				WithFilterGroup(
					NewFilterGroupWithFilters(
						[]*Filter{
							NewFilter("Id", PredicateIs),
						},
						LogicOperatorAnd),
				),
			),
			engines: map[string]engine.Engine{"MySQL": engine.NewMySQL()},
			wantResults: map[string]string{
				"MySQL": "SELECT COUNT(*) FROM `user` WHERE (`id` = ?)",
			},
			wantErr: false,
		},
		{
			name: "count multiple filter",
			query: New(
				SubjectCount,
				WithTable(NewTable("user")),
				WithFilterGroup(
					NewFilterGroupWithFilters(
						[]*Filter{
							NewFilter("Id", PredicateIs),
							NewFilter("Name", PredicateIs),
						},
						LogicOperatorAnd),
				),
			),
			engines: map[string]engine.Engine{"MySQL": engine.NewMySQL()},
			wantResults: map[string]string{
				"MySQL": "SELECT COUNT(*) FROM `user` WHERE ((`id` = ?) AND (`name` = ?))",
			},
			wantErr: false,
		},
		{
			name: "find multiple filter group",
			query: New(
				SubjectCount,
				WithTable(NewTable("user")),
				WithFilterGroup(
					NewFilterGroup(
						[]*FilterGroup{
							NewFilterGroupWithFilters(
								[]*Filter{
									NewFilter("Id", PredicateIs),
									NewFilter("Name", PredicateContains),
								},
								LogicOperatorAnd),
							NewFilterGroupWithFilters(
								[]*Filter{
									NewFilter("Age", PredicateGTE),
								},
								LogicOperatorAnd),
						},
						LogicOperatorOr),
				),
			),
			engines: map[string]engine.Engine{"MySQL": engine.NewMySQL()},
			wantResults: map[string]string{
				"MySQL": "SELECT COUNT(*) FROM `user` " +
					"WHERE (((`id` = ?) AND (`name` LIKE CONCAT('%',?,'%'))) OR (`age` >= ?))",
			},
			wantErr: false,
		},
		{
			name: "find distinct one filter",
			query: New(
				SubjectCount,
				WithSubjectModifier(SubjectModifierDistinct),
				WithTable(NewTable("user")),
				WithFilterGroup(
					NewFilterGroupWithFilters(
						[]*Filter{
							NewFilter("Id", PredicateIs),
						},
						LogicOperatorAnd),
				),
				WithSorts(
					[]*Sort{
						NewSort("Firstname", DirectionAsc),
					},
				),
				WithPager(NewPageRequest(1, 10, false)),
			),
			engines: map[string]engine.Engine{"MySQL": engine.NewMySQL()},
			wantResults: map[string]string{
				"MySQL": "SELECT COUNT(DISTINCT *) FROM `user` WHERE (`id` = ?)",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		for dialect, dbEngine := range tt.engines {
			t1.Run(tt.name, func(t1 *testing.T) {
				translator := NewRDBTranslator(dbEngine)
				gotResult, err := translator.Translate(context.Background(), tt.query)
				if (err != nil) != tt.wantErr {
					t1.Errorf("TranslateFind() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				wantResult := tt.wantResults[dialect]
				if gotResult != wantResult {
					t1.Errorf("TranslateFind() \nactual = %v, \nexpect = %v", gotResult, wantResult)
				}
			})
		}

	}
}

func TestRDBTranslator_TranslateExists(t1 *testing.T) {
	tests := []struct {
		name        string
		query       *Query
		engines     map[string]engine.Engine
		wantResults map[string]string
		wantErr     bool
	}{
		{
			name: "exists one filter one",
			query: New(
				SubjectExists,
				WithTable(NewTable("user")),
				WithFilterGroup(
					NewFilterGroupWithFilters(
						[]*Filter{
							NewFilter("Id", PredicateIs),
						},
						LogicOperatorAnd),
				),
			),
			engines: map[string]engine.Engine{"MySQL": engine.NewMySQL()},
			wantResults: map[string]string{
				"MySQL": "SELECT 1 FROM `user` WHERE (`id` = ?) LIMIT 0, 1",
			},
			wantErr: false,
		},
		{
			name: "exists multiple filter",
			query: New(
				SubjectExists,
				WithTable(NewTable("user")),
				WithFilterGroup(
					NewFilterGroupWithFilters(
						[]*Filter{
							NewFilter("Id", PredicateIs),
							NewFilter("Name", PredicateIs),
						},
						LogicOperatorAnd),
				),
			),
			engines: map[string]engine.Engine{"MySQL": engine.NewMySQL()},
			wantResults: map[string]string{
				"MySQL": "SELECT 1 FROM `user` WHERE ((`id` = ?) AND (`name` = ?)) LIMIT 0, 1",
			},
			wantErr: false,
		},
		{
			name: "exists multiple filter group",
			query: New(
				SubjectExists,
				WithTable(NewTable("user")),
				WithFilterGroup(
					NewFilterGroup(
						[]*FilterGroup{
							NewFilterGroupWithFilters(
								[]*Filter{
									NewFilter("Id", PredicateIs),
									NewFilter("Name", PredicateContains),
								},
								LogicOperatorAnd),
							NewFilterGroupWithFilters(
								[]*Filter{
									NewFilter("Age", PredicateGTE),
								},
								LogicOperatorAnd),
						},
						LogicOperatorOr),
				),
			),
			engines: map[string]engine.Engine{"MySQL": engine.NewMySQL()},
			wantResults: map[string]string{
				"MySQL": "SELECT 1 FROM `user` " +
					"WHERE (((`id` = ?) AND (`name` LIKE CONCAT('%',?,'%'))) OR (`age` >= ?)) " +
					"LIMIT 0, 1",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		for dialect, dbEngine := range tt.engines {
			t1.Run(tt.name, func(t1 *testing.T) {
				translator := NewRDBTranslator(dbEngine)
				gotResult, err := translator.Translate(context.Background(), tt.query)
				if (err != nil) != tt.wantErr {
					t1.Errorf("TranslateFind() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				wantResult := tt.wantResults[dialect]
				if gotResult != wantResult {
					t1.Errorf("TranslateFind() \nactual = %v, \nexpect = %v", gotResult, wantResult)
				}
			})
		}

	}
}

func TestRDBTranslator_TranslateDelete(t1 *testing.T) {
	tests := []struct {
		name        string
		query       *Query
		engines     map[string]engine.Engine
		wantResults map[string]string
		wantErr     bool
	}{
		{
			name: "delete one filter one sort pager",
			query: New(
				SubjectDelete,
				WithTable(NewTable("user")),
				WithFilterGroup(
					NewFilterGroupWithFilters(
						[]*Filter{
							NewFilter("Id", PredicateIs),
						},
						LogicOperatorAnd),
				),
				WithSorts(
					[]*Sort{
						NewSort("Firstname", DirectionAsc),
					},
				),
				WithPager(NewPageRequest(1, 10, false)),
			),
			engines: map[string]engine.Engine{"MySQL": engine.NewMySQL()},
			wantResults: map[string]string{
				"MySQL": "DELETE FROM `user` WHERE (`id` = ?) ORDER BY `firstname` ASC LIMIT ?, ?",
			},
			wantErr: false,
		},
		{
			name: "delete multiple filter one sort pager",
			query: New(
				SubjectDelete,
				WithTable(NewTable("user")),
				WithFilterGroup(
					NewFilterGroupWithFilters(
						[]*Filter{
							NewFilter("Id", PredicateIs),
							NewFilter("Name", PredicateIs),
						},
						LogicOperatorAnd),
				),
				WithSorts(
					[]*Sort{
						NewSort("Firstname", DirectionAsc),
					},
				),
				WithPager(NewPageRequest(1, 10, false)),
			),
			engines: map[string]engine.Engine{"MySQL": engine.NewMySQL()},
			wantResults: map[string]string{
				"MySQL": "DELETE FROM `user` WHERE ((`id` = ?) AND (`name` = ?)) ORDER BY `firstname` ASC LIMIT ?, ?",
			},
			wantErr: false,
		},
		{
			name: "delete multiple filter group multiple sort pager",
			query: New(
				SubjectDelete,
				WithTable(NewTable("user")),
				WithFilterGroup(
					NewFilterGroup(
						[]*FilterGroup{
							NewFilterGroupWithFilters(
								[]*Filter{
									NewFilter("Id", PredicateIs),
									NewFilter("Name", PredicateContains),
								},
								LogicOperatorAnd),
							NewFilterGroupWithFilters(
								[]*Filter{
									NewFilter("Age", PredicateGTE),
								},
								LogicOperatorAnd),
						},
						LogicOperatorOr),
				),
				WithSorts(
					[]*Sort{
						NewSort("Firstname", DirectionAsc),
						NewSort("Lastname", DirectionDesc),
					},
				),
				WithPager(NewPageRequest(1, 10, false)),
			),
			engines: map[string]engine.Engine{"MySQL": engine.NewMySQL()},
			wantResults: map[string]string{
				"MySQL": "DELETE FROM `user` " +
					"WHERE (((`id` = ?) AND (`name` LIKE CONCAT('%',?,'%'))) OR (`age` >= ?)) " +
					"ORDER BY `firstname` ASC, `lastname` DESC LIMIT ?, ?",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		for dialect, dbEngine := range tt.engines {
			t1.Run(tt.name, func(t1 *testing.T) {
				translator := NewRDBTranslator(dbEngine)
				gotResult, err := translator.Translate(context.Background(), tt.query)
				if (err != nil) != tt.wantErr {
					t1.Errorf("TranslateFind() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				wantResult := tt.wantResults[dialect]
				if gotResult != wantResult {
					t1.Errorf("TranslateFind() \nactual = %v, \nexpect = %v", gotResult, wantResult)
				}
			})
		}

	}
}
