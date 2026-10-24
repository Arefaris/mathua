> Content sourced from [OpenStax Precalculus 2e](https://openstax.org/details/books/precalculus-2e) by Jay Abramson, CC BY-NC-SA 4.0

## 6.1 Graphs of the Sine and Cosine Functions

### Learning Objectives

In this section, you will:

- Graph variations of   $y = \text{sin}(x)$  and $y = \text{cos}(x)$  .
- Use phase shifts of sine and cosine curves.

*Figure 1 Light can be separated into colors because of its wavelike properties. (credit: "wonderferret"/ Flickr)*

White light, such as the light from the sun, is not actually white at all. Instead, it is a composition of all the colors of the rainbow in the form of waves. The individual colors can be seen only when white light passes through an optical prism that separates the waves according to their wavelengths to form a rainbow.

Light waves can be represented graphically by the sine function. In the chapter on Trigonometric Functions, we examined trigonometric functions such as the sine function. In this section, we will interpret and create graphs of sine and cosine functions.

### Graphing Sine and Cosine Functions

Recall that the sine and cosine functions relate real number values to the *x*- and *y*-coordinates of a point on the unit circle. So what do they look like on a graph on a coordinate plane? Let’s start with the sine function. We can create a table of values and use them to sketch a graph. Table 1 lists some of the values for the sine function on a unit circle.

|  |  |  |  |  |  |  |  |  |  |
|----|----|----|----|----|----|----|----|----|----|
| **$x$** | 0 | $\frac{\pi}{6}$ | $\frac{\pi}{4}$ | $\frac{\pi}{3}$ | $\frac{\pi}{2}$ | $\frac{2\pi}{3}$ | $\frac{3\pi}{4}$ | $\frac{5\pi}{6}$ | $\pi$ |
| **$\sin(x)$** | 0 | $\frac{1}{2}$ | $\frac{\sqrt{2}}{2}$ | $\frac{\sqrt{3}}{2}$ | 1 | $\frac{\sqrt{3}}{2}$ | $\frac{\sqrt{2}}{2}$ | $\frac{1}{2}$ | 0 |

Table 1

Plotting the points from the table and continuing along the *x*-axis gives the shape of the sine function. See Figure 2.

*Figure 2 The sine function*

Notice how the sine values are positive between 0 and $\pi,$ which correspond to the values of the sine function in quadrants I and II on the unit circle, and the sine values are negative between $\pi$ and $2\pi,$ which correspond to the values of the sine function in quadrants III and IV on the unit circle. See Figure 3.

*Figure 3 Plotting values of the sine function*

Now let’s take a similar look at the cosine function. Again, we can create a table of values and use them to sketch a graph. Table 2 lists some of the values for the cosine function on a unit circle.

|  |  |  |  |  |  |  |  |  |  |
|----|----|----|----|----|----|----|----|----|----|
| $\mathbf{x}$ | 0 | $\frac{\pi}{6}$ | $\frac{\pi}{4}$ | $\frac{\pi}{3}$ | $\frac{\pi}{2}$ | $\frac{2\pi}{3}$ | $\frac{3\pi}{4}$ | $\frac{5\pi}{6}$ | $\pi$ |
| $\mathbf{\cos}\left( \mathbf{x} \right)$ | 1 | $\frac{\sqrt{3}}{2}$ | $\frac{\sqrt{2}}{2}$ | $\frac{1}{2}$ | 0 | $- \frac{1}{2}$ | $- \frac{\sqrt{2}}{2}$ | $- \frac{\sqrt{3}}{2}$ | $- 1$ |

Table 2

As with the sine function, we can plots points to create a graph of the cosine function as in Figure 4.

*Figure 4 The cosine function*

Because we can evaluate the sine and cosine of any real number, both of these functions are defined for all real numbers. By thinking of the sine and cosine values as coordinates of points on a unit circle, it becomes clear that the range of both functions must be the interval $\left\lbrack {- 1,1} \right\rbrack.$

In both graphs, the shape of the graph repeats after $2\pi,$ which means the functions are periodic with a period of $2\pi.$ A periodic function is a function for which a specific horizontal shift, *P*, results in a function equal to the original function: $f\left( {x + P} \right) = f(x)$ for all values of $x$ in the domain of $f.$ When this occurs, we call the smallest such horizontal shift with $P > 0$ the period of the function. Figure 5 shows several periods of the sine and cosine functions.

*Figure 5*

Looking again at the sine and cosine functions on a domain centered at the *y*-axis helps reveal symmetries. As we can see in Figure 6, the sine function is symmetric about the origin. Recall from The Other Trigonometric Functions that we determined from the unit circle that the sine function is an odd function because $\sin(-x) = -\sin\; x.$ Now we can clearly see this property from the graph.

*Figure 6 Odd symmetry of the sine function*

Figure 7 shows that the cosine function is symmetric about the *y*-axis. Again, we determined that the cosine function is an even function. Now we can see from the graph that $\cos(-x) = \cos\; x.$

*Figure 7 Even symmetry of the cosine function*

### Characteristics of Sine and Cosine Functions

The sine and cosine functions have several distinct characteristics:

- They are periodic functions with a period of $2\pi.$
- The domain of each function is $\left( {- \infty,\infty} \right)$ and the range is $\left\lbrack {- 1,1} \right\rbrack.$
- The graph of $y = \sin\; x$ is symmetric about the origin, because it is an odd function.
- The graph of $y = \cos\; x$ is symmetric about the $y$-axis, because it is an even function.

### Investigating Sinusoidal Functions

As we can see, sine and cosine functions have a regular period and range. If we watch ocean waves or ripples on a pond, we will see that they resemble the sine or cosine functions. However, they are not necessarily identical. Some are taller or longer than others. A function that has the same general shape as a sine or cosine function is known as a sinusoidal function. The general forms of sinusoidal functions are

$$
\begin{array}{l}
{y = A\sin\left( {Bx - C} \right) + D} \\
\text{and} \\
{y = A\cos\left( {Bx - C} \right) + D}
\end{array}
$$

#### Determining the Period of Sinusoidal Functions

Looking at the forms of sinusoidal functions, we can see that they are transformations of the sine and cosine functions. We can use what we know about transformations to determine the period.

In the general formula, $B$ is related to the period by $P = \frac{2\pi}{|B|}.$ If $|B| > 1,$ then the period is less than $2\pi$ and the function undergoes a horizontal compression, whereas if $|B| < 1,$ then the period is greater than $2\pi$ and the function undergoes a horizontal stretch. For example, $f(x) = \sin\left( x\text{),} \right.$ $B = 1,$ so the period is $2\pi,$ which we knew. If $f(x) = \sin\left( {2x} \right),$ then $B = 2,$ so the period is $\pi$ and the graph is compressed. If $f(x) = \sin\left( \frac{x}{2} \right),$ then $B = \frac{1}{2},$ so the period is $4\pi$ and the graph is stretched. Notice in Figure 8 how the period is indirectly related to $|B|.$

*Figure 8*

### Period of Sinusoidal Functions

If we let $C = 0$ and $D = 0$ in the general form equations of the sine and cosine functions, we obtain the forms

$$
y = A\sin\left( {Bx} \right)
$$

$$
y = A\cos\left( {Bx} \right)
$$

The period is $\frac{2\pi}{|B|}.$

### Example 1

#### Identifying the Period of a Sine or Cosine Function

Determine the period of the function $f(x) = \sin\left( {\frac{\pi}{6}x} \right).$

#### Solution

Let’s begin by comparing the equation to the general form $y = A\sin(Bx).$

In the given equation, $B = \frac{\pi}{6},$ so the period will be

$$
\begin{array}{l}
\begin{array}{l}
 \\
{P = \frac{2\pi}{|B|}}
\end{array} \\
{= \frac{2\pi}{\frac{\pi}{6}}} \\
{= 2\pi \cdot \frac{6}{\pi}} \\
{= 12}
\end{array}
$$

### Try It \#1

Determine the period of the function $g(x) = \cos\left( \frac{x}{3} \right).$

#### Determining Amplitude

Returning to the general formula for a sinusoidal function, we have analyzed how the variable $B$ relates to the period. Now let’s turn to the variable $A$ so we can analyze how it is related to the **amplitude**, or greatest distance from rest. $A$ represents the vertical stretch factor, and its absolute value $|A|$ is the amplitude. The local maxima will be a distance $|A|$ above the horizontal **midline** of the graph, which is the line $y = D;$ because $D = 0$ in this case, the midline is the *x*-axis. The local minima will be the same distance below the midline. If $|A| > 1,$ the function is stretched. For example, the amplitude of $f(x) = 4\;\sin\; x$ is twice the amplitude of
$$
f(x) = 2\;\sin\; x.
$$

If $|A| < 1,$ the function is compressed. Figure 9 compares several sine functions with different amplitudes.

*Figure 9*

### Amplitude of Sinusoidal Functions

If we let $C = 0$ and $D = 0$ in the general form equations of the sine and cosine functions, we obtain the forms

$$
y = A\sin\left( {Bx} \right)\text{~and~}y = A\cos\left( {Bx} \right)
$$

The amplitude is $\left| A\text{|}, \right.$ which is the vertical height from the midline $.$ In addition, notice in the example that

$$
|A|\text{~=~amplitude~=~}\frac{1}{2}\left| {\text{maximum~} - \text{~minimum}} \right|
$$

### Example 2

#### Identifying the Amplitude of a Sine or Cosine Function

What is the amplitude of the sinusoidal function $f(x) = -4\sin(x)?$ Is the function stretched or compressed vertically?

#### Solution

Let’s begin by comparing the function to the simplified form $y = A\sin(Bx).$

In the given function, $A = -4,$ so the amplitude is $|A| = |-4| = 4.$ The function is stretched.

#### Analysis

The negative value of $A$ results in a reflection across the *x*-axis of the sine function, as shown in Figure 10.

*Figure 10*

### Try It \#2

What is the amplitude of the sinusoidal function $f(x) = \frac{1}{2}\sin(x)?$ Is the function stretched or compressed vertically?

### Analyzing Graphs of Variations of *y* = sin *x* and *y* = cos *x*

Now that we understand how $A$ and $B$ relate to the general form equation for the sine and cosine functions, we will explore the variables $C$ and $D.$ Recall the general form:

$$
\begin{matrix}
{y = A\sin\left( {Bx - C} \right) + D\text{~and~}y = A\cos\left( {Bx - C} \right) + D} \\
{or} \\
{y = A\sin\left( {B\left( {x - \frac{C}{B}} \right)} \right) + D\text{~and~}y = A\cos\left( {B\left( {x - \frac{C}{B}} \right)} \right) + D}
\end{matrix}
$$

The value $\frac{C}{B}$ for a sinusoidal function is called the **phase shift**, or the horizontal displacement of the basic sine or cosine function. If $C > 0,$ the graph shifts to the right. If $C < 0,$ the graph shifts to the left. The greater the value of $|C|,$ the more the graph is shifted. Figure 11 shows that the graph of $f(x) = \sin\left( {x - \pi} \right)$ shifts to the right by $\pi$ units, which is more than we see in the graph of $f(x) = \sin\left( {x - \frac{\pi}{4}} \right),$ which shifts to the right by $\frac{\pi}{4}$ units.

*Figure 11*

While $C$ relates to the horizontal shift, $D$ indicates the vertical shift from the midline in the general formula for a sinusoidal function. See Figure 12. The function $y = \cos(x) + D$ has its midline at $y = D.$

*Figure 12*

Any value of $D$ other than zero shifts the graph up or down. Figure 13 compares $f(x) = \sin\;(x)$ with $f(x) = \sin\;(x) + 2,$ which is shifted 2 units up on a graph.

*Figure 13*

### Variations of Sine and Cosine Functions

Given an equation in the form $f(x) = A\sin\left( {Bx - C} \right) + D$ or $f(x) = A\cos\left( {Bx - C} \right) + D,$ $\frac{C}{B}$ is the phase shift and $D$ is the vertical shift.

### Example 3

#### Identifying the Phase Shift of a Function

Determine the direction and magnitude of the phase shift for $f(x) = \sin\left( {x + \frac{\pi}{6}} \right) - 2.$

#### Solution

Let’s begin by comparing the equation to the general form $y = A\sin(Bx - C) + D.$

In the given equation, notice that $B = 1$ and $C = - \frac{\pi}{6}.$ So the phase shift is

$$
\begin{array}{r}
 \\
{\frac{C}{B} = - \frac{\frac{\pi}{6}}{1}} \\
{= - \frac{\pi}{6}}
\end{array}
$$

or $\frac{\pi}{6}$ units to the left.

#### Analysis

We must pay attention to the sign in the equation for the general form of a sinusoidal function. The equation shows a minus sign before $C.$ Therefore $f(x) = \sin\left( {x + \frac{\pi}{6}} \right) - 2$ can be rewritten as $f(x) = \sin\left( {x - \left( {- \frac{\pi}{6}} \right)} \right) - 2.$ If the value of $C$ is negative, the shift is to the left.

### Try It \#3

Determine the direction and magnitude of the phase shift for $f(x) = 3\cos\left( {x - \frac{\pi}{2}} \right).$

### Example 4

#### Identifying the Vertical Shift of a Function

Determine the direction and magnitude of the vertical shift for $f(x) = \cos(x) - 3.$

#### Solution

Let’s begin by comparing the equation to the general form $y = A\cos(Bx - C) + D.$

In the given equation, $D = -3$ so the shift is 3 units downward.

### Try It \#4

Determine the direction and magnitude of the vertical shift for $f(x) = 3\sin(x) + 2.$

### How To

**Given a sinusoidal function in the form** $f(x) = A\sin\left( {Bx - C} \right) + D,$ **identify the midline, amplitude, period, and phase shift.**

1.  Determine the amplitude as $|A|.$
2.  Determine the period as $P = \frac{2\pi}{|B|}.$
3.  Determine the phase shift as $\frac{C}{B}.$
4.  Determine the midline as $y = D.$

### Example 5

#### Identifying the Variations of a Sinusoidal Function from an Equation

Determine the midline, amplitude, period, and phase shift of the function $y = 3\sin(2x) + 1.$

#### Solution

Let’s begin by comparing the equation to the general form $y = A\sin(Bx - C) + D.$

$A = 3,$ so the amplitude is $|A| = 3.$

Next, $B = 2,$ so the period is $P = \frac{2\pi}{|B|} = \frac{2\pi}{2} = \pi.$

There is no added constant inside the parentheses, so $C = 0$ and the phase shift is $\frac{C}{B} = \frac{0}{2} = 0.$

Finally, $D = 1,$ so the midline is $y = 1.$

#### Analysis

Inspecting the graph, we can determine that the period is $\pi,$ the midline is $y = 1,$ and the amplitude is 3. See Figure 14.

*Figure 14*

### Try It \#5

Determine the midline, amplitude, period, and phase shift of the function $y = \frac{1}{2}\cos\left( {\frac{x}{3} - \frac{\pi}{3}} \right).$

### Example 6

#### Identifying the Equation for a Sinusoidal Function from a Graph

Determine the formula for the cosine function in Figure 15.

*Figure 15*

#### Solution

To determine the equation, we need to identify each value in the general form of a sinusoidal function.

$$
\begin{array}{l}
{y = A\sin(Bx - C) + D} \\
{y = A\cos(Bx - C) + D}
\end{array}
$$

The graph could represent either a sine or a cosine function that is shifted and/or reflected. When $x = 0,$ the graph has an extreme point, $\left( {0,0} \right).$ Since the cosine function has an extreme point for $x = 0,$ let us write our equation in terms of a cosine function.

Let’s start with the midline. We can see that the graph rises and falls an equal distance above and below $y = 0.5.$ This value, which is the midline, is $D$ in the equation, so $D = 0.5.$

The greatest distance above and below the midline is the amplitude. The maxima are 0.5 units above the midline and the minima are 0.5 units below the midline. So $|A| = 0.5.$ Another way we could have determined the amplitude is by recognizing that the difference between the height of local maxima and minima is 1, so $|A| = \frac{1}{2} = 0.5.$ Also, the graph is reflected about the *x*-axis so that $A = - 0.5.$

The graph is not horizontally stretched or compressed, so $B = 1;$ and the graph is not shifted horizontally, so $C = 0.$

Putting this all together,

$$
g(x) = - 0.5\cos(x) + 0.5
$$

### Try It \#6

Determine the formula for the sine function in Figure 16.

*Figure 16*

### Example 7

#### Identifying the Equation for a Sinusoidal Function from a Graph

Determine the equation for the sinusoidal function in Figure 17.

*Figure 17*

#### Solution

With the highest value at 1 and the lowest value at $-5,$ the midline will be halfway between at $-2.$ So $D = -2.$

The distance from the midline to the highest or lowest value gives an amplitude of $|A| = 3.$

The period of the graph is 6, which can be measured from the peak at $x = 1$ to the next peak at $x = 7,$ or from the distance between the lowest points. Therefore, $P = \frac{2\pi}{|B|} = 6.$ Using the positive value for $B,$ we find that

$$
B = \frac{2\pi}{P} = \frac{2\pi}{6} = \frac{\pi}{3}
$$

So far, our equation is either $y = 3\sin\left( {\frac{\pi}{3}x - C} \right) - 2$ or $y = 3\cos\left( {\frac{\pi}{3}x - C} \right) - 2.$ For the shape and shift, we have more than one option. We could write this as any one of the following:

- a cosine shifted to the right
- a negative cosine shifted to the left
- a sine shifted to the left
- a negative sine shifted to the right

Choosing to use the cosine function, we observe that the peak, which would normally be at $x = 0$, is at $x = 1$, and given the horizontal compression factor of $\frac{\pi}{3}$, we get $C = 1 \cdot \frac{\pi}{3} = \frac{\pi}{3}$.

While any of these would be correct, the cosine shifts are easier to work with than the sine shifts in this case because they involve integer values. So our function becomes

$$
y = 3\cos\left( {\frac{\pi}{3}x - \frac{\pi}{3}} \right) - 2\text{~or~}y = - 3\cos\left( {\frac{\pi}{3}x + \frac{2\pi}{3}} \right) - 2
$$

Again, these functions are equivalent, so both yield the same graph.

### Try It \#7

Write a formula for the function graphed in Figure 18.

*Figure 18*

### Graphing Variations of *y* = sin *x* and *y* = cos *x*

Throughout this section, we have learned about types of variations of sine and cosine functions and used that information to write equations from graphs. Now we can use the same information to create graphs from equations.

Instead of focusing on the general form equations

$$
y = A\sin\left( {Bx - C} \right) + D\text{~and~}y = A\cos\left( {Bx - C} \right) + D,
$$

we will let $C = 0$ and $D = 0$ and work with a simplified form of the equations in the following examples.

### How To

**Given the function $y = A\sin\left( {Bx} \right),$ sketch its graph.**

1.  Identify the amplitude, $|A|.$
2.  Identify the period, $P = \frac{2\pi}{|B|}.$
3.  Start at the origin, with the function increasing to the right if $A$ is positive or decreasing if $A$ is negative.
4.  At $x = \frac{\pi}{2|B|}$ there is a local maximum for $A > 0$ or a minimum for $A < 0,$ with $y = A.$
5.  The curve returns to the *x*-axis at $x = \frac{\pi}{|B|}.$
6.  There is a local minimum for $A > 0$ (maximum for $A < 0$ ) at $x = \frac{3\pi}{2|B|}$ with $y = -A.$
7.  The curve returns again to the *x*-axis at $x = \frac{2\pi}{|B|}.$

CONTINUE HERE

### Example 8

#### Graphing a Function and Identifying the Amplitude and Period

Sketch a graph of $f(x) = - 2\sin\left( \frac{\pi x}{2} \right).$

#### Solution

Let’s begin by comparing the equation to the form $y = A\sin(Bx).$

- *Step 1.* We can see from the equation that $A = - 2,$ so the amplitude is 2.
  ``` math
  |A| = 2
  ```
- *Step 2.* The equation shows that $B = \frac{\pi}{2},$ so the period is
  ``` math
  \begin{array}{l}
  {P = \frac{2\pi}{\frac{\pi}{2}}} \\
  {= 2\pi \cdot \frac{2}{\pi}} \\
  {= 4}
  \end{array}
  ```
- *Step 3.* Because $A$ is negative, the graph descends as we move to the right of the origin.
- *Step 4–7.* The *x*-intercepts are at the beginning of one period, $x = 0,$ the horizontal midpoints are at $x = 2$ and at the end of one period at $x = 4.$

The quarter points include the minimum at $x = 1$ and the maximum at $x = 3.$ A local minimum will occur 2 units below the midline, at $x = 1,$ and a local maximum will occur at 2 units above the midline, at $x = 3.$ Figure 19 shows the graph of the function.

*Figure 19*

### Try It \#8

Sketch a graph of $g(x) = - 0.8\cos\left( {2x} \right).$ Determine the midline, amplitude, period, and phase shift.

### How To

**Given a sinusoidal function with a phase shift and a vertical shift, sketch its graph.**

1.  Express the function in the general form $y = A\sin(Bx - C) + D\text{~or~}y = A\cos(Bx - C) + D.$
2.  Identify the amplitude, $|A|.$
3.  Identify the period, $P = \frac{2\pi}{|B|}.$
4.  Identify the phase shift, $\frac{C}{B}.$
5.  Draw the graph of $f(x) = A\sin\left( {Bx} \right)$ shifted to the right or left by $\frac{C}{B}$ and up or down by $D.$

### Example 9

#### Graphing a Transformed Sinusoid

Sketch a graph of $f(x) = 3\sin\left( {\frac{\pi}{4}x - \frac{\pi}{4}} \right).$

#### Solution

- *Step 1.* The function is already written in general form: $f(x) = 3\sin\left( {\frac{\pi}{4}x - \frac{\pi}{4}} \right).$ This graph will have the shape of a sine function, starting at the midline and increasing to the right.

- *Step 2.* $|A| = |3| = 3.$ The amplitude is 3.

- *Step 3.* Since $|B| = \left| \frac{\pi}{4} \right| = \frac{\pi}{4},$ we determine the period as follows.
  ``` math
  P = \frac{2\pi}{|B|} = \frac{2\pi}{\frac{\pi}{4}} = 2\pi \cdot \frac{4}{\pi} = 8
  ```

  The period is 8.

- *Step 4.* Since $C = \frac{\pi}{4},$ the phase shift is
  ``` math
  \frac{C}{B} = \frac{\frac{\pi}{4}}{\frac{\pi}{4}} = 1.
  ```

  The phase shift is 1 unit.

- *Step 5.* Figure 20 shows the graph of the function.

  *Figure 20 A horizontally compressed, vertically stretched, and horizontally shifted sinusoid*

### Try It \#9

Draw a graph of $g(x) = - 2\cos\left( {\frac{\pi}{3}x + \frac{\pi}{6}} \right).$ Determine the midline, amplitude, period, and phase shift.

### Example 10

#### Identifying the Properties of a Sinusoidal Function

Given $y = - 2\cos\left( {\frac{\pi}{2}x + \pi} \right) + 3,$ determine the amplitude, period, phase shift, and vertical shift. Then graph the function.

#### Solution

Begin by comparing the equation to the general form and use the steps outlined in Example 9.

$$
y = A\cos\left( {Bx - C} \right) + D
$$

- *Step 1.* The function is already written in general form.
- *Step 2.* Since $A = - 2,$ the amplitude is $|A| = 2.$
- *Step 3.* $|B| = \frac{\pi}{2},$ so the period is $P = \frac{2\pi}{|B|} = \frac{2\pi}{\frac{\pi}{2}} = 2\pi \cdot \frac{2}{\pi} = 4.$ The period is 4.
- *Step 4.* $C = - \pi,$ so we calculate the phase shift as $\frac{C}{B} = \frac{- \pi,}{\frac{\pi}{2}} = - \pi \cdot \frac{2}{\pi} = - 2.$ The phase shift is $- 2.$
- *Step 5.* $D = 3,$ so the midline is $y = 3\text{,\,}$ and the vertical shift is up 3.

Since $A$ is negative, the graph of the cosine function has been reflected about the *x*-axis.

Figure 21 shows one cycle of the graph of the function.

*Figure 21*

### Using Transformations of Sine and Cosine Functions

We can use the transformations of sine and cosine functions in numerous applications. As mentioned at the beginning of the chapter, circular motion can be modeled using either the sine or cosine function.

### Example 11

#### Finding the Vertical Component of Circular Motion

A point rotates around a circle of radius 3 centered at the origin. Sketch a graph of the *y*-coordinate of the point as a function of the angle of rotation.

#### Solution

Recall that, for a point on a circle of radius *r*, the *y*-coordinate of the point is $y = r\;\sin(x),$ so in this case, we get the equation $y(x) = 3\;\sin(x).$ The constant 3 causes a vertical stretch of the *y*-values of the function by a factor of 3, which we can see in the graph in Figure 22.

*Figure 22*

#### Analysis

Notice that the period of the function is still $2\pi;$ as we travel around the circle, we return to the point $\left( {3,0} \right)$ for $x = 2\pi,4\pi,6\pi,...$ Because the outputs of the graph will now oscillate between $-3$ and 3, the amplitude of the sine wave is 3.

### Try It \#10

What is the amplitude of the function $f(x) = 7\cos(x)?$ Sketch a graph of this function.

### Example 12

#### Finding the Vertical Component of Circular Motion

A circle with radius 3 ft is mounted with its center 4 ft off the ground. The point closest to the ground is labeled *P*, as shown in Figure 23. Sketch a graph of the height above the ground of the point $P$ as the circle is rotated; then find a function that gives the height in terms of the angle of rotation.

*Figure 23*

#### Solution

Sketching the height, we note that it will start 1 ft above the ground, then increase up to 7 ft above the ground, and continue to oscillate 3 ft above and below the center value of 4 ft, as shown in Figure 24.

*Figure 24*

Although we could use a transformation of either the sine or cosine function, we start by looking for characteristics that would make one function easier to use than the other. Let’s use a cosine function because it starts at the highest or lowest value, while a sine function starts at the middle value. A standard cosine starts at the highest value, and this graph starts at the lowest value, so we need to incorporate a vertical reflection.

Second, we see that the graph oscillates 3 above and below the center, while a basic cosine has an amplitude of 1, so this graph has been vertically stretched by 3, as in the last example.

Finally, to move the center of the circle up to a height of 4, the graph has been vertically shifted up by 4. Putting these transformations together, we find that

$$
y = - 3\cos(x) + 4
$$

### Try It \#11

A weight is attached to a spring that is then hung from a board, as shown in Figure 25. As the spring oscillates up and down, the position $y$ of the weight relative to the board ranges from $-1$ in. (at time $x = 0)$ to $-7$ in. (at time $x = \pi)$ below the board. Assume the position of $y$ is given as a sinusoidal function of $x.$ Sketch a graph of the function, and then find a cosine function that gives the position $y$ in terms of $x.$

*Figure 25*

### Example 13

#### Determining a Rider’s Height on a Ferris Wheel

The London Eye is a huge Ferris wheel with a diameter of 135 meters (443 feet). It completes one rotation every 30 minutes. Riders board from a platform 2 meters above the ground. Express a rider’s height above ground as a function of time in minutes.

#### Solution

With a diameter of 135 m, the wheel has a radius of 67.5 m. The height will oscillate with amplitude 67.5 m above and below the center.

Passengers board 2 m above ground level, so the center of the wheel must be located $67.5 + 2 = 69.5$ m above ground level. The midline of the oscillation will be at 69.5 m.

The wheel takes 30 minutes to complete 1 revolution, so the height will oscillate with a period of 30 minutes.

Lastly, because the rider boards at the lowest point, the height will start at the smallest value and increase, following the shape of a vertically reflected cosine curve.

- Amplitude: 67.5, so $A = 67.5$
- Midline: 69.5, so $D = 69.5$
- Period: 30, so $B = \frac{2\pi}{30} = \frac{\pi}{15}$
- Shape: ${-cos}(t)$

An equation for the rider’s height would be

$$
y = - 67.5\cos\left( {\frac{\pi}{15}t} \right) + 69.5
$$

where $t$ is in minutes and $y$ is measured in meters.
