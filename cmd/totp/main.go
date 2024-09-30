package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/pquerna/otp/totp"
)

func promptForPasscode() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter Passcode: ")
	text, _ := reader.ReadString('\n')
	return text
}

// generateKey generates TOTP secret key
func generateKey(issuer, email string) (string, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      issuer,
		AccountName: email,
	})
	if err != nil {
		return "", fmt.Errorf("generate: %w", err)
	}
	return key.Secret(), nil
}

// generateTOTP generates TOTP token based on the given secret
func generateTOTP(secret string, now time.Time) (string, error) {
	otp, err := totp.GenerateCode(secret, now)
	if err != nil {
		return "", fmt.Errorf("generate code: %w", err)
	}
	return otp, nil
}

func main() {
	totpSecret, err := generateKey("MyApp", "ville@testmail.com")
	if err != nil {
		fmt.Printf("Error generating key: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Generated secret key: %s\n", totpSecret)

	secret := promptForPasscode()

	// Generate an OTP using the secret
	otp, err := generateTOTP(secret, time.Now())
	if err != nil {
		fmt.Printf("Error generating OTP: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Generated OTP: %s\n", otp)

	if totp.Validate(otp, secret) {
		fmt.Println("Valid OTP")
	} else {
		fmt.Println("Invalid OTP")
	}
}
