package base

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"mail/internal/config"
	"mail/pkg/logging"
	"os"
	"path"
	"strings"
	"time"
)

type DB struct {
	db     *sql.DB
	logger *logging.Logger
}

type Subscribers struct {
	Id         int
	FirstName  string
	SecondName string
	FatherName string
	Email      string
	BirthDay   time.Time
	Sex        string
}

func NewDataBase(cfg *config.Config, logger *logging.Logger) (DB, error) {

	err := os.MkdirAll(cfg.DataPath, 0777)
	if err != nil {
		logger.Error(err)
	}

	sqlite, err := sql.Open("sqlite3", path.Join(cfg.DataPath, "subscribers.db"))
	if err != nil {
		logger.Error(err)
	}
	_, err = sqlite.Exec(fmt.Sprintf("CREATE TABLE IF NOT EXISTS subscriber (id INTEGER NOT NULL PRIMARY KEY, " +
		"first_name VARCHAR (30) NOT NULL, " +
		"second_name VARCHAR (30) DEFAULT NONE, " +
		"father_name VARCHAR (30) DEFAULT NONE, " +
		"birth_day TIMESTAMP DEFAULT CURRENT_TIMESTAMP, " +
		"email VARCHAR (50) NOT NULL, " +
		"sex VARCHAR (20) DEFAULT NONE)"))
	if err != nil {
		logger.Error(err)
	}

	return DB{
		db:     sqlite,
		logger: logger,
	}, nil
}

func (d *DB) CreateNewSubscriber(u Subscribers) error {

	var user Subscribers

	query, err := d.db.Query(fmt.Sprintf("SELECT email FROM subscriber "))
	if err != nil {
		d.logger.Error(err)
		return err
	}

	for query.Next() {
		query.Scan(&user.Email)
		if user.Email == u.Email {
			//d.logger.Warn("Адрес уже существует")
			return nil
		}
	}

	q := fmt.Sprintf("INSERT INTO 'subscriber' (first_name, second_name, father_name, email, birth_day, sex) "+
		"VALUES ('%s', '%s', '%s', '%s', '%v', '%s')", u.FirstName, u.SecondName, u.FatherName, u.Email, u.BirthDay.UTC().Format(config.StructDateFormat), u.Sex)
	_, err = d.db.Exec(q)
	if err != nil {
		d.logger.Error(err)
		return err
	}

	return nil
}

func (d *DB) GetSubscribers() (subs []Subscribers, err error) {

	var sub Subscribers

	query, err := d.db.Query(fmt.Sprintf("SELECT * FROM subscriber "))
	if err != nil {
		d.logger.Error(err)
		return nil, err
	}

	for query.Next() {
		err = query.Scan(&sub.Id, &sub.FirstName, &sub.SecondName, &sub.FatherName, &sub.BirthDay, &sub.Email, &sub.Sex)
		if err != nil {
			d.logger.Error(err)
		}
		subs = append(subs, sub)
	}

	return subs, nil

}

func (d *DB) AddNewSubscribers(file []byte) (err error) {

	var s Subscribers
	var ss []Subscribers
	us := strings.Split(strings.TrimSpace(string(file)), ";")

	for _, v := range us {
		u := strings.Split(strings.TrimSpace(v), ",")
		if len(u) == 6 {
			s.FirstName = strings.TrimSpace(u[0])
			s.SecondName = strings.TrimSpace(u[1])
			s.FatherName = strings.TrimSpace(u[2])
			s.Email = strings.TrimSpace(u[3])

			s.BirthDay, err = time.Parse(config.StructDateFormat, strings.TrimSpace(u[4]))
			if err != nil {
				d.logger.Error(err)
				return err
			}

			s.Sex = strings.TrimSpace(u[5])
			ss = append(ss, s)
			if err = d.CreateNewSubscriber(s); err != nil {
				d.logger.Error(err)
				return err
			}
		}
	}
	return nil
}
