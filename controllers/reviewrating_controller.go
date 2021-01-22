package controller

import (
	"fmt"
	mysqldb "gogeek/connection"
	"gogeek/models"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func Reviewrating(c *gin.Context) {
	session, _ := store.Get(c.Request, "mysession")

	userid := session.Values["userid"]

	log.Println("session userid", userid)
	if nil == userid {
		c.Redirect(301, "/")
	}
	db := mysqldb.SetupDB()
	id := c.Request.URL.Query().Get("id")
	var pid int
	_ = db.QueryRow("SELECT id FROM tbl_products WHERE random_id = ? ", id).Scan(&pid)
	log.Println("product id is: ", pid)
	selDB, err := db.Query("SELECT tbl_products.id,tbl_products.title,tbl_products.price,tbl_products.image_path,tbl_products.random_id, sum(tbl_rating_reviews.rating),count(tbl_rating_reviews.user_id) FROM `tbl_products` INNER JOIN tbl_rating_reviews ON tbl_rating_reviews.product_id= tbl_products.id where tbl_products.id=? group by(tbl_products.id) ", pid)
	if err != nil {
		panic(err.Error())
	}
	user := models.Product{}
	res := []models.Product{}
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

		res = append(res, user)
	}
	// session, _ := store.Get(c.Request, "mysession")
	username := session.Values["name"]
	user1 := models.Review{}
	res1 := []models.Review{}
	selDB1, err := db.Query("SELECT rating,review from tbl_rating_reviews where product_id=? and user_id=? ", pid, userid)
	if err != nil {
		panic(err.Error())
	}
	for selDB1.Next() {
		var urating int
		var ureview string
		err = selDB1.Scan(&urating, &ureview)

		user1.Review = ureview
		remaining := 5 - urating

		user1.Updaterate = make([]int, urating)
		user1.Updaterem = make([]int, remaining)

		for k := urating; k > 0; k-- {

			user1.Updaterate[urating-k] = k
			log.Println("rate is :", user1.Updaterate[urating-k])

			// res1 = append(res1, rate)
		}
		b := 1
		for l := remaining; l > 0; l-- {
			log.Println("remde1 is :", remaining)

			a := urating + b
			user1.Updaterem[l-1] = a
			b++
			log.Println("remder is :", user1.Updaterem[l-1])

		}

		res1 = append(res1, user1)

	}
	log.Println("view review is ", user1.Review)

	c.HTML(200, "review-ratings.html", gin.H{"product": res, "PID": user.ID, "name": username, "update": res1})
}

func Reviewratingsave(c *gin.Context) {
	db := mysqldb.SetupDB()
	session, _ := store.Get(c.Request, "mysession")
	id := c.Query("id")
	// var pid int
	// _ = db.QueryRow("SELECT id FROM tbl_products WHERE random_id = ? ", id).Scan(&pid)
	rating := c.Request.PostFormValue("rate")
	review := c.Request.PostFormValue("review")
	userid := session.Values["userid"]
	log.Println("session userid", userid)
	log.Println("product id : ", id)
	if nil == userid {
		c.Redirect(301, "/")
	}
	createdate := time.Now().Format("Jan 2,2006 3:4:5 PM")
	var count int
	_ = db.QueryRow("SELECT count(product_id) cnt FROM tbl_rating_reviews WHERE user_id = ? AND product_id = ?", userid, id).Scan(&count)
	fmt.Println("count", count)
	if count != 1 {
		insForm, err := db.Prepare("INSERT INTO tbl_rating_reviews(product_id, user_id,rating,review,created_date) VALUES(?,?,?,?,?)")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(id, userid, rating, review, createdate)
	} else {
		updorder, err := db.Prepare("UPDATE tbl_rating_reviews SET rating = ?,review = ?,created_date = ? WHERE user_id = ? AND product_id = ?")
		if err != nil {
			panic(err.Error())
		}
		updorder.Exec(rating, review, createdate, userid, id)
	}
	log.Println("id :", id)
	log.Println("userid :", userid)
	log.Println("rating :", rating)
	log.Println("review :", review)
	log.Println("cdate :", createdate)
	// log.Println("INSERT: Userid: " + userid + "|Date: " + createdate)
	c.Redirect(301, "/index")
}
