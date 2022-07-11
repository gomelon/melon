package query

import "context"

type Translator interface {
	Translate(ctx context.Context, query *Query) (string, error)
	TranslateFind(ctx context.Context, query *Query) (string, error)
	TranslateCount(ctx context.Context, query *Query) (string, error)
	TranslateExists(ctx context.Context, query *Query) (string, error)
	TranslateDelete(ctx context.Context, query *Query) (string, error)
	TranslateTable(ctx context.Context, table Table) (string, error)
	TranslateFilterGroup(ctx context.Context, group *FilterGroup) (string, error)
	TranslateFilter(ctx context.Context, filter *Filter) (string, error)
	TranslateLogicOperator(ctx context.Context, operator LogicOperator) (string, error)
	TranslateSorts(ctx context.Context, sorts []*Sort) (string, error)
	TranslateSort(ctx context.Context, sort *Sort) (string, error)
	TranslatePager(ctx context.Context, pager Pager) (string, error)
}
