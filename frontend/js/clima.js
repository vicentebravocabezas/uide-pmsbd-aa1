import $ from "./jquery-4.0.0.module.min.js";

const btn = document.getElementById("btn-clima");
const result = document.getElementById("clima-result");

btn.addEventListener("click", () => {
  $.getJSON("/api/clima", (clima) => {
    result.classList.remove("hidden");
    document.getElementById("clima-temp").textContent =
      `${Number(clima.temperatura)} °C`;
    document.getElementById("clima-humedad").textContent = `${clima.humedad} %`;
    document.getElementById("clima-desc").textContent = clima.descripcion;
    document.getElementById("clima-viento").textContent =
      `${Number(clima.viento)} m/s`;
  });
});
