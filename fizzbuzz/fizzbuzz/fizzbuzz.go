package fizzbuzz

import (
	"errors"
	"strconv"
)

const DefaultFizz = 3
const DefaultBuzz = 5

type Fizzbuzz struct {
	fizzNumber int
	buzzNumber int
}

func Default() *Fizzbuzz {
	return &Fizzbuzz{DefaultFizz, DefaultBuzz}
}

func New(fizzNumber int, buzzNumber int) (*Fizzbuzz, error) {
	if fizzNumber == 0 {
		return nil, errors.New("Fizz number cannot be 0")
	}
	if buzzNumber == 0 {
		return nil, errors.New("Buzz number cannot be 0")
	}
	return &Fizzbuzz{fizzNumber, buzzNumber}, nil
}

func (fizzbuzz *Fizzbuzz) Eval(i int) string {
	output := ""
	if i%fizzbuzz.fizzNumber == 0 {
		output += "Fizz"
	}
	if i%fizzbuzz.buzzNumber == 0 {
		output += "Buzz"
	}

	if output == "" {
		return strconv.Itoa(i)
	} else {
		return output
	}
}
