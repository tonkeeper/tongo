package utils

import (
	"fmt"
	"strconv"
)

func HumanFriendlyCoinsRepr(amount int64) string {
	k := int64(1)
	currentAmount := amount
	for _, unit := range []string{"nanoTON", "microTON", "milliTON", "TON", "kiloTON", "megaTON", "gigaTON"} {
		if currentAmount < 1000 {
			fraction := amount - currentAmount*k
			if fraction == 0 {
				return fmt.Sprintf("%d %v", currentAmount, unit)
			}
			for amount%10 == 0 {
				amount /= 10
			}
			amountRepr := strconv.FormatInt(amount, 10)
			integerRepr := strconv.FormatInt(currentAmount, 10)
			repr := amountRepr[:len(integerRepr)] + "." + amountRepr[len(integerRepr):]
			return fmt.Sprintf("%s %v", repr, unit)
		}
		k *= 1000
		currentAmount /= 1000
	}
	return fmt.Sprintf("%v nanoTON", amount)
}
