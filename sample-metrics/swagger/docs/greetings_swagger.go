package docs

import "github.com/MaibornWolff/maDocK8s/exporter/sample-metrics/swagger/api"

// swagger:route POST /greetings greetings-tag idOfGreetingsEndpoint
// Greetings returns a greeting to the developer.
// responses:
//   description: GreetingResponse

// This text will appear as description of your response body.
// swagger:response GreetingResponse
type GreetingResponseWrapper struct {
	// in:body
	Body api.GreetingResponse
}

// swagger:parameters idOfGreetingsEndpoint
type GreetingParamsWrapper struct {
	// This text will appear as description of your request body.
	// in:body
	Body api.GreetingRequest
}
