package vrc_osc

const DefaultHost = "127.0.0.1"
const DefaultPort = 9000

type Info struct {
	Host string
	Port int
}

func (i *Info) Complete() bool {
	return i.Host != "" && i.Port != 0
}
