package api

type IssueResponse struct {
	Document []byte
	Error    error
}

type ValidateResponse struct {
	UserData []byte
	Valid    bool
	Error    error
}
