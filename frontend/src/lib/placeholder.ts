const FILL = "#e9ecef";

export function placeholder(width: number, height: number): string {
  const svg = `<svg xmlns="http://www.w3.org/2000/svg" width="${width}" height="${height}"><rect width="${width}" height="${height}" fill="${FILL}"/></svg>`;
  return `data:image/svg+xml,${encodeURIComponent(svg)}`;
}
