import { useQuery } from "@tanstack/react-query";

import { userQuery } from "../api/queries";
import type { User } from "../api/types";
import { useSession } from "./session";

export function useUser(): User | null {
  const session = useSession();

  const query = useQuery({
    ...userQuery(session?.token ?? "", session?.user.id ?? ""),
    enabled: session !== null,
  });

  return query.data ?? session?.user ?? null;
}
