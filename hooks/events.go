package hooks

import "fmt"

type Events []Event

type Event struct {
	Type             EventType
	Key          string
	Secondary string
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
	for i := 0; i < len(events); i++ {
		if events[i].Type != None {
			event := &events[i]
			// squash simple events of same kind together
			if (events[i].Type == Add) || events[i].Type == Delete || events[i].Type == Do || events[i].Type == Undo {
				j := i+1
				for j < len(events) && events[j].Type == events[i].Type {
					events[i].Key += ", " + events[j].Key
					j++
				}
				if i != j {
					i = j
				}
			}
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
		return fmt.Sprintf("Add %s", event.Key)
	case Edit:
		return fmt.Sprintf("Update %s to %q", event.Key, event.Secondary)
	case Delete:
		return fmt.Sprintf("Remove %s", event.Key)
	case DeleteAll:
		return "Remove all tasks"
	case Do:
		return fmt.Sprintf("chore: Do %s", event.Key)
	case Undo:
		return fmt.Sprintf("chore: Undo %s", event.Key)
	case UpdateKey:
		return fmt.Sprintf("Update key from %q to %q", event.Secondary, event.Key)
	case AddDep:
		return fmt.Sprintf("Add %q to %q dependencies", event.Secondary, event.Key)
	default:
		return event.Key
	}
}
