import {
  Alert,
  Button,
  CloseButton,
  Flex,
  Group,
  NativeSelect,
  SimpleGrid,
  Skeleton,
  Stack,
  Text,
  TextInput,
  Title,
} from "@mantine/core";
import { useDebouncedValue } from "@mantine/hooks";
import { IconSearch } from "@tabler/icons-react";
import { useQuery } from "@tanstack/react-query";
import { useSearchParams } from "react-router";

import { categoriesQuery, productsQuery } from "../../api/queries";
import { useSlowHint } from "../../lib/useSlowHint";
import { ProductCard } from "./ProductCard";

const SKELETONS = [0, 1, 2, 3, 4, 5];
const DEBOUNCE_MS = 300;

function positiveParam(raw: string | null): number | null {
  const value = Number(raw);
  return Number.isInteger(value) && value > 0 ? value : null;
}

export function ProductsPage() {
  const [params, setParams] = useSearchParams();
  const q = params.get("q") ?? "";
  const requested = positiveParam(params.get("category"));

  const categories = useQuery(categoriesQuery());

  const category = categories.data
    ? requested !== null &&
      categories.data.some((item) => item.id === requested)
      ? requested
      : null
    : requested;

  const [debounced] = useDebouncedValue(q, DEBOUNCE_MS);

  const { data, isPending, isError, error, refetch, isFetching } = useQuery(
    productsQuery(debounced, category),
  );
  const isSlow = useSlowHint(isPending);
  const isFiltered = Boolean(debounced) || category !== null;

  function setParam(key: string, value: string) {
    const next = new URLSearchParams(params);
    if (value) {
      next.set(key, value);
    } else {
      next.delete(key);
    }
    setParams(next, { replace: true });
  }

  function clearFilters() {
    const next = new URLSearchParams(params);
    next.delete("q");
    next.delete("category");
    setParams(next, { replace: true });
  }

  return (
    <Stack gap="lg">
      <div>
        <Title order={1} size="h2">
          Products
        </Title>
      </div>

      <Flex gap="sm" direction={{ base: "column", sm: "row" }}>
        <TextInput
          style={{ flex: 1 }}
          placeholder="Search products"
          aria-label="Search products"
          leftSection={<IconSearch size={16} stroke={1.5} />}
          rightSection={
            q ? (
              <CloseButton
                aria-label="Clear search"
                onClick={() => setParam("q", "")}
              />
            ) : null
          }
          value={q}
          onChange={(event) => setParam("q", event.currentTarget.value)}
        />

        <NativeSelect
          aria-label="Category"
          w={{ base: "100%", sm: 220 }}
          data={[
            { value: "", label: "All categories" },
            ...(categories.data ?? []).map((item) => ({
              value: String(item.id),
              label: item.name,
            })),
          ]}
          value={category === null ? "" : String(category)}
          onChange={(event) => setParam("category", event.currentTarget.value)}
        />
      </Flex>

      {isSlow && (
        <Alert color="yellow" variant="light" title="Waking the API up">
          The free instance (Render) sleeps after fifteen minutes. The first
          request can take about 20 sec. It is working, not broken.
        </Alert>
      )}

      {isError && (
        <Alert color="red" variant="light" title="Could not load the catalogue">
          <Text size="sm">{error.message}</Text>
          <Button
            mt="sm"
            size="xs"
            variant="light"
            color="red"
            onClick={() => refetch()}
            loading={isFetching}
          >
            Try again
          </Button>
        </Alert>
      )}

      {isPending && (
        <SimpleGrid cols={{ base: 1, xs: 2, md: 3 }} spacing="lg">
          {SKELETONS.map((skeleton) => (
            <Skeleton key={skeleton} height={300} radius="md" />
          ))}
        </SimpleGrid>
      )}

      {!isPending && data && data.length === 0 && (
        <Stack gap="xs" align="flex-start">
          <Text c="dimmed">
            {debounced
              ? `No products match "${debounced}".`
              : category !== null
                ? "No products in this category."
                : "The catalogue is empty."}
          </Text>
          {isFiltered && (
            <Button variant="default" size="xs" onClick={clearFilters}>
              Clear filters
            </Button>
          )}
        </Stack>
      )}

      {data && data.length > 0 && (
        <>
          <Group justify="space-between">
            <Text size="sm" c="dimmed">
              {data.length} {data.length === 1 ? "product" : "products"}
            </Text>
          </Group>

          <SimpleGrid cols={{ base: 1, xs: 2, md: 3 }} spacing="lg">
            {data.map((product) => (
              <ProductCard key={product.id} product={product} />
            ))}
          </SimpleGrid>
        </>
      )}
    </Stack>
  );
}
