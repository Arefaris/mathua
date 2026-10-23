package odes

import (
	"fmt"
	"math/rand"

	"github.com/chuma-beep/mathua/internal/generator"
)

func Register(reg *generator.Registry) {
	reg.Register("ode.basics.concept", &conceptGen{})
	reg.Register("ode.first_order.separable", &separableGen{})
	reg.Register("ode.first_order.linear", &linearFirstGen{})
	reg.Register("ode.first_order.exact", &exactGen{})
	reg.Register("ode.second_order.homogeneous", &homogeneousGen{})
	reg.Register("ode.second_order.nonhomogeneous", &nonhomogeneousGen{})
	reg.Register("ode.transforms.laplace", &laplaceGen{})
	reg.Register("ode.adv.systems", &systemsGen{})
	reg.Register("ode.second_order.variation_params", &variationParamsGen{})
	reg.Register("ode.transforms.convolution", &convolutionGen{})
	reg.Register("ode.systems.linear_phase", &linearPhaseGen{})
	reg.Register("ode.numerical.euler", &eulerGen{})
	reg.Register("ode.numerical.rk4", &rk4Gen{})
	reg.Register("ode.stability", &stabilityGen{})
	reg.Register("ode.pde.characteristics", &characteristicsGen{})
	reg.Register("ode.transforms.fourier", &fourierGen{})
	reg.Register("ode.nonlinear.bifurcation", &bifurcationGen{})
	reg.Register("ode.pde.wave", &waveGen{})
	reg.Register("ode.pde.heat", &heatGen{})
	reg.Register("ode.green.function", &greenFunctionGen{})
	reg.Register("ode.sturm_liouville", &sturmLiouvilleGen{})
	reg.Register("ode.boundary_value", &boundaryValueGen{})
	reg.Register("ode.nonlinear.lyapunov", &lyapunovGen{})
	reg.Register("ode.pde.laplace_eq", &laplaceEqGen{})
	reg.Register("ode.integral.volterra", &volterraGen{})
	registerNumeric(reg)
	registerBatch(reg)
}

type conceptGen struct{}

func (g *conceptGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	type entry struct {
		q string
		a string
		e string
	}
	easy := []entry{
		{q: "Is \\(y' = 2x\\) a differential equation? (yes/no)", a: "yes", e: "The equation contains a derivative \\(y'\\), so it is a differential equation."},
		{q: "Is \\(x^{2} + y^{2} = 5\\) a differential equation? (yes/no)", a: "no", e: "The equation contains no derivatives, so it is not a differential equation."},
		{q: "Does the equation \\(x^{3} - y^{3} = 10\\) contain a derivative? (yes/no)", a: "no", e: "The equation contains no derivative terms."},
	}
	orders := []struct{ q, a, e string }{
		{"What is the order of \\(y' + 2y = 0\\)? (enter a number)", "1", "The highest derivative is \\(y'\\), first order."},
		{"What is the order of \\(y'' + y = 0\\)? (enter a number)", "2", "The highest derivative is \\(y''\\), second order."},
		{"What is the order of \\(y''' + 2y'' - y' + y = 0\\)? (enter a number)", "3", "The highest derivative is \\(y'''\\), third order."},
		{"What is the order of \\(dy/dx = 3x^{2}\\)? (enter a number)", "1", "Only a first derivative appears."},
		{"What is the order of \\(d^{2}y/dx^{2} + y = 0\\)? (enter a number)", "2", "The highest derivative is second order."},
	}
	o := orders[rand.Intn(len(orders))]
	easy = append(easy, entry{q: o.q, a: o.a, e: o.e})
	hard := []entry{
		{q: "Is \\(\\sin(x) + \\cos(y) = 1\\) a differential equation? (yes/no)", a: "no", e: "The equation contains no derivatives, so it is not a differential equation."},
		{q: "Is \\(3x + 2y = 7\\) a differential equation? (yes/no)", a: "no", e: "The equation contains no derivatives, so it is not a differential equation."},
		{q: "Does \\(y' = 2x\\) have order 1? (yes/no)", a: "yes", e: "Only the first derivative appears."},
	}
	pool := easy
	if scale > 3 {
		pool = append(easy, hard...)
	}
	e := pool[rand.Intn(len(pool))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type separableGen struct{}

func (g *separableGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type qa struct {
		q string
		a string
		e string
	}
	templates := []qa{
		{
			q: "Solve: \\(dy/dx = x/y\\)",
			a: "y^2=x^2+C",
			e: "Separate: \\(y dy = x dx\\). Integrate: \\(y^{2}/2 = x^{2}/2 + C\\). Multiply by 2: \\(y^{2} = x^{2} + C\\).",
		},
		{
			q: "Solve: \\(dy/dx = x\\)",
			a: "y=x^2/2+C",
			e: "Rewrite as \\(dy = x dx\\). Integrate: \\(y = x^{2}/2 + C\\).",
		},
		{
			q: "Solve: \\(dy/dx = 1/y\\)",
			a: "y^2=2x+C",
			e: "Separate: \\(y dy = dx\\). Integrate: \\(y^{2}/2 = x + C\\). Multiply by 2: \\(y^{2} = 2x + C\\).",
		},
		{
			q: "Solve: \\(dy/dx = x^{2}/y\\)",
			a: "y^2=(2/3)x^3+C",
			e: "Separate: \\(y dy = x^{2} dx\\). Integrate: \\(y^{2}/2 = x^{3}/3 + C\\). Multiply by 2: \\(y^{2} = (2/3)x^{3} + C\\).",
		},
		{
			q: "Solve: \\(dy/dx = x^{2}\\)",
			a: "y=x^3/3+C",
			e: "Rewrite as \\(dy = x^{2} dx\\). Integrate: \\(y = x^{3}/3 + C\\).",
		},
		{
			q: "Solve: \\(dy/dx = 1/y^{2}\\)",
			a: "y^3=3x+C",
			e: "Separate: \\(y^{2} dy = dx\\). Integrate: \\(y^{3}/3 = x + C\\). Multiply by 3: \\(y^{3} = 3x + C\\).",
		},
	}
	t := templates[rand.Intn(len(templates))]
	return generator.Problem{
		Question:    t.q,
		Answer:      t.a,
		Explanation: t.e,
	}
}

type linearFirstGen struct{}

func (g *linearFirstGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type qa struct {
		q string
		a string
		e string
	}
	templates := []qa{
		{
			q: "Solve: \\(y' + 2y = 0\\)",
			a: "Ce^(-2x)",
			e: "This is a homogeneous first-order linear ODE. The integrating factor is \\(e^{2x}\\). \\(\\frac{d}{dx}(ye^{2x}) = 0\\), so \\(ye^{2x} = C\\), giving \\(y = Ce^{-2x}\\).",
		},
		{
			q: "Solve: \\(y' + y = 0\\)",
			a: "Ce^(-x)",
			e: "This is a homogeneous first-order linear ODE. The integrating factor is \\(e^{x}\\). \\(\\frac{d}{dx}(ye^{x}) = 0\\), so \\(ye^{x} = C\\), giving \\(y = Ce^{-x}\\).",
		},
		{
			q: "Solve: \\(y' + 3y = 0\\)",
			a: "Ce^(-3x)",
			e: "This is a homogeneous first-order linear ODE. The integrating factor is \\(e^{3x}\\). \\(\\frac{d}{dx}(ye^{3x}) = 0\\), so \\(ye^{3x} = C\\), giving \\(y = Ce^{-3x}\\).",
		},
		{
			q: "Solve: \\(y' + y = 1\\)",
			a: "1+Ce^(-x)",
			e: "The integrating factor is \\(e^{x}\\). \\(\\frac{d}{dx}(ye^{x}) = e^{x}\\). Integrate: \\(ye^{x} = e^{x} + C\\), so \\(y = 1 + Ce^{-x}\\).",
		},
		{
			q: "Solve: \\(y' + 2y = 4\\)",
			a: "2+Ce^(-2x)",
			e: "The integrating factor is \\(e^{2x}\\). \\(\\frac{d}{dx}(ye^{2x}) = 4e^{2x}\\). Integrate: \\(ye^{2x} = 2e^{2x} + C\\), so \\(y = 2 + Ce^{-2x}\\).",
		},
		{
			q: "Solve: \\(y' + y = e^{x}\\)",
			a: "0.5e^x+Ce^(-x)",
			e: "The integrating factor is \\(e^{x}\\). \\(\\frac{d}{dx}(ye^{x}) = e^{2x}\\). Integrate: \\(ye^{x} = \\frac{1}{2}e^{2x} + C\\), so \\(y = \\frac{1}{2}e^{x} + Ce^{-x}\\).",
		},
	}
	t := templates[rand.Intn(len(templates))]
	return generator.Problem{
		Question:    t.q,
		Answer:      t.a,
		Explanation: t.e,
	}
}

type exactGen struct{}

func (g *exactGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	type entry struct {
		q string
		a string
		e string
	}
	easy := []entry{
		{
			q: "Is \\(M dx + N dy = 0\\) exact when \\(M = y\\) and \\(N = x\\)? (yes/no)",
			a: "yes",
			e: "\\(\\partial M/\\partial y = 1\\) and \\(\\partial N/\\partial x = 1\\). Since \\(\\partial M/\\partial y = \\partial N/\\partial x\\), the equation is exact.",
		},
		{
			q: "Is \\(M dx + N dy = 0\\) exact when \\(M = x\\) and \\(N = x\\)? (yes/no)",
			a: "no",
			e: "\\(\\partial M/\\partial y = 0\\) and \\(\\partial N/\\partial x = 1\\). Since \\(\\partial M/\\partial y \\neq \\partial N/\\partial x\\) \\((0 \\neq 1)\\), the equation is not exact.",
		},
		{
			q: "Is \\(M dx + N dy = 0\\) exact when \\(M = xy\\) and \\(N = xy\\)? (yes/no)",
			a: "no",
			e: "\\(\\partial M/\\partial y = x\\) and \\(\\partial N/\\partial x = y\\). Since \\(\\partial M/\\partial y \\neq \\partial N/\\partial x\\) \\((x \\neq y\\) in general), the equation is not exact.",
		},
	}
	// Parameterized production: dM/dy at a point for M = a*x*y.
	aChoices := []int{2, 3}
	xChoices := []int{2, 3, 4, 5}
	a := aChoices[rand.Intn(len(aChoices))]
	x := xChoices[rand.Intn(len(xChoices))]
	easy = append(easy, entry{
		q: fmt.Sprintf("For \\(M = %dxy\\), what is \\(dM/dy\\) at \\(x = %d\\)? (enter a number)", a, x),
		a: fmt.Sprintf("%d", a*x),
		e: fmt.Sprintf("\\(dM/dy = %dx\\), so at \\(x = %d\\) the value is %d.", a, x, a*x),
	})
	hard := []entry{
		{
			q: "Is \\(M dx + N dy = 0\\) exact when \\(M = 2xy\\) and \\(N = x^{2}\\)? (yes/no)",
			a: "yes",
			e: "\\(\\partial M/\\partial y = 2x\\) and \\(\\partial N/\\partial x = 2x\\). Since \\(\\partial M/\\partial y = \\partial N/\\partial x\\), the equation is exact.",
		},
		{
			q: "Is \\(M dx + N dy = 0\\) exact when \\(M = y^{2}\\) and \\(N = x\\)? (yes/no)",
			a: "no",
			e: "\\(\\partial M/\\partial y = 2y\\) and \\(\\partial N/\\partial x = 1\\). Since \\(\\partial M/\\partial y \\neq \\partial N/\\partial x\\) \\((2y \\neq 1\\) in general), the equation is not exact.",
		},
		{
			q: "Is \\(M dx + N dy = 0\\) exact when \\(M = \\sin(y)\\) and \\(N = x\\cos(y)\\)? (yes/no)",
			a: "yes",
			e: "\\(\\partial M/\\partial y = \\cos(y)\\) and \\(\\partial N/\\partial x = \\cos(y)\\). Since \\(\\partial M/\\partial y = \\partial N/\\partial x\\), the equation is exact.",
		},
	}
	pool := easy
	if scale > 3 {
		pool = append(easy, hard...)
	}
	t := pool[rand.Intn(len(pool))]
	return generator.Problem{
		Question:    t.q,
		Answer:      t.a,
		Explanation: t.e,
	}
}

type homogeneousGen struct{}

func (g *homogeneousGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type qa struct {
		r1, r2 int
	}
	pairs := []qa{
		{1, 2}, {2, 3}, {1, -1}, {2, -2}, {3, -3},
		{1, 3}, {2, 4}, {-1, -2}, {3, 4}, {1, -2},
	}
	p := pairs[rand.Intn(len(pairs))]
	r1, r2 := p.r1, p.r2
	s := r1 + r2
	prod := r1 * r2

	charStr := "r^{2}"
	if s > 0 {
		charStr += fmt.Sprintf(" - %dr", s)
	} else if s < 0 {
		charStr += fmt.Sprintf(" + %dr", -s)
	}
	if prod > 0 {
		charStr += fmt.Sprintf(" + %d", prod)
	} else if prod < 0 {
		charStr += fmt.Sprintf(" - %d", -prod)
	}
	charStr += " = 0"

	eqStr := "y''"
	if s > 0 {
		eqStr += fmt.Sprintf(" - %dy'", s)
	} else if s < 0 {
		eqStr += fmt.Sprintf(" + %dy'", -s)
	}
	if prod > 0 {
		eqStr += fmt.Sprintf(" + %dy", prod)
	} else if prod < 0 {
		eqStr += fmt.Sprintf(" - %dy", -prod)
	}
	eqStr += " = 0"

	answer := fmt.Sprintf("C1*e^(%dx)+C2*e^(%dx)", r2, r1)
	if r1 > r2 {
		answer = fmt.Sprintf("C1*e^(%dx)+C2*e^(%dx)", r1, r2)
	}

	explanation := fmt.Sprintf("Characteristic equation: \\(%s\\). Roots: \\(r = %d\\) and \\(r = %d\\). The general solution is \\(y = C_{1}e^{%dx} + C_{2}e^{%dx}\\).", charStr, r1, r2, r1, r2)

	return generator.Problem{
		Question:    fmt.Sprintf("Solve: \\(%s\\)", eqStr),
		Answer:      answer,
		Explanation: explanation,
	}
}

type nonhomogeneousGen struct{}

func (g *nonhomogeneousGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type qa struct {
		q string
		a string
		e string
	}
	templates := []qa{
		{
			q: "Solve: \\(y'' + y = 1\\)",
			a: "C1*cos(x)+C2*sin(x)+1",
			e: "The homogeneous solution: \\(y_c = C_{1}\\cos(x) + C_{2}\\sin(x)\\) (characteristic \\(r^{2}+1=0\\), roots \\(\\pm i\\)). Particular solution: \\(y_p = 1\\). General: \\(y = C_{1}\\cos(x) + C_{2}\\sin(x) + 1\\).",
		},
		{
			q: "Solve: \\(y'' + y = 2\\)",
			a: "C1*cos(x)+C2*sin(x)+2",
			e: "The homogeneous solution: \\(y_c = C_{1}\\cos(x) + C_{2}\\sin(x)\\). Particular solution: \\(y_p = 2\\). General: \\(y = C_{1}\\cos(x) + C_{2}\\sin(x) + 2\\).",
		},
		{
			q: "Solve: \\(y'' - y = 1\\)",
			a: "C1*e^x+C2*e^(-x)-1",
			e: "The homogeneous solution: \\(y_c = C_{1}e^{x} + C_{2}e^{-x}\\) (characteristic \\(r^{2}-1=0\\), roots \\(\\pm 1\\)). Particular solution: \\(y_p = -1\\). General: \\(y = C_{1}e^{x} + C_{2}e^{-x} - 1\\).",
		},
		{
			q: "Solve: \\(y'' - y = 2\\)",
			a: "C1*e^x+C2*e^(-x)-2",
			e: "The homogeneous solution: \\(y_c = C_{1}e^{x} + C_{2}e^{-x}\\). Particular solution: \\(y_p = -2\\). General: \\(y = C_{1}e^{x} + C_{2}e^{-x} - 2\\).",
		},
		{
			q: "Solve: \\(y'' - 4y = 0\\)",
			a: "C1*e^(2x)+C2*e^(-2x)",
			e: "Characteristic equation: \\(r^{2} - 4 = 0\\), roots \\(r = \\pm 2\\). General solution: \\(y = C_{1}e^{2x} + C_{2}e^{-2x}\\).",
		},
		{
			q: "Solve: \\(y'' - 4y = 1\\)",
			a: "C1*e^(2x)+C2*e^(-2x)-0.25",
			e: "The homogeneous solution: \\(y_c = C_{1}e^{2x} + C_{2}e^{-2x}\\). Particular solution: \\(y_p = -1/4\\). General: \\(y = C_{1}e^{2x} + C_{2}e^{-2x} - 1/4\\).",
		},
	}
	t := templates[rand.Intn(len(templates))]
	return generator.Problem{
		Question:    t.q,
		Answer:      t.a,
		Explanation: t.e,
	}
}

type laplaceGen struct{}

func (g *laplaceGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type qa struct {
		q string
		a string
		e string
	}
	templates := []qa{
		{
			q: "Find \\(L\\{1\\}\\)",
			a: "1/s",
			e: "\\(L\\{1\\} = \\int_{0}^{\\infty} e^{-st} dt = 1/s\\) for \\(s > 0\\).",
		},
		{
			q: "Find \\(L\\{t\\}\\)",
			a: "1/s^2",
			e: "\\(L\\{t\\} = 1/s^{2}\\) for \\(s > 0\\).",
		},
		{
			q: "Find \\(L\\{t^{2}\\}\\)",
			a: "2/s^3",
			e: "\\(L\\{t^{2}\\} = 2/s^{3}\\) for \\(s > 0\\).",
		},
	}
	alphas := []int{1, 2, 3}
	a := alphas[rand.Intn(len(alphas))]
	templates = append(templates,
		qa{
			q: fmt.Sprintf("Find \\(L\\{e^{%dt}\\}\\)", a),
			a: fmt.Sprintf("1/(s-%d)", a),
			e: fmt.Sprintf("\\(L\\{e^{%dt}\\} = 1/(s - %d)\\) for \\(s > %d\\).", a, a, a),
		},
		qa{
			q: fmt.Sprintf("Find \\(L\\{\\sin(%dt)\\}\\)", a),
			a: fmt.Sprintf("%d/(s^2+%d)", a, a*a),
			e: fmt.Sprintf("\\(L\\{\\sin(%dt)\\} = %d/(s^{2} + %d^{2})\\) for \\(s > 0\\).", a, a, a),
		},
		qa{
			q: fmt.Sprintf("Find \\(L\\{\\cos(%dt)\\}\\)", a),
			a: fmt.Sprintf("s/(s^2+%d)", a*a),
			e: fmt.Sprintf("\\(L\\{\\cos(%dt)\\} = s/(s^{2} + %d^{2})\\) for \\(s > 0\\).", a, a),
		},
	)

	t := templates[rand.Intn(len(templates))]
	return generator.Problem{
		Question:    t.q,
		Answer:      t.a,
		Explanation: t.e,
	}
}

type systemsGen struct{}

func (g *systemsGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type qa struct {
		q string
		a string
		e string
	}
	templates := []qa{
		{
			q: "A system of 2 first-order ODEs represents how many unknown functions? (enter a number)",
			a: "2",
			e: "Each first-order ODE typically describes the derivative of one unknown function, so 2 equations means 2 unknown functions.",
		},
		{
			q: "A system of 3 first-order ODEs represents how many unknown functions? (enter a number)",
			a: "3",
			e: "Each first-order ODE typically describes the derivative of one unknown function, so 3 equations means 3 unknown functions.",
		},
		{
			q: "\\(dx/dt = 2x\\), \\(dy/dt = 3y\\). Are \\(x\\) and \\(y\\) independent? (yes/no)",
			a: "yes",
			e: "The equations are uncoupled: \\(dx/dt\\) depends only on \\(x\\), and \\(dy/dt\\) depends only on \\(y\\). Each can be solved independently.",
		},
		{
			q: "\\(dx/dt = x + y\\), \\(dy/dt = x - y\\). Is this a coupled system? (yes/no)",
			a: "yes",
			e: "Yes, because \\(dx/dt\\) depends on \\(y\\) and \\(dy/dt\\) depends on \\(x\\). The equations must be solved together.",
		},
		{
			q: "Can a system of two first-order ODEs be converted into a single second-order ODE? (yes/no)",
			a: "yes",
			e: "Yes, by differentiating one equation and substituting the other, you can eliminate one variable to obtain a single higher-order ODE.",
		},
		{
			q: "Does \\(dx/dt = 2x\\) represent a system of differential equations? (yes/no)",
			a: "no",
			e: "This is a single differential equation, not a system. A system requires at least two equations.",
		},
	}
	t := templates[rand.Intn(len(templates))]
	return generator.Problem{
		Question:    t.q,
		Answer:      t.a,
		Explanation: t.e,
	}
}

type variationParamsGen struct{}

func (g *variationParamsGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type qa struct {
		q string
		a string
		e string
	}
	templates := []qa{
		{q: "For \\(y'' + y = \\sec(x)\\), the homogeneous solutions are \\(y_1=\\cos x\\), \\(y_2=\\sin x\\). What is the Wronskian \\(W\\)? (enter a number)", a: "1", e: "\\(W = y_1y_2' - y_1'y_2 = \\cos x \\cdot \\cos x - (-\\sin x)\\sin x = 1\\)."},
		{q: "In variation of parameters, \\(y_p = u_1 y_1 + u_2 y_2\\) where \\(u_1' = -y_2 g/W\\). For \\(y''+y=1\\) with \\(W=1\\), \\(y_2=\\sin x\\), what is \\(u_1'\\)?", a: "-sin(x)", e: "\\(u_1' = -y_2 g/W = -\\sin x \\cdot 1/1 = -\\sin x\\)."},
		{q: "Does variation of parameters require knowing the homogeneous solution? (yes/no)", a: "yes", e: "Yes, the method builds the particular solution from the homogeneous basis \\(y_1, y_2\\) and the Wronskian."},
		{q: "For constant-coefficient ODEs, variation of parameters and undetermined coefficients give the same particular solution when both apply. Is this true? (yes/no)", a: "yes", e: "Both methods recover a particular solution; variation of parameters is more general and works even when undetermined coefficients does not."},
	}
	t := templates[rand.Intn(len(templates))]
	return generator.Problem{Question: t.q, Answer: t.a, Explanation: t.e}
}

type convolutionGen struct{}

func (g *convolutionGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type qa struct {
		q string
		a string
		e string
	}
	templates := []qa{
		{q: "The convolution \\((f*g)(t) = \\int_0^t f(\\tau)g(t-\\tau)d\\tau\\). What is \\((1*1)(t)\\)?", a: "t", e: "\\((1*1)(t)=\\int_0^t 1\\cdot1 d\\tau = t\\)."},
		{q: "\\(L\\{f*g\\} = L\\{f\\} \\cdot L\\{g\\}\\). If \\(L\\{f\\}=1/s\\) and \\(L\\{g\\}=1/s\\), what is \\(L\\{f*g\\}\\)?", a: "1/s^2", e: "By the convolution theorem, \\(L\\{f*g\\}= (1/s)(1/s)=1/s^2\\)."},
		{q: "Is convolution commutative: \\(f*g = g*f\\)? (yes/no)", a: "yes", e: "Yes, by substitution \\(\\tau \\to t-\\tau\\) the integral is symmetric."},
		{q: "If \\(f(t)=e^{t}\\) and \\(g(t)=1\\), \\((f*g)(t)=e^{t}-1\\). Does \\(L\\{f*g\\}=1/(s(s-1))\\) hold? (yes/no)", a: "yes", e: "\\(L\\{e^t\\}=1/(s-1)\\), \\(L\\{1\\}=1/s\\), product is \\(1/(s(s-1))=L\\{e^t-1\\}\\)."},
	}
	t := templates[rand.Intn(len(templates))]
	return generator.Problem{Question: t.q, Answer: t.a, Explanation: t.e}
}

type linearPhaseGen struct{}

func (g *linearPhaseGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	type qa struct {
		q string
		a string
		e string
	}
	templates := []qa{
		{q: "For \\(x' = Ax\\) with \\(A=[[2,0],[0,-3]]\\), the eigenvalues are 2 and -3. What type is the origin? (saddle/node/spiral)", a: "saddle", e: "One positive, one negative eigenvalue gives a saddle."},
		{q: "For \\(x' = Ax\\) with eigenvalues \\(-1, -2\\), what type is the origin? (saddle/node/spiral)", a: "node", e: "Both eigenvalues negative and real gives a stable node."},
		{q: "If eigenvalues are \\(0.5 \\pm 2i\\), what type is the origin? (saddle/node/spiral)", a: "spiral", e: "Complex pair with positive real part gives an unstable spiral."},
		{q: "Does the phase portrait of a linear system depend on eigenvalues of \\(A\\)? (yes/no)", a: "yes", e: "The eigenvalues classify the equilibrium: node, saddle, spiral, or centre."},
	}
	t := templates[rand.Intn(len(templates))]
	return generator.Problem{Question: t.q, Answer: t.a, Explanation: t.e}
}

type eulerGen struct{}

func (g *eulerGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*3)
	hChoices := []float64{0.1, 0.2, 0.5}
	h := hChoices[rand.Intn(len(hChoices))]
	x0 := 0.0
	y0 := float64(rand.Intn(max(1, scale*2)) + 1)
	// simple ODE y' = y, exact y = y0*e^x, Euler: y1 = y0 + h*y0
	y1 := y0 + h*y0
	_ = x0
	return generator.Problem{
		Question:    fmt.Sprintf("Use one step of Euler's method with \\(h=%.1f\\) for \\(y'=y\\), \\(y(0)=%.0f\\). Estimate \\(y(%.1f)\\).", h, y0, h),
		Answer:      fmt.Sprintf("%.4f", y1),
		Explanation: fmt.Sprintf("Euler: \\(y_1 = y_0 + h f(x_0,y_0) = %.0f + %.1f\\cdot%.0f = %.4f\\).", y0, h, y0, y1),
	}
}

type rk4Gen struct{}

func (g *rk4Gen) Generate(ctx generator.GeneratorContext) generator.Problem {
	mode := rand.Intn(5)
	switch mode {
	case 0:
		// Stage evaluations across N steps
		n := smallInt(2, 6+int(ctx.Difficulty*10))
		ans := 4 * n
		return generator.Problem{
			Question: fmt.Sprintf("How many stage function evaluations does the classic 4th-order Runge-Kutta method (RK4) require across %d steps?", n),
			Answer:   fmt.Sprintf("%d", ans),
			Explanation: fmt.Sprintf(
				"Each step of classic RK4 evaluates the slope at 4 stages (k_1, k_2, k_3 and k_4). For %d steps, the total number of evaluations is 4 * %d = %d.",
				n, n, ans,
			),
		}

	case 1:
		// Global error reduction factor under step size division
		kChoices := []int{2, 3, 4, 5}
		k := kChoices[rand.Intn(len(kChoices))]
		ans := k * k * k * k
		return generator.Problem{
			Question: fmt.Sprintf(
				"If the step size \\(h\\) in the fourth-order Runge-Kutta method (RK4) is reduced by a factor of %d, by what theoretical factor is the global truncation error reduced?",
				k,
			),
			Answer: fmt.Sprintf("%d", ans),
			Explanation: fmt.Sprintf(
				"RK4 has a global truncation error of \\(\\mathcal{O}(h^4)\\). Dividing the step size by %d reduces the global error by a factor of \\(%d^4 = %d\\).",
				k, k, ans,
			),
		}

	case 2:
		// First stage slope k_1 = f(x_0, y_0)
		a := smallInt(1, 4)
		b := smallInt(1, 4)
		x0 := smallInt(0, 3)
		y0 := smallInt(1, 4)
		k1 := a*x0 + b*y0
		return generator.Problem{
			Question: fmt.Sprintf(
				"For the initial value problem \\(y' = %dx + %dy\\) with \\(y(%d) = %d\\), compute the first stage slope \\(k_1\\) in the classic RK4 method.",
				a, b, x0, y0,
			),
			Answer: fmt.Sprintf("%d", k1),
			Explanation: fmt.Sprintf(
				"In the classic RK4 method, the first stage slope is evaluated at the initial point: \\(k_1 = f(x_0, y_0) = %d(%d) + %d(%d) = %d\\).",
				a, x0, b, y0, k1,
			),
		}

	case 3:
		// Truncation error order
		if rand.Float64() < 0.5 {
			return generator.Problem{
				Question: "What is the integer order \\(p\\) of the global truncation error \\(\\mathcal{O}(h^p)\\) for the classic Runge-Kutta method?",
				Answer:   "4",
				Explanation: "The classic fourth-order Runge-Kutta method achieves a global truncation error of order \\(\\mathcal{O}(h^4)\\).",
			}
		}
		return generator.Problem{
			Question: "What is the integer order \\(p\\) of the local truncation error \\(\\mathcal{O}(h^p)\\) for the classic fourth-order Runge-Kutta method?",
			Answer:   "5",
			Explanation: "Classic RK4 has local truncation error \\(\\mathcal{O}(h^5)\\), which accumulates over \\(\\mathcal{O}(1/h)\\) steps to produce a global error of \\(\\mathcal{O}(h^4)\\).",
		}

	default:
		// Sum of intermediate stage weights in butcher increment
		return generator.Problem{
			Question: "In the classic RK4 increment formula \\(\\frac{h}{6}(k_1 + w_2 k_2 + w_3 k_3 + k_4)\\), what is the sum of the inner weights \\(w_2 + w_3\\)?",
			Answer:   "4",
			Explanation: "The classic RK4 butcher tableau weights the mid-point slopes at 2 each (\\(w_2 = 2\\) and \\(w_3 = 2\\)), so \\(w_2 + w_3 = 2 + 2 = 4\\).",
		}
	}
}

type stabilityGen struct{}

func (g *stabilityGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	mode := rand.Intn(5)
	switch mode {
	case 0:
		// Real equilibria count for x' = mu - x^2
		mus := []int{-3, -2, -1, 1, 2, 4}
		mu := mus[rand.Intn(len(mus))]
		count := 0
		reason := fmt.Sprintf("Because %d < 0, no real solution exists and there are 0 real equilibria.", mu)
		if mu > 0 {
			count = 2
			reason = fmt.Sprintf("Because %d > 0, there are 2 real equilibria at \\(\\pm\\sqrt{%d}\\).", mu, mu)
		}
		return generator.Problem{
			Question: fmt.Sprintf("For the autonomous differential equation \\(x' = %d - x^2\\), how many real equilibrium points exist?", mu),
			Answer:   fmt.Sprintf("%d", count),
			Explanation: fmt.Sprintf(
				"Setting \\(x' = 0\\) gives \\(x^2 = %d\\). %s",
				mu, reason,
			),
		}

	case 1:
		// Value of non-zero equilibrium point
		a := smallInt(2, 6+int(ctx.Difficulty*4))
		return generator.Problem{
			Question: fmt.Sprintf("Find the non-zero equilibrium point \\(x^*\\) of the flow equation \\(x' = %dx - x^2\\).", a),
			Answer:   fmt.Sprintf("%d", a),
			Explanation: fmt.Sprintf(
				"Setting \\(x' = x(%d - x) = 0\\) yields equilibria at \\(x = 0\\) and \\(x = %d\\). The non-zero equilibrium is %d.",
				a, a, a,
			),
		}

	case 2:
		// Linearized derivative at stable equilibrium
		a := smallInt(2, 6+int(ctx.Difficulty*4))
		return generator.Problem{
			Question: fmt.Sprintf("For the flow \\(x' = %dx - x^2\\), evaluate the derivative \\(f'(x^*)\\) at the stable equilibrium \\(x^* = %d\\).", a, a),
			Answer:   fmt.Sprintf("%d", -a),
			Explanation: fmt.Sprintf(
				"The flow function is \\(f(x) = %dx - x^2\\) with derivative \\(f'(x) = %d - 2x\\). At \\(x^* = %d\\), \\(f'(%d) = %d - 2(%d) = %d\\). The negative derivative confirms asymptotic stability.",
				a, a, a, a, a, a, -a,
			),
		}

	case 3:
		// Number of asymptotically stable equilibria
		a := smallInt(2, 6+int(ctx.Difficulty*4))
		return generator.Problem{
			Question: fmt.Sprintf("How many asymptotically stable equilibrium points does the differential equation \\(x' = %dx - x^2\\) possess?", a),
			Answer:   "1",
			Explanation: fmt.Sprintf(
				"The equilibria are \\(x = 0\\) and \\(x = %d\\). The derivative \\(f'(x) = %d - 2x\\) gives \\(f'(0) = %d > 0\\) (unstable) and \\(f'(%d) = -%d < 0\\) (stable), so exactly 1 equilibrium is asymptotically stable.",
				a, a, a, a, a,
			),
		}

	default:
		// 2D linear system trace
		p := smallInt(1, 5)
		q := smallInt(1, 5)
		tr := -(p + q)
		return generator.Problem{
			Question: fmt.Sprintf("For the 2D linear system \\(x' = -%dx\\) and \\(y' = -%dy\\) governing decay to the origin, what is the trace of the system matrix \\(A\\)?", p, q),
			Answer:   fmt.Sprintf("%d", tr),
			Explanation: fmt.Sprintf(
				"The system matrix is diagonal with eigenvalues \\(\\lambda_1 = -%d\\) and \\(\\lambda_2 = -%d\\). The trace is the sum of the eigenvalues: \\(\\operatorname{tr}(A) = -%d + (-%d) = %d\\).",
				p, q, p, q, tr,
			),
		}
	}
}

type characteristicsGen struct{}

func (g *characteristicsGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	mode := rand.Intn(5)
	switch mode {
	case 0:
		// Propagation speed dx/dt for transport equation u_t + c u_x = 0
		c := smallInt(2, 6+int(ctx.Difficulty*4))
		return generator.Problem{
			Question: fmt.Sprintf("For the first-order linear PDE \\(u_t + %du_x = 0\\), what is the characteristic propagation speed \\(\\frac{dx}{dt}\\)?", c),
			Answer:   fmt.Sprintf("%d", c),
			Explanation: fmt.Sprintf(
				"The characteristic equations for the advection equation \\(u_t + %du_x = 0\\) are \\(\\frac{dt}{1} = \\frac{dx}{%d} = \\frac{du}{0}\\). Along the characteristic curves, \\(\\frac{dx}{dt} = %d\\).",
				c, c, c,
			),
		}

	case 1:
		// Evaluation of quadratic initial condition u(x, 0) = x^2
		c := smallInt(2, 4)
		t0 := smallInt(1, 3)
		delta := smallInt(1, 5)
		x0 := c*t0 + delta
		ans := delta * delta
		return generator.Problem{
			Question: fmt.Sprintf("Consider the transport equation \\(u_t + %du_x = 0\\) with initial profile \\(u(x, 0) = x^2\\). Evaluate \\(u(%d, %d)\\).", c, x0, t0),
			Answer:   fmt.Sprintf("%d", ans),
			Explanation: fmt.Sprintf(
				"Along the characteristic curves \\(x - %dt = C\\), the solution is constant: \\(u(x, t) = u(x - %dt, 0) = (x - %dt)^2\\). Substituting \\((x, t) = (%d, %d)\\) gives \\(u(%d, %d) = (%d - %d \\cdot %d)^2 = %d^2 = %d\\).",
				c, c, c, x0, t0, x0, t0, x0, c, t0, delta, ans,
			),
		}

	case 2:
		// Linear initial condition u(x, 0) = m*x + k
		c := smallInt(2, 4)
		m := smallInt(2, 5)
		k := smallInt(1, 9)
		t0 := smallInt(1, 3)
		x0 := smallInt(2, 8) + c*t0
		xBase := x0 - c*t0
		ans := m*xBase + k
		return generator.Problem{
			Question: fmt.Sprintf("For the advection equation \\(u_t + %du_x = 0\\) with initial condition \\(u(x, 0) = %dx + %d\\), evaluate \\(u(%d, %d)\\).", c, m, k, x0, t0),
			Answer:   fmt.Sprintf("%d", ans),
			Explanation: fmt.Sprintf(
				"The general solution is \\(u(x, t) = u(x - %dt, 0) = %d(x - %dt) + %d\\). Evaluating at \\((x, t) = (%d, %d)\\) yields \\(u(%d, %d) = %d(%d - %d \\cdot %d) + %d = %d(%d) + %d = %d\\).",
				c, m, c, k, x0, t0, x0, t0, m, x0, c, t0, k, m, xBase, k, ans,
			),
		}

	case 3:
		// Characteristic ray origin at t = 0
		c := smallInt(2, 5)
		t1 := smallInt(2, 5)
		xOrigin := smallInt(1, 10)
		x1 := xOrigin + c*t1
		return generator.Problem{
			Question: fmt.Sprintf("A characteristic curve for the PDE \\(u_t + %du_x = 0\\) passes through the point \\((x, t) = (%d, %d)\\). At what coordinate \\(x\\) does this characteristic curve intersect the initial axis \\(t = 0\\)?", c, x1, t1),
			Answer:   fmt.Sprintf("%d", xOrigin),
			Explanation: fmt.Sprintf(
				"The characteristic curves satisfy \\(\\frac{dx}{dt} = %d\\), which integrates to \\(x - %dt = x_0\\). Using the point \\((%d, %d)\\), the initial position is \\(x_0 = %d - %d \\cdot %d = %d\\).",
				c, c, x1, t1, x1, c, t1, xOrigin,
			),
		}

	default:
		// Inhomogeneous advection with zero initial state u(x, 0) = 0
		c := smallInt(2, 6)
		k := smallInt(2, 7)
		t0 := smallInt(2, 6)
		x0 := smallInt(1, 12)
		ans := k * t0
		return generator.Problem{
			Question: fmt.Sprintf("Consider the inhomogeneous transport equation \\(u_t + %du_x = %d\\) with initial condition \\(u(x, 0) = 0\\). What is the value of \\(u(%d, %d)\\)?", c, k, x0, t0),
			Answer:   fmt.Sprintf("%d", ans),
			Explanation: fmt.Sprintf(
				"Along the characteristics \\(\\frac{dx}{dt} = %d\\), the transport equation simplifies to \\(\\frac{du}{dt} = %d\\). Integrating from \\(t = 0\\) with \\(u(x - %dt, 0) = 0\\) yields \\(u(x, t) = %dt\\). Evaluating at \\(t = %d\\) gives \\(u(%d, %d) = %d \\cdot %d = %d\\).",
				c, k, c, k, t0, x0, t0, k, t0, ans,
			),
		}
	}
}

type fourierGen struct{}

func (g *fourierGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	mode := rand.Intn(5)
	switch mode {
	case 0:
		// Algebraic multiplier evaluation at specific angular frequency omega
		k := smallInt(12, 30)
		w := smallInt(1, 4)
		pw := k - w*w
		return generator.Problem{
			Question: fmt.Sprintf("Applying the Fourier transform to \\(y'' + %dy = f(t)\\) produces the algebraic relation \\(P(\\omega) \\hat{y}(\\omega) = \\hat{f}(\\omega)\\), where \\(P(\\omega) = %d - \\omega^2\\). Evaluate \\(P(%d)\\).", k, k, w),
			Answer:   fmt.Sprintf("%d", pw),
			Explanation: fmt.Sprintf(
				"The Fourier transform replaces \\(\\frac{d^2}{dt^2}\\) with \\((i\\omega)^2 = -\\omega^2\\). Therefore \\(P(\\omega) = %d - \\omega^2\\). At \\(\\omega = %d\\), \\(P(%d) = %d - %d^2 = %d - %d = %d\\).",
				k, w, w, k, w, k, w*w, pw,
			),
		}

	case 1:
		// Squared transfer denominator for first-order system y' + a*y = f(t)
		a := smallInt(2, 6)
		w := smallInt(1, 8)
		denom := w*w + a*a
		return generator.Problem{
			Question: fmt.Sprintf("For the relaxation system \\(y' + %dy = f(t)\\), the squared denominator of the frequency response function is \\(D(\\omega) = \\omega^2 + %d^2\\). Evaluate \\(D(%d)\\).", a, a, w),
			Answer:   fmt.Sprintf("%d", denom),
			Explanation: fmt.Sprintf(
				"The Fourier transform of the operator is \\(i\\omega + %d\\). The squared denominator is \\(|i\\omega + %d|^2 = \\omega^2 + %d^2 = \\omega^2 + %d\\). Evaluating at \\(\\omega = %d\\) gives \\(%d^2 + %d = %d + %d = %d\\).",
				a, a, a, a*a, w, w, a*a, w*w, a*a, denom,
			),
		}

	case 2:
		// Resonant frequency of undamped oscillator
		w0 := smallInt(2, 9)
		k := w0 * w0
		return generator.Problem{
			Question: fmt.Sprintf("At what positive angular frequency \\(\\omega > 0\\) does the undamped oscillator \\(y'' + %dy = f(t)\\) experience resonance?", k),
			Answer:   fmt.Sprintf("%d", w0),
			Explanation: fmt.Sprintf(
				"The Fourier transform gives \\((-\\omega^2 + %d)\\hat{y}(\\omega) = \\hat{f}(\\omega)\\). Resonance occurs when the transfer denominator vanishes: \\(\\omega^2 = %d\\), giving \\(\\omega = %d\\).",
				k, k, w0,
			),
		}

	case 3:
		// Fourth-order operator frequency factor
		c := smallInt(1, 20)
		w := smallInt(1, 3)
		w4 := w * w * w * w
		ans := w4 + c
		return generator.Problem{
			Question: fmt.Sprintf("For the fourth-order differential equation \\(y^{(4)} + %dy = f(t)\\), the Fourier transform yields \\((\\omega^4 + %d)\\hat{y}(\\omega) = \\hat{f}(\\omega)\\). Evaluate \\(\\omega^4 + %d\\) at \\(\\omega = %d\\).", c, c, c, w),
			Answer:   fmt.Sprintf("%d", ans),
			Explanation: fmt.Sprintf(
				"Because \\((i\\omega)^4 = i^4 \\omega^4 = \\omega^4\\), the operator multiplier is \\(\\omega^4 + %d\\). Substituting \\(\\omega = %d\\) gives \\(%d^4 + %d = %d + %d = %d\\).",
				c, w, w, c, w4, c, ans,
			),
		}

	default:
		// Plancherel total energy
		energy := smallInt(5, 50)
		return generator.Problem{
			Question: fmt.Sprintf("A signal \\(y(t)\\) has total time-domain energy \\(\\int_{-\\infty}^\\infty |y(t)|^2 dt = %d\\). What is the total spectral energy \\(\\frac{1}{2\\pi}\\int_{-\\infty}^\\infty |\\hat{y}(\\omega)|^2 d\\omega\\)?", energy),
			Answer:   fmt.Sprintf("%d", energy),
			Explanation: fmt.Sprintf(
				"By Plancherel's theorem for the Fourier transform, the total energy in the time domain is identical to the total energy in the frequency domain: \\(\\frac{1}{2\\pi}\\int_{-\\infty}^\\infty |\\hat{y}(\\omega)|^2 d\\omega = \\int_{-\\infty}^\\infty |y(t)|^2 dt = %d\\).",
				energy,
			),
		}
	}
}

type bifurcationGen struct{}

func (g *bifurcationGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	mode := rand.Intn(5)
	switch mode {
	case 0:
		// Real equilibria count for pitchfork bifurcation x' = mu*x - x^3
		mus := []int{-4, -3, -2, -1, 1, 2, 3, 5}
		mu := mus[rand.Intn(len(mus))]
		count := 1
		reason := fmt.Sprintf("Because \\(\\mu = %d < 0\\), the equation \\(x(\\mu - x^2) = 0\\) has only the origin \\(x = 0\\) as a real solution (1 real equilibrium).", mu)
		if mu > 0 {
			count = 3
			reason = fmt.Sprintf("Because \\(\\mu = %d > 0\\), the equation \\(x(\\mu - x^2) = 0\\) yields three real equilibria at \\(x = 0\\) and \\(x = \\pm\\sqrt{%d}\\).", mu, mu)
		}
		return generator.Problem{
			Question: fmt.Sprintf("For the pitchfork bifurcation system \\(x' = %dx - x^3\\), how many real equilibrium points exist?", mu),
			Answer:   fmt.Sprintf("%d", count),
			Explanation: fmt.Sprintf(
				"Setting \\(x' = x(%d - x^2) = 0\\) gives equilibria. %s",
				mu, reason,
			),
		}

	case 1:
		// Non-zero equilibrium point in transcritical bifurcation x' = ax - x^2
		a := smallInt(2, 6+int(ctx.Difficulty*4))
		return generator.Problem{
			Question: fmt.Sprintf("In the transcritical bifurcation system \\(x' = %dx - x^2\\), find the value of the non-zero equilibrium point \\(x^*\\).", a),
			Answer:   fmt.Sprintf("%d", a),
			Explanation: fmt.Sprintf(
				"Factoring \\(x' = x(%d - x) = 0\\) gives the trivial equilibrium \\(x = 0\\) and the non-zero branch \\(x^* = %d\\).",
				a, a,
			),
		}

	case 2:
		// Linearized derivative at the non-zero equilibrium of transcritical flow
		a := smallInt(2, 6+int(ctx.Difficulty*4))
		return generator.Problem{
			Question: fmt.Sprintf("For the transcritical flow \\(x' = %dx - x^2\\), evaluate the linearized derivative \\(f'(x^*)\\) at the non-zero equilibrium \\(x^* = %d\\).", a, a),
			Answer:   fmt.Sprintf("%d", -a),
			Explanation: fmt.Sprintf(
				"The flow function \\(f(x) = %dx - x^2\\) has derivative \\(f'(x) = %d - 2x\\). At \\(x^* = %d\\), \\(f'(%d) = %d - 2(%d) = %d\\). The negative derivative indicates exchange of stability.",
				a, a, a, a, a, a, -a,
			),
		}

	case 3:
		// Critical saddle-node parameter r_c for x' = r + 2kx - x^2
		k := smallInt(1, 5)
		twoK := 2 * k
		rc := -(k * k)
		return generator.Problem{
			Question: fmt.Sprintf("At what critical parameter value \\(r_c\\) does a saddle-node bifurcation occur for \\(x' = r + %dx - x^2\\)?", twoK),
			Answer:   fmt.Sprintf("%d", rc),
			Explanation: fmt.Sprintf(
				"Rewriting \\(f(x) = r + %dx - x^2 = r + %d - (x - %d)^2\\), equilibria satisfy \\((x - %d)^2 = r + %d\\). Real solutions appear when the discriminant \\(r + %d = 0\\), yielding critical parameter \\(r_c = %d\\).",
				twoK, k*k, k, k, k*k, k*k, rc,
			),
		}

	default:
		// Squared amplitude of non-trivial equilibria in pitchfork bifurcation
		k := smallInt(2, 9)
		return generator.Problem{
			Question: fmt.Sprintf("In the supercritical pitchfork bifurcation \\(x' = %dx - x^3\\), what is the squared magnitude \\((x^*)^2\\) of the non-trivial equilibrium points?", k),
			Answer:   fmt.Sprintf("%d", k),
			Explanation: fmt.Sprintf(
				"Setting \\(x' = x(%d - x^2) = 0\\) for non-zero equilibria gives \\(x^2 = %d\\). Therefore \\((x^*)^2 = %d\\).",
				k, k, k,
			),
		}
	}
}

type waveGen struct{}

func (g *waveGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	mode := rand.Intn(5)
	switch mode {
	case 0:
		// Propagation speed c from u_tt = c^2 u_xx
		c := smallInt(2, 6+int(ctx.Difficulty*4))
		c2 := c * c
		return generator.Problem{
			Question: fmt.Sprintf("For the one-dimensional wave equation \\(u_{tt} = %du_{xx}\\), what is the propagation speed \\(c > 0\\) of traveling wave disturbances?", c2),
			Answer:   fmt.Sprintf("%d", c),
			Explanation: fmt.Sprintf(
				"The wave equation is written as \\(u_{tt} = c^2 u_{xx}\\). Setting \\(c^2 = %d\\) gives the positive wave speed \\(c = \\sqrt{%d} = %d\\).",
				c2, c2, c,
			),
		}

	case 1:
		// d'Alembert with quadratic initial displacement u(x, 0) = x^2 and u_t(x, 0) = 0
		c := smallInt(1, 3)
		x0 := smallInt(1, 5)
		t0 := smallInt(1, 4)
		ans := x0*x0 + c*c*t0*t0
		return generator.Problem{
			Question: fmt.Sprintf("Consider the wave equation \\(u_{tt} = %du_{xx}\\) with initial conditions \\(u(x, 0) = x^2\\) and \\(u_t(x, 0) = 0\\). Evaluate \\(u(%d, %d)\\).", c*c, x0, t0),
			Answer:   fmt.Sprintf("%d", ans),
			Explanation: fmt.Sprintf(
				"By d'Alembert's formula with zero initial velocity, \\(u(x, t) = \\frac{1}{2}[(x - %dt)^2 + (x + %dt)^2] = x^2 + %dt^2\\). Substituting \\((x, t) = (%d, %d)\\) gives \\(u(%d, %d) = %d^2 + %d(%d^2) = %d + %d = %d\\).",
				c, c, c*c, x0, t0, x0, t0, x0, c*c, t0, x0*x0, c*c*t0*t0, ans,
			),
		}

	case 2:
		// Constant initial velocity u_t(x, 0) = v0 and u(x, 0) = 0
		c := smallInt(1, 4)
		v0 := smallInt(2, 7)
		x0 := smallInt(1, 8)
		t0 := smallInt(1, 5)
		ans := v0 * t0
		return generator.Problem{
			Question: fmt.Sprintf("Consider the wave equation \\(u_{tt} = %du_{xx}\\) with initial conditions \\(u(x, 0) = 0\\) and constant initial velocity \\(u_t(x, 0) = %d\\). Evaluate \\(u(%d, %d)\\).", c*c, v0, x0, t0),
			Answer:   fmt.Sprintf("%d", ans),
			Explanation: fmt.Sprintf(
				"Applying d'Alembert's velocity integral: \\(u(x, t) = \\frac{1}{2(%d)}\\int_{x - %dt}^{x + %dt} %d\\,ds = \\frac{%d}{2(%d)}(2 \\cdot %dt) = %dt\\). At \\(t = %d\\), \\(u(%d, %d) = %d \\cdot %d = %d\\).",
				c, c, c, v0, v0, c, c, v0, t0, x0, t0, v0, t0, ans,
			),
		}

	case 3:
		// Length of domain of dependence interval [x - ct, x + ct]
		c := smallInt(2, 5)
		t0 := smallInt(2, 5)
		length := 2 * c * t0
		return generator.Problem{
			Question: fmt.Sprintf("For the wave equation \\(u_{tt} = %du_{xx}\\), what is the length of the domain of dependence interval \\([x - ct, x + ct]\\) on the initial axis at time \\(t = %d\\)?", c*c, t0),
			Answer:   fmt.Sprintf("%d", length),
			Explanation: fmt.Sprintf(
				"The domain of dependence on the line \\(t = 0\\) for a point \\((x, t)\\) is the interval \\([x - ct, x + ct]\\). With wave velocity \\(c = %d\\) and time \\(t = %d\\), the length of this interval is \\(2ct = 2(%d)(%d) = %d\\).",
				c, t0, c, t0, length,
			),
		}

	default:
		// Linear initial displacement u(x, 0) = mx + k and u_t(x, 0) = 0
		c := smallInt(1, 4)
		m := smallInt(2, 5)
		k := smallInt(1, 9)
		x0 := smallInt(1, 8)
		t0 := smallInt(1, 4)
		ans := m*x0 + k
		return generator.Problem{
			Question: fmt.Sprintf("Consider the wave equation \\(u_{tt} = %du_{xx}\\) with initial conditions \\(u(x, 0) = %dx + %d\\) and \\(u_t(x, 0) = 0\\). Evaluate \\(u(%d, %d)\\).", c*c, m, k, x0, t0),
			Answer:   fmt.Sprintf("%d", ans),
			Explanation: fmt.Sprintf(
				"Using d'Alembert's formula with zero initial velocity gives \\(u(x, t) = \\frac{1}{2}[(%d(x - %dt) + %d) + (%d(x + %dt) + %d)] = %dx + %d\\). At \\((x, t) = (%d, %d)\\), \\(u(%d, %d) = %d(%d) + %d = %d\\).",
				m, c, k, m, c, k, m, k, x0, t0, x0, t0, m, x0, k, ans,
			),
		}
	}
}

type heatGen struct{}

func (g *heatGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	mode := rand.Intn(5)
	switch mode {
	case 0:
		// Spatial Dirichlet eigenvalue lambda_n on [0, pi]
		n := smallInt(1, 5+int(ctx.Difficulty*4))
		ans := n * n
		return generator.Problem{
			Question: fmt.Sprintf("For the heat equation on the spatial domain \\([0, \\pi]\\) with Dirichlet boundary conditions \\(u(0, t) = u(\\pi, t) = 0\\), what is the spatial eigenvalue \\(\\lambda_%d\\) corresponding to eigenfunction \\(\\sin(%dx)\\)?", n, n),
			Answer:   fmt.Sprintf("%d", ans),
			Explanation: fmt.Sprintf(
				"The spatial eigenvalue problem is \\(X'' + \\lambda X = 0\\) with \\(X(0) = X(\\pi) = 0\\). The eigenfunctions are \\(X_n(x) = \\sin(nx)\\) with eigenvalues \\(\\lambda_n = n^2\\). For \\(n = %d\\), \\(\\lambda_%d = %d^2 = %d\\).",
				n, n, n, ans,
			),
		}

	case 1:
		// Temporal decay rate beta = alpha * lambda_n
		alpha := smallInt(1, 5)
		n := smallInt(1, 4)
		lambda := n * n
		beta := alpha * lambda
		return generator.Problem{
			Question: fmt.Sprintf("For the diffusion equation \\(u_t = %du_{xx}\\) on \\([0, \\pi]\\) with Dirichlet ends, what is the exponential decay rate \\(\\beta\\) in the temporal factor \\(e^{-\\beta t}\\) for harmonic mode \\(n = %d\\)?", alpha, n),
			Answer:   fmt.Sprintf("%d", beta),
			Explanation: fmt.Sprintf(
				"Separation of variables gives \\(T'(t) = -\\alpha \\lambda_n T(t)\\). On \\([0, \\pi]\\), the \\(n = %d\\) spatial eigenvalue is \\(\\lambda_%d = %d^2 = %d\\). The decay rate is \\(\\beta = %d \\cdot %d = %d\\).",
				n, n, n, lambda, alpha, lambda, beta,
			),
		}

	case 2:
		// Steady-state linear temperature profile
		L := smallInt(2, 8)
		T0 := smallInt(0, 30)
		s := smallInt(2, 6)
		TL := T0 + s*L
		x0 := smallInt(1, L-1)
		ans := T0 + s*x0
		return generator.Problem{
			Question: fmt.Sprintf("A conductive rod of length \\(L = %d\\) satisfies the steady-state equation \\(u_{xx} = 0\\) with boundary temperatures \\(u(0) = %d\\) and \\(u(%d) = %d\\). Evaluate the temperature \\(u(%d)\\).", L, T0, L, TL, x0),
			Answer:   fmt.Sprintf("%d", ans),
			Explanation: fmt.Sprintf(
				"Steady-state heat distribution satisfies \\(u_{xx} = 0\\), producing the linear profile \\(u(x) = T_0 + \\frac{T_L - T_0}{L}x = %d + \\frac{%d - %d}{%d}x = %d + %dx\\). At \\(x = %d\\), \\(u(%d) = %d + %d(%d) = %d\\).",
				T0, TL, T0, L, T0, s, x0, x0, T0, s, x0, ans,
			),
		}

	case 3:
		// Count of interior nodal points
		n := smallInt(2, 8)
		nodalCount := n - 1
		return generator.Problem{
			Question: fmt.Sprintf("How many interior nodal points (coordinates where \\(u(x, t) = 0\\) strictly inside the open interval \\(0 < x < \\pi\\)) does the spatial harmonic \\(\\sin(%dx)\\) possess?", n),
			Answer:   fmt.Sprintf("%d", nodalCount),
			Explanation: fmt.Sprintf(
				"The eigenfunction \\(\\sin(%dx) = 0\\) occurs when \\(%dx = k\\pi\\) for integer \\(k\\). Inside the open domain \\(0 < x < \\pi\\), valid indices are \\(k = 1, 2, \\dots, %d\\). There are %d interior nodal points.",
				n, n, nodalCount, nodalCount,
			),
		}

	default:
		// Maximum principle bound
		L := smallInt(2, 8)
		tMax := smallInt(5, 15)
		M := smallInt(40, 95)
		return generator.Problem{
			Question: fmt.Sprintf("A heat diffusion process on the region \\([0, %d] \\times [0, %d]\\) has temperatures bounded by \\(M = %d\\) along the initial and spatial boundaries. By the maximum principle, what is the maximum possible temperature in the interior?", L, tMax, M),
			Answer:   fmt.Sprintf("%d", M),
			Explanation: fmt.Sprintf(
				"By the weak maximum principle for parabolic equations, the maximum of \\(u(x, t)\\) across a space-time domain is achieved on the parabolic boundary. Therefore the maximum interior temperature cannot exceed the boundary upper bound of %d.",
				M,
			),
		}
	}
}

type greenFunctionGen struct{}

func (g *greenFunctionGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	mode := rand.Intn(5)
	switch mode {
	case 0:
		// Scaled Green's function L * G(x, xi) = min(x, xi) * (L - max(x, xi))
		L := smallInt(6, 12)
		x := smallInt(1, L-1)
		xi := smallInt(1, L-1)
		mn := x
		mx := xi
		if xi < x {
			mn = xi
			mx = x
		}
		ans := mn * (L - mx)
		return generator.Problem{
			Question: fmt.Sprintf("For the boundary value problem \\(-y'' = f(x)\\) on \\([0, %d]\\) with Dirichlet conditions \\(y(0) = y(%d) = 0\\), the Green's function satisfies \\(L \\cdot G(x, \\xi) = \\min(x, \\xi)(L - \\max(x, \\xi))\\). Evaluate the integer quantity \\(%d \\cdot G(%d, %d)\\).", L, L, L, x, xi),
			Answer:   fmt.Sprintf("%d", ans),
			Explanation: fmt.Sprintf(
				"For \\(-y'' = f(x)\\) on \\([0, %d]\\), \\(G(x, \\xi) = \\frac{\\min(x, \\xi)(%d - \\max(x, \\xi))}{%d}\\). With \\(x = %d\\) and \\(\\xi = %d\\), the minimum is %d and the maximum is %d. Thus \\(%d \\cdot G(%d, %d) = %d(%d - %d) = %d \\cdot %d = %d\\).",
				L, L, L, x, xi, mn, mx, L, x, xi, mn, L, mx, mn, L-mx, ans,
			),
		}

	case 1:
		// Derivative jump discontinuity magnitude
		return generator.Problem{
			Question: "For the differential operator \\(L[y] = y''\\) governed by \\(y''(x) = \\delta(x - \\xi)\\), what is the jump discontinuity magnitude \\(\\left|\\left.\\frac{dG}{dx}\\right|_{x=\\xi^+} - \\left.\\frac{dG}{dx}\\right|_{x=\\xi^-}\\right|\\) across the point source \\(x = \\xi\\)?",
			Answer:   "1",
			Explanation: "Integrating the equation \\(G''(x, \\xi) = \\delta(x - \\xi)\\) over the infinitesimal interval \\((\\xi - \\epsilon, \\xi + \\epsilon)\\) yields \\(G'(\\xi^+) - G'(\\xi^-) = \\int_{\\xi-\\epsilon}^{\\xi+\\epsilon} \\delta(x-\\xi)\\,dx = 1\\). The jump magnitude is exactly 1.",
		}

	case 2:
		// Midpoint deflection under constant load: y(m) = c * m^2
		m := smallInt(1, 4)
		L := 2 * m
		c := smallInt(1, 5)
		k := 2 * c
		ans := c * m * m
		return generator.Problem{
			Question: fmt.Sprintf("A structural rod governed by \\(-y'' = %d\\) on \\([0, %d]\\) with boundary values \\(y(0) = y(%d) = 0\\) has deflection given by \\(y(x) = \\int_0^{%d} G(x, \\xi)(%d)\\,d\\xi\\). Evaluate the midpoint deflection \\(y(%d)\\).", k, L, L, L, k, m),
			Answer:   fmt.Sprintf("%d", ans),
			Explanation: fmt.Sprintf(
				"Integrating \\(-y'' = %d\\) twice with \\(y(0) = y(%d) = 0\\) yields \\(y(x) = \\frac{%d}{2}x(%d - x) = %d x(%d - x)\\). Evaluating at the midpoint \\(x = %d\\) gives \\(y(%d) = %d(%d)(%d - %d) = %d(%d)(%d) = %d\\).",
				k, L, k, L, c, L, m, m, c, m, L, m, c, m, m, ans,
			),
		}

	case 3:
		// First-order causal impulse jump
		a := smallInt(1, 5)
		k := smallInt(2, 9)
		t0 := smallInt(1, 5)
		return generator.Problem{
			Question: fmt.Sprintf("For the relaxation system \\(y' + %dy = %d \\delta(t - %d)\\) with initial rest state \\(y(0) = 0\\), what is the response value \\(y(%d^+)\\) immediately after the impulse?", a, k, t0, t0),
			Answer:   fmt.Sprintf("%d", k),
			Explanation: fmt.Sprintf(
				"Integrating \\(y' + %dy = %d\\delta(t - %d)\\) across the impulse yields a jump of \\(y(%d^+) - y(%d^-) = %d\\). Because the system is initially at rest with \\(y(%d^-) = 0\\), the post-impulse state is \\(y(%d^+) = %d\\).",
				a, k, t0, t0, t0, k, t0, t0, k,
			),
		}

	default:
		// Self-adjoint reciprocity relation
		x1 := smallInt(1, 4)
		xi1 := smallInt(5, 9)
		val := smallInt(2, 25)
		return generator.Problem{
			Question: fmt.Sprintf("A self-adjoint boundary value problem has symmetric Green's function \\(G(x, \\xi)\\). If \\(G(%d, %d) = %d\\), what is the value of \\(G(%d, %d)\\)?", x1, xi1, val, xi1, x1),
			Answer:   fmt.Sprintf("%d", val),
			Explanation: fmt.Sprintf(
				"By the reciprocity theorem for self-adjoint differential operators, the Green's function is symmetric: \\(G(x, \\xi) = G(\\xi, x)\\). Therefore \\(G(%d, %d) = G(%d, %d) = %d\\).",
				xi1, x1, x1, xi1, val,
			),
		}
	}
}

type sturmLiouvilleGen struct{}

func (g *sturmLiouvilleGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	type entry struct{ q, a, e string }
	easy := []entry{
		{"Are SL eigenfunctions orthogonal with weight? (yes/no)", "yes", "Orthogonal under the weight w."},
		{"Is Legendre equation Sturm-Liouville? (yes/no)", "yes", "With p=1-x^2."},
	}
	// Parameterized production: Dirichlet modes of y'' + λy = 0 on [0,π].
	ns := []int{1, 2, 3, 4}
	n := ns[rand.Intn(len(ns))]
	if rand.Intn(2) == 0 {
		easy = append(easy, entry{
			q: fmt.Sprintf("For y'' + λy = 0 with y(0) = y(π) = 0, λ_n = n^2. What is λ_%d? (enter a number)", n),
			a: fmt.Sprintf("%d", n*n),
			e: fmt.Sprintf("Eigenfunctions sin(nx) give λ_%d = %d^2 = %d.", n, n, n*n),
		})
	} else {
		easy = append(easy, entry{
			q: fmt.Sprintf("For y'' + λy = 0 with y(0) = y(π) = 0 and λ_%d = %d, what is the eigenfunction? (type like sin(2x))", n, n*n),
			a: fmt.Sprintf("sin(%dx)", n),
			e: fmt.Sprintf("The mode with λ_%d = %d is sin(%dx).", n, n*n, n),
		})
	}
	hard := []entry{{"Does SL on [a,b] with separated BC have discrete spectrum? (yes/no)", "yes", "Eigenvalues infinite discrete."}, {"Is Sturm-Liouville operator self-adjoint? (yes/no)", "yes", "SL is self-adjoint under weight."}, {"Does completeness of eigenfunctions hold? (yes/no)", "yes", "Expand in eigenfunctions."}}
	pool := easy
	if scale > 3 {
		pool = append(easy, hard...)
	}
	e := pool[rand.Intn(len(pool))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type boundaryValueGen struct{}

func (g *boundaryValueGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	scale := int(1 + ctx.Difficulty*4)
	type entry struct{ q, a, e string }
	easy := []entry{
		{"Is BVP different from IVP (conditions at two points)? (yes/no)", "yes", "BVP at boundaries."},
		{"Does BVP may have 0,1,∞ solutions? (yes/no)", "yes", "Unlike IVP."},
	}
	// Parameterized production: solve y'' = 0 with y(0) = 0, y(1) = b.
	bs := []int{1, 2, 3}
	b := bs[rand.Intn(len(bs))]
	sol := "y=x"
	if b != 1 {
		sol = fmt.Sprintf("y=%dx", b)
	}
	easy = append(easy, entry{
		q: fmt.Sprintf("Solve y'' = 0 with y(0) = 0 and y(1) = %d. (type like y=3x)", b),
		a: sol,
		e: fmt.Sprintf("y'' = 0 gives y = Cx+D; y(0) = 0 gives D = 0, y(1) = %d gives C = %d.", b, b),
	})
	hard := []entry{{"Does Green's function solve BVP? (yes/no)", "yes", "Via integral."}, {"Is shooting method for BVP? (yes/no)", "yes", "Convert BVP to IVP iteration."}, {"Does Fredholm alternative apply? (yes/no)", "yes", "Solvability condition."}}
	pool := easy
	if scale > 3 {
		pool = append(easy, hard...)
	}
	e := pool[rand.Intn(len(pool))]
	return generator.Problem{Question: e.q, Answer: e.a, Explanation: e.e}
}

type lyapunovGen struct{}

func (g *lyapunovGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	mode := rand.Intn(5)
	switch mode {
	case 0:
		// 1D orbital derivative dot(V) = -2 k x^2
		k := smallInt(1, 5)
		x0 := smallInt(1, 4)
		ans := -2 * k * x0 * x0
		return generator.Problem{
			Question: fmt.Sprintf("For the dynamical system \\(x' = -%dx\\) with candidate Lyapunov function \\(V(x) = x^2\\), evaluate the orbital time derivative \\(\\dot{V}(%d)\\).", k, x0),
			Answer:   fmt.Sprintf("%d", ans),
			Explanation: fmt.Sprintf(
				"The orbital derivative is \\(\\dot{V}(x) = \\frac{dV}{dx}x' = 2x(-%dx) = -%dx^2\\). Evaluating at \\(x = %d\\) yields \\(\\dot{V}(%d) = -%d(%d^2) = -%d(%d) = %d\\).",
				k, 2*k, x0, x0, 2*k, x0, 2*k, x0*x0, ans,
			),
		}

	case 1:
		// 2D linear system orbital derivative dot(V)
		a := smallInt(1, 4)
		b := smallInt(1, 4)
		x0 := smallInt(1, 3)
		y0 := smallInt(1, 3)
		ans := -2*a*x0*x0 - 2*b*y0*y0
		return generator.Problem{
			Question: fmt.Sprintf("For the 2D linear system \\(x' = -%dx\\) and \\(y' = -%dy\\) with Lyapunov function \\(V(x, y) = x^2 + y^2\\), evaluate the orbital time derivative \\(\\dot{V}(%d, %d)\\).", a, b, x0, y0),
			Answer:   fmt.Sprintf("%d", ans),
			Explanation: fmt.Sprintf(
				"The orbital time derivative is \\(\\dot{V} = \\nabla V \\cdot (x', y') = 2x(-%dx) + 2y(-%dy) = -%dx^2 - %dy^2\\). Substituting \\((x, y) = (%d, %d)\\) gives \\(\\dot{V}(%d, %d) = -%d(%d^2) - %d(%d^2) = -%d(%d) - %d(%d) = %d\\).",
				a, b, 2*a, 2*b, x0, y0, x0, y0, 2*a, x0, 2*b, y0, 2*a, x0*x0, 2*b, y0*y0, ans,
			),
		}

	case 2:
		// Basin radius r where dot(V) < 0 for x' = -a*x + b*x^3
		b := smallInt(1, 3)
		r := smallInt(2, 5)
		a := b * r * r
		return generator.Problem{
			Question: fmt.Sprintf("For the nonlinear system \\(x' = -%dx + %dx^3\\) with Lyapunov function \\(V(x) = x^2\\), the origin is strictly asymptotically stable on the symmetric open interval \\((-r, r)\\) where \\(\\dot{V}(x) < 0\\) for \\(x \\ne 0\\). What is the positive radius \\(r > 0\\)?", a, b),
			Answer:   fmt.Sprintf("%d", r),
			Explanation: fmt.Sprintf(
				"The orbital derivative is \\(\\dot{V}(x) = 2x(-%dx + %dx^3) = 2x^2(%dx^2 - %d)\\). For \\(x \\ne 0\\), strict negativity \\(\\dot{V}(x) < 0\\) requires \\(%dx^2 - %d < 0\\), which simplifies to \\(x^2 < %d\\), or \\(|x| < %d\\). Thus the basin radius is \\(r = %d\\).",
				a, b, b, a, b, a, r*r, r, r,
			),
		}

	case 3:
		// Quadratic Lyapunov energy evaluation V(x0, y0) = a*x0^2 + b*y0^2
		a := smallInt(1, 5)
		b := smallInt(1, 5)
		x0 := smallInt(1, 4)
		y0 := smallInt(1, 4)
		ans := a*x0*x0 + b*y0*y0
		return generator.Problem{
			Question: fmt.Sprintf("Given the positive definite Lyapunov energy function \\(V(x, y) = %dx^2 + %dy^2\\), evaluate the energy level \\(V(%d, %d)\\).", a, b, x0, y0),
			Answer:   fmt.Sprintf("%d", ans),
			Explanation: fmt.Sprintf(
				"Substituting the coordinates \\((x, y) = (%d, %d)\\) into the quadratic form gives \\(V(%d, %d) = %d(%d^2) + %d(%d^2) = %d(%d) + %d(%d) = %d + %d = %d\\).",
				x0, y0, x0, y0, a, x0, b, y0, a, x0*x0, b, y0*y0, a*x0*x0, b*y0*y0, ans,
			),
		}

	default:
		// Conservative Hamiltonian oscillator dot(E) = 0
		k := smallInt(2, 9)
		x0 := smallInt(1, 5)
		y0 := smallInt(1, 5)
		return generator.Problem{
			Question: fmt.Sprintf("For the undamped harmonic oscillator \\(x' = y\\) and \\(y' = -%dx\\) with total mechanical energy \\(E(x, y) = \\frac{%d}{2}x^2 + \\frac{1}{2}y^2\\), evaluate the orbital time derivative \\(\\dot{E}(%d, %d)\\).", k, k, x0, y0),
			Answer:   "0",
			Explanation: fmt.Sprintf(
				"Calculating the orbital derivative along trajectories: \\(\\dot{E} = \\frac{\\partial E}{\\partial x}x' + \\frac{\\partial E}{\\partial y}y' = (%dx)(y) + (y)(-%dx) = %dxy - %dxy = 0\\). Total energy is conserved along orbits, so \\(\\dot{E} = 0\\).",
				k, k, k, k,
			),
		}
	}
}

type laplaceEqGen struct{}

func (g *laplaceEqGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	mode := rand.Intn(5)
	switch mode {
	case 0:
		// Mean Value Property at circle center
		R := smallInt(2, 6)
		x0 := smallInt(1, 5)
		y0 := smallInt(1, 5)
		v := smallInt(3, 25)
		return generator.Problem{
			Question: fmt.Sprintf("A function \\(u(x, y)\\) is harmonic on a disk of radius \\(R = %d\\) centered at \\((%d, %d)\\). If its average boundary value satisfies \\(\\frac{1}{2\\pi R}\\oint_{\\partial B} u\\,ds = %d\\), what is the center value \\(u(%d, %d)\\)?", R, x0, y0, v, x0, y0),
			Answer:   fmt.Sprintf("%d", v),
			Explanation: fmt.Sprintf(
				"By the Gauss Mean Value Property for harmonic functions, the value at the center of a disk equals the average over its boundary circle: \\(u(%d, %d) = \\frac{1}{2\\pi R}\\oint_{\\partial B} u\\,ds = %d\\).",
				x0, y0, v,
			),
		}

	case 1:
		// Harmonic cubic polynomial coefficient b for u = a*x^3 + b*x*y^2
		a := smallInt(1, 5)
		b := -3 * a
		return generator.Problem{
			Question: fmt.Sprintf("For what integer coefficient \\(b\\) is the polynomial \\(u(x, y) = %dx^3 + bx y^2\\) a harmonic function (\\(\\Delta u = 0\\))?", a),
			Answer:   fmt.Sprintf("%d", b),
			Explanation: fmt.Sprintf(
				"Calculating second partial derivatives: \\(u_{xx} = 6(%d)x = %dx\\) and \\(u_{yy} = 2bx\\). Requiring \\(\\Delta u = u_{xx} + u_{yy} = %dx + 2bx = 0\\) yields \\(2b = -%d\\), so \\(b = %d\\).",
				a, 6*a, 6*a, 6*a, b,
			),
		}

	case 2:
		// Harmonic polynomial evaluation u(x, y) = x^2 - y^2
		x0 := smallInt(3, 8)
		y0 := smallInt(1, x0-1)
		ans := x0*x0 - y0*y0
		return generator.Problem{
			Question: fmt.Sprintf("Evaluate the harmonic polynomial \\(u(x, y) = x^2 - y^2\\) at the point \\((%d, %d)\\).", x0, y0),
			Answer:   fmt.Sprintf("%d", ans),
			Explanation: fmt.Sprintf(
				"The polynomial is harmonic because \\(\\Delta u = u_{xx} + u_{yy} = 2 - 2 = 0\\). Evaluating at \\((%d, %d)\\) gives \\(u(%d, %d) = %d^2 - %d^2 = %d - %d = %d\\).",
				x0, y0, x0, y0, x0, y0, x0*x0, y0*y0, ans,
			),
		}

	case 3:
		// Elliptic maximum principle interior bound
		M := smallInt(15, 60)
		return generator.Problem{
			Question: fmt.Sprintf("A harmonic function \\(u(x, y)\\) on the unit square satisfies boundary temperatures bounded by \\(u \\le %d\\). By the maximum principle, what is the maximum possible temperature at the interior center \\((0.5, 0.5)\\)?", M),
			Answer:   fmt.Sprintf("%d", M),
			Explanation: fmt.Sprintf(
				"By the strong maximum principle for elliptic equations, a non-constant harmonic function achieves its maximum strictly on the boundary. Therefore the interior temperature cannot exceed the boundary upper bound of %d.",
				M,
			),
		}

	default:
		// Total outward normal flux of harmonic function vanishes
		return generator.Problem{
			Question: "For a smooth harmonic function \\(u(x, y)\\) on a bounded planar domain \\(\\Omega\\), evaluate the total outward boundary normal flux \\(\\oint_{\\partial \\Omega} \\frac{\\partial u}{\\partial n}\\,ds\\).",
			Answer:   "0",
			Explanation: "By Green's first identity (the divergence theorem), \\(\\oint_{\\partial \\Omega} \\frac{\\partial u}{\\partial n}\\,ds = \\iint_\\Omega \\Delta u\\,dA\\). Because \\(u\\) is harmonic (\\(\\Delta u = 0\\)), the integral vanishes identically, yielding 0.",
		}
	}
}

type volterraGen struct{}

func (g *volterraGen) Generate(ctx generator.GeneratorContext) generator.Problem {
	mode := rand.Intn(5)
	switch mode {
	case 0:
		// Initial value y(0) = y0
		y0 := smallInt(2, 20)
		c := smallInt(1, 9)
		return generator.Problem{
			Question: fmt.Sprintf("Consider the Volterra integral equation of the second kind \\(y(x) = %d + %dx + \\int_0^x (x - t)y(t)\\,dt\\). Evaluate the initial value \\(y(0)\\).", y0, c),
			Answer:   fmt.Sprintf("%d", y0),
			Explanation: fmt.Sprintf(
				"Evaluating at the lower integration limit \\(x = 0\\), the definite integral \\(\\int_0^0 (0 - t)y(t)\\,dt = 0\\). Therefore \\(y(0) = %d + %d(0) + 0 = %d\\).",
				y0, c, y0,
			),
		}

	case 1:
		// Initial derivative y'(0) = m + c * k
		m := smallInt(1, 8)
		c := smallInt(1, 5)
		k := smallInt(1, 6)
		ans := m + c*k
		return generator.Problem{
			Question: fmt.Sprintf("For the Volterra integral equation \\(y(x) = %d + %dx + \\int_0^x %dy(t)\\,dt\\), evaluate the derivative value \\(y'(0)\\).", k, m, c),
			Answer:   fmt.Sprintf("%d", ans),
			Explanation: fmt.Sprintf(
				"Differentiating both sides with respect to \\(x\\) via the Leibniz integral rule gives \\(y'(x) = %d + %dy(x)\\). Evaluating at \\(x = 0\\) with \\(y(0) = %d\\) yields \\(y'(0) = %d + %d(%d) = %d + %d = %d\\).",
				m, c, k, m, c, k, m, c*k, ans,
			),
		}

	case 2:
		// Direct polynomial evaluation (a + c)*x0 + b
		a := smallInt(1, 5)
		b := smallInt(1, 9)
		c := smallInt(1, 5)
		x0 := smallInt(1, 6)
		ans := (a+c)*x0 + b
		return generator.Problem{
			Question: fmt.Sprintf("For the integral equation \\(y(x) = %dx + %d + \\int_0^x %d\\,dt\\), evaluate \\(y(%d)\\).", a, b, c, x0),
			Answer:   fmt.Sprintf("%d", ans),
			Explanation: fmt.Sprintf(
				"Evaluating the definite integral gives \\(\\int_0^x %d\\,dt = %dx\\). Thus \\(y(x) = %dx + %d + %dx = %dx + %d\\). Substituting \\(x = %d\\) yields \\(y(%d) = %d(%d) + %d = %d + %d = %d\\).",
				c, c, a, b, c, a+c, b, x0, x0, a+c, x0, b, (a+c)*x0, b, ans,
			),
		}

	case 3:
		// Equivalent ODE growth rate constant a = k
		y0 := smallInt(1, 9)
		k := smallInt(2, 8)
		return generator.Problem{
			Question: fmt.Sprintf("Differentiating the Volterra integral equation \\(y(x) = %d + \\int_0^x %dy(t)\\,dt\\) converts it into the initial value problem \\(y' = ay\\) with \\(y(0) = %d\\). What is the value of coefficient \\(a\\)?", y0, k, y0),
			Answer:   fmt.Sprintf("%d", k),
			Explanation: fmt.Sprintf(
				"By the fundamental theorem of calculus, differentiating \\(y(x) = %d + \\int_0^x %dy(t)\\,dt\\) yields \\(y'(x) = %dy(x)\\). Therefore the growth coefficient is \\(a = %d\\).",
				y0, k, k, k,
			),
		}

	default:
		// Second derivative y''(0) = c * k
		k := smallInt(2, 7)
		c := smallInt(2, 6)
		ans := c * k
		return generator.Problem{
			Question: fmt.Sprintf("Consider the Volterra integral equation with displacement kernel \\(y(x) = %d + \\int_0^x %d(x - t)y(t)\\,dt\\). Evaluate the second derivative \\(y''(0)\\).", k, c),
			Answer:   fmt.Sprintf("%d", ans),
			Explanation: fmt.Sprintf(
				"Differentiating twice with respect to \\(x\\) via the Leibniz rule: \\(y'(x) = %d\\int_0^x y(t)\\,dt\\) and \\(y''(x) = %dy(x)\\). At \\(x = 0\\), the integral equation gives \\(y(0) = %d\\). Therefore \\(y''(0) = %dy(0) = %d(%d) = %d\\).",
				c, c, k, c, c, k, ans,
			),
		}
	}
}
