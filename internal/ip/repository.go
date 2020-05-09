package ip

import "time"

type IpRepoInter interface {
	ResetIpCoolDown(string) error
	CreateFirstIpMask(string) error
	GetMaskData(string) (int, string, error)
	SetMaskDataToDefault(string) error
	SetMaskCoolDown(string, int, time.Time) error
	UpdateMaskData(string) error
}
