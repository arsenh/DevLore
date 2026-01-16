function debounce(fn, delay) {
  let timeoutId;
  return function (...args) {
    clearTimeout(timeoutId);
    timeoutId = setTimeout(() => {
      fn.apply(this, args);
    }, delay);
  };
}

function emailExistsFetch(email) {
  if (typeof email != "string" || email.length < 1)
    return Promise.resolve({ exists: false });

  const url = `/email-exists?email=${encodeURIComponent(email)}`;

  return fetch(url)
    .then((response) => {
      if (!response.ok) throw new Error("Network response was not ok");
      return response.json();
    })
    .catch((error) => {
      console.error("Error checking email:", error);
      return { exists: false }; // Fail gracefully
    });
}

document.addEventListener("DOMContentLoaded", function () {
  // Delete modal logic (existing code)
  const modal = document.getElementById("delete-modal");
  const openBtn = document.getElementById("open-delete-modal");
  const cancelBtn = document.getElementById("cancel-delete");
  const confirmBtn = document.getElementById("confirm-delete");
  const form = document.getElementById("delete-form");

  if (modal && openBtn && cancelBtn && confirmBtn && form) {
    openBtn.addEventListener("click", () => (modal.style.display = "flex"));
    cancelBtn.addEventListener("click", () => (modal.style.display = "none"));
    confirmBtn.addEventListener("click", () => form.submit());
    modal.addEventListener("click", (e) => {
      if (e.target === modal) modal.style.display = "none";
    });
  }

  // Helper functions (shared by both forms)
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

  function clearError(input) {
    const error = input.parentElement.querySelector(".error-message");
    if (error) error.remove();
    input.style.borderColor = "";
  }

  // Registration form validation
  const registerForm = document.querySelector('form[action="/auth/register"]');
  if (registerForm) {
    const emailInput = document.getElementById("email");
    const fullNameInput = document.getElementById("full_name");
    const passwordInput = document.getElementById("password");
    const passwordConfirmInput = document.getElementById("password_confirm");

    // Track validation states
    let emailCheckInProgress = false;
    let emailExists = false;
    let fieldValidation = {
      email: false,
      fullName: false,
      password: false,
      passwordConfirm: false,
    };

    // Validation functions
    function validateEmail(showErrorMsg = true) {
      const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;

      if (!emailInput.value.trim()) {
        if (showErrorMsg) showError(emailInput, "Email is required");
        fieldValidation.email = false;
        return false;
      } else if (!emailRegex.test(emailInput.value)) {
        if (showErrorMsg)
          showError(emailInput, "Please enter a valid email address");
        fieldValidation.email = false;
        return false;
      } else if (emailExists) {
        if (showErrorMsg) showError(emailInput, "This email is already in use");
        fieldValidation.email = false;
        return false;
      } else {
        if (showErrorMsg) clearError(emailInput);
        fieldValidation.email = true;
        return true;
      }
    }

    function validateFullName(showErrorMsg = true) {
      if (!fullNameInput.value.trim()) {
        if (showErrorMsg) showError(fullNameInput, "Full name is required");
        fieldValidation.fullName = false;
        return false;
      } else if (fullNameInput.value.trim().length < 2) {
        if (showErrorMsg)
          showError(fullNameInput, "Full name must be at least 2 characters");
        fieldValidation.fullName = false;
        return false;
      } else {
        if (showErrorMsg) clearError(fullNameInput);
        fieldValidation.fullName = true;
        return true;
      }
    }

    function validatePassword(showErrorMsg = true) {
      if (!passwordInput.value) {
        if (showErrorMsg) showError(passwordInput, "Password is required");
        fieldValidation.password = false;
        return false;
      } else if (passwordInput.value.length < 8) {
        if (showErrorMsg)
          showError(passwordInput, "Password must be at least 8 characters");
        fieldValidation.password = false;
        return false;
      } else {
        if (showErrorMsg) clearError(passwordInput);
        fieldValidation.password = true;
        return true;
      }
    }

    function validatePasswordConfirm(showErrorMsg = true) {
      if (!passwordConfirmInput.value) {
        if (showErrorMsg)
          showError(passwordConfirmInput, "Please confirm your password");
        fieldValidation.passwordConfirm = false;
        return false;
      } else if (passwordConfirmInput.value !== passwordInput.value) {
        if (showErrorMsg)
          showError(passwordConfirmInput, "Passwords do not match");
        fieldValidation.passwordConfirm = false;
        return false;
      } else {
        if (showErrorMsg) clearError(passwordConfirmInput);
        fieldValidation.passwordConfirm = true;
        return true;
      }
    }

    // Email validation on blur
    emailInput.addEventListener("blur", function () {
      validateEmail(true);
    });

    // Clear error on input (when user starts typing)
    emailInput.addEventListener("input", function () {
      if (this.value.trim()) {
        clearError(this);
        emailExists = false; // Reset the flag when user types
      }
      validateEmail(false); // Validate silently to update state
    });

    // Debounced email existence check
    const debouncedFetch = debounce((email) => {
      const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;

      // Only check if email format is valid
      if (!emailRegex.test(email)) return;

      emailCheckInProgress = true;
      emailExistsFetch(email).then((status) => {
        emailCheckInProgress = false;

        // Only show error if the email value hasn't changed
        if (emailInput.value === email) {
          if (status.exists) {
            emailExists = true;
            showError(emailInput, "This email is already in use");
            fieldValidation.email = false;
          } else {
            emailExists = false;
            // Only clear if there's an "email exists" error
            const errorMsg =
              emailInput.parentElement.querySelector(".error-message");
            if (
              errorMsg &&
              errorMsg.textContent === "This email is already in use"
            ) {
              clearError(emailInput);
              validateEmail(false); // Revalidate silently
            }
          }
        }
      });
    }, 500);

    // Use 'input' event for email existence check
    emailInput.addEventListener("input", function () {
      debouncedFetch(this.value);
    });

    // Full name validation
    fullNameInput.addEventListener("blur", function () {
      validateFullName(true);
    });

    fullNameInput.addEventListener("input", function () {
      if (this.value.trim()) clearError(this);
      validateFullName(false); // Validate silently to update state
    });

    // Password validation
    passwordInput.addEventListener("blur", function () {
      validatePassword(true);
      // Re-validate password confirmation if it has a value
      if (passwordConfirmInput.value) {
        validatePasswordConfirm(true);
      }
    });

    passwordInput.addEventListener("input", function () {
      if (this.value) clearError(this);
      validatePassword(false); // Validate silently to update state
      // Re-validate confirmation on password change
      if (passwordConfirmInput.value) {
        validatePasswordConfirm(false);
      }
    });

    // Password confirmation validation
    passwordConfirmInput.addEventListener("blur", function () {
      validatePasswordConfirm(true);
    });

    passwordConfirmInput.addEventListener("input", function () {
      if (this.value) clearError(this);
      validatePasswordConfirm(false); // Validate silently to update state
    });

    // Form submission validation
    registerForm.addEventListener("submit", function (e) {
      e.preventDefault(); // Always prevent default first

      // Check if email check is still in progress
      if (emailCheckInProgress) {
        showError(emailInput, "Please wait while we verify your email");
        return;
      }

      // Validate all fields and show errors
      const emailValid = validateEmail(true);
      const fullNameValid = validateFullName(true);
      const passwordValid = validatePassword(true);
      const passwordConfirmValid = validatePasswordConfirm(true);

      const isValid =
        emailValid && fullNameValid && passwordValid && passwordConfirmValid;

      if (isValid) {
        // All fields are valid, submit the form
        this.submit();
      } else {
        // Focus on first invalid field
        const firstError = registerForm.querySelector(".error-message");
        if (firstError) {
          firstError.previousElementSibling.focus();
        }
      }
    });
  }

  // Login form validation
  const loginForm = document.querySelector('form[action="/auth/login"]');
  if (loginForm) {
    const loginEmailInput = document.getElementById("email");
    const loginPasswordInput = document.getElementById("password");

    // Add real-time validation for login form
    if (loginEmailInput) {
      loginEmailInput.addEventListener("input", function () {
        if (this.value.trim()) clearError(this);
      });
    }

    if (loginPasswordInput) {
      loginPasswordInput.addEventListener("input", function () {
        if (this.value) clearError(this);
      });
    }

    loginForm.addEventListener("submit", function (e) {
      e.preventDefault(); // Always prevent default first

      let isValid = true;

      // Validate email
      const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
      if (!loginEmailInput.value.trim()) {
        showError(loginEmailInput, "Email is required");
        isValid = false;
      } else if (!emailRegex.test(loginEmailInput.value)) {
        showError(loginEmailInput, "Please enter a valid email address");
        isValid = false;
      }

      // Validate password
      if (!loginPasswordInput.value) {
        showError(loginPasswordInput, "Password is required");
        isValid = false;
      }

      if (isValid) {
        // All fields are valid, submit the form
        this.submit();
      } else {
        // Focus on first invalid field
        const firstError = loginForm.querySelector(".error-message");
        if (firstError) {
          firstError.previousElementSibling.focus();
        }
      }
    });
  }
});
