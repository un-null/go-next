// @vitest-environment jsdom
import { MantineProvider } from "@mantine/core";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { act, render, screen, waitFor } from "@testing-library/react";
import { MemoryRouter, Route, Routes } from "react-router";
import { afterEach, beforeEach, describe, expect, it } from "vitest";

import { fetchMock, jsonResponse, lastCall, mockApi } from "../../api/testing";
import type { CoinTransaction, User } from "../../api/types";
import { clearSession, setSession } from "../../lib/session";
import { AccountPage } from "./AccountPage";

mockApi();

const SIGNED_IN: User = {
  id: "0d6f3c1e-6a1b-4f2a-9f0e-2b7c5d8a1e33",
  name: "Ada",
  email: "ada@example.com",
  coins: 1000,
  created_at: "2026-11-01T12:00:00Z",
  updated_at: "2026-11-01T12:00:00Z",
};

const SERVER: User = { ...SIGNED_IN, coins: 250 };

const ACTIVITY: CoinTransaction[] = [
  {
    id: 2,
    user_id: SIGNED_IN.id,
    transaction_type: "purchase",
    amount: -950,
    balance_after: 50,
    order_id: 7,
    description: "Order ORD-20261117-0007",
    created_at: "2026-11-17T12:00:00Z",
  },
  {
    id: 1,
    user_id: SIGNED_IN.id,
    transaction_type: "charge",
    amount: 1000,
    balance_after: 1000,
    description: "Demo top-up",
    created_at: "2026-11-16T12:00:00Z",
  },
];

function serve(
  transactions: CoinTransaction[] = ACTIVITY,
  activityStatus = 200,
) {
  fetchMock.mockImplementation(async (input) => {
    const url = String(input);

    if (url.includes("/api/coins/transactions")) {
      return jsonResponse(
        { transactions, page: 1, limit: 100 },
        activityStatus,
      );
    }

    return jsonResponse(SERVER);
  });
}

function renderAccount() {
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false } },
  });

  return render(
    <MantineProvider>
      <QueryClientProvider client={queryClient}>
        <MemoryRouter initialEntries={["/account"]}>
          <Routes>
            <Route path="/account" element={<AccountPage />} />
          </Routes>
        </MemoryRouter>
      </QueryClientProvider>
    </MantineProvider>,
  );
}

beforeEach(() => {
  setSession({ token: "jwt-value", user: SIGNED_IN });
});

afterEach(() => {
  act(() => clearSession());
});

describe("AccountPage", () => {
  it("shows the balance the server reports, not the one from sign-in", async () => {
    serve();
    renderAccount();

    expect(await screen.findByText("250 coins")).toBeTruthy();
    expect(screen.getByText("ada@example.com")).toBeTruthy();
    expect(screen.getByText("Nov 1, 2026")).toBeTruthy();
    expect(screen.getByText("Demo top-up")).toBeTruthy();
    expect(screen.getByText("+1,000")).toBeTruthy();
    expect(screen.getByText("Order ORD-20261117-0007")).toBeTruthy();
    expect(screen.getByText("−950")).toBeTruthy();
    expect(lastCall().headers.Authorization).toBe("Bearer jwt-value");
  });

  it("asks for one page of activity", async () => {
    serve();
    renderAccount();

    await screen.findByText("Demo top-up");

    const urls = fetchMock.mock.calls.map(([url]) => String(url));
    expect(
      urls.some((url) =>
        url.endsWith("/api/coins/transactions?page=1&limit=100"),
      ),
    ).toBe(true);
  });

  it("says so when the ledger is empty", async () => {
    serve([]);
    renderAccount();

    expect(await screen.findByText("No coin activity yet.")).toBeTruthy();
  });

  it("keeps the balance when the ledger fails", async () => {
    serve(ACTIVITY, 500);
    renderAccount();

    expect(await screen.findByText("250 coins")).toBeTruthy();
    expect(
      await screen.findByText("Could not load your activity"),
    ).toBeTruthy();
    await waitFor(() => expect(screen.getByText("Try again")).toBeTruthy());
  });
});
