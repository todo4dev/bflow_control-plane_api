package mailer

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"mime"
	"mime/quotedprintable"
	"net/smtp"
	"strings"
	"time"

	adapter "src/application/adapter/mailer"
)

type SMTPMailerAdapter struct {
	config *adapter.MailerConfig
}

var _ adapter.IMailerAdapter = (*SMTPMailerAdapter)(nil)

func NewSMTPMailerAdapter(config *adapter.MailerConfig) *SMTPMailerAdapter {
	if config == nil {
		panic("smtp mailer: config cannot be nil")
	}

	if config.SMTPHostname == "" {
		panic("smtp mailer: SMTP.Hostname is required")
	}

	if config.SMTPPort == 0 {
		config.SMTPPort = 587
	}

	return &SMTPMailerAdapter{config: config}
}

func (a *SMTPMailerAdapter) Send(ctx context.Context, input adapter.MailPayload) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	from := input.From
	if from == "" {
		from = a.config.DefaultFrom
	}
	if from == "" {
		return fmt.Errorf("smtp mailer: from address is empty")
	}

	msg, err := buildMIMEMessage(from, input)
	if err != nil {
		return err
	}

	host := a.config.SMTPHostname
	addr := fmt.Sprintf("%s:%d", host, a.config.SMTPPort)

	recipients := dedupeRecipients(input.To, input.CC, input.BCC)
	if len(recipients) == 0 {
		return fmt.Errorf("smtp mailer: no recipients")
	}

	client, err := a.newSMTPClient(addr, host)
	if err != nil {
		return err
	}
	defer client.Close()

	// auth se tiver usuário/senha
	if a.config.SMTPUsername != "" || a.config.SMTPPassword != "" {
		auth := smtp.PlainAuth("", a.config.SMTPUsername, a.config.SMTPPassword, host)
		if err := client.Auth(auth); err != nil {
			return fmt.Errorf("smtp mailer: auth failed: %w", err)
		}
	}

	envelopeFrom := extractEmailAddress(from)

	if err := client.Mail(envelopeFrom); err != nil {
		return fmt.Errorf("smtp mailer: MAIL FROM failed: %w", err)
	}

	for _, rcpt := range recipients {
		if err := client.Rcpt(rcpt); err != nil {
			return fmt.Errorf("smtp mailer: RCPT TO %s failed: %w", rcpt, err)
		}
	}

	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("smtp mailer: DATA failed: %w", err)
	}

	if _, err := w.Write(msg); err != nil {
		_ = w.Close()
		return fmt.Errorf("smtp mailer: write failed: %w", err)
	}

	if err := w.Close(); err != nil {
		return fmt.Errorf("smtp mailer: close failed: %w", err)
	}

	return client.Quit()
}

func (a *SMTPMailerAdapter) newSMTPClient(addr, host string) (*smtp.Client, error) {
	if a.config.SMTPSecure {
		tlsConfig := &tls.Config{
			ServerName:         host,
			InsecureSkipVerify: a.config.SMTPIgnoreTLS,
		}

		conn, err := tls.Dial("tcp", addr, tlsConfig)
		if err != nil {
			return nil, fmt.Errorf("smtp mailer: tls dial failed: %w", err)
		}

		client, err := smtp.NewClient(conn, host)
		if err != nil {
			_ = conn.Close()
			return nil, fmt.Errorf("smtp mailer: new client failed: %w", err)
		}

		return client, nil
	}

	client, err := smtp.Dial(addr)
	if err != nil {
		return nil, fmt.Errorf("smtp mailer: dial failed: %w", err)
	}

	// tenta STARTTLS se disponível
	if ok, _ := client.Extension("STARTTLS"); ok {
		tlsConfig := &tls.Config{
			ServerName:         host,
			InsecureSkipVerify: a.config.SMTPIgnoreTLS,
		}
		if err := client.StartTLS(tlsConfig); err != nil {
			client.Close()
			return nil, fmt.Errorf("smtp mailer: STARTTLS failed: %w", err)
		}
	}

	return client, nil
}

func buildMIMEMessage(from string, input adapter.MailPayload) ([]byte, error) {
	var buf bytes.Buffer

	buf.WriteString("From: " + from + "\r\n")
	if len(input.To) > 0 {
		buf.WriteString("To: " + strings.Join(input.To, ", ") + "\r\n")
	}
	if len(input.CC) > 0 {
		buf.WriteString("Cc: " + strings.Join(input.CC, ", ") + "\r\n")
	}
	if len(input.ReplyTo) > 0 {
		buf.WriteString("Reply-To: " + strings.Join(input.ReplyTo, ", ") + "\r\n")
	}
	if input.Subject != "" {
		encodedSubject := mime.QEncoding.Encode("utf-8", input.Subject)
		buf.WriteString("Subject: " + encodedSubject + "\r\n")
	}
	buf.WriteString("MIME-Version: 1.0\r\n")

	hasAttachments := len(input.Attachments) > 0
	hasHTML := input.HTML != ""
	hasText := input.Text != ""

	// caso simples: sem anexo e só um tipo de corpo
	if !hasAttachments && (!hasHTML || !hasText) {
		if hasHTML {
			buf.WriteString("Content-Type: text/html; charset=\"utf-8\"\r\n")
		} else {
			buf.WriteString("Content-Type: text/plain; charset=\"utf-8\"\r\n")
		}
		buf.WriteString("Content-Transfer-Encoding: quoted-printable\r\n\r\n")

		qp := quotedprintable.NewWriter(&buf)
		body := input.Text
		if hasHTML {
			body = input.HTML
		}
		if _, err := qp.Write([]byte(body)); err != nil {
			return nil, err
		}
		if err := qp.Close(); err != nil {
			return nil, err
		}

		return buf.Bytes(), nil
	}

	// multipart
	boundaryMixed := fmt.Sprintf("mixed_%d", time.Now().UnixNano())

	if hasHTML && hasText {
		// multipart/mixed + multipart/alternative
		buf.WriteString("Content-Type: multipart/mixed; boundary=" + boundaryMixed + "\r\n\r\n")

		boundaryAlt := fmt.Sprintf("alt_%d", time.Now().UnixNano())
		fmt.Fprintf(&buf, "--%s\r\n", boundaryMixed)
		fmt.Fprintf(&buf, "Content-Type: multipart/alternative; boundary=%s\r\n\r\n", boundaryAlt)

		// text/plain
		fmt.Fprintf(&buf, "--%s\r\n", boundaryAlt)
		buf.WriteString("Content-Type: text/plain; charset=\"utf-8\"\r\n")
		buf.WriteString("Content-Transfer-Encoding: quoted-printable\r\n\r\n")
		qp := quotedprintable.NewWriter(&buf)
		if _, err := qp.Write([]byte(input.Text)); err != nil {
			return nil, err
		}
		if err := qp.Close(); err != nil {
			return nil, err
		}
		buf.WriteString("\r\n")

		// text/html
		fmt.Fprintf(&buf, "--%s\r\n", boundaryAlt)
		buf.WriteString("Content-Type: text/html; charset=\"utf-8\"\r\n")
		buf.WriteString("Content-Transfer-Encoding: quoted-printable\r\n\r\n")
		qp = quotedprintable.NewWriter(&buf)
		if _, err := qp.Write([]byte(input.HTML)); err != nil {
			return nil, err
		}
		if err := qp.Close(); err != nil {
			return nil, err
		}
		buf.WriteString("\r\n")

		fmt.Fprintf(&buf, "--%s--\r\n", boundaryAlt)
	} else {
		// só um tipo de corpo, mas com anexo
		buf.WriteString("Content-Type: multipart/mixed; boundary=" + boundaryMixed + "\r\n\r\n")

		fmt.Fprintf(&buf, "--%s\r\n", boundaryMixed)
		if hasHTML {
			buf.WriteString("Content-Type: text/html; charset=\"utf-8\"\r\n")
		} else {
			buf.WriteString("Content-Type: text/plain; charset=\"utf-8\"\r\n")
		}
		buf.WriteString("Content-Transfer-Encoding: quoted-printable\r\n\r\n")

		qp := quotedprintable.NewWriter(&buf)
		body := input.Text
		if hasHTML {
			body = input.HTML
		}
		if _, err := qp.Write([]byte(body)); err != nil {
			return nil, err
		}
		if err := qp.Close(); err != nil {
			return nil, err
		}
		buf.WriteString("\r\n")
	}

	// anexos
	for _, att := range input.Attachments {
		if len(att.Content) == 0 {
			continue
		}

		filename := att.Filename
		if filename == "" {
			filename = "attachment"
		}
		contentType := att.ContentType
		if contentType == "" {
			contentType = "app/octet-stream"
		}

		encodedName := mime.QEncoding.Encode("utf-8", filename)

		fmt.Fprintf(&buf, "--%s\r\n", boundaryMixed)
		fmt.Fprintf(&buf, "Content-Type: %s; name=\"%s\"\r\n", contentType, encodedName)
		buf.WriteString("Content-Transfer-Encoding: base64\r\n")
		fmt.Fprintf(&buf, "Content-Disposition: attachment; filename=\"%s\"\r\n\r\n", encodedName)

		encoded := base64.StdEncoding.EncodeToString(att.Content)
		for i := 0; i < len(encoded); i += 76 {
			end := i + 76
			if end > len(encoded) {
				end = len(encoded)
			}
			buf.WriteString(encoded[i:end] + "\r\n")
		}
	}

	fmt.Fprintf(&buf, "--%s--\r\n", boundaryMixed)

	return buf.Bytes(), nil
}

func dedupeRecipients(groups ...[]string) []string {
	seen := make(map[string]struct{})
	var result []string

	for _, list := range groups {
		for _, addr := range list {
			addr = strings.TrimSpace(addr)
			if addr == "" {
				continue
			}
			if _, ok := seen[addr]; ok {
				continue
			}
			seen[addr] = struct{}{}
			result = append(result, addr)
		}
	}

	return result
}

func extractEmailAddress(addr string) string {
	addr = strings.TrimSpace(addr)
	if addr == "" {
		return addr
	}

	if i := strings.LastIndex(addr, "<"); i != -1 && strings.Contains(addr[i:], "@") {
		j := strings.Index(addr[i:], ">")
		if j != -1 {
			return strings.TrimSpace(addr[i+1 : i+j])
		}
	}

	return addr
}
