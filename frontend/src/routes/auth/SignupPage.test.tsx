// @vitest-environment jsdom
import { MantineProvider } from "@mantine/core";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { MemoryRouter, Route, Routes } from "react-router";
import { afterEach, beforeEach, describe, expect, it } from "vitest";

import { deferred, fetchMock, jsonResponse, mockApi } from "../../api/testing";
import type { User } from "../../api/types";
import { clearSession, getSession } from "../../lib/session";
import { SignupPage } from "./SignupPage";

mockApi();

const USER: User = {
  id: "0d6f3c1e-6a1b-4f2a-9f0e-2b7c5d8a1e33",
  name: "Ada",
  email: "ada@example.com",
  coins: 1000,
  created_at: "2026-11-17T00:00:00Z",
  updated_at: "2026-11-17T00:00:00Z",
};

function renderSignup(path = "/signup") {
  const queryClient = new QueryClient({
    defaultOptions: { queries: { retry: false } },
  });

  render(
    <MantineProvider>
      <QueryClientProvider client={queryClient}>
        <MemoryRouter initialEntries={[path]}>
          <Routes>
            <Route path="/signup" element={<SignupPage />} />
            <Route path="/login" element={<p>the sign in page</p>} />
            <Route path="/cart" element={<p>the cart</p>} />
            <Route path="/products" element={<p>the catalogue</p>} />
          </Routes>
        </MemoryRouter>
      </QueryClientProvider>
    </MantineProvider>,
  );
}

function serveSignup() {
  fetchMock.mockImplementation((input) => {
    const url = String(input);

    if (url.endsWith("/api/signup")) {
      return Promise.resolve(
        jsonResponse({ message: "User created successfully", user: USER }, 201),
      );
    }

    return Promise.resolve(jsonResponse({ token: "jwt-value", user: USER }));
  });
}

beforeEach(() => {
  clearSession();
});

afterEach(() => {
  clearSession();
});

describe("SignupPage", () => {
  it("creates the account, signs in, and returns to the page that asked", async () => {
    serveSignup();

    renderSignup("/signup?next=%2Fcart");

    const user = userEvent.setup();
    await user.type(screen.getByLabelText(/^Name/), "Ada");
    await user.type(screen.getByLabelText(/^Email/), "ada@example.com");
    await user.type(screen.getByLabelText(/^Password/), "correct horse");
    await user.click(screen.getByRole("button", { name: "Create account" }));

    expect(await screen.findByText("the cart")).not.toBeNull();

    const [signup, login] = fetchMock.mock.calls;
    expect(String(signup[0])).toMatch(/\/api\/signup$/);
    expect(JSON.parse(String(signup[1]?.body))).toEqual({
      name: "Ada",
      email: "ada@example.com",
      password: "correct horse",
    });
    expect(String(login[0])).toMatch(/\/api\/login$/);
    expect(getSession()?.token).toBe("jwt-value");
  });

  it("shows why the account was refused and does not sign in", async () => {
    fetchMock.mockResolvedValue(
      jsonResponse({ message: "user already exists" }, 400),
    );

    renderSignup();

    const user = userEvent.setup();
    await user.type(screen.getByLabelText(/^Name/), "Ada");
    await user.type(screen.getByLabelText(/^Email/), "ada@example.com");
    await user.type(screen.getByLabelText(/^Password/), "correct horse");
    await user.click(screen.getByRole("button", { name: "Create account" }));

    expect(await screen.findByText("user already exists")).not.toBeNull();
    expect(fetchMock).toHaveBeenCalledTimes(1);
    expect(getSession()).toBeNull();
  });

  it("wants eight characters before it sends anything", async () => {
    renderSignup();

    const user = userEvent.setup();
    await user.type(screen.getByLabelText(/^Name/), "Ada");
    await user.type(screen.getByLabelText(/^Email/), "ada@example.com");
    await user.type(screen.getByLabelText(/^Password/), "short");
    await user.click(screen.getByRole("button", { name: "Create account" }));

    expect(await screen.findByText("At least 8 characters")).not.toBeNull();
    expect(fetchMock).not.toHaveBeenCalled();
  });

  it("holds the form while the request is in flight", async () => {
    const gate = deferred<Response>();
    fetchMock.mockImplementation((input) =>
      String(input).endsWith("/api/signup")
        ? gate.promise
        : Promise.resolve(jsonResponse({ token: "jwt-value", user: USER })),
    );

    renderSignup("/signup?next=%2Fcart");

    const user = userEvent.setup();
    await user.type(screen.getByLabelText(/^Name/), "Ada");
    await user.type(screen.getByLabelText(/^Email/), "ada@example.com");
    const password = screen.getByLabelText(/^Password/);
    await user.type(password, "correct horse");

    const button = screen.getByRole("button", { name: "Create account" });
    await user.click(button);

    expect(button.hasAttribute("disabled")).toBe(true);

    await user.click(button);
    await user.type(password, "{Enter}");

    expect(fetchMock).toHaveBeenCalledTimes(1);

    gate.resolve(
      jsonResponse({ message: "User created successfully", user: USER }, 201),
    );

    expect(await screen.findByText("the cart")).not.toBeNull();
    expect(fetchMock).toHaveBeenCalledTimes(2);
  });

  it("submits from the keyboard when nothing is pending", async () => {
    serveSignup();

    renderSignup("/signup?next=%2Fcart");

    const user = userEvent.setup();
    await user.type(screen.getByLabelText(/^Name/), "Ada");
    await user.type(screen.getByLabelText(/^Email/), "ada@example.com");
    await user.type(screen.getByLabelText(/^Password/), "correct horse{Enter}");

    expect(await screen.findByText("the cart")).not.toBeNull();
    expect(fetchMock).toHaveBeenCalledTimes(2);
  });
});
