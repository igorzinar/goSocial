package mailer

import (
	"bytes"
	"fmt"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"html/template"
	"time"
)

type SendGridMailer struct {
	fromEmail string
	apiKey    string
	client    *sendgrid.Client
}

func NewSendGridMailer(apiKey, fromEmail string) *SendGridMailer {
	client := sendgrid.NewSendClient(apiKey)

	return &SendGridMailer{
		fromEmail: fromEmail,
		apiKey:    apiKey,
		client:    client,
	}
}

func (m *SendGridMailer) Send(templateFile, username, email string, data any, isSandbox bool) error {

	from := mail.NewEmail(FromName, m.fromEmail)
	to := mail.NewEmail(username, email)
	// parse template
	tmpl, err := template.ParseFS(FS, "templates/"+templateFile)
	if err != nil {
		return err
	}

	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return err
	}
	body := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(body, "body", data)
	if err != nil {
		return err
	}
	message := mail.NewSingleEmail(from, subject.String(), to, "", body.String())

	message.SetMailSettings(&mail.MailSettings{
		SandboxMode: &mail.Setting{
			Enable: &isSandbox,
		},
	})

	for i := 0; i < maxRetry; i++ {
		response, err := m.client.Send(message)
		if err != nil {
			fmt.Printf("Failed to send email %v, attempt %d of %d", email, i+1, maxRetry)
			fmt.Printf("Error : %v", err.Error())

			// exponential backoff
			time.Sleep(time.Second * time.Duration(i+1))
			continue
		}
		fmt.Printf("Email sent with status code %v", response.StatusCode)
		return nil
	}

	return fmt.Errorf("failed to send email after %v attempts", maxRetry)
}
