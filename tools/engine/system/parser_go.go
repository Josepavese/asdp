package system

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/Josepavese/asdp/engine/domain"
)

type GoASTParser struct {
	fs *RealFileSystem
}

func NewGoASTParser() *GoASTParser {
	return &GoASTParser{fs: NewRealFileSystem()}
}

func (p *GoASTParser) ParseDir(root string) ([]domain.Symbol, error) {
	var symbols []domain.Symbol
	fset := token.NewFileSet()

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			return nil
		}
		// Skip hidden dirs
		if strings.HasPrefix(info.Name(), ".") && path != root {
			return filepath.SkipDir
		}

		// Parse just this directory
		pkgs, err := parser.ParseDir(fset, path, func(fi os.FileInfo) bool {
			return !strings.HasSuffix(fi.Name(), "_test.go") && !strings.HasPrefix(fi.Name(), ".")
		}, parser.ParseComments)

		if err != nil {
			return nil // ignore parse errors in subdirs for now
		}

		for _, pkg := range pkgs {
			for _, file := range pkg.Files {
				for _, decl := range file.Decls {
					// 1. Functions
					if fn, ok := decl.(*ast.FuncDecl); ok {
						sym := domain.Symbol{
							Name:      fn.Name.Name,
							Kind:      "function",
							Exported:  fn.Name.IsExported(),
							Docstring: strings.TrimSpace(fn.Doc.Text()),
							Signature: formatFuncSignature(fn),
						}
						// Check if it's a method
						if fn.Recv != nil {
							sym.Kind = "method"
							for _, field := range fn.Recv.List {
								if star, ok := field.Type.(*ast.StarExpr); ok {
									if ident, ok := star.X.(*ast.Ident); ok {
										sym.Parent = ident.Name
									}
								} else if ident, ok := field.Type.(*ast.Ident); ok {
									sym.Parent = ident.Name
								}
							}
						}
						symbols = append(symbols, sym)
					}

					// 2. Types (Structs/Interfaces)
					if gen, ok := decl.(*ast.GenDecl); ok && gen.Tok == token.TYPE {
						for _, spec := range gen.Specs {
							if typeSpec, ok := spec.(*ast.TypeSpec); ok {
								sym := domain.Symbol{
									Name:      typeSpec.Name.Name,
									Exported:  typeSpec.Name.IsExported(),
									Docstring: strings.TrimSpace(gen.Doc.Text()),
									Signature: fmt.Sprintf("type %s", typeSpec.Name.Name),
								}

								switch typeSpec.Type.(type) {
								case *ast.StructType:
									sym.Kind = "struct"
								case *ast.InterfaceType:
									sym.Kind = "interface"
								default:
									sym.Kind = "type"
								}
								symbols = append(symbols, sym)
							}
						}
					}
				}
			}
		}
		return nil
	})

	return symbols, err
}

// Helper to reconstruct signature string roughly
func formatFuncSignature(fn *ast.FuncDecl) string {
	// Ideally we use printer.Fprint, but simple reconstruction is fine for now
	sig := "func "
	if fn.Recv != nil {
		sig += "(...)" // Method
	}
	sig += fn.Name.Name + "(...)"
	return sig
}
