package system

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
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

	// Parse just THIS directory (non-recursive)
	pkgs, err := parser.ParseDir(fset, root, func(fi os.FileInfo) bool {
		return !strings.HasSuffix(fi.Name(), "_test.go") && !strings.HasPrefix(fi.Name(), ".")
	}, parser.ParseComments)

	if err != nil {
		return nil, fmt.Errorf("failed to parse go directory %s: %w", root, err)
	}

	for _, pkg := range pkgs {
		for _, file := range pkg.Files {
			for _, decl := range file.Decls {
				// 1. Functions
				if fn, ok := decl.(*ast.FuncDecl); ok {
					pos := fset.Position(fn.Pos())
					sym := domain.Symbol{
						Name:      fn.Name.Name,
						Kind:      "function",
						Exported:  fn.Name.IsExported(),
						Line:      pos.Line,
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
							pos := fset.Position(typeSpec.Pos())
							sym := domain.Symbol{
								Name:      typeSpec.Name.Name,
								Exported:  typeSpec.Name.IsExported(),
								Line:      pos.Line,
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

	return symbols, nil
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
