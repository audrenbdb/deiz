package email

import (
	"bytes"
	"embed"
	"gopkg.in/gomail.v2"
	"io"
	"log"
	"os"
	"os/exec"
	"text/template"
)

const noReplyAddress = "noreply@deiz.fr"
const postFixBinPath = "/usr/sbin/sendmail"

//go:embed *.html
var emailFiles embed.FS

//HTML email data required to send an html email
type HTML struct {
	To               string
	From             string
	Subject          string
	TemplateFilePath string
	DataToBind       interface{}
	PlainBody        string
	Attachment       *bytes.Buffer
}

type Send = func(m HTML) error

func SendWithPostfix() Send {
	createMail := createGoMailFn()
	return func(m HTML) error {
		goMsg, err := createMail(m)
		if err != nil {
			return err
		}
		cmd := exec.Command(postFixBinPath, "-t")
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
		_, errs[0] = goMsg.WriteTo(pw)
		errs[1] = pw.Close()
		errs[2] = cmd.Wait()
		for _, err = range errs {
			if err != nil {
				return err
			}
		}
		return nil
	}
}

func SendWithGmail() Send {
	address := os.Getenv("SMTP_EMAIL")
	pw := os.Getenv("SMTP_PASSWORD")
	server := os.Getenv("SMTP_SERVER")
	port := 465
	createMail := createGoMailFn()
	return func(m HTML) error {
		goMsg, err := createMail(m)
		if err != nil {
			return err
		}
		dialer := gomail.NewDialer(server, port, address, pw)
		return dialer.DialAndSend(goMsg)
	}
}

type createGoMail = func(mail HTML) (*gomail.Message, error)

func createGoMailFn() createGoMail {
	templates, err := template.ParseFS(emailFiles, "*")
	if err != nil {
		log.Fatalf("unable To create createGoMailFn: %v", err)
	}
	return func(mail HTML) (*gomail.Message, error) {
		b, err := tmplToBuffer(templates, mail.TemplateFilePath, mail.DataToBind)
		if err != nil {
			return nil, err
		}
		body := b.String()
		m := gomail.NewMessage()
		m.SetAddressHeader("From", noReplyAddress, "Deiz")
		m.SetHeader("Reply-To", mail.From)
		m.SetHeader("To", mail.To)
		m.SetHeader("Subject", mail.Subject)
		m.SetBody("text/plain", mail.PlainBody)
		m.AddAlternative("text/html", body)
		if mail.Attachment != nil {
			m.Attach("doc.pdf", gomail.SetCopyFunc(func(w io.Writer) error {
				_, err := io.Copy(w, mail.Attachment)
				return err
			}))
		}
		return m, nil
	}
}

func tmplToBuffer(templates *template.Template, file string, dataToBind interface{}) (*bytes.Buffer, error) {
	var b bytes.Buffer
	return &b, templates.ExecuteTemplate(&b, file, dataToBind)
}
