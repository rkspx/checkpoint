package checkpoint

type LoginResponse struct {
	SID              string `json:"sid"`
	APIServerVersion string `json:"api-server-version"`
	LastLogin        struct {
		ISO8601 string `json:"iso-8601"`
		Posix   int    `json:"posix"`
	} `json:"last-login-was-at"`
	LoginMessage   string `json:"loginMessage"`
	ReadOnly       bool   `json:"read-only"`
	SessionTimeout int    `json:"session-timeout"`
	Standby        bool   `json:"standby"`
	UID            string `json:"uid"`
	URL            string `json:"string"`
}
