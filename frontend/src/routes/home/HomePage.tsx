import {
  Alert,
  Anchor,
  Button,
  Group,
  SimpleGrid,
  Skeleton,
  Stack,
  Text,
  ThemeIcon,
  Title,
} from "@mantine/core";
import {
  type Icon,
  IconReceipt,
  IconShoppingBag,
  IconShoppingCart,
  IconUserPlus,
} from "@tabler/icons-react";
import { useQuery } from "@tanstack/react-query";
import { Link } from "react-router";

import { productsQuery } from "../../api/queries";
import { useSession } from "../../lib/session";
import { useSlowHint } from "../../lib/useSlowHint";
import { ProductCard } from "../products/ProductCard";

const FEATURED = 3;
const SKELETONS = [0, 1, 2];

const STEPS: {
  icon: Icon;
  color: "teal" | "blue" | "grape" | "orange";
  title: string;
  text: string;
}[] = [
  {
    icon: IconUserPlus,
    color: "teal",
    title: "Create an account",
    text: "You start with 1,000 coins. There is no card to enter.",
  },
  {
    icon: IconShoppingBag,
    color: "blue",
    title: "Pick your items",
    text: "Every price is in dollars and in coins. One coin buys one dollar, so $49.99 costs 50.",
  },
  {
    icon: IconShoppingCart,
    color: "grape",
    title: "Fill your cart",
    text: "Your cart is saved with your account, so it is still there when you come back.",
  },
  {
    icon: IconReceipt,
    color: "orange",
    title: "Check out",
    text: "Spend the coins and the item is yours. If someone else got the last one, your coins stay.",
  },
];

export function HomePage() {
  const session = useSession();

  return (
    <Stack gap={56} pb="xl">
      <HowItWorks />
      <Stack gap="sm" align="center">
        {session ? (
          <Button component={Link} to="/products" size="md">
            Start shopping
          </Button>
        ) : (
          <>
            <Button component={Link} to="/signup" size="md">
              Create an account
            </Button>
            <Text size="sm">
              Already have an account?{" "}
              <Anchor component={Link} to="/login">
                Sign in
              </Anchor>
            </Text>
          </>
        )}
      </Stack>
      <Featured />
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
      <Title order={2} size="h4">
        Featured
      </Title>

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
            <Skeleton key={key} height={300} radius="sm" />
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

      <Group justify="center" mt="xs">
        <Anchor component={Link} to="/products">
          See all products
        </Anchor>
      </Group>
    </Stack>
  );
}

function HowItWorks() {
  return (
    <Stack gap="xl">
      <Title order={1} size="h3">
        How buying works
      </Title>
      <SimpleGrid cols={{ base: 1, sm: 2, md: 4 }} spacing="xl">
        {STEPS.map((step) => (
          <Stack key={step.title} gap="xs">
            <ThemeIcon variant="light" color={step.color} size="lg" radius="sm">
              <step.icon size={20} stroke={1.5} />
            </ThemeIcon>
            <Text fw={600}>{step.title}</Text>
            <Text size="sm" c="dimmed">
              {step.text}
            </Text>
          </Stack>
        ))}
      </SimpleGrid>
    </Stack>
  );
}
