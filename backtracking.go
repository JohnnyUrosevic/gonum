// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package optimize

const (
	defaultBacktrackingDecrease  = 0.5
	defaultBacktrackingFuncConst = 1e-4
	minimumBacktrackingStepSize  = 1e-20
)

// Backtracking is a Linesearcher that uses a backtracking to find a point that
// satisfies the Armijo condition with the given function constant FuncConst. If
// the Armijo condition has not been met, the step size is decreased by a
// factor of Decrease.
//
// The Armijo condition only requires the gradient at the beginning of each
// major iteration (not at successive step locations), and so Backtracking may
// be a good linesearch for functions with expensive gradients. Backtracking is
// not appropriate for optimizers that require the Wolfe conditions to be met,
// such as BFGS.
//
// Both FuncConst and Decrease must be between zero and one, and Backtracking will
// panic otherwise. If either FuncConst or Decrease are zero, it will be set to a
// reasonable default.
type Backtracking struct {
	FuncConst float64 // Necessary function descrease for Armijo condition.
	Decrease  float64 // Step size multiplier at each iteration (stepSize *= Decrease).

	stepSize float64
	initF    float64
	initG    float64
}

func (b *Backtracking) Init(f, g float64, step float64) EvaluationType {
	if step <= 0 {
		panic("backtracking: bad step size")
	}
	if g >= 0 {
		panic("backtracking: initial derivative is non-negative")
	}

	if b.Decrease == 0 {
		b.Decrease = defaultBacktrackingDecrease
	}
	if b.FuncConst == 0 {
		b.FuncConst = defaultBacktrackingFuncConst
	}
	if b.Decrease <= 0 || b.Decrease >= 1 {
		panic("backtracking: Decrease must be between 0 and 1")
	}
	if b.FuncConst <= 0 || b.FuncConst >= 1 {
		panic("backtracking: FuncConst must be between 0 and 1")
	}

	b.stepSize = step
	b.initF = f
	b.initG = g
	return FuncEvaluation
}

func (b *Backtracking) Finished(f, _ float64) bool {
	return ArmijoConditionMet(f, b.initF, b.initG, b.stepSize, b.FuncConst)
}

func (b *Backtracking) Iterate(_, _ float64) (float64, EvaluationType, error) {
	b.stepSize *= b.Decrease
	if b.stepSize < minimumBacktrackingStepSize {
		return 0, NoEvaluation, ErrLinesearchFailure
	}
	return b.stepSize, FuncEvaluation, nil
}
