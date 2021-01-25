package terraform

//go:generate go run golang.org/x/tools/cmd/stringer -type=GraphType context_graph_type.go

// GraphType is an enum of the type of graph to create with a Context.
// The values of the constants may change so they shouldn't be depended on;
// always use the constant name.
type GraphType byte

const (
	GraphTypeInvalid GraphType = iota
	GraphTypePlan
	GraphTypePlanDestroy
	GraphTypeApply
	GraphTypeValidate
	GraphTypeEval // only visits in-memory elements such as variables, locals, and outputs.
)

// GraphTypeMap is a mapping of human-readable string to GraphType. This
// is useful to use as the mechanism for human input for configurable
// graph types.
var GraphTypeMap = map[string]GraphType{
	"apply":        GraphTypeApply,
	"plan":         GraphTypePlan,
	"plan-destroy": GraphTypePlanDestroy,
	"validate":     GraphTypeValidate,
	"eval":         GraphTypeEval,
}
