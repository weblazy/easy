package dsn

import (
	"net/url"
	"strings"

	"emperror.dev/errors"
	"github.com/weblazy/easy/utils/db/mysql/manager"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	errInvalidDSNUnescaped                   = errors.New("invalid DSN: did you forget to escape a param value")
	errInvalidDSNAddr                        = errors.New("invalid DSN: network address not terminated (missing closing brace)")
	errInvalidDSNNoSlash                     = errors.New("invalid DSN: missing the slash separating the database name")
	_                      manager.DSNParser = (*MysqlDSNParser)(nil)
)

type MysqlDSNParser struct {
}

func init() { //nolint:gochecknoinits
	manager.Register(&MysqlDSNParser{})
}

func (m *MysqlDSNParser) Scheme() string {
	return "mysql"
}

func (m *MysqlDSNParser) GetDialector(dsn string) gorm.Dialector {
	return mysql.Open(dsn)
}

func (m *MysqlDSNParser) ParseDSN(dsn string) (cfg *manager.DSN, err error) {
	// New config with some default values
	cfg = new(manager.DSN)

	// [user[:password]@][net[(addr)]]/dbname[?param1=value1&paramN=valueN]
	// Find the last '/' (since the password or the net addr might contain a '/')
	foundSlash := false
	for i := len(dsn) - 1; i >= 0; i-- {
		if dsn[i] == '/' { //nolint
			foundSlash = true
			var j int

			// left part is empty if i <= 0
			if i > 0 {
				// [username[:password]@][protocol[(address)]]
				// Find the last '@' in dsn[:i]
				for j = i; j >= 0; j-- {
					if dsn[j] == '@' {
						parseUsernamePassword(cfg, dsn[:j])
						break
					}
				}

				// [protocol[(address)]]
				// Find the first '(' in dsn[j+1:i]
				if err = parseAddrNet(cfg, dsn[j:i]); err != nil {
					return
				}
			}

			// dbname[?param1=value1&...&paramN=valueN]
			// Find the first '?' in dsn[i+1:]
			for j = i + 1; j < len(dsn); j++ {
				if dsn[j] == '?' {
					if err = parseDSNParams(cfg, dsn[j+1:]); err != nil {
						return
					}
					break
				}
			}
			cfg.DBName = dsn[i+1 : j]

			break
		}
	}
	if !foundSlash && len(dsn) > 0 {
		return nil, errInvalidDSNNoSlash
	}
	return
}

// username[:password]
func parseUsernamePassword(cfg *manager.DSN, userPassStr string) {
	for i := 0; i < len(userPassStr); i++ {
		if userPassStr[i] == ':' {
			cfg.Password = userPassStr[i+1:]
			cfg.User = userPassStr[:i]
			break
		}
	}
}

// [protocol[(address)]]
func parseAddrNet(cfg *manager.DSN, addrNetStr string) error {
	for i := 0; i < len(addrNetStr); i++ {
		if addrNetStr[i] == '(' {
			// dsn[i-1] must be == ')' if an address is specified
			if addrNetStr[len(addrNetStr)-1] != ')' {
				if strings.ContainsRune(addrNetStr[i+1:], ')') {
					return errInvalidDSNUnescaped
				}
				return errInvalidDSNAddr
			}
			cfg.Addr = addrNetStr[i+1 : len(addrNetStr)-1]
			cfg.Net = addrNetStr[1:i]
			break
		}
	}
	return nil
}

// param1=value1&...&paramN=valueN
func parseDSNParams(cfg *manager.DSN, params string) (err error) {
	for _, v := range strings.Split(params, "&") {
		param := strings.SplitN(v, "=", 2)
		if len(param) != 2 {
			continue
		}
		// lazy init
		if cfg.Params == nil {
			cfg.Params = make(map[string]string)
		}
		value := param[1]
		if cfg.Params[param[0]], err = url.QueryUnescape(value); err != nil {
			return
		}
	}
	return
}
