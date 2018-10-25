package rmgo

const (
	Eq Operator = iota
	Gt
	Gte
	Lt
	Lte
)

var opMap = map[Operator]func(target, v Message) bool{
	Eq: EqFn,
}

func (op Operator) Exec(target, v Message) bool {
	fn := opMap[op]
	return fn(target, v)
}

func EqFn(target, v Message) bool {
	return target == v
}

//Is func returns an eq expression
func Is(v Message) Exp {
	return Exp{Eq: v}
}
