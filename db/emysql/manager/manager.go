package manager

var (
	// m is a map from scheme to dsn builder.
	m = make(map[string]DSNParser)
)

func Register(b DSNParser) {
	m[b.Scheme()] = b
}

// Get returns the dsn builder registered with the given scheme.
//
// If no builder is register with the scheme, nil will be returned.
func Get(scheme string) DSNParser {
	if b, ok := m[scheme]; ok {
		return b
	}
	return nil
}
