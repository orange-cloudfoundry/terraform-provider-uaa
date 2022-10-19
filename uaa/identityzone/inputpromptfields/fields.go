package inputpromptfields

type InputPromptField int64

const (
	Name InputPromptField = iota
	Type
	Value
)

func (s InputPromptField) String() string {
	switch s {
	case Name:
		return "name"
	case Type:
		return "type"
	case Value:
		return "value"
	}
	return "unknown"
}
