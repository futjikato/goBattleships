package bots

import (
	"github.com/patrikeh/go-deep/training"
)

func getHitSurroundingExamples() training.Examples {
	examples := make(training.Examples, 0)
	for i := 0; i < 100; i++ {
		scenario := map[int]string{i: " XX "}
		ex := training.Example{Input: matrixToInput(scenario), Response: matrixToOutput(getSurroundingShotsMatrix(scenario))}
		examples = append(examples, ex)
		if i > 9 {
			o := i - 10
			scenario := map[int]string{i: " XX ", o: " XX "}
			ex := training.Example{Input: matrixToInput(scenario), Response: matrixToOutput(getSurroundingShotsMatrix(scenario))}
			examples = append(examples, ex)
		}
		if i%10 > 0 {
			o := i - 1
			scenario := map[int]string{i: " XX ", o: " XX "}
			ex := training.Example{Input: matrixToInput(scenario), Response: matrixToOutput(getSurroundingShotsMatrix(scenario))}
			examples = append(examples, ex)
		}
		if i%10 < 9 {
			o := i + 1
			scenario := map[int]string{i: " XX ", o: " XX "}
			ex := training.Example{Input: matrixToInput(scenario), Response: matrixToOutput(getSurroundingShotsMatrix(scenario))}
			examples = append(examples, ex)
		}
		if i < 89 {
			o := i + 10
			scenario := map[int]string{i: " XX ", o: " XX "}
			ex := training.Example{Input: matrixToInput(scenario), Response: matrixToOutput(getSurroundingShotsMatrix(scenario))}
			examples = append(examples, ex)
		}
	}

	return examples
}

func getSurroundingShotsMatrix(matrix map[int]string) map[int]string {
	surroundings := make(map[int]string)
	for i := 0; i < 100; i++ {
		if v, hasV := matrix[i]; hasV == true && v == " XX " {
			if i > 9 {
				surroundings[i-10] = " XX "
			}
			if i%10 > 0 {
				surroundings[i-1] = " XX "
			}
			if i%10 < 9 {
				surroundings[i+1] = " XX "
			}
			if i < 89 {
				surroundings[i+10] = " XX "
			}
		}
	}
	for i, _ := range matrix {
		surroundings[i] = " OO "
	}

	return surroundings
}

func getFreeForAllExamples() training.Examples {
	examples := make(training.Examples, 100)
	for i := 0; i < 100; i++ {
		scenario := map[int]string{i: " OO "}
		examples[i] = training.Example{Input: matrixToInput(scenario), Response: matrixToOutput(getFreeShotsMatrix(scenario))}
	}

	return examples
}

func getFreeShotsMatrix(matrix map[int]string) map[int]string {
	free := make(map[int]string)
	for i := 0; i < 100; i++ {
		if _, hasV := matrix[i]; hasV == true {
			free[i] = " OO "
		} else {
			free[i] = " XX "
		}
	}

	return free
}

func getValidation1() training.Example {
	input := map[int]string{15: " XX "}
	output := map[int]string{5: " XX ", 14: " XX ", 16: " XX ", 25: " XX "}
	return training.Example{Input: matrixToInput(input), Response: matrixToOutput(output)}
}

func getValidation2() training.Example {
	input := map[int]string{15: " XX ", 16: " OO "}
	output := map[int]string{5: " XX ", 14: " XX ", 25: " XX "}
	return training.Example{Input: matrixToInput(input), Response: matrixToOutput(output)}
}

func getValidation3() training.Example {
	input := map[int]string{6: " XX "}
	output := map[int]string{5: " XX ", 7: " XX ", 16: " XX "}
	return training.Example{Input: matrixToInput(input), Response: matrixToOutput(output)}
}
