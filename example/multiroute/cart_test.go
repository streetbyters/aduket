package cart

import (
	"github.com/streetbyters/aduket"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTotalPrice(t *testing.T) {
	productServer, productServerRequestRecorder := aduket.NewServer(
		http.MethodGet, "/product/:productid",
		aduket.StatusCode(http.StatusOK),
		aduket.JSONBody(Product{ID: 123, Price: 100.0}),
	)
	defer productServer.Close()

	cartRoute := aduket.Route{http.MethodGet, "/user/:userid/cart"}
	discountRoute := aduket.Route{http.MethodGet, "/user/:userid/discount"}

	cartServer, cartServerRequestRecorder := aduket.NewMultiRouteServer(
		map[aduket.Route][]aduket.ResponseRuleOption{
			cartRoute: {
				aduket.StatusCode(http.StatusOK),
				aduket.JSONBody(
					Cart{Items: []CartItem{
						{ProductID: 123, Count: 2},
					}},
				),
			},
			discountRoute: {
				aduket.StatusCode(http.StatusOK),
				aduket.JSONBody(CartDiscount{Percent: 0.5}),
			},
		},
	)
	defer cartServer.Close()

	token := "testToken"

	productClient := ProductClient{Token: token, ProductServiceURL: productServer.URL}
	shoppingCartClient := ShoppingCartClient{Token: token, CartServiceURL: cartServer.URL}

	shoppingCart := ShoppingCart{productClient: productClient, shoppingCartClient: shoppingCartClient}

	actualPrice, err := shoppingCart.GetTotalPrice(111)

	authHeader := http.Header{}
	authHeader.Add("Authorization", token)

	cartServerRequestRecorder[cartRoute].AssertParamEqual(t, "userid", "111")
	cartServerRequestRecorder[cartRoute].AssertHeaderEqual(t, authHeader)

	cartServerRequestRecorder[discountRoute].AssertParamEqual(t, "userid", "111")
	cartServerRequestRecorder[discountRoute].AssertHeaderEqual(t, authHeader)

	productServerRequestRecorder.AssertParamEqual(t, "productid", "123")
	productServerRequestRecorder.AssertQueryParamEqual(t, "short", []string{"true"})
	productServerRequestRecorder.AssertHeaderEqual(t, authHeader)

	assert.Nil(t, err)
	assert.Equal(t, 100.0, actualPrice)
}
