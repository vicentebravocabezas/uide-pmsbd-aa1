package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
)

func MiddlewareHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache")
		next.ServeHTTP(w, r)
	})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func handleCategorias(w http.ResponseWriter, r *http.Request) {
	cats, err := listarCategorias()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, cats)
}

func handleProductos(w http.ResponseWriter, r *http.Request) {
	productos, err := listarProductos(r.URL.Query().Get("oferta") == "1")
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, productos)
}

func handleCrearProductos(w http.ResponseWriter, r *http.Request) {
	var p Producto
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		writeError(w, http.StatusBadRequest, "JSON inválido")
		return
	}
	if err := validarProducto(p); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	id, err := crearProducto(p)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	p.ID = int(id)
	writeJSON(w, http.StatusCreated, p)
}

func handleModificarProducto(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "id inválido")
		return
	}

	var p Producto
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		writeError(w, http.StatusBadRequest, "JSON inválido")
		return
	}
	if err := validarProducto(p); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	p.ID = id
	if err := actualizarProducto(p); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, p)
}

func handleBorrarProducto(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "id inválido")
		return
	}

	if err := eliminarProducto(id); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func handleProducto(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "id inválido")
		return
	}

	p, err := obtenerProducto(id)
	if err != nil {
		writeError(w, http.StatusNotFound, "producto no encontrado")
		return
	}
	writeJSON(w, http.StatusOK, p)
}

func validarProducto(p Producto) error {
	switch {
	case strings.TrimSpace(p.Nombre) == "":
		return errors.New("el nombre es obligatorio")
	case p.IDCategoria <= 0:
		return errors.New("la categoría es obligatoria")
	case p.Precio < 0:
		return errors.New("el precio no puede ser negativo")
	case p.Descuento < 0 || p.Descuento > 100:
		return errors.New("el descuento debe estar entre 0 y 100")
	}
	return nil
}
