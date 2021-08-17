package definition

type TargetType string

const (
	TargetStruct TargetType = "struct"
	TargetField  TargetType = "field"
	TargetMethod TargetType = "method"
	TargetAll    TargetType = "all"
)

type Annotation struct {
	Target  TargetType
	Name    string
	Payload string
}

func NewAnnotation(name string, payload string) *Annotation {
	return &Annotation{
		Target:  TargetAll,
		Name:    name,
		Payload: payload,
	}
}

func NewAnnotationStruct(name string, payload string) *Annotation {
	a := NewAnnotation(name, payload)
	a.Target = TargetStruct
	return a
}

func NewAnnotationField(name string, payload string) *Annotation {
	a := NewAnnotation(name, payload)
	a.Target = TargetField
	return a
}

func NewAnnotationMethod(name string, payload string) *Annotation {
	a := NewAnnotation(name, payload)
	a.Target = TargetMethod
	return a
}
