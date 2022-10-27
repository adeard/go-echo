package controller

import (
	"encoding/base64"
	"fmt"
	"go-echo/config"
	"go-echo/model"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func UserCreate() echo.HandlerFunc {
	return func(c echo.Context) error {
		user := new(model.Users)
		response := new(Response)

		c.Bind(user)

		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 8)
		user.Password = string(hashedPassword)

		if err := c.Validate(user); err != nil {

			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		test, isNull := model.GetOneByEmail(user.Email)
		if isNull == nil {

			config.MakeLogEntry(nil).Info("User Exist : " + test.Nama)

			return echo.NewHTTPError(http.StatusBadRequest, "User Already Exist")
		}

		contentType := c.Request().Header.Get("Content-type")
		if contentType == "application/json" {
			fmt.Println("Request dari json")
		} else if strings.Contains(contentType, "multipart/form-data") || contentType == "application/x-www-form-urlencoded" {
			file, err := c.FormFile("ktp")
			if err != nil {
				fmt.Println("Ktp kosong")
			} else {
				src, err := file.Open()
				if err != nil {
					return err
				}
				defer src.Close()
				dst, err := os.Create(file.Filename)
				if err != nil {
					return err
				}
				defer dst.Close()
				if _, err = io.Copy(dst, src); err != nil {
					return err
				}

				user.Ktp = file.Filename
				fmt.Println("Ada file, akan disimpan")
			}
		}

		if user.CreateUser() != nil { // method create user
			return echo.NewHTTPError(http.StatusBadRequest, "Failed create data user")
		} else {
			response.ErrorCode = 0
			response.Message = "Success create data user"
			response.Data = *user
		}
		return c.JSON(http.StatusOK, response)
	}
}

func UserUpdate() echo.HandlerFunc {
	return func(c echo.Context) error {
		user := new(model.Users)
		c.Bind(user)
		response := new(Response)
		if user.UpdateUser(c.Param("email")) != nil { // method update user
			response.ErrorCode = 10
			response.Message = "Gagal update data user"
		} else {
			response.ErrorCode = 0
			response.Message = "Sukses update data user"
			response.Data = *user
		}
		return c.JSON(http.StatusOK, response)
	}
}

func UserDelete() echo.HandlerFunc {
	return func(c echo.Context) error {
		user, _ := model.GetOneByEmail(c.Param("email")) // method get by email
		response := new(Response)

		if user.DeleteUser() != nil { // method update user
			response.ErrorCode = 10
			response.Message = "Gagal menghapus data user"
		} else {
			response.ErrorCode = 0
			response.Message = "Sukses menghapus data user"
		}
		return c.JSON(http.StatusOK, response)
	}
}

func UserGet() echo.HandlerFunc {
	return func(c echo.Context) error {
		response := new(Response)
		users, err := model.GetAll(c.QueryParam("keywords")) // method get all
		if err != nil {
			response.ErrorCode = 10
			response.Message = "Gagal melihat data user"
		} else {
			response.ErrorCode = 0
			response.Message = "Sukses melihat data user"
			response.Data = users
		}

		return c.JSON(http.StatusOK, response)
	}
}

func ProtectedVerify() echo.HandlerFunc {
	return func(c echo.Context) error {
		response := new(Response)
		data, err := base64.StdEncoding.DecodeString(c.FormValue("u"))
		if err != nil {
			return err
		}

		email := string(data)

		user, err := model.GetOneByEmail(email)
		if err != nil {
			response.ErrorCode = http.StatusNotFound
			response.Message = "User Not Found"
		} else {

			result := map[string]int{"login_user_id": user.Id}

			response.ErrorCode = 0
			response.Message = "User Detected"
			response.Data = result
		}

		return c.JSON(http.StatusOK, response)
	}
}
