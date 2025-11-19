# TODO

## Stage 1: MVP

The MVP is just a Core + CLI that can convert gfm to HTML, but also supports syntax highlighting, mathml, hashtags, footnote, callouts.
Mathml is fine for now with my PR merged into treeblood. In the future, when encountering problems, we might want to use a more stable Rust implementation and call it with FFI.

- [ ] Custom Goldmark module to render callout as `<wa-callout>` (or, dynamically, and add to Goldmark readme)
- [ ] Hashtags just do `class=hastag` yucky yuck: either make one that renders `wa-tag size="small" pill>`, or also make dynamic and add to Goldmark readme.
- [ ] Test footnote
- [ ] Test quoteblock
- [ ] Test anchor (https://github.com/abhinav/goldmark-anchor); copy anchor visuals from web awesome documentation
- [ ] Let's just wrap render output with a template which includes webawesome import and my css file, and just a single `<main>` with a max-width or whatever so its readable. Can still update later, but without wa-page
- [ ] table of contents in the right-hand aside. https://github.com/abhinav/goldmark-toc could be used. Also copy javascript to highlight current heading from web awesome documentation.
  - apply minDepth to skip first heading? Does this skip h2 if no h1
  - We need to change the ast.Node tree returned by toc.RenderList to apply `data-level:wq`
- [x] fix FOUC and FOUCE
- [ ] Make page layout adaptive



## Stage 2

Now we must make it reactive to local file system changes. That means making a server that ships a little bit of JS to the doc,
so that it can rehydrate whenever the server wants it to. Ideally this is very lightweight.

All we really need is for the browser to react to a "refresh" signal. We could just have it show some `/tmp/blablapreview.html` file.
The refresh signal is just some new, or the same, `/tmp/` file. We are probably best served implementing this with SSE

## Stage 3

Now make the full application. with caching and metadata and login and page sharing and link archival and blabla


***






- every filename must be unique within a vault.
    - Path is considered metadata?
    - In case of duplicates, and thus, ambiguity, probably throw warnings. Which one will be rendered? undefined behavior.
- every filename must only contain alphanumerics `[0-9a-zA-Z]` and the special characters `-_.+`. Thus, a subset of the [URL specification](https://www.rfc-editor.org/rfc/rfc1738.txt).
- All markdown files must end with their respective extension, `.md`.
- Wikilinks allowed, syntax [[filename|alttext]] and embedded ![[my_attachment]]





md.fnle.be/dir/filename?decoration=false
md.fnle.be/dir/filename?decoration=false&format=raw
This is the default, then -> md.fnle./be/dir/filename?decoration=true&format=md
md.fnle.be/login
md.fnle.be/share/2349jlksjfklsdjjb987sdf79sd7f7s98df987sdf987



