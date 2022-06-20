package controllers

import "github.com/arvinpaundra/ecommerce-api/api/middlewares"

func (s *Server) InitializeRoutes() {
	// Home route
	s.Router.HandleFunc("/", middlewares.SetMiddlewareJSON(s.Home)).Methods("Get")

	// Auth routes
	s.Router.HandleFunc("/login", middlewares.SetMiddlewareJSON(s.Login)).Methods("POST")
	s.Router.HandleFunc("/register", middlewares.SetMiddlewareJSON(s.Register)).Methods("POST")

	// Products routes
	s.Router.HandleFunc("/products", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.AddProduct))).Methods("POST")
	s.Router.HandleFunc("/products", middlewares.SetMiddlewareJSON(s.GetAllProducts)).Methods("GET")
	s.Router.HandleFunc("/products/{id}", middlewares.SetMiddlewareJSON(s.GetSingleProduct)).Methods("GET")
	s.Router.HandleFunc("/products/{id}/edit", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdateProduct))).Methods("PUT")
	s.Router.HandleFunc("/products/{id}/delete", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.DeleteProduct))).Methods("DELETE")

	// Categories routes
	s.Router.HandleFunc("/categories", middlewares.SetMiddlewareJSON(s.AddCategory)).Methods("POST")
	s.Router.HandleFunc("/categories", middlewares.SetMiddlewareJSON(s.GetAllCategories)).Methods("GET")
	s.Router.HandleFunc("/categories/{id}", middlewares.SetMiddlewareJSON(s.GetSingleCategory)).Methods("GET")
	s.Router.HandleFunc("/categories/{id}/edit", middlewares.SetMiddlewareJSON(s.UpdateCategory)).Methods("PUT")
	s.Router.HandleFunc("/categories/{id}/delete", middlewares.SetMiddlewareJSON(s.DeleteCategory)).Methods("DELETE")

	// Payments routes
	s.Router.HandleFunc("/payments", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.AddPayment))).Methods("POST")
	s.Router.HandleFunc("/payments", middlewares.SetMiddlewareJSON(s.GetAllPayments)).Methods("GET")
	s.Router.HandleFunc("/payments/{id}/edit", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdatePayment))).Methods("PUT")
	s.Router.HandleFunc("/payments/{id}/delete", middlewares.SetMiddlewareJSON(s.DeletePayment)).Methods("DELETE")

	// Carts routes
	s.Router.HandleFunc("/customers/{id}/carts", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.AddToCart))).Methods("POST")
	s.Router.HandleFunc("/customers/{id}/carts", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.GetCustomerCarts))).Methods("GET")
	s.Router.HandleFunc("/customers/{id}/carts/{cartId}/delete", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.DeleteCustomerCart))).Methods("DELETE")

	// Checkouts routes
	s.Router.HandleFunc("/customers/{id}/checkouts", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.CustomerCreateCheckout))).Methods("POST")
	s.Router.HandleFunc("/customers/{id}/checkouts", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.GetCustomerCheckouts))).Methods("GET")
}
