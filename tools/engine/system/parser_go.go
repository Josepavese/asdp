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

	// 1. Walk RECURSIVELY (Boundary-Aware)
	ignoredPatterns := []string{
		"node_modules", "vendor", "dist", "build",
		".vscode", ".idea", ".git", ".cache",
	}

	isIgnored := func(name string) bool {
		name = strings.ToLower(name)
		for _, p := range ignoredPatterns {
			if strings.Contains(name, p) {
				return true
			}
		}
		return false
	}

	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if path == root {
				return nil
			}
			if isIgnored(d.Name()) {
				return filepath.SkipDir
			}
			// Boundary Check
			if _, err := os.Stat(filepath.Join(path, "codespec.md")); err == nil {
				return filepath.SkipDir
			}
			if _, err := os.Stat(filepath.Join(path, "codemodel.md")); err == nil {
				return filepath.SkipDir
			}
			return nil
		}

		if !strings.HasSuffix(d.Name(), ".go") || strings.HasSuffix(d.Name(), "_test.go") {
			return nil
		}

		// Parse individual file
		f, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
		if err != nil {
			// Log error but continue? Or fail?
			// For robustness, let's log and continue
			return nil
		}

		// Extract symbols from file
		for _, decl := range f.Decls {
			// 1. Functions
			if fn, ok := decl.(*ast.FuncDecl); ok {
				pos := fset.Position(fn.Pos())
				sym := domain.Symbol{
					Name:      fn.Name.Name,
					Kind:      "function",
					Exported:  fn.Name.IsExported(),
					Line:      pos.Line,
					LineEnd:   fset.Position(fn.End()).Line,
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
							LineEnd:   fset.Position(typeSpec.End()).Line,
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

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk go directory %s: %w", root, err)
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
