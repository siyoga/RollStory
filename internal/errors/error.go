package errors

type Error struct {
	Code    int64  `json:"code"`
	Reason  string `json:"reason"`
	Details error  `json:"details"`
}

func WD(err *Error, details error) *Error {
	e := *err
	e.Details = details
	return &e
}

func DatabaseError(details error) *Error {
	return WD(&Error{Code: 500, Reason: "database failed"}, details)
}

func AdapterError(details error) *Error {
	return WD(&Error{Code: 500, Reason: "adapter failed"}, details)
}
