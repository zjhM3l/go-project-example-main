package firststep

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

func guess() {
	max := 100
	secretNumber := rand.Intn(max)
	fmt.Println("The secret number is", secretNumber)

	fmt.Println("Guess the number")
	reader := bufio.NewReader(os.Stdin)
	for {
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("An error occurred while reading input. Please try again", err)
			continue
		}
		input = strings.TrimSuffix(input, "\r\n")

		guess, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println(err)
			fmt.Println("Invalid input. Please enter a number")
			continue
		}
		fmt.Println("You guessed", guess)
		if guess < secretNumber {
			fmt.Println("Your guess is less than the secret number")
		} else if guess > secretNumber {
			fmt.Println("Your guess is greater than the secret number")
		} else {
			fmt.Println("Congratulations! You guessed the secret number")
			break
		}
	}
}
