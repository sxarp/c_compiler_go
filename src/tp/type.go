package tp

type Type struct {
	ptrOf *Type
	basic string
}

var Int = Type{basic: "int"}

func (tp Type) Eq(rhs Type) bool {
	if lhsp := tp.ptrOf; lhsp != nil {
		rhsp := rhs.ptrOf

		if rhsp != nil {
			return lhsp.Eq(*rhsp)
		}

		return false
	}

	if rhs.ptrOf != nil {
		return false
	}

	return tp.basic == rhs.basic
}

func (tp Type) Ptr() Type {
	return Type{ptrOf: &tp}
}

func (tp Type) DeRef() (Type, bool) {
	if tp.ptrOf != nil {
		return *(tp.ptrOf), true
	}

	return Type{}, false
}

func (tp Type) Size() int {
	return 8
}

func (tp Type) AddUnit() int {
	switch {
	case tp.Eq(Int):
		return 1
	case tp.ptrOf != nil:
		return tp.ptrOf.Size()
	default:
		panic("AddUnit cannot be determined.")
	}
}
