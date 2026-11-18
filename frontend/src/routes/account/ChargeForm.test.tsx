// @vitest-environment jsdom
import { MantineProvider } from "@mantine/core";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { act } from "react";
import { MemoryRouter, Route, Routes } from "react-router";
import { afterEach, beforeEach, describe, expect, it } from "vitest";

import { fetchMock, jsonResponse, mockApi } from "../../api/testing";
import type { User } from "../../api/types";
import { clearSession, getSession, setSession } from "../../lib/session";
import { AccountPage } from "./AccountPage";

const USER: User = {
  id: "0d6f3c1e-6a1b-4f2a-9f0e-2b7c5d8a1e33",
  name: "Ada",
  email: "ada@example.com",
  coins: 1000,
  created_at: "2026-11-01T12:00:00Z",
  updated_at: "2026-11-01T12:00:00Z",
};

function chargeResponse(coins: number, amount: number) {
  return {
    message: "Coins charged successfully",
    user: { ...USER, coins },
    transaction: {
      id: 9,
      user_id: USER.id,
      transaction_type: "charge",
      amount,
      balance_after: coins,
      description: "Demo top-up",
      created_at: "2026-11-17T12:00:00Z",
    },
  };
}

function serve({ fail }: { fail?: () => Response } = {}) {
  let coins = USER.coins;

  fetchMock.mockImplementation(async (input) => {
    const url = String(input);

    if (url.endsWith("/api/coins/charge")) {
      if (fail) {
        return fail();
      }
      coins += 500;
      return jsonResponse(chargeResponse(coins, 500));
    }

    if (url.includes("/api/coins/transactions")) {
      return jsonResponse({ transactions: [], page: 1, limit: 100 });
    }

    return jsonResponse({ ...USER, coins });
  });
}

function postCount() {
  return fetchMock.mock.calls.filter(([input]) =>
    String(input).endsWith("/api/coins/charge"),
  ).length;
}

function renderAccount() {
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false } },
  });

  render(
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

describe("ChargeForm", () => {
  mockApi();

  beforeEach(() => {
    setSession({ token: "jwt-value", user: USER });
    serve();
  });

  afterEach(() => {
    act(() => clearSession());
  });

  it("charges the amount in the form and shows the new balance", async () => {
    const user = userEvent.setup();
    renderAccount();

    await screen.findByText("1,000 coins");
    await user.click(screen.getByRole("button", { name: "Add coins" }));

    await screen.findByText("Charged 500 coins.");

    const charge = fetchMock.mock.calls.find(([input]) =>
      String(input).endsWith("/api/coins/charge"),
    );
    expect(JSON.parse(String(charge?.[1]?.body))).toEqual({
      amount: 500,
      description: "Demo top-up",
    });
    expect((charge?.[1]?.headers as Record<string, string>).Authorization).toBe(
      "Bearer jwt-value",
    );

    await waitFor(() => expect(screen.getByText("1,500 coins")).toBeDefined());
  });

  it("refuses an empty description without calling the API", async () => {
    const user = userEvent.setup();
    renderAccount();

    const description = screen.getByLabelText(/^Description/);
    await user.clear(description);
    await user.click(screen.getByRole("button", { name: "Add coins" }));

    expect(await screen.findByText("Say what the coins are for")).toBeDefined();
    expect(postCount()).toBe(0);
  });

  it("shows the server's message when the charge fails", async () => {
    serve({
      fail: () => jsonResponse({ message: "amount too large" }, 400),
    });

    const user = userEvent.setup();
    renderAccount();

    await user.click(screen.getByRole("button", { name: "Add coins" }));

    expect(await screen.findByText("amount too large")).toBeDefined();
    expect(getSession()).not.toBeNull();
  });
});
