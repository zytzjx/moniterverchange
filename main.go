package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/go-resty/resty/v2"
	"github.com/gookit/color"
	"github.com/joho/godotenv"
	"howett.net/plist"

	gomail "gopkg.in/mail.v2"
)

func SendReport(v1, v2 string) {
	m := gomail.NewMessage()

	// Set E-Mail sender
	m.SetHeader("From", "jefferyz#futuredial.com")

	// Set E-Mail receivers
	m.SetHeader("To", "jefferyz#futuredial.com")

	// Set E-Mail subject
	m.SetHeader("Subject", "iTunes updated")

	// Set E-Mail body. You can set plain text or html with text/html
	m.SetBody("text/plain", fmt.Sprintf("iTunes version has updated from %s to %s\n", v1, v2))

	// Settings for SMTP server
	d := gomail.NewDialer("smtp.office365.com", 587, "jefferyz#futuredial.com", "passwordinput")

	// This is only needed when SSL/TLS certificate is not valid on server.
	// In production this should be set to false.
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Now send E-Mail
	if err := d.DialAndSend(m); err != nil {
		color.Redf("send mail error: %v\n", err)
	}

}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	iTunesVer := os.Getenv("version")
	color.Greenln("current Version: " + iTunesVer)

	//color.Yellowln(string(resp.Body()))

	s := gocron.NewScheduler(time.UTC)
	s.Every(1).Day().At("7:30").Do(func() {
		client := resty.New()
		resp, err := client.R().Get("https://s.mzstatic.com/version?machineID=c2d9cd28d8c9349c")
		if err != nil {
			color.Redf("error: %v\n", err)
			return
		}
		var data map[string]interface{}
		buf := bytes.NewReader(resp.Body())
		decoder := plist.NewDecoder(buf)
		err = decoder.Decode(&data)
		//_, err = plist.Unmarshal(resp.Body(), data)
		if err != nil {
			color.Redf("error: %v\n", err)
			return
		}

		if ver, ok := data["iTunesWindows10Version"]; ok {
			color.Cyanf("server version: %v\n", ver)
			vers := ver.(string)
			if iTunesVer != vers {
				SendReport(iTunesVer, vers)
				iTunesVer = vers
				os.Setenv("version", iTunesVer)
				//Send Mail
			}
		}
	})

	s.StartBlocking()
}
