/*
Copyright © 2023 Daniel Jay Haskin <me@djha.skin>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	math/bits
	math/rand
)

const (
	INITIAL_MAX_SIZE = 16
)

struct Solution {
	bitstring []uint8
	fitness float64
}

struct Population {
	solutions []Solution
	max_size int
	best_fitness float64
	best_fitness_index int
	iterations_without_improvement int
}

interface FitnessFunction {
	func rank([]uint8) float64
}


func flip_bit(bitstring []uint8, position int) {
	uint8_index := position / 8
	bit_index := position % 8
	bitstring[uint8_index] ^= 1 << uint(bit_index)
}

func flip_random_bits(bitstring []uint8, num_bits int) {
	for i := 0; i < num_bits; i++ {
		flip_bit(bitstring, rand.Intn(len(bitstring)*8))
	}
}

func single_crossover(parent1 []uint8, parent2 []uint8) uint8[] {
	child := make([]uint8, len(parent1))

	crossover_point := rand.Intn(len(parent1)*8)
	// First, set the uint8s before the crossover point
	before_crossover := crossover_point / 8
	for i := 0; i < before_crossover; i++ {
		child[i] = parent1[i]
	}
	// Second, set the uint8s after the crossover point
	for i := before_crossover+1; i < len(parent1); i++ {
		child[i] = parent2[i]
	}
	// Third, set the bits in the crossover uint8
	uint8_crossover_point = crossover_point % 8
	var mask uint8 = 0
	for i := 0; i < uint8_crossover_point; i++ {
		mask |= 1 << uint(i)
	}
	child[before_crossover] = (parent1[before_crossover] & mask) | (parent2[before_crossover] & ^mask)
	return child
}

func same_bits(a uint8, b uint8) int {
	distance := 0
	the_same := ~a ^ b
	return bits.OnesCount8(the_same)
}

func same_bits_in_bytes(bitstring1 []uint8, bitstring2 []uint8) int {
	distance := 0
	for i := 0; i < len(bitstring1); i++ {
		distance += same_bits(bitstring1[i], bitstring2[i])
	}
	return distance
}

func mutate_solution(spot1 int, spot2 int, population *Population, fitness *FitnessFunction) {
	// First, we determine "how strongly" we wish to mutate the solution.
	// This function is based on the entropy equation for two bitstrings.
	// It is bowl-shaped, which is nice.
	mutation_urge_base := 2*(.5*same_bits_in_bytes(population.solutions[spot1].bitstring,
		population.solutions[spot2].bitstring)-.5)
	mutation_urge := mutation_urge_base * mutation_urge_base


	var mutation_spot int
	if spot1 == population.best_fitness_index {
		// Don't mutate the best solution, if only to ensure that
		// the `best_fitness_index` is always correct.
		// Otherwise, if we mutate the best solution at spot 1,
		// and the new solution isn't better than the old one, then
		// the `best_fitness_index` will not get updated,
		// but will be pointing at a solution that is no longer the best.
		mutation_spot = spot2
	} else {
		mutation_spot = spot1
	}
	// If the urge is close to 1 (SUPER STRONG), flip 10% of the bits.
	// If the urge is close to 0 (SUPER WEAK), flip close to 0% of the bits,
	// except that we must flip at least 1 bit.
	mutated_bits = max(int(mutation_urge * .1 * len(population.solutions[spot1])),
		1)
	flip_random_bits(population.solutions[spot1].bitstring, mutated_bits)
	population.solutions[spot1].fitness = fitness.rank(population.solutions[spot1].bitstring)
	// Update the best seen solution if necessary.
	if population.solutions[spot1].fitness > population.best_fitness {
		population.best_fitness = population.solutions[spot1].fitness
		population.best_fitness_index = spot1
	}
}

func breed_and_kill(spot1 int, spot2 int, population *Population, fitness *FitnessFunction) {
	// Make a kid.
	kids_DNA := single_crossover(population.solutions[spot1], population.solutions[spot2])
	kids_fitness := fitness.rank(kids_DNA)

	kid := Solution{bitstring: kids_DNA, fitness: kids_fitness}
	if kids_fitness > population.solutions[spot1].fitness and
	   kids_fitness > population.solutions[spot2].fitness {

	// If our population is already full, we need to make room for the kid.
	if len(population.solutions) >= population.max_size {
		// If the kid is better than the worst solution, we replace the worst
		// solution with the kid.
		if kids_fitness > population.solutions[spot1].fitness and
		   kids_fitness > population.solutions[spot2].fitness {
			if population.solutions[spot1].fitness < population.solutions[spot2].fitness {
				population.solutions[spot1].bitstring = kids_DNA
				population.solutions[spot1].fitness = kids_fitness
			} else {
				population.solutions[spot2].bitstring = kids_DNA
				population.solutions[spot2].fitness = kids_fitness
			}
			population.best_fitness = kids_fitness
			population.best_fitness_index = spot1
		} else if kids_fitness > population.solutions[spot1].fitness {
			population.solutions[spot1].bitstring = kids_DNA
			population.solutions[spot1].fitness = kids_fitness
		} else if kids_fitness > population.solutions[spot2].fitness {
			population.solutions[spot2].bitstring = kids_DNA
			population.solutions[spot2].fitness = kids_fitness
		}

	} else {

		// If the population isn't full, we just add the kid to the population.
		population.solutions = append(population.solutions, Solution{bitstring: kids_DNA, fitness: kids_fitness})
		if kids_fitness > population.best_fitness {
			population.best_fitness = kids_fitness
			population.best_fitness_index = len(population.solutions)-1
		}
	}
}

func generate_random_solution(num_bits int, fitness *FitnessFunction) *Solution {
	bitstring := make([]uint8, bits.Len8(num_bits)
	for i := 0; i < num_bits; i++ {
		// Shuffle good and proper.
		flip_random_bits(bitstring, len(bitstring) * 3)
	}
	new_solution := &Solution{bitstring: bitstring, fitness: fitness.rank(bitstring)}
	return solution
}

func find_best_solution(num_bits int, fitness *FitnessFunction,
	max_iterations_without_improvement int) *Solution {

	var solutions []Solution
	var best_fitness float64 = 0
	var best_fitness_index int = 0

	for i := 0; i < INITIAL_MAX_SIZE; i++ {
		solution := generate_random_solution(num_bits, fitness)
		solutions = append(solutions, solution, fitness))
		if solution.fitness > best_fitness {
			best_fitness = solution.fitness
			best_fitness_index = i
		}
	}

	population = Population{
		solutions: solutions,
		max_size: INITIAL_MAX_SIZE,
		best_fitness: best_fitness,
		best_fitness_index: best_fitness_index,
		iterations_without_improvement: 0,
	}

	var spot1, spot2 int
	var previous_record = population.best_fitness
	for population.iterations_without_improvement < max_iterations_without_improvement {
		spot1 = rand.Intn(len(population.solutions))
		spot2 = rand.Intn(len(population.solutions))
		for spot1 == spot2 {
			spot2 = rand.Intn(len(population.solutions))
		}
		// First, we mutate one of the parents, maybe, depending on how similar
		// the parent are.
		mutate_solution(spot1, spot2, population, fitness)
		breed_and_kill(spot1, spot2, population, fitness)
		if population.best_fitness > previous_record {
			population.iterations_without_improvement = 0
			previous_record = population.best_fitness
		} else {
			population.iterations_without_improvement++
		}
		if population.iterations_without_improvement > 3*len(population.solutions) {
			// We've been stuck for a while. Let's try to get unstuck.
			// We'll do this by expanding the population size.
			// This effectively descreases the learning rate,
			// Allowing the algorithm to "fine tune" what it has learned.
			population.max_size = population.max_size * 2
			population.iterations_without_improvement = 0
		}
	}
	return &population.solutions[population.best_fitness_index]
}
