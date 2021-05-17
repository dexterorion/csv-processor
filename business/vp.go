package business

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/csv-processor/model"
	"github.com/csv-processor/mongo"
	"go.uber.org/zap"
)

const (
	VALID = iota
	DEVIATION
	INVALID
)

type VP interface {
	Process(ctx context.Context) error
}

type vpImpl struct {
	dbAcess  *mongo.DB
	reader   *csv.Reader
	filetype string
	parking  model.Parking
	cn       chan *model.Line
	logger   *zap.Logger
}

func NewVP(dbAcess *mongo.DB, reader *csv.Reader, filetype string, parking model.Parking) VP {
	log, _ := zap.NewProduction()

	service := &vpImpl{
		dbAcess:  dbAcess,
		reader:   reader,
		filetype: filetype,
		parking:  parking,
		cn:       make(chan *model.Line),
		logger:   log,
	}

	go service.processLine()

	return service
}

func (s *vpImpl) Process(ctx context.Context) error {
	switch s.filetype {
	case "transactions":
		return s.transactionsProcess(ctx)
	default:
		return fmt.Errorf("filetype [%s] does not exists for parking", s.filetype)
	}
}

func (s *vpImpl) transactionsProcess(ctx context.Context) error {

	for {
		line, err := s.reader.Read()

		if err != nil {
			if err == io.EOF {
				s.logger.Info("Done")
				return nil
			} else {
				return err
			}
		}

		cin, parsed := parseDate(line[6])
		if !parsed {
			fmt.Println(line[6])
			continue
		}

		cout, parsed := parseDate(line[7])
		if !parsed {
			fmt.Println(line[7])
			continue
		}

		paid, _ := strconv.ParseFloat(line[10], 64)

		data := &model.Line{
			Unit:          line[0],
			Ticket:        line[1],
			Identity:      line[2],
			Matricula:     line[3],
			UseType:       line[4],
			CheckIn:       cin,
			CheckOut:      cout,
			Duration:      int64(cout.Sub(cin).Minutes()),
			PaidValue:     paid,
			PaymentMethod: line[11],
			Table:         line[12],
		}

		s.cn <- data
	}
}

func (s *vpImpl) processLine() {
	s.logger.Info("Starting process line")
	for {
		line := <-s.cn

		ci := line.CheckIn
		co := line.CheckOut

		transaction := &model.Transaction{
			CheckinDate:   ci,
			CheckoutDate:  co,
			Sequence:      line.Ticket,
			FareAmount:    line.PaidValue,
			PaidAmount:    line.PaidValue,
			Matricula:     line.Matricula,
			IsValid:       true,
			UseType:       getUseType(line.Table),
			OfferType:     "On-demand",
			PaymentMethod: getPaymentMethod(line.PaymentMethod),
		}

		transaction.TimeIntervalHour = []time.Time{}
		transaction.Duration = 0

		co = time.Date(co.Year(), co.Month(), co.Day(), co.Hour(), 0, 0, 0, co.Location())

		if ci.Minute() > 0 || ci.Second() > 0 {
			ci = time.Date(ci.Year(), ci.Month(), ci.Day(), ci.Hour()+1, 0, 0, 0, ci.Location())
		}

		for {
			if ci.After(co) {
				break
			}

			transaction.Duration = transaction.Duration + 1
			transaction.TimeIntervalHour = append(transaction.TimeIntervalHour, ci)
			ci = ci.Add(1 * time.Hour)
		}

		s.dbAcess.TransactionCollection.Create(context.Background(), transaction)
	}
}

func getUseType(value string) string {
	if value == "Mensalista" {
		return "Mensalista"
	}

	return "Avulso"
}

func getPaymentMethod(value string) string {
	switch value {
	case "CC":
		return "Creditcard"
	case "CD":
		return "Debitcard"
	case "DI":
		return "Dinheiro"
	case "CA":
		return "Cancelado"
	case "IS":
		return "Isento"
	case "CH":
		return "Cheque"
	default:
		return "Outro"
	}
}

func parseDate(dt string) (date time.Time, parsed bool) {
	d := strings.Split(dt, " ")
	if len(d) != 2 {
		return time.Now(), false
	}

	days := strings.Split(d[0], "/")
	if len(days) != 3 {
		return time.Now(), false
	}
	day, _ := strconv.Atoi(days[1])
	month, _ := strconv.Atoi(days[0])
	year, _ := strconv.Atoi(days[2])

	hours := strings.Split(d[1], ":")
	if len(hours) != 2 {
		return time.Now(), false
	}

	hour, _ := strconv.Atoi(hours[0])
	min, _ := strconv.Atoi(hours[1])

	return time.Date(year+2000, time.Month(month), day, hour, min, 0, 0, time.UTC), true
}
