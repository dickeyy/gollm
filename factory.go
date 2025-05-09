package gollm

import (
	"fmt"
	"log"
	"sync"
)

// ModelConstructor defines the function signature for creating an LLM instance.
// It takes apiKey and an optional apiURL.
type ModelConstructor func(modelName string, apiKey string, apiURL string) (LLM, error)

var (
	modelRegistry = make(map[string]ModelConstructor)
	registryMutex = &sync.RWMutex{}
)

// RegisterModel allows LLM providers to register their constructors.
// This function is intended to be called from the init() function of provider packages.
func RegisterModel(modelName string, constructor ModelConstructor) {
	registryMutex.Lock()
	defer registryMutex.Unlock()
	if _, exists := modelRegistry[modelName]; exists {
		log.Printf("Warning: Overriding registered model constructor for model '%s'", modelName)
	}
	modelRegistry[modelName] = constructor
}

// InitializeModel creates an LLM instance based on the model name using the registry.
// It takes modelName and apiKey. The apiURL is handled by the constructor if needed.
func InitializeModel(modelName string, apiKey string) (LLM, error) {
	registryMutex.RLock()
	constructor, exists := modelRegistry[modelName]
	registryMutex.RUnlock()

	if !exists {
		return nil, fmt.Errorf("unsupported model: %s (not registered)", modelName)
	}

	// Pass an empty string for apiURL; constructors can use a default if it's empty.
	return constructor(modelName, apiKey, "")
}
