export function getToken() {
  return localStorage.getItem("access_token");
}

export function getPayload() {
  const token = getToken();
  if (!token) return null;

  try {
    const payload = token.split(".")[1];
    return JSON.parse(atob(payload));
  } catch {
    return null;
  }
}

export function getPermissions() {
  const payload = getPayload();
  return payload?.permissions || [];
}

export function hasPermission(permission) {
  return getPermissions().includes(permission);
}
