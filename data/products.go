package data

import (
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"time"

	"github.com/go-playground/validator"
)

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	Price       float32 `json:"price" validate:"gt=0"`
	SKU         string  `json:"sku" validate:"required,sku"`
	CreatedOn   string  `json:"-"`
	UpdatedOn   string  `json:"-"`
	DeletedOn   string  `json:"-"`
}

type Products []*Product

func (p *Product) Validate() error {
	validate := validator.New()

	validate.RegisterValidation("sku", ValidateSKU)
	return validate.Struct(p)

}

func ValidateSKU(fl validator.FieldLevel) bool {

	//sku is of format abc-sds-cxc

	re := regexp.MustCompile("[a-z]+-[a-z]+-[a-z]+")
	mathces := re.FindAllString(fl.Field().String(), -1)
	if len(mathces) != 1 {
		return false
	}
	return true
}

func (p *Products) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)

	return e.Encode(p)

}

func (p *Product) FromJSON(r io.Reader) error {

	d := json.NewDecoder(r)
	return d.Decode(p)

}

func GetProducts() Products {

	return ProductList
}

func AddProduct(p *Product) {
	p.ID = getNextId()
	ProductList = append(ProductList, p)

}

func getNextId() int {
	lp := ProductList[len(ProductList)-1]
	return lp.ID + 1

}

func UpdateProd(id int, p *Product) error {
	_, pos, err := FindProduct(id)
	if err != nil {

		return err
	}
	p.ID = id
	ProductList[pos] = p
	return nil
}

func DeleteProd(id int, p *Product) error {
	_, pos, err := FindProduct(id)
	if err != nil {
		return err
	}

	p.ID = id

	ProductList = append(ProductList[:pos], ProductList[pos+1:]...)
	return nil
}

var ErrorNotFound = fmt.Errorf("Product Not found")

func FindProduct(id int) (*Product, int, error) {
	for i, p := range ProductList {
		if p.ID == id {

			return p, i, nil
		}
	}
	return nil, -1, ErrorNotFound
}

var ProductList = []*Product{
	&Product{
		ID:          1,
		Name:        "Hot Coffee",
		Description: "Hot enough to get Hot",
		Price:       5.3,
		SKU:         "abc23",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
	&Product{
		ID:          2,
		Name:        "Cold Coffee",
		Description: "Frothy milky cold coffee",
		Price:       9.3,
		SKU:         "abc653",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
}
