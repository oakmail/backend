package errors

var (
	DatabaseError                = &Error{Code: 1, Message: "Database error."}
	DatabaseInconsistency        = &Error{Code: 2, Message: "Database inconsistency."}
	InvalidAuthorizationHeader   = &Error{Code: 3, Message: "Invalid Authorization header."}
	InvalidJSONInput             = &Error{Code: 4, Message: "Invalid JSON input."}
	TokenNotFound                = &Error{Code: 5, Message: "Token not found."}
	InvalidTokenTypeMustBeAuth   = &Error{Code: 6, Message: "Invalid token type, must be \"auth\"."}
	AuthenticationTokenExpired   = &Error{Code: 7, Message: "Authentication token expired."}
	AccountIsBlocked             = &Error{Code: 8, Message: "Account is blocked."}
	InsufficientTokenPermissions = &Error{Code: 9, Message: "Insufficient token permissions."}
	InvalidCallbackFormat        = &Error{Code: 10, Message: "Invalid callback format."}
	ApplicationNameIsInvalid     = &Error{Code: 11, Message: "Application name is invalid."}
	InvalidEmailFormat           = &Error{Code: 12, Message: "Invalid e-mail format."}
	InvalidHomePageFormat        = &Error{Code: 13, Message: "Invalid home page format."}
	InvalidIDFormat              = &Error{Code: 14, Message: "Invalid ID format."}
	ApplicationNotFound          = &Error{Code: 15, Message: "Application not found."}
	InvalidOwnerID               = &Error{Code: 16, Message: "Invalid owner ID."}
	InvalidDateCreatedStart      = &Error{Code: 17, Message: "Invalid date created start."}
	InvalidDateCreatedEnd        = &Error{Code: 18, Message: "Invalid date created end."}
	InvalidDateModifiedStart     = &Error{Code: 19, Message: "Invalid date modified start."}
	InvalidDateModifiedEnd       = &Error{Code: 20, Message: "Invalid date modified end."}
	InvalidSkipFormat            = &Error{Code: 21, Message: "Invalid skip format."}
	InvalidLimitFormat           = &Error{Code: 22, Message: "Invalid limit format."}
	OutdatedObjectUsedInPUT      = &Error{Code: 23, Message: "Outdated object used in PUT."}
	InvalidEmailDomain           = &Error{Code: 24, Message: "Invalid email domain."}
	AddressIsTaken               = &Error{Code: 25, Message: "Address is taken."}
	InvalidPasswordFormat        = &Error{Code: 26, Message: "Invalid password format."}
	AccountNotFound              = &Error{Code: 27, Message: "Account not found."}
	InvalidApplicationSecret     = &Error{Code: 28, Message: "Invalid application secret."}
	AuthorizationCodeNotFound    = &Error{Code: 29, Message: "Authorization code not found."}
	InvalidGrantType             = &Error{Code: 30, Message: "Invalid grant type."}
	InvalidExpiryDate            = &Error{Code: 31, Message: "Invalid expiry date."}
	InvalidPassword              = &Error{Code: 32, Message: "Invalid password."}
)
