package log

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/BigBoulard/api-error-handling-and-logging/microservice1/src/rest_errors"
	"github.com/BigBoulard/api-error-handling-and-logging/microservice2/src/conf"
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

const microservice = "microservice2"
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
			zerolog.ConsoleWriter{
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

func Error(err error, pkg string, method string, msg string) {
	var restErr rest_errors.RestErr
	if errors.As(err, &restErr) {
		l.logger.Error().
			Err(fmt.Errorf((fmt.Sprintf("%d - %s", restErr.Status(), restErr.Title())))).
			Str("method", method).
			Str("pkg", pkg).
			Msg(msg)
	} else {
		l.logger.Error().
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

		// add correlation_id to the request and logger context
		correlationID := c.Request.Header.Get(CorrelationHeader)
		if correlationID == "" {
			correlationID := xid.New().String()
			c.Request.Header.Add(CorrelationHeader, correlationID)
		}
		ctx := context.WithValue(c.Request.Context(), "correlation_id", correlationID)
		c.Request = c.Request.WithContext(ctx)

		// l.logger panics
		l.logger.UpdateContext(func(c zerolog.Context) zerolog.Context {
			return c.Str("correlation_id", correlationID)
		})

		c.Next()
		// after request is handled
		latency := time.Since(t)
		msg := c.Errors.String()
		if msg == "" {
			msg = "incoming request"
		}

		switch status := c.Writer.Status(); {
		case status >= 400:
			defer func() { // defer func makes the request be logged even if the handler panics
				l.logger.Error().
					Str("correlation_id", c.Request.Header.Get(CorrelationHeader)).
					Str("method", c.Request.Method).
					Str("url", c.Request.URL.Path).
					Str("host", c.Request.Host).
					Str("user_agent", c.Request.UserAgent()).
					Int("status", c.Writer.Status()).
					Dur("lat", latency).
					Msg(msg)
			}()
			break
		case status >= 300:
			defer func() {
				l.logger.Warn().
					Str("correlation_id", c.Request.Header.Get(CorrelationHeader)).
					Str("method", c.Request.Method).
					Str("url", c.Request.URL.Path).
					Str("host", c.Request.Host).
					Str("user_agent", c.Request.UserAgent()).
					Int("status", c.Writer.Status()).
					Dur("lat", latency).
					Msg(msg)
			}()
			break
		default:
			defer func() {
				l.logger.Info().
					Str("correlation_id", c.Request.Header.Get(CorrelationHeader)).
					Str("method", c.Request.Method).
					Str("url", c.Request.URL.Path).
					Str("host", c.Request.Host).
					Str("user_agent", c.Request.UserAgent()).
					Int("status", c.Writer.Status()).
					Dur("lat", latency).
					Msg(msg)
			}()
		}
	}
}
