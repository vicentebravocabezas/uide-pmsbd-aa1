package main

import (
	"database/sql"
	_ "embed"
	"math"

	_ "modernc.org/sqlite"
)

var db *sql.DB

type Categoria struct {
	ID          int    `json:"id_categoria"`
	Nombre      string `json:"nombre"`
	Descripcion string `json:"descripcion"`
}

type Producto struct {
	ID           int     `json:"id_producto"`
	IDCategoria  int     `json:"id_categoria"`
	Categoria    string  `json:"categoria"`
	Nombre       string  `json:"nombre"`
	Descripcion  string  `json:"descripcion"`
	Precio       float64 `json:"precio"`
	Imagen       string  `json:"imagen"`
	EsOfertaMes  bool    `json:"es_oferta_mes"`
	Descuento    float64 `json:"descuento"`
	PrecioOferta float64 `json:"precio_oferta"`
}

func listarCategorias() ([]Categoria, error) {
	rows, err := db.Query(`SELECT id_categoria, nombre, descripcion FROM categoria ORDER BY nombre`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Categoria
	for rows.Next() {
		var c Categoria
		if err := rows.Scan(&c.ID, &c.Nombre, &c.Descripcion); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

func scanProducto(row interface{ Scan(...any) error }) (Producto, error) {
	var p Producto
	var oferta int
	err := row.Scan(&p.ID, &p.IDCategoria, &p.Categoria, &p.Nombre, &p.Descripcion,
		&p.Precio, &p.Imagen, &oferta, &p.Descuento)
	p.EsOfertaMes = oferta == 1
	p.PrecioOferta = math.Round(p.Precio*(1-p.Descuento/100)*100) / 100
	return p, err
}

func listarProductos(soloOferta bool) ([]Producto, error) {
	query := `SELECT p.id_producto, p.id_categoria, c.nombre, p.nombre, p.descripcion, p.precio, p.imagen, p.es_oferta_mes, p.descuento
FROM producto p
INNER JOIN categoria c ON c.id_categoria = p.id_categoria`
	if soloOferta {
		query += ` WHERE p.es_oferta_mes = 1`
	}

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := []Producto{}
	for rows.Next() {
		var p Producto
		var oferta int
		if err := rows.Scan(&p.ID, &p.IDCategoria, &p.Categoria, &p.Nombre, &p.Descripcion, &p.Precio, &p.Imagen, &oferta, &p.Descuento); err != nil {
			return nil, err
		}

		p.EsOfertaMes = oferta == 1
		p.PrecioOferta = math.Round(p.Precio*(1-p.Descuento/100)*100) / 100
		out = append(out, p)
	}

	return out, rows.Err()
}

func obtenerProducto(id int) (Producto, error) {
	row := db.QueryRow(`SELECT p.id_producto, p.id_categoria, c.nombre, p.nombre, p.descripcion, p.precio, p.imagen, p.es_oferta_mes, p.descuento
FROM producto p
INNER JOIN categoria c ON c.id_categoria = p.id_categoria WHERE p.id_producto = ?`, id)

	var p Producto

	var oferta int

	if err := row.Scan(&p.ID, &p.IDCategoria, &p.Categoria, &p.Nombre, &p.Descripcion, &p.Precio, &p.Imagen, &oferta, &p.Descuento); err != nil {
		return Producto{}, err
	}

	p.EsOfertaMes = oferta == 1
	p.PrecioOferta = math.Round(p.Precio*(100-p.Descuento)) / 100
	return p, nil
}

func crearProducto(p Producto) (int64, error) {
	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	if p.EsOfertaMes {
		if _, err := tx.Exec(`UPDATE producto SET es_oferta_mes = 0, descuento = 0 WHERE es_oferta_mes = 1`); err != nil {
			return 0, err
		}
	}
	res, err := tx.Exec(
		`INSERT INTO producto (id_categoria, nombre, descripcion, precio, imagen, es_oferta_mes, descuento)
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		p.IDCategoria, p.Nombre, p.Descripcion, p.Precio, p.Imagen, p.EsOfertaMes, p.Descuento)
	if err != nil {
		return 0, err
	}
	if err := tx.Commit(); err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func actualizarProducto(p Producto) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if p.EsOfertaMes {
		if _, err := tx.Exec(`UPDATE producto SET es_oferta_mes = 0, descuento = 0 WHERE es_oferta_mes = 1 AND id_producto != ?`, p.ID); err != nil {
			return err
		}
	}
	_, err = tx.Exec(
		`UPDATE producto
		 SET id_categoria = ?, nombre = ?, descripcion = ?, precio = ?, imagen = ?, es_oferta_mes = ?, descuento = ?
		 WHERE id_producto = ?`,
		p.IDCategoria, p.Nombre, p.Descripcion, p.Precio, p.Imagen, p.EsOfertaMes, p.Descuento, p.ID)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func eliminarProducto(id int) error {
	_, err := db.Exec(`DELETE FROM producto WHERE id_producto = ?`, id)
	return err
}
