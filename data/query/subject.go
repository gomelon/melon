package query

type Subject struct {
	keywords []string
	sortable bool
}

func (s *Subject) Keywords() []string {
	return s.keywords
}

func (s *Subject) Sortable() bool {
	return s.sortable
}

func (s Subject) String() string {
	return s.keywords[0]
}

var (
	SubjectFind   = &Subject{keywords: []string{"Find", "Query", "Get", "Search"}, sortable: true}
	SubjectCount  = &Subject{keywords: []string{"Count"}, sortable: false}
	SubjectExists = &Subject{keywords: []string{"Exists"}, sortable: false}
	SubjectDelete = &Subject{keywords: []string{"Delete", "Remove"}, sortable: true}
)

var Subjects = []*Subject{
	SubjectFind, SubjectCount, SubjectExists, SubjectDelete,
}

type SubjectModifier struct {
	keywords []string
	subjects map[*Subject]bool //the modifier supported Subject
}

func (s *SubjectModifier) Keywords() []string {
	return s.keywords
}

func (s *SubjectModifier) Subjects() map[*Subject]bool {
	return s.subjects
}

func (s SubjectModifier) String() string {
	return s.keywords[0]
}

var (
	SubjectModifierDistinct = &SubjectModifier{
		keywords: []string{"Distinct"},
		subjects: map[*Subject]bool{SubjectFind: true, SubjectCount: true},
	}
	SubjectModifierTop = &SubjectModifier{
		keywords: []string{"Top"},
		subjects: map[*Subject]bool{SubjectFind: true, SubjectDelete: true},
	}
)

var SubjectModifiers = []*SubjectModifier{
	SubjectModifierDistinct, SubjectModifierTop,
}

type SubjectModifierArg string

const (
	SubjectModifierArgLimit = "Limit"
)
