package email

import (
	"bytes"
	"embed"
	"fmt"
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

//Body of an email
type Body struct {
	HTMLFileName string
	DataToBind   interface{}
	Plain        string
}

//Header of a given email
type Header struct {
	To      string
	From    string
	Subject string
}

//HTMLMail Email data struct
type HTMLMail struct {
	Header     Header
	Body       Body
	Attachment *bytes.Buffer
}

type Send = func(m HTMLMail) error

func SendWithPostfix() Send {
	createMail := createGoMailFunc()
	return func(m HTMLMail) error {
		goMsg, err := createMail(m)
		if err != nil {
			return err
		}
		cmd := getPostfixCmd()
		writer, err := startPostfixCmdWriter(cmd)
		if err != nil {
			return err
		}
		if _, err := goMsg.WriteTo(writer); err != nil {
			return fmt.Errorf("unable to write go msg to postfix writer: %v", err)
		}
		if err := terminatePostfixCmd(writer, cmd); err != nil {
			return err
		}
		return nil
	}
}

func startPostfixCmdWriter(cmd *exec.Cmd) (io.WriteCloser, error) {
	writer, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("unable to writer to stdin pipe: %v", err)
	}
	err = cmd.Start()
	if err != nil {
		return nil, fmt.Errorf("unable to start cmd: %v", err)
	}
	return writer, nil
}

func terminatePostfixCmd(writer io.WriteCloser, cmd *exec.Cmd) error {
	if err := writer.Close(); err != nil {
		return fmt.Errorf("unable to close system writer: %v", err)
	}
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("unable to wait for cmd stdin to be done: %v", err)
	}
	return nil
}

func getPostfixCmd() *exec.Cmd {
	cmd := exec.Command(postFixBinPath, "-t")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd
}

type smtpCfg struct {
	address  string
	password string
	server   string
	port     int
}

func getSmtpCfg() smtpCfg {
	return smtpCfg{
		address:  os.Getenv("SMTP_EMAIL"),
		password: os.Getenv("SMTP_PASSWORD"),
		server:   os.Getenv("SMTP_SERVER"),
		port:     465,
	}
}

func SendWithGmail() Send {
	cfg := getSmtpCfg()
	createMail := createGoMailFunc()
	return func(m HTMLMail) error {
		goMsg, err := createMail(m)
		if err != nil {
			return err
		}
		dialer := gomail.NewDialer(cfg.server, cfg.port, cfg.address, cfg.password)
		return dialer.DialAndSend(goMsg)
	}
}

type createGoMail = func(mail HTMLMail) (*gomail.Message, error)

func createGoMailFunc() createGoMail {
	templates := parseFSTemplates()
	return func(mail HTMLMail) (*gomail.Message, error) {
		m := gomail.NewMessage()
		setGoMailHeader(m, mail.Header)
		if err := setGoMailBody(templates, m, mail.Body); err != nil {
			return nil, err
		}
		setGoMailAttachment(m, mail.Attachment)
		return m, nil
	}
}

func parseFSTemplates() *template.Template {
	templates, err := template.ParseFS(emailFiles, "*")
	if err != nil {
		log.Fatalf("unable To create createGoMailFunc: %v", err)
	}
	return templates
}

func setGoMailBody(templates *template.Template, m *gomail.Message, body Body) error {
	htmlBody, err := parseHTMLBody(templates, body)
	if err != nil {
		return err
	}
	m.SetBody("text/plain", htmlBody)
	m.AddAlternative("text/html", body.Plain)
	return nil
}

func parseHTMLBody(templates *template.Template, b Body) (string, error) {
	bytes, err := tmplToBuffer(templates, b.HTMLFileName, b.DataToBind)
	if err != nil {
		return "", err
	}
	return bytes.String(), nil
}

func setGoMailAttachment(m *gomail.Message, attachment *bytes.Buffer) {
	m.Attach("doc.pdf", gomail.SetCopyFunc(func(w io.Writer) error {
		_, err := io.Copy(w, attachment)
		return err
	}))
}

func setGoMailHeader(m *gomail.Message, header Header) {
	m.SetAddressHeader("From", noReplyAddress, "Deiz")
	m.SetHeader("Reply-To", header.From)
	m.SetHeader("To", header.To)
	m.SetHeader("Subject", header.Subject)
}

func tmplToBuffer(templates *template.Template, file string, dataToBind interface{}) (*bytes.Buffer, error) {
	var b bytes.Buffer
	return &b, templates.ExecuteTemplate(&b, file, dataToBind)
}
