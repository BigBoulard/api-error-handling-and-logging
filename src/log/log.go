package log

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/BigBoulard/api-error-handling-and-logging/src/rest_errors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

const API = "gw"

var l *logger = NewLogger()

type logger struct {
	logger zerolog.Logger
	// see https://github.com/rs/zerolog#leveled-logging
}

func NewLogger() *logger {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	var zlog zerolog.Logger
	if os.Getenv("APP_ENV") == "dev" {
		zlog = zerolog.New(
			zerolog.ConsoleWriter{
				Out:        os.Stderr,
				TimeFormat: time.RFC3339,
				FormatLevel: func(i interface{}) string {
					return strings.ToUpper(fmt.Sprintf("[%s]", i))
				},
				FormatMessage: func(i interface{}) string {
					return fmt.Sprintf("| %s |", i)
				},
			}).
			Level(zerolog.TraceLevel).
			With().
			Str("api", API).
			Timestamp().
			Logger()

	} else { // prod
		zlog = zerolog.New(os.Stdout).
			Level(zerolog.InfoLevel).
			With().
			Timestamp().
			Logger()
	}
	return &logger{
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
			Err(fmt.Errorf((fmt.Sprintf("%d - %s - %s", restErr.Status(), restErr.Causes(), restErr.Message())))).
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

func Error(err error, pkg string, method string, msg string) {
	var restErr rest_errors.RestErr
	if errors.As(err, &restErr) {
		l.logger.
			Error().
			Stack().
			Err(fmt.Errorf((fmt.Sprintf("%d - %s", restErr.Status(), restErr.Causes())))).
			Str("method", method).
			Str("pkg", pkg).
			Msg(msg)
	} else {
		l.logger.
			Error().
			Stack().
			Err(err).
			Str("method", method).
			Str("pkg", pkg).
			Msg(msg)
	}
}

func Fatal(err error, pkg string, method string, msg string) {
	var restErr rest_errors.RestErr

	if errors.As(err, &restErr) {
		l.logger.
			Fatal().
			Stack().
			Err(fmt.Errorf((fmt.Sprintf("%d - %s - %s", restErr.Status(), restErr.Causes(), restErr.Message())))).
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
		c.Next()
		// after request
		latency := time.Since(t)
		msg := c.Errors.String()
		if msg == "" {
			msg = "Request"
		}

		switch status := c.Writer.Status(); {
		case status >= 400:
			l.logger.Error().
				Str("method", c.Request.Method).
				Str("host", c.Request.Host).
				Str("url", c.Request.URL.Path).
				Int("status", c.Writer.Status()).
				Dur("lat", latency).
				Msg(msg)
			break
		case status >= 300:
			l.logger.Warn().
				Str("method", c.Request.Method).
				Str("host", c.Request.Host).
				Str("url", c.Request.URL.Path).
				Int("status", c.Writer.Status()).
				Dur("lat", latency).
				Msg(msg)
			break
		default:
			l.logger.Info().
				Str("method", c.Request.Method).
				Str("host", c.Request.Host).
				Str("url", c.Request.URL.Path).
				Int("status", c.Writer.Status()).
				Dur("lat", latency).
				Msg(msg)
		}
	}
}
