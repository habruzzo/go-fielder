package fielder

type ConditionalField interface {
	Field
	Conditional
}

type Prerequisite struct {
	// function that determines if we should apply the checks
	IsCandidate Enforceable
	// set of required functions to run
	Gauntlet []Question
}

type Enforceable func(f any) bool

var (
	EnforceableTrue = func(f any) bool {
		return true
	}
	EnforceableFalse = func(f any) bool {
		return false
	}
)

type Question func() Enforceable

type Conditional interface {
	// returns the set of keys and values that have to be present for the new field to be set
	Prerequisites() []Prerequisite
	Meets(any) bool
}

type conditional struct {
	prereqs []Prerequisite
}

func Conditions(in ...Prerequisite) Conditional {
	return &conditional{prereqs: in}
}

func (c *conditional) Prerequisites() []Prerequisite {
	return c.prereqs
}

func (c *conditional) Meets(toSet any) bool {
	for _, v := range c.prereqs {
		// if its a candidate for this prerequisite, then we test
		if v.IsCandidate(toSet) {
			// run through all the tests
			for _, w := range v.Gauntlet {
				if !w()(toSet) {
					// if one of the tests fails, we reject
					return false
				}
			}
		}
	}
	return true
}

type FieldConditional struct {
	Field
	Conditional
}

func NewConditionalField(field Field, cond Conditional) ConditionalField {
	return &FieldConditional{
		Field:       field,
		Conditional: cond,
	}
}

func (s *FieldConditional) SetValue(intendedToSet FieldValue) {
	// first we do the safety check and convert to a field
	fieldIntended := intendedToSet.(Field)
	if s.Conditional.Meets(fieldIntended) {
		s.Field.SetValue(fieldIntended)
		return
	}
	return
}

// example vars to illustrate the idea
/*
our type, ExampleParent, looks like this
type ExampleParent struct {
	Grass bool
	Green bool
	HasColor bool
}
in this example, we have an object that represents our state
state := ExampleParent{Grass: true, Green: true, HasColor: false,}
we want to change "has color" to true
HasColor is a ConditionalField , the field inside is specifically a BoolField, with BoolField{Key:"HasColor", value: false, set: true,}
the Conditional is defined to protect the inner BoolField from being in a bad state. if we want to change the BoolField, we have to pass the Questions in the Conditional
a ConditionalField is a field, so its value can be set by calling SetValue
we will override SetValue and make it validate the conditional first, and pass any checks that are needed, if we can, then we set the value
we want to set that field to true, so we will attempt the conditional. you can see inside "SetValue" we call "Meets" which runs through the functions to do an analysis of the parent and the intended outcome
to set the field, we must pass all the criteria in the gauntlet, and we need to define our gauntlet for "HasColor" field
if we want to set HasColor to true, then Green should also be true. if we can only detect the color green, then for something to have color,
it has to be green. so if we want to set HasColor, we need to check that Green is true
the example prerequisite shows one of the tests in th
the example gauntlet item is a representation of a question we will ask to see if we can set the new value
*/

//type ExampleParent struct {
//	Grass    bool `field:"Grass"`
//	Green    bool `field:"Green"`
//	HasColor bool `field:"HasColor"`
//}
//
//var (
//	EP = ExampleParent{
//		Grass:    false,
//		Green:    true,
//		HasColor: false,
//	}
//	// you can define the specific constants or variables you need inline and then pass them
//	// for flexibility, the Question type allows a Parent as a parameter, but if you are using default
//	// field keys then you can leave the Parent empty/nil
//	ExampleGauntletItem = func() Enforceable {
//		// maybe we want to compare two totally separate fields
//		return func(f Field) bool {
//			pri := GetResultItemFieldFromKeyDefault[ExampleParent](EP, NewDefaultFieldKey("Green"))
//			if pri.ToString() == "true" {
//				return true
//			}
//			return false
//		}
//	}
//	ExamplePrerequisite = Prerequisite{
//		// if the "incoming" value (the value our field will become) is true, then
//		// we will enforce the tests in the gauntlet
//		IsCandidate: func(f Field) bool {
//			return f.Key().Name == "HasColor" && f.ToString() == "true"
//		},
//		Gauntlet: []Question{
//			ExampleGauntletItem,
//		},
//	}
//)
//
//func ExampleMain() {
//	intendedField := NewConditionalField(&BoolField{KeyField: NewDefaultFieldKey("HasColor"), ValueField: false, Set: true,}, Conditions(ExamplePrerequisite))
//	intendedField.SetValue(true)
//}
