package cart

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Cart struct {
	Items []CartItem `json:"items"`
}

type CartItem struct {
	ProductID uint64 `json:"productID"`
	Count     int    `json:"count"`
}

type CartDiscount struct {
	Percent float64 `json:"percent"`
}

type Product struct {
	ID    uint64  `json:"id"`
	Price float64 `json:"price"`
}

type ProductClient struct {
	Token             string
	ProductServiceURL string
}

func (p *ProductClient) GetProduct(id uint64) (Product, error) {
	url := fmt.Sprintf("%s/%s/%d?short=true", p.ProductServiceURL, "product", id)
	req, _ := http.NewRequest(http.MethodGet, url, http.NoBody)
	req.Header.Add("Authorization", p.Token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return Product{}, err
	}

	if res.StatusCode == http.StatusOK {
		body, _ := ioutil.ReadAll(res.Body)
		product := Product{}
		json.Unmarshal(body, &product)

		return product, nil
	}

	return Product{}, errors.New("Couldn't get the cart")
}

type ShoppingCartClient struct {
	Token          string
	CartServiceURL string
}

func (s *ShoppingCartClient) GetUserCart(userID uint64) (Cart, error) {
	url := fmt.Sprintf("%s/%s/%d/%s", s.CartServiceURL, "user", userID, "cart")
	req, _ := http.NewRequest(http.MethodGet, url, http.NoBody)
	req.Header.Add("Authorization", s.Token)

	res, _ := http.DefaultClient.Do(req)

	if res.StatusCode == http.StatusOK {
		body, _ := ioutil.ReadAll(res.Body)
		cart := Cart{}
		json.Unmarshal(body, &cart)

		return cart, nil
	}

	return Cart{}, errors.New("Couldn't get the cart")
}

func (s *ShoppingCartClient) GetUserCartDiscount(userID uint64) (CartDiscount, error) {
	url := fmt.Sprintf("%s/%s/%d/%s", s.CartServiceURL, "user", userID, "discount")
	req, _ := http.NewRequest(http.MethodGet, url, http.NoBody)
	req.Header.Add("Authorization", s.Token)

	res, _ := http.DefaultClient.Do(req)

	if res.StatusCode == http.StatusOK {
		body, _ := ioutil.ReadAll(res.Body)
		discount := CartDiscount{}
		json.Unmarshal(body, &discount)

		return discount, nil
	}

	return CartDiscount{}, errors.New("Couldn't get the cart discount")

}

type ShoppingCart struct {
	productClient      ProductClient
	shoppingCartClient ShoppingCartClient
}

func (s *ShoppingCart) GetTotalPrice(userID uint64) (float64, error) {
	cart, _ := s.shoppingCartClient.GetUserCart(userID)

	totalPrice := 0.0
	for _, item := range cart.Items {
		product, _ := s.productClient.GetProduct(item.ProductID)
		totalPrice += product.Price * float64(item.Count)
	}

	cartDiscount, _ := s.shoppingCartClient.GetUserCartDiscount(userID)

	return totalPrice * cartDiscount.Percent, nil
}
