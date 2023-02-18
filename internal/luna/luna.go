package luna
 // Valid check number is valid or not based on Luhn algorithm
import (
	"fmt"
	"math/rand"
	"time"
)

// Valid check number is valid or not based on Luhn algorithm
func Valid(number int) bool {
	return (number%10+checksum(number/10))%10 == 0
}

func checksum(number int) int {
	var luhn int

	for i := 0; number > 0; i++ {
		cur := number % 10

		if i%2 == 0 { // even
			cur = cur * 2
			if cur > 9 {
				cur = cur%10 + cur/10
			}
		}

		luhn += cur
		number = number / 10
	}
	return luhn % 10
}
func GetOrderNumber() int {
	rand.Seed(time.Now().UnixNano())
	i := 0
	n := 0
	j := 0
	j = rand.Intn(100)
	fmt.Println("Answer: ", Valid(n))
	for {
		i++
		n = rand.Intn(1000000)
		if Valid(n) && n > 100000 && i > j{
			fmt.Println("Number is: ", n, "; index is: ", i)
			break
		}
	}
	return n
}