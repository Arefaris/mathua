# Reclaiming the dropped figures

This is the work plan for the figures the lesson ingestion dropped while keeping their captions.
It is an **analytical-tier** doc (AGENTS.md section 6), so `scripts/check_ste.py` does not cover
it. The measurement belongs to `scripts/audit_lesson_media.py`, which writes
`data/lessons/figure-audit.json` and the report in `docs/lesson-media.md`; this file turns the
`certain` tier of that report into a plan.

## 1. What is missing

244 paragraphs across 100 lessons read as a caption whose image was dropped:
140 certain, 33 strong and 71 likely. Every one sits
in `data/lessons/teaching/`, and until the first reclaim below those lessons contained no image at
all, so the caption was the only remaining trace of the figure. The `certain` tier opens with an
auto-numbered `Figure 3.2` label that the PreTeXt conversion emits for a figure and keeps when it
drops the image, so it needs no judgement.

## 2. The upstream is OpenStax, not Algebrica

Every teaching lesson opens with a credit line naming its source, for example:

```
> Content sourced from [OpenStax Prealgebra](https://openstax.org/books/prealgebra/pages/1-introduction) by Marecek & Anthony-Smith - CC BY 4.0
```

Two books account for every `certain` caption: **OpenStax Calculus Volume 1** (69 figures) and
**OpenStax Prealgebra** (33). Both are CC BY 4.0, which permits reuse and adaptation with
attribution, and the attribution is already the lesson's own header line.

That changes what "a direct look-alike" means for these figures. Their visual target is their own
upstream, not Algebrica. The Algebrica dialect in `docs/diagram-style.md` applies when a figure is
redrawn by hand; a figure re-fetched from OpenStax keeps the OpenStax look and does not need to
match Algebrica.

## 3. The captions are 140 rows and 70 distinct figures

Teaching lessons share source sections, so one figure appears in several concept lessons: the same
`Figure 1.13` is the caption in five of them, and `Figure 1.11` in six. Deduplicating on the pair
(figure number, caption) leaves **70 distinct figures** still to reclaim. Section 6 gives
the method.

## 4. Dispositions

| Disposition | Figures | Action |
| --- | --- | --- |
| schematic, redraw | 21 | Draw in the Algebrica cool dialect from section 4 of the style doc. Graphs, number lines, regions and solids are fully determined by their captions, so no upstream fetch is needed and the result is Mathua's own work. |
| illustration, fetch | 44 | Fetch the OpenStax original. The license is CC BY 4.0 and the credit line is already in the lesson header. |
| unidentified | 0 | Empty since 2026-10-16: every bare caption was classified by opening its upstream section and looking at the image, and each of the seventeen is a `schematic, redraw` or an `illustration, fetch` now. |
| photo, third-party credit | 5 | The caption names a photographer or agency (`credit: modification of work by Steve Jurvetson, Wikimedia Commons` and similar). These cannot be redrawn; each needs a licensed fetch or a replacement. |

**One correction the first batch forced.** Figures 1.7, 1.8 and 1.9 were assigned
`illustration, fetch` when this table was drawn, because the assignment read the caption and
nothing else: *76 is closer to 80 than to 70* describes a picture of a number line, and a picture
of a number line is an illustration. Opening Prealgebra 1.1 shows all three are the same 70..80
number line as Figure 1.2 under a different mark, which makes them schematics. They were redrawn
(section 7), so the redraw row gained three figures and then lost the same three, and the fetch
row lost three: when that batch landed the table above read 43 redraw and 37 fetch, where reading
the caption alone had suggested 35 and 38, and the classification below moved those counts on
again to 53 and 44, of which section 7 has since landed thirty-two.

**The column is a guess until someone opens upstream, and the caption is the wrong thing to
read.** The reliable text is the upstream `alt` attribute, which describes the picture rather than
the argument the picture supports: Figure 1.7's caption is a sentence about rounding, while its
alt is *An image of a number line from 70 to 80 with increments of one*. That is how Figure 1.6 was
placed here as a `schematic, redraw` even though its caption is a bare `Figure 1.6` — Prealgebra
1.1 gives it an alt of *A chart titled 'Place value' with fifteen columns and 4 rows* — and it is
how the seventeen `unidentified` rows were classified on 2026-10-16.

**The seventeen were read twice, and the first reading was not enough.** The disposition is a
reading of the description, so the description has to be the whole one: the scrape behind this plan
truncated every alt at 110 characters, and the 23 upstream sections were re-fetched to read
`data-alt` in full (Figure 6.34's is 1,244 characters). The text alone still does not settle it,
because an alt describes a picture in words and the words can hide the shape: Figure 1.4 opens *An
image with three items* and the three items are a cube, a rod and a flat, which is what the lesson
calls base-10 blocks, and nothing in *An image of 28 cookies placed at random* says whether that is
a drawing or a photograph. So each of the seventeen images was downloaded and looked at, and the
disposition below is the reading of the picture with the alt as its check.

**Where the line falls: ten redraw, seven fetch.** A figure is a schematic when its content is a
number line, a graph, a geometric block, tile, circle or region, or a table of words and numbers —
a redraw of that is faithful and is Mathua's own work. It is an illustration when the subject is a
photograph (three stacks of dollar bills, 24 pennies, 28 cookies in three of the figures, two
pizzas) or when it is a table whose cells are formulas and small plots (6.34), which an SVG cannot
typeset. Two calls are worth recording because a later reader will meet them again: Figure 5.7 is a
calendar page with seven distances written on it, which is a table of words and numbers and is
therefore a redraw even though the upstream image is a picture of a calendar; and 1.3 is a
photograph of currency with no credit line, so it is `illustration, fetch` and not the
`photo, third-party credit` row, which exists for a third party whose name must travel with the
image rather than for every photograph.

**What lands with them.** All seventeen keep `(none)` in the Caption column, because upstream's
`figcaption` holds only the label for these figures — the alt is the only description there is. It
is also what a landed figure will carry: the rule that a lesson figure's alt is its caption's
words, never a paraphrase, has nothing to copy here, and the same document's alt is still the
author's words.

**Ten of the figures below were invisible until 2026-10-15**, because of a defect in this plan's own
measurement rather than in the corpus. `section_a_candidates` skipped any standalone-italic
paragraph whose neighbour started with `-`, `*`, `|` or `>`, which is right for a line inside a
list and wrong for a caption that simply sits above one: Figure 1.5 is above the place-value table,
Figure 1.6 above the bulleted list of digit values, Figures 4.50 and 4.51 are two captions in a
row. Exempting auto-numbered labels from that test added 24 rows and these ten figures — 1.5, 1.6,
3.20, 4.40, 4.50, 4.51, 4.52, 4.58, 4.59 and 5.3 — every one of which rates `certain` and none of
which is a list item. The tier is 220 rows now rather than 196, which is a fix and not a
regression: the same corpus, counted without the exclusion.

## 5. The captions were damaged too, and are repaired

The caption text was ingested with `figcaption.get_text(' ', strip=True)`. OpenStax writes a
mathematical expression as `<math><semantics><mrow>…</mrow><annotation-xml>…</annotation-xml>
</semantics></math>`: the expression is present **twice**, and the flattening joins every token
with a space. So a caption arrived as `… the number 138 . 138 .` — the expression, then the
expression again. The fault is in `scripts/ingest/ingest_openstax.py` and not in the document.
A figure restored under a caption like that still reads as broken, so repairing the caption was
part of the work rather than a separate cleanup.

**Scale.** The corpus has 244 lines that carry a figure caption; 160 of them come from a book the
modern ingester fetches. Of those 160, 119 were repaired, 40 already matched the pipeline's output,
and 1 was held back for a hand. The repair landed on 2026-10-15, and the re-run is idempotent:
159 already correct, 0 to repair, 1 held back.

**The repair took the text from upstream**, which is the only reliable source for it: each caption
was regenerated through the fixed ingester, so it arrives in the same form as the lesson body
around it, and then passed through the corpus `dollar-number` rule
(`scripts/lint_math_layout.py`) that bare numerals are upright text rather than inline math. That
rule is why the repaired captions read `… worth 35 cents, or $\frac{35}{100}$ of a dollar.` and not
`… worth $35$ cents`. The test that the repair is complete is a comparison against upstream, not a
pattern: the numbers alone were never the whole damage, and letters, expressions, percentages and
dollar amounts were doubled in the same way.

Two things the repair had to know:

- **A bare caption is not damage.** `Figure 1.3`, `1.4` and `1.6` have no caption upstream either —
  their `figcaption` holds only the label — so they are reclaim targets rather than repair targets,
  and section 4 records what each of them turned out to be.
- **One caption needed a hand.** `Figure 2.18` carries a piecewise function as an inline MathML
  table. Pandoc renders it with literal newlines, and the line-scoped dollar balancing then escapes
  the span's opener, so the ingester cannot emit it correctly today; the text in the lesson is the
  pipeline's intended output written on one line.

**The guard** is rule 16 in `scripts/audit_lessons.py`: a caption may not contain an immediately
repeated token run (`doubled_caption_run`). Measured against the damaged corpus it catches 63 of the
67 distinct damaged captions and fires on none of the repaired ones, and it has no baseline, so
re-ingesting with the old figure handler fails the audit. Rule 16 covers captions only, because a
corpus-wide doubled-token test is not usable: it reports 1009 legitimately repeated pairs, mostly
`& &` inside math and paired table cells.

**One measurement it moved.** The repaired captions that gained inner `*x*` emphasis — OpenStax
marks a variable with `<em>` — stopped matching `audit_lesson_media.py`'s standalone-emphasis test,
which removed eight imageless captions from this plan. The test now allows a balanced inner span
and still rejects a thematic break (`* * *`); the report stood at 324 rows and 220 in the
certain tier on 2026-10-15, after two reclaims in section 7 had removed eighteen rows and the label
exemption in section 4 had restored twenty-four. Section 7 records what each batch since has
removed, and section 1 states what is left.

**One glyph normalisation, in the same text.** OpenStax marks the parts of a two-part example with
circled letters (`ⓐ` … `ⓔ`), and the ingester copied the codepoint through. At lesson size `ⓐ` is an
`a` inside a circle — which is what `@` is — so a caption reads as if an email sign preceded it, and
no theme font carries the codepoint, so the browser substitutes a face at a different size and
baseline. The corpus already spelled the parts out in prose (`part (a)`, 22 times, next to a marker
that said `ⓐ`), so all 2340 markers across 72 lessons, several of them inside captions, were written
as `(a)` … `(e)` on 2026-10-15. AGENTS.md section 1 records the reasoning and rule 17 of
`scripts/audit_lessons.py` enforces it; the captions in section 6 are the repaired text with that
normalisation applied.

**Not repaired, and a different fault.** `trig.period` and `trig.radians` are the only two lessons
that came from the retired hand-rolled converter (`scripts/extract_openstax.py`), and their text is
mangled rather than merely doubled: `_underline Send underline_ine` stands for `Sine`,
`sin(t)= opposite hypotenuse` stands for $\sin t = \frac{\text{opposite}}{\text{hypotenuse}}$, and
every expression appears twice. The modern pipeline renders both sections correctly, but it produces
602 and 1071 lines against the committed 280 and 558: the committed files are curated, with back
matter stripped and the wiki-lesson fixes applied, so re-ingesting them would replace the lessons
rather than repair them. That is a content decision and this record leaves it open.

## 6. The figures

Sorted by disposition, then figure number compared as text, which is why 1.12 comes before 1.5 and
2.10 before 2.6. `Lessons` is how many concept lessons carry the
caption; one figure is restored once and every lesson that cites it gains the image. The caption is
the text as repaired in section 5, truncated at 85 characters, so it is what the figure will carry
when it lands. A figure that has landed leaves this table and is recorded in section 7 instead, so
every count here is work still to do.

| Figure | Caption (repaired) | Book | Section | Lessons | Disposition |
| --- | --- | --- | --- | --- | --- |
| 4.3 | An airplane is flying at a constant height of 4000 ft. The distance between the perso... | Calculus Volume 1 | 4 Related Rates | 1 | schematic, redraw |
| 4.53 | The graph of this rational function approaches a horizontal asymptote as $x\rightarro... | Calculus Volume 1 | 4 Limits at Infinity and Asymptotes | 1 | schematic, redraw |
| 4.54 | The graph of this rational function approaches the horizontal asymptote $y = 0$ as $x... | Calculus Volume 1 | 4 Limits at Infinity and Asymptotes | 1 | schematic, redraw |
| 4.56 | The graph of the rational function $f(x) = \left( {3x^{2} + 4x} \right)\text{/}\left(... | Calculus Volume 1 | 4 Limits at Infinity and Asymptotes | 1 | schematic, redraw |
| 4.58 | The function $f(x) = \text{sin} x$ oscillates between $1\ \text{and}\ -1$ as $x\right... | Calculus Volume 1 | 4 Limits at Infinity and Asymptotes | 1 | schematic, redraw |
| 4.59 | The function $f(x) = \text{tan} x$ does not approach a limit and does not approach $\... | Calculus Volume 1 | 4 Limits at Infinity and Asymptotes | 1 | schematic, redraw |
| 5.32 | The graph shows speed versus time for the given motion of a car. | Calculus Volume 1 | 5 Integration Formulas and the Net Change Theorem | 1 | schematic, redraw |
| 5.33 | The graph shows velocity versus time for a particle moving with a linear velocity fun... | Calculus Volume 1 | 5 Integration Formulas and the Net Change Theorem | 1 | schematic, redraw |
| 5.35 | Graph (a) shows the positive area between the curve and the *x*-axis, whereas graph (... | Calculus Volume 1 | 5 Integration Formulas and the Net Change Theorem | 1 | schematic, redraw |
| 5.36 | The graph shows areas between a curve and the *x*-axis for an odd function. | Calculus Volume 1 | 5 Integration Formulas and the Net Change Theorem | 1 | schematic, redraw |
| 6.25 | (a) A region bounded by the graph of a function of $x.$ (b) The solid of revolution f... | Calculus Volume 1 | 6 Volumes of Revolution: Cylindrical Shells | 1 | schematic, redraw |
| 6.26 | (a) A representative rectangle. (b) When this rectangle is revolved around the $y\tex... | Calculus Volume 1 | 6 Volumes of Revolution: Cylindrical Shells | 1 | schematic, redraw |
| 6.27 | Calculating the volume of the shell. | Calculus Volume 1 | 6 Volumes of Revolution: Cylindrical Shells | 1 | schematic, redraw |
| 6.28 | (a) Make a vertical cut in a representative shell. (b) Open the shell up to form a fl... | Calculus Volume 1 | 6 Volumes of Revolution: Cylindrical Shells | 1 | schematic, redraw |
| 6.29 | (a) The region $R$ under the graph of $f(x) = 1\text{/}x$ over the interval $\left\lb... | Calculus Volume 1 | 6 Volumes of Revolution: Cylindrical Shells | 1 | schematic, redraw |
| 6.30 | (a) The region $R$ under the graph of $f(x) = 2x - x^{2}$ over the interval $\left\lb... | Calculus Volume 1 | 6 Volumes of Revolution: Cylindrical Shells | 1 | schematic, redraw |
| 6.31 | (a) The region $Q$ to the left of the function $g(y)$ over the interval $\left\lbrack... | Calculus Volume 1 | 6 Volumes of Revolution: Cylindrical Shells | 1 | schematic, redraw |
| 6.32 | (a) The region $R$ between the graph of $f(x)$ and the $x\text{-axis}$ over the inter... | Calculus Volume 1 | 6 Volumes of Revolution: Cylindrical Shells | 1 | schematic, redraw |
| 6.33 | (a) The region $R$ between the graph of $f(x)$ and the graph of $g(x)$ over the inter... | Calculus Volume 1 | 6 Volumes of Revolution: Cylindrical Shells | 1 | schematic, redraw |
| 6.35 | (a) The region $R$ bounded by two lines and the $x\text{-axis}.$ (b) The solid of rev... | Calculus Volume 1 | 6 Volumes of Revolution: Cylindrical Shells | 1 | schematic, redraw |
| 6.36 | (a) The region $R$ between the curve and the $x\text{-axis}.$ (b) The solid of revolu... | Calculus Volume 1 | 6 Volumes of Revolution: Cylindrical Shells | 1 | schematic, redraw |
| 1.11 | (none) | Prealgebra | 1 Multiply Whole Numbers | 6 | illustration, fetch |
| 1.13 | (none) | Prealgebra | 1 Divide Whole Numbers | 5 | illustration, fetch |
| 1.14 | (none) | Prealgebra | 1 Divide Whole Numbers | 5 | illustration, fetch |
| 1.15 | (none) | Prealgebra | 1 Divide Whole Numbers | 5 | illustration, fetch |
| 1.3 | (none) | Prealgebra | 1 Introduction to Whole Numbers | 6 | illustration, fetch |
| 2.10 | Factors of the counting numbers from 2 through 20, with prime numbers highlighted | Prealgebra | 2 Find Multiples and Factors | 4 | illustration, fetch |
| 2.20 | The function $f(x) = {1\text{/}\left( {x - a} \right)^{n}}$ has infinite limits at *a... | Calculus Volume 1 | 2 The Limit of a Function | 1 | illustration, fetch |
| 2.6 | Multiples of 2 between 1 and 50 | Prealgebra | 2 Find Multiples and Factors | 4 | illustration, fetch |
| 2.7 | Multiples of 5 between 1 and 50 | Prealgebra | 2 Find Multiples and Factors | 4 | illustration, fetch |
| 2.8 | Multiples of 10 between 1 and 50 | Prealgebra | 2 Find Multiples and Factors | 4 | illustration, fetch |
| 2.9 | Multiples of 3 between 1 and 50 | Prealgebra | 2 Find Multiples and Factors | 4 | illustration, fetch |
| 3.18 | The derivative of $g(x)$ is 3 times the derivative of $f(x).$ | Calculus Volume 1 | 3 Differentiation Rules | 3 | illustration, fetch |
| 3.19 | This function has horizontal tangent lines at *x* = 2/3 and *x* = 4. | Calculus Volume 1 | 3 Differentiation Rules | 3 | illustration, fetch |
| 3.20 | The grandstand next to a straightaway of the Circuit de Barcelona-Catalunya race trac... | Calculus Volume 1 | 3 Differentiation Rules | 3 | illustration, fetch |
| 3.21 | (a) One section of the racetrack can be modeled by the function $f(x) = x^{3} + 3x^{2... | Calculus Volume 1 | 3 Differentiation Rules | 3 | illustration, fetch |
| 3.30 | The equation $x^{2} + y^{2} = 25$ defines many functions implicitly. | Calculus Volume 1 | 3 Implicit Differentiation | 1 | illustration, fetch |
| 3.31 | The line $y = \frac{3}{4}x - \frac{25}{4}$ is tangent to $x^{2} + y^{2} = 25$ at the ... | Calculus Volume 1 | 3 Implicit Differentiation | 1 | illustration, fetch |
| 3.32 | Finding the tangent line to the folium of Descartes at $\left( {\frac{3}{2},\frac{3}{... | Calculus Volume 1 | 3 Implicit Differentiation | 1 | illustration, fetch |
| 3.34 | The tangent line to $E(x) = e^{x}$ at $x = 0$ has slope 1. | Calculus Volume 1 | 3 Derivatives of Exponential and Logarithmic Functions | 1 | illustration, fetch |
| 3.35 | $\text{The function}\ y = \text{ln} x$ is increasing on $(0,\text{+}\infty).$ Its der... | Calculus Volume 1 | 3 Derivatives of Exponential and Logarithmic Functions | 1 | illustration, fetch |
| 4.2 | As the balloon is being filled with air, both the radius and the volume are increasin... | Calculus Volume 1 | 4 Related Rates | 1 | illustration, fetch |
| 4.4 | (none) | Prealgebra | 4 Visualize Fractions | 7 | illustration, fetch |
| 4.41 | (a) As $x\rightarrow\infty,$ the values of $f$ are getting arbitrarily close to $L.$ ... | Calculus Volume 1 | 4 Limits at Infinity and Asymptotes | 1 | illustration, fetch |
| 4.43 | This function approaches a horizontal asymptote as $x\rightarrow\text{\pm}\infty.$ | Calculus Volume 1 | 4 Limits at Infinity and Asymptotes | 1 | illustration, fetch |
| 4.44 | This function crosses its horizontal asymptote multiple times. | Calculus Volume 1 | 4 Limits at Infinity and Asymptotes | 1 | illustration, fetch |
| 4.46 | This function has two horizontal asymptotes. | Calculus Volume 1 | 4 Limits at Infinity and Asymptotes | 1 | illustration, fetch |
| 4.47 | For this function, the functional values approach infinity as $x\rightarrow\text{\pm}... | Calculus Volume 1 | 4 Limits at Infinity and Asymptotes | 1 | illustration, fetch |
| 4.48 | For a function with a limit at infinity, for all $x > N,$ $\left\| {f(x) - L} \right\| ... | Calculus Volume 1 | 4 Limits at Infinity and Asymptotes | 1 | illustration, fetch |
| 4.49 | For a function with an infinite limit at infinity, for all $x > N,$ $f(x) > M.$ | Calculus Volume 1 | 4 Limits at Infinity and Asymptotes | 1 | illustration, fetch |
| 4.5 | A camera is positioned 5000 ft from the launch pad of the rocket. The height of the r... | Calculus Volume 1 | 4 Related Rates | 1 | illustration, fetch |
| 4.55 | As $x\rightarrow\infty,$ the values $f(x)\rightarrow\infty.$ As $x\rightarrow\text{-}... | Calculus Volume 1 | 4 Limits at Infinity and Asymptotes | 1 | illustration, fetch |
| 4.57 | This function has two horizontal asymptotes and it crosses one of the asymptotes. | Calculus Volume 1 | 4 Limits at Infinity and Asymptotes | 1 | illustration, fetch |
| 4.6 | Water is draining from a funnel of height 2 ft and radius 1 ft. The height of the wat... | Calculus Volume 1 | 4 Related Rates | 1 | illustration, fetch |
| 4.6 | The U.S. coin called a quarter is worth one-fourth of a dollar. | Prealgebra | 4 Multiply and Divide Fractions | 4 | illustration, fetch |
| 4.60 | The exponential function approaches zero as $x\rightarrow\text{-}\infty$ and approach... | Calculus Volume 1 | 4 Limits at Infinity and Asymptotes | 1 | illustration, fetch |
| 4.61 | The natural logarithm function approaches $\infty$ as $x\rightarrow\infty.$ | Calculus Volume 1 | 4 Limits at Infinity and Asymptotes | 1 | illustration, fetch |
| 4.7 | Together, a quarter and a dime are worth 35 cents, or $\frac{35}{100}$ of a dollar. | Prealgebra | 4 Add and Subtract Fractions with Different Denominators | 3 | illustration, fetch |
| 5.2 | This chart illustrates place values to the left and right of the decimal point. | Prealgebra | 5 Decimals | 3 | illustration, fetch |
| 5.3 | When we write a check, we write the amount as a decimal number as well as in words. T... | Prealgebra | 5 Decimals | 3 | illustration, fetch |
| 5.4 | (a) We see that 2.72 is closer to 3 than to 2. So, 2.72 rounded to the nearest whole ... | Prealgebra | 5 Decimals | 3 | illustration, fetch |
| 6.2 | Among every 100 community college students, 57 are female. | Prealgebra | 6 Understand Percent | 3 | illustration, fetch |
| 6.34 | (none) | Calculus Volume 1 | 6 Volumes of Revolution: Cylindrical Shells | 1 | illustration, fetch |
| 6.6 | A $\text{20\%}$ tip for an $\text{\$80}$ restaurant bill comes out to $\text{\$16}.$ | Prealgebra | 6 Solve General Applications of Percent | 3 | illustration, fetch |
| 6.7 | The sales tax is calculated as a percent of the purchase price. | Prealgebra | 6 Solve Sales Tax, Commission, and Discount Applications | 2 | illustration, fetch |
| 2.22 | (credit: NASA) | Calculus Volume 1 | 2 The Limit of a Function | 1 | photo, third-party credit |
| 4.4 | (credit: modification of work by Steve Jurvetson, Wikimedia Commons) | Calculus Volume 1 | 4 Related Rates | 1 | photo, third-party credit |
| 5.34 | (credit: modification of work by Carter Brown, Flickr) | Calculus Volume 1 | 5 Integration Formulas and the Net Change Theorem | 1 | photo, third-party credit |
| 6.4 | (credit: Riles32807, Wikimedia Commons) | Prealgebra | 6 Understand Percent | 3 | photo, third-party credit |
| 6.8 | Applications of discounts are common in everyday life. (credit: Charleston's TheDigit... | Prealgebra | 6 Solve Sales Tax, Commission, and Discount Applications | 2 | photo, third-party credit |

## 7. Progress

**Figure 1.2** ("The numbers on the number line increase from left to right, and decrease from
right to left") landed on 2026-10-15 as the first reclaim and the proof of the pipeline. It is a
`schematic, redraw`: drawn as `web/next-app/public/diagrams/arithmetic/number-line-counting.svg`,
embedded in all six lessons that cite it (`arith.place.hundreds`, `arith.place.tens`,
`arith.place.thousands`, `arith.round.hundreds`, `arith.round.tens`, `arith.round.thousands`), and
the regenerated audit dropped those six rows from the `certain` tier. It is a lesson figure, so it
needs no `data/diagrams/meta.json` entry (section 5 of the style doc explains why).

**Figures 1.7, 1.8 and 1.9** (the same 70..80 number line under three marks: 76 nearer 80, 72
nearer 70, 75 exactly midway) landed on 2026-10-15 as one batch, in those six lessons. They are
`number-line-round-76.svg`, `-72.svg` and `-75.svg`, beside Figure 1.2's
`number-line-counting.svg` and on its canvas to the pixel — 460x130, the same ink, the same
tick style — because the four sit within a few paragraphs of each other and a figure that does
not match its neighbours is worse than no figure. Each draws the distance to both neighbouring
tens to scale, so the shorter span is the argument for which way the number rounds; the ink weight
on that span's arrow agrees with it, and at 75 both spans are the same length and both are inked.
The alt text is the caption minus its `Figure 1.7 ` label, which is the rule
`scripts/place_lesson_figures.py` follows and the exact text Figure 1.2 carries.

Those three captions held eighteen rows, and the regenerated audit dropped all of them: 318
candidates to 300, and 214 in the certain tier to 196. They were the first batch drawn from a
disposition that had been assigned by reading a caption, and the reading was wrong — section 4
records the correction.

**Figures 1.4, 1.5 and 1.6** (the base-10 vocabulary, then 138 modeled with it, then 5,278,194 in a
place value chart) landed on 2026-10-16 as one batch, in those six lessons, as
`base-ten-blocks.svg`, `base-ten-138.svg` and `place-value-chart.svg`. 1.5 and 1.6 had been
`schematic, redraw` since the plan was written, read out of their upstream `alt`; 1.4 was one of the
seventeen that section 4 classified, and the three are one subject, so they land together. 1.4 and
1.5 share a cell size of 9px, which is what makes a rod ten cubes and a flat ten rods by
construction rather than by three separate drawings that happen to agree.

The upstream 1.4 stacks its three items vertically with each label beside its item, and the lesson
figure is bounded to 180px tall (section 3 of the style doc), so the batch is landscape and each
label sits under its item in two lines instead. That is a transposition of the layout and not of the
content: the same three blocks, the same three labels.

The alt text differs by figure, and the difference is the point. 1.5 has a caption, so its alt is
the caption minus its `Figure 1.5 ` label, the rule Figure 1.2 follows. 1.4 and 1.6 have only the
bare label upstream, so there is nothing of the author's to reuse except the `alt` itself: both
carry upstream's `data-alt` verbatim, which is the same document's accessibility text and not a
paraphrase of it.

Those three captions held eighteen rows, and the regenerated audit dropped all of them: 324
candidates to 306, and 220 in the certain tier to 202. The ceiling for orphan captions in
`internal/lessons/assets_test.go` ratchets down with it, to 306.

**Figure 1.12** (a 2-by-3 rectangle of unit squares, its area 6 square feet) landed on 2026-10-16 in
the six `arith.mult.*` lessons, as `rectangle-area.svg`. It is upstream's layout kept whole: `2 ft`
against the left edge, `3 ft` under the bottom edge, and `2 · 3 = 6 ft²` to the right, which is where
upstream puts it. Upstream's words are set at the family's 13px label size and the blocks are sized
to fill the canvas, so the blocks are larger than upstream's against the same words — a re-drawn
figure cannot carry upstream's text-to-shape ratio without breaking the size the other arithmetic
figures use, and the dial is turned towards the drawing. The cell size is one number, so the
rectangle is six unit squares and the two labeled edges are the drawing's own edges rather than a
coincidence of it. The alt is the caption minus its `Figure 1.12 ` label, the rule Figure 1.2
follows.

These six lessons are the first the plan has put a figure into that had none: the `arith.mult.*`
bodies carried no image at all, where the earlier batches went into lessons that already carried
Figure 1.2 or 1.7-1.9. So this is also the first reclaim whose lesson `lessonFigures.spec.ts` could
not measure before, and the case for it is added with the figure rather than after it.

Those six captions held six rows, and the regenerated audit dropped all of them: 306 rows to 300,
and 202 in the certain tier to 196. The ceiling for orphan captions in
`internal/lessons/assets_test.go` follows, to 300.

**Figures 4.2, 4.3 and 4.5** (a circle cut into thirds; five tiles of one length showing a whole,
halves, thirds, fourths and sixths; and a half above the same rectangle in sixths) landed on
2026-10-16 in the eleven `frac.*` lessons that cite them, as `fraction-circle.svg`,
`fraction-tiles.svg` and `fraction-sixths.svg`. 4.2 and 4.3 are the same seven lessons, where the
two figures sit a few paragraphs apart, so they are drawn to the same rules and the same label
size; 4.5 is in the four `frac.ops.*` and `frac.*.whole` lessons.

These are the first figures in the plan whose labels *are* mathematics. Upstream writes every label
as a stacked fraction — numerator, rule, denominator — and a figure that teaches fractions while
writing `1/3` on one line is showing the answer instead of the notation the lesson is teaching, so
each label here is three primitives, two `<text>` and a `<line>`, on a shared vertical rhythm. There
is no KaTeX inside an SVG file and no `foreignObject` survives the sanitizer, so the fraction is
drawn rather than typeset.

One deliberate departure, and it is worth naming because it is visible. Upstream colors the five
tiles differently — blue, grey, yellow, grey, peach — and the house palette has one fill,
`#e6e6e6`, with ink at `#333333`. A tile is distinguished by how many pieces it is cut into, which
is what the figure is about; color would carry no information the division count does not, so all
five tiles are the same fill. The division count is also what makes the tiles one figure rather than
five: they share a length, so the pieces of the sixths row are visibly half the width of the thirds.

The canvases are sized to their content, which for these three does not mean one width. The circle
is `240x180` — a lone circle in a `460`-wide canvas would be mostly empty panel — and the two bar
figures are `420x180`, the tiles and the rectangles being `360` wide with the same margin the number
lines leave. All three carry the opaque white background rect, so they read on the dark panel.

Their alt text is upstream's `data-alt` verbatim, all three being bare-label captions, and it was
read out of the scrape rather than retyped. That scrape is double-encoded — decoded as latin-1 while
holding UTF-8, so a curly quote arrives as three characters — and the wiring inverts exactly that
error rather than transcribing around it, which is also why Figure 4.3's alt is one 530-character
sentence that no hand copy would have survived intact.

Those eighteen captions held eighteen rows, and the regenerated audit dropped all of them: 300 to
282, and the tier that opens on an auto-numbered label from 196 to 178. The ceiling for orphan
captions in `internal/lessons/assets_test.go` follows, to 282. This batch leaves no lesson in the
`frac.*` set without a figure, since all eleven had none.

**Figures 6.3 and 6.5** (one mechanism shown in both directions: the percents of Table 6.1 as
decimals, and the decimals of Table 6.2 as percents) landed on 2026-10-16 in the three `pct.*`
lessons that cite them, as `percent-to-decimal.svg` and `decimal-to-percent.svg`.

This is the first place in the plan where the upstream is inconsistent with itself rather than
merely silent. The section page references **one image hash for both figures** — the same table — and
that table is the `Percent | Decimal` one with `006.%` through `012.5%`. Figure 6.5's alt names the
columns in the other order and lists five decimals where that image's decimal column has four
values, so the alt and the image cannot both be describing Figure 6.5. The lesson prose decides it:
"Figure 6.5 uses the decimal numbers in Table 6.2 and shows visually to convert them to percents by
moving the decimal point two places to the right". Two tables are drawn, each following its own
figure's alt, because for 6.5 the alt is the only description of what the figure is for. The four
values of Table 6.2 are in it and so is the fifth the alt names, and the prose stays true either way
since it claims the figure uses those numbers rather than only those.

The notation `006.%` is upstream's and not a slip: it writes the percent as a number with its
decimal point visible at the end, so the two hops under it have something to move. The hops are two
small arcs under the right end of the value, where the point arrives; they are not themselves
directional, and the prose says which way the point travels.

Each table is `250` wide, which is the width of its content: two columns of `115` holding labels of
up to six characters at the family's 13px. The header row is filled `#e6e6e6` and the body cells are
white, ruled with guides and framed in ink, which is the same division of structure the place value
chart uses. The white background rect is there as in every local figure, so the tables read on the
dark panel.

Their alt text is upstream's `data-alt` verbatim, read out of the scrape with its encoding inverted.
Those six captions held six rows, and the regenerated audit dropped all of them: 282 to 276, and 172
in the certain tier. The ceiling for orphan captions in `internal/lessons/assets_test.go` follows, to
276.

**Figure 5.5** (a number line from 0 to 1 marked every 0.2, with a measurement bracket over each
step) landed on 2026-10-16 in the four `dec.ops.*` lessons, as `decimal-number-line.svg`. It is the
first reclaim in the plan that reuses a dialect already in the tree rather than starting one: the
`number-line-*.svg` files have drawn the same axis, ticks, arrowheads and 13px labels since the
arithmetic set was placed, so this adds one primitive to that vocabulary, the bracket, and inherits
the rest.

The bracket is where the drawing and upstream disagree about the figure's own text. Upstream's image
braces four of the five steps -- the one from 0.8 to 1 has none -- while its alt says the braces show
"a distance of 0.2 between each adjacent set of 2 numbers", which is five. Five are drawn. That is
the rule section 7 records for Figure 6.5, and the reason is the same: the alt is the description
that ships with the figure, and a reader told about a bracket that is not there has been told
something false. It also happens to be the better figure, since the claim being made is that every
step is worth the same amount.

The bracket itself is one unfilled path: a bar whose ends turn down toward the axis, with a small
peak at the center, drawn in `#636363` because it annotates the axis rather than being part of it.
The five `0.2` labels over it are ink, not annotation, because they are the measurement the figure
exists to make; the six axis numbers stay in the annotation color, as they do in every other number
line here. The canvas is 460 by 130, the size the other number lines use, so the figures that share
lessons with this one line up.

Its alt text is upstream's `data-alt` verbatim, read out of the scrape with its encoding inverted.
Those four captions held four rows, and the regenerated audit dropped all of them: 276 to 272, and
168 in the certain tier. The ceiling for orphan captions in `internal/lessons/assets_test.go`
follows, to 272.

**Figure 5.6** (a circle with the three parts named: a dashed diameter across the widest portion, a
dashed radius from the center to the rim, and the circumference along the edge) landed on 2026-10-16
in the three lessons that cite it, as `circle-parts.svg`. Two decisions are worth naming because
both are places the dialect did not already have an answer.

The dashed lines are ink at `stroke-dasharray="4"`. The style doc ties dashes to guides, and a guide
is `#BABABA`; the Algebrica set also uses an ink dashed line for an annotation that is not
scaffolding, which is what these two are -- they are what the words point at, so drawing them in the
guide color would make the two labeled objects the faintest thing in the figure.

`Circumference` is drawn one glyph at a time: thirteen `<text>` elements, each positioned on a
radius-76 arc under the circle and rotated to that arc's tangent. The word is placed *along* the
edge, and that is what along means. `<textPath>` would do it in one element, but it needs an `href`
to a path `id`, and an `href` on a text element is exactly the kind of reference a sanitizer drops --
the failure mode would be a figure with no such word in it at all, not a figure with a straight one.
The host font is monospace and every advance is the same 7.8px, so thirteen positions computed on an
arc are the layout rather than an approximation of the layout. The canvas is `220x180`: a lone
circle, sized to itself, with the height the renderer's bound allows.

Its alt text is upstream's `data-alt` verbatim, read out of the scrape with its encoding inverted.
Those three captions held three rows, and the regenerated audit dropped all of them: 272 to 269, and
165 in the certain tier. Three lessons leave the orphan list, so the report's lesson count goes 109
to 106, and the ceiling for orphan captions in `internal/lessons/assets_test.go` follows, to 269.

**Figure 5.7** (a page of a January calendar with the week's distances written into the days: 2, 15,
8, 3, 8, 5 and 8 miles) landed on 2026-10-16 in the five `stat.prob.*` lessons, as
`statistics/calendar-distances.svg`. It is the first reclaim outside the arithmetic directory, and
that is the only part of the wiring that is new: a lesson figure is a markdown image on an absolute
path, so a subject directory costs nothing and the figure lands where a reader looking for it would
look.

Upstream's calendar is blue-ruled, its header row is the darker of two shades and its distances are
grey italic. The house has one fill, so the header takes `#e6e6e6` and the body stays white, which is
the division the percent tables landed with -- the relationship survives (the header is the darker
strip) without a second fill. The distances are ink italic rather than grey, because they are what
the lesson is about: the mean of those seven numbers is the exercise, and the dates are scaffolding
the page around them. `New Year's day` is the one `#636363` label, because it names a day rather than
contributing a value.

The canvas is `460x170`, and the width decides the layout. Seven columns of 64 hold `Wednesday` at
the 11px the stylesheet forces (59.4px of a monospace advance), and the tallest row is the first,
which carries a date, a distance and the event name. That name is set on two lines -- `New Year's`
and `day` -- because at 8.5px it needs 76 of the 64 columns available, and a label that overflows its
cell in a calendar is a label a reader has to guess at.

Its alt text is upstream's `data-alt` verbatim, read out of the scrape with its encoding inverted.
Those five captions held five rows, and the regenerated audit dropped all of them: 269 to 264, and
160 in the certain tier. Five lessons leave the orphan list, so the report's lesson count goes 106 to
101, and the ceiling for orphan captions in `internal/lessons/assets_test.go` follows, to 264.

**Figure 2.12** (three graphs of three functions at `x = 2`: a line with a hole where its denominator
vanishes, the half-lines `y = 1` and `y = -1` whose ends are both holes, and `1/(x - 2)^2` running off
the top either side of the asymptote) landed on 2026-10-16 in `calc.limit.numeric`, as
`calculus/limit-behavior.svg`. It is the first graph and the first y axis in the local set, and the
first figure under the calculus directory.

The dialect had what a graph needs. Axes and their arrowheads are `#BABABA` at 1.5, their unit ticks
belong to the axis so they are the same color, the curve is ink at 2, and the hollow marker is a
`fill="white"` circle with an ink stroke -- the shape the vendored set already uses for a point that
is not on the graph, and the one place a figure carries white inside itself.

Three things upstream does not settle, and the choices are visible. The three panels share one y
range, so their x axes land on one row and they read as small multiples of a single question. Each
panel's x range is its own, because panel (b) has to carry the segment out to `x = 4` without leaving
its hole at the frame's edge, and panel (c) has nothing to show left of `x = -2.5`. And no tick
carries a number: the caption names the functions, the holes are drawn where they are, and nine
labeled ticks in a 153px panel would be a table rather than a picture.

This figure has a caption, so its alt is that caption minus its `Figure 2.12 ` label, the rule Figures
1.2 and 1.5 follow. Upstream's `data-alt` also exists, is longer, and is wrong in a way that matters:
it puts the second graph's holes at `(2, 1)` and `(-2, 1)`, and the second of those is not a point of
that graph, whose lower half-line is `y = -1`. Figures 5.5 and 6.5 record the same disagreement
between a drawing and its own description and resolve it toward the description, because there the
description was a figure that could be drawn. Here it cannot: no hole belongs at `(-2, 1)`. The
drawing follows the mathematics and the alt is the author's other sentence about it.

That one caption held one row, and the regenerated audit dropped it: 264 to 263, and 159 in the
certain tier. The lesson keeps other orphan captions, so the report's lesson count stays at 101, and
the ceiling for orphan captions in `internal/lessons/assets_test.go` follows the certain tier, to
263.

**Figures 2.13 and 2.14** (the two graphs whose captions say they confirm the table estimates:
`sin(x)/x` over `[-6, 6]`, empty at the origin because the function is not defined there and 1 is the
value it approaches; and `(sqrt(x) - 2)/(x - 4)` over `[0, 8]`, empty at `(4, 0.25)` for the same
reason) landed on 2026-10-16 in `calc.limit.numeric`, as `calculus/sine-over-x.svg` and
`calculus/radical-over-x.svg`. They are the first local figures whose ticks carry numbers, so the
labels are the only new thing the dialect needed: 11px, ink, every even integer on the x axis and
every twentieth on the y axis of 2.13.

The break in 2.14's y axis is the figure's real decision. Upstream runs 0.25 to 0.50 through about
two thirds of the frame, jogs under a zigzag, and gives the last thirteenth to 0.25 down to 0; the
curve is drawn through the same two-piece mapping, 364px per unit above the break and 72 below it.
The alternative, one linear axis from 0 to 0.5, puts the whole curve in the top sixth of an empty
frame and hides the flat tail -- which is the part that says the limit is 0. A break says the scale
changes here, which is true, and the zigzag is where it says it.

Both captions are mostly math, and an alt is plain text, so the caption's math is rendered into the
words it is made of by a small converter: `\frac{a}{b}` becomes `(a)/(b)`, `\sqrt{x}` becomes
`sqrt(x)`, `\text{sin}` becomes `sin`, and grouping braces and `$` come off. The converter asserts
on a macro it does not know, so no caption reaches a lesson as a backslash. For 2.14 the result is
`(sqrt(x) - 2)/(x - 4)`, which is how upstream's own `data-alt` writes that same function, and an
independent rendering that agrees with the author's is the check that this is a rendering and not a
paraphrase. 2.13's alt reads "The graph of f(x) = (sin x)/x confirms the estimate from Table 2.2."

Those two captions held two rows, and the regenerated audit dropped both: 263 to 261, and 157 in the
certain tier. The lesson keeps other orphan captions, so the report's lesson count stays at 101, and
the ceiling for orphan captions in `internal/lessons/assets_test.go` follows, to 261.

**Figures 2.15 and 2.16** (a curve whose value at `x = -1` is 4 while the curve itself passes through
3 there, drawn as a filled marker one unit above a hollow one; and a parabola whose vertex at
`(2, -1)` is the one point missing from it) landed on 2026-10-16 in `calc.limit.numeric`, as
`calculus/one-value-off-the-curve.svg` and `calculus/missing-vertex.svg`.

Neither caption gives a formula, which is the release to draw the shape the caption describes rather
than a function it does not name. 2.15's curve is therefore points chosen by hand and smoothed as a
Catmull-Rom chain -- steep out of the lower left, flat through the marker pair, steep again to the
upper right, which is the shape upstream's has -- and 2.16's is `0.25(x - 2)^2 - 1`, a parabola wide
enough that both arms and the vertex are inside a frame 180px tall. They are the first local figures
with a filled marker beside a hollow one, and their widths are set by the ratio the curve wants, 300
and 330, because a curve that rises five units over five units of x should not be drawn 3:1 wide.

**The batch's finding is that the audit could not see 2.16.** Its alt is upstream's `data-alt`
verbatim and that text contains "over [-2.5, 5]"; `IMAGE_OPEN_RE` in
`scripts/audit_lesson_media.py` stopped at the first `]`, so the audit did not read the line above
the caption as an image and counted the caption as an orphan anyway. `imageRe` in
`internal/lessons/assets.go` has carried a `\]` alternative since it was written -- its comment names
`numerical-integration.md`'s "The interval [a, b]" and `[n, n+1)`, whose bracket is never closed at
all -- so the Go test counted the same caption as wired. The two owners of one number disagreed by
one: 259 there and 260 here, on the same corpus. Section 1 quotes the audit, so the report was the
wrong one. The pattern now mirrors `imageRe` in both `IMAGE_RE` and `IMAGE_OPEN_RE`, and the four
vendored bracketed alts move no count at all: none of their captions is a row that a bracket is the
reason for. That the drift was invisible is the other half of the finding -- the ceiling is an upper
bound, so 259 and 260 both pass it, and the disagreement was only found by recomputing both plans
while adding a figure.

2.15 has a caption, so its alt is that caption minus its `Figure 2.15 ` label, with the `$g(x)$`
rendered as plain text by the converter 2.13 and 2.14 needed. 2.16 has the bare label upstream, so
its alt is upstream's `data-alt` verbatim, brackets and all.

Those two captions held two rows, and with the pattern fixed the regenerated audit drops both: 261 to
259, and 155 in the certain tier. The lesson keeps other orphan captions, so the report's lesson
count stays at 101, and the ceiling for orphan captions in `internal/lessons/assets_test.go` follows,
to 259.

**Figures 2.17, 2.18 and 2.19** (the last three graphs of the section: `sin(1/x)`, which oscillates
without settling as `x` approaches 0; a piecewise function with a break at `x = 2`, whose first piece
ends at an open circle at `(2, 3)` and whose second starts at a filled dot at `(2, 0)`; and `1/x` in
two branches) landed on 2026-10-16 in `calc.limit.numeric`, as `calculus/oscillation.svg`,
`calculus/piecewise-jump.svg` and `calculus/reciprocal.svg`.

The helper had to grow for them, and the growth is the reason they are one batch: this is the first
group needing ticks at tenths with labels at fifths (2.17), no labels at all (2.18, which upstream
also leaves unnumbered), and labels on both axes (2.19). The axes primitives take a tick list now
instead of looping over integers.

2.17 is the one with a method in it. `sin(1/x)` has no last oscillation, so sampling uniformly in `x`
spends its points where the curve is already smooth and starves the band. It is sampled uniformly in
`1/x` from 2 to 500 instead: every oscillation gets the same number of samples, the band closes into
the solid block upstream's has, and `x` down to 0.002 is as far in as a 280px frame shows anything
but a block. The count is 1,500 points per branch rather than 4,000 because the figure is inlined
into the lesson page — 4,000 was 124KB of path data and 1,500 is 47KB for the same picture. 2.18
draws its two pieces as two paths and marks both endpoints, which is the whole content of its
caption: the formula with the break is not at the left piece's end and is at the right piece's start.
2.19's branches leave the frame at `y = ±3.6` instead of being clipped into a chord, the same trick
2.12's third panel used.

The alts split three ways here, and the third way is new. 2.17 and 2.19 have captions that are prose
with a formula in them, so their alts are the caption minus the label with the formula rendered by
the converter 2.13 and 2.14 needed; the `*x*` stays, because the emphasis is the author's and the
renderer flattens it in the attribute. 2.18's caption is not prose at all: it is
`f(x) = \left\{ ... \right.`, a piecewise definition inside a matrix, and the converter refuses a
`\begin` rather than inventing words for a formula, which is the assert doing what it was written to
do. So 2.18 takes upstream's `data-alt` verbatim, and that alt describes the figure rather than
restating it: "The first piece is f(x) = x+1 if x < 2. ... There is an open circle at (2,3), where the
endpoint would be." The rule that follows is that a caption is a source for an alt while it is a
sentence about the figure; when the caption *is* the figure's formula, the `data-alt` is the sentence
and the caption is not.

Those three captions held three rows, and the regenerated audit drops all three: 259 to 256, and 152
in the certain tier. The lesson keeps other orphan captions, so the report's lesson count stays at
101, and the ceiling for orphan captions in `internal/lessons/assets_test.go` follows, to 256.

**Figures 2.21 and 2.23** (the figure an exercise is read off -- four limits and four values asked at
`x = -4, -2, 1` and `3`, with a curve interrupted twice and a vertical asymptote in the middle of it
-- and Einstein's mass ratio against speed ratio) landed on 2026-10-16 in `calc.limit.numeric`, as
`calculus/interrupted-curve.svg` and `calculus/mass-ratio.svg`. They are the last two redrawable
figures of Calculus 2.2.

2.21 is drawn in the alt's order rather than as one function, because that is what it is: a spline
below `x = -2` ending at an open circle at `(-2, 3)`, a second spline from there to a filled dot at
`(1, 6)`, an open circle at `(1, 3)` with a branch running down to `-infinity` along `x = 3`, and the
mirror branch on the other side of the asymptote, which comes back up through the concavity change
the alt names in the first quadrant. Only the two interesting y values are numbered, 3 and 6, and
only the four x values the exercise asks about -- upstream's choice, and the right one: numbers
everywhere would have a reader checking the graph instead of reading it.

2.23's axes are named for its quantities, `v/c` and `m/m0`, and the formula is written inside the
frame the way upstream writes it, as a fraction with a bar and the root below it. That covers the
whole figure: there is no tick label above 5 and no second series. `m0` is written with a plain
zero and not U+2080, because no theme font carries a subscript zero and a substituted face would sit
at a different size and baseline inside the drawing -- which is the objection AGENTS.md section 1
makes about a codepoint the fonts do not carry, and it applies to a subscript as much as to a
pictograph. The two superscripts are U+00B2, the one the arithmetic figures already use.

Both have captions that are prose, so both alts are the caption minus its label: 2.21's is "The graph
shows f(x)." and 2.23's is the sentence about Einstein's equation, kept with the author's
typographic apostrophe.

Those two captions held two rows, and the regenerated audit drops both: 256 to 254, and 150 in the
certain tier. The lesson keeps other orphan captions, so the report's lesson count stays at 101, and
the ceiling for orphan captions in `internal/lessons/assets_test.go` follows, to 254. With 2.23
landed, every redrawable figure in Calculus 2.2 is drawn: 2.12 through 2.19 and 2.21, with 2.20 and
2.22 remaining as fetched illustrations and a photograph.

**Figures 3.25, 3.26, 3.27 and 3.33** (the difference quotient the sine rule is guessed from; the
two limits its proof needs; sine and cosine in one frame with dots where the tangent is
horizontal; and the family `2^x`, `2.7^x`, `2.8^x`, `3^x` that places e) landed on 2026-10-16 in
`calc.deriv.trig` and `calc.deriv.exp_log`, as `calculus/difference-quotient.svg`,
`calculus/two-limits.svg`, `calculus/sine-and-cosine.svg` and `calculus/exponential-family.svg`.
They are the first figures outside section 2, and the first two lessons here to gain a figure at
all.

Three of the four needed the same layout fix, which makes it a rule: the frame has to run past its
last tick label, or the label, the arrowhead and the axis name pile up on top of each other. 3.25
runs to `2pi + 0.7`, 3.27 to `2pi + 0.55` and 3.33 to `0.72`.

3.25 draws `D(x) = (sin(x + 0.01) - sin x)/0.01` over `[-2pi, 2pi]`, and repeats the formula below
the frame rather than inside it, because the lesson writes it as a display equation and upstream
prints it under the plot.

3.26 is two panels with one curve each. Upstream gives each panel its own y range; this gives them
one, `[-0.85, 1.05]`, so that their x axes land in one row and the two limits can be read against
each other, which is the only reason to show them together. The cost is real and worth naming:
`(cos h - 1)/h` has amplitude 0.73 and is drawn at 82% of the size its own range would give it, and
the shared range then holds labels every 0.5 instead of upstream's every 0.2. The alt says the two
curves "have discontinuities on the y-axis", which is what the two hollow markers at `(0, 1)` and
`(0, 0)` are.

3.27 is the first local figure to use the trig accents section 4 of the style doc reserves for a
function family, `#2980b9` for the sine and `#e74c3c` for the cosine: two curves that cross four
times in one frame cannot be told apart by weight alone. The dots are the caption's statement --
the two extrema of the sine and the two zeros of the cosine at those same x values -- plus the
cosine's maxima at `x = 0` and `x = 2pi`, which are upstream's and are kept because a look-alike is
the goal. Upstream's dot set is not symmetric, and this record says so rather than inventing a dot
at `(pi, -1)` to match: there is no dot there.

3.33 has the x axis on the frame's bottom edge instead of at `y = 0`, because its y range starts at
0.5, and the y axis carries a break glyph to say so -- the glyph 2.14 introduced for a two-piece
scale, here for one that simply starts above zero. Its four curves all end at the right edge so
each label sits at its own curve's end, and the two the caption names (`2^x` and `3^x`) are ink at
width 2 while the pair that brackets e is annotation grey at 1.5. The labels read `3^x`, not a
raised superscript, for the same reason 2.23's `m0` is a plain zero.

**No new e2e case.** None of the four introduces a primitive the sanitizer has not already been
shown to pass: a colored stroke is the same `stroke` attribute the ink curves carry, and the break
glyph and the bottom x axis are shapes 2.14 and 2.12 already draw. The screenshot check carries
them, and `calc.deriv.trig` is a lesson that had no figure before today, which the
`arith.mult.tables` case already covers as a pattern.

**The lesson count moved, not only the row count.** Figure 3.33 was `calc.deriv.exp_log`'s only
orphan caption, so with it wired that lesson leaves the report entirely: **101 lessons to 100**.
Every earlier batch took rows off the list and left the lesson count alone, so this is the first
time the third number in section 1 has moved, and the first batch where the Go test's own log line
and the audit disagree with the docs unless all three are updated together.

Those four captions held four rows: the report went from 254 to 250, and the certain tier from
150 to 146. The ceiling for orphan captions in `internal/lessons/assets_test.go` follows, to
250. (The numbers are written in this order rather than as one live claim because the checker
reads a count followed by the word it counts: a dated record must not look like today's
measurement, or every later batch has to rewrite history to keep the checker quiet.)

**Figures 4.40, 4.42 and 4.45** (a curve approaching the horizontal asymptote `y = 2`; a curve that
crosses its horizontal asymptote an infinite number of times; and `tan x` with vertical asymptotes
at `x = +-pi/2`) landed on 2026-10-16 in `calc.limit.infinity`, as
`calculus/horizontal-asymptote.svg`, `calculus/crossing-asymptote.svg` and
`calculus/tangent-asymptotes.svg`. This is the first batch whose subject is the asymptote itself,
and the first whose third figure goes into a numbered list item rather than a paragraph.

**The dashed guide is new to the local set, and already proven.** No local calculus figure used
`stroke-dasharray` before today; the vendored Algebrica set does, and three of the lessons the e2e
list already checks carry one (`alg.log.concept`, `calc.deriv.cauchy_mvt`, `calc.seq.euler`), so the
attribute is one the sanitizer has been shown to pass for months. That is why this batch adds no
e2e case: the new shape is new *here* and not new to the sanitizer, which is the distinction the
rule was written for. The guides themselves are annotation grey `#636363` rather than the axes'
`#BABABA`: a labeled asymptote is an annotation, not an axis.

**4.40's caption does not name its function, and the figure still had to be drawable.** "The
function approaches the asymptote `y = 2` as `x` approaches `+-oo`" is prose about a curve, so the
alt is that sentence and the curve has to come from somewhere else: the lesson's own prose later in
the section ("we used graphical evidence in Figure 4.40 and numerical evidence in Table 4.2 to
conclude that `lim (2 + 1/x) = 2`) and upstream's `data-alt`, which reads "The function f(x) 2 +
1/x is graphed". Taking the function from the lesson rather than the caption is a decision, and it
is the reason this row was classified redraw while its neighbour 4.43, whose caption names no
function *and* whose prose names it only in the surrounding paragraph, was classified fetch.

**Both 4.40 and 4.42 draw one path per branch, from analytic domains, with no clip path.** The
right branch of `2 + 1/x` starts where the curve leaves the top of the frame, at `x = 1/3`, and the
left branch ends where it leaves the bottom, at `x = -1/7`; 4.42's branches end where the curve
leaves the frame at `y = 2.2` and `y = -1.4`, found by bisection. That is how `reciprocal.svg` drew
`1/x` and it means no path ever leaves the box, which is visible rather than merely contained.
4.40 also needed a wider x range than its ticks: labels every 2 instead of upstream's every 1,
because at 14.4px per unit a two-character label does not fit between two ticks. 4.26's record
already made that trade; this is the second time.

**4.45 is the first figure in the set that a numbered list item carries.** Its caption is indented
four spaces upstream, so the insertion preserves the indentation; an image at column 0 would have
ended the list item and left the caption orphaned in the markdown in a way the audit would still
count as wired. The label placement is worth recording too: `y = tan x` sits *left* of the y axis
between the two asymptotes, not right of it as upstream has it, because at 45.9px per unit the free
span right of the axis ends where the curve crosses the label's own row, and the label is 59px
wide. The `x = pi/2` label had to drop below the x axis for the same reason in reverse: at row 97 it
landed on the "2" and "3" tick labels, which sit at row+15.

**Three new macros and one rule in the alt converter.** 4.40's caption ends
`$\text{\pm}\infty$` and 4.45's reads `$x = \text{\pm}\frac{\pi}{2}$`, so the converter now
renders `\pm`, `\infty` and `\pi` instead of asserting on them, and it no longer brackets a
fraction whose numerator and denominator are each a single token: 4.45's alt reads `x = +-pi/2`,
not `x = +-(pi)/(2)`. A part with an operator inside keeps its parentheses, so 2.14's committed alt
is unchanged and the wire script checks that against the string in the lesson rather than against a
comment. One further fix came out of the same caption: upstream writes `\left( {\text{cos} x}`,
and the space after the opening parenthesis leaked into 4.42's alt as `( cos x )/x`, so the
converter now closes up whitespace inside parentheses.

**The checker caught a stale count in the previous record, which is worth writing down.** After
this batch's numbers moved, `scripts/check_docs.py` failed on the 3.25-3.33 record, which ended its
row count with the number 250 and the word "candidates" straight after it: the checker reads any
count followed by the word it counts as a live claim, and a dated record is not one. The fix was
to phrase the history as a sequence ("went from 254 to 250") rather than as a measurement, and then
to write this paragraph so that it did not re-create the pattern while describing it. Every later batch has to keep that distinction, or the records have
to be rewritten to match today's corpus.

Those three captions held three rows: the report went from 250 to 247, and the certain tier from
146 to 143. The lesson keeps nineteen other orphan captions, so the lesson count stays at 100, and
the ceiling in `internal/lessons/assets_test.go` follows, to 247.

**Figures 4.50, 4.51 and 4.52** (the even powers `x^2`, `x^4` and `x^6`; the odd powers `x`, `x^3`
and `x^5`; and a polynomial over its leading term) landed on 2026-10-16 in `calc.limit.infinity`,
as `calculus/power-even.svg`, `calculus/power-odd.svg` and `calculus/polynomial-leading-term.svg`.

**4.50 and 4.51 are the first local figures to use the accents for something that is not
trigonometric.** Section 4 of the style doc gives the rule's scope as "where a function is compared
across a family" and names the three colors after the functions that first needed them. This batch
needs them for a measured reason: at the frame's top edge 4.50's `x^4` and `x^6` curves exit 12px
apart, and they are closer at every row below it, so no set of positional labels can identify them —
and their left-to-right order flips at `x = +-1`, so a stacked set cannot either. 3.33 manages four
exponential curves without color because its labels sit in the right margin, each at the height
where its own curve leaves the frame; these curves all leave through the top. Upstream encodes the
same three curves by color, and the vendored set colors text as well as strokes (`gen-diagrams.py`
writes `#2980b9` labels), so here the label and its curve carry the same accent.

**4.52 takes no accent, and that is the same rule read the other way.** Its two curves coincide at
both ends — that is the caption's statement — so they are told apart by the ladder
`exponential-family.svg` uses: `#333333` at width 2 for `f(x) = 5x^3 - 3x^2 + 4`, and `#636363` at
1.5 for `g(x) = 5x^3`, the term it behaves like. A reader who cannot see the color still has two
labels, and each label says which function it belongs to.

**The pair shares one x scale, and that is what buys the labels their room.** Both frames are 420px
wide over `[-5, 5]`, or 36.6px per unit against upstream's 38.4. 4.51 needed it: at `y = 10` the
`x^5` curve's right branch is 58px from the y axis, so a 46px label fits beside it, where at a
360px canvas the gap was 45px and the label had to overlap the axis's own name. 4.50's three labels
are stacked in the bottom-right wedge instead of upstream's top-right corner, by the same kind of
measurement: that corner has 14px between the `x^2` curve's exit and the frame's edge, and the wedge
has 35px at the top row of the stack. One more row moved after looking at the render, which is the
only way this class of defect shows up: 4.51's `y = x^5` label sat at baseline 27, on the rows the y
axis's own name and arrowhead occupy, and it is at 34 now.

**The converter learned two macros, and one of them was silently wrong first.** 4.50 and 4.51 state
their end behavior in OpenStax's display notation, `\underset{x\rightarrow\infty}{\text{lim}}x^{n}`.
`\underset{A}{B}` renders as `B(A)`: dropping A would lose which limit is meant, and appending it
bare reads as `lim x -> oo x^n`, ambiguous about where the argument starts. `\rightarrow` renders as
` -> ` with LaTeX's own relation spacing, which the captions do not write. The first attempt put the
`\rightarrow` branch after the `\left`/`\right` branch, and since `\rightarrow` starts with
`\right`, that branch ate it and produced `xarrowoo` in all three alts. The wire script asserts the
three expected strings rather than printing them, which is why the defect did not reach the corpus;
that assert is the unknown-macro assert one level up.

Those three captions held three rows: the report went from 247 to 244, and the certain tier from
143 to 140. The lesson keeps sixteen other orphan captions, so the lesson count stays at 100, and
the ceiling in `internal/lessons/assets_test.go` follows, to 244.

## 8. Order of work

The list is three tiers now rather than four: the `unidentified` tier used to lead here as its own
step, and section 4 records that it was emptied on 2026-10-16, so those seventeen figures sit inside
items 1 and 2 by their new dispositions.

1. **The schematics.** Highest value and no license question: they are redrawn as Mathua's own
   work, in the dialect the style doc fixes, and they land under `source: local`.
2. **The illustrations.** Fetch, credit and record the upstream link in the metadata entry.
3. **The photos.** Decide per figure between a licensed fetch and a replacement, because a
   photograph cannot be redrawn as a line figure.

Each figure is wired by the paths in section 5 of `docs/diagram-style.md`, and the published
`web/next-app/public/diagrams/meta.json` copy is refreshed before the change is called done.

## 9. Method

```bash
python3 scripts/audit_lesson_media.py --json data/lessons/figure-audit.json > docs/lesson-media.md
```
