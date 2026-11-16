import { useSyncExternalStore } from "react";

import type { User } from "../api/types";

export interface Session {
  token: string;
  user: User;
}

const KEY = "go-next.session";

const listeners = new Set<() => void>();

let snapshot: Session | null | undefined;

export function getSession(): Session | null {
  if (snapshot === undefined) {
    snapshot = parse(storage()?.getItem(KEY) ?? null);
  }
  return snapshot;
}

export function setSession(session: Session): void {
  commit(session);
}

export function clearSession(): void {
  commit(null);
}

export function useSession(): Session | null {
  return useSyncExternalStore(subscribe, getSession, getSession);
}

function commit(next: Session | null): void {
  snapshot = next;

  try {
    if (next) {
      storage()?.setItem(KEY, JSON.stringify(next));
    } else {
      storage()?.removeItem(KEY);
    }
  } catch {
  }

  for (const listener of listeners) {
    listener();
  }
}

export function subscribe(listener: () => void): () => void {
  listeners.add(listener);
  return () => {
    listeners.delete(listener);
  };
}

function storage(): Storage | null {
  return typeof sessionStorage === "undefined" ? null : sessionStorage;
}

function parse(raw: string | null): Session | null {
  if (!raw) {
    return null;
  }

  try {
    const value = JSON.parse(raw) as { token?: unknown; user?: unknown };

    if (typeof value.token !== "string" || !value.token) {
      return null;
    }
    if (typeof value.user !== "object" || value.user === null) {
      return null;
    }

    return { token: value.token, user: value.user as User };
  } catch {
    return null;
  }
}
