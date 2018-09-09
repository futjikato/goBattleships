package bots

import (
	"fmt"
	"math/rand"

	"github.com/futjikato/goBattleships/battleships"
	deep "github.com/patrikeh/go-deep"
	"github.com/patrikeh/go-deep/training"
)

type Bot struct {
	// neural network
	n *deep.Neural
	// history of predictions
	h []int
	// history of random values chosen
	rh []int
}

func New() *Bot {
	n := getNetwork()
	train(n)
	return &Bot{
		n:  n,
		h:  make([]int, 0),
		rh: make([]int, 0),
	}
}

func (bot *Bot) NextHit(board *battleships.Board) (int, int) {
	m := board.GetHitMatrix()
	input := matrixToInput(m)
	out := bot.n.Predict(input)

	max := 0.4
	maxi := -1
	for i := 0; i < 100; i++ {
		if out[i] > max {
			max = out[i]
			maxi = i
		}
	}

	maxi = bot.checkHitIndex(maxi)

	x := maxi % 10
	y := maxi / 10
	fmt.Printf("%d => (%d, %d)\n", maxi, x, y)
	return x, y
}

func (bot *Bot) checkHitIndex(maxi int) int {
	bot.h = append(bot.h, maxi)

	if maxi == -1 {
		ri := rand.Intn(100)
		bot.rh = append(bot.rh, ri)
		return ri
	}

	return maxi
}

func (bot *Bot) PrintHistory() {
	fmt.Printf("Predictions (%d): %+v\n", len(bot.h), bot.h)
	fmt.Printf("Randoms (%d): %+v\n", len(bot.rh), bot.rh)
}

func matrixToInput(matrix map[int]string) []float64 {
	input := make([]float64, 200)
	for i := 0; i < 100; i++ {
		if h, ok := matrix[i]; ok == true {
			if h == " XX " {
				input[i] = 1.0
			} else {
				input[100+i] = 1.0
			}
		}
	}

	return input
}

func matrixToOutput(matrix map[int]string) []float64 {
	output := make([]float64, 100)
	for i := 0; i < 100; i++ {
		if h, ok := matrix[i]; ok == true && h == " XX " {
			output[i] = 1.0
		} else {
			output[i] = 0
		}
	}

	return output
}

func getNetwork() *deep.Neural {
	return deep.NewNeural(&deep.Config{
		/* Input dimensionality: first 100 for previous hits next 100 for known misses */
		Inputs: 200,
		/* One hidden row and one output row */
		Layout: []int{200, 100},
		/* Activation functions: {deep.Sigmoid, deep.Tanh, deep.ReLU, deep.Linear} */
		Activation: deep.ActivationTanh,
		/* Determines output layer activation & loss function:
		ModeRegression: linear outputs with MSE loss
		ModeMultiClass: softmax output with Cross Entropy loss
		ModeMultiLabel: sigmoid output with Cross Entropy loss
		ModeBinary: sigmoid output with binary CE loss */
		Mode: deep.ModeBinary,
		/* Weight initializers: {deep.NewNormal(μ, σ), deep.NewUniform(μ, σ)} */
		Weight: deep.NewNormal(0.5, 0.0),
		/* Apply bias */
		Bias: false,
	})
}

func train(n *deep.Neural) {
	optimizer := training.NewSGD(0.03, 0.1, 1e-6, true)
	trainer := training.NewTrainer(optimizer, 100)

	examples := append([]training.Example{}, getFreeForAllExamples()...)
	examples = append(examples, getHitSurroundingExamples()...)

	validations := []training.Example{
		getValidation1(),
		getValidation2(),
		getValidation3(),
	}

	trainer.Train(n, examples, validations, 500)
}
