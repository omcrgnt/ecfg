package example

type Config struct {
	Some1 string `ecfg:"some1" usage:"some1 usage"`
	Some2 struct {
		Another1 string
		Another2 int64
	} `ecfg:"some2" usage:"some2 usage"`
}
