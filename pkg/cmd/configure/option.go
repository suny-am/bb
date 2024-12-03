package configure

import (
	"fmt"

	"github.com/suny-am/bb/internal/config"
	"github.com/suny-am/bb/internal/textinput"
)

func configureOption(option string) {
	currentConfigValue := config.GetConfigOption(option)
	textinput.ConfigListen(fmt.Sprintf("Enter new value for %s", option), option, currentConfigValue)
}
