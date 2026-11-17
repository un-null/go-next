// @vitest-environment jsdom
import { MantineProvider } from "@mantine/core";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { act, render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { MemoryRouter, Route, Routes } from "react-router";
import { afterEach, describe, expect, it } from "vitest";

import type { User } from "../api/types";
import { clearSession, setSession } from "../lib/session";
import { AppLayout } from "./AppLayout";

const USER: User = {
  id: "0d6f3c1e-6a1b-4f2a-9f0e-2b7c5d8a1e33",
  name: "Ada",
  email: "ada@example.com",
  coins: 1000,
  created_at: "2026-11-01T00:00:00Z",
  updated_at: "2026-11-01T00:00:00Z",
};

function renderLayout() {
  return render(
    <MantineProvider>
      <QueryClientProvider client={new QueryClient()}>
        <MemoryRouter initialEntries={["/products"]}>
          <Routes>
            <Route element={<AppLayout />}>
              <Route path="/products" element={<p>the catalogue</p>} />
            </Route>
          </Routes>
        </MemoryRouter>
      </QueryClientProvider>
    </MantineProvider>,
  );
}

afterEach(() => {
  act(() => clearSession());
});

describe("AppLayout", () => {
  it("keeps the cart out of the header for a guest", () => {
    renderLayout();

    expect(screen.queryByRole("link", { name: "Cart" })).toBeNull();
    expect(
      screen.getByRole("link", { name: "Sign in" }).getAttribute("href"),
    ).toBe("/login");
  });

  it("links to the cart once there is a session", () => {
    setSession({ token: "jwt-value", user: USER });
    renderLayout();

    expect(
      screen.getByRole("link", { name: "Cart" }).getAttribute("href"),
    ).toBe("/cart");
    expect(screen.queryByRole("link", { name: "Sign in" })).toBeNull();
    expect(screen.getByRole("button", { name: "Ada" })).toBeTruthy();
  });

  it("offers the order history from the account menu", async () => {
    const user = userEvent.setup();
    setSession({ token: "jwt-value", user: USER });
    renderLayout();

    await user.click(screen.getByRole("button", { name: "Ada" }));

    const item = await screen.findByRole("menuitem", { name: "Orders" });
    expect(item.getAttribute("href")).toBe("/orders");
  });
});
