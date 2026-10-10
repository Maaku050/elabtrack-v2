export function shortReference(value: string) { return value.length > 22 ? `${value.slice(0, 11)}…${value.slice(-5)}` : value }
