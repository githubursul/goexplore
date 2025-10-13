package main

import (
	"fmt"

	"github.com/githubursul/goexplore/gonative"
)

func main() {
	fmt.Println("Hello, World!")
	factorialValue := gonative.Factorial(5)
	fmt.Printf("Factorial of 5 is: %d\n", factorialValue)
}
