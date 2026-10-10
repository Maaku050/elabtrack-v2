export const pesos = (minor: number) => new Intl.NumberFormat('en-PH', { style: 'currency', currency: 'PHP' }).format(minor / 100)
