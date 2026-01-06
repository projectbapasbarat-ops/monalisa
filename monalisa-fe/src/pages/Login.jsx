import { useState } from "react";
import api from "../api/axios";
import { setToken } from "../auth/authStorage";

export default function Login() {
  const [nip, setNip] = useState("");
  const [error, setError] = useState("");

  const submit = async (e) => {
    e.preventDefault();
    try {
      const res = await api.post("/auth/login", { nip });
      setToken(res.data.token);
      window.location.href = "/admin/users";
    } catch {
      setError("Login gagal");
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
        {error && <p style={{ color: "red" }}>{error}</p>}
      </form>
    </div>
  );
}
