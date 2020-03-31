package constant

// define common use word for middleware and handler
const (
	// for xorm engine
	Db = "db"
	// for xorm session
	Session = "session"
	/*
	for response
	 */
	StatusCode = "statusCode"
	Error      = "error"
	Output     = "output"
	// for patch or put method
	Update     = "update"
	// for transfer userId in middleware and handler
	UserId 	   = "userId"
)
