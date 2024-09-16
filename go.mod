module github.com/suny-am/bitbucket-cli

go 1.22.7

require (
	github.com/go-resty/resty/v2 v2.15.0
	github.com/joho/godotenv v1.5.1
)

require golang.org/x/net v0.27.0 // indirect

replace example.com/greetings => ../greetings
