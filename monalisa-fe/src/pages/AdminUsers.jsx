import { useEffect, useState } from "react";
import api from "../api/axios";

export default function AdminUsers() {
  const [users, setUsers] = useState([]);
  const [roles, setRoles] = useState([]);
  const [selectedRole, setSelectedRole] = useState({});

  const loadData = async () => {
    const usersRes = await api.get("/admin/users");
    const rolesRes = await api.get("/admin/roles");

    setUsers(usersRes.data.data);
    setRoles(rolesRes.data.data);
  };

  useEffect(() => {
    loadData().catch(() => alert("Forbidden"));
  }, []);

  const assignRole = async (userId) => {
    const role = selectedRole[userId];
    if (!role) return alert("Pilih role dulu");

    await api.post(`/admin/users/${userId}/roles`, {
      role_code: role,
    });

    loadData();
  };

  const removeRole = async (userId, role) => {
    await api.delete(`/admin/users/${userId}/roles/${role}`);
    loadData();
  };

  return (
    <div style={{ padding: 40 }}>
      <h2>Admin - Users</h2>

      <table border="1" cellPadding="8">
        <thead>
          <tr>
            <th>NIP</th>
            <th>Nama</th>
            <th>Jabatan</th>
            <th>Roles</th>
            <th>Manage</th>
          </tr>
        </thead>

        <tbody>
          {users.map((u) => (
            <tr key={u.id}>
              <td>{u.nip}</td>
              <td>{u.nama}</td>
              <td>{u.jabatan}</td>

              <td>
                {u.roles.map((r) => (
                  <span key={r} style={{ marginRight: 6 }}>
                    {r}
                    <button
                      style={{ marginLeft: 4 }}
                      onClick={() => removeRole(u.id, r)}
                    >
                      ❌
                    </button>
                  </span>
                ))}
              </td>

              <td>
                <select
                  onChange={(e) =>
                    setSelectedRole({
                      ...selectedRole,
                      [u.id]: e.target.value,
                    })
                  }
                >
                  <option value="">-- pilih role --</option>
                  {roles.map((r) => (
                    <option key={r} value={r}>
                      {r}
                    </option>
                  ))}
                </select>

                <button onClick={() => assignRole(u.id)}>➕</button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
