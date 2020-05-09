package tools

import "errors"

var (
	IpCoolDown = errors.New("ip got cooldown")
	WrongIpFormat = errors.New("wrong ip format")
)
