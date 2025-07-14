package hooks

import "fmt"

type Events []Event

type Event struct {
	Type             EventType
	Literal          string
	SecondaryLiteral string
}

type EventType int

// any other type is literal, uses the Literal field as String() output try
// to have explicit types instead of implicit ones with a literal type
const (
	None EventType = iota
	Add
	Edit
	UpdateKey
	Delete
	DeleteAll
	Do
	Undo
	AddDep
)

func (events Events) String() (output string) {
	for _, event := range events {
		if event.Type != None {
			output += event.String() + "\n"
		}
	}
	return
}

func (events Events) Empty() bool {
	for _, event := range events {
		if event.Type != None {
			return false
		}
	}
	return true
}


func (event *Event) String() string {
	switch event.Type {
	case None:
		return ""
	case Add:
		return fmt.Sprintf("Add %q", event.Literal)
	case Edit:
		return fmt.Sprintf("Update %q to %q", event.SecondaryLiteral, event.Literal)
	case Delete:
		return fmt.Sprintf("Remove %q", event.Literal)
	case DeleteAll:
		return "Remove all tasks"
	case Do:
		return fmt.Sprintf("chore: Do %q", event.Literal)
	case Undo:
		return fmt.Sprintf("chore: Undo %q", event.Literal)
	case UpdateKey:
		return fmt.Sprintf("Update key from %q to %q", event.SecondaryLiteral, event.Literal)
	case AddDep:
		return fmt.Sprintf("Add %q to %q dependencies", event.SecondaryLiteral, event.Literal)
	default:
		return event.Literal
	}
}
