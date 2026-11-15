import { request } from "./client";
import type { Category, Product } from "./types";

export interface ListProductsParams {
  page?: number;
  limit?: number;
}

export function listProducts({
  page = 1,
  limit = 10,
}: ListProductsParams = {}): Promise<Product[]> {
  const query = new URLSearchParams({
    page: String(page),
    limit: String(limit),
  });
  return request<Product[]>(`/api/products?${query}`);
}

export function getProduct(id: number): Promise<Product> {
  return request<Product>(`/api/products/${id}`);
}

export function listCategories(): Promise<Category[]> {
  return request<Category[]>("/api/categories");
}
