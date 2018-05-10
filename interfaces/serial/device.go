package serial

import "time"

type Device interface {
	GetBaud() int
	GetPort() string
	GetTimeout() time.Duration
}