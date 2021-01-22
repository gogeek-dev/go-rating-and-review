package controller

import (
	"fmt"
	mysqldb "gogeek/connection"
	"gogeek/models"
	"log"

	"github.com/gin-gonic/gin"
)

func Index(c *gin.Context) {
	user := models.Product{}
	res := []models.Product{}
	session, _ := store.Get(c.Request, "mysession")
	// session.Options = &sessions.Options{
	// 	Path:     "/",
	// 	MaxAge:   86400 * 7,
	// 	HttpOnly: true,
	// }
	userid := session.Values["userid"]

	log.Println("session userid", userid)
	if nil != userid {

		db := mysqldb.SetupDB()
		// sle, err := db.Query("select id,title,price,image_path from tbl_products")
		// if err != nil {
		// 	panic(err.Error())
		// }

		// for sle.Next() {
		// 	var id int
		// 	var title, imagepath string
		// 	var price float32
		// 	err = sle.Scan(&id, &title, &price, &imagepath)
		// 	if err != nil {
		// 		panic(err.Error())
		// 	}
		// 	user.ID = id
		// 	user.Title = title
		// 	user.Price = price
		// 	user.Imagepath = imagepath
		// 	res1 = append(res1, user)
		// }

		selDB, err := db.Query("SELECT tbl_products.id,tbl_products.title,tbl_products.price,tbl_products.image_path,tbl_products.random_id, sum(tbl_rating_reviews.rating),count(tbl_rating_reviews.user_id) FROM `tbl_products` INNER JOIN tbl_rating_reviews ON tbl_rating_reviews.product_id= tbl_products.id group by(tbl_products.id) ")
		if err != nil {
			panic(err.Error())
		}

		for selDB.Next() {
			var id, overallrating, rcount, randomid int
			var title, imagepath string
			var price float32
			err = selDB.Scan(&id, &title, &price, &imagepath, &randomid, &overallrating, &rcount)
			if err != nil {
				panic(err.Error())
			}
			user.ID = id
			user.Title = title
			user.Price = price
			user.Imagepath = imagepath
			user.Randomid = randomid
			user.Reviewcount = rcount
			allrating := overallrating / rcount
			log.Println("rat value is :", allrating)
			remaining := 5 - allrating
			user.Overallrating = make([]int, allrating)
			user.Remain = make([]int, remaining)
			// user.Overallrating = overallrating / rcount
			for i := 0; i < allrating; i++ {

				user.Overallrating[i] = i + 1
				log.Println("rat is :", user.Overallrating[i])
				// res1 = append(res1, rate)

			}
			log.Println("rem is :", remaining)
			for j := 0; j < remaining; j++ {
				user.Remain[j] = j + 1
				log.Println("rem is :", user.Remain[j])

			}
			updt, err := db.Prepare("UPDATE  tbl_products SET overall_rating=? WHERE id=?")
			if err != nil {
				panic(err.Error())
			}
			updt.Exec(allrating, id)
			res = append(res, user)
		}
	} else {
		c.Redirect(301, "/")
		// c.HTML(200, "login.html", nil)
	}
	// session, _ := store.Get(c.Request, "mysession")
	username := session.Values["name"]

	c.HTML(200, "index.html", gin.H{"index": res, "name": username})

}

func Productdetail(c *gin.Context) {
	db := mysqldb.SetupDB()
	session, _ := store.Get(c.Request, "mysession")
	username := session.Values["name"]

	userid := session.Values["userid"]
	log.Println("session username", username)
	log.Println("session userid", userid)
	if nil == userid {
		c.Redirect(301, "/")
		// c.HTML(200, "login.html", nil)
	}
	id := c.Request.URL.Query().Get("id")
	var pid int
	_ = db.QueryRow("SELECT id FROM tbl_products WHERE random_id = ? ", id).Scan(&pid)
	selDB, err := db.Query("SELECT tbl_products.id,tbl_products.title,tbl_products.price,tbl_products.image_path, tbl_products.description,tbl_products.random_id, sum(tbl_rating_reviews.rating),count(tbl_rating_reviews.user_id) FROM `tbl_products` INNER JOIN tbl_rating_reviews ON tbl_rating_reviews.product_id= tbl_products.id where tbl_products.id =? group by(tbl_products.id)  ", pid)
	if err != nil {
		panic(err.Error())
	}
	user := models.Product{}
	res := []models.Product{}
	for selDB.Next() {
		var pid, overallrating, rcount, randomid int
		var title, imagepath, description string
		var price float32
		err = selDB.Scan(&pid, &title, &price, &imagepath, &description, &randomid, &overallrating, &rcount)
		if err != nil {
			panic(err.Error())
		}
		user.ID = pid
		user.Title = title
		user.Price = price
		user.Imagepath = imagepath
		user.Description = description
		user.Randomid = randomid
		user.Reviewcount = rcount
		allrating := overallrating / rcount
		log.Println("rat value is :", allrating)
		remaining := 5 - allrating
		user.Overallrating = make([]int, allrating)
		user.Remain = make([]int, remaining)
		// user.Overallrating = overallrating / rcount
		for i := 0; i < allrating; i++ {

			user.Overallrating[i] = i + 1
			log.Println("rat is :", user.Overallrating[i])
			// res1 = append(res1, rate)

		}
		log.Println("rem is :", remaining)
		for j := 0; j < remaining; j++ {
			user.Remain[j] = j + 1
			log.Println("rem is :", user.Remain[j])

		}
		userid := session.Values["userid"]

		log.Println("productid is :", pid)
		log.Println("userid is :", userid)
		var count int
		_ = db.QueryRow("SELECT count(product_id) cnt FROM tbl_orders WHERE user_id = ? AND product_id = ?", userid, pid).Scan(&count)
		fmt.Println("count", count)

		user.Value = count

		res = append(res, user)
	}
	log.Println("res is : ", res)

	id1 := c.Request.URL.Query().Get("id")
	var pid2 int
	_ = db.QueryRow("SELECT id FROM tbl_products WHERE random_id = ? ", id1).Scan(&pid2)
	selDB1, err := db.Query("SELECT tbl_user.first_name,tbl_user.location, tbl_rating_reviews.rating,tbl_rating_reviews.review,tbl_rating_reviews.product_id,tbl_rating_reviews.user_id FROM `tbl_user` INNER JOIN tbl_rating_reviews ON tbl_user.id = tbl_rating_reviews.user_id WHERE product_id=?", pid2)

	// id1 := c.Request.URL.Query().Get("id")
	// selDB1, err := db.Query("SELECT rating,review FROM tbl_rating_reviews WHERE product_id=?", id1)
	// if err != nil {
	// 	panic(err.Error())
	// }
	rate := models.Review{}
	res1 := []models.Review{}
	// rat := models.Reviewrating{}
	// rat1 := []models.Reviewrating{}
	for selDB1.Next() {
		var rating, productid, userid int
		var firstname, location, review string

		// var rem []int
		err = selDB1.Scan(&firstname, &location, &rating, &review, &productid, &userid)
		if err != nil {
			panic(err.Error())
		}
		rate.Name = firstname
		rate.Location = location
		// rate.Rating = rating
		rate.Review = review
		rate.Userid = userid
		remaining := 5 - rating
		rate.Rat = make([]int, rating)
		rate.Rem = make([]int, remaining)

		for i := 0; i < rating; i++ {

			rate.Rat[i] = i + 1
			log.Println("rat is :", rate.Rat[i])
			// res1 = append(res1, rate)

		}

		log.Println("rem is :", remaining)
		for j := 0; j < remaining; j++ {
			rate.Rem[j] = j + 1
			log.Println("rem is :", rate.Rem[j])

		}
		res1 = append(res1, rate)
	}

	log.Println("res1 is : ", res1)

	m := map[string]interface{}{
		"product": res,
		"review":  res1,
		"name":    username,
	}

	defer db.Close()
	c.HTML(200, "product-detail.html", m)
}

// SELECT tbl_products.id,tbl_products.title,tbl_products.price,tbl_products.image_path, tbl_rating_reviews.rating,tbl_rating_reviews.review,tbl_rating_reviews.product_id,tbl_rating_reviews.user_id FROM `tbl_products` INNER JOIN tbl_rating_reviews ON tbl_rating_reviews.product_id= tbl_products.id WHERE product_id =3

// SELECT tbl_products.id,tbl_products.title,tbl_products.price,tbl_products.image_path, sum(tbl_rating_reviews.rating),count(tbl_rating_reviews.user_id) FROM `tbl_products` INNER JOIN tbl_rating_reviews ON tbl_rating_reviews.product_id= tbl_products.id group by(tbl_products.id)
