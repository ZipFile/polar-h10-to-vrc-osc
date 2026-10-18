package vrc_osc

import "github.com/hypebeast/go-osc/osc"

func fmt[V any](name string, value V) (msg *osc.Message) {
	if name == "" {
		return nil
	}

	return osc.NewMessage("/avatar/parameters/"+name, value)
}
