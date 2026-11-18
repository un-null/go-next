// @vitest-environment jsdom
import { MantineProvider } from "@mantine/core";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import {
  act,
  fireEvent,
  render,
  screen,
  waitFor,
} from "@testing-library/react";
import { MemoryRouter, Route, Routes } from "react-router";
import { afterEach, beforeEach, describe, expect, it } from "vitest";

import { fetchMock, jsonResponse, mockApi } from "../api/testing";
import type { User } from "../api/types";
import { clearSession, setSession } from "../lib/session";
import { AppLayout } from "./AppLayout";

mockApi();

const USER: User = {
  id: "0d6f3c1e-6a1b-4f2a-9f0e-2b7c5d8a1e33",
  name: "Ada",
  email: "ada@example.com",
  coins: 1000,
  created_at: "2026-11-01T00:00:00Z",
  updated_at: "2026-11-01T00:00:00Z",
};

function serve(user: User = USER, cartCount = 0) {
  fetchMock.mockImplementation((input) => {
    const url = new URL(String(input));
    if (url.pathname === "/api/cart") {
      return Promise.resolve(
        jsonResponse({
          items: [],
          total_items: cartCount,
          success: true,
        }),
      );
    }
    return Promise.resolve(jsonResponse(user));
  });
}

beforeEach(() => {
  serve();
});

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

async function openMenu() {
  const trigger = screen.getByRole("button", { name: "Ada" });
  fireEvent.click(trigger);
  await waitFor(() => {
    expect(trigger.getAttribute("aria-expanded")).toBe("true");
  });
}

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
    setSession({ token: "jwt-value", user: USER });
    renderLayout();

    await openMenu();

    const item = await screen.findByRole(
      "menuitem",
      { name: "Orders" },
      { timeout: 8000 },
    );
    expect(item.getAttribute("href")).toBe("/orders");
  });

  it("shows the balance the server reports, not the one from sign-in", async () => {
    serve({ ...USER, coins: 250 });
    setSession({ token: "jwt-value", user: USER });
    renderLayout();

    await openMenu();

    expect(await screen.findByText("250 coins")).toBeTruthy();
  });

  it("offers the account page from the menu", async () => {
    setSession({ token: "jwt-value", user: USER });
    renderLayout();

    await openMenu();

    const item = await screen.findByRole(
      "menuitem",
      { name: "Account" },
      { timeout: 8000 },
    );
    expect(item.getAttribute("href")).toBe("/account");
  });

  it("shows how many lines are in the cart", async () => {
    serve(USER, 3);
    setSession({ token: "jwt-value", user: USER });
    renderLayout();

    const link = await screen.findByRole("link", { name: "Cart, 3 items" });
    expect(link.getAttribute("href")).toBe("/cart");
    expect(screen.queryByRole("link", { name: "Cart" })).toBeNull();
  });
});
