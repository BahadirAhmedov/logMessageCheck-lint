package plugin

import (
	"github.com/BahadirAhmedov/LogMessageCheck/lint/pkg/analyzer"
	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"
)

func init() {
	register.Plugin("mycustomlinter", New)
}

func New(conf any) (register.LinterPlugin, error) {
	return &logMessagePlugin{}, nil
}

type logMessagePlugin struct{}

func (*logMessagePlugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{analyzer.Analyzer}, nil
}

func (*logMessagePlugin) GetLoadMode() string {
	return register.LoadModeTypesInfo
}
