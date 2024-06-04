package consts

import "time"

type IPVersion int
type Protocol string

const DefaultTimeout = 3 * time.Second

const (
	IPV4 IPVersion = 4
	IPV6 IPVersion = 6
)

const (
	HTTP  Protocol = "http"
	HTTPS Protocol = "https"
)
