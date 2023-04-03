package main

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/fev0ks/ydx-goadv-metrics/cmd/staticlint/checkers"
)

func TestCheckMainExit(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), checkers.CheckMainExit, "./...")
}
