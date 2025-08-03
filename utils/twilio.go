package utils

import (
	"log"

	"github.com/twilio/twilio-go"
	verify "github.com/twilio/twilio-go/rest/verify/v2"
)

var (
	accountSID = GetEnv("TWILIO_ACCOUNT_SID", "")
	authToken  = GetEnv("TWILIO_AUTH_TOKEN", "")
	serviceSID = GetEnv("TWILIO_VERIFY_SERVICE_SID", "")
	client     = twilio.NewRestClientWithParams(twilio.ClientParams{Username: accountSID, Password: authToken})
)

func SendOTP(phone string) error {
	params := &verify.CreateVerificationParams{}
	params.SetTo(phone)
	params.SetChannel("sms")

	_, err := client.VerifyV2.CreateVerification(serviceSID, params)
	if err != nil {
		log.Printf("Failed to send verification code: %v", err)
		return err
	}

	log.Printf("OTP sent to %s", phone)
	return nil
}

func VerifyOTP(phone, code string) (bool, error) {
	params := &verify.CreateVerificationCheckParams{
		To:   &phone,
		Code: &code,
	}

	resp, err := client.VerifyV2.CreateVerificationCheck(serviceSID, params)
	if err != nil {
		log.Printf("OTP verification failed: %v", err)
		return false, err
	}

	return resp.Status != nil && *resp.Status == "approved", nil
}
