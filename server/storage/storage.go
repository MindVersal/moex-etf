// Package storage описывает общие требования к поставщику хранилища и используемые типы данных
package storage

// Security - ценная бумага
type Security struct {
	ID        string  // ticker
	Name      string  // полное имя бумаги
	IssueDate int64   // дата выпуска в обращение
	Quotes    []Quote // котировки
}

// Quote - котировка ценной бумаги (цена 'close')
type Quote struct {
	SecurityID string  // ticker
	Num        int     // номер измерения (номер месяца)
	TimeStamp  int64   // отметка времени в формате Unix Time
	Price      float64 // цена закрытия
}

// Interface - контракт для драйвера хранилища котировок
type Interface interface {
	InitData() error                 // инициализирует хранилище данными с сервера Мосбиржи
	Securities() ([]Security, error) // получить список бумаг с котировками
}
