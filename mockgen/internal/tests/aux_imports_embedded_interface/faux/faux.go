package faux

type Foreign interface {
	Method() Return
	Embedded
	error
}

type Embedded any

type Return any
