package parser

// Priority is the predefined
// parsing order of operations
type Priority int

const (
	_ Priority = iota
	// LOWEST priority
	LOWEST
	// EQUALS eval equality; ie. '=='
	EQUALS
	// LESSGREATER eval; ie. '>' or '<'
	LESSGREATER
	// SUM addition; ie. '+'
	SUM
	// PRODUCT multiplication; ie. '*'
	PRODUCT
	// PREFIX operator in front of operand; ie. '-X'
	PREFIX
	// CALL function invocations; ie 'myfunc(X)'
	CALL
)

func (p Priority) String() string {
	switch p {
	case LOWEST:
		return "LOWEST"
	case EQUALS:
		return "EQUALS"
	case LESSGREATER:
		return "LESSGREATER"
	case SUM:
		return "SUM"
	case PRODUCT:
		return "PRODUCT"
	case PREFIX:
		return "PREFIX"
	case CALL:
		return "CALL"
	default:
		return "UNKNOWN"
	}
}
