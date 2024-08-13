package fielder

type ConditionalFieldWDefault interface {
	Conditional
	Field
	Default
}

type conditionalFieldWDefault struct {
	Conditional
	Default
	Field
}

func NewCFWD(f Field, prereqs []Prerequisite, defaultField Field) ConditionalFieldWDefault {
	return &conditionalFieldWDefault{
		Conditional: Conditions(prereqs...),
		Default:     NewDefault(true, defaultField),
		Field:       f,
	}
}

func NewEmptyCFWD(prereqs []Prerequisite, defaultField Field) ConditionalFieldWDefault {
	return &conditionalFieldWDefault{
		Conditional: Conditions(prereqs...),
		Default:     NewDefault(false, defaultField),
		Field:       defaultField,
	}
}
