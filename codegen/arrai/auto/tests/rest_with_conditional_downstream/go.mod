module rest_with_conditional_downstream

go 1.14

require (
	github.com/anz-bank/pkg v0.0.28
	github.com/anz-bank/sysl-go v0.84.0
	github.com/go-chi/chi v4.1.2+incompatible
	github.com/go-playground/universal-translator v0.17.0 // indirect
	github.com/leodido/go-urn v1.2.0 // indirect
	github.com/mattn/go-sqlite3 v2.0.3+incompatible
	github.com/rickb777/date v1.14.0
	github.com/sethvargo/go-retry v0.1.0
	github.com/spf13/afero v1.4.0
	github.com/stretchr/testify v1.6.1
	gopkg.in/go-playground/assert.v1 v1.2.1 // indirect
)

replace github.com/anz-bank/sysl-go => ../../../../..
