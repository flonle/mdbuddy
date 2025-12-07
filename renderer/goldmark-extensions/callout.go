package goldmarkextension

import (
	"regexp"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

// AST Node
type CalloutNode struct {
	ast.BaseBlock
	CalloutType string
	Title       string
}

var KindCallout = ast.NewNodeKind("Callout")

func (n *CalloutNode) Kind() ast.NodeKind { return KindCallout }
func (n *CalloutNode) Dump(source []byte, level int) {
	ast.DumpHelper(n, source, level, map[string]string{
		"CalloutType": n.CalloutType,
		"Title":       n.Title,
	}, nil)
}

// Transformer that converts blockquotes starting with [!type] into callouts
type CalloutTransformer struct{}

var calloutRegex = regexp.MustCompile(`^\[!([a-zA-Z-]+)\](?:\s+(.*))?`)

func (t *CalloutTransformer) Transform(node *ast.Document, reader text.Reader, pc parser.Context) {
	source := reader.Source()

	ast.Walk(node, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		bq, ok := n.(*ast.Blockquote)
		if !ok {
			return ast.WalkContinue, nil
		}

		// Check if first child is a paragraph starting with [!type]
		firstChild := bq.FirstChild()
		if firstChild == nil {
			return ast.WalkContinue, nil
		}

		para, ok := firstChild.(*ast.Paragraph)
		if !ok {
			return ast.WalkContinue, nil
		}

		// Get text content of first line
		if para.Lines().Len() == 0 {
			return ast.WalkContinue, nil
		}

		firstLine := para.Lines().At(0)
		lineText := firstLine.Value(source)

		matches := calloutRegex.FindSubmatch(lineText)
		if matches == nil {
			return ast.WalkContinue, nil
		}

		// Create callout node
		callout := &CalloutNode{
			CalloutType: string(matches[1]),
		}
		if len(matches) > 2 {
			callout.Title = string(matches[2])
		}

		// Remove inline nodes belonging to first line (title)
		firstLineEnd := firstLine.Stop
		var toRemove []ast.Node

		for child := para.FirstChild(); child != nil; child = child.NextSibling() {
			// Check if this node belongs to the first line
			shouldRemove := false

			switch n := child.(type) {
			case *ast.Text:
				// Text node: check if its segment is within first line
				if n.Segment.Start < firstLineEnd {
					shouldRemove = true
				}
			default:
				// For other inline nodes (Emphasis, Link, etc.),
				// check if they contain text from the first line
				// by looking at their first text child
				ast.Walk(child, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
					if entering {
						if textNode, ok := n.(*ast.Text); ok {
							if textNode.Segment.Start < firstLineEnd {
								shouldRemove = true
								return ast.WalkStop, nil
							}
						}
					}
					return ast.WalkContinue, nil
				})
			}

			if shouldRemove {
				toRemove = append(toRemove, child)
			} else {
				// Once we hit a node past the first line, stop checking
				break
			}
		}

		// Actually remove the nodes (can't remove while iterating)
		for _, node := range toRemove {
			para.RemoveChild(para, node)
		}

		// If paragraph is now empty, remove it
		if para.Lines().Len() == 0 {
			bq.RemoveChild(bq, para)
		}

		// Move all children from blockquote to callout
		for child := bq.FirstChild(); child != nil; {
			next := child.NextSibling()
			bq.RemoveChild(bq, child)
			callout.AppendChild(callout, child)
			child = next
		}

		// Replace blockquote with callout
		parent := bq.Parent()
		parent.InsertBefore(parent, bq, callout)
		parent.RemoveChild(parent, bq)

		return ast.WalkContinue, nil
	})
}

// Renderer (same as before)
type CalloutHTMLRenderer struct {
	html.Config
}

// Mapping of callout types to <wa-callout> variants and icons
var calloutMapping = map[string]struct {
	Variant string
	Icon    string
}{
	"note":      {"brand", "circle-info"},
	"abstract":  {"neutral", "clipboard"},
	"summary":   {"neutral", "clipboard"},
	"tldr":      {"neutral", "clipboard"},
	"info":      {"brand", "circle-info"},
	"todo":      {"brand", "circle-check"},
	"tip":       {"success", "lightbulb"},
	"hint":      {"success", "lightbulb"},
	"important": {"brand", "star"},
	"success":   {"success", "check"},
	"check":     {"success", "check"},
	"done":      {"success", "check"},
	"question":  {"warning", "circle-question"},
	"help":      {"warning", "circle-question"},
	"faq":       {"warning", "circle-question"},
	"warning":   {"warning", "triangle-exclamation"},
	"caution":   {"warning", "triangle-exclamation"},
	"attention": {"warning", "triangle-exclamation"},
	"failure":   {"danger", "xmark"},
	"fail":      {"danger", "xmark"},
	"missing":   {"danger", "xmark"},
	"danger":    {"danger", "bolt"},
	"error":     {"danger", "bolt"},
	"bug":       {"danger", "bug"},
	"example":   {"neutral", "list"},
	"quote":     {"neutral", "quote-left"},
	"cite":      {"neutral", "quote-left"},
}

func (r *CalloutHTMLRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(KindCallout, r.renderCallout)
}

func (r *CalloutHTMLRenderer) renderCallout(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n := node.(*CalloutNode)

	if entering {
		mapping, ok := calloutMapping[n.CalloutType]
		if !ok {
			mapping = calloutMapping["note"]
		}

		w.WriteString(`<wa-callout variant="`)
		w.WriteString(mapping.Variant)
		w.WriteString(`"><wa-icon slot="icon" name="`)
		w.WriteString(mapping.Icon)
		w.WriteString(`" variant="regular"></wa-icon>`)

		if n.Title != "" {
			w.WriteString("<strong>")
			w.Write(util.EscapeHTML([]byte(n.Title)))
			w.WriteString("</strong><br />")
		}
	} else {
		w.WriteString(`</wa-callout>`)
	}

	return ast.WalkContinue, nil
}

// Extender
type CalloutExtender struct{}

func (c *CalloutExtender) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(
		parser.WithASTTransformers(
			util.Prioritized(&CalloutTransformer{}, 100),
		),
	)
	m.Renderer().AddOptions(
		renderer.WithNodeRenderers(
			util.Prioritized(&CalloutHTMLRenderer{}, 100),
		),
	)
}
