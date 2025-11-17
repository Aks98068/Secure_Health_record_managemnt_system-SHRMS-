document.addEventListener("DOMContentLoaded", () => {

    // ------------------------
    // Registration
    // ------------------------
    const registrationForm = document.getElementById("registration-form");
    if (registrationForm) {
        const patientBtn = document.getElementById("patient-btn");
        const doctorBtn = document.getElementById("doctor-btn");
        const currentType = document.getElementById("current-type");
        const submitBtn = document.getElementById("submit-btn");

        // Toggle account type
        patientBtn.addEventListener("click", () => {
            patientBtn.classList.add("active");
            doctorBtn.classList.remove("active");
            currentType.innerText = "Patient";
            submitBtn.innerText = "Create Patient Account";
            document.getElementById("admin-notice").style.display = "none";
        });

        doctorBtn.addEventListener("click", () => {
            doctorBtn.classList.add("active");
            patientBtn.classList.remove("active");
            currentType.innerText = "Doctor";
            submitBtn.innerText = "Create Doctor Account";
            document.getElementById("admin-notice").style.display = "block";
        });

        // Submit registration form
        registrationForm.addEventListener("submit", async (e) => {
            e.preventDefault();

            let role = currentType.innerText.toLowerCase();
            let apiURL = (role === "patient") ? "/api/auth/register/patient" : "/api/auth/register/doctor";

            let data = {
                fullname: document.getElementById("fullName").value.trim(),
                email: document.getElementById("email").value.trim(),
                phonenumber: document.getElementById("phone").value.trim(),
                password: document.getElementById("password").value,
                confirmpassword: document.getElementById("confirmPassword").value
            };

            try {
                let response = await fetch(apiURL, {
                    method: "POST",
                    headers: { "Content-Type": "application/json" },
                    body: JSON.stringify(data)
                });

                let result = await response.json();

                if (response.ok) {
                    alert(result.message);
                    window.location.href = "/login";
                } else {
                    alert(result.error || "Registration failed");
                }
            } catch (err) {
                console.error(err);
                alert("Something went wrong.");
            }
        });
    }

    // ------------------------
    // Login
    // ------------------------
    const loginForm = document.getElementById("login-form");
    if (loginForm) {
        loginForm.addEventListener("submit", async (e) => {
            e.preventDefault();

            const role = document.getElementById("role-select").value.toLowerCase();

            let data = {
                email: document.getElementById("email").value.trim(),
                password: document.getElementById("password").value,
                role: role
            };

            try {
                let response = await fetch("/api/auth/login", {
                    method: "POST",
                    headers: { "Content-Type": "application/json" },
                    body: JSON.stringify(data)
                });

                let result = await response.json();

                if (response.ok) {
                    alert(result.message);

                    // Redirect based on role
                    if (role === "patient") window.location.href = "/patient/dashboard";
                    else if (role === "doctor") window.location.href = "/doctor/dashboard";
                    else window.location.href = "/admin/dashboard";

                } else {
                    alert(result.error || "Login failed");
                }
            } catch (err) {
                console.error(err);
                alert("Something went wrong.");
            }
        });
    }

});
