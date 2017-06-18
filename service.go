package service

import (
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Service interface {
	Start() error
	Stop() error
	Run() error
}

type service struct {
	opts *Options
}

func New(opts ...Option) (Service, error) {
	options := NewOptions(opts...)
	return &service{
		opts: options,
	}, nil
}

func (s *service) Start() error {
	for _, fn := range s.opts.BeforeStart {
		if err := fn(); err != nil {
			return err
		}
	}

	// ...

	for _, fn := range s.opts.AfterStart {
		if err := fn(); err != nil {
			return err
		}
	}
	return nil
}

func (s *service) Stop() error {

	var gerr error

	for _, fn := range s.opts.BeforeStop {
		if err := fn(); err != nil {
			gerr = err
		}
	}

	// ...

	if s.opts.Registry != nil {
		s.opts.Registry.Close()
	}

	if s.opts.Tracer != nil {
		s.opts.Tracer.Close()
	}

	for _, fn := range s.opts.AfterStop {
		if err := fn(); err != nil {
			gerr = err
		}
	}
	return gerr
}

func (s *service) Run() error {
	if err := s.Start(); err != nil {
		return err
	}

	// start reg loop
	ex := make(chan bool)
	go s.run(ex)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)

	select {
	// wait on kill signal
	case <-ch:
	// wait on context cancel
	case <-s.opts.Context.Done():
	}

	// exit reg loop
	close(ex)

	if err := s.Stop(); err != nil {
		return err
	}

	return nil
}

func (s *service) register() {

}

func (s *service) run(exit chan bool) {
	if s.opts.RegisterInterval <= time.Duration(0) {
		return
	}

	t := time.NewTicker(s.opts.RegisterInterval)

	for {
		select {
		case <-t.C:
			s.register()
		case <-exit:
			t.Stop()
			return
		}
	}
}

func (s *service) String() string {
	return s.opts.Name
}

func (s *service) GoString() string {
	return s.String()
}
