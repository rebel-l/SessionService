package session

const (
	BadRequestText          = "A body needs to be send as JSON and it needs at least a data field"
	InternalServerErrorText = "Session could not be stored because of internal error. Contact administrator or retry it later."
	SessionNotFoundText		= "Session was not found or has expired."
	UUIDLENGTH              = 36	// includes dashes: 8-4-4-4-12
)
