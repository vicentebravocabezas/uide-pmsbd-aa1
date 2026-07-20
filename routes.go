package main

import "net/http"

func registerRoutes() {
	http.Handle("GET /", http.FileServer(http.Dir("frontend")))

	http.HandleFunc("GET /api/categorias", handleCategorias)
	http.HandleFunc("GET /api/productos", handleProductos)
	http.HandleFunc("POST /api/productos", handleCrearProductos)
	http.HandleFunc("GET /api/productos/{id}", handleProducto)
	http.HandleFunc("PUT /api/productos/{id}", handleModificarProducto)
	http.HandleFunc("DELETE /api/productos/{id}", handleBorrarProducto)
}
