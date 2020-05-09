package ip

import "time"

type IpUseInter interface {
	ResetIpCoolDown(string) error
	RegisterIp(string, time.Time) (int, error)
}
