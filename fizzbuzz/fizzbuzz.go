package main

import (
	"bufio"
	"fmt"
	"linusphoenix/fizzbuzz/fizzbuzz"
	"os"
	"strconv"
)

func main() {
	stdin := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter the divisor for Fizz: ")
	stdin.Scan()
	fizzNumber, err := strconv.Atoi(stdin.Text())
	if err != nil {
		fmt.Println("The divisor for Fizz must be an integer")
		os.Exit(1)
	}

	fmt.Print("Enter the divisor for Buzz: ")
	stdin.Scan()
	buzzNumber, err := strconv.Atoi(stdin.Text())
	if err != nil {
		fmt.Println("The divisor for Buzz must be an integer")
		os.Exit(1)
	}

	f, err := fizzbuzz.New(fizzNumber, buzzNumber)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for i := 1; i <= 100; i++ {
		fmt.Println(f.Eval(i))
	}
}
