import type { EquipmentPresentation } from '@/components/application/equipment'
import spoon from './assets/wooden-spoon.png'
import chafing from './assets/chafing-dish.png'
import muffin from './assets/muffin-tray.png'
import blender from './assets/blender.png'
import tongs from './assets/tongs.png'
import pan from './assets/food-pan.png'

// Synthetic, development-only presentation fixtures. Never fetch or persist these.
export const equipment: readonly EquipmentPresentation[] = [
  { id: 'spoon', name: 'Wooden Spoon', category: 'Utensils', unitLabel: '1 pc', tags: ['Wooden', 'Standard 12″'], image: spoon, available: 24, availability: 'available' },
  { id: 'chafing', name: 'Chafing Dish Set', category: 'Cooking', unitLabel: '2 pcs', tags: ['Stainless steel', 'Fuel holders'], image: chafing, available: 6, availability: 'available' },
  { id: 'muffin', name: 'Muffin Tray', category: 'Baking', unitLabel: '1 pc', tags: ['Non-stick', '12 cup'], image: muffin, available: 2, availability: 'attention' },
  { id: 'blender', name: 'Commercial Blender', category: 'Appliances', unitLabel: '1 pc', tags: ['High performance', '64 oz'], image: blender, available: 5, availability: 'available' },
  { id: 'tongs', name: 'Serving Tongs', category: 'Serving', unitLabel: '1 pc', tags: ['Stainless steel', '12″'], image: tongs, available: 18, availability: 'available' },
  { id: 'pan', name: 'Food Pan Set', category: 'Serving', unitLabel: '6 pcs', tags: ['Stainless steel', '1/3 size'], image: pan, available: 0, availability: 'unavailable' },
]
export const categories = ['All', 'Cooking', 'Baking', 'Serving', 'Utensils', 'Appliances'] as const
export const requestTrends = [
  { day: 'Oct 2', issued: 4, pending: 2, denied: 1 }, { day: 'Oct 3', issued: 6, pending: 3, denied: 1 },
  { day: 'Oct 4', issued: 5, pending: 2, denied: 0 }, { day: 'Oct 5', issued: 9, pending: 3, denied: 1 },
  { day: 'Oct 6', issued: 7, pending: 2, denied: 1 }, { day: 'Oct 7', issued: 4, pending: 1, denied: 1 },
  { day: 'Oct 8', issued: 3, pending: 1, denied: 0 },
] as const
export const stockSummary = [
  { label: 'Available', count: 82, tone: 'success' }, { label: 'Checked out', count: 28, tone: 'primary' },
  { label: 'Damaged held', count: 8, tone: 'danger' }, { label: 'Reserved', count: 6, tone: 'warning' },
] as const
export const pendingRequests = [
  { reference: 'PREVIEW-0215', name: 'Alex Rivera', initials: 'AR', category: 'Student', program: 'Culinary Arts', submitted: '08 Oct 2026', time: '09:15 AM', context: 'Kitchen laboratory', itemIds: ['spoon', 'chafing'], units: 4, due: '09 Oct, 4:00 PM', fine: 'PHP 10' },
  { reference: 'PREVIEW-0214', name: 'Jordan Cruz', initials: 'JC', category: 'Faculty', program: 'Hospitality', submitted: '08 Oct 2026', time: '08:40 AM', context: 'Training session', itemIds: ['muffin', 'blender'], units: 3, due: '10 Oct, 2:00 PM', fine: null },
  { reference: 'PREVIEW-0213', name: 'Casey Santos', initials: 'CS', category: 'Student', program: 'Hospitality', submitted: '08 Oct 2026', time: '08:10 AM', context: 'Food preparation', itemIds: ['tongs', 'chafing'], units: 6, due: '09 Oct, 3:00 PM', fine: null },
  { reference: 'PREVIEW-0212', name: 'Taylor Reyes', initials: 'TR', category: 'Student', program: 'Culinary Arts', submitted: '07 Oct 2026', time: '03:20 PM', context: 'Baking practice', itemIds: ['muffin', 'spoon'], units: 5, due: '09 Oct, 1:00 PM', fine: null },
  { reference: 'PREVIEW-0211', name: 'Morgan Diaz', initials: 'MD', category: 'Faculty', program: 'Culinary Arts', submitted: '07 Oct 2026', time: '02:30 PM', context: 'Demonstration', itemIds: ['blender', 'tongs'], units: 3, due: '10 Oct, 4:00 PM', fine: null },
  { reference: 'PREVIEW-0210', name: 'Riley Flores', initials: 'RF', category: 'Student', program: 'Hospitality', submitted: '07 Oct 2026', time: '02:15 PM', context: 'Kitchen laboratory', itemIds: ['spoon'], units: 2, due: '09 Oct, 4:00 PM', fine: null },
  { reference: 'PREVIEW-0209', name: 'Avery Lee', initials: 'AL', category: 'Student', program: 'Culinary Arts', submitted: '07 Oct 2026', time: '01:45 PM', context: 'Food preparation', itemIds: ['tongs'], units: 3, due: '09 Oct, 5:00 PM', fine: null },
  { reference: 'PREVIEW-0208', name: 'Jamie Silva', initials: 'JS', category: 'Faculty', program: 'Hospitality', submitted: '07 Oct 2026', time: '01:00 PM', context: 'Demonstration', itemIds: ['chafing'], units: 1, due: '10 Oct, 2:00 PM', fine: null },
] as const
