package logger

import (
	"context"
	"fmt"
	"io"
	"log"
	"log/slog"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/green-ecolution/green-ecolution-backend/internal/utils"
)

type LogFormat string

const (
	JSON    LogFormat = "json"
	Text    LogFormat = "text"
	Console LogFormat = "console"
)

type LogLevel string

const (
	Debug LogLevel = "debug"
	Info  LogLevel = "info"
	Warn  LogLevel = "warn"
	Error LogLevel = "error"
)

func (l LogLevel) ToSLog() slog.Level {
	switch l {
	case Debug:
		return slog.LevelDebug
	case Info:
		return slog.LevelInfo
	case Warn:
		return slog.LevelWarn
	case Error:
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

type GroupOrAttrs struct {
	attr  slog.Attr
	group string
}

type PrettyHandlerOptions struct {
	SlogOpts slog.HandlerOptions
}

type PrettyHandler struct {
	opts PrettyHandlerOptions
	slog.Handler
	l   *log.Logger
	goa []GroupOrAttrs
}

func NewPrettyHandler(out io.Writer, opts PrettyHandlerOptions) *PrettyHandler {
	h := &PrettyHandler{
		Handler: slog.NewTextHandler(out, &opts.SlogOpts),
		l:       log.New(out, "", 0),
		opts:    opts,
		goa:     make([]GroupOrAttrs, 0),
	}

	return h
}

func (h *PrettyHandler) WithGroup(group string) slog.Handler {
	return &PrettyHandler{
		Handler: h.Handler.WithGroup(group),
		l:       &*h.l,
		opts:    h.opts,
		goa:     append(h.goa, GroupOrAttrs{group: group}),
	}
}

func (h *PrettyHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newAttrs := make([]GroupOrAttrs, len(attrs))

	for i, attr := range attrs {
		newAttrs[i] = GroupOrAttrs{attr: attr}
	}

	return &PrettyHandler{
		Handler: h.Handler.WithAttrs(attrs),
		l:       &*h.l,
		opts:    h.opts,
		goa:     append(h.goa, newAttrs...),
	}
}

func (h *PrettyHandler) Handle(ctx context.Context, r slog.Record) error {
	level := r.Level.String()
	switch r.Level {
	case slog.LevelDebug:
		level = fmt.Sprintf("\033[35;1m%s\033[0m", level) // Magenta
	case slog.LevelInfo:
		level = fmt.Sprintf("\033[34;1m%s\033[0m", level) // Blue
	case slog.LevelWarn:
		level = fmt.Sprintf("\033[33;1m%s\033[0m", level) // Yellow
	case slog.LevelError:
		level = fmt.Sprintf("\033[31;1m%s\033[0m", level) // Red
	}

	fields := make(map[string]interface{})
	lastGroup := ""
	for _, goa := range h.goa {
		if goa.group != "" {
			lastGroup += goa.group + "."
		} else {
			attr := goa.attr
			if lastGroup != "" {
				attr.Key = lastGroup + attr.Key
			}

			fields[attr.Key] = attr.Value.Any()
		}
	}

	r.Attrs(func(a slog.Attr) bool {
		a.Value = a.Value.Resolve()
		if lastGroup != "" {
			a.Key = lastGroup + a.Key
		}
		fields[a.Key] = a.Value.Any()
		return true
	})

	logFields := make([]string, 0, len(fields))
	for k, v := range fields {
		keyColorSeq := "\033[38;5;249;1m" // gray bold
		valueColorSeq := "\033[38;5;249m" // gray

		if k == "err" || k == "error" {
			keyColorSeq = "\033[31;1m" // red bold
			valueColorSeq = "\033[31m" // red
		}

		logFields = append(logFields, fmt.Sprintf("%s%s=\033[0m%s\"%v\"\033[0m", keyColorSeq, k, valueColorSeq, v))
	}

	var source string
	if h.opts.SlogOpts.AddSource {
		frame, _ := runtime.CallersFrames([]uintptr{r.PC}).Next()
		dir, file := filepath.Split(frame.File)
		rootDir := utils.RootDir()
		relativeDir, err := filepath.Rel(rootDir, filepath.Dir(dir))
		if err != nil {
			return err
		}

		source = fmt.Sprintf("\033[38;5;98;1;4m%s:%d\033[0m", filepath.Join(relativeDir, file), frame.Line)
	}

	timeStr := r.Time.Format(time.Stamp)
	msg := fmt.Sprintf("\033[36m%s\033[0m", r.Message) // Cyan

	h.l.Println(timeStr, level, source, msg, strings.Join(logFields, " "))

	return nil
}

func CreateLogger(out io.Writer, format LogFormat, level LogLevel) func() *slog.Logger {
	return func() *slog.Logger {
		options := &slog.HandlerOptions{
			AddSource: true,
			Level:     level.ToSLog(),
		}

		var handler slog.Handler
		switch format {
		case JSON:
			handler = slog.NewJSONHandler(out, options)
		case Text:
			handler = slog.NewTextHandler(out, options)
		case Console:
			handler = NewPrettyHandler(out, PrettyHandlerOptions{SlogOpts: *options})
		default:
			handler = slog.NewTextHandler(out, options)
		}

		log := slog.New(handler)
		return log
	}
}
