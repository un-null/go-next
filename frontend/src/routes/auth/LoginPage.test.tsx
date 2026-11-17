// @vitest-environment jsdom
import { MantineProvider } from "@mantine/core";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { MemoryRouter, Route, Routes } from "react-router";
import { afterEach, beforeEach, describe, expect, it } from "vitest";

import { fetchMock, jsonResponse, lastCall, mockApi } from "../../api/testing";
import type { User } from "../../api/types";
import { clearSession, getSession } from "../../lib/session";
import { DEMO, LoginPage } from "./LoginPage";

mockApi();

const USER: User = {
  id: "0d6f3c1e-6a1b-4f2a-9f0e-2b7c5d8a1e33",
  name: "Ada",
  email: "ada@example.com",
  coins: 1000,
  created_at: "2026-11-17T00:00:00Z",
  updated_at: "2026-11-17T00:00:00Z",
};

function renderLogin(path = "/login") {
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false } },
  });

  render(
    <MantineProvider>
      <QueryClientProvider client={queryClient}>
        <MemoryRouter initialEntries={[path]}>
          <Routes>
            <Route path="/login" element={<LoginPage />} />
            <Route path="/signup" element={<p>the signup page</p>} />
            <Route path="/cart" element={<p>the cart</p>} />
            <Route path="/products" element={<p>the catalogue</p>} />
          </Routes>
        </MemoryRouter>
      </QueryClientProvider>
    </MantineProvider>,
  );
}

beforeEach(() => {
  clearSession();
});

afterEach(() => {
  clearSession();
});

describe("LoginPage", () => {
  it("signs in and returns to the page that asked", async () => {
    fetchMock.mockResolvedValue(
      jsonResponse({ token: "jwt-value", user: USER }),
    );

    renderLogin("/login?next=%2Fcart");

    const user = userEvent.setup();
    await user.type(screen.getByLabelText(/^Email/), "ada@example.com");
    await user.type(screen.getByLabelText(/^Password/), "correct horse");
    await user.click(screen.getByRole("button", { name: "Sign in" }));

    expect(await screen.findByText("the cart")).not.toBeNull();
    expect(lastCall().body).toEqual({
      email: "ada@example.com",
      password: "correct horse",
    });
    expect(getSession()?.token).toBe("jwt-value");
  });

  it("shows the server's reason and keeps the form", async () => {
    fetchMock.mockResolvedValue(
      jsonResponse({ message: "Invalid credentials" }, 401),
    );

    renderLogin();

    const user = userEvent.setup();
    await user.type(screen.getByLabelText(/^Email/), "ada@example.com");
    await user.type(screen.getByLabelText(/^Password/), "wrong");
    await user.click(screen.getByRole("button", { name: "Sign in" }));

    expect(await screen.findByText("Invalid credentials")).not.toBeNull();
    expect(getSession()).toBeNull();
    expect(screen.getByRole("button", { name: "Sign in" })).not.toBeNull();
  });

  it("signs in as the demo with one click", async () => {
    fetchMock.mockResolvedValue(
      jsonResponse({ token: "jwt-value", user: USER }),
    );

    renderLogin();

    const user = userEvent.setup();
    await user.click(
      screen.getByRole("button", { name: "Use the demo account" }),
    );

    expect(await screen.findByText("the catalogue")).not.toBeNull();
    expect(lastCall().body).toEqual(DEMO);
    expect(getSession()?.user.email).toBe(USER.email);
  });

  it("does not call the API until the form is valid", async () => {
    renderLogin();

    const user = userEvent.setup();
    await user.click(screen.getByRole("button", { name: "Sign in" }));

    expect(await screen.findByText("Enter an email address")).not.toBeNull();
    expect(fetchMock).not.toHaveBeenCalled();
  });
});
