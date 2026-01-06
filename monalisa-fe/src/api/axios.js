import axios from "axios";

const baseURL = import.meta.env.VITE_API_URL;

if (!baseURL) {
  console.error("VITE_API_URL is undefined");
}

const api = axios.create({
  baseURL,
  headers: {
    "Content-Type": "application/json",
  },
});

// REQUEST INTERCEPTOR
api.interceptors.request.use((config) => {
  const token = localStorage.getItem("access_token");
  console.log("AXIOS TOKEN:", token);
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});


// RESPONSE INTERCEPTOR
api.interceptors.response.use(
  (response) => response,
  (error) => {
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
