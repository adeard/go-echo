package controller

import (
	"go-echo/model"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func InventoryGet() echo.HandlerFunc {
	return func(c echo.Context) error {
		response := new(Response)
		inventory, err := model.InventoryGetAll(c.QueryString()) // method get all
		if err != nil {
			response.ErrorCode = 10
			response.Message = "Gagal melihat data inventory"
		} else {
			response.ErrorCode = 0
			response.Message = "Sukses melihat data inventory"
			response.Data = inventory
		}

		return c.JSON(http.StatusOK, response)
	}
}

func InventoryCreate() echo.HandlerFunc {
	return func(c echo.Context) error {
		inventory := new(model.Inventory)
		response := new(Response)

		c.Bind(inventory)

		if inventory.CreateInventory() != nil { // method create user
			return echo.NewHTTPError(http.StatusBadRequest, "Failed create data inventory")
		} else {
			response.ErrorCode = 0
			response.Message = "Success create data inventory"
			response.Data = *inventory
		}
		return c.JSON(http.StatusOK, response)
	}
}

func InventoryUpdate() echo.HandlerFunc {
	return func(c echo.Context) error {
		inventory := new(model.Inventory)

		id, _ := strconv.Atoi(c.Param("id"))

		c.Bind(inventory)
		response := new(Response)
		if inventory.UpdateInventory(id) != nil { // method update inventory
			response.ErrorCode = 10
			response.Message = "Gagal update data inventory"
		} else {
			response.ErrorCode = 0
			response.Message = "Sukses update data inventory"
			response.Data = *inventory
		}
		return c.JSON(http.StatusOK, response)
	}
}

func InventoryDelete() echo.HandlerFunc {
	return func(c echo.Context) error {
		inventory, _ := model.InventoryGetOneById(c.Param("id")) // method get by email
		response := new(Response)

		if inventory.DeleteInventory() != nil { // method update user
			response.ErrorCode = 10
			response.Message = "Gagal menghapus data user"
		} else {
			response.ErrorCode = 0
			response.Message = "Sukses menghapus data user"
		}
		return c.JSON(http.StatusOK, response)
	}
}
