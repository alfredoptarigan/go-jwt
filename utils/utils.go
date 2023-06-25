package utils

func customErrorMessage(tag string) string {
	switch tag {
	case "required":
		return "The %s field is required"
	case "email":
		return "The %s field must be a valid email address"
	case "min":
		return "The %s field must be at least %s characters"
	case "max":
		return "The %s field must be at most %s characters"
	case "eqfield":
		return "The %s field must be equal to the %s field"
		// unique
	case "unique":
		return "The %s field must be unique"
	default:
		return "The %s field is invalid"
	}
}
