import { Button } from "@/components/ui/button";

export function Welcome({ message }: { message: string }) {
  return (
    <main className="flex items-center justify-center pt-16 pb-4">
      <Button>{message}</Button>
    </main>
  );
}
