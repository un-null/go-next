// @vitest-environment jsdom
import { act, render, screen } from "@testing-library/react";
import { MemoryRouter, Route, Routes, useLocation } from "react-router";
import { afterEach, beforeEach, describe, expect, it } from "vitest";

import type { User } from "../api/types";
import { clearSession, setSession } from "../lib/session";
import { RequireAuth } from "./RequireAuth";

const user: User = {
  id: "0d6f3c1e-6a1b-4f2a-9f0e-2b7c5d8a1e33",
  name: "Ada",
  email: "ada@example.com",
  coins: 1000,
  created_at: "2026-11-16T00:00:00Z",
  updated_at: "2026-11-16T00:00:00Z",
};

function LoginPage() {
  const location = useLocation();
  return <p>{`${location.pathname}${location.search}`}</p>;
}

function renderAt(path: string) {
  return render(
    <MemoryRouter initialEntries={[path]}>
      <Routes>
        <Route element={<RequireAuth />}>
          <Route path="/cart" element={<p>the cart</p>} />
        </Route>
        <Route path="/login" element={<LoginPage />} />
      </Routes>
    </MemoryRouter>,
  );
}

beforeEach(() => {
  act(() => clearSession());
});

afterEach(() => {
  act(() => clearSession());
});

describe("RequireAuth", () => {
  it("shows the page to a signed-in visitor", () => {
    setSession({ token: "test-token", user });

    renderAt("/cart");

    expect(screen.getByText("the cart")).toBeDefined();
  });

  it("sends a guest to sign in, keeping the way back", () => {
    renderAt("/cart");

    expect(screen.getByText("/login?next=%2Fcart")).toBeDefined();
  });

  it("keeps the query string of the page it protected", () => {
    renderAt("/cart?coupon=x");

    expect(screen.getByText("/login?next=%2Fcart%3Fcoupon%3Dx")).toBeDefined();
  });

  it("leaves the page as soon as the session ends", () => {
    setSession({ token: "test-token", user });
    renderAt("/cart");

    act(() => clearSession());

    expect(screen.getByText("/login?next=%2Fcart")).toBeDefined();
  });
});
