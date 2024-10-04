package markdown

import (
	"fmt"

	markdown "github.com/MichaelMure/go-term-markdown"
)

func Render(readme string) {
	markdown := markdown.Render(readme, 240, 6)

	fmt.Println(string(markdown))
}
