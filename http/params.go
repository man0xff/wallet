package http

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"

	"example/core"
)

var errParam = errors.New("incorrect parameter")

var nameRegexp = regexp.MustCompile(`[a-zA-Z0-9_]{1,}`)

func (s *Server) getName(c echo.Context, key string) (string, error) {
	name := c.FormValue(key)
	if !nameRegexp.MatchString(name) {
		return "", fmt.Errorf("%w (key:'%s')", errParam, key)
	}
	return name, nil
}

func (s *Server) getUint(c echo.Context, key string) (uint64, error) {
	str := c.FormValue(key)
	val, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("%w (key:'%s')", errParam, key)
	}
	return val, nil
}

func (s *Server) getTime(c echo.Context, key string, defaultx time.Time) (time.Time, error) {
	str := c.FormValue(key)
	if str == "" {
		return defaultx, nil
	}
	val, err := time.Parse("2006-01-02T15:04:05", str)
	if err != nil {
		return time.Time{}, fmt.Errorf("%w (key:'%s')", errParam, key)
	}
	return val, nil
}

func (s *Server) getDirection(c echo.Context, key string) (core.Direction, error) {
	switch c.FormValue(key) {
	case "":
		return core.DirAny, nil
	case "in":
		return core.DirIn, nil
	case "out":
		return core.DirOut, nil
	}
	return core.Direction(0), fmt.Errorf("%w (key:'%s')", errParam, key)
}
