module github.com/BahadirAhmedov/LogMessageCheck/lint

go 1.24.7

require (
	github.com/golangci/plugin-module-register v0.1.2
	go.uber.org/zap v1.27.1
	golang.org/x/tools v0.42.0
)

require (
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/stretchr/testify v1.11.1 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/mod v0.33.0 // indirect
	golang.org/x/sync v0.19.0 // indirect
)

replace github.com/BahadirAhmedov/logMessageCheck-lint => ./

replace github.com/BahadirAhmedov/LogMessageCheck/lint/pkg/analyzer => ./pkg/analyzer

replace github.com/BahadirAhmedov/LogMessageCheck/lint => ./
