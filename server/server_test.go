// юнит-тесты
package main

import (
	"encoding/json"
	"moex_etf/server/storage"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

// для целей тестирования бизнес-логики создаём заглушку хранилища
type stub int // тип данных не имеет значения

var securities []storage.Security // имитация хранилища данных

// *******************************
// Выполняем контракт на хранилище
// InitData инициализирует фейковое хранилище фейковыми данными
func (s *stub) InitData() (err error) {

	// добавив в хранилище-заглушку одну запись
	var security = storage.Security{
		ID:        "MSFT",
		Name:      "Microsoft",
		IssueDate: 1514764800, // 01/01/2018

	}

	var quote = storage.Quote{
		SecurityID: "MSFT",
		Num:        0,
		TimeStamp:  1514764800,
		Price:      100,
	}

	security.Quotes = append(security.Quotes, quote)

	securities = append(securities, security)

	return err
}

// Securities возвращает список бумаг с котировками
func (s *stub) Securities() (data []storage.Security, err error) {
	return securities, err
}

// контракт выполнен
// *****************

// подготавливаем тестовую среду - инициализируем данные
func TestMain(m *testing.M) {

	// присваиваем указатель на экземпляр хранилища-заглушки глобальной переменной хранилища
	db = new(stub)

	// инициализируем данные (ничем)
	db.InitData()

	// выполняем все тесты пакета
	os.Exit(m.Run())
}

// тестируем отдачу котировок
func TestSecuritiesHandler(t *testing.T) {

	// проверяем обработчик запроса котировок
	req, err := http.NewRequest(http.MethodGet, "/api/v1/securities", nil)
	if err != nil {
		t.Fatal(err)
	}

	// ResponseRecorder записывает ответ сервера
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(securitiesHandler)

	// вызываем обработчик и передаём ему запрос
	handler.ServeHTTP(rr, req)

	// проверяем HTTP-код ответа
	if rr.Code != http.StatusOK {
		t.Errorf("код неверен: получили %v, а хотели %v", rr.Code, http.StatusOK)
	}

	// десериализуем (раскодируем) ответ сервера из формата json в структуру данных
	var ss []storage.Security

	err = json.NewDecoder(rr.Body).Decode(&ss)
	if err != nil {
		t.Fatal(err)
	}

	// выведем данные на экран для наглядности
	t.Logf("\n%+v\n\n", ss)

}

// тестируем отдачу данных об инфляции
func TestInflationHandler(t *testing.T) {

	// проверяем обработчик запроса данных об инфляции
	req, err := http.NewRequest(http.MethodGet, "/api/v1/inflation", nil)
	if err != nil {
		t.Fatal(err)
	}

	// ResponseRecorder записывает ответ сервера
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(inflationHandler)

	// вызываем обработчик и передаём ему запрос
	handler.ServeHTTP(rr, req)

	// проверяем HTTP-код ответа
	if rr.Code != http.StatusOK {
		t.Errorf("код неверен: получили %v, а хотели %v", rr.Code, http.StatusOK)
	}

	// десериализуем (раскодируем) ответ сервера из формата json в структуру данных
	var infl []inflationType

	err = json.NewDecoder(rr.Body).Decode(&infl)
	if err != nil {
		t.Fatal(err)
	}

	// выведем данные на экран для наглядности
	t.Logf("\n%+v\n\n", infl)

}
