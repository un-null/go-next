import { AspectRatio, Badge, Card, Group, Image, Text } from "@mantine/core";
import { Link } from "react-router";

import type { Product } from "../../api/types";
import { formatPrice, thumbnailUrl } from "../../lib/format";
import { placeholder } from "../../lib/placeholder";
import classes from "./ProductCard.module.css";

const PLACEHOLDER = placeholder(600, 400);

interface ProductCardProps {
  product: Product;
}

export function ProductCard({ product }: ProductCardProps) {
  const inStock = product.stock_quantity > 0;

  return (
    <Card
      className={classes.card}
      component={Link}
      to={`/products/${product.id}`}
      padding="md"
      withBorder
    >
      <Card.Section>
        <AspectRatio ratio={3 / 2} bg="gray.1">
          <Image
            src={thumbnailUrl(product.image_url)}
            fallbackSrc={PLACEHOLDER}
            alt=""
            w="100%"
            h="100%"
            fit="cover"
            loading="lazy"
          />
        </AspectRatio>
      </Card.Section>

      <Group
        justify="space-between"
        align="flex-start"
        mt="md"
        gap="xs"
        wrap="nowrap"
      >
        <Text fw={600} lineClamp={1} style={{ flex: 1, minWidth: 0 }}>
          {product.name}
        </Text>
        <Text fw={700} ff="monospace" style={{ flexShrink: 0 }}>
          {formatPrice(product.price)}
        </Text>
      </Group>

      <Text size="sm" c="dimmed" lineClamp={2} mt={4}>
        {product.description}
      </Text>

      <Group justify="space-between" mt="auto" pt="md">
        <Badge variant="light" color={inStock ? "teal" : "gray"}>
          {inStock ? `${product.stock_quantity} in stock` : "Out of stock"}
        </Badge>
        <Text size="sm" c="dimmed">
          {product.average_rating.toFixed(1)} stars
        </Text>
      </Group>
    </Card>
  );
}
