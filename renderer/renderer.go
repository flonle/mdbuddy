package renderer

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"io"
	"log"
	"strings"

	"github.com/flonle/mdbuddy/assets"

	chroma "github.com/alecthomas/chroma/v2"
	treeblood "github.com/wyatt915/goldmark-treeblood"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	// ast "github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	gmHtml "github.com/yuin/goldmark/renderer/html"
	gmText "github.com/yuin/goldmark/text"
	"go.abhg.dev/goldmark/hashtag"
	toc "go.abhg.dev/goldmark/toc"
	"go.abhg.dev/goldmark/wikilink"
)

var myStyle = chroma.MustNewStyle(
	"catpuccin-frappe-without-bg",
	chroma.StyleEntries{
		chroma.Background:               "", // <- !
		chroma.CodeLine:                 "#c6d0f5",
		chroma.Error:                    "#e78284",
		chroma.Other:                    "#c6d0f5",
		chroma.LineTableTD:              "",
		chroma.LineTable:                "",
		chroma.LineHighlight:            "bg:#51576d",
		chroma.LineNumbersTable:         "#838ba7",
		chroma.LineNumbers:              "#838ba7",
		chroma.Keyword:                  "#ca9ee6",
		chroma.KeywordReserved:          "#ca9ee6",
		chroma.KeywordPseudo:            "#ca9ee6",
		chroma.KeywordConstant:          "#ef9f76",
		chroma.KeywordDeclaration:       "#e78284",
		chroma.KeywordNamespace:         "#81c8be",
		chroma.KeywordType:              "#e78284",
		chroma.Name:                     "#c6d0f5",
		chroma.NameClass:                "#e5c890",
		chroma.NameConstant:             "#e5c890",
		chroma.NameDecorator:            "bold #8caaee",
		chroma.NameEntity:               "#81c8be",
		chroma.NameException:            "#ef9f76",
		chroma.NameFunction:             "#8caaee",
		chroma.NameFunctionMagic:        "#8caaee",
		chroma.NameLabel:                "#99d1db",
		chroma.NameNamespace:            "#ef9f76",
		chroma.NameProperty:             "#ef9f76",
		chroma.NameTag:                  "#ca9ee6",
		chroma.NameVariable:             "#f2d5cf",
		chroma.NameVariableClass:        "#f2d5cf",
		chroma.NameVariableGlobal:       "#f2d5cf",
		chroma.NameVariableInstance:     "#f2d5cf",
		chroma.NameVariableMagic:        "#f2d5cf",
		chroma.NameAttribute:            "#8caaee",
		chroma.NameBuiltin:              "#99d1db",
		chroma.NameBuiltinPseudo:        "#99d1db",
		chroma.NameOther:                "#c6d0f5",
		chroma.Literal:                  "#c6d0f5",
		chroma.LiteralDate:              "#c6d0f5",
		chroma.LiteralString:            "#a6d189",
		chroma.LiteralStringChar:        "#a6d189",
		chroma.LiteralStringSingle:      "#a6d189",
		chroma.LiteralStringDouble:      "#a6d189",
		chroma.LiteralStringBacktick:    "#a6d189",
		chroma.LiteralStringOther:       "#a6d189",
		chroma.LiteralStringSymbol:      "#a6d189",
		chroma.LiteralStringInterpol:    "#a6d189",
		chroma.LiteralStringAffix:       "#e78284",
		chroma.LiteralStringDelimiter:   "#8caaee",
		chroma.LiteralStringEscape:      "#8caaee",
		chroma.LiteralStringRegex:       "#81c8be",
		chroma.LiteralStringDoc:         "#737994",
		chroma.LiteralStringHeredoc:     "#737994",
		chroma.LiteralNumber:            "#ef9f76",
		chroma.LiteralNumberBin:         "#ef9f76",
		chroma.LiteralNumberHex:         "#ef9f76",
		chroma.LiteralNumberInteger:     "#ef9f76",
		chroma.LiteralNumberFloat:       "#ef9f76",
		chroma.LiteralNumberIntegerLong: "#ef9f76",
		chroma.LiteralNumberOct:         "#ef9f76",
		chroma.Operator:                 "bold #99d1db",
		chroma.OperatorWord:             "bold #99d1db",
		chroma.Comment:                  "italic #737994",
		chroma.CommentSingle:            "italic #737994",
		chroma.CommentMultiline:         "italic #737994",
		chroma.CommentSpecial:           "italic #737994",
		chroma.CommentHashbang:          "italic #626880",
		chroma.CommentPreproc:           "italic #737994",
		chroma.CommentPreprocFile:       "bold #737994",
		chroma.Generic:                  "#c6d0f5",
		chroma.GenericInserted:          "bg:#414559 #a6d189",
		chroma.GenericDeleted:           "bg:#414559 #e78284",
		chroma.GenericEmph:              "italic #c6d0f5",
		chroma.GenericStrong:            "bold #c6d0f5",
		chroma.GenericUnderline:         "underline #c6d0f5",
		chroma.GenericHeading:           "bold #ef9f76",
		chroma.GenericSubheading:        "bold #ef9f76",
		chroma.GenericOutput:            "#c6d0f5",
		chroma.GenericPrompt:            "#c6d0f5",
		chroma.GenericError:             "#e78284",
		chroma.GenericTraceback:         "#e78284",
	},
)

func RenderTest(input []byte, output io.Writer) (bytes.Buffer, error) {
	var x bytes.Buffer
	fmt.Println("test")
	return x, nil
}

// g-o:embed templates/*.html
// var tmplFS embed.FS
var tmpl = template.Must(template.ParseFS(assets.FS, "templates/*.html"))

//g-o:embed css/*.css
// var cssFS embed.FS

type StandaloneNotePage struct {
	Content template.HTML // Main Content
	TOC     template.HTML // Table Of Contents
	CSS     template.CSS
	JS      template.JS
}

// TODO add input; what mode? file object?
func Render(input []byte, output io.Writer) error {
	// Goldmark
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.Footnote,
			// mathjax.MathJax,
			treeblood.MathML(),
			highlighting.NewHighlighting(
				highlighting.WithCustomStyle(myStyle),
			),
			&wikilink.Extender{},
			&hashtag.Extender{},
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			gmHtml.WithHardWraps(),
			gmHtml.WithXHTML(),
			gmHtml.WithUnsafe(),
		),
	)

	// Render note
	noteRootNode := md.Parser().Parse(gmText.NewReader(input))
	var noteBuf bytes.Buffer
	md.Renderer().Render(&noteBuf, input, noteRootNode)

	// Render TOC
	tocTree, err := toc.Inspect(noteRootNode, input, toc.Compact(true))
	if err != nil {
		return fmt.Errorf("failed to create table of contents")
	}
	// fmt.Fprintln(os.Stderr, string(tocTree.Items[0].Items[0].Items[1].Title))
	var tocBuf bytes.Buffer
	md.Renderer().Render(&tocBuf, input, toc.RenderList(tocTree))
	// fmt.Fprintln(os.Stderr, tocBuf.String())

	// hashtags := make(map[string]struct{})
	// ast.Walk(doc, func(node ast.Node, enter bool) (ast.WalkStatus, error) {
	// 	fmt.Println(node.Type(), node.Kind())
	// 	if node.Kind() == treeblood.KindMathBlock {
	// 		fmt.Println("Block!")
	// 	}
	// 	if node.Kind() == treeblood.KindMathInline {
	// 		fmt.Println("Inline!")
	// 	}

	// 	if n, ok := node.(*hashtag.Node); ok && enter {
	// 		hashtags[string(n.Tag)] = struct{}{}
	// 	}
	// 	return ast.WalkContinue, nil
	// })
	// fmt.Println(hashtags)

	// Fetch and concatenate all CSS files
	// cssDir, err := cssFS.ReadDir("css")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// var cssBuilder strings.Builder
	// for _, file := range cssDir {
	// 	// fmt.Println(file.Name())
	// 	if file.IsDir() {
	// 		continue // we don't traverse nested
	// 	}
	// 	css, err := cssFS.ReadFile("css/" + file.Name())
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	// fmt.Println(css)
	// 	cssBuilder.Write(css)
	// 	cssBuilder.WriteRune('\n')
	// }

	// Render template
	// css, err := cssFS.ReadFile("css/style.css")
	// fmt.Println(string(css))
	// if err != nil {
	// 	log.Fatal(err)
	// }
	note := StandaloneNotePage{
		Content: template.HTML(noteBuf.String()),
		TOC:     template.HTML(tocBuf.String()),
		// CSS:     template.CSS(cssBuilder.String()),
		// JS:      template.JS(""),
	}
	err = tmpl.ExecuteTemplate(output, "note.html", note)
	if err != nil {
		return fmt.Errorf("failed to render template: %w", err)
	}

	return nil
}
