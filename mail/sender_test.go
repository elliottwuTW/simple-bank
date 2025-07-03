package mail

import (
	"testing"

	"github.com/simple_bank/config"
	"github.com/stretchr/testify/require"
)

func TestSendEmailWithGmail(t *testing.T) {
	cfg, err := config.LoadConfig("..")
	require.NoError(t, err)

	sender := NewGmailSender(cfg.Email.Name, cfg.Email.Address, cfg.Email.Password)

	subject := "A test email"
	content := `
	<h1> Hello World </h1>
	<p>This is a test message.</p>
	`
	attachFiles := []string{"../Makefile"}
	to := []string{"elliott10009@gmail.com"}

	err = sender.SendEmail(subject, content, to, nil, nil, attachFiles)
	require.NoError(t, err)
}
