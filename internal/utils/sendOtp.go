package utils

import (
	"dubai-auto/internal/config"
	"fmt"
	"log"

	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

func SendOtp(toNumber string, otp int) {

	// Twilio credentials (use env variables in production)
	accountSID := config.ENV.TWILIO_ACCOUNT_SID
	authToken := config.ENV.TWILIO_AUTH_TOKEN
	fromNumber := config.ENV.TWILIO_PHONE_NUMBER // e.g. "+15076273628"

	// Recipient number

	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSID,
		Password: authToken,
	})

	message := fmt.Sprintf("Your OTP code is %d. It expires in 5 minutes.", otp)

	params := &openapi.CreateMessageParams{}
	params.SetTo(toNumber)
	params.SetFrom(fromNumber)
	params.SetBody(message)

	resp, err := client.Api.CreateMessage(params)

	if err != nil {
		log.Fatal("Error sending SMS:", err)
	}

	fmt.Println("OTP Sent!")

	if resp.Sid != nil {
		fmt.Println("Message SID:", *resp.Sid)
	}
}
