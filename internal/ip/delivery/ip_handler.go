package delivery

import (
	"github.com/labstack/echo"
	"go.uber.org/zap"
	"main/internal/ip"
	"net/http"
	"strconv"
	"time"
)

type IpHandler struct {
	ipUse  ip.IpUseInter
	logger *zap.SugaredLogger
}

func NewIpHandler(use ip.IpUseInter, logs *zap.SugaredLogger) IpHandler {
	return IpHandler{ipUse: use, logger: logs}
}

func (Ip IpHandler) RegisterTime(rwContext echo.Context) error {

	rId := rwContext.Get("REQUEST_ID").(string)
	timer := time.Now()
	ip := rwContext.Request().Header.Get("X-Forwarded-For")

	retryAfter, err := Ip.ipUse.RegisterIp(ip, timer)

	if err != nil {
		Ip.logger.Debug(
			zap.String("ID", rId),
			zap.String("ERROR", err.Error()),
			zap.Int("ANSWER STATUS", http.StatusTooManyRequests),
		)

		rwContext.Response().Header().Set("Retry-After", strconv.Itoa(retryAfter))
		return rwContext.NoContent(http.StatusTooManyRequests)
	}

	Ip.logger.Info(
		zap.String("ID", rId),
		zap.Int("ANSWER STATUS", http.StatusCreated),
	)

	return rwContext.JSON(http.StatusCreated, `{"ip" :  "registered"}`)
}

func (Ip IpHandler) Reset(rwContext echo.Context) error {

	rId := rwContext.Get("REQUEST_ID").(string)
	ip := rwContext.Request().Header.Get("X-Forwarded-For")

	err := Ip.ipUse.ResetIpCoolDown(ip)

	if err != nil {
		Ip.logger.Debug(
			zap.String("ID", rId),
			zap.String("ERROR", err.Error()),
			zap.Int("ANSWER STATUS", http.StatusConflict),
		)

		return rwContext.NoContent(http.StatusConflict)
	}

	Ip.logger.Info(
		zap.String("ID", rId),
		zap.Int("ANSWER STATUS", http.StatusOK),
	)

	return rwContext.JSON(http.StatusOK, `{"ip" :  "reset"}`)
}

func (Ip IpHandler) InitHandlers(server *echo.Echo) {
	server.POST("/", Ip.RegisterTime)
	server.PUT("/reset", Ip.Reset)
}
