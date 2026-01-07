document.addEventListener("DOMContentLoaded", function () {
    const modal = document.getElementById("delete-modal");
    const openBtn = document.getElementById("open-delete-modal");
    const cancelBtn = document.getElementById("cancel-delete");
    const confirmBtn = document.getElementById("confirm-delete");
    const form = document.getElementById("delete-form");

    // Guard: stop if any element is missing (for pages without delete modal)
    if (!modal || !openBtn || !cancelBtn || !confirmBtn || !form) return;

    openBtn.addEventListener("click", () => modal.style.display = "flex");
    cancelBtn.addEventListener("click", () => modal.style.display = "none");
    confirmBtn.addEventListener("click", () => form.submit());

    modal.addEventListener("click", (e) => {
        if (e.target === modal) modal.style.display = "none";
    });
});
