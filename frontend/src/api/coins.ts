import { request } from "./client";
import type { ChargeResponse, TransactionsResponse } from "./types";

export interface ListTransactionsParams {
  page?: number;
  limit?: number;
}

export interface ChargeCoinsInput {
  amount: number;
  description: string;
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

export function chargeCoins(
  token: string,
  input: ChargeCoinsInput,
): Promise<ChargeResponse> {
  return request<ChargeResponse>("/api/coins/charge", {
    method: "POST",
    token,
    body: input,
  });
}
