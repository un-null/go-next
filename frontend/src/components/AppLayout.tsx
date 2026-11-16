import { Anchor, AppShell, Container, Group, Menu, Text } from "@mantine/core";
import { useQueryClient } from "@tanstack/react-query";
import { Link, Outlet, useNavigate } from "react-router";

import { clearSession, useSession } from "../lib/session";

const NAV = {
  size: "xs",
  tt: "uppercase",
  ff: "monospace",
  c: "dimmed",
} as const;

export function AppLayout() {
  const session = useSession();
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

              {session ? (
                <Menu position="bottom-end" width={180}>
                  <Menu.Target>
                    <Anchor
                      component="button"
                      type="button"
                      {...NAV}
                      c="dark.9"
                    >
                      {session.user.name}
                    </Anchor>
                  </Menu.Target>
                  <Menu.Dropdown>
                    <Menu.Label>
                      {session.user.coins.toLocaleString()} coins
                    </Menu.Label>
                    <Menu.Item onClick={signOut}>Sign out</Menu.Item>
                  </Menu.Dropdown>
                </Menu>
              ) : (
                <>
                  <Anchor component={Link} to="/login" {...NAV}>
                    Sign in
                  </Anchor>
                  <Anchor component={Link} to="/signup" {...NAV}>
                    Sign up
                  </Anchor>
                </>
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
