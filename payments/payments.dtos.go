package payments

type AuthorizationRequest struct {
	Type         string
	RequiresAuth bool
	URL          string
}
