package docs

import _ "embed"

//go:embed openapi.yaml
var openAPISpec string

func OpenAPISpec() string {
	return openAPISpec
}
