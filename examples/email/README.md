# Email Template Example

This example demonstrates how to use Fíth to generate professional email templates in both HTML and plain text formats.

## Features Demonstrated

- Welcome emails
- Order confirmations
- Password reset emails
- HTML email styling
- Plain text alternatives
- Responsive email design
- Transactional email patterns

## Running the Example

```bash
go run main.go
```

## Email Types

### 1. Welcome Email
- User onboarding
- Next steps guidance
- Call-to-action button
- Unsubscribe link

### 2. Order Confirmation
- Order details table
- Itemized pricing
- Shipping information
- Tracking link

### 3. Plain Text Version
- Accessible alternative
- Better deliverability
- Mobile-friendly

### 4. Password Reset
- Security alert
- Expiring link
- Clear instructions

## Key Concepts

### HTML Email Best Practices

```html
<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8">
  <style>
    /* Inline CSS for email compatibility */
    body { font-family: Arial, sans-serif; }
    .container { max-width: 600px; margin: 0 auto; }
  </style>
</head>
<body>
  <!-- Content here -->
</body>
</html>
```

### Plain Text Alternative

Always provide a plain text version:

```go
htmlVersion, _ := engine.RenderString(htmlTemplate, data)
plainVersion, _ := engine.RenderString(plainTemplate, data)

// Send both versions
sendEmail(to, subject, htmlVersion, plainVersion)
```

### Dynamic Content

```html
{{range .Items}}
<tr>
  <td>{{.Product}}</td>
  <td>{{.Quantity}}</td>
  <td>${{.Price}}</td>
</tr>
{{end}}
```

## Production Integration

### With SMTP

```go
import "net/smtp"

// Generate email
html, _ := engine.Render("welcome-email", data)
plain, _ := engine.Render("welcome-email-plain", data)

// Prepare MIME message
msg := fmt.Sprintf(`From: %s
To: %s
Subject: %s
MIME-Version: 1.0
Content-Type: multipart/alternative; boundary="boundary"

--boundary
Content-Type: text/plain; charset="utf-8"

%s

--boundary
Content-Type: text/html; charset="utf-8"

%s

--boundary--`, from, to, subject, plain, html)

// Send via SMTP
smtp.SendMail(smtpServer, auth, from, []string{to}, []byte(msg))
```

### With Email Service (SendGrid, Mailgun, etc.)

```go
// Generate content
htmlContent, _ := engine.Render("template", data)

// Send via service
client.SendEmail(&Email{
    To:      recipient,
    Subject: subject,
    HTML:    htmlContent,
})
```

## Email Styling Tips

1. **Use inline CSS** - Many email clients strip `<style>` tags
2. **Keep width ≤ 600px** - For mobile compatibility
3. **Use tables for layout** - Better email client support
4. **Test thoroughly** - Different clients render differently
5. **Include alt text** - For images
6. **Use web-safe fonts** - Arial, Helvetica, Times New Roman

## Template Organization

```
templates/emails/
├── base.html           # Shared layout
├── welcome.html        # Welcome email
├── welcome-plain.txt   # Plain text version
├── order.html          # Order confirmation
├── order-plain.txt     # Plain text version
├── reset.html          # Password reset
└── partials/
    ├── header.html
    ├── footer.html
    └── button.html
```

## Testing

```go
// Test email rendering
func TestWelcomeEmail(t *testing.T) {
    engine, _ := fith.NewWithDefaults()
    
    data := map[string]interface{}{
        "UserName": "Test User",
        "AppName": "TestApp",
    }
    
    output, err := engine.RenderString(welcomeTemplate, data)
    require.NoError(t, err)
    assert.Contains(t, output, "Test User")
    assert.Contains(t, output, "TestApp")
}
```

## Use Cases

- User registration confirmations
- Order and shipping notifications
- Password reset flows
- Newsletter campaigns
- Notification emails
- Marketing emails
- Transactional receipts
- Account alerts

## Security Considerations

1. **Sanitize user input** - Use `htmlEscape` for user data
2. **Use HTTPS links** - For all URLs
3. **Implement rate limiting** - Prevent abuse
4. **Add unsubscribe links** - Required by law (CAN-SPAM)
5. **Use secure tokens** - For password resets

## Deliverability Tips

1. Include both HTML and plain text
2. Keep email size < 100KB
3. Use a proper "From" address
4. Include an unsubscribe link
5. Authenticate with SPF/DKIM
6. Avoid spammy words
7. Test with spam checkers
