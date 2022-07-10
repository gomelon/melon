package query

import "fmt"

type Sort struct {
	fieldName string
	direction Direction
}

func NewSort(fieldName string, direction Direction) *Sort {
	return &Sort{fieldName: fieldName, direction: direction}
}

func (s Sort) String() string {
	return fmt.Sprintf("%s %s", s.fieldName, s.direction)
}

func (s *Sort) FieldName() string {
	return s.fieldName
}

func (s *Sort) Direction() Direction {
	return s.direction
}

type Direction string

func (d Direction) String() string {
	return string(d)
}

const (
	DirectionDesc Direction = "Desc"
	DirectionAsc  Direction = "Asc"
)
