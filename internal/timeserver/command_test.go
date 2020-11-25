package timeserver_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/komem3/goalarm/internal/timeserver"
)

func TestAllCommands(t *testing.T) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "command.go", nil, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}
	var wantCommand []string
	for _, d := range file.Decls {
		if gen, ok := d.(*ast.GenDecl); ok {
			for _, s := range gen.Specs {
				vs, ok := s.(*ast.ValueSpec)
				if !ok {
					break
				}
				if value, ok := vs.Values[0].(*ast.BasicLit); ok {
					wantCommand = append(wantCommand, strings.Trim(value.Value, "\""))
				}
			}
		}
	}
	var givenCommand []string
	for _, s := range timeserver.AllCommands() {
		givenCommand = append(givenCommand, string(s.Command))
	}
	if diff := cmp.Diff(givenCommand, wantCommand); diff != "" {
		t.Errorf("all commands, given(-), want(+)\n%s\n", diff)
	}
}
