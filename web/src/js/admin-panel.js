document.addEventListener("DOMContentLoaded", function () {
  const addUserForm = document.getElementById("addUserForm");
  const editUserForm = document.getElementById("editUserForm");
  const listUsersButton = document.getElementById("listUsersButton");

  if (addUserForm) {
    addUserForm.addEventListener("submit", async (e) => {
      e.preventDefault();
      const username = document.getElementById("username").value;
      const password = document.getElementById("password").value;
      const response = await fetch("/addUser", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ username, password }),
      });
      if (response.ok) {
        alert("User added successfully");
      } else {
        alert("Error adding user");
      }
    });
  }

  if (editUserForm) {
    editUserForm.addEventListener("submit", async (e) => {
      e.preventDefault();
      const userID = document.getElementById("userId").value;
      const newUsername = document.getElementById("newUsername").value;
      const newPassword = document.getElementById("newPassword").value;
      const response = await fetch("/editUser", {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ userID, newUsername, newPassword }),
      });
      if (response.ok) {
        alert("User edited successfully");
      } else {
        alert("Error editing user");
      }
    });
  }

  if (listUsersButton) {
    listUsersButton.addEventListener("click", async () => {
      const response = await fetch("/listUsers");
      if (response.ok) {
        const users = await response.json();
        const userTableBody = document.querySelector("#userTable tbody");
        userTableBody.innerHTML = "";
        users.forEach((user) => {
          const tr = document.createElement("tr");
          tr.className = "border-t border-gray-300 hover:bg-gray-200"; // Apply Tailwind CSS classes to the row

          const idCell = document.createElement("td");
          idCell.className = "border-t border-gray-300 px-4 py-3"; // Apply Tailwind CSS classes to the cell
          idCell.textContent = user.id || user.user_id;
          tr.appendChild(idCell);

          const usernameCell = document.createElement("td");
          usernameCell.className = "border-t border-gray-300 px-4 py-3"; // Apply Tailwind CSS classes to the cell
          usernameCell.textContent = user.username || user.Username;
          tr.appendChild(usernameCell);

          const passwordCell = document.createElement("td");
          passwordCell.className = "border-t border-gray-300 px-4 py-3"; // Apply Tailwind CSS classes to the cell
          passwordCell.textContent = user.password || user.Password;
          tr.appendChild(passwordCell);

          userTableBody.appendChild(tr);
        });
      } else {
        alert("Error listing users");
      }
    });
  }
});
