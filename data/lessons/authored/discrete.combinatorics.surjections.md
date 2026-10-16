# Counting onto functions

## Choosing the image set

A function from an \(m\)-element set to an \(n\)-element set is onto when every element of the codomain is hit at least once. Counting them means removing the cases that miss an image.

## Inclusion and exclusion

Start from all \(n^m\) functions and subtract those that miss at least one image. Missing a specific image leaves \((n-1)^m\) functions, there are \(n\) images to miss, and adding back the ones that miss two keeps the accounting right:

\[\sum_{i=0}^{n} (-1)^i \binom{n}{i} (n-i)^m.\]

For \(m = 4\) and \(n = 2\) the value is \(2^4 - 2 \cdot 1^4 = 14\). The sum alternates because a function missing \(k\) images has been subtracted \(k\) times and added back \(\binom{k}{2}\) times.

## The case that matters

The count is zero when \(n > m\), which the formula produces on its own rather than needing a separate rule, and it is \(n!\) when \(m = n\), since an onto function between equal finite sets is a bijection.
