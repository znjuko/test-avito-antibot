package usecase

import (
	"main/internal/ip"
	"main/internal/tools"
	"net"
	"time"
)

type IpUseCase struct {
	ipRepo  ip.IpRepoInter
	mask    string
	limit   int
	timeout int
}

func NewIpUseCase(rep ip.IpRepoInter, msk string, limit int, outer int) IpUseCase {
	return IpUseCase{ipRepo: rep, mask: msk, limit: limit, timeout: outer}
}

func (Ip IpUseCase) RegisterIp(ip string, reqTime time.Time) (int, error) {

	_, netIp, err := net.ParseCIDR(ip + "/" + Ip.mask)

	if err != nil {
		return 0, err
	}

	err = Ip.ipRepo.CreateFirstIpMask(netIp.String())

	if err == nil {
		return 0, nil
	}

	ipCounter, dbTime, err := Ip.ipRepo.GetMaskData(netIp.String())

	if err != nil {
		return 0, err
	}

	if dbTime != "" {
		timer, err := time.Parse(time.RFC3339, dbTime)

		if err != nil {
			return 0, err
		}

		if timer.String() > reqTime.String() {
			return int(timer.Sub(reqTime).Seconds()), tools.IpCoolDown
		}

		err = Ip.ipRepo.SetMaskDataToDefault(netIp.String())

		return 0, err
	}

	if ipCounter+1 > Ip.limit {

		err = Ip.ipRepo.SetMaskCoolDown(netIp.String(), Ip.limit, reqTime.Add(time.Duration(Ip.timeout)*time.Minute))

		if err != nil {
			return 0, err
		}

		return Ip.timeout * 60, tools.IpCoolDown
	}

	return 0, Ip.ipRepo.UpdateMaskData(netIp.String())
}

func (Ip IpUseCase) ResetIpCoolDown(ip string) error {
	_, netIp, err := net.ParseCIDR(ip + "/" + Ip.mask)

	if err != nil {
		return err
	}

	return Ip.ipRepo.ResetIpCoolDown(netIp.String())
}
