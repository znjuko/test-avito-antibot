package delivery

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"main/internal/ip/delivery/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIpHandler_Reset(t *testing.T) {
	ctrl := gomock.NewController(t)
	ipUseMock := mock_ip.NewMockIpUseInter(ctrl)
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	prLogger, _ := config.Build()
	logger := prLogger.Sugar()
	defer prLogger.Sync()
	ipHandler := NewIpHandler(ipUseMock, logger)

	errs := []error{errors.New("smth happend"), nil}
	expectedBehaviour := []int{http.StatusConflict, http.StatusOK}

	for iter, _ := range expectedBehaviour {
		ip := uuid.NewV4()
		ipUseMock.EXPECT().ResetIpCoolDown(ip.String()).Return(errs[iter])

		e := echo.New()
		req := httptest.NewRequest(echo.POST, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/")
		c.Set("REQUEST_ID", "123")
		req.Header.Set("X-Forwarded-For", ip.String())

		if assert.NoError(t, ipHandler.Reset(c)) {
			assert.Equal(t, expectedBehaviour[iter], rec.Code)
		}
	}
}

func TestIpHandler_RegisterTime(t *testing.T) {
	ctrl := gomock.NewController(t)
	ipUseMock := mock_ip.NewMockIpUseInter(ctrl)
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	prLogger, _ := config.Build()
	logger := prLogger.Sugar()
	defer prLogger.Sync()
	ipHandler := NewIpHandler(ipUseMock, logger)

	errs := []error{errors.New("smth happend"), nil}
	expectedBehaviour := []int{http.StatusTooManyRequests, http.StatusOK}

	for iter, _ := range expectedBehaviour {
		ip := uuid.NewV4()
		ipUseMock.EXPECT().RegisterIp(ip.String(), gomock.Any()).Return(23, errs[iter])

		e := echo.New()
		req := httptest.NewRequest(echo.POST, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/")
		c.Set("REQUEST_ID", "123")
		req.Header.Set("X-Forwarded-For", ip.String())

		if assert.NoError(t, ipHandler.RegisterTime(c)) {
			assert.Equal(t, expectedBehaviour[iter], rec.Code)
		}
	}
}
