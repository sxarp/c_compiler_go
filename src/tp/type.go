package tp

type Type struct {
	ptrOf *Type
	basic string
}

var Int = Type{basic: "int"}

func (lhs Type) Eq(rhs Type) bool {
	if lhsp := lhs.ptrOf; lhsp != nil {
		if rhsp := rhs.ptrOf; rhsp != nil {
			return lhsp.Eq(*rhsp)

		} else {
			return false
		}
	}

	if rhs.ptrOf != nil {
		return false
	}

	return lhs.basic == rhs.basic
}

func (org Type) Ptr() Type {
	return Type{ptrOf: &org}
}
