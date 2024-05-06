package main

import (
	"database/sql"
	"errors"
	"fmt"
)

type ParcelStore struct {
	db *sql.DB
}

// конструктрор ParselStore. на вход нужно подать готовую db,
// которая к моменту вызова должна быть открыта sql.Open("sqlite", "tracker.db")
// вся работа с базой через ParselStore.db
func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

// Метод Add структуры ParselStore.
// возвращает номер добавленной записи и ошибку
func (s ParcelStore) Add(p Parcel) (int, error) {
	// реализуйте добавление строки в таблицу parcel, используйте данные из переменной p
	res, err := s.db.Exec("INSERT INTO parcel (client, status, address, created_at) VALUES (:client, :status, :address, :created_at)",
		sql.Named("client", p.Client),
		sql.Named("status", p.Status),
		sql.Named("address", p.Address),
		sql.Named("created_at", p.CreatedAt))
	// верните идентификатор последней добавленной записи
	// сохраняем номер последней добавленной записи в id
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	// при возврате приводим к int т к тип id int64
	return int(id), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	// реализуйте чтение строки по заданному number
	// здесь из таблицы должна вернуться только одна строка
	// запрашиваем запись с number
	row := s.db.QueryRow("SELECT number, client, status, address, created_at FROM parcel WHERE number = :number",
		sql.Named("number", number))
	// заполните объект Parcel данными из таблицы
	p := Parcel{} // это прекод
	err := row.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
	if err != nil {
		return p, err
		// вариант return Parsel{}, err // а как правильно??
		// если можно то прокомментируйте как правильно и почему
		// желательно подробно)))
	}

	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	// реализуйте чтение строк из таблицы parcel по заданному client
	// здесь из таблицы может вернуться несколько строк
	rows, err := s.db.Query("SELECT number, client, status, address, created_at FROM parcel WHERE client = :client", sql.Named("client", client))
	if err != nil {
		return nil, err
	}

	// заполните срез Parcel данными из таблицы
	var res []Parcel
	for rows.Next() {
		var p Parcel
		err := rows.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
		if err != nil {
			return nil, err
		}
		// добавить в res
		res = append(res, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	// реализуйте обновление статуса в таблице parcel
	_, err := s.db.Exec("UPDATE parcel SET status = :status WHERE number = :number",
		sql.Named("status", status),
		sql.Named("number", number))

	return err
}

func (s ParcelStore) SetAddress(number int, address string) error {
	// реализуйте обновление адреса в таблице parcel

	// запрашиваем с помощье метода Get
	p, err := s.Get(number)
	if err != nil {
		wrapsErr := fmt.Errorf("error reading a record with a number %d, error %w ", number, err)
		return wrapsErr
		// *** было return err
		// *** Ревьюер: Тут нужно обработать ошибку и вернуть ее. Это крит.
		// *** 			У вас в err в этой области видимости по сути новая ошибка, нужно также ее вернуть
		// *** Обёртывание р...ра устроило. В каких случаях прямо нужно обертывать ошибку
		// *** а в каких можно возвращать то, что есть?
	}
	// менять адрес можно только если значение статуса registered
	// тут ошибок нет. Проверяем значение p.Status
	if p.Status == ParcelStatusRegistered {
		// статус соответствует. меняем значение адреса
		_, err = s.db.Exec("UPDATE parcel SET address = :address WHERE number = :number",
			sql.Named("address", address),
			sql.Named("number", number))
		if err != nil {
			wrapsErr := fmt.Errorf("error updating the address of a record with a number %d, error %w ", number, err)
			return wrapsErr
			// *** а можно было return err
			// *** как правильно ???
		}
		return nil
	} else {
		// тут генерируем собственную ощибку, т к статус контролинуем сами
		wrapsErr := fmt.Errorf("error updating the address of a record with a number %d, status not registered", number)
		return wrapsErr
		// *** а можно return errors.New("Status not registered")
		// *** как правильно ?
	}
}

func (s ParcelStore) Delete(number int) error {
	// реализуйте удаление строки из таблицы parcel
	// удалять строку можно только если значение статуса registered
	p, err := s.Get(number)
	if err != nil {
		wrapsErr := fmt.Errorf("Error reading a record with a number %d, error %w ", number, err)
		return wrapsErr
		// *** а можно было return err
		// *** как правильно ???
	}
	if p.Status == ParcelStatusRegistered {
		// удаляем
		_, err := s.db.Exec("DELETE FROM parcel WHERE number = :number", sql.Named("number", number))
		if err != nil {
			// ошибка удаления
			wrapsErr := fmt.Errorf("Error deleting a record with a number %d, error %w ", number, err)
			return wrapsErr
			// *** а можно было return err
			// *** как правильно ???
		}
		// удалили
		return nil
	} else {
		// статус не соответствует
		return errors.New("Nothing has been deleted. Status not registered.")
	}
}
