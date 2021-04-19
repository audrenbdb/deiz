package mail

import (
	"bytes"
	"github.com/audrenbdb/deiz/intl"
	"gopkg.in/gomail.v2"
	"io"
	"os"
	"os/exec"
	"text/template"
)

const noReplyAddress = "noreply@deiz.fr"

type Mailer struct {
	tmpl   *template.Template
	intl   *intl.Parser
	client client
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

type client interface {
	Send(m *gomail.Message) error
}

type Deps struct {
	Templates *template.Template
	Client    client
	Intl      *intl.Parser
}

func NewService(deps Deps) *Mailer {
	return &Mailer{
		tmpl:   deps.Templates,
		client: deps.Client,
		intl:   deps.Intl,
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

func (mailer *PostFix) Send(m *gomail.Message) error {
	cmd := exec.Command(mailer.BinPath, "-t")
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

func (mailer *Gmail) Send(m *gomail.Message) error {
	dialer := gomail.NewDialer(mailer.Server, mailer.Port, mailer.Address, mailer.Password)
	return dialer.DialAndSend(m)
}

type mail struct {
	to         string
	from       string
	subject    string
	template   *bytes.Buffer
	plainBody  string
	attachment *bytes.Buffer
}

func createMail(mail mail) *gomail.Message {
	body := mail.template.String()
	m := gomail.NewMessage()
	m.SetAddressHeader("From", noReplyAddress, "Deiz")
	m.SetHeader("Reply-To", mail.from)
	m.SetHeader("To", mail.to)
	m.SetHeader("Subject", mail.subject)
	m.SetBody("text/plain", mail.plainBody)
	m.AddAlternative("text/html", body)
	if mail.attachment != nil {
		m.Attach("doc.pdf", gomail.SetCopyFunc(func(w io.Writer) error {
			_, err := io.Copy(w, mail.attachment)
			return err
		}))
	}
	return m
}

func (m *Mailer) htmlTemplate(name string, details interface{}) (*bytes.Buffer, error) {
	var emailBuffer bytes.Buffer
	return &emailBuffer, m.tmpl.ExecuteTemplate(&emailBuffer, name, details)
}
