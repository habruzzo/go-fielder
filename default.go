package fielder

type FieldWDefault interface {
	Field
	Default
	IsDefault() bool
}

type Default interface {
	ExplicitlySet() bool
	MatchesDefault(Field) bool
	DefaultField() Field
}

func NewDefault(explicitly bool, f Field) Default {
	return &defaulter{
		ExplicitlySetField: explicitly,
		Value:              f,
	}
}

type defaulter struct {
	ExplicitlySetField bool `dynamodbav:"explicitly_set" json:"explicitly_set"`
	Value              Field
}

func (d *defaulter) ExplicitlySet() bool {
	return d.ExplicitlySetField
}

func (d *defaulter) MatchesDefault(f Field) bool {
	return d.Value.Equal(f)
}

func (d *defaulter) DefaultField() Field {
	return d.Value
}

type FieldWDefaultImpl struct {
	Field
	Default
}

func (s *FieldWDefaultImpl) IsDefault() bool {
	return s.Default.MatchesDefault(s.Field) && !s.Default.ExplicitlySet()
}

func NewFieldWDefault(f Field, d Default) FieldWDefault {
	return &FieldWDefaultImpl{
		Field:   f,
		Default: d,
	}
}
