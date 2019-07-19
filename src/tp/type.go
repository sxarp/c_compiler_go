package tp

import "fmt"

type TypeElm interface {
	ptrOf() TypeElm
	literal() string
	Size() int
}

type eightBytesType string

func (it eightBytesType) ptrOf() TypeElm {
	return nil
}

func (it eightBytesType) literal() string {
	return string(it)
}

func (it eightBytesType) Size() int {
	return 8
}

type ptrType struct {
	base TypeElm
}

func (pt ptrType) ptrOf() TypeElm {
	return pt.base
}

func (pt ptrType) literal() string {
	return fmt.Sprintf("ptr of %s", pt.base.literal())
}

func (pt ptrType) Size() int {
	return 8
}

var Int = Type{eightBytesType("int")}

type Type struct {
	TypeElm
}

func (tp Type) Eq(rhs Type) bool {
	if lhsp := tp.ptrOf(); lhsp != nil {
		rhsp := rhs.ptrOf()

		if rhsp != nil {
			return Type{lhsp}.Eq(Type{rhsp})
		}

		return false
	}

	if rhs.ptrOf() != nil {
		return false
	}

	return tp.literal() == rhs.literal()
}

func (tp Type) Ptr() Type {
	return Type{ptrType{tp.TypeElm}}
}

func (tp Type) DeRef() (Type, bool) {
	if tp.ptrOf() != nil {
		return Type{tp.ptrOf()}, true
	}

	return Type{}, false
}

func (tp Type) Size() int {
	return tp.TypeElm.Size()
}

func (tp Type) AddUnit() int {
	switch {
	case tp.Eq(Int):
		return 1
	case tp.ptrOf() != nil:
		return tp.ptrOf().Size()
	default:
		panic("AddUnit cannot be determined.")
	}
}
