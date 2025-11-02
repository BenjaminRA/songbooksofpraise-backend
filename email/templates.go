package email

import (
	"fmt"
	"strings"
)

// EmailTemplate generates a professional HTML email template
// Params:
//   - title: Email title/heading
//   - preheader: Short preview text (optional)
//   - content: Main email content (can include HTML)
//   - buttonText: CTA button text (optional)
//   - buttonLink: CTA button link (optional)
//   - footerText: Additional footer text (optional)
func EmailTemplate(title, preheader, content, buttonText, buttonLink, footerText string) string {
	// If no button is needed, remove the button section
	buttonHTML := ""
	if buttonText != "" && buttonLink != "" {
		buttonHTML = fmt.Sprintf(`
			<tr>
				<td style="padding: 30px 0;">
					<table role="presentation" cellspacing="0" cellpadding="0" border="0" align="center">
						<tr>
							<td style="border-radius: 4px; background: #5048E5;">
								<a href="%s" target="_blank" style="background: #5048E5; border: none; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif; font-size: 16px; color: #ffffff; text-decoration: none; border-radius: 4px; padding: 12px 40px; display: inline-block; font-weight: 600;">
									%s
								</a>
							</td>
						</tr>
					</table>
				</td>
			</tr>
		`, buttonLink, buttonText)
	}

	// Optional footer text (appears outside the card, above copyright)
	additionalFooter := ""
	if footerText != "" {
		additionalFooter = fmt.Sprintf(`
				<!-- Additional Footer Message -->
				<table role="presentation" cellspacing="0" cellpadding="0" border="0" width="100%%" style="max-width: 600px; margin: 20px auto 0 auto;">
					<tr>
						<td style="padding: 0 20px; text-align: center; color: #65748B; font-size: 13px; line-height: 1.6;">
							<p style="margin: 0;">
								%s
							</p>
						</td>
					</tr>
				</table>
		`, footerText)
	}

	// Preheader text (shows in email preview)
	preheaderHTML := ""
	if preheader != "" {
		preheaderHTML = fmt.Sprintf(`
			<div style="display: none; max-height: 0; overflow: hidden;">
				%s
			</div>
		`, preheader)
	}

	// Main template
	template := fmt.Sprintf(`
<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<meta http-equiv="X-UA-Compatible" content="IE=edge">
	<title>%s</title>
	<!--[if mso]>
	<style type="text/css">
		body, table, td {font-family: Arial, Helvetica, sans-serif !important;}
	</style>
	<![endif]-->
</head>
<body style="margin: 0; padding: 0; background-color: #F9FAFC; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif;">
	%s
	
	<!-- Main Container -->
	<table role="presentation" cellspacing="0" cellpadding="0" border="0" width="100%%" style="background-color: #F9FAFC;">
		<tr>
			<td style="padding: 40px 20px;">
				
				<!-- Email Content -->
				<table role="presentation" cellspacing="0" cellpadding="0" border="0" width="100%%" style="max-width: 600px; margin: 0 auto; background-color: #FFFFFF; border-radius: 8px; box-shadow: 0 2px 8px rgba(0,0,0,0.05);">
					
					<!-- Header with Logo/Brand -->
					<tr>
						<td style="background: linear-gradient(135deg, #5048E5 0%%, #3832A0 100%%); padding: 40px 40px 30px 40px; text-align: center; border-radius: 8px 8px 0 0;">
							<h1 style="margin: 0; color: #FFFFFF; font-size: 28px; font-weight: 700; letter-spacing: -0.5px;">
								Songbooks of Praise
							</h1>
							<p style="margin: 8px 0 0 0; color: rgba(255,255,255,0.9); font-size: 14px; font-weight: 500;">
								🎵 Worship Management System
							</p>
						</td>
					</tr>
					
					<!-- Email Title -->
					<tr>
						<td style="padding: 40px 40px 20px 40px;">
							<h2 style="margin: 0; color: #121828; font-size: 24px; font-weight: 600; line-height: 1.3;">
								%s
							</h2>
						</td>
					</tr>
					
					<!-- Email Content -->
					<tr>
						<td style="padding: 0 40px; color: #374151; font-size: 16px; line-height: 1.6;">
							%s
						</td>
					</tr>
					
					<!-- CTA Button (if provided) -->
					%s
					
					<!-- Divider -->
					<tr>
						<td style="padding: 30px 40px;">
							<div style="border-top: 1px solid #E6E8F0;"></div>
						</td>
					</tr>
					
					<!-- Footer -->
					<tr>
						<td style="padding: 0 40px 40px 40px; color: #65748B; font-size: 14px; line-height: 1.6;">
							<p style="margin: 0 0 12px 0;">
								<strong style="color: #374151;">Songbooks of Praise Team</strong>
							</p>
							<p style="margin: 0;">
								This email was sent from the Songbooks of Praise worship management system.
							</p>
						</td>
					</tr>
					
				</table>
				
				<!-- Additional Footer Message (if provided) -->
				%s
				
				<!-- Footer Note -->
				<table role="presentation" cellspacing="0" cellpadding="0" border="0" width="100%%" style="max-width: 600px; margin: 20px auto 0 auto;">
					<tr>
						<td style="padding: 0 20px; text-align: center; color: #9CA3AF; font-size: 12px; line-height: 1.5;">
							<p style="margin: 0;">
								© %d Songbooks of Praise. All rights reserved.
							</p>
						</td>
					</tr>
				</table>
				
			</td>
		</tr>
	</table>
	
</body>
</html>
	`, title, preheaderHTML, title, content, buttonHTML, additionalFooter, getCurrentYear())

	return strings.TrimSpace(template)
}

// getCurrentYear returns the current year for the copyright notice
func getCurrentYear() int {
	return 2024 // You can make this dynamic if needed: time.Now().Year()
}

// Helper function to create simple paragraph content
func Paragraph(text string) string {
	return fmt.Sprintf(`<p style="margin: 0 0 16px 0;">%s</p>`, text)
}

// Helper function to create a list
func List(items []string) string {
	listItems := ""
	for _, item := range items {
		listItems += fmt.Sprintf(`<li style="margin: 0 0 8px 0;">%s</li>`, item)
	}
	return fmt.Sprintf(`<ul style="margin: 0 0 16px 0; padding-left: 20px;">%s</ul>`, listItems)
}

// Helper function to create highlighted text
func Highlight(text string) string {
	return fmt.Sprintf(`<strong style="color: #5048E5;">%s</strong>`, text)
}

// Helper function to create a warning/info box
func InfoBox(text string, boxType string) string {
	var bgColor, borderColor, textColor string
	
	switch boxType {
	case "success":
		bgColor = "#ECFDF5"
		borderColor = "#10B981"
		textColor = "#065F46"
	case "warning":
		bgColor = "#FFF7ED"
		borderColor = "#FFB020"
		textColor = "#B27B16"
	case "error":
		bgColor = "#FEF2F2"
		borderColor = "#D14343"
		textColor = "#922E2E"
	default: // info
		bgColor = "#EFF6FF"
		borderColor = "#2196F3"
		textColor = "#0B79D0"
	}
	
	return fmt.Sprintf(`
		<div style="margin: 0 0 16px 0; padding: 16px; background-color: %s; border-left: 4px solid %s; border-radius: 4px;">
			<p style="margin: 0; color: %s; font-size: 14px; line-height: 1.5;">%s</p>
		</div>
	`, bgColor, borderColor, textColor, text)
}
