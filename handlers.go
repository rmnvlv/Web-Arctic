package main

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/gofiber/fiber/v2"
)

const (
	ErrorMessage   = "Some form fields are entered incorrectly. Please change them."
	SuccessMessage = "Thank you for registration for AMTC 2022!"

	hCaptchaAPIURL = "https://hcaptcha.com/siteverify"
)

var (
	hCaptchaSecretKey = os.Getenv("HCAPTCHA_SECRET_KEY")
	hCaptchaSiteKey   = os.Getenv("HCAPTCHA_SITE_KEY")
	ErrCaptchaEmpty   = errors.New("captcha is empty")
)

// Response is the hcaptcha JSON response.
type Response struct {
	ChallengeTS string   `json:"challenge_ts"`
	Hostname    string   `json:"hostname"`
	ErrorCodes  []string `json:"error-codes,omitempty"`
	Success     bool     `json:"success"`
	Credit      bool     `json:"credit,omitempty"`
}

func verifyCaptcha(captcha string) (bool, error) {
	if captcha == "" {
		return false, ErrCaptchaEmpty
	}

	form := url.Values{}
	form.Add("secret", hCaptchaSecretKey)
	form.Add("response", captcha)
	form.Add("sitekey", hCaptchaSiteKey)

	resp, err := http.DefaultClient.PostForm(hCaptchaAPIURL, form)
	if err != nil {
		return false, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	defer resp.Body.Close()

	var response Response
	if err := json.Unmarshal(body, &response); err != nil {
		return false, err
	}

	if !response.Success {
		return false, fmt.Errorf("hCaptcha: %v", response.ErrorCodes)
	}

	return true, nil
}

func registerNewParticipant(c *fiber.Ctx) error {
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
	}

	hCaptcha := c.FormValue("h-captcha-response")

	formErrors := make(map[string]string)

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
	val, err = regexp.MatchString(`[^@\s]+@[^@\s]+\.[^@\s]+$`, participant.Email)
	if err != nil || !val {
		formErrors["Email"] = "Wrong email format. Example: mail@example.com"
	}

	data := fiber.Map{}
	messages := make(map[string]string)

	if len(formErrors) > 0 {
		messages["Error"] = ErrorMessage
		data["Values"] = participant
	} else {
		DB.Create(&participant)
		messages["Success"] = SuccessMessage
	}

	data["Title"] = "Registration and submission"
	data["Links"] = Links
	data["Errors"] = formErrors
	data["Message"] = messages

	return c.Render("registration", data)
}

func downloadFile(c *fiber.Ctx) error {
	var participants []Participant

	DB.Find(&participants)

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

	// "ABCDEFGHI"
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
	}

	fileNameExcel := strconv.FormatInt(time.Now().Unix(), 10) + ".xlsx" //  "./" +

	if err = fileExcel.SaveAs(fileNameExcel); err != nil {
		return err
	}

	return c.SendFile("./" + fileNameExcel)
}

func UploadFile(c *fiber.Ctx) error {
	return nil
}
