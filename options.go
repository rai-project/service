package service

import (
	"context"
	"time"

	registrystore "github.com/rai-project/libkv/store"
	"github.com/rai-project/registry"
	"github.com/rai-project/tracer"
	"github.com/rai-project/tracer/zipkin"
	"google.golang.org/grpc"
)

type Options struct {
	Name        string
	Description grpc.ServiceDesc

	// Register loop interval
	RegisterInterval time.Duration

	// Tracer
	Tracer tracer.Tracer

	// Registry
	Registry registrystore.Store

	// Before and After funcs
	BeforeStart []func() error
	BeforeStop  []func() error
	AfterStart  []func() error
	AfterStop   []func() error

	Context context.Context
}

type Option func(*Options)

var (
	DefaultName = "rai-project/service[default]"
)

func NewOptions(opts ...Option) *Options {
	rgs, err := registry.New()
	if err != nil {
		rgs = nil
	}
	trcr, err := zipkin.New(DefaultName)
	if err != nil {
		trcr = nil
	}
	options := &Options{
		Name:     DefaultName,
		Registry: rgs,
		Tracer:   trcr,
		Context:  context.Background(),
	}
	for _, o := range opts {
		o(options)
	}
	return options
}

// overrides the default option being used
// Must be at the beginning of the options
func Using(e *Options) Option {
	return func(o *Options) {
		*o = *e
	}
}

func Name(s string) Option {
	return func(o *Options) {
		o.Name = s
	}
}

func ServiceDescription(s grpc.ServiceDesc) Option {
	return func(o *Options) {
		o.Description = s
		o.Name = o.Description.ServiceName
	}
}

func Registry(s registrystore.Store) Option {
	return func(o *Options) {
		if o.Registry != nil {
			o.Registry.Close()
		}
		o.Registry = s
	}
}

func Tracer(s tracer.Tracer) Option {
	return func(o *Options) {
		if o.Tracer != nil {
			o.Tracer.Close()
		}
		o.Tracer = s
	}
}

func Context(ctx context.Context) Option {
	return func(o *Options) {
		o.Context = ctx
	}
}

func BeforeStart(fn func() error) Option {
	return func(o *Options) {
		o.BeforeStart = append(o.BeforeStart, fn)
	}
}

func BeforeStop(fn func() error) Option {
	return func(o *Options) {
		o.BeforeStop = append(o.BeforeStop, fn)
	}
}

func AfterStart(fn func() error) Option {
	return func(o *Options) {
		o.AfterStart = append(o.AfterStart, fn)
	}
}

func AfterStop(fn func() error) Option {
	return func(o *Options) {
		o.AfterStop = append(o.AfterStop, fn)
	}
}
