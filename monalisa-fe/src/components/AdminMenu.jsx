import { hasPermission } from "../utils/auth";

export default function AdminMenu() {
  return (
    <div style={{ marginBottom: 20 }}>
      <h3>Admin Menu</h3>

      <ul>
        {hasPermission("user.read") && (
          <li>
            <a href="/admin/users">Manajemen User</a>
          </li>
        )}

        {hasPermission("audit.read") && (
          <li>
            <a href="/admin/audit-logs">Audit Log</a>
          </li>
        )}
      </ul>
    </div>
  );
}
