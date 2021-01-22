package controller

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	mysqldb "gogeek/connection"
	"io"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte("mysession"))

func Loginview(c *gin.Context) {
	c.HTML(200, "login.html", nil)
}

func Login(c *gin.Context) {
	db := mysqldb.SetupDB()
	key := hex.EncodeToString(make([]byte, 32))
	name := c.Request.PostFormValue("name")
	cpwds := c.Request.PostFormValue("password")
	var id int
	var password, firstname string
	err1 := db.QueryRow("select id,first_name,password from tbl_user where email_id=?", name).Scan(&id, &firstname, &password)
	if err1 != nil {
		c.HTML(200, "login.html", gin.H{"error": "**You are not register"})
	}

	decrypted := decrypt(password, key)
	fmt.Printf("decrypted : %s\n", decrypted)

	if cpwds == decrypted {
		session, _ := store.Get(c.Request, "mysession")
		session.Values["name"] = firstname
		session.Values["userid"] = id
		session.Save(c.Request, c.Writer)

		last_login_date := time.Now().Format("Jan 2,2006 3:4:5 PM")
		updt, err := db.Prepare("UPDATE tbl_user SET last_login_date=? WHERE email_id=?")
		if err != nil {
			panic(err.Error())
		}
		updt.Exec(last_login_date, name)
		c.Redirect(301, "/index")
	} else {
		c.HTML(200, "login.html", gin.H{"error": "**Password is incorrect"})
	}

}

func Register(c *gin.Context) {
	c.HTML(200, "register.html", nil)
}

func Regsave(c *gin.Context) {

	db := mysqldb.SetupDB()
	key := hex.EncodeToString(make([]byte, 32))
	name := c.Request.PostFormValue("username")
	email := c.Request.PostFormValue("emailid")
	mobilenumber := c.Request.PostFormValue("phonenumber")
	pwd := c.Request.PostFormValue("password")
	encryptpwd := encrypt(pwd, key)

	// currentTime := time.Now()
	createdate := time.Now().Format("Jan 2,2006 3:4:5 PM")
	var id, emailid int
	var firstname string
	_ = db.QueryRow("select count(email_id) from users where email_id=?", email).Scan(&emailid)
	log.Println("emailid is", emailid)
	if emailid < 1 {

		insForm, err := db.Prepare("INSERT INTO users(user_name, email_id,password,phone_number,created_date) VALUES(?,?,?,?,?)")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(name, email, encryptpwd, mobilenumber, createdate)
		log.Println("INSERT: Name: " + name + " | City: " + createdate + "|Password: " + encryptpwd)

		_ = db.QueryRow("select id,first_name from users where email_id=?", email).Scan(&id, &firstname)
		session, _ := store.Get(c.Request, "mysession")
		session.Values["name"] = firstname
		session.Values["userid"] = id
		log.Println("session name", firstname)
		session.Save(c.Request, c.Writer)

		last_login_date := time.Now().Format("Jan 2,2006 3:4:5 PM")
		updt, err := db.Prepare("UPDATE users SET last_login_date=? WHERE email_id=?")
		if err != nil {
			panic(err.Error())
		}
		updt.Exec(last_login_date, name)

		c.Redirect(301, "/index")

	} else {

		c.HTML(200, "register.html", gin.H{"error": "**Already have a mail id give another one"})
	}

}

func encrypt(stringToEncrypt string, keyString string) (encryptedString string) {

	//Since the key is in string, we need to convert decode it to bytes
	key, _ := hex.DecodeString(keyString)
	plaintext := []byte(stringToEncrypt)

	//Create a new Cipher Block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	//Create a new GCM - https://en.wikipedia.org/wiki/Galois/Counter_Mode
	//https://golang.org/pkg/crypto/cipher/#NewGCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	//Create a nonce. Nonce should be from GCM
	nonce := make([]byte, aesGCM.NonceSize())
	log.Println("nounce size is", nonce)
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}

	//Encrypt the data using aesGCM.Seal
	//Since we don't want to save the nonce somewhere else in this case, we add it as a prefix to the encrypted data. The first nonce argument in Seal is the prefix.
	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)
	return fmt.Sprintf("%x", ciphertext)
}

func decrypt(encryptedString string, keyString string) (decryptedString string) {

	key, _ := hex.DecodeString(keyString)
	enc, _ := hex.DecodeString(encryptedString)

	//Create a new Cipher Block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	//Create a new GCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	//Get the nonce size
	nonceSize := aesGCM.NonceSize()

	//Extract the nonce from the encrypted data
	nonce, ciphertext := enc[:nonceSize], enc[nonceSize:]
	log.Println("nonce,cipher", nonce, ciphertext)
	//Decrypt the data
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}

	return fmt.Sprintf("%s", plaintext)
}

func DeleteSession(c *gin.Context) {

	// C := http.Cookie{
	// 	Name:   "mysession",
	// 	MaxAge: -1}
	// http.SetCookie(c.Writer, &C)

	session, _ := store.Get(c.Request, "mysession")
	session.Values["userid"] = " "
	session.Values["name"] = " "
	session.Options.MaxAge = -1
	log.Println("logout session username", session.Values["userid"])
	log.Println("logout session userid", session.Values["name"])
	// delete(session.Values, "")
	// delete(session.Values, "userid")
	session.Save(c.Request, c.Writer)

	// fmt.Fprintf(c.Writer, "Sanjai")
	time.Sleep(1 * time.Second)
	c.HTML(200, "login.html", nil)
	// c.Redirect(301, "/")
}
