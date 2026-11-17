import { Navigate, Outlet, useLocation } from "react-router";

import { useSession } from "../lib/session";

export function RequireAuth() {
  const session = useSession();
  const location = useLocation();

  if (!session) {
    const next = `${location.pathname}${location.search}`;
    return <Navigate to={`/login?next=${encodeURIComponent(next)}`} replace />;
  }

  return <Outlet />;
}
