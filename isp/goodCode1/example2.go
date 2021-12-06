package goodCode1

import (
	"context"
	"errors"
)

type Money struct {
	// какие-то поля
}

type Product struct {
	// какие-то поля
}

type Wallet struct {
	// какие-то поля
}

func (w *Wallet) Deduct(money Money) error {
	return nil
}

type DiscountPolicy struct {
	// какие-то поля
}

func (d *DiscountPolicy) IsApplicableFor(customer *PremiumCustomer, product Product) bool {
	return true
}

type ShoppingCart struct {
}

func (s ShoppingCart) Add(product Product) {

}

type User interface {
	AddToShoppingCart(product Product)
	//
	// некоторые дополнительные методы
	//
}

type LoggedInUser interface {
	User
	Pay(money Money) error
	//
	// некоторые дополнительные методы
	//
}

type PremiumUser interface {
	LoggedInUser
	HasDiscountFor(product Product) bool
	//
	// некоторые дополнительные методы
	//
}

type Guest struct {
	cart ShoppingCart
	//
	// некоторые дополнительные поля
	//
}

func (g *Guest) AddToShoppingCart(product Product) {
	g.cart.Add(product)
}

type NormalCustomer struct {
	cart   ShoppingCart
	wallet Wallet
	//
	// некоторые дополнительные поля
	//
}

func (c *NormalCustomer) AddToShoppingCart(product Product) {
	c.cart.Add(product)
}

func (c *NormalCustomer) Pay(money Money) error {
	return c.wallet.Deduct(money)
}

type PremiumCustomer struct {
	cart     ShoppingCart
	wallet   Wallet
	policies []DiscountPolicy
	//
	// некоторые дополнительные поля
	//
}

func (c *PremiumCustomer) AddToShoppingCart(product Product) {
	c.cart.Add(product)
}

func (c *PremiumCustomer) Pay(money Money) error {
	return c.wallet.Deduct(money)
}

func (c *PremiumCustomer) HasDiscountFor(product Product) bool {
	for _, p := range c.policies {
		if p.IsApplicableFor(c, product) {
			return true
		}
	}

	return false
}

type UserService struct {
	//
	// какие-то поля
	//
}

func (u *UserService) Checkout(ctx context.Context, user User, product Product) error {
	loggedIn, ok := user.(LoggedInUser)
	if !ok {
		return errors.New("user is not logged in")
	}

	var money Money
	//
	// какие-то вычисления
	//
	if premium, ok := loggedIn.(PremiumUser); ok && premium.HasDiscountFor(product) {
		//
		// применить скидку
		//
	}

	return loggedIn.Pay(money)
}
