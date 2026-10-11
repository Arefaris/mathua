# Glossary audit — research notes (working artifact, keep until approved)

Each batch records: claims checked, sources actually opened (URL + what was
verified there), and findings. "Searched" is not "verified": an entry counts
as source-verified only when a claim was checked against an opened source.
Render-testing is covered separately by the KaTeX gate
(`test/glossaryData.test.ts` "renders every example through KaTeX") and the
Playwright spec; the ledger distinguishes the two.

## Batch: advanced (27 entries)

### Research log (all pages opened, not snippets)
- https://dlmf.nist.gov/5.2 (DLMF gamma definitions): 5.2.1 Euler integral
  for Gamma; analytic-continuation paragraph. Supports `gamma-function`
  (integral definition; integer values follow the standard recurrence).
- https://dlmf.nist.gov/5.12 (DLMF beta function): 5.12.1
  B(a,b) = Gamma(a)Gamma(b)/Gamma(a+b). Exact confirmation of the
  `beta-function` example.
- https://dlmf.nist.gov/7.2 (DLMF error functions): 7.2.1
  erf z = (2/sqrt(pi)) integral_0^z e^{-t^2} dt. Exact confirmation of the
  `error-function` example.
- https://dlmf.nist.gov/10.4 and https://dlmf.nist.gov/10.8 (DLMF Bessel):
  connection formulas and power-series framework (10.8 references the J
  series 10.2.2). J_0(0) = 1 follows from the defining series (k=0 term is
  1/(0! Gamma(1)) = 1); recorded as derivation, not lookup.
- https://mathworld.wolfram.com/BetaFunction.html (opened): on-topic page
  citing A&S section 6.2; equation glyphs did not render in text fetch, so
  the relation claim rests on DLMF 5.12.1 above.
- Number-set, group, isomorphism, congruence, gcd/lcm, totient, tensor,
  direct-sum, argmin, Fourier, Laplace and big-O examples all recomputed by
  hand: phi(9) = 6 (1,2,4,5,7,8); (R,+) iso (R>0,x) via exp; L{f'} formula;
  all correct.

### Findings and changes
- No corrections required. Status: verified.

## Batch: linear-discrete + probability (37 entries)

### Research log (all pages opened, not snippets)
- https://openstax.org/books/introductory-statistics/pages/6-1-the-standard-normal-distribution
  (OpenStax Introductory Statistics, section 6.1): z = (x-mu)/sigma;
  Z ~ N(0,1); Empirical Rule 68/95/99.7%. Confirms `z-score`, `normal-distribution`
  (modulo the convention note below) and `standard-deviation` entries.
- Same page footer (licensing): "Except where otherwise noted, textbooks on
  this site are licensed under CC BY-NC-SA 4.0", plus an LLM-training
  restriction (not applicable: no text reused anywhere in the glossary).
  This corrects the previously recorded "CC BY 4.0 (OpenStax)" on 105
  source lines across the dataset and in THIRD_PARTY_NOTICES.md.
- Linear-algebra examples (transpose, 2x2 determinant, trace 5, rank 1,
  eigenvalue equation, eigenvector check, inner/dot/cross products,
  Kronecker sift, P(5,2) = 20, C(5,2) = 10, mod remainders, floor/ceiling)
  and probability examples (die mean 3.5, Var identity, z = 2, Binomial pmf,
  p-hat 0.6, y-hat 14, p = 0.03 decision) recomputed by hand: all correct.

### Findings and changes
- `normal-distribution`: added a convention note — probability theory
  writes N(mu, sigma-squared) with the variance, while OpenStax and some
  texts write N(mu, sigma) with the standard deviation. The entry previously
  presented one convention as universal.
- Licensing: 105 source license strings corrected to CC BY-NC-SA 4.0 with
  "consulted as reference, original wording, no text reused"; ledger entry
  updated. No text was ever copied, so no ShareAlike trigger; the record is
  now accurate.
- No other corrections. Status: verified (`normal-distribution`:
  corrected and verified).

## Batch: geometry-trig + calculus (42 entries)

### Research log (all pages opened, not snippets)
- https://openstax.org/books/calculus-volume-3/pages/6-5-divergence-and-curl
  (OpenStax Calculus Vol 3, section 6.5): div F = Px+Qy+Rz (scalar);
  curl F = (Ry-Qz)i + (Pz-Rx)j + (Qx-Py)k (vector); curl-as-rotation with
  paddlewheel interpretation; div(curl F) = 0. Confirms the `divergence`
  example (1+1 = 2), the `curl` example (1-(-1) = 2k), and both definitions.
- https://tutorial.math.lamar.edu/classes/calciii/curldivergence.aspx
  (Lamar University, Paul's Online Notes): nabla operator, curl as cross
  product, divergence as dot product. Second confirmation of the same
  formulas and of the `nabla-gradient` framing.
- Geometry/trig examples (3-4-5 distance, midpoint (2,3), circle center and
  radius, sin30/cos60/tan45, arcsin principal value, vector length 5) and
  all calculus examples (limits, derivatives, partials, integrals,
  epsilon-delta statement, sup/inf of (0,1), Euler limit, Riemann limit)
  recomputed by hand during review: all correct.

### Findings and changes
- No corrections required. Every definition, example, explanation and
  cross-link reviewed field by field. Status: verified.

## Batch: algebra + logic-sets (41 entries)

### Research log (all pages opened, not snippets)
- Metamath Proof Explorer `qnnen` (formal proof, `ℚ ≈ ℕ`), Math SE #1311
  (rationals equicardinal with integers/naturals), and a cardinality
  exposition with explicit bijections: all confirm |N| = |Q|. Supports the
  `cardinality` example `|N| = |Q|`.
- Algebra examples (quadratic formula roots x = -1, -2 for x^2+3x+2;
  summation 1+4+9+16 = 30; log_2(8) = 3; arcsin(1/2) = 30 degrees principal
  value; (a+b)^3 expansion) verified by direct computation during review.

### Findings and changes
- `disjunction`: example explanation suggested "both conditions" could hold
  for (x=0) v (x=1), which is impossible there. Reworded: inclusive-or allows
  both in general, though not in this instance.
- Remaining algebra entries (function notation, inverse, composition,
  mapping arrow, exp/log/ln, sequences, summation, product, binomial,
  proportion, delta, ellipsis, slope-intercept, quadratic formula) and logic
  entries (membership, subsets, union/intersection/difference/complement,
  empty and power sets, set-builder, quantifiers, connectives, implication,
  iff, therefore/because, cardinality, QED, predicates): every field
  reviewed; examples recomputed; cross-links checked. Status: verified
  (`disjunction`: corrected and verified).

## Batch: greek + arithmetic (56 entries)

### Research log (all pages opened, not snippets)
- https://dlmf.nist.gov/25.6 (DLMF, Riemann zeta integer values): confirms
  25.6.1 zeta(2) = pi^2/6. Supports `zeta` example and `zeta-function` example.
- https://mathworld.wolfram.com/MangoldtFunction.html (MathWorld, specialist
  reference): von Mangoldt function "sometimes also called the lambda
  function", denoted by capital lambda. Supports `lambda` correction.
- https://leanprover-community.github.io/mathlib4_docs/Mathlib/NumberTheory/ArithmeticFunction/VonMangoldt.html
  (Lean Mathlib docs): "We use the standard notation Λ." Second confirmation.
- https://www.tauday.com/tau-manifesto (Hartl, primary source of the proposal)
  and Wikipedia "Tau (mathematics)": tau-as-circle-constant is a 2010
  proposal; "π is used almost exclusively in mainstream mathematical
  education and practice". Supports `tau` correction.
- Wikipedia "Correlation ratio" + https://encyclopediaofmath.org/wiki/Correlation_ratio
  (specialist encyclopedia): correlation ratio denoted eta. Supports `eta`
  correction (was "partial regression coefficient").

### Findings and changes
- `beta`: claimed lowercase beta denotes "the beta function" — wrong, that is
  capital B. Reworded to regression coefficients / beta-distribution shape.
  Related link `beta-function` replaced with `regression-note`.
- `eta`: "partial regression coefficient" unverified — replaced with
  correlation ratio (eta), verified above.
- `lambda`: "uppercase Lambda is used in logic and set theory" unverified —
  replaced with von Mangoldt function (number theory), verified above.
- `tau`: presented the tau proposal as standard usage — qualified as a
  proposal; mainstream uses pi. Example explanation qualified.
- `upsilon`: overstated mathematical use — corrected to rarely used.
- `alpha`: "root of an equation" is unusual — replaced with regression
  intercept / significance level.
- `power`: x^{-1} = 1/x now states nonzero x.
- `nth-root`: x^{1/n} now states the domain qualification.
- Remaining Greek letters (gamma, delta, epsilon, theta, iota, kappa, mu, nu,
  xi, omicron, pi, rho, sigma, phi, chi, psi, omega) and all other arithmetic
  entries: symbol, definition, examples and cross-links reviewed field by
  field against the entries' cited sources (Unicode charts for glyph identity,
  OpenStax Prealgebra/Elementary/Intermediate Algebra for meaning); no issues.
  Status: verified.
