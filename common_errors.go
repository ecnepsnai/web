package web

// CommonErrors are common errors types suitable for API endpoints
var CommonErrors = struct {
	NotFound        *Error
	BadRequest      *Error
	Unauthorized    *Error
	Forbidden       *Error
	ServerError     *Error
	TooManyRequests *Error
}{
	NotFound: &Error{
		Code:    404,
		Message: "Not Found",
	},
	BadRequest: &Error{
		Code:    400,
		Message: "Bad Request",
	},
	Unauthorized: &Error{
		Code:    403,
		Message: "Unauthorized",
	},
	Forbidden: &Error{
		Code:    403,
		Message: "Forbidden",
	},
	ServerError: &Error{
		Code:    500,
		Message: "Server Error",
	},
	TooManyRequests: &Error{
		Code:    429,
		Message: "Too Many Requests",
	},
}
