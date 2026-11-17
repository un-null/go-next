// @vitest-environment jsdom
import { MantineProvider } from "@mantine/core";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { act, render, screen, waitFor } from "@testing-library/react";
import { MemoryRouter, Route, Routes } from "react-router";
import { afterEach, beforeEach, describe, expect, it } from "vitest";

import { fetchMock, jsonResponse, lastCall, mockApi } from "../../api/testing";
import type { Order, User } from "../../api/types";
import { clearSession, setSession } from "../../lib/session";
import { OrdersPage } from "./OrdersPage";

mockApi();

const USER: User = {
  id: "0d6f3c1e-6a1b-4f2a-9f0e-2b7c5d8a1e33",
  name: "Ada",
  email: "ada@example.com",
  coins: 1000,
  created_at: "2026-11-01T00:00:00Z",
  updated_at: "2026-11-01T00:00:00Z",
};

function order(id: number, number: string, status: Order["status"]): Order {
  return {
    id,
    user_id: USER.id,
    order_number: number,
    total_amount: 949.97,
    total_coins_used: 950,
    status,
    created_at: "2026-11-17T12:00:00Z",
    updated_at: "2026-11-17T12:00:00Z",
  };
}

const ORDERS = [
  order(7, "ORD-20261117-0007", "completed"),
  order(8, "ORD-20261117-0008", "pending"),
];

function renderOrders() {
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false } },
  });

  return render(
    <MantineProvider>
      <QueryClientProvider client={queryClient}>
        <MemoryRouter initialEntries={["/orders"]}>
          <Routes>
            <Route path="/orders" element={<OrdersPage />} />
            <Route path="/orders/:id" element={<p>the order</p>} />
            <Route path="/products" element={<p>the catalogue</p>} />
          </Routes>
        </MemoryRouter>
      </QueryClientProvider>
    </MantineProvider>,
  );
}

beforeEach(() => {
  setSession({ token: "jwt-value", user: USER });
});

afterEach(() => {
  act(() => clearSession());
});

describe("OrdersPage", () => {
  it("asks the API for this user's orders and lists them", async () => {
    fetchMock.mockResolvedValue(jsonResponse(ORDERS));

    renderOrders();

    expect(await screen.findByText("ORD-20261117-0007")).toBeTruthy();
    expect(screen.getByText("ORD-20261117-0008")).toBeTruthy();
    // Two orders share the same date, price and coins.
    expect(screen.getAllByText("Nov 17, 2026")).toHaveLength(2);
    expect(screen.getAllByText("$949.97")).toHaveLength(2);
    expect(screen.getAllByText("950 coins")).toHaveLength(2);

    const call = lastCall();
    expect(call.url).toBe(`https://api.test/api/users/${USER.id}/orders`);
    expect(call.headers.Authorization).toBe("Bearer jwt-value");

    const links = screen.getAllByRole("link", { name: "View order" });
    expect(links.map((link) => link.getAttribute("href"))).toEqual([
      "/orders/7",
      "/orders/8",
    ]);
  });

  it("says so when there is nothing to show", async () => {
    fetchMock.mockResolvedValue(jsonResponse([]));

    renderOrders();

    expect(await screen.findByText(/No orders yet/)).toBeTruthy();
    expect(
      screen
        .getByRole("link", { name: "Browse products" })
        .getAttribute("href"),
    ).toBe("/products");
  });

  it("shows the server's message when the list cannot be read", async () => {
    fetchMock.mockResolvedValue(
      jsonResponse({ message: "orders are unavailable" }, 500),
    );

    renderOrders();

    expect(await screen.findByText("orders are unavailable")).toBeTruthy();
    await waitFor(() =>
      expect(screen.getByRole("button", { name: "Try again" })).toBeTruthy(),
    );
  });
});
