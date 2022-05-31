package service

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"market-bot/internal/dao"
	"market-bot/internal/service/model"
	"market-bot/sdk/tgbot"
)

type EmployeeService struct {
	employeeRepository dao.EmployeeRepository
	userRepository     dao.UserRepository
}

func NewEmployeeService(repository *dao.Repository) *EmployeeService {
	return &EmployeeService{repository, repository}
}

func (us EmployeeService) SaveEmployee(employee model.Employee) error {
	err := validation.ValidateStruct(&employee,
		validation.Field(&employee.Status, validation.Required),
		validation.Field(&employee.TelegramId, validation.Required),
	)
	if err != nil {
		return fmt.Errorf("validate employee: %w", err)
	}
	user := tgbot.User{UserId: &employee.TelegramId, Status: tgbot.Active}
	if employee.Status == model.NewEmployee {
		if err = us.userRepository.SaveUser(user); err != nil {
			return err
		}
	} else if employee.Status == model.ActiveEmployee {
		if err = us.userRepository.UpdateUser(user); err != nil {
			return err
		}
	}

	if err = us.employeeRepository.SaveEmployee(employee); err != nil {
		return err
	}
	return nil
}

func (us EmployeeService) UpdateEmployee(employee model.Employee, employeeId string) error {
	errors := validation.Errors{"employee": validation.ValidateStruct(&employee,
		validation.Field(&employee.Status, validation.Required),
		validation.Field(&employee.TelegramId, validation.Required),
	), "employeeId": validation.Validate(employeeId, validation.Required),
	}.Filter()
	if errors != nil {
		return fmt.Errorf("validate employee: %v", errors)
	}
	if err := us.employeeRepository.UpdateEmployee(employee, employeeId); err != nil {
		return err
	}
	return nil
}

func (us EmployeeService) GetEmployeeById(employeeId string) (model.Employee, error) {
	employee, err := us.employeeRepository.GetEmployeeById(employeeId)
	if err != nil {
		return employee, err
	}
	return employee, nil
}

func (us EmployeeService) GetEmployeeByTelegramId(telegramId int64) (model.Employee, error) {
	employee, err := us.employeeRepository.GetEmployeeByTelegramId(telegramId)
	if err != nil {
		return employee, err
	}
	return employee, nil
}

func (us EmployeeService) GetEmployeeByTelegramUserName(userName string) (model.Employee, error) {
	employee, err := us.employeeRepository.GetEmployeeByTelegramUserName(userName)
	if err != nil {
		return employee, err
	}
	return employee, nil
}

func (us EmployeeService) GetEmployees(search string) ([]model.Employee, error) {
	employees, err := us.employeeRepository.GetEmployees(search)
	if err != nil {
		return nil, err
	}
	return employees, nil
}

func (s EmployeeService) GetRandomUser() (tgbot.User, error) {
	return s.userRepository.GetRandomUser()
}
