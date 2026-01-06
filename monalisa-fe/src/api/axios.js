import axios from "axios";

// ============================
// DEBUG: pastikan env terbaca
// ============================
const baseURL = import.meta.env.VITE_API_URL;

if (!baseURL) {
  console.error(
    "VITE_API_URL is undefined. Pastikan file .env ada di root frontend dan dev server direstart."
  );
} else {
  console.log("API BASE URL:", baseURL);
}

// ============================
// AXIOS INSTANCE
// ============================
const api = axios.create({
  baseURL: baseURL, // contoh: http://localhost:8080/api/v1
  headers: {
    "Content-Type": "application/json",
  },
});

// ============================
// REQUEST INTERCEPTOR
// ============================
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem("token");

    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }

    return config;
  },
  (error) => Promise.reject(error)
);

// ============================
// RESPONSE INTERCEPTOR
// ============================
api.interceptors.response.use(
  (response) => response,
  (error) => {
    // DEBUG ERROR RESPONSE
    if (error.response) {
      console.error(
        "API ERROR:",
        error.response.status,
        error.response.data
      );
    } else {
      console.error("API ERROR:", error.message);
    }

    return Promise.reject(error);
  }
);

export default api;
