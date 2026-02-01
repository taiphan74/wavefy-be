package mail

import (
	"bytes"
	"embed"
	"html/template"
)

//go:embed templates/reset_password.html templates/verify_email.html
var templatesFS embed.FS

var resetPasswordTemplate = template.Must(
	template.ParseFS(templatesFS, "templates/reset_password.html"),
)

var verifyEmailTemplate = template.Must(
	template.ParseFS(templatesFS, "templates/verify_email.html"),
)

type resetPasswordData struct {
	ResetURL string
}

func RenderResetPasswordHTML(resetURL string) (string, error) {
	var buf bytes.Buffer
	if err := resetPasswordTemplate.Execute(&buf, resetPasswordData{ResetURL: resetURL}); err != nil {
		return "", err
	}
	return buf.String(), nil
}

type verifyEmailData struct {
	VerifyURL string
}

func RenderVerifyEmailHTML(verifyURL string) (string, error) {
	var buf bytes.Buffer
	if err := verifyEmailTemplate.Execute(&buf, verifyEmailData{VerifyURL: verifyURL}); err != nil {
		return "", err
	}
	return buf.String(), nil
}
