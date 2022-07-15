package handler

// Option is an interface for dependency injection.
type Option interface {
	apply(c *Handler)
}

type Config struct {
	MysqlDSN    string
	CdnEndpoint string
	UserJwtSalt	string
	GroupSalt   string
	UserSalt    string
}

// OptionFunc is a function for Option interface.
type OptionFunc func(h *Handler)

// anonymous func, it has only input
//

func (o OptionFunc) apply(h *Handler) { o(h) }

// WithConfig returns a function for setting config.
func WithConfig(c *Config) OptionFunc {
	return func(h *Handler) { h.Cfg = c }
}

//return handler setting function
//annonymous func

/*
func PrettyPrint(target ...interface{}) {
	replacer := strings.NewReplacer("\t", " ", "\n", "")

	for _, t := range target {
		if _t, ok := t.(string); ok {
			t = replacer.Replace(_t)
		}
		s, _ := json.MarshalIndent(t, "", "\t")
		fmt.Println(string(s))
	}
}
*/
// Valid checks to be correct config or not.
func (cfg *Config) Valid() (ok bool) {
	if cfg.CdnEndpoint == "" {
		// temporary remove mySQLDSN
		return
	}
	// Check a salt for encrypting data exists or not.
	if cfg.UserJwtSalt == "" || cfg.GroupSalt == "" || cfg.UserSalt == "" {
		return
	}

	ok = true
	return
}
