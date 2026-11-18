import { Anchor, AppShell, Container, Group, Menu, Text } from "@mantine/core";
import { useQueryClient } from "@tanstack/react-query";
import BoringAvatar from "boring-avatars";
import { Link, Outlet, useNavigate } from "react-router";

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
                <Anchor component={Link} to="/cart" {...NAV}>
                  Cart
                </Anchor>
              )}

              {session ? (
                <Menu position="bottom-end" width={180}>
                  <Menu.Target>
                    <Anchor
                      component="button"
                      type="button"
                      aria-label={user?.name ?? "Account"}
                      style={{ display: "flex" }}
                    >
                      <BoringAvatar
                        name={user?.name ?? ""}
                        size={28}
                        variant="marble"
                      />
                    </Anchor>
                  </Menu.Target>
                  <Menu.Dropdown>
                    <Menu.Label>
                      {user?.coins.toLocaleString()} coins
                    </Menu.Label>
                    <Menu.Item component={Link} to="/orders">
                      Orders
                    </Menu.Item>
                    <Menu.Item component={Link} to="/account">
                      Account
                    </Menu.Item>
                    <Menu.Item onClick={signOut}>Sign out</Menu.Item>
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
