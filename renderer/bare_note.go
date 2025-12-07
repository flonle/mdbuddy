package renderer

import (
	"bytes"
	"fmt"
	"html/template"
	"io"

	"github.com/flonle/mdbuddy/assets"
	customExtensions "github.com/flonle/mdbuddy/renderer/goldmark-extensions"

	treeblood "github.com/wyatt915/goldmark-treeblood"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	gmHtml "github.com/yuin/goldmark/renderer/html"
	gmText "github.com/yuin/goldmark/text"
	anchor "go.abhg.dev/goldmark/anchor"
	toc "go.abhg.dev/goldmark/toc"
	"go.abhg.dev/goldmark/wikilink"
)

var tmpl = template.Must(template.ParseFS(assets.FS, "static/templates/*.html"))

type BareNotePage struct {
	Title   string
	Content template.HTML // Main Content
	TOC     template.HTML // Table Of Contents
	CSS     template.CSS
	JS      template.JS
}

func RenderBareNote(input []byte, output io.Writer) error {
	// Goldmark
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.Footnote,
			treeblood.MathML(),
			highlighting.NewHighlighting(
				highlighting.WithCustomStyle(catpuccinFrappeNoBg),
			),
			&wikilink.Extender{},
			&customExtensions.HashtagExtension{},
			&anchor.Extender{
				Texter: anchor.Text("#"),
			},
			&customExtensions.CalloutExtender{},
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
		return fmt.Errorf("failed to create table of contents: %w", err)
	}
	var tocBuf bytes.Buffer
	md.Renderer().Render(&tocBuf, input, toc.RenderList(tocTree))

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
	css, err := assets.FS.ReadFile("static/css/bare_note_layout.css")
	if err != nil {
		panic(err)
	}
	tocCSS, err := assets.FS.ReadFile("static/css/table_of_contents.css")
	if err != nil {
		panic(err)
	}
	tocJS, err := assets.FS.ReadFile("static/js/table_of_contents.js")
	if err != nil {
		panic(err)
	}
	sseRefreshJS, err := assets.FS.ReadFile("static/js/sse_refresh.js") // TODO make this conditional ofcourse, we don't always want this
	if err != nil {
		panic(err)
	}
	note := BareNotePage{
		Title:   "My Note",
		Content: template.HTML(noteBuf.String()),
		TOC:     template.HTML(tocBuf.String()),
		CSS:     template.CSS(append(css, tocCSS...)),
		JS:      template.JS(append(tocJS, sseRefreshJS...)),
	}
	err = tmpl.ExecuteTemplate(output, "bare_note.html", note)
	if err != nil {
		return fmt.Errorf("failed to render template: %w", err)
	}

	return nil
}
