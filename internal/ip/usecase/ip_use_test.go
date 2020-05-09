package usecase

import (
	"errors"
	"github.com/golang/mock/gomock"
	"main/internal/ip/usecase/mock"
	_errors "main/internal/tools"
	"testing"
	"time"
)

func TestIpUseCase_ResetIpCoolDown(t *testing.T) {
	ctrl := gomock.NewController(t)
	ipRepoMock := mock_ip.NewMockIpRepoInter(ctrl)
	ipUse := NewIpUseCase(ipRepoMock, "24", 2, 2)
	wrongIp := "123123123123"

	if err := ipUse.ResetIpCoolDown(wrongIp); err == nil {
		t.Error(" not expected behaviour at converting ip")
	}

	goodIp := "123.45.67.89"
	netIp := "123.45.67.0/24"

	ipRepoMock.EXPECT().ResetIpCoolDown(netIp).Return(nil)

	if err := ipUse.ResetIpCoolDown(goodIp); err != nil {
		t.Error(" not expected behaviour at cooldown resetting of ip", err)
	}
}

func TestIpUseCase_RegisterIp(t *testing.T) {
	ctrl := gomock.NewController(t)
	ipRepoMock := mock_ip.NewMockIpRepoInter(ctrl)
	ipUse := NewIpUseCase(ipRepoMock, "24", 2, 2)
	customErr := errors.New("msth happend")
	_time := time.Now()

	wrongIp := "123123123123"
	if val, err := ipUse.RegisterIp(wrongIp, _time); err == nil || val != 0 {
		t.Error(" not expected behaviour at converting ip", err, val)
	}

	goodIp := "123.45.67.89"
	netIp := "123.45.67.0/24"

	ipRepoMock.EXPECT().CreateFirstIpMask(netIp).Return(nil)

	if val, err := ipUse.RegisterIp(goodIp, _time); err != nil || val != 0 {
		t.Error(" not expected behaviour at CreateFirstIpMask", err, val)
	}

	ipRepoMock.EXPECT().CreateFirstIpMask(netIp).Return(customErr)
	ipRepoMock.EXPECT().GetMaskData(netIp).Return(0, "", customErr)
	if val, err := ipUse.RegisterIp(goodIp, _time); err == nil || val != 0 {
		t.Error(" not expected behaviour at GetMaskData", err, val)
	}

	ipRepoMock.EXPECT().CreateFirstIpMask(netIp).Return(customErr)
	ipRepoMock.EXPECT().GetMaskData(netIp).Return(0, "", nil)
	ipRepoMock.EXPECT().UpdateMaskData(netIp).Return(customErr)
	if val, err := ipUse.RegisterIp(goodIp, _time); err != customErr || val != 0 {
		t.Error(" not expected behaviour at UpdateMaskData", err, val)
	}

	ipRepoMock.EXPECT().CreateFirstIpMask(netIp).Return(customErr)
	ipRepoMock.EXPECT().GetMaskData(netIp).Return(2, "", nil)
	ipRepoMock.EXPECT().SetMaskCoolDown(netIp, gomock.Any(), gomock.Any()).Return(nil)
	if val, err := ipUse.RegisterIp(goodIp, _time); err != _errors.IpCoolDown || val != 120 {
		t.Error(" not expected behaviour at SetMaskCoolDown", err, val)
	}

	ipRepoMock.EXPECT().CreateFirstIpMask(netIp).Return(customErr)
	ipRepoMock.EXPECT().GetMaskData(netIp).Return(2, "", nil)
	ipRepoMock.EXPECT().SetMaskCoolDown(netIp, gomock.Any(), gomock.Any()).Return(customErr)
	if val, err := ipUse.RegisterIp(goodIp, _time); err != customErr || val != 0 {
		t.Error(" not expected behaviour at SetMaskCoolDown", err, val)
	}

	formatTime := _time.Add(2000000000).Format(time.RFC3339)
	ipRepoMock.EXPECT().CreateFirstIpMask(netIp).Return(customErr)
	ipRepoMock.EXPECT().GetMaskData(netIp).Return(2, formatTime, nil)
	if val, err := ipUse.RegisterIp(goodIp, _time); err == nil || val != 1 {
		t.Error(" not expected behaviour at timeout ", err, val)
	}

	formatTime = _time.Add(2).Format(time.RFC3339)
	ipRepoMock.EXPECT().CreateFirstIpMask(netIp).Return(customErr)
	ipRepoMock.EXPECT().GetMaskData(netIp).Return(2, formatTime, nil)
	ipRepoMock.EXPECT().SetMaskDataToDefault(netIp).Return(customErr)
	if val, err := ipUse.RegisterIp(goodIp, _time); err != customErr || val != 0 {
		t.Error(" not expected behaviour at timeout ", err, val)
	}

}
