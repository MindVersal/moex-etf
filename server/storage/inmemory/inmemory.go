// Package inmemory реализует хранение данных в памяти
package inmemory

import (
	"encoding/json"
	"fmt"
	"moex_etf/server/storage"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// объект в памяти, хранящий данные котировок
var securities []storage.Security

type jsonSecurities struct {
	Securities struct {
		Data [][]string `json:"data"`
	} `json:"securities"`
}

type jsonQuote struct {
	History struct {
		Data [][]interface{} `json:"data"`
	} `json:"history"`
}

// Storage - тип данных хранилища
type Storage struct {
	Name string
}

// New создаёт экземпляр хранилища и возвращает ссылку на него
func New() storage.Interface {
	s := &Storage{}
	s.Name = "InMemory Storage Backend"
	return s
}

// Securities возвращает список бумаг с котировками
func (s *Storage) Securities() (data []storage.Security, err error) {
	return securities, err
}

// Add добавляет бумагу в список
func (s *Storage) Add(item storage.Security) (err error) {
	securities = append(securities, item)
	return err
}

// InitData инициализирует хранилище данными с сервера Мосбиржи
func (s *Storage) InitData() (err error) {

	securities, err := getSecurities()
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	wg.Add(len(securities))

	for _, security := range securities {

		go func(item storage.Security) {

			defer wg.Done()

			var quotes []storage.Quote
			quotes, err = getSecurityQuotes(item)
			if err != nil {
				fmt.Println(item, err)
				return
			}

			item.Quotes = quotes

			err = s.Add(item)
			if err != nil {
				return
			}

		}(security)

	}

	wg.Wait()

	return err

}

// получаем список бумаг биржевых фондов на Мосбирже
func getSecurities() (data []storage.Security, err error) {

	const securitiesListURL = "https://iss.moex.com/iss/securitygroups/stock_etf/collections/stock_etf_all/securities.json"

	var list jsonSecurities
	var tm time.Time

	resp, err := http.Get(securitiesListURL)
	if err != nil {
		return data, err
	}

	err = json.NewDecoder(resp.Body).Decode(&list)
	if err != nil {
		return data, err
	}

	for _, item := range list.Securities.Data {
		var s storage.Security
		s.ID = item[0]
		s.Name = item[2]
		tm, err = time.Parse("2006-01-02", item[5])
		if err != nil {
			return data, err
		}
		s.IssueDate = tm.Unix()

		data = append(data, s)

	}

	return data, err

}

// получаем котировки для бумаги
func getSecurityQuotes(security storage.Security) (quotes []storage.Quote, err error) {

	var year, month, day int

	// начало измерения котировок
	start := time.Unix(security.IssueDate, 0)
	// устанавливаем первое число месяца
	start = start.AddDate(0, 0, -(start.Day() - 1))
	now := time.Now()

	counter := 0 // номер измерения

	for d := start; d.Unix() < now.Unix(); d = d.AddDate(0, 1, 0) {

		var q storage.Quote
		var jq jsonQuote
		var tm time.Time

		q.Num = counter

		year = d.Year()
		month = int(d.Month())

		// находим последний рабочий день месяца - последний день торгов
		curDate := d
		curMonth := curDate.Month()

		for {
			if int(curDate.Month()) != int(curMonth) {
				break
			}

			if int(curDate.Weekday()) != 0 && int(curDate.Weekday()) != 6 {
				day = curDate.Day()
			}

			curDate = curDate.AddDate(0, 0, 1)
		}

		var quotesURL = "https://iss.moex.com/iss/history/engines/stock/markets/shares/securities/" +
			security.ID + ".json?from=" + strconv.Itoa(year) + "-" + strconv.Itoa(month) + "-" + strconv.Itoa(day) +
			"&till=" + strconv.Itoa(year) + "-" + strconv.Itoa(month) + "-" + strconv.Itoa(day)

		resp, err := http.Get(quotesURL)
		if err != nil {
			return quotes, err
		}

		err = json.NewDecoder(resp.Body).Decode(&jq)
		if err != nil {
			fmt.Println("quote decoding error:", err)
			continue
		}

		if len(jq.History.Data) < 1 || len(jq.History.Data[0]) < 10 {
			continue
		}

		lastIndex := len(jq.History.Data) - 1

		q.Num = counter
		q.SecurityID = security.ID
		tm, err = time.Parse("2006-01-02", jq.History.Data[lastIndex][1].(string))
		q.TimeStamp = tm.Unix()
		if err != nil {
			fmt.Println("time parsing error:", err)
			continue
		}

		// legal close price
		q.Price = jq.History.Data[lastIndex][9].(float64)

		quotes = append(quotes, q)

		counter++

	}

	return quotes, err
}
