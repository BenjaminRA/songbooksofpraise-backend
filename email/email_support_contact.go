package email

import (
	"fmt"
	"html"
	"os"
	"strings"
)

// SendSupportContactEmail sends a support inquiry from the support page contact form
// to the support inbox. It uses inline styles matching the support page's
// crimson / gold / cream palette.
func SendSupportContactEmail(name, fromEmail, topic, message, lang string) error {
	subject := fmt.Sprintf("Support Request: %s", topic)
	if lang == "es" {
		subject = fmt.Sprintf("Solicitud de Soporte: %s", topic)
	}

	to := os.Getenv("SUPPORT_EMAIL")
	if to == "" {
		to = os.Getenv("MAIL_USERNAME")
	}

	return SendEmail(to, subject, supportContactEmailTemplate(name, fromEmail, topic, message, lang))
}

// supportContactEmailTemplate builds a fully inline-styled HTML email that
// mirrors the crimson / gold / cream design of support-page/index.html.
func supportContactEmailTemplate(name, fromEmail, topic, message, lang string) string {
	// Sanitize all user inputs – prevents HTML injection in the email body.
	safeName := html.EscapeString(name)
	safeEmail := html.EscapeString(fromEmail)
	safeTopic := html.EscapeString(topic)
	// Escape then convert newlines so multi-line messages render correctly.
	safeMessage := strings.ReplaceAll(html.EscapeString(message), "\n", "<br>")

	// Default to English labels; override for Spanish.
	langAttr := "en"
	heading := "New Support Message"
	subheading := "A user has submitted a message through the Songbooks of Praise support form."
	labelName := "Name"
	labelEmail := "Email"
	labelTopic := "Topic"
	labelMessage := "Message"
	footerNote := "This message was submitted via the Songbooks of Praise support page."

	if lang == "es" {
		langAttr = "es"
		heading = "Nuevo Mensaje de Soporte"
		subheading = "Un usuario ha enviado un mensaje a través del formulario de soporte de Songbooks of Praise."
		labelName = "Nombre"
		labelEmail = "Correo Electrónico"
		labelTopic = "Tema"
		labelMessage = "Mensaje"
		footerNote = "Este mensaje fue enviado desde la página de soporte de Songbooks of Praise."
	}

	// fieldRow renders one labeled field block (label on top, value box below).
	fieldRow := func(label, value string) string {
		return fmt.Sprintf(`
              <table role="presentation" cellspacing="0" cellpadding="0" border="0" width="100%%" style="margin-bottom:20px;">
                <tr>
                  <td style="padding-bottom:4px;">
                    <span style="font-size:11px;font-weight:700;letter-spacing:0.10em;text-transform:uppercase;color:#8B1A1A;">%s</span>
                  </td>
                </tr>
                <tr>
                  <td style="background-color:#F5F0E8;border-left:3px solid #C9A84C;padding:10px 14px;border-radius:0 4px 4px 0;">
                    <span style="font-size:15px;color:#1A1A1A;font-weight:400;line-height:1.6;">%s</span>
                  </td>
                </tr>
              </table>`, label, value)
	}

	// Email field gets a mailto: link so the recipient can reply in one click.
	emailRow := fmt.Sprintf(`
              <table role="presentation" cellspacing="0" cellpadding="0" border="0" width="100%%" style="margin-bottom:20px;">
                <tr>
                  <td style="padding-bottom:4px;">
                    <span style="font-size:11px;font-weight:700;letter-spacing:0.10em;text-transform:uppercase;color:#8B1A1A;">%s</span>
                  </td>
                </tr>
                <tr>
                  <td style="background-color:#F5F0E8;border-left:3px solid #C9A84C;padding:10px 14px;border-radius:0 4px 4px 0;">
                    <a href="mailto:%s" style="font-size:15px;color:#6B1212;font-weight:500;text-decoration:none;">%s</a>
                  </td>
                </tr>
              </table>`, labelEmail, safeEmail, safeEmail)

	// Message field keeps the left-border style but allows taller content.
	messageRow := fmt.Sprintf(`
              <table role="presentation" cellspacing="0" cellpadding="0" border="0" width="100%%" style="margin-bottom:28px;">
                <tr>
                  <td style="padding-bottom:4px;">
                    <span style="font-size:11px;font-weight:700;letter-spacing:0.10em;text-transform:uppercase;color:#8B1A1A;">%s</span>
                  </td>
                </tr>
                <tr>
                  <td style="background-color:#F5F0E8;border-left:3px solid #C9A84C;padding:14px;border-radius:0 4px 4px 0;">
                    <span style="font-size:15px;color:#4A4A4A;line-height:1.7;">%s</span>
                  </td>
                </tr>
              </table>`, labelMessage, safeMessage)

	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="%s">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <meta http-equiv="X-UA-Compatible" content="IE=edge">
  <title>%s</title>
</head>
<body style="margin:0;padding:0;background-color:#F5F0E8;font-family:Arial,Helvetica,sans-serif;">

  <table role="presentation" cellspacing="0" cellpadding="0" border="0" width="100%%" style="background-color:#F5F0E8;">
    <tr>
      <td style="padding:40px 20px;">

        <!-- Card -->
        <table role="presentation" cellspacing="0" cellpadding="0" border="0" width="100%%" style="max-width:600px;margin:0 auto;background-color:#FFFFFF;border-radius:8px;box-shadow:0 2px 12px rgba(107,18,18,0.10);border:1px solid rgba(139,26,26,0.12);">

          <!-- Header band -->
          <tr>
            <td style="background-color:#6B1212;padding:36px 40px 28px 40px;text-align:center;border-radius:8px 8px 0 0;">
              <div style="display:inline-block;background-color:#C9A84C;border-radius:8px;width:44px;height:44px;line-height:44px;text-align:center;font-size:22px;margin-bottom:14px;">&#127925;</div>
              <h1 style="margin:0;color:#FFFFFF;font-family:Georgia,'Times New Roman',serif;font-size:24px;font-weight:700;letter-spacing:0.01em;line-height:1.2;">Songbooks of Praise</h1>
              <p style="margin:8px 0 0 0;color:rgba(255,255,255,0.65);font-size:11px;font-weight:400;letter-spacing:0.14em;text-transform:uppercase;">Support Center</p>
            </td>
          </tr>

          <!-- Gold accent line -->
          <tr>
            <td style="background-color:#C9A84C;height:3px;font-size:0;line-height:0;">&nbsp;</td>
          </tr>

          <!-- Heading -->
          <tr>
            <td style="padding:36px 40px 8px 40px;">
              <h2 style="margin:0;color:#6B1212;font-family:Georgia,'Times New Roman',serif;font-size:20px;font-weight:700;line-height:1.3;">%s</h2>
              <p style="margin:10px 0 0 0;color:#7A7A7A;font-size:14px;line-height:1.5;">%s</p>
            </td>
          </tr>

          <!-- Divider -->
          <tr>
            <td style="padding:16px 40px 0 40px;">
              <div style="border-top:1px solid #EDE6D6;"></div>
            </td>
          </tr>

          <!-- Fields -->
          <tr>
            <td style="padding:24px 40px 0 40px;">
              %s
              %s
              %s
              %s
            </td>
          </tr>

          <!-- Footer inside card -->
          <tr>
            <td style="background-color:#EDE6D6;padding:20px 40px;border-radius:0 0 8px 8px;border-top:1px solid rgba(139,26,26,0.08);">
              <p style="margin:0;color:#7A7A7A;font-size:12px;line-height:1.6;text-align:center;">%s</p>
            </td>
          </tr>

        </table>

        <!-- Copyright line -->
        <table role="presentation" cellspacing="0" cellpadding="0" border="0" width="100%%" style="max-width:600px;margin:16px auto 0 auto;">
          <tr>
            <td style="text-align:center;color:#7A7A7A;font-size:11px;padding:0 20px;">
              &#169; 2026 Songbooks of Praise. All rights reserved.
            </td>
          </tr>
        </table>

      </td>
    </tr>
  </table>

</body>
</html>`,
		langAttr,
		heading,
		heading,
		subheading,
		fieldRow(labelName, safeName),
		emailRow,
		fieldRow(labelTopic, safeTopic),
		messageRow,
		footerNote,
	)
}
