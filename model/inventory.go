package model

import (
	"go-echo/config"
	"net/url"
)

type Inventory struct {
	Id           int    `json:"id" form:"id"`
	Name         string `json:"name" form:"name" validate:"required,name"`
	CategoryType string `json:"category_type" form:"category_type" validate:"required"`
	Stock        int    `json:"stock" form:"stock" validate:"required"`
	Price        int    `json:"price" form:"price" validate:"required"`
}

func (inventory *Inventory) CreateInventory() error {
	if err := config.DB.Create(inventory).Error; err != nil {
		return err
	}
	return nil
}

func (inventory *Inventory) UpdateItem(id int) error {
	if err := config.DB.Model(&Inventory{}).Where("id = ?", id).Updates(inventory).Error; err != nil {
		return err
	}
	return nil
}

func (inventory *Inventory) DeleteItem() error {
	if err := config.DB.Delete(inventory).Error; err != nil {
		return err
	}
	return nil
}

func InventoryGetOneById(id string) (Inventory, error) {
	var inventory Inventory
	result := config.DB.Where("id = ?", id).First(&inventory)
	return inventory, result.Error
}

func InventoryGetAll(params string) ([]Inventory, error) {
	var inventory []Inventory

	q, err := url.ParseQuery(params)
	if err != nil {
		panic(err)
	}

	var name = q.Get("name")
	var category_type = q.Get("category_type")
	var stock = q.Get("stock")
	var price = q.Get("price")

	m := make(map[string]interface{})

	if category_type != "" {
		m["category_type"] = q.Get("category_type")
	}

	if stock != "" {
		m["stock"] = q.Get("stock")
	}

	if price != "" {
		m["price"] = q.Get("price")
	}

	if name != "" {
		m["name"] = q.Get("name")
	}

	result := config.DB.Where(m).Find(&inventory)

	return inventory, result.Error
}
