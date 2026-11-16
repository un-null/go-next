import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

import type { User } from "../api/types";
import { clearSession, getSession, setSession, subscribe } from "./session";

const KEY = "go-next.session";

const user: User = {
  id: "0d6f3c1e-6a1b-4f2a-9f0e-2b7c5d8a1e33",
  name: "Ada",
  email: "ada@example.com",
  coins: 1000,
  created_at: "2026-11-16T09:00:00Z",
  updated_at: "2026-11-16T09:00:00Z",
};

const session = { token: "token-1", user };

let values: Map<string, string>;

function fakeStorage(): Storage {
  return {
    getItem: (key) => values.get(key) ?? null,
    setItem: (key, value) => {
      values.set(key, value);
    },
    removeItem: (key) => {
      values.delete(key);
    },
    clear: () => {
      values.clear();
    },
    key: (index) => [...values.keys()][index] ?? null,
    get length() {
      return values.size;
    },
  };
}

// The module caches its snapshot, so going back to storage needs a fresh module.
// This is what a page load does.
async function reload() {
  vi.resetModules();
  return import("./session");
}

beforeEach(() => {
  values = new Map();
  vi.stubGlobal("sessionStorage", fakeStorage());
  clearSession();
});

afterEach(() => {
  vi.unstubAllGlobals();
});

describe("session", () => {
  it("starts empty", () => {
    expect(getSession()).toBeNull();
    expect(values.size).toBe(0);
  });

  it("returns the same object until the session changes", () => {
    setSession(session);
    const first = getSession();

    expect(getSession()).toBe(first);

    setSession({ ...session, token: "token-2" });
    expect(getSession()).not.toBe(first);

    clearSession();
    expect(getSession()).toBeNull();
  });

  it("reads what the previous page load left behind", async () => {
    setSession(session);

    const fresh = await reload();

    expect(fresh.getSession()).toEqual(session);
  });

  it("ignores a broken value", async () => {
    values.set(KEY, "{ not json");

    const fresh = await reload();

    expect(fresh.getSession()).toBeNull();
  });

  it("ignores a value that is missing half of the session", async () => {
    values.set(KEY, JSON.stringify({ user }));
    const withoutToken = await reload();

    values.set(KEY, JSON.stringify({ token: "token-1" }));
    const withoutUser = await reload();

    expect(withoutToken.getSession()).toBeNull();
    expect(withoutUser.getSession()).toBeNull();
  });

  it("notifies subscribers until they unsubscribe", () => {
    const listener = vi.fn();
    const unsubscribe = subscribe(listener);

    setSession(session);
    expect(listener).toHaveBeenCalledTimes(1);

    unsubscribe();
    clearSession();
    expect(listener).toHaveBeenCalledTimes(1);
  });

  it("keeps the session in memory when storage rejects the write", () => {
    vi.stubGlobal("sessionStorage", {
      ...fakeStorage(),
      setItem: () => {
        throw new Error("storage is blocked");
      },
    });

    setSession(session);

    expect(getSession()).toEqual(session);
  });
});
