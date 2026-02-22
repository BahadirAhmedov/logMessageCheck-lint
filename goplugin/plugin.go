package main

import (
	"github.com/BahadirAhmedov/LogMessageCheck/lint/pkg/analyzer"
	"golang.org/x/tools/go/analysis"
)

func New(conf any) ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{analyzer.Analyzer}, nil
}
