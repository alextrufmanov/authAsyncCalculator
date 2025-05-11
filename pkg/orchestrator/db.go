package orchestrator

import (
	"context"
	"database/sql"
	"log"

	"github.com/alextrufmanov/asyncCalculator/pkg/models"
	_ "github.com/mattn/go-sqlite3"
)

type DBStorage struct {
	db  *sql.DB
	ctx context.Context
}

// Функция создает новое БД хранилище
func NewDBStorage() (*DBStorage, error) {

	ctx := context.TODO()

	db, err := sql.Open("sqlite3", "asyncCalculator.db")
	if err != nil {
		return nil, err
	}

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	expressionsSQL := `
 	CREATE TABLE IF NOT EXISTS expressions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER,
		expression TEXT NOT NULL,
		result FLOAT64 DEFAULT 0,
		Status TEXT NOT NULL DEFAULT "ready"
 	);`
	_, err = db.Exec(expressionsSQL)
	if err != nil {
		return nil, err
	}

	usersSQL := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		login TEXT NOT NULL UNIQUE,
		psw_hash TEXT NOT NULL DEFAULT ""
	);`
	_, err = db.Exec(usersSQL)
	if err != nil {
		return nil, err
	}

	return &DBStorage{
		db:  db,
		ctx: ctx,
	}, nil
}

// Функция закрывает БД
func (dbs *DBStorage) Close() {
	dbs.db.Close()
}

// Функция добавляет в БД нового пользователя
func (dbs *DBStorage) InsertUser(login string, password string) (int, bool) {
	var sql = `
	INSERT INTO users (login, psw_hash) values ($1, $2)
	`
	pswHash, res := Hash(login, password)
	if !res {
		return 0, false
	}
	result, err := dbs.db.ExecContext(dbs.ctx, sql, login, pswHash)
	if err != nil {
		return 0, false
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, false
	}
	return int(id), true
}

// Функция проверяет в БД пользователя, возвращает JWT
func (dbs *DBStorage) GetToken(login string, password string) (string, bool) {
	var id int
	var pswHashBd string
	var sql = `
	SELECT id, psw_hash FROM users WHERE login = $1
	`
	err := dbs.db.QueryRowContext(dbs.ctx, sql, login).Scan(&id, &pswHashBd)
	if err != nil {
		log.Printf("  GetToken 002 -> %v", err)
		return "", false
	}
	if !TestHash(login, password, pswHashBd) {
		return "", false
	}
	jwt, res := GetJWT(login, id)
	if res {
		return jwt, true
	}
	return "", false
}

// Функция добавляет в БД новое арифметическое выражение
func (dbs *DBStorage) AppendExpression(userId int32, expression string) (int32, error) {
	var sql = `
	INSERT INTO expressions (user_id, expression) values ($1, $2)
	`
	result, err := dbs.db.ExecContext(dbs.ctx, sql, userId, expression)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int32(id), nil
}

// Функция возвращает из хранилища все арифметические выражения
func (dbs *DBStorage) GetAllExpressions(userId int32) []models.Expression {
	var id int32
	var user_id int32
	var expression string
	var status string
	var result float64
	r := make([]models.Expression, 0)
	var sql = `
	SELECT id, user_id, expression, status, result FROM expressions WHERE user_id = $1
	`
	rows, err := dbs.db.QueryContext(dbs.ctx, sql, userId)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			err := rows.Scan(&id, &user_id, &expression, &status, &result)
			if err == nil {
				r = append(r, models.Expression{
					Id:         id,
					UserId:     user_id,
					Expression: expression,
					Status:     status,
					Result:     result,
				})
			}
		}
	}
	return r
}

// Функция возвращает из хранилища арифметическое выражение с указанным Id
func (dbs *DBStorage) GetExpressionByID(id int32, userId int32) (models.Expression, bool) {
	var expression_id int32
	var user_id int32
	var expression string
	var status string
	var result float64
	var sql = `
	SELECT id, user_id, expression, status, result FROM expressions WHERE id = $1 and user_id = $2
	`
	err := dbs.db.QueryRowContext(dbs.ctx, sql, id, userId).Scan(&expression_id, &user_id, &expression, &status, &result)
	if err != nil {
		log.Printf("  GetToken 002 -> %v", err)
		return models.Expression{}, false
	}
	return models.Expression{
		Id:         expression_id,
		UserId:     user_id,
		Expression: expression,
		Status:     status,
		Result:     result,
	}, true
}

// Функция возвращает из хранилища арифметические выражения которые надо посчитать
func (dbs *DBStorage) GetActiveExpressions() []models.Expression {
	var id int32
	var user_id int32
	var expression string
	var status string
	var result float64
	r := make([]models.Expression, 0)
	var sql = `
	SELECT id, user_id, expression, status, result FROM expressions WHERE status IN ("ready", "calculate")
	`
	rows, err := dbs.db.QueryContext(dbs.ctx, sql)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			err := rows.Scan(&id, &user_id, &expression, &status, &result)
			if err == nil {
				r = append(r, models.Expression{
					Id:         id,
					UserId:     user_id,
					Expression: expression,
					Status:     status,
					Result:     result,
				})
			}
		}
	}
	return r
}

// Функция устанавливает результат и статус вычисления указанного арифметического выражения
func (dbs *DBStorage) UpdateExpressionStatus(id int32, state string, result float64) error {
	var sql = `
	UPDATE expressions SET Status = $1, result = $2 WHERE id = $3
	`
	_, err := dbs.db.ExecContext(dbs.ctx, sql, state, result, id)
	return err
}
