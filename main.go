package main

import (
	"fmt"
	"log"
	"os"

	"errors"
	"github.com/sendgrid/rest"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"net/http"
)

func main() {

	createServer()

}

func createServer() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `
	<html>
		<body>
			<h1>Welcome to SendGrid Test App</h1>
		</body>
	</html>`)
	})

	http.HandleFunc("/email", handleEmail)

	port := os.Getenv("HTTP_PLATFORM_PORT")

	if port == "" {
		port = "8082"
	} else {
		f, err := os.OpenFile("D:\\home\\site\\wwwroot\\testlogfile",
			os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)

		if err != nil {
			log.Fatalf("error opening log file: %v", err)
		}

		defer f.Close()

		log.SetOutput(f)
	}

	log.Printf("port is: %v", port)
	log.Fatal(http.ListenAndServe(":" + port, nil))
}

func handleEmail(w http.ResponseWriter, r *http.Request) {

	//validations
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query()
	to := query.Get("to")
	if to == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "request must include 'to' query param with valid email address")
		return
	}
	subject := query.Get("subject")
	message := query.Get("message")

	log.Printf("sending email. to: %v, subject: %v, message: %v", to, subject, message)
	response, err := sendEmail(to, subject, message)
	if err != nil {
		log.Println(err)
		w.WriteHeader(response.StatusCode)
		fmt.Fprint(w, err)
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
		w.WriteHeader(response.StatusCode)

		for k, v := range response.Headers {
			for _, sv := range v {
				w.Header().Add(k, sv)
			}
		}

		fmt.Fprint(w, response.Body)
	}
}
func sendEmail(to string, subject string, message string) (*rest.Response, error) {
	key := os.Getenv("SENDGRID_API_KEY")
	if key == "" {
		return &rest.Response{}, errors.New("missing api key")
	}
	fromEmail := mail.NewEmail("Example User", "test@example.com")
	sbj := "Sending with SendGrid is Fun"
	if subject != "" {
		sbj = subject
	}
	toEmail := mail.NewEmail("Example User", to)
	plainTextContent := "and easy to do anywhere, even with Go"
	htmlContent := "<strong>and easy to do anywhere, even with Go</strong>"

	if message != "" {
		plainTextContent = message
		htmlContent = fmt.Sprintf("<strong>%v</strong>", message)
	}

	msg := mail.NewSingleEmail(fromEmail, sbj, toEmail, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(key)
	response, err := client.Send(msg)

	return response, err
}
