package errno

type Errno struct {
	State int    `json:"state"`
	Msg   string `json:"msg"`
}

func (e *Errno) Error() string {
	return e.Msg
}

func (e Errno) Add(s string) *Errno {
	e.Msg += ": " + s
	return &e
}

type ResponseErrno struct {
	State      int    `json:"state"`
	Msg        string `json:"msg"`
	HttpCode   int    `json:"http_code"`
	OriginBody string `json:"origin_body"`
}

func (e *ResponseErrno) Error() string {
	return e.Msg
}

func (e ResponseErrno) SetCode(code int, s string) *ResponseErrno {
	e.HttpCode = code
	e.OriginBody = s
	return &e
}

func (e ResponseErrno) Add(s string) *ResponseErrno {
	e.Msg += ": " + s
	return &e
}

var (
	NotSupportChainType = &Errno{10001, "Not support this chain"}
	InvalidTx           = &Errno{10002, "Invalid Tx"}
	InvalidTypeAssert   = &Errno{20001, "Invalid type asset"}
)
