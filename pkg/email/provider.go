package email

type AddEmailInput struct {
	Email     string
	ListID    string
}

type Provider interface {
	AddEmailToList(AddEmailInput) error
}
