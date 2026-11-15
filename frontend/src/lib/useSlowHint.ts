import { useEffect, useState } from "react";

export function useSlowHint(isPending: boolean, delay = 3000): boolean {
  const [isSlow, setIsSlow] = useState(false);

  useEffect(() => {
    if (!isPending) {
      setIsSlow(false);
      return;
    }

    const timer = setTimeout(() => setIsSlow(true), delay);
    return () => clearTimeout(timer);
  }, [isPending, delay]);

  return isSlow;
}
