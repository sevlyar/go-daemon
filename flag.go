package daemon

import (
	"os"
)

// Not implemented !!!
func AddCommand(f Flag, sig os.Signal, handler SignalHandlerFunc) {
	if f != nil {
		AddFlag(f, sig)
	}
	if handler != nil {
		SetSigHandler(handler, sig)
	}
}

type Flag interface {
	IsSet() bool
}

func BoolFlag(f *bool) Flag {
	return &boolFlag{f}
}

func StringFlag(f *string, v string) Flag {
	return &stringFlag{f, v}
}

type boolFlag struct {
	b *bool
}

func (f *boolFlag) IsSet() bool {
	if f == nil {
		return false
	}
	return *f.b
}

type stringFlag struct {
	s *string
	v string
}

func (f *stringFlag) IsSet() bool {
	if f == nil {
		return false
	}
	return *f.s == f.v
}

var flags = make(map[Flag]os.Signal)

func AddFlag(f Flag, sig os.Signal) {
	flags[f] = sig
}

func SendCommands(p *os.Process) (err error) {
	for _, sig := range signals() {
		if err = p.Signal(sig); err != nil {
			return
		}
	}
	return
}

func signals() (ret []os.Signal) {
	ret = make([]os.Signal, 0, 1)
	for f, sig := range flags {
		if f.IsSet() {
			ret = append(ret, sig)
		}
	}
	return
}
