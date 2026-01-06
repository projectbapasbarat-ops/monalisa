import { useEffect, useState } from "react";
import api from "../api/axios";
import AdminMenu from "../components/AdminMenu";

export default function AuditLogs() {
  const [logs, setLogs] = useState([]);
  const [page, setPage] = useState(1);
  const [total, setTotal] = useState(0);

  const loadData = async () => {
    const res = await api.get("/admin/audit-logs", {
      params: { page, limit: 10 },
    });

    setLogs(res.data.data);
    setTotal(res.data.meta?.total || 0);
  };

  useEffect(() => {
    loadData();
  }, [page]);

  return (
    <div style={{ padding: 40 }}>
      <AdminMenu />
      <h2>Audit Logs</h2>

      <table border="1" width="100%" cellPadding="8">
        <thead>
          <tr>
            <th>Waktu</th>
            <th>Aktor</th>
            <th>Aksi</th>
            <th>Target</th>
          </tr>
        </thead>

        <tbody>
          {logs.map((l) => (
            <tr key={l.id}>
              <td>{new Date(l.created_at).toLocaleString()}</td>
              <td>{l.actor_name}</td>
              <td>{l.action}</td>
              <td>{l.target}</td>
            </tr>
          ))}
        </tbody>
      </table>

      <div style={{ marginTop: 20 }}>
        <button disabled={page === 1} onClick={() => setPage(page - 1)}>Prev</button>
        <span style={{ margin: "0 10px" }}>Page {page}</span>
        <button disabled={page * 10 >= total} onClick={() => setPage(page + 1)}>Next</button>
      </div>
    </div>
  );
}
