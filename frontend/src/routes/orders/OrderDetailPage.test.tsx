// @vitest-environment jsdom
import { MantineProvider } from "@mantine/core";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { act, render, screen, waitFor } from "@testing-library/react";
import { MemoryRouter, Route, Routes } from "react-router";
import { afterEach, beforeEach, describe, expect, it } from "vitest";

import { fetchMock, jsonResponse, lastCall, mockApi } from "../../api/testing";
import type { Order, User } from "../../api/types";
import { clearSession, setSession } from "../../lib/session";
import { OrderDetailPage } from "./OrderDetailPage";

mockApi();

const USER: User = {
  id: "0d6f3c1e-6a1b-4f2a-9f0e-2b7c5d8a1e33",
  name: "Ada",
  email: "ada@example.com",
  coins: 1000,
  created_at: "2026-11-01T00:00:00Z",
  updated_at: "2026-11-01T00:00:00Z",
};

const ORDER: Order = {
  id: 7,
  user_id: USER.id,
  order_number: "ORD-20261117-0007",
  total_amount: 949.97,
  total_coins_used: 950,
  status: "completed",
  created_at: "2026-11-17T12:00:00Z",
  updated_at: "2026-11-17T12:00:00Z",
  items: [
    {
      id: 1,
      order_id: 7,
      product_id: 1,
      product_name: "4K Monitor UltraSharp",
      product_price: 399.99,
      quantity: 2,
      subtotal: 799.98,
    },
    {
      id: 2,
      order_id: 7,
      product_id: 2,
      product_name: "Wireless Earbuds Pro",
      product_price: 149.99,
      quantity: 1,
      subtotal: 149.99,
    },
  ],
};

function renderOrder(path = "/orders/7") {
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false } },
  });

  return render(
    <MantineProvider>
      <QueryClientProvider client={queryClient}>
        <MemoryRouter initialEntries={[path]}>
          <Routes>
            <Route path="/orders/:id" element={<OrderDetailPage />} />
            <Route path="/orders" element={<p>the list</p>} />
            <Route path="/products/:id" element={<p>the product</p>} />
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

describe("OrderDetailPage", () => {
  it("shows every line the server sent, with its own subtotal", async () => {
    fetchMock.mockImplementation(async () => jsonResponse(ORDER));

    renderOrder();

    expect(await screen.findByText("ORD-20261117-0007")).toBeTruthy();
    expect(screen.getByText(/Nov 17, 2026/)).toBeTruthy();

    expect(screen.getByText("4K Monitor UltraSharp")).toBeTruthy();
    expect(screen.getByText("2 × $399.99")).toBeTruthy();
    expect(screen.getByText("Wireless Earbuds Pro")).toBeTruthy();
    expect(screen.getByText("1 × $149.99")).toBeTruthy();
    expect(screen.getByText("$799.98")).toBeTruthy();

    expect(screen.getByText("$949.97")).toBeTruthy();
    expect(screen.getByText("950 coins")).toBeTruthy();

    expect(screen.getByText("4K Monitor UltraSharp").getAttribute("href")).toBe(
      "/products/1",
    );

    expect(lastCall().url).toBe("https://api.test/api/orders/7");
    expect(lastCall().headers.Authorization).toBe("Bearer jwt-value");
  });

  it("answers a 404 with the same page as a bad id, and no retry", async () => {
    fetchMock.mockImplementation(async () =>
      jsonResponse({ message: "order not found" }, 404),
    );

    renderOrder();

    expect(await screen.findByText("Order not found")).toBeTruthy();
    expect(screen.queryByRole("button", { name: "Try again" })).toBeNull();
  });

  it("does not ask the server about an id that cannot exist", () => {
    renderOrder("/orders/abc");

    expect(screen.getByText("Order not found")).toBeTruthy();
    expect(fetchMock).not.toHaveBeenCalled();
  });

  it("keeps the server's message when the request fails", async () => {
    fetchMock.mockImplementation(async () =>
      jsonResponse({ message: "orders are unavailable" }, 500),
    );

    renderOrder();

    await waitFor(() => {
      expect(screen.getByText("orders are unavailable")).toBeTruthy();
    });
    expect(screen.getByRole("button", { name: "Try again" })).toBeTruthy();
  });
});
