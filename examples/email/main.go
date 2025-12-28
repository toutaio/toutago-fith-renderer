// Package main demonstrates email template generation with Fíth.
// Shows how to create HTML and plain text email templates.
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/toutaio/toutago-fith-renderer"
)

func main() {
	engine, err := fith.NewWithDefaults()
	if err != nil {
		log.Fatal(err)
	}

	// Example 1: Welcome email
	fmt.Println("=== Example 1: Welcome Email (HTML) ===")
	welcomeHTML := `<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8">
  <style>
    body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
    .container { max-width: 600px; margin: 0 auto; padding: 20px; }
    .header { background: #4CAF50; color: white; padding: 20px; text-align: center; }
    .content { padding: 20px; background: #f9f9f9; }
    .button { display: inline-block; padding: 10px 20px; background: #4CAF50; color: white; text-decoration: none; border-radius: 5px; }
    .footer { text-align: center; padding: 20px; font-size: 12px; color: #666; }
  </style>
</head>
<body>
  <div class="container">
    <div class="header">
      <h1>Welcome to {{.AppName}}!</h1>
    </div>
    
    <div class="content">
      <p>Hi {{.UserName}},</p>
      
      <p>Thank you for signing up! We're excited to have you on board.</p>
      
      <p>Here's what you can do next:</p>
      <ul>
        {{range .NextSteps}}
        <li>{{.}}</li>
        {{end}}
      </ul>
      
      <p style="text-align: center; margin: 30px 0;">
        <a href="{{.ActivationLink}}" class="button">Activate Your Account</a>
      </p>
      
      <p>If you have any questions, feel free to reply to this email.</p>
      
      <p>Best regards,<br>The {{.AppName}} Team</p>
    </div>
    
    <div class="footer">
      <p>{{.CompanyName}} | {{.CompanyAddress}}</p>
      <p><a href="{{.UnsubscribeLink}}">Unsubscribe</a></p>
    </div>
  </div>
</body>
</html>`

	welcomeData := map[string]interface{}{
		"AppName":      "MyApp",
		"UserName":     "Alice",
		"CompanyName":  "Tech Corp",
		"CompanyAddress": "123 Tech St, San Francisco, CA",
		"NextSteps": []string{
			"Complete your profile",
			"Connect with friends",
			"Explore features",
		},
		"ActivationLink":   "https://example.com/activate/abc123",
		"UnsubscribeLink":  "https://example.com/unsubscribe/abc123",
	}

	output, err := engine.RenderString(welcomeHTML, welcomeData)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(output)
	fmt.Println()

	// Example 2: Order confirmation
	fmt.Println("=== Example 2: Order Confirmation ===")
	orderHTML := `<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8">
  <style>
    body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
    .container { max-width: 600px; margin: 0 auto; }
    .header { background: #2196F3; color: white; padding: 20px; }
    .order-details { padding: 20px; }
    table { width: 100%; border-collapse: collapse; margin: 20px 0; }
    th, td { padding: 10px; text-align: left; border-bottom: 1px solid #ddd; }
    th { background: #f5f5f5; }
    .total { font-size: 18px; font-weight: bold; }
    .footer { text-align: center; padding: 20px; background: #f5f5f5; }
  </style>
</head>
<body>
  <div class="container">
    <div class="header">
      <h1>Order Confirmation</h1>
      <p>Order #{{.OrderID}}</p>
    </div>
    
    <div class="order-details">
      <p>Hi {{.CustomerName}},</p>
      
      <p>Thank you for your order! Your order has been confirmed and will be shipped soon.</p>
      
      <h2>Order Details</h2>
      <p><strong>Order Date:</strong> {{date "January 2, 2006" .OrderDate}}<br>
         <strong>Shipping Address:</strong> {{.ShippingAddress}}</p>
      
      <table>
        <thead>
          <tr>
            <th>Product</th>
            <th>Quantity</th>
            <th>Price</th>
            <th>Subtotal</th>
          </tr>
        </thead>
        <tbody>
          {{range .Items}}
          <tr>
            <td>{{.Product}}</td>
            <td>{{.Quantity}}</td>
            <td>${{.Price}}</td>
            <td>${{.Subtotal}}</td>
          </tr>
          {{end}}
        </tbody>
        <tfoot>
          <tr>
            <td colspan="3" style="text-align: right;"><strong>Shipping:</strong></td>
            <td>${{.Shipping}}</td>
          </tr>
          <tr>
            <td colspan="3" style="text-align: right;"><strong>Tax:</strong></td>
            <td>${{.Tax}}</td>
          </tr>
          <tr class="total">
            <td colspan="3" style="text-align: right;">Total:</td>
            <td>${{.Total}}</td>
          </tr>
        </tfoot>
      </table>
      
      <p>You can track your order at: <a href="{{.TrackingLink}}">{{.TrackingLink}}</a></p>
    </div>
    
    <div class="footer">
      <p>Questions? Contact us at support@example.com</p>
    </div>
  </div>
</body>
</html>`

	orderData := map[string]interface{}{
		"OrderID":       "12345",
		"CustomerName":  "Bob Smith",
		"OrderDate":     time.Now(),
		"ShippingAddress": "456 Main St, New York, NY 10001",
		"Items": []map[string]interface{}{
			{"Product": "Laptop", "Quantity": 1, "Price": 999.99, "Subtotal": 999.99},
			{"Product": "Mouse", "Quantity": 2, "Price": 29.99, "Subtotal": 59.98},
			{"Product": "Keyboard", "Quantity": 1, "Price": 79.99, "Subtotal": 79.99},
		},
		"Shipping":     10.00,
		"Tax":         113.00,
		"Total":       1263.96,
		"TrackingLink": "https://example.com/track/12345",
	}

	output, err = engine.RenderString(orderHTML, orderData)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(output)
	fmt.Println()

	// Example 3: Plain text version
	fmt.Println("=== Example 3: Plain Text Email ===")
	plainText := `Welcome to {{.AppName}}!

Hi {{.UserName}},

Thank you for signing up! We're excited to have you on board.

Here's what you can do next:
{{range .NextSteps}}
• {{.}}
{{end}}

To activate your account, visit:
{{.ActivationLink}}

If you have any questions, feel free to reply to this email.

Best regards,
The {{.AppName}} Team

---
{{.CompanyName}} | {{.CompanyAddress}}
To unsubscribe: {{.UnsubscribeLink}}`

	output, err = engine.RenderString(plainText, welcomeData)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(output)
	fmt.Println()

	// Example 4: Password reset email
	fmt.Println("=== Example 4: Password Reset ===")
	resetHTML := `<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8">
  <style>
    body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
    .container { max-width: 600px; margin: 0 auto; padding: 20px; }
    .alert { background: #fff3cd; border-left: 4px solid #ffc107; padding: 15px; margin: 20px 0; }
    .button { display: inline-block; padding: 12px 24px; background: #f44336; color: white; text-decoration: none; border-radius: 5px; }
    .footer { text-align: center; padding: 20px; font-size: 12px; color: #666; }
  </style>
</head>
<body>
  <div class="container">
    <h1>Password Reset Request</h1>
    
    <p>Hi {{.UserName}},</p>
    
    <p>We received a request to reset your password for your {{.AppName}} account.</p>
    
    <div class="alert">
      <strong>⚠ Security Notice:</strong> If you didn't request this, please ignore this email.
    </div>
    
    <p style="text-align: center; margin: 30px 0;">
      <a href="{{.ResetLink}}" class="button">Reset Password</a>
    </p>
    
    <p><small>This link will expire in {{.ExpiryMinutes}} minutes.</small></p>
    
    <p>Or copy and paste this URL into your browser:</p>
    <p style="word-break: break-all; background: #f5f5f5; padding: 10px;">{{.ResetLink}}</p>
    
    <div class="footer">
      <p>{{.CompanyName}} Security Team</p>
    </div>
  </div>
</body>
</html>`

	resetData := map[string]interface{}{
		"AppName":       "MyApp",
		"UserName":      "Charlie",
		"ResetLink":     "https://example.com/reset/xyz789",
		"ExpiryMinutes": 30,
		"CompanyName":   "Tech Corp",
	}

	output, err = engine.RenderString(resetHTML, resetData)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(output)
}
