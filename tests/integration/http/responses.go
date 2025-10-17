package http_test

var (
	ValidationError_ProductNameMissing = `{
		"code": 400,
		"err": "bad_request",
		"message": "Invalid field values",
		"causes": [{"field":"ProductName","message":"ProductName is a required field"}]
	}`

	InvalidAuctionError = `{
		"code": 400,
		"err": "bad_request",
		"message": "invalid auction object",
		"causes": null
	}`

	ValidationError_DescriptionTooShort = `{
		"code": 400,
		"err": "bad_request",
		"message": "Invalid field values",
		"causes": [{"field":"Description","message":"Description must be at least 10 characters in length"}]
	}`

	InvalidTypeError = `{
		"code": 404,
		"err": "not_found",
		"message": "Invalid type error",
		"causes": null
	}`

	MalformedJSONError = `{
		"code": 404,
		"err": "not_found",
		"message": "Invalid type error",
		"causes": null
	}`

	ValidationError_InvalidCondition = `{
		"code": 400,
		"err": "bad_request",
		"message": "Invalid field values",
		"causes": [{"field":"Condition","message":"Condition must be one of [1 2 3]"}]
	}`

	InsertOneError = `{
		"code": 500,
		"err": "internal_server",
		"message": "Error trying to insert auction",
		"causes": null
	}`
)
