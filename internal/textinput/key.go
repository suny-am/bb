package textinput

import (
	"bufio"
	"fmt"
	"os"
)

func ReadKey() rune {
	reader := bufio.NewReader(os.Stdin)
	char, _, err := reader.ReadRune()
	if err != nil {
		fmt.Println("Error reading key: ", err)
	}
	return char
}

func ConfirmKey() {
loop:
	for {
		_ = ReadKey()
		break loop
	}
}
