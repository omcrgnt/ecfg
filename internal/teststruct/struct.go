package teststruct

type Config struct {
	Port int    `ecfg:"PORT" usage:"Server port"`
	User string `ecfg:"USER" usage:"User name"`
}
