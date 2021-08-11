package service

import (
	"fmt"
	"github.com/cookienyancloud/back/internal/config"
	emailProvider "github.com/cookienyancloud/back/pkg/email"
)

const (
	nameField            = "name"
	verificationLinkTmpl = "%s/verification?code=%s" // <frontend_url>/verification?code=<verification_code>

)

type EmailsService struct {
	provider    emailProvider.Provider
	sender      emailProvider.Sender
	config      config.EmailConfig
	frontendUrl string
}

// Structures used for templates.
type verificationEmailInput struct {
	VerificationLink string
}

type purchaseSuccessfulEmailInput struct {
	Name       string
	CourseName string
}

func NewEmailsService(
	provider emailProvider.Provider,
	sender emailProvider.Sender,
	config config.EmailConfig,
	frontendUrl string,
) *EmailsService {
	return &EmailsService{
		provider,
		sender,
		config,
		frontendUrl,
	}
}

func (s *EmailsService) SendUserVerificationEmail(input VerificationEmailInput) error {
	//s.config.Subjects.Verification = "./templates/verification_email.html"
	//subject := fmt.Sprintf(s.config.Subjects.Verification, input.Name)
	//templateInput := verificationEmailInput{s.createVerificationLink(input.VerificationCode)}
	//if err := sendInput.GenerateBodyFromHTML(s.config.Templates.Verification, templateInput); err != nil {
	//	return err
	//}
	subject := "Subject: Email Verification Code\r\n\r\n"
	body := s.createVerificationLink(input.VerificationCode)
	sendInput := emailProvider.SendEmailInput{Subject: subject, To: input.Email,Body: body}
	return s.sender.Send(sendInput)

}

func (s *EmailsService) createVerificationLink(code string) string {
	return fmt.Sprintf(verificationLinkTmpl, s.frontendUrl, code)
}

