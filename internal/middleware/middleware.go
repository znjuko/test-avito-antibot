package middleware

import (
	"github.com/labstack/echo"
	"github.com/satori/go.uuid"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type MiddlewareHandler struct {
	logger *zap.SugaredLogger
}

func NewMiddlewareHandler(logger *zap.SugaredLogger) MiddlewareHandler {
	return MiddlewareHandler{logger: logger}
}

func (mh MiddlewareHandler) SetMiddleware(server *echo.Echo) {

	logFunc := mh.AccessLog()
	server.Use(mh.PanicMiddleWare)

	server.Use(logFunc)
}

func (mh MiddlewareHandler) PanicMiddleWare(next echo.HandlerFunc) echo.HandlerFunc {

	return func(c echo.Context) error {

		defer func() error {
			if err := recover(); err != nil {
				rId := c.Get("REQUEST_ID").(string)
				mh.logger.Info(
					zap.String("ID", rId),
					zap.String("ERROR", err.(error).Error()),
					zap.Int("ANSWER STATUS", http.StatusInternalServerError),
				)

				return c.JSON(http.StatusInternalServerError, `{panic}`)
			}
			return nil
		}()

		return next(c)
	}
}

func (mh MiddlewareHandler) AccessLog() echo.MiddlewareFunc {

	return func(next echo.HandlerFunc) echo.HandlerFunc {

		return func(rwContext echo.Context) error {

			uniqueID := uuid.NewV4()
			start := time.Now()
			rwContext.Set("REQUEST_ID", uniqueID.String())

			mh.logger.Info(
				zap.String("ID", uniqueID.String()),
				zap.String("URL", rwContext.Request().URL.Path),
				zap.String("METHOD", rwContext.Request().Method),
			)

			err := next(rwContext)

			respTime := time.Since(start)
			mh.logger.Info(
				zap.String("ID", uniqueID.String()),
				zap.Duration("TIME FOR ANSWER", respTime),
			)

			return err

		}
	}
}
