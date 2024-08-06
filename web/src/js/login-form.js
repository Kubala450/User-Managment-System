document.addEventListener("DOMContentLoaded", function () {
  document
    .getElementById("loginForm")
    .addEventListener("submit", async function (event) {
      event.preventDefault();

      const username = document.getElementById("username").value;
      const password = document.getElementById("password").value;

      const response = await fetch("/authAdmin", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ username, password }),
      });

      const result = await response.json();
      const messageElement = document.getElementById("responseMessage");

      if (messageElement) {
        if (response.ok) {
          messageElement.textContent = result.message;
          messageElement.style.color = "green";

          // Redirect to the subpage after successful validation
          window.location.href =
            "http://localhost:8080/src/html/admin-panel.html"; // Change to your desired subpage URL
        } else {
          messageElement.textContent = result.error || "Login failed";
          messageElement.style.color = "red";
        }
      } else {
        console.error("Element with ID 'responseMessage' not found.");
      }
    });
});
