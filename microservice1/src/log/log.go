package log

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/BigBoulard/api-error-handling-and-logging/microservice1/src/conf"
	"github.com/BigBoulard/api-error-handling-and-logging/microservice1/src/rest_errors"
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

const microservice = "microservice1"
const CorrelationHeader = "Correlation-ID"

var l *logger

type logger struct {
	logger zerolog.Logger
	// see https://github.com/rs/zerolog#leveled-logging
}

func NewLogger() {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	var zlog zerolog.Logger
	if conf.Env.AppMode == "dev" {
		zlog = zerolog.New(
			zerolog.ConsoleWriter{ // ConsoleWriter is too slow to be used in prod
				Out:        os.Stderr,
				TimeFormat: time.RFC3339,
				FormatMessage: func(i interface{}) string {
					return fmt.Sprintf("| %s |", i)
				},
				FieldsExclude: []string{
					"api",
					"user_agent",
					"host",
					"lat",
				},
			}).
			Level(zerolog.InfoLevel).
			With().
			Str("microservice", microservice).
			Timestamp().
			Logger()
	} else { // prod
		zlog = zerolog.New(os.Stdout).
			Level(zerolog.InfoLevel).
			With().
			Timestamp().
			Logger()
	}
	l = &logger{
		logger: zlog,
	}
}

func Info(pkg string, method string, msg string) {
	l.logger.
		Info().
		Str("method", method).
		Str("pkg", pkg).
		Msg(msg)
}

func Warn(err error, pkg string, method string, msg string) {
	var restErr rest_errors.RestErr

	if errors.As(err, &restErr) {
		l.logger.
			Warn().
			Stack().
			Err(fmt.Errorf((fmt.Sprintf("%d - %s - %s", restErr.Status(), restErr.Title(), restErr.Message())))).
			Str("method", method).
			Str("pkg", pkg).
			Msg(msg)
	} else {
		l.logger.
			Warn().
			Err(err).
			Str("method", method).
			Str("pkg", pkg).
			Msg(msg)
	}
}

func Error(restErr rest_errors.RestErr, msg string) {
	l.logger.Error().
		Err(fmt.Errorf((fmt.Sprintf("%s - %s", restErr.Code(), restErr.Message())))).
		Str("path", restErr.Path()).
		Msg(msg)
}

func Fatal(err error, pkg string, method string, msg string) {
	var restErr rest_errors.RestErr

	if errors.As(err, &restErr) {
		l.logger.
			Fatal().
			Stack().
			Err(fmt.Errorf((fmt.Sprintf("%d - %s - %s", restErr.Status(), restErr.Title(), restErr.Message())))).
			Str("method", method).
			Str("pkg", pkg).
			Msg(msg)
	} else {
		l.logger.
			Fatal().
			Stack().
			Err(err).
			Str("method", method).
			Str("pkg", pkg).
			Msg(msg)
	}
}

func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()
		// before request

		// created a correlation_id in the header if not present
		correlationID := c.Request.Header.Get(CorrelationHeader)
		if correlationID == "" {
			correlationID = xid.New().String()
			c.Request.Header.Add(CorrelationHeader, correlationID)
		}
		// add correlationID to the logger context
		l.logger.UpdateContext(func(c zerolog.Context) zerolog.Context {
			return c.Str("correlation_id", correlationID)
		})

		c.Next()
		// after request

		latency := time.Since(t)
		msg := c.Errors.String()
		if msg == "" {
			msg = "request"
		}

		// defer func makes the request be logged even if the handler panics
		defer func() {
			logSwitch(c, latency, msg)
		}()

	}
}

func logSwitch(c *gin.Context, latency time.Duration, msg string) {
	switch status := c.Writer.Status(); {
	case status >= 400:
		l.logger.Error().
			Str("correlation_id", c.Request.Header.Get(CorrelationHeader)).
			Str("method", c.Request.Method).
			Str("url", c.Request.URL.Path).
			Str("host", c.Request.Host).
			Str("user_agent", c.Request.UserAgent()).
			Int("status", c.Writer.Status()).
			Dur("lat", latency).
			Msg(msg)
		break
	case status >= 300:
		l.logger.Warn().
			Str("correlation_id", c.Request.Header.Get(CorrelationHeader)).
			Str("method", c.Request.Method).
			Str("url", c.Request.URL.Path).
			Str("host", c.Request.Host).
			Str("user_agent", c.Request.UserAgent()).
			Int("status", c.Writer.Status()).
			Dur("lat", latency).
			Msg(msg)
		break
	default:
		l.logger.Info().
			Str("correlation_id", c.Request.Header.Get(CorrelationHeader)).
			Str("method", c.Request.Method).
			Str("url", c.Request.URL.Path).
			Str("host", c.Request.Host).
			Str("user_agent", c.Request.UserAgent()).
			Int("status", c.Writer.Status()).
			Dur("lat", latency).
			Msg(msg)
	}
}
