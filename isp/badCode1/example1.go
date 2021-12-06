package badCode1

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

type User interface {
	AddToShoppingCart(product Product)
	IsLoggedIn() bool
	Pay(money Money) error
	HasPremium() bool
	HasDiscountFor(product Product) bool
	//
	// некоторые дополнительные методы
	//
}

type ShoppingCart struct {

}

func (s ShoppingCart) Add(product Product) {

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

func (g *Guest) IsLoggedIn() bool {
	return false
}

func (g *Guest) Pay(Money) error {
	return errors.New("user is not logged in")
}

func (g *Guest) HasPremium() bool {
	return false
}

func (g *Guest) HasDiscountFor(product Product) bool {
	return false
}

type NormalCustomer struct {
	cart ShoppingCart
	wallet Wallet
	//
	// некоторые дополнительные поля
	//
}

func (c *NormalCustomer) AddToShoppingCart(product Product) {
	c.cart.Add(product)
}

func (c *NormalCustomer) IsLoggedIn() bool {
	return true
}

func (c *NormalCustomer) Pay(money Money) error {
	return c.wallet.Deduct(money)
}

func (c *NormalCustomer) HasPremium() bool {
	return false
}

func (c *NormalCustomer) HasDiscountFor(product Product) bool {
	return false
}

type PremiumCustomer struct {
	cart ShoppingCart
	wallet Wallet
	policies []DiscountPolicy
	//
	// некоторые дополнительные поля
	//
}

func (c *PremiumCustomer) AddToShoppingCart(product Product) {
	c.cart.Add(product)
}

func (c *PremiumCustomer) IsLoggedIn() bool {
	return true
}

func (c *PremiumCustomer) Pay(money Money) error {
	return c.wallet.Deduct(money)
}

func (c *PremiumCustomer) HasPremium() bool {
	return true
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
	if !user.IsLoggedIn() {
		return errors.New("user is not logged in")
	}

	var money Money
	//
	// какие-то вычисления
	//
	if user.HasDiscountFor(product) {
		//
		// применить скидку
		//
	}
	return user.Pay(money)
}