import $ from "./jquery-4.0.0.module.min.js";

// Carrusel que muestra dinámicamente un producto cada 3 segundos
const carousel = document.querySelector("section.carousel");
if (carousel) {
  $.getJSON("/api/productos", (productos) => {
    if (!productos.length) return;
    const titulo = carousel.querySelector("h3");
    const descripcion = carousel.querySelector("p");
    let idx = 0;

    const actualizar = () => {
      titulo.textContent = productos[idx].nombre;
      descripcion.textContent = productos[idx].descripcion;
      idx = (idx + 1) % productos.length;
    };

    actualizar();
    setInterval(actualizar, 3000);
  });
}
