package diag

import (
	"net/http"
	"strings"

	"github.com/prasannavl/mchain/hconv"

	"github.com/prasannavl/mchain"

	"github.com/prasannavl/go-gluons/log"
)

func LogLevelSwitcher(opts *LogSwitcherOpts) func(*http.ServeMux) {
	if opts == nil {
		o := DefaultLogSwitcherOpts()
		opts = &o
	}
	return func(mux *http.ServeMux) {
		mux.Handle(opts.Path, hconv.ToHttp(LogLevelSwitchHandlerFunc(opts), nil))
	}
}

type LogSwitcherOpts struct {
	Path       string
	LevelParam string
	FlushParam string
}

func DefaultLogSwitcherOpts() LogSwitcherOpts {
	return LogSwitcherOpts{
		Path:       "/log",
		LevelParam: "set-level",
		FlushParam: "flush",
	}
}

func LogLevelSwitchHandlerFunc(opts *LogSwitcherOpts) mchain.HandlerFunc {
	if opts == nil {
		o := DefaultLogSwitcherOpts()
		opts = &o
	}
	f := func(w http.ResponseWriter, r *http.Request) error {
		flush := r.URL.Query().Get(opts.FlushParam)
		lvl := r.URL.Query().Get(opts.LevelParam)

		if lvl != "" {
			lvl = strings.ToLower(strings.TrimSpace(lvl))
			level := log.LogLevelFromString(lvl)
			if log.IsValidLevel(level) {
				log.SetFilter(log.GetLogger(), log.LogFilterForLevel(level))
			} else {
				w.WriteHeader(http.StatusBadRequest)
				return nil
			}
		}

		if flush != "" {
			log.Flush()
		}

		w.Write([]byte("OK"))

		return nil
	}
	return f
}
