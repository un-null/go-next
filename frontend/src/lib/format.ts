const priceFormatter = new Intl.NumberFormat("en-US", {
  style: "currency",
  currency: "USD",
});

export function formatPrice(amount: number): string {
  return priceFormatter.format(amount);
}

export function thumbnailUrl(url: string, width = 600): string {
  if (!url) {
    return "";
  }
  try {
    const parsed = new URL(url);
    parsed.searchParams.set("w", String(width));
    parsed.searchParams.set("auto", "format");
    parsed.searchParams.set("q", "70");
    return parsed.toString();
  } catch {
    return url;
  }
}
