import { shortReference } from '@/lib/record-reference'

export function RecordReference({ value }: { value: string }) { return <span className="record-reference" title={value} aria-label={value}>{shortReference(value)}</span> }
