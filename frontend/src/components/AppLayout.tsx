import { AppShell, Container, Group, Text } from "@mantine/core";
import { Link, Outlet } from "react-router";

export function AppLayout() {
  return (
    <AppShell header={{ height: 56 }} padding="md">
      <AppShell.Header>
        <Container size="lg" h="100%">
          <Group h="100%" justify="space-between">
            <Link to="/" style={{ textDecoration: "none" }}>
              <Text fw={700} c="teal.8">
                go-next
              </Text>
            </Link>
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
