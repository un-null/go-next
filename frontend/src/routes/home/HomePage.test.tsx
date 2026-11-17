// @vitest-environment jsdom
import { MantineProvider } from "@mantine/core";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { render, screen } from "@testing-library/react";
import { MemoryRouter, Route, Routes } from "react-router";
import { afterEach, beforeEach, describe, expect, it } from "vitest";

import { fetchMock, jsonResponse, mockApi } from "../../api/testing";
import type { User } from "../../api/types";
import { clearSession, setSession } from "../../lib/session";
import { HomePage } from "./HomePage";

mockApi();

const USER: User = {
  id: "0d6f3c1e-6a1b-4f2a-9f0e-2b7c5d8a1e33",
  name: "Ada",
  email: "ada@example.com",
  coins: 1000,
  created_at: "2026-11-17T00:00:00Z",
  updated_at: "2026-11-17T00:00:00Z",
};

function renderHome() {
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false } },
  });

  render(
    <MantineProvider>
      <QueryClientProvider client={queryClient}>
        <MemoryRouter initialEntries={["/"]}>
          <Routes>
            <Route path="/" element={<HomePage />} />
          </Routes>
        </MemoryRouter>
      </QueryClientProvider>
    </MantineProvider>,
  );
}

beforeEach(() => {
  clearSession();
  fetchMock.mockResolvedValue(jsonResponse([]));
});

afterEach(() => {
  clearSession();
});

describe("HomePage", () => {
  it("sends guests to signup", () => {
    renderHome();

    expect(
      screen
        .getByRole("link", { name: "Create an account" })
        .getAttribute("href"),
    ).toBe("/signup");
    expect(
      screen.getByRole("link", { name: "Sign in" }).getAttribute("href"),
    ).toBe("/login");
  });

  it("sends signed-in visitors to the catalogue", () => {
    setSession({ token: "jwt-value", user: USER });

    renderHome();

    expect(
      screen.getByRole("link", { name: "Start shopping" }).getAttribute("href"),
    ).toBe("/products");
    expect(
      screen.queryByRole("link", { name: "Create an account" }),
    ).toBeNull();
  });
});
