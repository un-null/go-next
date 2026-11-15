import {
  Alert,
  Anchor,
  Badge,
  Box,
  Button,
  Divider,
  Group,
  Image,
  SimpleGrid,
  Skeleton,
  Stack,
  Text,
  Title,
} from "@mantine/core";
import { useQuery } from "@tanstack/react-query";
import { Link, useParams } from "react-router";

import { ApiError } from "../../api/client";
import { productQuery } from "../../api/queries";
import { formatPrice, thumbnailUrl } from "../../lib/format";
import { useSlowHint } from "../../lib/useSlowHint";
import classes from "./ProductDetailPage.module.css";

const PLACEHOLDER = `data:image/svg+xml,${encodeURIComponent(
  '<svg xmlns="http://www.w3.org/2000/svg" width="900" height="600"><rect width="900" height="600" fill="#e9ecef"/></svg>',
)}`;

export function ProductDetailPage() {
  const { id } = useParams();
  const productId = Number(id);

  if (!Number.isInteger(productId) || productId <= 0) {
    return <NotFound />;
  }

  return <ProductDetail productId={productId} />;
}

function ProductDetail({ productId }: { productId: number }) {
  const { data, isPending, isError, error, refetch, isFetching } = useQuery(
    productQuery(productId),
  );
  const isSlow = useSlowHint(isPending);

  if (isError && error instanceof ApiError && error.status === 404) {
    return <NotFound />;
  }

  return (
    <Stack gap="lg">
      <Anchor component={Link} to="/products" size="sm" c="dimmed">
        All products
      </Anchor>

      {isSlow && (
        <Alert color="yellow" variant="light" title="Waking the API up">
          The free instance sleeps after fifteen minutes. This can take about
          twenty seconds.
        </Alert>
      )}

      {isError && (
        <Alert color="red" variant="light" title="Could not load this product">
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
        <SimpleGrid cols={{ base: 1, md: 2 }} spacing="xl">
          <Skeleton height={360} radius="md" />
          <Skeleton height={280} radius="md" />
        </SimpleGrid>
      )}

      {data && (
        <SimpleGrid
          cols={{ base: 1, md: 2 }}
          spacing="xl"
          className={classes.layout}
        >
          <Box h={360} bg="gray.1" className={classes.media}>
            <Image
              src={thumbnailUrl(data.image_url, 900)}
              fallbackSrc={PLACEHOLDER}
              alt=""
              h="100%"
              fit="cover"
            />
          </Box>

          <Stack gap="md">
            <div>
              <Title order={1} size="h2">
                {data.name}
              </Title>
              <Text c="dimmed" mt={4}>
                {data.description}
              </Text>
            </div>

            <Group gap="xs">
              <Text fw={700} fz="h3" c="teal.8">
                {formatPrice(data.price)}
              </Text>
              <Text size="sm" c="dimmed">
                from {Math.ceil(data.price)} coins
              </Text>
            </Group>

            <Group gap="xs">
              <Badge
                variant="light"
                color={data.stock_quantity > 0 ? "teal" : "gray"}
              >
                {data.stock_quantity > 0
                  ? `${data.stock_quantity} in stock`
                  : "Out of stock"}
              </Badge>
              <Text size="sm" c="dimmed">
                {data.average_rating.toFixed(1)} stars
              </Text>
            </Group>

            <Divider />

            <Text size="sm" c="dimmed">
              Checkout arrives with the cart. Coins are demo currency — signup
              starts you at 1,000.
            </Text>
          </Stack>
        </SimpleGrid>
      )}
    </Stack>
  );
}

function NotFound() {
  return (
    <Stack gap="md">
      <Title order={1} size="h2">
        Product not found
      </Title>
      <Text c="dimmed">This item does not exist, or it was removed.</Text>
      <Group>
        <Button component={Link} to="/products" variant="light" color="teal">
          Back to the catalogue
        </Button>
      </Group>
    </Stack>
  );
}
