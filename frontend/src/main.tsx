import "@mantine/core/styles.css";
import "./global.css";

import { MantineProvider } from "@mantine/core";
import {
  QueryCache,
  QueryClient,
  QueryClientProvider,
} from "@tanstack/react-query";
import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import { BrowserRouter } from "react-router";

import App from "./App";
import { ApiError } from "./api/client";
import { retryUnlessClientError } from "./api/retry";
import { clearSession } from "./lib/session";
import { theme } from "./theme";

const queryClient = new QueryClient({
  queryCache: new QueryCache({
    onError: (error) => {
      if (error instanceof ApiError && error.status === 401) {
        clearSession();
      }
    },
  }),
  defaultOptions: {
    queries: {
      retry: retryUnlessClientError,
      refetchOnWindowFocus: false,
    },
  },
});

const root = document.getElementById("root");
if (!root) {
  throw new Error("#root not found");
}

createRoot(root).render(
  <StrictMode>
    <MantineProvider theme={theme} defaultColorScheme="light">
      <QueryClientProvider client={queryClient}>
        <BrowserRouter>
          <App />
        </BrowserRouter>
      </QueryClientProvider>
    </MantineProvider>
  </StrictMode>,
);
