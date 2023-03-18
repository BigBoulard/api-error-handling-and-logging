package log

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/BigBoulard/api-error-handling-and-logging/src/rest_errors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

var l *logger = NewLogger()

type logger struct {
	logger zerolog.Logger
	// see https://github.com/rs/zerolog#leveled-logging
}

func NewLogger() *logger {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	var zlog zerolog.Logger
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
		Timestamp().
		Logger()

	return &logger{
		logger: zlog,
	}
}

func Error(err error, pkg string, funct string, msg string) {
	var restErr rest_errors.RestErr
	if errors.As(err, &restErr) {
		l.logger.
			Error().
			Stack().
			Err(restErr).
			Str("code", strconv.Itoa(restErr.Status())).
			Str("function", funct).
			Str("pkg", pkg).
			Msg(msg)
	} else {
		l.logger.
			Error().
			Stack().
			Err(err).
			Str("function", funct).
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
