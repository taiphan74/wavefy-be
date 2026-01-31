package mail

import (
	"bytes"
	"embed"
	"html/template"
)

//go:embed templates/reset_password.html
var templatesFS embed.FS

var resetPasswordTemplate = template.Must(
	template.ParseFS(templatesFS, "templates/reset_password.html"),
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
