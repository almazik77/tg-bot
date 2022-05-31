package dao

import (
	"github.com/pkg/errors"
	"market-bot/internal/service/model"
)

type EmployeeRepository interface {
	SaveEmployee(user model.Employee) error
	UpdateEmployee(employee model.Employee, employeeId string) error
	GetEmployeeById(employeeId string) (model.Employee, error)
	GetEmployeeByTelegramId(telegramId int64) (model.Employee, error)
	GetEmployeeByTelegramUserName(userName string) (model.Employee, error)
	GetEmployees(search string) ([]model.Employee, error)
}

func (r *Repository) SaveEmployee(user model.Employee) error {
	insert := `INSERT INTO employee (id, status, first_name, last_name, phone, telegram_id, photo_id, start_message, district, created_date) 
			   VALUES (:id, :status, :first_name, :last_name, :phone, :telegram_id, :photo_id, :start_message, :district, :created_date)`

	if _, err := r.db.NamedExec(insert, user); err != nil {
		return err
	}
	return nil
}

func (r *Repository) UpdateEmployee(employee model.Employee, employeeId string) error {
	query := `UPDATE employee SET status = $1, telegram_id = $2 WHERE id = $3`

	if _, err := r.db.Query(query, employee.Status, employee.TelegramId, employeeId); err != nil {
		return errors.Wrapf(err, "unable update employee, id: %v", employeeId)
	}
	return nil
}

func (r *Repository) GetEmployeeById(employeeId string) (model.Employee, error) {
	row := r.db.QueryRowx("SELECT * FROM employee WHERE id = $1", employeeId)

	employee := model.Employee{}
	if err := row.StructScan(employee); err != nil {
		return employee, errors.Wrapf(err, "unable to get employee, id: %v", employeeId)
	}
	return employee, nil
}

func (r *Repository) GetEmployeeByTelegramId(telegramId int64) (model.Employee, error) {
	row := r.db.QueryRowx("SELECT * FROM employee WHERE telegram_id = $1", telegramId)

	employee := model.Employee{}
	if err := row.StructScan(&employee); err != nil {
		return employee, errors.Wrapf(err, "unable to get employee, telegramId: %v", telegramId)
	}
	return employee, nil
}

func (r *Repository) GetEmployeeByTelegramUserName(userName string) (model.Employee, error) {
	row := r.db.QueryRowx("SELECT * FROM employee WHERE telegram_id = (SELECT user_id FROM profile WHERE user_name = $1)", userName)

	employee := model.Employee{}
	if err := row.StructScan(&employee); err != nil {
		return employee, errors.Wrapf(err, "unable to get employee, userName: %v", userName)
	}
	return employee, nil
}

func (r *Repository) GetEmployees(search string) ([]model.Employee, error) {
	query := `SELECT * FROM employee
				WHERE first_name ilike concat('%', $1::varchar , '%')
				OR last_name ilike concat('%', $1::varchar  , '%')   
				OR phone ilike concat('%', $1::varchar , '%')
				   OR district ilike concat('%', $1::varchar , '%')`
	rows, err := r.db.Queryx(query, search)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to get employees, search: %v", search)
	}
	defer rows.Close()
	res := make([]model.Employee, 0)
	for rows.Next() {
		e := model.Employee{}
		if err = rows.StructScan(&e); err != nil {
			return nil, errors.Wrapf(err, "unable struct scan, rows: %v", rows)
		}
		res = append(res, e)
	}
	return res, nil
}
