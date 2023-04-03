package main

import (
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/asmdecl"
	"golang.org/x/tools/go/analysis/passes/assign"
	"golang.org/x/tools/go/analysis/passes/atomic"
	"golang.org/x/tools/go/analysis/passes/bools"
	"golang.org/x/tools/go/analysis/passes/buildtag"
	"golang.org/x/tools/go/analysis/passes/cgocall"
	"golang.org/x/tools/go/analysis/passes/composite"
	"golang.org/x/tools/go/analysis/passes/copylock"
	"golang.org/x/tools/go/analysis/passes/errorsas"
	"golang.org/x/tools/go/analysis/passes/httpresponse"
	"golang.org/x/tools/go/analysis/passes/loopclosure"
	"golang.org/x/tools/go/analysis/passes/lostcancel"
	"golang.org/x/tools/go/analysis/passes/nilfunc"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shift"
	"golang.org/x/tools/go/analysis/passes/stdmethods"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"golang.org/x/tools/go/analysis/passes/tests"
	"golang.org/x/tools/go/analysis/passes/unmarshal"
	"golang.org/x/tools/go/analysis/passes/unreachable"
	"golang.org/x/tools/go/analysis/passes/unsafeptr"
	"golang.org/x/tools/go/analysis/passes/unusedresult"
	"honnef.co/go/tools/simple"
	"honnef.co/go/tools/staticcheck"

	"github.com/fev0ks/ydx-goadv-metrics/cmd/staticlint/checkers"
)

// Создайте свой multichecker, состоящий из:
// * стандартных статических анализаторов пакета golang.org/x/tools/go/analysis/passes;
// * всех анализаторов класса SA пакета staticcheck.io;
// * не менее одного анализатора остальных классов пакета staticcheck.io;
// * двух или более любых публичных анализаторов на ваш выбор.
// * Напишите и добавьте в multichecker собственный анализатор, запрещающий использовать прямой вызов os.Exit в функции main пакета main.
// При необходимости перепишите код своего проекта так, чтобы он удовлетворял данному анализатору.
// Поместите анализатор в директорию cmd/staticlint вашего проекта.
// Добавьте документацию в формате godoc, подробно опишите в ней механизм запуска multichecker, а также каждый анализатор и его назначение.
// Исходный код вашего проекта должен проходить статический анализ созданного multichecker.

// main go run cmd/staticlint/prj_multichecker.go ./...
func main() {
	myChecks := make([]*analysis.Analyzer, 0)
	myChecks = appendPassesChecks(myChecks)
	//I have no idea why honnef.co/go/tools@v0.0.1-2020.1.4/ir/methods.go:239
	//*types.TypeParam
	// panic: T
	//myChecks = appendSAStaticChecks(myChecks)
	//myChecks = appendSimpleStaticChecks(myChecks)
	myChecks = appendCustomChecks(myChecks)
	multichecker.Main(
		myChecks...,
	)
}

func appendPassesChecks(myChecks []*analysis.Analyzer) []*analysis.Analyzer {
	passesChecks := []*analysis.Analyzer{
		asmdecl.Analyzer,
		assign.Analyzer,
		atomic.Analyzer,
		bools.Analyzer,
		buildtag.Analyzer,
		cgocall.Analyzer,
		composite.Analyzer,
		copylock.Analyzer,
		errorsas.Analyzer,
		httpresponse.Analyzer,
		loopclosure.Analyzer,
		lostcancel.Analyzer,
		nilfunc.Analyzer,
		printf.Analyzer,
		shift.Analyzer,
		stdmethods.Analyzer,
		structtag.Analyzer,
		tests.Analyzer,
		unmarshal.Analyzer,
		unreachable.Analyzer,
		unsafeptr.Analyzer,
		unusedresult.Analyzer,
	}
	return append(myChecks, passesChecks...)
}

func appendSAStaticChecks(myChecks []*analysis.Analyzer) []*analysis.Analyzer {
	for _, v := range staticcheck.Analyzers {
		myChecks = append(myChecks, v)
	}
	return myChecks
}

func appendSimpleStaticChecks(myChecks []*analysis.Analyzer) []*analysis.Analyzer {
	for _, v := range simple.Analyzers {
		myChecks = append(myChecks, v)
	}
	return myChecks
}

func appendCustomChecks(myChecks []*analysis.Analyzer) []*analysis.Analyzer {
	return append(myChecks, checkers.CheckMainExit)
}
