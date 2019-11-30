package model

type Step interface {
	GetKind() string
	GetDescription() string
	Execute()
}

type RunnableStep struct {
	Kind        string
	Description string
}
