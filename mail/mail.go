package mail

import (
	"bytes"
	"context"
	"gopkg.in/gomail.v2"
	"io"
	"os"
	"os/exec"
	"text/template"
	"time"
)

const noReplyAddress = "noreply@deiz.fr"

type mailer struct {
	tmpl   *template.Template
	loc    *time.Location
	sender sender
}

type PostFix struct {
	BinPath string
}

type Gmail struct {
	Address  string
	Password string
	Server   string
	Port     int
	ReplyTo  string
}

type sender interface {
	Send(ctx context.Context, m *gomail.Message) error
}

func NewService(tmpl *template.Template, sender sender, loc *time.Location) *mailer {
	return &mailer{
		tmpl:   tmpl,
		sender: sender,
		loc:    loc,
	}
}

func NewPostFixClient() *PostFix {
	return &PostFix{
		BinPath: "/usr/sbin/sendmail",
	}
}

func NewGmailClient() *Gmail {
	return &Gmail{
		Address:  os.Getenv("SMTP_EMAIL"),
		Password: os.Getenv("SMTP_PASSWORD"),
		Server:   os.Getenv("SMTP_SERVER"),
		Port:     465,
		ReplyTo:  noReplyAddress,
	}
}

func (mail *PostFix) Send(ctx context.Context, m *gomail.Message) error {
	cmd := exec.Command(mail.BinPath, "-t")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	pw, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	err = cmd.Start()
	if err != nil {
		return err
	}

	var errs [3]error
	_, errs[0] = m.WriteTo(pw)
	errs[1] = pw.Close()
	errs[2] = cmd.Wait()
	for _, err = range errs {
		if err != nil {
			return err
		}
	}
	return err
}

func (mail *Gmail) Send(ctx context.Context, m *gomail.Message) error {
	dialer := gomail.NewDialer(mail.Server, mail.Port, mail.Address, mail.Password)
	return dialer.DialAndSend(m)
}

func createMail(to string, from string, subject string, tmpl *bytes.Buffer, plainBody string, attachment *bytes.Buffer) *gomail.Message {
	body := tmpl.String()
	m := gomail.NewMessage()
	m.SetAddressHeader("From", noReplyAddress, "Deiz")
	m.SetHeader("Reply-To", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", plainBody)
	m.AddAlternative("text/html", body)
	if attachment != nil {
		m.Attach("doc.pdf", gomail.SetCopyFunc(func(w io.Writer) error {
			_, err := io.Copy(w, attachment)
			return err
		}))
	}
	return m
}
