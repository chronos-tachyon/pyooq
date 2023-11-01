package pyooq

func Identity[T ValueExpr](x T) OpExpr {
	t := x.Type()
	return makeUnaryOp(OpIdentity, t, x)
}

func Add[TA ValueExpr, TB ValueExpr](a TA, b TB) OpExpr {
	at := a.Type()
	bt := b.Type()
	t := InvalidType
	if at.WidensTo(bt) {
		t = bt
	}
	if bt.WidensTo(at) {
		t = at
	}
	return makeBinaryOp(OpAdd, t, a, b)
}

// Not()  => 1-op NOT
// And()  => 2-op AND
// Nand() => 2-op NAND
// Or()   => 2-op OR
// Nor()  => 2-op NOR
// Xor()  => 2-op XOR
// Xnor() => 2-op XNOR
// All()  => multi-op AND
// Any()  => multi-op OR
// EQ()   => 2-op EQ
// NE()   => 2-op NE
// LT()   => 2-op LT
// LE()   => 2-op LE
// GE()   => 2-op GE
// GT()   => 2-op GT
