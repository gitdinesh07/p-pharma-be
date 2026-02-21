package auth

type ContextKey string

const (
	ContextPrincipalKey ContextKey = "principal"
	ContextAPIKeyKey    ContextKey = "api_key_principal"
)
