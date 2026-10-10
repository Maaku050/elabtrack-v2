import { jsPDF } from 'jspdf'
import { autoTable } from 'jspdf-autotable'
import fontURL from '@/assets/fonts/DejaVuSans.ttf?url'
import { manilaTimeCompact } from '@/features/borrowing/lib/time'
import type { ReportFilter, ReportPage } from '../types'
let font: Promise<string> | undefined
async function loadFont() {
 font ??= fetch(fontURL).then(async r=>{if(!r.ok)throw new Error('The local PDF font could not be loaded.');const bytes=new Uint8Array(await r.arrayBuffer());let out='';for(let i=0;i<bytes.length;i+=8192)out+=String.fromCharCode(...bytes.subarray(i,i+8192));return btoa(out)}).catch(e=>{font=undefined;throw e})
 return font
}
/** Full authorized snapshot, never current-page rows; strings are PDF text, never HTML. */
export async function reportPDF(page: ReportPage,title:string,filter:ReportFilter):Promise<Blob>{
 if(page.total>5000||page.rows.length!==page.total)throw new Error('The complete bounded report snapshot is required. Narrow filters and retry.')
 const doc=new jsPDF({orientation:page.columns.length>6?'landscape':'portrait',unit:'mm',format:'a4',putOnlyUsedFonts:true,compress:true})
 doc.addFileToVFS('DejaVuSans.ttf',await loadFont());doc.addFont('DejaVuSans.ttf','FSMO','normal');doc.setFont('FSMO')
 const width=doc.internal.pageSize.getWidth(),height=doc.internal.pageSize.getHeight(),filters=Object.entries(filter).filter(([k,v])=>!['page','per_page'].includes(k)&&v!==undefined&&v!=='').map(([k,v])=>`${k.replaceAll('_',' ')}: ${v}`).join(' · ')||'All authorized records'
 const filterLines=doc.setFontSize(8).splitTextToSize(filters,width-28),top=35+filterLines.length*4
 const head=page.columns.map(v=>v.endsWith(' UTC')?v.replace(' UTC',' (Manila)'):v),rows=page.rows.map(row=>row.map((cell,i)=>page.columns[i]?.endsWith(' UTC')&&cell?manilaTimeCompact(cell):cell))
 autoTable(doc,{head:[head],body:rows,startY:top,margin:{left:14,right:14,top,bottom:15},showHead:'everyPage',rowPageBreak:'avoid',theme:'striped',styles:{font:'FSMO',fontStyle:'normal',fontSize:page.columns.length>12?6:8,cellPadding:2,overflow:'linebreak',textColor:[25,31,60]},headStyles:{font:'FSMO',fontStyle:'normal',fillColor:[43,16,187],textColor:[255,255,255]},alternateRowStyles:{fillColor:[245,246,252]},willDrawPage:()=>{doc.setFont('FSMO').setTextColor(25,31,60).setFontSize(10);doc.text('FSMO · eLabTrack V2',14,14);doc.setFontSize(14);doc.text(title,14,22);doc.setFontSize(8);doc.text(`Generated ${manilaTimeCompact(new Date().toISOString())} · Asia/Manila | Snapshot ${manilaTimeCompact(page.as_of)} | ${page.total} records`,14,29);doc.text(filterLines,14,35)}})
 for(let n=1;n<=doc.getNumberOfPages();n++){doc.setPage(n);doc.setFont('FSMO').setFontSize(8).setTextColor(80,85,108);doc.text(`Page ${n} of ${doc.getNumberOfPages()}`,width-14,height-7,{align:'right'})}
 if(!page.rows.length)doc.text('No matching records.',14,top+7)
 return doc.output('blob')
}
