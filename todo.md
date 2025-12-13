# TODO

- [ ] Fix checkboxes centering
- [x] Custom Goldmark module to render callout as `<wa-callout>` (or, dynamically, and add to Goldmark readme)
- [x] Hashtags just do `class=hastag` yucky yuck: either make one that renders `wa-tag size="small" pill>`, or also make dynamic and add to Goldmark readme.
- [x] Test footnote
- [x] Test quoteblock
- [x] Test anchor (https://github.com/abhinav/goldmark-anchor); copy anchor visuals from web awesome documentation
- [x] Let's just wrap render output with a template which includes webawesome import and my css file, and just a single `<main>` with a max-width or whatever so its readable. Can still update later, but without wa-page
- [x] table of contents in the right-hand aside. https://github.com/abhinav/goldmark-toc could be used. Also make JS to highlight current visible heading
- [x] fix FOUC and FOUCE
- [x] Make page layout adaptive
- Add styling to override native.css for the following:
  - [x] Links: why so ugly? Maybe with a class?
  - [x] blockquote padding needs to be `var(--wa-space-xs)` or `-s` instead of `-xl`.
  - [x] Maybe also remove weight font and font size from blockquote. Low priority.
- [x] is class scroll-content even defined anywhere?
- [ ] Only the first callout in a file works as intended, all others are skipped?
- [ ] contents of code blocks can exceed my 80ch horizontal limit on `div class='main-content'` and I don't know why (text just goes off into the distance)
- [ ] MathML with Treeblood seems to fail fucking horribly all the time
- [ ] Migration scripts for the vault. This time; with clearly defined rules around syntax & structure. I already kinda started this at the bottom of this file.
- [ ] Git hooks / CI pipeline that enforces certain invariants, like vault/repo uniqueness of filenames, and valid filenames, and maybe even that all wikilinks are valid.
- [ ] Make the 'render' command produce an actually standalone HTML file (by distributing the webawesome components and css myself)
- [ ] Crashes when no headings in file
- [ ] Crashed when you create a new file in a watched directory

## Next Stage

Now make the full application. with caching and metadata and login and page sharing and link archival and blabla

- Adds tags metadata
- Add "edit this page on Gitea" or whatever
- "x-minute read" somewhere automatically?
- (last) publish date somewhere?
- Light/dark mode switch (just switch class `wa-dark`/`wa-light` on `<html>` tag) (catpuccin frappe doesn't work very well in light mode)


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



