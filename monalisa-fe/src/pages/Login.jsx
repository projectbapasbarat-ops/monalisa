import { useState } from "react";
import { useNavigate } from "react-router-dom";
import api from "../api/axios";

export default function Login() {
  const [nip, setNip] = useState("");
  const [error, setError] = useState("");
  const navigate = useNavigate();

  const submit = async (e) => {
    e.preventDefault(); // 🔴 WAJIB

    setError("");

    try {
      const res = await api.post("/auth/login", { nip });

      // 🔐 SIMPAN TOKEN
      localStorage.setItem("access_token", res.data.access_token);
      localStorage.setItem("refresh_token", res.data.refresh_token);

      // ✅ NAVIGATE SETELAH LOGIN BERHASIL
      navigate("/admin/users");
    } catch (err) {
      setError("Login gagal");
      console.error(err);
    }
  };

  return (
    <div style={{ padding: 40 }}>
      <h2>Login</h2>

      <form onSubmit={submit}>
        <input
          placeholder="NIP"
          value={nip}
          onChange={(e) => setNip(e.target.value)}
        />

        <button type="submit">Login</button>
      </form>

      {error && <p style={{ color: "red" }}>{error}</p>}
    </div>
  );
}
