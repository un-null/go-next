import { keepPreviousData, queryOptions } from "@tanstack/react-query";

import { getCart } from "./cart";
import { getProduct, listCategories, listProducts } from "./catalog";
import { listTransactions } from "./coins";
import { getOrder, listOrders } from "./orders";
import { getUser } from "./users";

export const productKeys = {
  all: ["products"] as const,
  list: (q = "", category: number | null = null) =>
    [...productKeys.all, "list", q, category] as const,
  detail: (id: number) => [...productKeys.all, "detail", id] as const,
};

export const productsQuery = (q = "", category: number | null = null) =>
  queryOptions({
    queryKey: productKeys.list(q, category),
    queryFn: () =>
      listProducts({ limit: 100, q, category: category ?? undefined }),
    placeholderData: keepPreviousData,
  });

export const productQuery = (id: number) =>
  queryOptions({
    queryKey: productKeys.detail(id),
    queryFn: () => getProduct(id),
  });

export const categoryKeys = {
  all: ["categories"] as const,
};

export const categoriesQuery = () =>
  queryOptions({
    queryKey: categoryKeys.all,
    queryFn: listCategories,
  });

export const cartKeys = {
  all: ["cart"] as const,
};

export const cartQuery = (token: string) =>
  queryOptions({
    queryKey: cartKeys.all,
    queryFn: () => getCart(token),
  });

export const orderKeys = {
  all: ["orders"] as const,
  list: (userId: string) => [...orderKeys.all, "list", userId] as const,
  detail: (id: number) => [...orderKeys.all, "detail", id] as const,
};

export const ordersQuery = (token: string, userId: string) =>
  queryOptions({
    queryKey: orderKeys.list(userId),
    queryFn: () => listOrders(token, userId),
  });

export const orderQuery = (token: string, id: number) =>
  queryOptions({
    queryKey: orderKeys.detail(id),
    queryFn: () => getOrder(token, id),
  });

export const userKeys = {
  all: ["user"] as const,
  detail: (id: string) => [...userKeys.all, id] as const,
};

export const userQuery = (token: string, id: string) =>
  queryOptions({
    queryKey: userKeys.detail(id),
    queryFn: () => getUser(token, id),
  });

export const transactionKeys = {
  all: ["transactions"] as const,
};

export const transactionsQuery = (token: string) =>
  queryOptions({
    queryKey: transactionKeys.all,
    queryFn: () => listTransactions(token, { limit: 100 }),
  });
