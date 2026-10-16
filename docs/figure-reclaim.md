# Reclaiming the dropped figures

This is the work plan for the figures the lesson ingestion dropped while keeping their captions.
It is an **analytical-tier** doc (AGENTS.md section 6), so `scripts/check_ste.py` does not cover
it. The measurement belongs to `scripts/audit_lesson_media.py`, which writes
`data/lessons/figure-audit.json` and the report in `docs/lesson-media.md`; this file turns the
`certain` tier of that report into a plan.

## 1. What is missing

276 paragraphs across 113 lessons read as a caption whose image was dropped:
172 certain, 33 strong and 71 likely. Every one sits
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

## 3. The captions are 172 rows and 93 distinct figures

Teaching lessons share source sections, so one figure appears in several concept lessons: the same
`Figure 1.13` is the caption in five of them, and `Figure 1.11` in six. Deduplicating on the pair
(figure number, caption) leaves **93 distinct figures** still to reclaim. Section 6 gives
the method.

## 4. Dispositions

| Disposition | Figures | Action |
| --- | --- | --- |
| schematic, redraw | 44 | Draw in the Algebrica cool dialect from section 4 of the style doc. Graphs, number lines, regions and solids are fully determined by their captions, so no upstream fetch is needed and the result is Mathua's own work. |
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
again to 53 and 44, of which section 7 has since landed nine.

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
| 2.12 | These graphs show the behavior of three different functions around $x = 2.$ | Calculus Volume 1 | 2 The Limit of a Function | 1 | schematic, redraw |
| 2.13 | The graph of $f(x) = {{(\text{sin} x)}\text{/}x}$ confirms the estimate from Table 2.... | Calculus Volume 1 | 2 The Limit of a Function | 1 | schematic, redraw |
| 2.14 | The graph of $f(x) = \frac{\sqrt{x} - 2}{x - 4}$ confirms the estimate from Table 2.3... | Calculus Volume 1 | 2 The Limit of a Function | 1 | schematic, redraw |
| 2.15 | The graph of $g(x)$ includes one value not on a smooth curve. | Calculus Volume 1 | 2 The Limit of a Function | 1 | schematic, redraw |
| 2.16 | (none) | Calculus Volume 1 | 2 The Limit of a Function | 1 | schematic, redraw |
| 2.17 | The graph of $f(x) = \text{sin}({1\text{/}{x)}}$ oscillates rapidly between −1 and 1 ... | Calculus Volume 1 | 2 The Limit of a Function | 1 | schematic, redraw |
| 2.18 | The graph of $f(x) = \left\{ \begin{matrix} {x + 1\ \text{if}\ x < 2} \\ {x^{2} - 4\ ... | Calculus Volume 1 | 2 The Limit of a Function | 1 | schematic, redraw |
| 2.19 | The graph of $f(x) = {1\text{/}x}$ confirms that the limit as *x* approaches 0 does n... | Calculus Volume 1 | 2 The Limit of a Function | 1 | schematic, redraw |
| 2.21 | The graph shows $f(x).$ | Calculus Volume 1 | 2 The Limit of a Function | 1 | schematic, redraw |
| 2.23 | This graph shows the ratio of masses as a function of the ratio of speeds in Einstein... | Calculus Volume 1 | 2 The Limit of a Function | 1 | schematic, redraw |
| 3.25 | The graph of the function $D(x)$ looks a lot like a cosine curve. | Calculus Volume 1 | 3 Derivatives of Trigonometric Functions | 1 | schematic, redraw |
| 3.26 | These graphs show two important limits needed to establish the derivative formulas fo... | Calculus Volume 1 | 3 Derivatives of Trigonometric Functions | 1 | schematic, redraw |
| 3.27 | Where $f(x)$ has a maximum or a minimum, $f'(x) = 0$ that is, $f'(x) = 0$ where $f(x)... | Calculus Volume 1 | 3 Derivatives of Trigonometric Functions | 1 | schematic, redraw |
| 3.33 | The graph of $E(x) = e^{x}$ is between $y = 2^{x}$ and $y = 3^{x}.$ | Calculus Volume 1 | 3 Derivatives of Exponential and Logarithmic Functions | 1 | schematic, redraw |
| 4.3 | An airplane is flying at a constant height of 4000 ft. The distance between the perso... | Calculus Volume 1 | 4 Related Rates | 1 | schematic, redraw |
| 4.40 | The function approaches the asymptote $y = 2$ as $x$ approaches $\text{\pm}\infty.$ | Calculus Volume 1 | 4 Limits at Infinity and Asymptotes | 1 | schematic, redraw |
| 4.42 | The graph of $f(x) = \left( {\text{cos} x} \right)\text{/}x + 1$ crosses its horizont... | Calculus Volume 1 | 4 Limits at Infinity and Asymptotes | 1 | schematic, redraw |
| 4.45 | The graph of $\text{tan} x$ has vertical asymptotes at $x = \text{\pm}\frac{\pi}{2}$ | Calculus Volume 1 | 4 Limits at Infinity and Asymptotes | 1 | schematic, redraw |
| 4.50 | For power functions with an even exponent of $n,$ $\underset{x\rightarrow\infty}{\tex... | Calculus Volume 1 | 4 Limits at Infinity and Asymptotes | 1 | schematic, redraw |
| 4.51 | For power functions with an odd exponent of $n,$ $\underset{x\rightarrow\infty}{\text... | Calculus Volume 1 | 4 Limits at Infinity and Asymptotes | 1 | schematic, redraw |
| 4.52 | The end behavior of a polynomial is determined by the behavior of the term with the l... | Calculus Volume 1 | 4 Limits at Infinity and Asymptotes | 1 | schematic, redraw |
| 4.53 | The graph of this rational function approaches a horizontal asymptote as $x\rightarro... | Calculus Volume 1 | 4 Limits at Infinity and Asymptotes | 1 | schematic, redraw |
| 4.54 | The graph of this rational function approaches the horizontal asymptote $y = 0$ as $x... | Calculus Volume 1 | 4 Limits at Infinity and Asymptotes | 1 | schematic, redraw |
| 4.56 | The graph of the rational function $f(x) = \left( {3x^{2} + 4x} \right)\text{/}\left(... | Calculus Volume 1 | 4 Limits at Infinity and Asymptotes | 1 | schematic, redraw |
| 4.58 | The function $f(x) = \text{sin} x$ oscillates between $1\ \text{and}\ -1$ as $x\right... | Calculus Volume 1 | 4 Limits at Infinity and Asymptotes | 1 | schematic, redraw |
| 4.59 | The function $f(x) = \text{tan} x$ does not approach a limit and does not approach $\... | Calculus Volume 1 | 4 Limits at Infinity and Asymptotes | 1 | schematic, redraw |
| 5.32 | The graph shows speed versus time for the given motion of a car. | Calculus Volume 1 | 5 Integration Formulas and the Net Change Theorem | 1 | schematic, redraw |
| 5.33 | The graph shows velocity versus time for a particle moving with a linear velocity fun... | Calculus Volume 1 | 5 Integration Formulas and the Net Change Theorem | 1 | schematic, redraw |
| 5.35 | Graph (a) shows the positive area between the curve and the *x*-axis, whereas graph (... | Calculus Volume 1 | 5 Integration Formulas and the Net Change Theorem | 1 | schematic, redraw |
| 5.36 | The graph shows areas between a curve and the *x*-axis for an odd function. | Calculus Volume 1 | 5 Integration Formulas and the Net Change Theorem | 1 | schematic, redraw |
| 5.5 | (none) | Prealgebra | 5 Decimal Operations | 4 | schematic, redraw |
| 5.6 | (none) | Prealgebra | 5 Decimals and Fractions | 3 | schematic, redraw |
| 5.7 | (none) | Prealgebra | 5 Averages and Probability | 5 | schematic, redraw |
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
