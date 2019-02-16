package inmemory

import (
	"moex_etf/server/storage"
	"reflect"
	"testing"
)

// тестируем конструктор хранилища
func TestNew(t *testing.T) {

	// вызываем тестируемый метод-фабрику
	memoryStorage := New()
	// создаём переменную для сравнения
	var s *Storage

	// сравниваем тип результата вызова функции с типом в модуле. просто так
	if reflect.TypeOf(memoryStorage) != reflect.TypeOf(s) {
		t.Errorf("тип неверен: получили %v, а хотели %v", reflect.TypeOf(memoryStorage), reflect.TypeOf(s))
	}

	// для наглядности выводим результат
	t.Logf("\n%+v\n\n", memoryStorage)

}

// проверяем отдачу котировок
func TestSecurities(t *testing.T) {

	// экземпляр хранилища в памяти
	var s *Storage

	// вызываем тестируемый метод
	ss, err := s.Securities()
	if err != nil {
		t.Error(err)
	}

	// для наглядности выводим результат
	t.Logf("\n%+v\n\n", ss)

}

// проверяем добавление котировки
func TestAdd(t *testing.T) {

	// экземпляр хранилища в памяти
	var s *Storage

	var security = storage.Security{
		ID: "MSFT",
	}

	var tt = []struct {
		s      storage.Security // добавляемая бумага
		length int              // длина массива (среза)
	}{
		{
			s:      security,
			length: 1,
		},
		{
			s:      security,
			length: 2,
		},
	}

	var ss []storage.Security

	// tc - test case, tt - table tests
	for _, tc := range tt {

		// вызываем тестируемый метод
		err := s.Add(security)
		if err != nil {
			t.Error(err)
		}

		ss, err = s.Securities()
		if err != nil {
			t.Error(err)
		}

		if len(ss) != tc.length {
			t.Errorf("невереная длина среза: получили %d, а хотели %d", len(ss), tc.length)
		}

	}
	// для наглядности выводим результат
	t.Logf("\n%+v\n\n", ss)

}

// проверяем инициализацию данных
func TestInitData(t *testing.T) {

	// экземпляр хранилища в памяти
	var s *Storage

	// вызываем тестируемый метод
	err := s.InitData()
	if err != nil {
		t.Error(err)
	}

	ss, err := s.Securities()
	if err != nil {
		t.Error(err)
	}

	if len(ss) < 1 {
		t.Errorf("невереный результат: получили %d, а хотели '> 1'", len(ss))
	}

	// для наглядности выводим результат
	t.Logf("\n%+v\n\n", ss[0])

}
