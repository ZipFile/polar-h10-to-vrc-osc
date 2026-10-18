package vrc_osc

import (
	"github.com/ZipFile/polar-h10-to-osc/internal/common"
)

type Params struct {
	Min, Max         int
	StatusConnected  string `json:"status_connected,omitempty"`
	StatusActive     string `json:"status_active,omitempty"`
	StatusBeat       string `json:"status_beat,omitempty"`
	ValueAbsolute    string `json:"value_absolute,omitempty"`
	ValuePercent     string `json:"value_percent,omitempty"`
	ValuePercentFull string `json:"value_percent_full,omitempty"`
}

func (p *Params) Percent(hr int) float64 {
	return common.Percent(hr, p.Min, p.Max)
}

func (p *Params) PercentFull(hr int) float64 {
	return common.PercentFull(hr, p.Min, p.Max)
}
