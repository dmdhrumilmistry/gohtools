package metasploit

type loginReq struct {
	_msgpack struct{} `msgpack:",asArray"`
	Method   string
	Username string
	Password string
}

type loginRes struct {
	Result       string `msgpack:"result"`
	Token        string `msgpack:"token"`
	Error        bool   `msgpack:"error"`
	ErrorClass   string `msgpack:"error_class"`
	ErrorMessage string `msgpack:"error_message"`
}

type logoutReq struct {
	_msgpack    struct{} `msgpack:",asArray"`
	Method      string
	Token       string
	LogoutToken string
}

type logoutRes struct {
	Result string `msgpack:"result"`
}

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
