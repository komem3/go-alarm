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

func TestAllStatuses(t *testing.T) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "status.go", nil, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}
	var wantStatus []string
	for _, d := range file.Decls {
		if gen, ok := d.(*ast.GenDecl); ok {
			for _, s := range gen.Specs {
				vs, ok := s.(*ast.ValueSpec)
				if !ok {
					break
				}
				if value, ok := vs.Values[0].(*ast.BasicLit); ok {
					wantStatus = append(wantStatus, strings.Trim(value.Value, "\""))
				}
			}
		}
	}
	var givenStatus []string
	for _, s := range timeserver.AllStatuses() {
		givenStatus = append(givenStatus, string(s.Status))
	}
	if diff := cmp.Diff(givenStatus, wantStatus); diff != "" {
		t.Errorf("all status, given(-), want(+)\n%s\n", diff)
	}
}
