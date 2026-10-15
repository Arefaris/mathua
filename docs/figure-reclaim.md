# Reclaiming the dropped figures

This is the work plan for the figures the lesson ingestion dropped while keeping their captions.
It is an **analytical-tier** doc (AGENTS.md section 6), so `scripts/check_ste.py` does not cover
it. The measurement belongs to `scripts/audit_lesson_media.py`, which writes
`data/lessons/figure-audit.json` and the report in `docs/lesson-media.md`; this file turns the
`certain` tier of that report into a plan.

## 1. What is missing

318 paragraphs across 113 lessons read as a caption whose image was dropped:
214 certain, 33 strong and 71 likely. Every one sits
in `data/lessons/teaching/`, and until the first reclaim below those lessons contained no image at
all, so the caption was the only remaining trace of the figure. The `certain` tier opens with an
auto-numbered `Figure 3.2` label that the PreTeXt conversion emits for a figure and keeps when it
drops the image, so it needs no judgement.

## 2. The upstream is OpenStax, not Algebrica

Every teaching lesson opens with a credit line naming its source, for example:

```
> Content sourced from [OpenStax Prealgebra](https://openstax.org/books/prealgebra/pages/1-introduction) by Marecek & Anthony-Smith - CC BY 4.0
```

Two books account for every `certain` caption: **OpenStax Calculus Volume 1** (62 figures) and
**OpenStax Prealgebra** (33). Both are CC BY 4.0, which permits reuse and adaptation with
attribution, and the attribution is already the lesson's own header line.

That changes what "a direct look-alike" means for these figures. Their visual target is their own
upstream, not Algebrica. The Algebrica dialect in `docs/diagram-style.md` applies when a figure is
redrawn by hand; a figure re-fetched from OpenStax keeps the OpenStax look and does not need to
match Algebrica.

## 3. The captions are 214 rows and 95 distinct figures

Teaching lessons share source sections, so one figure appears in several concept lessons: the same
`Figure 1.13` is the caption in five of them, and `Figure 1.11` in six. Deduplicating on the pair
(figure number, caption) leaves **95 distinct figures** still to reclaim. Section 6 gives
the method.

## 4. Dispositions

| Disposition | Figures | Action |
| --- | --- | --- |
| schematic, redraw | 35 | Draw in the Algebrica cool dialect from section 4 of the style doc. Graphs, number lines, regions and solids are fully determined by their captions, so no upstream fetch is needed and the result is Mathua's own work. |
| illustration, fetch | 38 | Fetch the OpenStax original. The license is CC BY 4.0 and the credit line is already in the lesson header. |
| unidentified | 17 | The caption is a bare `Figure 1.13` with no description, so the upstream section has to be opened to find out what the figure is before anything can be drawn or fetched. |
| photo, third-party credit | 5 | The caption names a photographer or agency (`credit: modification of work by Steve Jurvetson, Wikimedia Commons` and similar). These cannot be redrawn; each needs a licensed fetch or a replacement. |

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
  their `figcaption` holds only the label — so they stay `unidentified` in section 6 and are reclaim
  targets, not repair targets.
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
and still rejects a thematic break (`* * *`); the report is back at 318 candidates, 214 certain.

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

Sorted by disposition, then figure number. `Lessons` is how many concept lessons carry the
caption; one figure is restored once and every lesson that cites it gains the image. The caption is
the text as repaired in section 5, truncated at 85 characters, so it is what the figure will carry
when it lands.

| Figure | Caption (repaired) | Book | Section | Lessons | Disposition |
| --- | --- | --- | --- | --- | --- |
| 1.12 | The area of a rectangle is the product of its length and its width, or 6 square feet. | Prealgebra | 1 Multiply Whole Numbers | 6 | schematic, redraw |
| 2.12 | These graphs show the behavior of three different functions around $x = 2.$ | Calculus Volume 1 | 2 The Limit of a Function | 1 | schematic, redraw |
| 2.13 | The graph of $f(x) = {{(\text{sin} x)}\text{/}x}$ confirms the estimate from Table 2.... | Calculus Volume 1 | 2 The Limit of a Function | 1 | schematic, redraw |
| 2.14 | The graph of $f(x) = \frac{\sqrt{x} - 2}{x - 4}$ confirms the estimate from Table 2.3... | Calculus Volume 1 | 2 The Limit of a Function | 1 | schematic, redraw |
| 2.15 | The graph of $g(x)$ includes one value not on a smooth curve. | Calculus Volume 1 | 2 The Limit of a Function | 1 | schematic, redraw |
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
| 4.42 | The graph of $f(x) = \left( {\text{cos} x} \right)\text{/}x + 1$ crosses its horizont... | Calculus Volume 1 | 4 Limits at Infinity and Asymptotes | 1 | schematic, redraw |
| 4.45 | The graph of $\text{tan} x$ has vertical asymptotes at $x = \text{\pm}\frac{\pi}{2}$ | Calculus Volume 1 | 4 Limits at Infinity and Asymptotes | 1 | schematic, redraw |
| 4.53 | The graph of this rational function approaches a horizontal asymptote as $x\rightarro... | Calculus Volume 1 | 4 Limits at Infinity and Asymptotes | 1 | schematic, redraw |
| 4.54 | The graph of this rational function approaches the horizontal asymptote $y = 0$ as $x... | Calculus Volume 1 | 4 Limits at Infinity and Asymptotes | 1 | schematic, redraw |
| 4.56 | The graph of the rational function $f(x) = \left( {3x^{2} + 4x} \right)\text{/}\left(... | Calculus Volume 1 | 4 Limits at Infinity and Asymptotes | 1 | schematic, redraw |
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
| 1.7 | We can see that 76 is closer to 80 than to 70. So 76 rounded to the nearest ten is 80... | Prealgebra | 1 Introduction to Whole Numbers | 6 | illustration, fetch |
| 1.8 | We can see that 72 is closer to 70, so 72 rounded to the nearest ten is 70. | Prealgebra | 1 Introduction to Whole Numbers | 6 | illustration, fetch |
| 1.9 | The number 75 is exactly midway between 70 and 80. | Prealgebra | 1 Introduction to Whole Numbers | 6 | illustration, fetch |
| 2.10 | Factors of the counting numbers from 2 through 20, with prime numbers highlighted | Prealgebra | 2 Find Multiples and Factors | 4 | illustration, fetch |
| 2.20 | The function $f(x) = {1\text{/}\left( {x - a} \right)^{n}}$ has infinite limits at *a... | Calculus Volume 1 | 2 The Limit of a Function | 1 | illustration, fetch |
| 2.6 | Multiples of 2 between 1 and 50 | Prealgebra | 2 Find Multiples and Factors | 4 | illustration, fetch |
| 2.7 | Multiples of 5 between 1 and 50 | Prealgebra | 2 Find Multiples and Factors | 4 | illustration, fetch |
| 2.8 | Multiples of 10 between 1 and 50 | Prealgebra | 2 Find Multiples and Factors | 4 | illustration, fetch |
| 2.9 | Multiples of 3 between 1 and 50 | Prealgebra | 2 Find Multiples and Factors | 4 | illustration, fetch |
| 3.18 | The derivative of $g(x)$ is 3 times the derivative of $f(x).$ | Calculus Volume 1 | 3 Differentiation Rules | 3 | illustration, fetch |
| 3.19 | This function has horizontal tangent lines at *x* = 2/3 and *x* = 4. | Calculus Volume 1 | 3 Differentiation Rules | 3 | illustration, fetch |
| 3.21 | (a) One section of the racetrack can be modeled by the function $f(x) = x^{3} + 3x^{2... | Calculus Volume 1 | 3 Differentiation Rules | 3 | illustration, fetch |
| 3.30 | The equation $x^{2} + y^{2} = 25$ defines many functions implicitly. | Calculus Volume 1 | 3 Implicit Differentiation | 1 | illustration, fetch |
| 3.31 | The line $y = \frac{3}{4}x - \frac{25}{4}$ is tangent to $x^{2} + y^{2} = 25$ at the ... | Calculus Volume 1 | 3 Implicit Differentiation | 1 | illustration, fetch |
| 3.32 | Finding the tangent line to the folium of Descartes at $\left( {\frac{3}{2},\frac{3}{... | Calculus Volume 1 | 3 Implicit Differentiation | 1 | illustration, fetch |
| 3.34 | The tangent line to $E(x) = e^{x}$ at $x = 0$ has slope 1. | Calculus Volume 1 | 3 Derivatives of Exponential and Logarithmic Functions | 1 | illustration, fetch |
| 3.35 | $\text{The function}\ y = \text{ln} x$ is increasing on $(0,\text{+}\infty).$ Its der... | Calculus Volume 1 | 3 Derivatives of Exponential and Logarithmic Functions | 1 | illustration, fetch |
| 4.2 | As the balloon is being filled with air, both the radius and the volume are increasin... | Calculus Volume 1 | 4 Related Rates | 1 | illustration, fetch |
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
| 5.4 | (a) We see that 2.72 is closer to 3 than to 2. So, 2.72 rounded to the nearest whole ... | Prealgebra | 5 Decimals | 3 | illustration, fetch |
| 6.2 | Among every 100 community college students, 57 are female. | Prealgebra | 6 Understand Percent | 3 | illustration, fetch |
| 6.6 | A $\text{20\%}$ tip for an $\text{\$80}$ restaurant bill comes out to $\text{\$16}.$ | Prealgebra | 6 Solve General Applications of Percent | 3 | illustration, fetch |
| 6.7 | The sales tax is calculated as a percent of the purchase price. | Prealgebra | 6 Solve Sales Tax, Commission, and Discount Applications | 2 | illustration, fetch |
| 1.11 | (none) | Prealgebra | 1 Multiply Whole Numbers | 6 | unidentified |
| 1.13 | (none) | Prealgebra | 1 Divide Whole Numbers | 5 | unidentified |
| 1.14 | (none) | Prealgebra | 1 Divide Whole Numbers | 5 | unidentified |
| 1.15 | (none) | Prealgebra | 1 Divide Whole Numbers | 5 | unidentified |
| 1.3 | (none) | Prealgebra | 1 Introduction to Whole Numbers | 6 | unidentified |
| 1.4 | (none) | Prealgebra | 1 Introduction to Whole Numbers | 6 | unidentified |
| 2.16 | (none) | Calculus Volume 1 | 2 The Limit of a Function | 1 | unidentified |
| 4.2 | (none) | Prealgebra | 4 Visualize Fractions | 7 | unidentified |
| 4.3 | (none) | Prealgebra | 4 Visualize Fractions | 7 | unidentified |
| 4.4 | (none) | Prealgebra | 4 Visualize Fractions | 7 | unidentified |
| 4.5 | (none) | Prealgebra | 4 Multiply and Divide Fractions | 4 | unidentified |
| 5.5 | (none) | Prealgebra | 5 Decimal Operations | 4 | unidentified |
| 5.6 | (none) | Prealgebra | 5 Decimals and Fractions | 3 | unidentified |
| 5.7 | (none) | Prealgebra | 5 Averages and Probability | 5 | unidentified |
| 6.3 | (none) | Prealgebra | 6 Understand Percent | 3 | unidentified |
| 6.34 | (none) | Calculus Volume 1 | 6 Volumes of Revolution: Cylindrical Shells | 1 | unidentified |
| 6.5 | (none) | Prealgebra | 6 Understand Percent | 3 | unidentified |
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

## 8. Order of work

1. **The schematics.** Highest value and no license question: they are redrawn as Mathua's own
   work, in the dialect the style doc fixes, and they land under `source: local`.
2. **The illustrations.** Fetch, credit and record the upstream link in the metadata entry.
3. **The unidentified.** Open the upstream section named in the table and identify each one
   first; they then fall into one of the other dispositions.
4. **The photos.** Decide per figure between a licensed fetch and a replacement, because a
   photograph cannot be redrawn as a line figure.

Each figure is wired by the paths in section 5 of `docs/diagram-style.md`, and the published
`web/next-app/public/diagrams/meta.json` copy is refreshed before the change is called done.

## 9. Method

```bash
python3 scripts/audit_lesson_media.py --json data/lessons/figure-audit.json > docs/lesson-media.md
```
