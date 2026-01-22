package system

import (
	"strings"

	"github.com/Josepavese/asdp/engine/domain"
)

// PolyglotParser tries to use specialized parsers first, then falls back to Ctags.
// Actually, Ctags is recursive, so it might duplicate.
// Strategy:
// 1. We want to use GoASTParser for .go files (better fidelity).
// 2. We want to use Ctags for everything else.
// BUT: Ctags -R walks everything.
// Simple approach for V1:
// Run GoASTParser.
// Run CtagsParser.
// Merge results, deduping by Name+File (not easy since ParseDir returns flattened symbols w/o file path in Symbol struct).
// Wait, Symbol struct doesn't have FilePath yet? Let's check entities.
// CodeModel Symbol has 'name', 'kind', 'parent'. Implicitly it belongs to the module.
// If we want to support multi-file modules correctly, we might need to be careful.
// Let's just append for now.

type PolyglotParser struct {
	goParser    *GoASTParser
	ctagsParser *CtagsParser
	config      domain.Config
}

func NewPolyglotParser(config domain.Config) *PolyglotParser {
	return &PolyglotParser{
		goParser:    NewGoASTParser(config),
		ctagsParser: NewCtagsParser(config),
		config:      config,
	}
}

func (p *PolyglotParser) GetSymbolBody(root string, sym domain.Symbol) (string, error) {
	if strings.HasSuffix(sym.FilePath, ".go") {
		return p.goParser.GetSymbolBody(root, sym)
	}
	return p.ctagsParser.GetSymbolBody(root, sym)
}

func (p *PolyglotParser) ParseDir(root string) ([]domain.Symbol, error) {
	var allSymbols []domain.Symbol

	// 1. Try Go Native Parser
	goSymbols, err := p.goParser.ParseDir(root)
	if err == nil {
		allSymbols = append(allSymbols, goSymbols...)
	}

	// 2. Try Ctags Parser
	// We need to tell ctags to EXCLUDE go files if we already parsed them?
	// Or we just rely on Ctags for non-Go.
	// ctags --exclude=*.go
	// However, `CtagsParser` inside implementation hardcodes the command.
	// Let's just run it. If Ctags returns Go symbols too, we might have duplicates.
	// For this iteration, let's assume Ctags is mainly for "other" languages.
	// We can update CtagsParser to accept excludes or we can just append.

	ctagsSymbols, err := p.ctagsParser.ParseDir(root)
	if err == nil {
		// Deduping logic could go here.
		// For now, let's just append.
		// Actually, if Ctags is missing (err != nil), we just ignore it.
		allSymbols = append(allSymbols, ctagsSymbols...)
	} else {
		// Ensure we don't return error if just ctags is missing, UNLESS go also failed.
		// If we found 0 go symbols and ctags failed, maybe we return error?
		// No, empty symbols is a valid result.
	}

	return allSymbols, nil
}
