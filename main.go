package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Matrix struct {
	Birght_year  int
	Birght_month int
	Birght_day   int
	// Years_number: sum of digits in year
	Years_number int
	// N3: Sum of digits in Birght_year
	N3 int
	// N4: Birght_day + Birght_month
	N4 int
	// N5: N2 + N3
	N5 int
	// N6: N4 + N5
	N6 int
	// N7: N4 + Years_number
	N7 int
	// N8: Plot + Lesson
	N8 int
	// N9: N7 + N8
	N9 int
	// N10: N5 + Years_number
	N10 int
	// N11: Lesson + Exam
	N11 int
	// N12: N10 + N11
	N12 int
	// Plot: Birght_day + Years_number
	Plot int
	// Lesson: Birght_month + Years_number
	Lesson int
	// Exam: N3 + Years_number
	Exam int
	// Star: N6 + Years_number
	Star int
	// Help: Lesson + Exam
	Help int
}

func sumDigitsInNumber(num int) int {
	res := 0
	for num > 0 {
		res += num % 10
		num /= 10
	}
	return res
}

func reduceUnder22(num int) int {
	for num > 22 {
		num = num - 22
	}
	return num
}

func (matrix *Matrix) calc() {
	year := time.Now().Year()
	matrix.Years_number = reduceUnder22(sumDigitsInNumber(year))
	matrix.N3 = reduceUnder22(sumDigitsInNumber(matrix.Birght_year))
	matrix.Plot = reduceUnder22(matrix.Birght_day + matrix.Years_number)
	matrix.Lesson = reduceUnder22(matrix.Birght_month + matrix.Years_number)
	matrix.Exam = reduceUnder22(matrix.N3 + matrix.Years_number)
	matrix.N4 = reduceUnder22(matrix.Birght_day + matrix.Birght_month)
	matrix.N5 = reduceUnder22(matrix.Birght_month + matrix.N3)
	matrix.N6 = reduceUnder22(matrix.N4 + matrix.N5)
	matrix.Star = reduceUnder22(matrix.N6 + matrix.Years_number)
	matrix.Help = reduceUnder22(matrix.Lesson + matrix.Star)
	matrix.N7 = reduceUnder22(matrix.N4 + matrix.Years_number)
	matrix.N8 = reduceUnder22(matrix.Plot + matrix.Lesson)
	matrix.N9 = reduceUnder22(matrix.N7 + matrix.N8)
	matrix.N10 = reduceUnder22(matrix.N5 + matrix.Years_number)
	matrix.N11 = reduceUnder22(matrix.Lesson + matrix.Exam)
	matrix.N12 = reduceUnder22(matrix.N10 + matrix.N11)
}

func (matrix *Matrix) Process(birghtday string) {
	re := regexp.MustCompile(`(?P<day>\d+)-(?P<month>\d+)-(?P<year>\d+)`)
	r2 := re.FindAllStringSubmatch(birghtday, -1)[0]

	year, _ := strconv.Atoi(r2[len(r2)-1])
	month, _ := strconv.Atoi(r2[len(r2)-2])
	day, _ := strconv.Atoi(r2[len(r2)-3])

	matrix.Birght_year = year
	matrix.Birght_month = month
	matrix.Birght_day = day

	matrix.calc()
}

func main() {
	birghtday := flag.String("b", "", "birghtday [26-04-1985]")
	daemon := flag.Bool("d", false, "run as daemon")
	flag.Parse()
	//fmt.Println(*birghtday)
	if *daemon {
		runBot()
		os.Exit(0)
	}
	if *birghtday == "" {
		flag.Usage()
		os.Exit(1)
	}

	matrix := Matrix{}
	matrix.Process(*birghtday)

	fmt.Printf("%+v\n", matrix)

}

func runBot() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("MATRIX_TG_API_TOKEN"))
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			re := regexp.MustCompile(`(?P<day>\d+)-(?P<month>\d+)-(?P<year>\d+)`)
			if !re.Match([]byte(update.Message.Text)) {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Please use this form 31-02-1985")
				bot.Send(msg)
				continue
			}
			matrix := Matrix{}
			matrix.Process(update.Message.Text)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("%+v\n", matrix))
			msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)
		}
	}
}
