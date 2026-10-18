package vrc_osc

const HRtoVRChat = "hr-to-vrchat"
const ZF = "zf"
const DefaultPreset = HRtoVRChat

var PresetHRtoVRChat = Params{
	Min: 0,
	Max: 150,

	StatusConnected:  "isHRConnected",
	StatusActive:     "isHRActive",
	StatusBeat:       "isHRBeat",
	ValueAbsolute:    "HR",
	ValuePercent:     "HRPercent",
	ValuePercentFull: "FullHRPercent",
}

var PresetDefault = &PresetHRtoVRChat

var PresetZF = Params{
	Min: 32,
	Max: 192,

	StatusConnected: "hr/connected",
	StatusActive:    "hr/active",
	StatusBeat:      "hr/beat",
	ValueAbsolute:   "hr/value",
	ValuePercent:    "hr/percent",
}

var Presets = map[string]Params{
	HRtoVRChat: PresetHRtoVRChat,
	ZF:         PresetZF,
}
