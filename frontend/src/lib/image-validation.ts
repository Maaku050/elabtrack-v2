export async function validateImageFile(file: File, description='image') {
 if(file.size===0||file.size>512*1024)throw Error('Choose a nonempty image no larger than 512 KiB.')
 if(!['image/png','image/jpeg'].includes(file.type))throw Error(`Choose a PNG or JPEG ${description}.`)
 const bytes=new Uint8Array(await file.slice(0,8).arrayBuffer()),png=bytes.length>=8&&[137,80,78,71,13,10,26,10].every((n,i)=>bytes[i]===n),jpeg=bytes[0]===255&&bytes[1]===216&&bytes[2]===255
 if((file.type==='image/png'&&!png)||(file.type==='image/jpeg'&&!jpeg))throw Error('The file contents must match its PNG or JPEG format.')
 let bitmap:ImageBitmap
 try{bitmap=await createImageBitmap(file)}catch{throw Error('This image could not be decoded. Choose a valid PNG or JPEG.')}
 const valid=bitmap.width>0&&bitmap.height>0&&bitmap.width<=2048&&bitmap.height<=2048;bitmap.close()
 if(!valid)throw Error('Images must be at most 2048 pixels per dimension.')
}
