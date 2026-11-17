import { request } from "./client";
import type { TransactionsResponse } from "./types";

export interface ListTransactionsParams {
  page?: number;
  limit?: number;
}

export function listTransactions(
  token: string,
  { page = 1, limit = 100 }: ListTransactionsParams = {},
): Promise<TransactionsResponse> {
  const query = new URLSearchParams({
    page: String(page),
    limit: String(limit),
  });
  return request<TransactionsResponse>(`/api/coins/transactions?${query}`, {
    token,
  });
}
