module rest_jwt_authorization

go 1.14

replace github.com/anz-bank/sysl-go => ../../../../..

require (
	github.com/anz-bank/pkg v0.0.28
	github.com/anz-bank/sysl-go v0.0.0-20200325045908-46c4ce0a2736
	github.com/go-chi/chi v4.1.2+incompatible
	github.com/go-playground/universal-translator v0.17.0 // indirect
	github.com/golang/protobuf v1.4.2
	github.com/leodido/go-urn v1.2.0 // indirect
	github.com/rickb777/date v1.12.4
	github.com/sethvargo/go-retry v0.1.0
	github.com/spf13/afero v1.4.0
	github.com/stretchr/testify v1.6.1
	google.golang.org/grpc v1.32.0
	google.golang.org/protobuf v1.25.0
	gopkg.in/go-playground/assert.v1 v1.2.1 // indirect
	gopkg.in/yaml.v2 v2.3.0
)
