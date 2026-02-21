package factory

import "fmt"

// ProviderFactory is a lightweight registry to resolve AI providers by name.
// It is intentionally generic so domain/usecase layers stay provider-agnostic.
type ProviderFactory[T any] struct {
	providers map[string]T
}

func NewProviderFactory[T any]() *ProviderFactory[T] {
	return &ProviderFactory[T]{providers: map[string]T{}}
}

func (f *ProviderFactory[T]) Register(name string, provider T) {
	f.providers[name] = provider
}

func (f *ProviderFactory[T]) Resolve(name string) (T, error) {
	provider, ok := f.providers[name]
	if !ok {
		var zero T
		return zero, fmt.Errorf("ai provider not registered: %s", name)
	}
	return provider, nil
}
