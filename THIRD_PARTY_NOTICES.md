# Third-Party Notices

Mathua is MIT licensed (see `LICENSE`). It also depends on, ships, or redistributes
content from the projects below. This file records those obligations in one place so
they are not scattered across the README and forgotten when a dependency changes.

## Software dependencies

### MathLive — `mathlive@0.111.0`

Used for the answer editor in `web/next-app/components/math/`. Provides the
`<math-field>` custom element: mathematical input, structure navigation, and the
virtual keyboard. Licensed under the MIT License.

MathLive is **used, not vendored**. Nothing in it is copied into this repository; the
integration goes through its public API, configuration, custom keyboard layouts, CSS
variables and `::part()`. MathLive brings its own webfonts (the KaTeX family), which
`web/next-app/scripts/copy-mathlive-fonts.mjs` copies from `node_modules` into
`public/fonts` at build time. Those files are build output and are not committed.

```
Copyright (c) 2017 - present Arno Gourdol. All rights reserved.

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```

Upstream: <https://github.com/arnog/mathlive> · Documentation: <https://mathlive.io/mathfield/>

### Compute Engine — `@cortex-js/compute-engine@0.147.0`

Used for one thing: `web/next-app/components/math/SelfCheck.tsx`, which tells a learner
what their own expression simplifies to. Licensed under the MIT License.

Compute Engine is **used, not vendored**, and is a browser dependency only — it is
dynamically imported so it never enters any route's initial JS, and it is skipped on
metered connections. It is *not* used for grading: the Go engine and SymPy remain the
only graders (ADR-005), and the self-check cannot reach a submitted answer.

Note for anyone extending this: importing the library also switches on MathLive's
MathJSON paste support, because Compute Engine registers itself under
`globalThis[Symbol.for("io.cortexjs.compute-engine")]` as a side effect of being
imported. That is an implicit coupling and is asserted in
`web/next-app/test/equivalence.test.ts`.

```
MIT License

Copyright (c) 2019 CortexJS

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```

Upstream: <https://cortexjs.io/compute-engine/>

### Other npm dependencies

KaTeX (MIT) renders question text, answers and lesson content; three.js and
`@react-three/fiber` (MIT) render the landing concept graph; dagre, lucide-react,
remark-math and rehype-katex (MIT) support layout, icons and Markdown. Their licence
texts ship inside each package under `node_modules/`, and are reproduced here on
upgrade by the package manager.

## Go dependencies

The Go module's dependency licences are recorded in `go.sum` and in each module's own
repository. SymPy (MIT) is the grading runtime; it is pinned in
`grading/requirements.txt` and installed into the image.

## Content

### Algebrica — CC BY-NC 4.0

Lesson content and diagrams under `data/lessons/` and `figures/` are sourced from
[Algebrica](https://algebrica.org) by Antonio Lupetti, used under
[CC BY-NC 4.0](https://creativecommons.org/licenses/by-nc/4.0/).

**This is a non-commercial licence.** Mathua's own MIT licence does not extend to this
content: the lesson material and diagrams remain CC BY-NC 4.0 and carry that
restriction with them.

### DiceBear — per-style licences

Profile avatars are served by [DiceBear](https://www.dicebear.com). Individual avatar
styles carry their own licences; see their
[license overview](https://www.dicebear.com/licenses/).

### Footer artwork — UNRESOLVED, do not ship without settling this

The eight prints in `web/next-app/public/footer/` are shipped to every visitor: two
desktop landscapes and six phone portraits, drawn full-bleed behind the footer.

**Their provenance and licence are not established.** They were collected from image
search without a recorded source, author, or licence, and nothing in the repository
says otherwise. This entry is here because the obligation is real and unsatisfied,
not because it is settled.

What follows from that, stated plainly:

- The files are **not** covered by Mathua's MIT licence. A licence covers what the
  licensor may give; nothing here establishes that anyone did.
- Redistributing them is the risk. Mathua publishes this repository and deploys the
  static site from it, so the prints leave the machine.
- Replacing them is the clean fix. The footer's contract is only that each file
  exists, declares the real dimensions, and carries a tone — `test/footerArt.test.tsx`
  holds that contract, and it names no particular picture. Dropping in known-provenance
  artwork needs no code change.

To resolve this: record the author, source URL and licence for each print below, or
replace the set. Until one of those happens, treat every print in that directory as
unlicensed.

### OpenStax, MIT OpenCourseWare, Art of Problem Solving

Referenced for content design and curriculum structure. No code or content is taken
from any of them. See the Acknowledgements section of `README.md`.

### Mathematical glossary sources — consulted, original wording

The glossary dataset (`web/next-app/data/glossary/`) was written in original
wording for Mathua. Its sources were consulted as references, not copied: no
text was adapted from any of them, so no ShareAlike or attribution clause is
triggered by adaptation. Each entry still records the source consulted, with
the reuse terms as verified at authoring time:

- **OpenStax textbooks** (openstax.org) for elementary through university
  mathematics. The openstax.org footer states its textbooks are licensed
  under CC BY-NC-SA 4.0 (verified on the Introductory Statistics section
  page, whose footer also restricts LLM training ingestion — not applicable
  here since no text is reused). Entries cite the book consulted and use
  original wording with no text reused, so no ShareAlike obligation is
  triggered; the license is recorded for accuracy.
- **NIST Digital Library of Mathematical Functions** (https://dlmf.nist.gov/)
  for special functions and advanced notation (gamma, error functions, Bessel,
  zeta, combinatorics, number theory). The DLMF is produced by NIST; entries
  cite the chapter consulted and users should follow the DLMF's own citing
  guidance. Definitions are original.
- **Wikibooks Subject:Mathematics**
  (https://en.wikibooks.org/wiki/Subject:Mathematics) for general notation,
  logic and set theory. The site footer states text is available under the
  Creative Commons Attribution-ShareAlike license; entries use original wording
  and cite the subject page rather than adapting text.
- **Unicode code charts** (unicode.org) for glyph identity of Greek letters
  and mathematical operators (Greek and Coptic, Mathematical Operators),
  under the Unicode License Agreement.
- **NIST Dictionary of Algorithms and Data Structures** (DADS,
  xlinux.nist.gov/dads) for big-Theta and big-Omega definitions. A U.S.
  federal reference; entries cite the term page and use original wording
  with no text reused.
- **Encyclopaedia Britannica** (britannica.com) for the aleph-null
  definition. Proprietary content consulted as a reference only; the entry
  uses original wording with no text reused.
- **Wolfram MathWorld** (mathworld.wolfram.com) for the pi continued
  fraction values. Proprietary content consulted as a reference only; the
  entry uses original wording with no text reused.
- **Mathematics LibreTexts** (math.libretexts.org, CC BY-NC-SA) for
  differential-equation terminology and the topology axioms page (the
  latter page carries a public-domain notice, recorded per page).
  Consulted as references; entries use original wording with no text reused.
- **University course notes** (UC San Diego, University of British Columbia,
  University of Chicago authors' pages)
  for generating functions, metric-space completeness, and topology
  formulations. Consulted as references; entries use original wording with
  no text reused.
- **MacTutor History of Mathematics** (University of St Andrews) for
  Fermat's Last Theorem and Diophantine context. Consulted as a reference;
  entries use original wording with no text reused.
- **Stanford Encyclopedia of Philosophy** (open access) for the Turing
  machine definition. Consulted as a reference; the entry uses original
  wording with no text reused.
- **Scholarpedia** (CC BY-NC-SA 3.0 per the article footer) for Shannon
  entropy, cited via its DOI. Consulted as a reference; the entry uses
  original wording with no text reused.
- **PlanetMath** (CC BY-SA, per the site license page) for the
  sigma-algebra definition. Consulted as a reference; the entry uses
  original wording with no text reused.
- **NIST Dictionary of Algorithms and Data Structures** (DADS,
  xlinux.nist.gov/dads) for big-Theta and big-Omega definitions. A U.S.
  federal reference; entries cite the term page and use original wording
  with no text reused.
- **Encyclopaedia Britannica** (britannica.com) for the aleph-null
  definition. Proprietary content consulted as a reference only; the entry
  uses original wording with no text reused.
- **Wolfram MathWorld** (mathworld.wolfram.com) for the pi continued
  fraction values. Proprietary content consulted as a reference only; the
  entry uses original wording with no text reused.
- **Wikipedia** (CC BY-SA 4.0, per article footers) for specialist
  definitions across homological algebra, operator theory, algebraic
  geometry, representation theory, and differential geometry. Individual
  articles opened and verified; entries use original wording with no text
  reused.
- **Stanford Encyclopedia of Philosophy** (open access) for category
  theory definitions alongside the Turing machine entry above. Consulted
  as a reference; entries use original wording with no text reused.

If an entry could not be verified against an appropriate source, it was
excluded rather than published with invented attribution; the excluded
candidates are listed in `web/next-app/scripts/glossary-validation-report.json`,
which is not bundled with the app.