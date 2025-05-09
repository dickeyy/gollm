package gollm

type LLM interface {
	Chat(structure ChatStructure) (*ChatResponse, error)
}
