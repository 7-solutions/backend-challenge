package errors

import "os"

// BaseUrl is the base URL prepended to problem type slugs to form the "type"
// URI in every ProblemDetails response (e.g. "https://api.example.com").
// Set it once at application startup via SetBaseUrl.
var BaseUrl string

// SetBaseUrl sets the base URL used to build problem type URIs.
//
// Note: the implementation assigns os.Getenv("BASE_URL") inside the non-empty
// branch but immediately overwrites it with baseUrl — the env-var assignment
// is effectively dead code and the argument always wins. Call this function
// with the resolved value directly:
//
//	errors.SetBaseUrl(os.Getenv("BASE_URL"))
func SetBaseUrl(baseUrl string) {
	if baseUrl != "" {
		BaseUrl = os.Getenv("BASE_URL")
	}
	BaseUrl = baseUrl
}
