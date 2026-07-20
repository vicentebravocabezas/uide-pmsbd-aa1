import $ from "./jquery-4.0.0.module.min.js";
// Granito Coffee Shop — página Oferta del mes
// Carga el producto en oferta desde la base de datos (AJAX con jQuery)
// y lo renderiza con JavaScript moderno.

const contenedor = document.querySelector("#oferta");

$.getJSON("/api/productos?oferta=1", (productos) => {
  const p = productos[0];

  const panel = document.createElement("div");
  panel.className = "offer-panel";

  const img = document.createElement("img");
  img.src = `media/${p.imagen}`;
  img.alt = p.nombre;

  const body = document.createElement("div");
  body.className = "offer-body";
  body.innerHTML = `
    <h3>${p.nombre}</h3>
    <p class="offer-desc">${p.descripcion}</p>
    <p class="offer-price">
      <span class="price-old">$${p.precio.toFixed(2)}</span>
      <span class="price">$${p.precio_oferta.toFixed(2)}</span>
    </p>`;

  const cualidades = [
    "Tostado en casa cada semana, frescura garantizada.",
    "Grano 100% arábica de altura, de comercio directo.",
    "Molienda a pedido según tu método de preparación.",
    "Disponible en el local y para pedidos al por mayor.",
  ];
  const lista = document.createElement("ul");
  cualidades.forEach((c) => {
    const li = document.createElement("li");
    li.textContent = c;
    lista.append(li);
  });
  body.append(lista);

  panel.append(img, body);
  contenedor.append(panel);
});
