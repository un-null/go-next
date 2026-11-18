import {
  Alert,
  Button,
  CloseButton,
  Group,
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

import { productsQuery } from "../../api/queries";
import { useSlowHint } from "../../lib/useSlowHint";
import { ProductCard } from "./ProductCard";

const SKELETONS = [0, 1, 2, 3, 4, 5];
const DEBOUNCE_MS = 300;

export function ProductsPage() {
  const [params, setParams] = useSearchParams();
  const q = params.get("q") ?? "";

  const [debounced] = useDebouncedValue(q, DEBOUNCE_MS);

  const { data, isPending, isError, error, refetch, isFetching } = useQuery(
    productsQuery(debounced),
  );
  const isSlow = useSlowHint(isPending);

  function setQ(value: string) {
    const next = new URLSearchParams(params);
    if (value) {
      next.set("q", value);
    } else {
      next.delete("q");
    }
    setParams(next, { replace: true });
  }

  return (
    <Stack gap="lg">
      <div>
        <Title order={1} size="h2">
          Products
        </Title>
      </div>

      <TextInput
        placeholder="Search products"
        aria-label="Search products"
        leftSection={<IconSearch size={16} stroke={1.5} />}
        rightSection={
          q ? (
            <CloseButton aria-label="Clear search" onClick={() => setQ("")} />
          ) : null
        }
        value={q}
        onChange={(event) => setQ(event.currentTarget.value)}
      />

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
              : "The catalogue is empty."}
          </Text>
          {debounced && (
            <Button variant="default" size="xs" onClick={() => setQ("")}>
              Clear search
            </Button>
          )}
        </Stack>
      )}

      {data && data.length > 0 && (
        <>
          <Group justify="space-between">
            <Text size="sm" c="dimmed">
              {data.length} products
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
