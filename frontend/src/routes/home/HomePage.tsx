import {
  Alert,
  Anchor,
  Button,
  Divider,
  Group,
  SimpleGrid,
  Skeleton,
  Stack,
  Text,
  Title,
} from "@mantine/core";
import { useQuery } from "@tanstack/react-query";
import { Link } from "react-router";

import { productsQuery } from "../../api/queries";
import { useSlowHint } from "../../lib/useSlowHint";
import { ProductCard } from "../products/ProductCard";

const FEATURED = 3;
const SKELETONS = [0, 1, 2];

const POINTS = [
  {
    title: "Priced on the server",
    text: "The browser sends product IDs and quantities. Amounts, coin totals and the order number are computed in Go.",
  },
  {
    title: "Stock and coins move together",
    text: "Checkout locks the user row first, then the products by ascending ID, so the last unit cannot be sold twice.",
  },
  {
    title: "Tested against Postgres",
    text: "Integration tests run on a real database, including two concurrent checkouts of the final unit.",
  },
];

export function HomePage() {
  return (
    <Stack gap={56} pb="xl">
      <Hero />
      <Featured />
      <Stack gap="lg">
        <Divider />
        <SimpleGrid cols={{ base: 1, sm: 3 }} spacing="xl">
          {POINTS.map((point) => (
            <Stack key={point.title} gap={6}>
              <Text fw={600}>{point.title}</Text>
              <Text size="sm" c="dimmed">
                {point.text}
              </Text>
            </Stack>
          ))}
        </SimpleGrid>
      </Stack>
      <Text size="xs" c="dimmed">
        The API runs on Render&apos;s free plan, so the first request after
        fifteen idle minutes takes about twenty seconds.
      </Text>
    </Stack>
  );
}

function Hero() {
  return (
    <Stack gap="lg" maw={680} pt="xl">
      <Text size="sm" fw={600} c="teal.8" tt="uppercase">
        go-next
      </Text>
      <Title order={1}>A small shop, priced on the server</Title>
      <Text c="dimmed" size="lg">
        A Go API behind a React client. The server owns the prices, the stock
        and the coin balance; the browser only chooses what to buy.
      </Text>
      <Group>
        <Button component={Link} to="/products" color="teal" size="md">
          Browse the catalogue
        </Button>
        <Text size="sm" c="dimmed">
          Signup starts you at 1,000 coins.
        </Text>
      </Group>
    </Stack>
  );
}

function Featured() {
  const { data, isPending, isError, error, refetch, isFetching } = useQuery(
    productsQuery(),
  );
  const isSlow = useSlowHint(isPending);
  const featured = data?.slice(0, FEATURED);

  return (
    <Stack gap="md">
      <Group justify="space-between" align="baseline">
        <Title order={2} size="h3">
          Featured
        </Title>
        <Anchor component={Link} to="/products" size="sm" c="dimmed">
          All products
        </Anchor>
      </Group>

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
          {SKELETONS.map((key) => (
            <Skeleton key={key} height={300} radius="md" />
          ))}
        </SimpleGrid>
      )}

      {featured && featured.length > 0 && (
        <SimpleGrid cols={{ base: 1, xs: 2, md: 3 }} spacing="lg">
          {featured.map((product) => (
            <ProductCard key={product.id} product={product} />
          ))}
        </SimpleGrid>
      )}
    </Stack>
  );
}
