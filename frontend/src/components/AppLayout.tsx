import {
  ActionIcon,
  Anchor,
  AppShell,
  Container,
  Group,
  Indicator,
  Menu,
  Text,
} from "@mantine/core";
import {
  IconListDetails,
  IconLogout,
  IconShoppingCart,
  IconUserCircle,
} from "@tabler/icons-react";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import BoringAvatar from "boring-avatars";
import { Link, Outlet, useNavigate } from "react-router";

import { cartQuery } from "../api/queries";
import { clearSession, useSession } from "../lib/session";
import { useUser } from "../lib/useUser";

const NAV = {
  size: "xs",
  tt: "uppercase",
  ff: "monospace",
  c: "dimmed",
} as const;

export function AppLayout() {
  const session = useSession();
  const user = useUser();
  const navigate = useNavigate();
  const queryClient = useQueryClient();

  const cart = useQuery({
    ...cartQuery(session?.token ?? ""),
    enabled: session !== null,
  });
  const cartCount = cart.data?.total_items ?? 0;

  function signOut() {
    clearSession();
    queryClient.clear();
    navigate("/");
  }

  return (
    <AppShell header={{ height: 56 }} padding="md">
      <AppShell.Header>
        <Container size="lg" h="100%">
          <Group h="100%" justify="space-between">
            <Link to="/" style={{ textDecoration: "none" }}>
              <Text fw={700} c="dark.9">
                go-next
              </Text>
            </Link>

            <Group gap="lg">
              <Anchor component={Link} to="/products" {...NAV}>
                Products
              </Anchor>

              {session && (
                <Indicator
                  label={cartCount}
                  size={16}
                  offset={4}
                  disabled={cartCount === 0}
                  color="dark"
                  withBorder
                  styles={{
                    root: { display: "flex", alignItems: "center" },
                    indicator: {
                      width: 18,
                      height: 18,
                      minWidth: 18,
                      padding: 0,
                      borderRadius: "50%",
                      display: "flex",
                      alignItems: "center",
                      justifyContent: "center",
                    },
                  }}
                >
                  <ActionIcon
                    component={Link}
                    to="/cart"
                    variant="subtle"
                    color="dark"
                    size="md"
                    aria-label={
                      cartCount > 0 ? `Cart, ${cartCount} items` : "Cart"
                    }
                  >
                    <IconShoppingCart size={20} stroke={1.5} />
                  </ActionIcon>
                </Indicator>
              )}

              {session ? (
                <Menu position="bottom-end" width={180}>
                  <Menu.Target>
                    <ActionIcon
                      component="button"
                      type="button"
                      variant="subtle"
                      color="dark"
                      size="md"
                      aria-label={user?.name ?? "Account"}
                    >
                      <BoringAvatar
                        name={user?.name ?? ""}
                        size={24}
                        variant="beam"
                      />
                    </ActionIcon>
                  </Menu.Target>
                  <Menu.Dropdown>
                    <Menu.Label>
                      {user?.coins.toLocaleString()} coins
                    </Menu.Label>
                    <Menu.Item
                      component={Link}
                      to="/orders"
                      leftSection={<IconListDetails size={16} stroke={1.5} />}
                    >
                      Orders
                    </Menu.Item>
                    <Menu.Item
                      component={Link}
                      to="/account"
                      leftSection={<IconUserCircle size={16} stroke={1.5} />}
                    >
                      Account
                    </Menu.Item>
                    <Menu.Divider />
                    <Menu.Item
                      color="red"
                      leftSection={<IconLogout size={16} stroke={1.5} />}
                      onClick={signOut}
                    >
                      Sign out
                    </Menu.Item>
                  </Menu.Dropdown>
                </Menu>
              ) : (
                <Anchor component={Link} to="/login" {...NAV}>
                  Sign in
                </Anchor>
              )}
            </Group>
          </Group>
        </Container>
      </AppShell.Header>

      <AppShell.Main>
        <Container size="lg">
          <Outlet />
        </Container>
      </AppShell.Main>
    </AppShell>
  );
}
