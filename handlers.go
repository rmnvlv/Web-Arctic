package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"net/mail"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	emailverifier "github.com/AfterShip/email-verifier"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const (
	ErrorMessage   = "Some form fields are entered incorrectly. Please change them."
	SuccessMessage = "Thank you for registration for AMTC 2022!"
)

// var mailQueue []Participant = nil

func (a *App) registerNewParticipant(c *fiber.Ctx) error {
	participant := Participant{
		Surname:             c.FormValue("surname"),
		Name:                c.FormValue("name"),
		Organization:        c.FormValue("organization"),
		Position:            c.FormValue("position"),
		Phone:               c.FormValue("phone"),
		Email:               c.FormValue("email"),
		PresentationForm:    c.FormValue("presentation-form"),
		PresentationSection: c.FormValue("presentation-section"),
		PresentationTitle:   c.FormValue("presentation-title"),
		Code:                "",
		Attempts:            3,
	}

	formErrors := make(map[string]string)

	if os.Getenv("APP_ENV") == "prod" {
		hCaptcha := c.FormValue("h-captcha-response")

		if ok, err := verifyCaptcha(hCaptcha); !ok {
			if errors.Is(err, ErrCaptchaEmpty) {
				formErrors["Captcha"] = "Сaptcha is not passed"
				// formError.Captcha = "Сaptcha is not passed"
			} else {
				// formError.Captcha = "Please try again"
				formErrors["Captcha"] = "Please try again"
			}
			log.Println(err)
		}
	}

	// Validate phone
	val, err := regexp.MatchString(`^((8|\+7)[\- ]?)?(\(?\d{3}\)?[\- ]?)?[\d\- ]{7,10}$`, participant.Phone)
	if err != nil && participant.Phone != "" || !val && participant.Phone != "" {
		formErrors["Phone"] = "Phone number should be valid format."
	}
	//Validate surname
	val, err = regexp.MatchString(`^[a-zA-Z]+$`, participant.Surname)
	if err != nil || !val {
		formErrors["Surname"] = "Surname can only be a-zA-Z."
	}
	//validate name
	val, err = regexp.MatchString(`^[a-zA-Z]+$`, participant.Name)
	if err != nil || !val {
		formErrors["Name"] = "Name can only be a-zA-Z."
	}
	//validate email
	if _, err := mail.ParseAddress(participant.Email); err != nil {
		formErrors["Email"] = "Wrong email format. Example: mail@example.com"
	}
	// val, err = regexp.MatchString(`[^@\s]+@[^@\s]+\.[^@\s]+$`, participant.Email)
	// if err != nil || !val {
	// 	formErrors["Email"] = "Wrong email format. Example: mail@example.com"
	// }
	verifier := emailverifier.NewVerifier()
	if _, err = verifier.Verify(participant.Email); err != nil {
		formErrors["Email"] = "Email does not exists."
	}

	data := fiber.Map{}
	messages := make(map[string]string)

	participant.Code = uuid.New().String()

	if len(formErrors) > 0 {
		messages["Error"] = ErrorMessage
		data["Values"] = participant
	} else {
		a.db.Create(&participant)

		if err := a.sendEmail(
			To{strings.Join([]string{participant.Name, participant.Surname}, " "), participant.Email},
			Message{EmailSubject, EmailRegistrationTemplate},
		); err != nil {
			// log
			fmt.Println(err)
		}

		// mailQueue = append(mailQueue, participant)

		// ch := make(chan []Participant, 1)
		// go worker(ch)
		// ch <- mailQueue
		// close(ch)

		messages["Success"] = SuccessMessage

	}

	data["Title"] = "Registration and submission"
	data["Links"] = Links
	data["Errors"] = formErrors
	data["Message"] = messages

	return c.Render("registration", data)
}

func (a *App) downloadFile(c *fiber.Ctx) error {
	var participants []Participant

	a.db.Find(&participants)

	fileName := "./" + strconv.FormatInt(time.Now().Unix(), 10) + ".xlsx"
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)

	headers := []string{
		"Name",
		"Surname",
		"Organization",
		"Position",
		"Phone",
		"Email",
		"Presentation Form",
		"Presentation Section",
		"Presentation Title",
		"Code",
	}

	if err := writer.Write(headers); err != nil {
		panic(err)
	}

	for _, participan := range participants {
		row := []string{
			participan.Name,
			participan.Surname,
			participan.Organization,
			participan.Position,
			participan.Phone,
			participan.Email,
			participan.PresentationForm,
			participan.PresentationSection,
			participan.PresentationTitle,
			participan.Code,
		}
		if err := writer.Write(row); err != nil {
			panic(err)
		}
	}

	writer.Flush()

	fileExcel := excelize.NewFile()

	_ = fileExcel.NewSheet("Sheet1")

	fileExcel.SetCellValue("Sheet1", "A1", "Name")
	fileExcel.SetCellValue("Sheet1", "B1", "Surname")
	fileExcel.SetCellValue("Sheet1", "C1", "Organization")
	fileExcel.SetCellValue("Sheet1", "D1", "Position")
	fileExcel.SetCellValue("Sheet1", "E1", "Phone")
	fileExcel.SetCellValue("Sheet1", "F1", "Email")
	fileExcel.SetCellValue("Sheet1", "G1", "PresentationForm")
	fileExcel.SetCellValue("Sheet1", "H1", "PresentationSection")
	fileExcel.SetCellValue("Sheet1", "I1", "PresentationTitle")
	fileExcel.SetCellValue("Sheet1", "J1", "Unique code")

	// "ABCDEFGHIJ"
	for counter, participan := range participants {
		number := strconv.Itoa(counter + 2)
		fileExcel.SetCellValue("Sheet1", "A"+number, participan.Name)
		fileExcel.SetCellValue("Sheet1", "B"+number, participan.Surname)
		fileExcel.SetCellValue("Sheet1", "C"+number, participan.Organization)
		fileExcel.SetCellValue("Sheet1", "D"+number, participan.Position)
		fileExcel.SetCellValue("Sheet1", "E"+number, participan.Phone)
		fileExcel.SetCellValue("Sheet1", "F"+number, participan.Email)
		fileExcel.SetCellValue("Sheet1", "G"+number, participan.PresentationForm)
		fileExcel.SetCellValue("Sheet1", "H"+number, participan.PresentationSection)
		fileExcel.SetCellValue("Sheet1", "I"+number, participan.PresentationTitle)
		fileExcel.SetCellValue("Sheet1", "J"+number, participan.Code)
	}

	fileNameExcel := strconv.FormatInt(time.Now().Unix(), 10) + ".xlsx" //  "./" +

	if err = fileExcel.SaveAs(fileNameExcel); err != nil {
		return err
	}

	return c.SendFile("./" + fileNameExcel)
}

func (a *App) mainView(c *fiber.Ctx) error {
	data := IndexPage
	data["Links"] = Links
	data["Header"] = true
	return c.Render("index", data)
}

func (a *App) programOverviewView(c *fiber.Ctx) error {
	data := fiber.Map{}
	data["Links"] = Links
	data["Title"] = "Programme Overview"
	return c.Render("programm-overview", data)
}

func (a *App) keynoteSpeakersView(c *fiber.Ctx) error {
	data := fiber.Map{}
	data["Title"] = "Keynote Speakers"
	data["Links"] = Links
	data["Content"] = "Key speakers to be determined later."
	return c.Render("basic", data)
}

func (a *App) requirementsView(c *fiber.Ctx) error {
	data := fiber.Map{}
	data["Title"] = "Requirements"
	data["Links"] = Links
	data["Content"] = "Article template will be posted later."
	return c.Render("basic", data)
}

func (a *App) generalInfoView(c *fiber.Ctx) error {
	data := fiber.Map{}
	data["Links"] = Links
	return c.Render("general-information", data)
}

func (a *App) registrationView(c *fiber.Ctx) error {
	data := fiber.Map{}
	data["Title"] = "Registration and submission"
	data["Links"] = Links
	return c.Render("registration", data)
}

func (a *App) adminView(c *fiber.Ctx) error {
	var participants []Participant
	a.db.Find(&participants)
	data := fiber.Map{}
	data["Users"] = participants
	data["Title"] = "Admin"
	data["Links"] = Links
	data["Content"] = "Admin panel"
	return c.Render("admin", data)
}

func (a *App) uploadView(c *fiber.Ctx) error {
	code := c.Query("code")
	fmt.Println(code)

	if code == "" {
		return c.Redirect("not-found")
	}

	var person Participant
	result := a.db.First(&person, "code = ?", code)

	if result.Error != nil {
		log.Fatal(result.Error)
	}

	data := fiber.Map{}
	data["Title"] = "Upload"
	data["User"] = person
	return c.Render("upload", data)
}

func (a *App) uploadArticleOrTezisi(c *fiber.Ctx) error {
	article, err := c.FormFile("article")
	if err != nil {
		return err
	}

	articleFile, err := article.Open()
	if err != nil {
		log.Fatal(err)
	}

	defer articleFile.Close()

	if err := saveToYandexDisk(articleFile, "Test/"+article.Filename); err != nil {
		log.Default().Panicln(err)
		return err
	}

	// disable upload files

	return c.Render("upload", fiber.Map{})
}

func (a *App) mailing(c *fiber.Ctx) error {
	var participants []Participant
	a.db.Find(&participants)
	for _, participant := range participants {
		if participant.PresentationForm == "Speaker" || participant.PresentationForm == "Publication" {
			//worker <-- chanel <-- participants ?
			break
		}
	}
	return nil
}

func (a *App) notFoundView(c *fiber.Ctx) error {
	data := fiber.Map{}
	data["Title"] = "Page Not Found"
	data["Links"] = Links
	data["Content"] = "Page Not Found"
	return c.Render("basic", data)
}

// func worker(ch <-chan []Participant) {
// 	//participan := Participant{}
// 	mailQueue := <-ch
// 	for len(mailQueue) > 0 {
// 		participan := mailQueue[0]
// 		if len(mailQueue) < 2 {
// 			mailQueue = nil
// 		} else {
// 			mailQueue = mailQueue[1:]
// 		}
// 		go func() {
// 			err := SendMail(To{participan.Name, participan.Email}, Message{EmailSubject, EmailTemplate})
// 			if err != nil && participan.Attempts > 0 {
// 				participan.Attempts--
// 				mailQueue = append(mailQueue, participan)
// 			}
// 		}()
// 	}
// }
