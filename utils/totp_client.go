package utils

import (
	"github.com/skip2/go-qrcode"
	"github.com/xlzd/gotp"
	"log"
	"time"
	"darkoo/apperrors"
)


const (
	TotpSecretLength int = 16

	DarkRoom = "DarkRoom"

	QRCodeSize = 256
)

func GenerateTOTPSecret() string {
	return gotp.RandomSecret(TotpSecretLength)
}


func GenerateTOTPQRCode(totpSecret, userEmail string) ([]byte, error) {
	totpClient := gotp.NewDefaultTOTP(totpSecret)

	uri := totpClient.ProvisioningUri(userEmail, DarkRoom)

	qrcode, err := qrcode.Encode(uri, qrcode.Medium, QRCodeSize)
	if err != nil {
		log.Printf("error generating totp QR code: %v\n", err)
		return nil, apperrors.NewInternalWithMessage("Unable to generate TOTP QR code. Please try again.")
	}

	return qrcode, nil
}


func VerifyTOTP(totpSecret, userOtp string) bool {
	totpClient := gotp.NewDefaultTOTP(totpSecret)
	return totpClient.Verify(userOtp, time.Now().Unix())
}


