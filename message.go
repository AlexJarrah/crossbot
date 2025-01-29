package crossbot

type Message struct {
	Content           string
	Title             string
	Description       string
	URL               string
	Color             int
	ThumbnailImageURL string
	Footer            Footer
	Fields            []Field
	Buttons           [][]Button
}

type Footer struct {
	Text    string
	IconURL string
}

type Field struct {
	Name   string
	Value  string
	Inline bool
}

type Button struct {
	Label    string
	Emoji    string
	Callback Callback
}

type Callback struct {
	Action       CallbackAction
	Fields       string
	Function     func(map[string]string) *Message
	Prompt       Prompt
	AlertMessage func(map[string]string) string
}

type Prompt struct {
	Prefix string
	Fields []CallbackPromptField
}

type CallbackPromptField struct {
	Key         string
	Value       string
	Placeholder string
}

type CallbackAction uint8

const (
	CallbackActionEditMessage CallbackAction = iota
	CallbackActionCreateMessage
	CallbackActionDeleteMessage
	CallbackActionPrompt
	CallbackActionAlert
)
