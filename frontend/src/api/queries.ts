import { queryOptions } from "@tanstack/react-query";

import { getProduct, listProducts } from "./catalog";

export const productKeys = {
  all: ["products"] as const,
  list: () => [...productKeys.all, "list"] as const,
  detail: (id: number) => [...productKeys.all, "detail", id] as const,
};

export const productsQuery = () =>
  queryOptions({
    queryKey: productKeys.list(),
    queryFn: () => listProducts({ limit: 100 }),
  });

export const productQuery = (id: number) =>
  queryOptions({
    queryKey: productKeys.detail(id),
    queryFn: () => getProduct(id),
  });
