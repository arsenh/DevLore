document.addEventListener("DOMContentLoaded", function () {
    // Delete modal logic (existing code)
    const modal = document.getElementById("delete-modal");
    const openBtn = document.getElementById("open-delete-modal");
    const cancelBtn = document.getElementById("cancel-delete");
    const confirmBtn = document.getElementById("confirm-delete");
    const form = document.getElementById("delete-form");

    if (modal && openBtn && cancelBtn && confirmBtn && form) {
        openBtn.addEventListener("click", () => (modal.style.display = "flex"));
        cancelBtn.addEventListener(
            "click",
            () => (modal.style.display = "none"),
        );
        confirmBtn.addEventListener("click", () => form.submit());
        modal.addEventListener("click", (e) => {
            if (e.target === modal) modal.style.display = "none";
        });
    }

    // Registration form validation
    const registerForm = document.querySelector(
        'form[action="/auth/register"]',
    );
    if (!registerForm) return;

    const emailInput = document.getElementById("email");
    const fullNameInput = document.getElementById("full_name");
    const passwordInput = document.getElementById("password");
    const passwordConfirmInput = document.getElementById("password_confirm");

    // Helper function to show error
    function showError(input, message) {
        clearError(input);
        const error = document.createElement("span");
        error.className = "error-message";
        error.style.color = "#d32f2f";
        error.style.fontSize = "0.875rem";
        error.style.marginTop = "0.25rem";
        error.style.display = "block";
        error.textContent = message;
        input.parentElement.appendChild(error);
        input.style.borderColor = "#d32f2f";
    }

    // Helper function to clear error
    function clearError(input) {
        const error = input.parentElement.querySelector(".error-message");
        if (error) error.remove();
        input.style.borderColor = "";
    }

    // Email validation
    emailInput.addEventListener("blur", function () {
        const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
        if (!this.value.trim()) {
            showError(this, "Email is required");
        } else if (!emailRegex.test(this.value)) {
            showError(this, "Please enter a valid email address");
        } else {
            clearError(this);
        }
    });

    emailInput.addEventListener("input", function () {
        if (this.value.trim()) clearError(this);
    });

    // Full name validation
    fullNameInput.addEventListener("blur", function () {
        if (!this.value.trim()) {
            showError(this, "Full name is required");
        } else if (this.value.trim().length < 2) {
            showError(this, "Full name must be at least 2 characters");
        } else {
            clearError(this);
        }
    });

    fullNameInput.addEventListener("input", function () {
        if (this.value.trim()) clearError(this);
    });

    // Password validation
    passwordInput.addEventListener("blur", function () {
        if (!this.value) {
            showError(this, "Password is required");
        } else if (this.value.length < 8) {
            showError(this, "Password must be at least 8 characters");
        } else {
            clearError(this);
            // Re-validate password confirmation if it has a value
            if (passwordConfirmInput.value) {
                passwordConfirmInput.dispatchEvent(new Event("blur"));
            }
        }
    });

    passwordInput.addEventListener("input", function () {
        if (this.value) clearError(this);
        // Re-validate confirmation on password change
        if (passwordConfirmInput.value) {
            clearError(passwordConfirmInput);
        }
    });

    // Password confirmation validation
    passwordConfirmInput.addEventListener("blur", function () {
        if (!this.value) {
            showError(this, "Please confirm your password");
        } else if (this.value !== passwordInput.value) {
            showError(this, "Passwords do not match");
        } else {
            clearError(this);
        }
    });

    passwordConfirmInput.addEventListener("input", function () {
        if (this.value) clearError(this);
    });

    // Form submission validation
    registerForm.addEventListener("submit", function (e) {
        let isValid = true;

        // Validate email
        const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
        if (!emailInput.value.trim()) {
            showError(emailInput, "Email is required");
            isValid = false;
        } else if (!emailRegex.test(emailInput.value)) {
            showError(emailInput, "Please enter a valid email address");
            isValid = false;
        }

        // Validate full name
        if (!fullNameInput.value.trim()) {
            showError(fullNameInput, "Full name is required");
            isValid = false;
        } else if (fullNameInput.value.trim().length < 2) {
            showError(fullNameInput, "Full name must be at least 2 characters");
            isValid = false;
        }

        // Validate password
        if (!passwordInput.value) {
            showError(passwordInput, "Password is required");
            isValid = false;
        } else if (passwordInput.value.length < 8) {
            showError(passwordInput, "Password must be at least 8 characters");
            isValid = false;
        }

        // Validate password confirmation
        if (!passwordConfirmInput.value) {
            showError(passwordConfirmInput, "Please confirm your password");
            isValid = false;
        } else if (passwordConfirmInput.value !== passwordInput.value) {
            showError(passwordConfirmInput, "Passwords do not match");
            isValid = false;
        }

        if (!isValid) {
            e.preventDefault();
            // Focus on first invalid field
            const firstError = registerForm.querySelector(".error-message");
            if (firstError) {
                firstError.previousElementSibling.focus();
            }
        }
    });
});
