import $ from "./jquery-4.0.0.module.min.js";

const catalogo = document.querySelector("#catalogo");
const form = document.querySelector("#form-producto");
const tabla = document.querySelector("#tabla-productos tbody");

const cardProducto = (p) => {
  const card = document.createElement("article");
  card.className = "product-card";

  const img = document.createElement("img");
  img.src = `/media/${p.imagen}`;
  img.loading = "lazy";

  const body = document.createElement("div");
  body.className = "product-body";

  const titulo = document.createElement("h4");
  titulo.textContent = p.nombre;

  const desc = document.createElement("p");
  desc.className = "product-desc";
  desc.textContent = p.descripcion;

  const precio = document.createElement("p");
  precio.innerHTML = p.es_oferta_mes
    ? `<span class="price-old">$${p.precio.toFixed(2)}</span> <span class="price">$${p.precio_oferta.toFixed(2)}</span> <span class="badge-oferta">-${p.descuento}%</span>`
    : `<span class="price">$${p.precio.toFixed(2)}</span>`;

  body.append(titulo, desc, precio);
  card.append(img, body);
  return card;
};

const cargarCatalogo = () => {
  $.getJSON("/api/productos", (productos) => {
    catalogo.replaceChildren(); // para no duplicar elementos al cargar nuevamente

    // agrupar los productos por categoría
    const grupos = {};
    productos.forEach((p) => {
      const g = grupos[p.categoria];

      if (g) {
        g.push(p);
      } else {
        grupos[p.categoria] = [p];
      }
    });

    Object.entries(grupos).forEach(([categoria, items]) => {
      const seccion = document.createElement("section");
      seccion.className = "catalog-category";

      const h3 = document.createElement("h3");
      h3.textContent = categoria;

      const grid = document.createElement("div");
      grid.className = "product-grid";
      items.forEach((p) => grid.append(cardProducto(p)));

      seccion.append(h3, grid);
      catalogo.append(seccion);
    });
  });
};

// tabla de administración
const filaProducto = (p) => {
  const fila = document.createElement("tr");

  const celdaNombre = document.createElement("td");
  celdaNombre.textContent = p.nombre;
  fila.append(celdaNombre);

  const celdaCategoria = document.createElement("td");
  celdaCategoria.textContent = p.categoria;
  fila.append(celdaCategoria);

  const celdaPrecio = document.createElement("td");
  celdaPrecio.textContent = `$${p.precio.toFixed(2)}`;
  fila.append(celdaPrecio);

  const celdaOferta = document.createElement("td");
  celdaOferta.textContent = p.es_oferta_mes ? `Sí (-${p.descuento}%)` : "No";
  fila.append(celdaOferta);

  const acciones = document.createElement("td");

  const btnEditar = document.createElement("button");
  btnEditar.className = "btn-small btn-edit";
  btnEditar.textContent = "Editar";

  btnEditar.addEventListener("click", () => {
    document.querySelector("#p-id").value = p.id_producto;
    document.querySelector("#p-nombre").value = p.nombre;
    document.querySelector("#p-categoria").value = p.id_categoria;
    document.querySelector("#p-precio").value = p.precio;
    document.querySelector("#p-imagen").value = p.imagen;
    document.querySelector("#p-descuento").value = p.descuento;
    document.querySelector("#p-oferta").checked = p.es_oferta_mes;
    document.querySelector("#p-descripcion").value = p.descripcion;
    document.querySelector("#btn-cancelar").classList.remove("hidden");
    form.scrollIntoView({ behavior: "smooth", block: "center" });
  });

  const btnEliminar = document.createElement("button");
  btnEliminar.className = "btn-small btn-delete";
  btnEliminar.textContent = "Eliminar";
  btnEliminar.addEventListener("click", () => eliminar(p));

  acciones.append(btnEditar, " ", btnEliminar);
  fila.append(acciones);
  return fila;
};

const cargarTabla = () => {
  $.getJSON("/api/productos", (productos) => {
    tabla.replaceChildren();
    productos.forEach((p) => tabla.append(filaProducto(p)));
  });
};

const recargar = () => {
  cargarCatalogo();
  cargarTabla();
};

// operaciones CRUD
const eliminar = (p) => {
  if (!confirm(`¿Eliminar el producto "${p.nombre}"?`)) return;
  $.ajax({
    url: `/api/productos/${p.id_producto}`,
    method: "DELETE",
    success: recargar,
  });
};

form.addEventListener("submit", (e) => {
  e.preventDefault(); // evitar que el navegador envie el formulario normalmente para poder hacerlo con jquery

  const id = document.querySelector("#p-id").value;
  const producto = {
    nombre: document.querySelector("#p-nombre").value.trim(),
    id_categoria: Number(document.querySelector("#p-categoria").value),
    precio: Number(document.querySelector("#p-precio").value),
    imagen:
      document.querySelector("#p-imagen").value.trim() ||
      "logo-horizontal-white.svg",
    descripcion: document.querySelector("#p-descripcion").value.trim(),
    es_oferta_mes: document.querySelector("#p-oferta").checked,
    descuento: Number(document.querySelector("#p-descuento").value) || 0,
  };

  $.ajax({
    url: id ? `/api/productos/${id}` : "/api/productos",
    method: id ? "PUT" : "POST",
    contentType: "application/json",
    data: JSON.stringify(producto),
    success: () => {
      form.reset();
      document.querySelector("#p-id").value = "";
      document.querySelector("#btn-cancelar").classList.add("hidden");
      recargar();
    },
  });
});

document.querySelector("#btn-cancelar").addEventListener("click", () => {
  form.reset();
  document.querySelector("#p-id").value = "";
  document.querySelector("#btn-cancelar").classList.add("hidden");
});

$.getJSON("/api/categorias", (categorias) => {
  const select = document.querySelector("#p-categoria");
  select.append(new Option("— Selecciona una categoría —", ""));
  categorias.forEach((c) =>
    select.append(new Option(c.nombre, c.id_categoria)),
  );
});

recargar();
