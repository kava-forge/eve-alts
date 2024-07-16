package esi

import (
	"net/http"
	stdhttp "net/http"
	"sync"
	"time"

	"github.com/kava-forge/eve-alts/lib/logging"
	"github.com/kava-forge/eve-alts/lib/logging/level"
)

const (
	CodeKey  = "code"
	StateKey = "state"
)

type CallbackServer struct {
	*stdhttp.Server

	logger        logging.Logger
	stateChannels *sync.Map
}

func (s *CallbackServer) handler(w stdhttp.ResponseWriter, r *stdhttp.Request) {
	q := r.URL.Query()
	code := q.Get(CodeKey)
	state := q.Get(StateKey)

	if code == "" {
		level.Error(s.logger).Message("could not retrieve response - empty code")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("error"))
		return
	}

	target, ok := s.stateChannels.Load(state)
	if !ok {
		level.Error(s.logger).Message("could not retrieve response - unexpected state")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("error"))
		return
	}

	targetChan, ok := target.(chan<- CodeState)
	if !ok {
		level.Error(s.logger).Message("could not retrieve response - unexpected state")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("error"))
		return
	}

	targetChan <- CodeState{
		Code:  code,
		State: state,
		Valid: true,
	}
	_, _ = w.Write([]byte("ok"))
}

func (s *CallbackServer) Expect(state string, target chan<- CodeState) {
	s.stateChannels.Store(state, target)
}

func (s *CallbackServer) Remove(state string) {
	s.stateChannels.Delete(state)
}

type CodeState struct {
	Code  string
	State string
	Valid bool
}

func NewCallbackServer(logger logging.Logger, serveAddr, callbackPath string) *CallbackServer {
	csrv := &CallbackServer{
		stateChannels: &sync.Map{},
	}

	mux := stdhttp.NewServeMux()
	mux.HandleFunc(callbackPath, csrv.handler)

	srv := &stdhttp.Server{
		Addr:         serveAddr,
		Handler:      mux,
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 2 * time.Second,
	}

	csrv.Server = srv

	return csrv
}
