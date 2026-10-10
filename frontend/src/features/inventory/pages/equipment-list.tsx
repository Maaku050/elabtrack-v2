import { CatalogCart } from '@/features/borrowing/pages/catalog-cart'
import { InventoryDirectory } from './inventory-directory'
export function BorrowerCatalog() { return <CatalogCart /> }
export function EquipmentList({borrower=false}:{borrower?:boolean}){return borrower?<BorrowerCatalog />:<InventoryDirectory />}
