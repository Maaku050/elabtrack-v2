#!/usr/bin/env python3
"""Original vector illustrations; no network, external photos or real identities."""
from pathlib import Path
import json
root=Path(__file__).resolve().parent/'assets'
entries=[
('Stainless Steel Ladle','Utensils','ladle'),('Wooden Spoon','Utensils','spoon'),('Slotted Serving Spoon','Utensils','spoon'),('Silicone Spatula','Utensils','spatula'),('Stainless Turner','Utensils','spatula'),('Balloon Whisk','Utensils','whisk'),('Kitchen Tongs','Utensils','tongs'),('Soup Skimmer','Utensils','skimmer'),
('Small Saucepan','Cookware','pot'),('Medium Cooking Pot','Cookware','pot'),('Large Stockpot','Cookware','pot'),('Nonstick Frying Pan','Cookware','pan'),('Saute Pan','Cookware','pan'),('Grill Pan','Cookware','pan'),('Carbon Steel Wok','Cookware','wok'),('Steamer Basket','Cookware','strainer'),
('Measuring Cup Set','Measuring','cup'),('Measuring Spoon Set','Measuring','spoon'),('Digital Kitchen Scale','Measuring','scale'),('Food Thermometer','Measuring','thermometer'),('Glass Measuring Jug','Measuring','jug'),('Portion Scoop','Measuring','ladle'),
('Stainless Mixing Bowl','Preparation','bowl'),('Mesh Strainer','Preparation','strainer'),('Kitchen Colander','Preparation','colander'),('Food Prep Cutting Board','Preparation','board'),('Wooden Cutting Board','Preparation','board'),('Chef Knife','Preparation','knife'),('Bread Knife','Preparation','knife'),
('Stainless Food Pan','Serving','tray'),('Chafing Dish','Serving','chafing'),('Serving Tray','Serving','tray'),('Beverage Pitcher','Serving','jug'),('Large Serving Scoop','Serving','ladle'),
('Muffin Tray','Baking','muffin'),('Baking Sheet','Baking','tray'),('Rolling Pin','Baking','rolling'),('Pastry Brush','Baking','brush'),('Countertop Blender','Appliances','blender'),('Stand Mixer','Appliances','mixer')]
shapes={
'ladle':'<path d="M125 43v107q0 20-17 20" fill="none" stroke="{metal}" stroke-width="13"/><ellipse cx="96" cy="178" rx="35" ry="25" fill="{metal}"/><ellipse cx="96" cy="174" rx="25" ry="14" fill="#cbd5e1"/>',
'spoon':'<rect x="119" y="43" width="13" height="115" rx="6" fill="{metal}"/><ellipse cx="125" cy="178" rx="28" ry="38" fill="{metal}"/><ellipse cx="125" cy="179" rx="17" ry="25" fill="#d7dde6"/>',
'spatula':'<rect x="117" y="39" width="17" height="108" rx="6" fill="#293b59"/><rect x="96" y="137" width="59" height="73" rx="10" fill="{metal}"/><path d="M110 151v44m15-44v44m15-44v44" stroke="#d8e1ed" stroke-width="5"/>',
'whisk':'<rect x="115" y="34" width="20" height="70" rx="8" fill="#293b59"/><g fill="none" stroke="{metal}" stroke-width="4"><path d="M124 101C48 198 83 219 124 209C165 220 201 198 124 101Z"/><path d="M124 101C88 185 103 221 124 209C146 221 160 183 124 101Z"/><path d="M124 101v108"/></g>',
 'tongs':'<path d="M81 190L118 53q8-17 15 0l40 137" fill="none" stroke="{metal}" stroke-width="13" stroke-linecap="round"/><path d="M77 184l-8 23m106-23l10 23" stroke="#293b59" stroke-width="22" stroke-linecap="round"/>',
 'skimmer':'<rect x="121" y="35" width="10" height="110" rx="5" fill="{metal}"/><circle cx="126" cy="174" r="40" fill="{metal}"/><g fill="#e9edf3">'+''.join(f'<circle cx="{x}" cy="{y}" r="3"/>' for x in [107,126,145] for y in [156,174,191])+'</g>',
 'pot':'<path d="M67 102h122v67q0 27-26 27H93q-26 0-26-27Z" fill="{metal}"/><path d="M66 115H47v30h20m122-30h19v30h-19" stroke="#35465f" stroke-width="8" fill="none"/><ellipse cx="128" cy="102" rx="61" ry="17" fill="#d8e0ea"/><rect x="110" y="75" width="36" height="12" rx="6" fill="#35465f"/>',
 'pan':'<ellipse cx="105" cy="155" rx="65" ry="48" fill="#293747"/><ellipse cx="105" cy="149" rx="54" ry="35" fill="#4b5c71"/><path d="M156 124l57-61" stroke="#223148" stroke-width="20" stroke-linecap="round"/>',
 'wok':'<path d="M46 103q80 160 161 0Z" fill="#3b4a5e"/><ellipse cx="126" cy="103" rx="80" ry="29" fill="#697b8d"/><path d="M43 114L29 105m181 11l15-14" stroke="#293b59" stroke-width="13"/>',
 'strainer':'<path d="M63 97q62 144 125 0Z" fill="#c8d1dc" stroke="{metal}" stroke-width="4"/><ellipse cx="125" cy="97" rx="62" ry="20" fill="#e1e7ee" stroke="{metal}" stroke-width="5"/><path d="M181 93l46-36" stroke="#35465f" stroke-width="11"/><path d="M85 119l70 48m-62-23l43 28m-26-58l54 32" stroke="#9aa8b7" stroke-width="2"/>',
 'cup':'<path d="M73 80h107l-8 116H82Z" fill="#b9c9dc" stroke="{metal}" stroke-width="4"/><path d="M179 95q59 0 12 57h-15" stroke="{metal}" stroke-width="9" fill="none"/><path d="M95 106h35m-35 24h25m-25 24h35m-35 23h25" stroke="#546780" stroke-width="4"/>',
 'scale':'<rect x="55" y="97" width="145" height="99" rx="18" fill="{metal}"/><ellipse cx="127" cy="99" rx="70" ry="20" fill="#d9e1e9"/><rect x="97" y="150" width="60" height="25" rx="5" fill="#294259"/><text x="127" y="168" text-anchor="middle" fill="#d6efe6" font-family="monospace" font-size="13">0.00</text>',
 'thermometer':'<rect x="95" y="43" width="61" height="67" rx="12" fill="#293b59"/><rect x="106" y="58" width="40" height="24" rx="4" fill="#d2eadc"/><path d="M125 109v103" stroke="{metal}" stroke-width="8"/>',
 'jug':'<path d="M76 68h95l13 132H64Z" fill="#dce7f2" stroke="#8ea4ba" stroke-width="4"/><path d="M176 89q60 0 14 74h-9" stroke="#8ea4ba" stroke-width="8" fill="none"/><path d="M97 99h29m-29 27h20m-20 27h29m-29 27h20" stroke="#546780" stroke-width="4"/><path d="M73 71l-16-11 21 29" fill="#dce7f2"/>',
 'bowl':'<path d="M47 101q17 104 80 104q64 0 81-104Z" fill="{metal}"/><ellipse cx="127" cy="101" rx="80" ry="22" fill="#d9e2eb"/><path d="M72 135q20 55 56 56" stroke="#e7edf4" stroke-width="6" fill="none"/>',
 'colander':'<path d="M51 98q17 105 76 105q58 0 78-105Z" fill="{metal}"/><ellipse cx="128" cy="98" rx="76" ry="21" fill="#d7e1eb"/><g fill="#e1e8ef">'+''.join(f'<circle cx="{x}" cy="{y}" r="4"/>' for x in [91,116,141,166] for y in [125,148,170])+'</g>',
 'board':'<rect x="55" y="48" width="146" height="167" rx="16" fill="{metal}"/><rect x="88" y="61" width="79" height="10" rx="5" fill="#e9edf3"/><path d="M71 89v106m20-106v106m20-106v106m20-106v106m20-106v106m20-106v106m14-106v106" stroke="#b88653" stroke-width="2" opacity=".5"/>',
 'knife':'<path d="M111 107V36h28v85Z" fill="#28364b"/><path d="M108 115h32v44q0 34-43 62Z" fill="{metal}"/><path d="M108 119l-8 91" stroke="#e4eaf0" stroke-width="4"/>',
 'tray':'<rect x="35" y="74" width="186" height="123" rx="16" fill="{metal}"/><rect x="48" y="87" width="160" height="97" rx="10" fill="#d8e1ec"/><rect x="61" y="98" width="134" height="75" rx="8" fill="#a7b8c9"/>',
 'chafing':'<rect x="47" y="124" width="162" height="39" rx="6" fill="{metal}"/><path d="M50 123q75-108 155 0Z" fill="#d4dfe9" stroke="{metal}" stroke-width="4"/><rect x="109" y="63" width="38" height="9" rx="4" fill="#293b59"/><path d="M65 163v40m126-40v40" stroke="#293b59" stroke-width="8"/><ellipse cx="128" cy="189" rx="20" ry="9" fill="#b7c9db"/>',
 'muffin':'<rect x="36" y="56" width="184" height="154" rx="13" fill="{metal}"/>'+''.join(f'<ellipse cx="{x}" cy="{y}" rx="20" ry="18" fill="#536980" stroke="#d6e0ea" stroke-width="4"/>' for x in [72,127,183] for y in [93,143,185]),
 'rolling':'<path d="M53 143h151" stroke="#945f34" stroke-width="18" stroke-linecap="round"/><rect x="78" y="112" width="101" height="61" rx="19" fill="#c99460"/><path d="M88 121h76" stroke="#e8bf91" stroke-width="4"/>',
 'brush':'<rect x="116" y="39" width="23" height="95" rx="9" fill="#b4804c"/><path d="M94 140h65l-4 58H98Z" fill="{metal}"/><path d="M106 147v47m10-47v47m10-47v47m10-47v47m10-47v47" stroke="#e1c490" stroke-width="4"/>',
 'blender':'<path d="M82 45h91l-10 101H91Z" fill="#cdddeb" stroke="#839bb2" stroke-width="4"/><path d="M174 66q53 0 14 51h-20" stroke="#839bb2" stroke-width="9" fill="none"/><rect x="76" y="143" width="104" height="63" rx="15" fill="{metal}"/><circle cx="128" cy="169" r="11" fill="#e5edf3"/><path d="M99 111l51 14m-43 0l34-16" stroke="#687f97" stroke-width="5"/><rect x="77" y="36" width="101" height="12" rx="5" fill="#293b59"/>',
 'mixer':'<path d="M63 200h130v-22h-42V95h-47v93H63Z" fill="{metal}"/><path d="M76 62h72q42 0 42 39H76Z" fill="{metal}"/><path d="M158 107v38" stroke="#6c8199" stroke-width="8"/><path d="M131 133h72q-3 52-36 52q-33 0-36-52Z" fill="#c7d4e3"/>'}
manifest=[]
for i,(name,category,kind) in enumerate(entries):
 slug=name.lower().replace(' ','-');metal='#b7c6d6'
 if 'Wooden' in name or kind in ('rolling','brush'):metal='#c7955c'
 elif kind in ('spatula','mixer','blender'):metal=['#7786bd','#bf8893','#87aaa0'][i%3]
 art=shapes[kind].replace('{metal}',metal)
 svg=f'<svg xmlns="http://www.w3.org/2000/svg" width="256" height="256" viewBox="0 0 256 256"><title>{name} — original catalog illustration</title><rect width="256" height="256" rx="24" fill="#f1f4f9"/><ellipse cx="128" cy="222" rx="80" ry="9" fill="#d9e1ee"/>{art}</svg>'
 (root/'equipment'/f'{slug}.svg').write_text(svg+'\n');manifest.append({'name':name,'category':category,'description':f'{name} for food-service laboratory preparation and supervised FSMO use.','asset':f'equipment/{slug}.png','status':'INACTIVE' if i>=38 else 'ACTIVE','opening':0 if i==37 else 8+i%17})
for i in range(28):
 bg=['#ded9f8','#d5e8ef','#e7ddca','#d9e8df'][i%4];skin=['#d29b74','#bb825e','#e3ba92','#95613f'][i%4];shirt=['#5d63a9','#437c7c','#936575','#8a743f'][i%4];hair=['#303140','#4b362a','#6a4a32','#29282c'][i%4]
 svg=f'<svg xmlns="http://www.w3.org/2000/svg" width="256" height="256" viewBox="0 0 256 256"><title>Fictional illustrated presentation avatar {i+1:02d}</title><rect width="256" height="256" rx="128" fill="{bg}"/><path d="M51 256q0-84 77-84t77 84" fill="{shirt}"/><path d="M102 152v31q26 23 52 0v-31" fill="{skin}"/><ellipse cx="128" cy="112" rx="48" ry="60" fill="{skin}"/><path d="M79 112V82q0-51 51-51q53 0 50 68l-20-23q-33 22-64-4Z" fill="{hair}"/><path d="M101 108h12m30 0h12" stroke="{hair}" stroke-width="5" stroke-linecap="round"/><circle cx="109" cy="117" r="3" fill="#30303a"/><circle cx="149" cy="117" r="3" fill="#30303a"/><path d="M117 145q12 9 24 0" stroke="#8b4e42" stroke-width="4" fill="none" stroke-linecap="round"/><path d="M94 183l34 27 34-27" stroke="#ffffff" stroke-width="5" fill="none" opacity=".6"/></svg>'
 accessory=''
 if i%3==0:accessory='<g fill="none" stroke="#414963" stroke-width="4"><rect x="94" y="107" width="29" height="22" rx="6"/><rect x="134" y="107" width="29" height="22" rx="6"/><path d="M123 116h11"/></g>'
 elif i%3==1:accessory='<path d="M100 152q28 31 57-1l-8 18h-39Z" fill="'+hair+'" opacity=".8"/>'
 if i%7 in (2,4):accessory+='<circle cx="86" cy="143" r="5" fill="#bd985d"/><circle cx="170" cy="143" r="5" fill="#bd985d"/>'
 if i%7 in (3,5):accessory+='<path d="M76 118q-20 45 9 68l-8-80m102 0q24 59-4 77" fill="'+hair+'"/>'
 svg=svg.replace('</svg>',accessory+'</svg>')
 svg=svg.replace('q53 0 50 68',f'q{48+i%7} 0 50 {60+i%7*3}').replace('M117 145q12 9 24 0',f'M117 145q12 {5+i%7} 24 0')
 (root/'avatars'/f'avatar-{i+1:02d}.svg').write_text(svg+'\n')
(root/'equipment.json').write_text(json.dumps(manifest,indent=2)+'\n')
print('Generated40 original equipment vectors and28 illustrated fictional-avatar vectors. Render PNGs with the local Chromium asset tool.')
