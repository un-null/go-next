import {
  Alert,
  Button,
  Group,
  SimpleGrid,
  Skeleton,
  Stack,
  Text,
  Title,
} from "@mantine/core";
import { useQuery } from "@tanstack/react-query";

import { productsQuery } from "../../api/queries";
import { useSlowHint } from "../../lib/useSlowHint";
import { ProductCard } from "./ProductCard";

const SKELETONS = [0, 1, 2, 3, 4, 5];

export function ProductsPage() {
  const { data, isPending, isError, error, refetch, isFetching } = useQuery(
    productsQuery(),
  );
  const isSlow = useSlowHint(isPending);

  return (
    <Stack gap="lg">
      <div>
        <Title order={1} size="h2">
          Products
        </Title>
        <Text c="dimmed" mt={4}>
          Twenty items, priced in dollars and paid in coins.
        </Text>
      </div>

      {isSlow && (
        <Alert color="yellow" variant="light" title="Waking the API up">
          The free instance (Render) sleeps after fifteen minutes. The first request can take about 20 sec. It is working, not broken.
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
          {SKELETONS.map((key) => (
            <Skeleton key={key} height={300} radius="md" />
          ))}
        </SimpleGrid>
      )}

      {data && data.length === 0 && (
        <Text c="dimmed">The catalogue is empty.</Text>
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
