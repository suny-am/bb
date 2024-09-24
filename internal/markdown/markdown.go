package markdown

import (
	"fmt"

	markdown "github.com/MichaelMure/go-term-markdown"
)

func Render(readme string) {
	markdown := markdown.Render(readme, 80, 6)

	fmt.Println(string(markdown))
}
