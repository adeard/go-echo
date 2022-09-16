package model

import (
	"go-echo/config"
)

type Users struct {
	Email       string `json:"email" form:"email" validate:"required,email"`
	Nama        string `json:"nama" form:"nama" validate:"required"`
	NoHandphone string `json:"no_handphone" form:"no_handphone" validate:"required"`
	Alamat      string `json:"alamat" form:"alamat" validate:"required"`
	Ktp         string `json:"ktp" form:"ktp" validate:"required"`
	Password    string `json:"password" form:"password" validate:"required"`
}

func (user *Users) CreateUser() error {
	if err := config.DB.Create(user).Error; err != nil {
		return err
	}
	return nil
}

func (user *Users) UpdateUser(email string) error {
	if err := config.DB.Model(&Users{}).Where("email = ?", email).Updates(user).Error; err != nil {
		return err
	}
	return nil
}

func (user *Users) DeleteUser() error {
	if err := config.DB.Delete(user).Error; err != nil {
		return err
	}
	return nil
}

func GetOneByEmail(email string) (Users, error) {
	var user Users
	result := config.DB.Where("email = ?", email).First(&user)
	return user, result.Error
}

func GetAll(keywords string) ([]Users, error) {
	var users []Users
	result := config.DB.Where("email LIKE ? OR nama LIKE ?", "%"+keywords+"%", "%"+keywords+"%").Find(&users)

	return users, result.Error
}
