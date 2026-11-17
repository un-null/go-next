export const DEFAULT_NEXT = "/products";

export function safeNext(value: string | null): string {
  if (value?.startsWith("/") && !value.startsWith("//")) {
    return value;
  }
  return DEFAULT_NEXT;
}

export function withNext(path: string, next: string): string {
  if (next === DEFAULT_NEXT) {
    return path;
  }
  return `${path}?next=${encodeURIComponent(next)}`;
}
