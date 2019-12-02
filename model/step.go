package model

type Step interface {
	GetKind() string
	GetDescription() string
	Execute() Result
}

type RunnableStep struct {
	Kind        string
	Description string
}
