package metasploit

type sessionListReq struct {
	_msgpack struct{} `msgpack:",asArray"`
	Method   string
	Token    string
}

type sessionListRes struct {
	Id          uint32 `msgpack:",omitempty"`
	Type        string `msgpack:",omitempty"`
	TunnelLocal string `msgpack:",omitempty"`
	TunnelPeer  string `msgpack:",omitempty"`
	ViaExploit  string `msgpack:",omitempty"`
	ViaPayload  string `msgpack:",omitempty"`
	Description string `msgpack:",omitempty"`
	Info        string `msgpack:",omitempty"`
	Workplace   string `msgpack:",omitempty"`
	SessionHost string `msgpack:",omitempty"`
	SessionPort int    `msgpack:",omitempty"`
	Username    string `msgpack:",omitempty"`
	UUID        string `msgpack:",omitempty"`
	ExploitUUID string `msgpack:",omitempty"`
}
