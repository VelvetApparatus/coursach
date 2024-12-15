package exps

import (
	"course/optimizer"
	"course/optimizer/bruteforce"
	"course/optimizer/gen_algorithm"
	"course/optimizer/greedy"
)

func Optimizers() map[string]optimizer.Optimizer {
	return map[string]optimizer.Optimizer{
		"bruteforce":       bruteforce.NewBrutForceOptimizer(),
		"greedy":           greedy.NewGreedyOptimizer(),
		"gen_algorithm_bf": gen_algorithm.New(bruteforce.NewBrutForceOptimizer()),
		"gen_algorithm_gr": gen_algorithm.New(greedy.NewGreedyOptimizer()),
	}
}
