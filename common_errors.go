package web

// CommonErrors common errors
var CommonErrors = struct {
	NotFound     *Error
	BadRequest   *Error
	Forbidden    *Error
	ServerError  *Error
	Unauthorized *Error
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
}
