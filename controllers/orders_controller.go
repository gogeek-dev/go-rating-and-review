package controller

import (
	"fmt"
	mysqldb "gogeek/connection"
	"gogeek/models"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func Orders(c *gin.Context) {
	db := mysqldb.SetupDB()
	session, _ := store.Get(c.Request, "mysession")
	id := c.Query("id")
	var productid int
	_ = db.QueryRow("SELECT id FROM tbl_products WHERE random_id = ? ", id).Scan(&productid)
	qty := c.Request.PostFormValue("qty")
	username := session.Values["name"]
	userid := session.Values["userid"]
	log.Println("session username", username)
	log.Println("productid userid", productid)
	log.Println("session userid", userid)
	if nil == userid {
		c.Redirect(301, "/")
	}
	createdate := time.Now().Format("Jan 2,2006 3:4:5 PM")
	var count int
	_ = db.QueryRow("SELECT count(product_id) cnt FROM tbl_orders WHERE user_id = ? AND product_id = ?", userid, productid).Scan(&count)
	fmt.Println("count", count)
	if count != 1 {

		insForm, err := db.Prepare("INSERT INTO tbl_orders(product_id, user_id,quantity,created_date) VALUES(?,?,?,?)")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(productid, userid, qty, createdate)
	}

	log.Println("productid is :", productid)
	log.Println("Quantity is :", qty)
	log.Println("userid is :", userid)
	log.Println("cdate is :", createdate)
	selDB, err := db.Query("SELECT product_id,created_date FROM  tbl_orders where user_id=?", userid)
	if err != nil {
		panic(err.Error())
	}
	// product := models.Order{}
	// res := []models.Order{}
	product1 := models.Product{}
	res := []models.Product{}
	for selDB.Next() {
		var productids int
		var orderdate string

		err = selDB.Scan(&productids, &orderdate)
		if err != nil {
			panic(err.Error())
		}
		log.Println("productid is :", productids)
		// product.ID = id
		// product.Productid = productid

		selDB1, err := db.Query("SELECT id,title,price,image_path,random_id FROM  tbl_products where id=?", productids)
		if err != nil {
			panic(err.Error())
		}

		for selDB1.Next() {
			var id, randomid int
			var price float32
			var title, imagepath string
			err = selDB1.Scan(&id, &title, &price, &imagepath, &randomid)
			if err != nil {
				panic(err.Error())
			}
			product1.ID = id
			product1.Title = title
			product1.Price = price
			product1.Imagepath = imagepath
			product1.Randomid = randomid
			product1.Created = orderdate
			// product1.Value = count
			var count int
			_ = db.QueryRow("SELECT count(product_id) cnt FROM tbl_rating_reviews WHERE user_id = ? AND product_id = ?", userid, id).Scan(&count)
			fmt.Println("count", count)

			product1.Value = count

			res = append(res, product1)
			log.Println("productid is :", res)
			log.Println("productid is :", product1.Value)
		}
	}
	c.HTML(200, "my-orders.html", gin.H{"orders": res, "name": username})
}
func Ordersview(c *gin.Context) {
	db := mysqldb.SetupDB()
	session, _ := store.Get(c.Request, "mysession")
	username := session.Values["name"]
	userid := session.Values["userid"]
	log.Println("session username", username)
	log.Println("session userid", userid)
	if nil == userid {
		c.Redirect(301, "/")
	}
	selDB, err := db.Query("SELECT product_id,created_date FROM  tbl_orders where user_id=?", userid)
	if err != nil {
		panic(err.Error())
	}
	// product := models.Order{}
	// res := []models.Order{}
	product1 := models.Product{}
	res := []models.Product{}
	for selDB.Next() {
		var productids int
		var orderdate string

		err = selDB.Scan(&productids, &orderdate)
		if err != nil {
			panic(err.Error())
		}
		log.Println("productid is :", productids)
		log.Println("order is :", orderdate)
		// product.ID = id
		// product.Productid = productid

		selDB1, err := db.Query("SELECT id,title,price,image_path,random_id FROM  tbl_products where id=?", productids)
		if err != nil {
			panic(err.Error())
		}

		for selDB1.Next() {
			var id, randomid int
			var price float32
			var title, imagepath string
			err = selDB1.Scan(&id, &title, &price, &imagepath, &randomid)
			if err != nil {
				panic(err.Error())
			}
			product1.ID = id
			product1.Title = title
			product1.Price = price
			product1.Imagepath = imagepath
			product1.Randomid = randomid
			product1.Created = orderdate
			// product1.Value = count
			var count int
			_ = db.QueryRow("SELECT count(product_id) cnt FROM tbl_rating_reviews WHERE user_id = ? AND product_id = ?", userid, id).Scan(&count)
			fmt.Println("count", count)

			product1.Value = count

			res = append(res, product1)
			log.Println("productid is :", res)
			log.Println("created date is :", product1.Created)
		}
	}
	c.HTML(200, "my-orders.html", gin.H{"orders": res, "name": username})
}
