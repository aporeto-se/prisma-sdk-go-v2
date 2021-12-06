package token

// Config config
type Config struct {
	TokenString string
}

// NewConfig returns new Config
func NewConfig() *Config {
	return &Config{}
}

// SetTokenString sets attribute and returns self
func (t *Config) SetTokenString(v string) *Config {
	t.TokenString = v
	return t
}
