package models

type Product struct {
	ID            int
	Title         string
	Description   string
	Price         float32
	Image         string
	Imagepath     string
	Overallrating []int
	Whishlist     int
	Created       string
	Modified      string
	Name          string
	Reviewcount   int
	Remain        []int
	Randomid      int
	Value         int
}
