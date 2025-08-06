package api

type IssueRequest struct {
	UserData []byte
	Nonce    []byte
	Options  any
}

type ValidateRequest struct {
	Document []byte
	Nonce    []byte
	Options  any
}

type IssueRequestWrapper struct {
	Request  *IssueRequest
	Response chan *IssueResponse
}

type ValidateRequestWrapper struct {
	Request  *ValidateRequest
	Response chan *ValidateResponse
}
