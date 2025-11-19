#dev 

The Bayer-Moore string-search algorithm is an efficient algorithm for searching strings for certain literals.

[[Modelling Downwards Vertical Motion.md]]
[title](<Modelling Downwards Vertical Motion.md>)

> [!important] regex
> Regex search tools like grep will make heavy use of this algorithm. It, and similar tools, will do something called "literal optimizations".
> For example, If the input pattern for a grep search is "test", then grep does not need to invoke the [[Regular Expressions|regex]] engine at all: it can just scan the text corpus for the literal string "test" using Bayer-Moore.
> 
> A lot of regular expressions contain literals in some way, and are thus subject to literal optimizations. A search engine just needs to extract them:
> - `foo|bar` -> detects  `foo` and `bar`
> - `(a|b)c` -> detects `ac` and `bc`
> - `[ab]foo[yz]` -> detects `afooy`, `afooz`, `bfooy`, and `bfooz`
> - `(foo{3,6})` -> detects `foofoofoo`
> 
> Once extracted, the program can search for (one of) the literal(s). Only when it detects a match, will it drop down into the regex engine to verify the candidate against the full regular expression.

It preprocesses the string being searched for (the pattern), but not the string being searched (the text). This makes it particularly useful for searching long texts for a short pattern. In general, the algorithm runs faster as the pattern length increases.

## Mechanics

- It matches on the tail of the pattern rather than the head
- It skips along the text in jumps of multiple characters, instead of iterating over every character

First, you align your pattern to the start of the text.
If your pattern has a length of 5, you start by checking the fifth character in the text (the tail). If it is not a match, *you can safely ignore the first 4 characters*, because they simply cannot be part of a complete match with the pattern.

If a character does not match *any* of the characters in the pattern, we can skip ahead by $m$ characters, where $m$ is the length of the pattern. 

If a character in the text *is* in the pattern, then a partial shift of the pattern along the text is done, as to line up the matching character.

```
A N P A N M A N
---------------
P A N - - - - -  <- P in in pan, shift by 2
- P A N - - - -
- - P A N - - -  <- match! & shift by 1
- - - P A N - -  <- m is not in pan, shift by 3 (not enough characters left)
- - - - P A N -
- - - - - P A N
```

> Alignments of pattern "PAN" to text "ANPANMAN".
> Alignments without an arrow (`<-`) are skipped by the Bayer-Moore algorithm.

The shift rules are implemented as constant-time table lookups, using tables generated during the preprocessing of the pattern.

### Shift Rules

#### The Bad-character Rule

The most obvious. If a character did not match, but it is part of the pattern; shift the pattern so that the the text characters aligns with the relevant pattern character.
if a character did not match, and is not part of the pattern; shift the entire pattern past the point of mismatch.

#### The Good-suffix Rule

Markedly more complex than the bad-character rule. Omitted from the document, for now.

### Example Shift Table

```
Index| Mismatch | Shift
-----+----------+------
 0   |         N|   1    
 1   |        AN|   8    
 2   |       MAN|   3    
 3   |      NMAN|   6   
 4   |     ANMAN|   6   
 5   |    PANMAN|   6  
 6   |   NPANMAN|   6  
 7   |  ANPANMAN|   6
```

> A shift table as a result of preprocessing the pattern "ANPANMAN". Both the bad-character rule as well as the good-suffix rule were used to generate it.

