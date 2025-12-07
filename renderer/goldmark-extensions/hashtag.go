package goldmarkextension

import (
	"regexp"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

// Define AST Node
type Hashtag struct {
	ast.BaseInline
}

func (n *Hashtag) Dump(source []byte, level int) {
	ast.DumpHelper(n, source, level, nil, nil)
}

var KindHashtag = ast.NewNodeKind("Hashtag")

func (n *Hashtag) Kind() ast.NodeKind {
	return KindHashtag
}

// Create Parser
type hashtagParser struct{}

func (s *hashtagParser) Trigger() []byte {
	return []byte{'#'}
}

var hashtagRegex = regexp.MustCompile(`^#[a-zA-Z0-9_]+`)

func (s *hashtagParser) Parse(parent ast.Node, block text.Reader, pc parser.Context) ast.Node {
	before := block.PrecendingCharacter()

	// Only match if preceded by whitespace or start of line
	if before != ' ' && before != '\n' && before != -1 {
		return nil
	}

	line, lineSegment := block.PeekLine()
	match := hashtagRegex.FindSubmatch(line)

	if match == nil {
		return nil
	}

	// Add child TextSegment to store the tag name
	node := &Hashtag{}
	textNode := ast.NewTextSegment(
		text.NewSegment(lineSegment.Start+1, lineSegment.Start+len(match[0])),
	)
	node.AppendChild(node, textNode)
	block.Advance(len(match[0]))

	return node
}

// Create Renderer
type hashtagHTMLRenderer struct{}

func (r *hashtagHTMLRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(KindHashtag, r.renderHashtag)
}

func (r *hashtagHTMLRenderer) renderHashtag(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		escapedTag := util.EscapeHTML(node.Text(source))

		_, _ = w.WriteString(`<wa-tag size="small" appearance="filled" pill><a href="/tags/"`)
		_, _ = w.Write(escapedTag)
		_, _ = w.WriteString(`">#`)
		_, _ = w.Write(escapedTag)
		_, _ = w.WriteString(`</a></wa-tag>`)

		return ast.WalkSkipChildren, nil
	}
	return ast.WalkContinue, nil
}

// Create Extension
type HashtagExtension struct{}

func (e *HashtagExtension) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(
		parser.WithInlineParsers(
			util.Prioritized(&hashtagParser{}, 500),
		),
	)
	m.Renderer().AddOptions(
		renderer.WithNodeRenderers(
			util.Prioritized(&hashtagHTMLRenderer{}, 500),
		),
	)
}
