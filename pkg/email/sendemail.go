package email

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"net/smtp"
)

func SenSecretCode(to []string) (string, error) {
	from := "diyordev3@gmail.com"
	password := "uvcy ksvb yxyy gsmi"

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	randInt, err := rand.Int(rand.Reader, big.NewInt(100000000))
	if err != nil {
		return "", err
	}

	subject := "Subject: Your Secret Code\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Secret Code</title>
    <style>
        body {
            font-family: 'Arial', sans-serif;
            background-color: #2c3e50;
            display: flex;
            justify-content: center;
            align-items: center;
            height: 100vh;
            margin: 0;
        }
        .container {
            background-color: #ecf0f1;
            padding: 20px;
            border-radius: 10px;
            box-shadow: 0 0 20px rgba(0, 0, 0, 0.1);
            text-align: center;
            width: 300px;
        }
        .container h1 {
            color: #3498db;
            font-size: 24px;
            margin-bottom: 10px;
        }
        .container p {
            color: #34495e;
            font-size: 18px;
            margin: 20px 0;
        }
        .container .code {
            font-size: 32px;
            color: #e74c3c;
            font-weight: bold;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>Smart Home Control System</h1>
        <p>Your secret code is:</p>
        <div class="code">` + fmt.Sprintf("%d", randInt) + `</div>
    </div>
</body>
</html>`

	message := []byte(subject + mime + body)

	auth := smtp.PlainAuth("", from, password, smtpHost)

	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%d", randInt), nil
}
