package cc

// OperatorOnlyKinds are the message kinds a customer must never receive: notes
// operators leave for each other, and the operator's conversation with the AI
// assistant. Everything that decides visibility keys off this one list, so a
// kind added to it is hidden everywhere at once rather than in three of the four
// places that filter.
//
// Kept as a hand-written file next to the generated ones: `buf generate`
// rewrites cc.pb.go and would drop this.
var OperatorOnlyKinds = []Kind{Kind_ADMIN_ONLY, Kind_COPILOT}

// IsOperatorOnly reports whether this kind stays inside the support team.
func IsOperatorOnly(k Kind) bool {
	for _, o := range OperatorOnlyKinds {
		if k == o {
			return true
		}
	}
	return false
}
