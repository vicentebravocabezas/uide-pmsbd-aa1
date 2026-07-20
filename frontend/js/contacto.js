const contactForm = document.getElementById("contact-form");
if (contactForm) {
  const showError = (inputElement, errorId, message) => {
    inputElement.classList.add("invalid");
    document.getElementById(errorId).textContent = message;
  };

  contactForm.addEventListener("submit", (e) => {
    // prevenir que el formulario se resetee al presionar el boton de submit
    e.preventDefault();

    // limpiar errores
    const inputs = document.querySelectorAll(
      ".form-group input, .form-group select, .form-group textarea",
    );
    inputs.forEach((i) => {
      i.classList.remove("invalid");
    });

    const errors = document.querySelectorAll(".error-message");
    errors.forEach((e) => {
      e.textContent = "";
    });

    const name = document.getElementById("name");
    const email = document.getElementById("email");
    const subject = document.getElementById("subject");
    const message = document.getElementById("message");

    let isValid = true;

    // validar nombre
    if (name.value.trim() === "") {
      showError(name, "name-error", "El nombre es obligatorio.");
      isValid = false;
    }

    // validar email
    if (email.value.trim() === "") {
      showError(email, "email-error", "El correo es obligatorio.");
      isValid = false;
      // regex para validar email
    } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email.value.trim())) {
      showError(email, "email-error", "Ingresa un correo electrónico válido.");
      isValid = false;
    }

    // validar asunto
    if (subject.value === "") {
      showError(subject, "subject-error", "Selecciona un asunto.");
      isValid = false;
    }

    // validar mensaje
    if (message.value.trim() === "") {
      showError(message, "message-error", "El mensaje es obligatorio.");
      isValid = false;
    } else if (message.value.trim().length < 20) {
      showError(
        message,
        "message-error",
        "El mensaje debe tener al menos 20 caracteres.",
      );
      isValid = false;
    }

    // mostrar mensaje de exito
    if (isValid) {
      contactForm.style.display = "none";
      document.getElementById("success-message").classList.remove("hidden");
    }
  });
}
