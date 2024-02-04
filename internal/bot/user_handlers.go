package bot

import (
	"errors"
	"fmt"
	"github.com/enescakir/emoji"
	"github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
	"regexp"
	"strings"
	"telegram-bot/internal/bot/internal/button"
	"telegram-bot/internal/bot/internal/template"
	"telegram-bot/internal/bot/internal/validator"
	"telegram-bot/internal/domain"
	"telegram-bot/internal/requester"
	"telegram-bot/internal/utils/formatter"
)

const (
	messageTag    = "Message"
	hoursTag      = "Hours"
	startDateTag  = "StartDate"
	endDateTag    = "EndDate"
	notApplicable = "N/A"
)

var tryAgainAlarmMessage = "Try again editing the form message or execute /setAlarm to start again"

// start this endpoint has two possible flows
// 1. If the user is registered, awaits for other commands
// 2. If is not registered, gives to the user the options to create an account
func (tb *TelegramBot) start(c tele.Context) error {
	senderInfo := c.Sender()
	if senderInfo == nil {
		_ = c.Send(errUserInfoNotFound.Error())
		return errUserInfoNotFound
	}

	isRegistered, userInfo, err := tb.IsUserRegistered(senderInfo.ID)

	if err != nil {
		return c.Send("Oops, something went wrong searching your info. Please try again")
	}

	if !isRegistered {
		button.Menu.Inline(
			button.Menu.Row(button.CreateAccount),
			button.Menu.Row(button.DontCreateAccount),
		)
		message := fmt.Sprintf("You are not registered in %s, "+
			"do you want to create an account now?",
			formatter.Bold(appName),
		)
		message += fmt.Sprintf(
			"\n\n%s If you don't have an account you will not be able to perform operations with %s",
			emoji.Eyes,
			formatter.Italic(botName),
		)
		return c.Send(message, button.Menu)
	}

	welcomeMessage := template.WelcomeMessage(userInfo.FullName)
	return c.Send(welcomeMessage)
}

// createAccount returns a URL to register for appName
func (tb *TelegramBot) createAccount(c tele.Context) error {
	senderInfo := c.Sender()
	if senderInfo == nil {
		_ = c.Send(errUserInfoNotFound.Error())
		return errUserInfoNotFound
	}

	signUpButton := button.SignUpButton(senderInfo.ID)
	message := fmt.Sprintf("Click below to sign up %s", emoji.BackhandIndexPointingDown)
	err := c.Send(message, signUpButton)
	if err != nil {
		return fmt.Errorf("%w: %w", errSendingSignUpLink, err)
	}

	afterCreationMessage := fmt.Sprintf("After creating the account perform /start again %s", emoji.GrinningCatWithSmilingEyes)
	return c.Send(afterCreationMessage)
}

// omitAccountCreation byd dude, good luck
func (tb *TelegramBot) omitAccountCreation(c tele.Context) error {
	p := &tele.Photo{File: tele.FromURL("https://pbs.twimg.com/media/FRxJVLYXwAAlGPk?format=jpg&name=small")}
	_, err := p.Send(tb.bot, c.Recipient(), nil)
	return err
}

// setAlarm sends a form to the user so the alarm can be registered
func (tb *TelegramBot) setAlarm(c tele.Context) error {
	senderInfo := c.Sender()
	if senderInfo == nil {
		_ = c.Send(errUserInfoNotFound.Error())
		return errUserInfoNotFound
	}

	alarmMenu := tb.bot.NewMarkup()
	helpButton := alarmMenu.Text("Click here to display the alarm form")

	alarmForm := fmt.Sprintf("%s\n\n", registerAlarmEndpoint)
	alarmForm += template.Alarm()

	helpButton.InlineQueryChat = alarmForm

	alarmMenu.Inline(
		alarmMenu.Row(helpButton),
	)

	return c.Send("Please, enter the information about the alarm", alarmMenu)
}

// IsUserRegistered returns three elements:
//
// + First: a boolean to know if the user is registered or not
//
// + Second: the user information
//
// + Third: an error if something occurs requesting the user information
func (tb *TelegramBot) IsUserRegistered(telegramID int64) (bool, domain.UserInfo, error) {
	userInfo, err := tb.requester.GetUserData(telegramID)

	var requestError requester.RequestError
	isRequestError := errors.As(err, &requestError)
	if isRequestError && requestError.IsNotFound() {
		logrus.Infof("user with telegramID %v not found", telegramID)
		return false, domain.UserInfo{}, nil
	}

	if err != nil {
		return false, domain.UserInfo{}, err
	}

	return true, userInfo, nil
}

// registerAlarm register an alarm for the user with the provided data
func (tb *TelegramBot) registerAlarm(c tele.Context) error {
	senderInfo := c.Sender()
	if senderInfo == nil {
		_ = c.Send(errUserInfoNotFound.Error())
		return errUserInfoNotFound
	}

	alarmData, err := extractAlarmData(c.Message().Text, messageTag, hoursTag, startDateTag, endDateTag)
	if err != nil && errors.Is(err, errInvalidForm) {
		return c.Send(fmt.Sprintf("%v Invalid form, you don't have to modify the structure, only the field values. %s",
			emoji.PoliceCarLight,
			tryAgainAlarmMessage,
		))
	}

	if err != nil && errors.Is(err, errMissingFormField) {
		return c.Send("%v %v. %s", emoji.PoliceCarLight, err, tryAgainAlarmMessage)
	}

	if len(alarmData[hoursTag]) < 5 {
		return c.Send(fmt.Sprintf("Message must have at least 5 characters. %s", tryAgainAlarmMessage))
	}

	for _, hour := range strings.Split(alarmData[hoursTag], ",") {
		if err := validator.ValidateHour(hour); err != nil {
			return c.Send(fmt.Sprintf("%v. %s", err, tryAgainAlarmMessage))
		}
	}

	if err := validator.ValidateDateType(alarmData[startDateTag]); err != nil {
		return c.Send(fmt.Sprintf("Invalid start date: format must be year/month/day. %s", tryAgainAlarmMessage))
	}

	if err := validator.ValidateDateType(alarmData[endDateTag]); alarmData[endDateTag] != notApplicable && err != nil {
		return c.Send(fmt.Sprintf("Invalid end date: format must be year/month/day. %s", tryAgainAlarmMessage))
	}

	// ToDo: add request to ticker service. Licha

	return c.Send("Your alarm was set correctly")
}

// extractAlarmData extracts alarm data from the given message. Does not validate the fields, it only ensures that they are all present
func extractAlarmData(alarmDataRaw string, fields ...string) (map[string]string, error) {
	regex := regexp.MustCompile(`Message:\s*(?P<Message>([^\n]*))\s*Hours:\s*(?P<Hours>[^\n]*)\s+Start Date:\s*(?P<StartDate>([^\n]*))\s+End Date:\s*(?P<EndDate>([^\n]*|N/A))`)
	match := regex.FindStringSubmatch(alarmDataRaw)
	if match == nil {
		return nil, fmt.Errorf("%w", errInvalidForm)
	}

	// groupName are capture from the regex expression
	petData := make(map[string]string)
	for idx, groupName := range regex.SubexpNames() {
		if idx != 0 && groupName != "" {
			petData[groupName] = strings.TrimRight(match[idx], " ")
		}
	}

	for _, field := range fields {
		if _, found := petData[field]; !found {
			return nil, fmt.Errorf("%w: %s", errMissingFormField, field)
		}
	}

	return petData, nil
}
