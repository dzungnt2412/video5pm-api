package services

import (
	"video5pm-api/models/entity"

	"github.com/jinzhu/gorm"
)

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

type UserService struct {
	db *gorm.DB
}

type ShopPoint struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Orders string `json:"orders"`
	Point  string `json:"point"`
}

type UserInfo struct {
	UserName    string `json:"username" gorm:"column:username"`
	ID          string `json:"id"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	FullName    string `json:"full_name"`
	PackageQt   string `json:"package_qt"`
	PackageVn   string `json:"package_vn"`
}

//FindUser - service query find user by id
func (c *UserService) FindUser(uid int64) (entity.User, error) {
	var user entity.User
	err := c.db.First(&user, uid).Error
	if err != nil {
		return user, err
	}
	return user, nil
}

//Find shop by user_id - service query find user by id
func (c *UserService) FindUserByUsername(username string) (UserInfo, error) {

	var userInfo UserInfo

	err := c.db.Table("users").Select("users.id, users.username, users.full_name, users.phone_number, users.email, vn.name as package_vn, qt.name as package_qt").Joins("LEFT JOIN packages AS vn ON vn.id = users.package_vn_id").Joins("LEFT JOIN packages AS qt ON qt.id = users.package_id").Where("users.username = ?", username).Or("users.email = ?", username).Find(&userInfo).Error
	if err != nil {
		return userInfo, err
	}

	return userInfo, nil
}

func (c *UserService) GetShopPointVN(username string) ([]ShopPoint, error) {

	var shopPointVN []ShopPoint

	err := c.db.Table("shops").Select("shops.id, shops.name, count(orders.id) AS orders, sum(order_items.quantity*product_variants.point) AS point").Joins("LEFT JOIN users ON users.id=shops.user_id").Joins("LEFT JOIN orders ON orders.shop_id=shops.id").Joins("LEFT JOIN order_items ON order_items.order_id=orders.id").Joins("LEFT JOIN product_variants ON product_variants.id=order_items.product_variant_id").Joins("LEFT JOIN consignments ON order_items.consignment_id=consignments.id").Where("consignments.supplier_id = ? AND (users.username = ? OR users.email = ?) AND orders.status IN (?)", 730, username, username, []string{"process", "paid", "fulfilled"}).Group("shops.name, shops.id").Find(&shopPointVN).Error
	// err := database.DB.Table("shops").Select("shops.id, shops.name, count(orders.id) AS orders").Joins("LEFT JOIN users ON users.id=shops.user_id").Joins("LEFT JOIN orders ON orders.shop_id=shops.id").Joins("LEFT JOIN order_items ON order_items.order_id=orders.id").Joins("LEFT JOIN consignments ON order_items.consignment_id=consignments.id").Where("consignments.supplier_id = ? AND users.username = ?", 730, username).Group("shops.name, shops.id").Find(&shopPointVN).Error
	if err != nil {
		return shopPointVN, err
	}

	return shopPointVN, nil
}

func (c *UserService) GetShopPointGlobal(username string) ([]ShopPoint, error) {

	var shopPointGlobal []ShopPoint

	err := c.db.Table("shops").Select("shops.id, shops.name, count(orders.id) AS orders, sum(order_items.quantity*product_variants.point) AS point").Joins("LEFT JOIN users ON users.id=shops.user_id").Joins("LEFT JOIN orders ON orders.shop_id=shops.id").Joins("LEFT JOIN order_items ON order_items.order_id=orders.id").Joins("LEFT JOIN product_variants ON product_variants.id=order_items.product_variant_id").Joins("LEFT JOIN consignments ON order_items.consignment_id=consignments.id").Where("consignments.supplier_id != ? AND (users.username = ? OR users.email = ?) AND orders.status IN (?)", 730, username, username, []string{"process", "paid", "fulfilled"}).Group("shops.name, shops.id").Find(&shopPointGlobal).Error
	// err := database.DB.Table("shops").Select("shops.id, shops.name, count(orders.id) AS orders").Joins("LEFT JOIN users ON users.id=shops.user_id").Joins("LEFT JOIN orders ON orders.shop_id=shops.id").Joins("LEFT JOIN order_items ON order_items.order_id=orders.id").Joins("LEFT JOIN consignments ON order_items.consignment_id=consignments.id").Where("consignments.supplier_id != ? AND users.username = ?", 730, username).Group("shops.name, shops.id").Find(&shopPointGlobal).Error
	if err != nil {
		return shopPointGlobal, err
	}

	return shopPointGlobal, nil
}
