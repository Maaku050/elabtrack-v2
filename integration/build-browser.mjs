import { mkdir, writeFile } from 'node:fs/promises'
import path from 'node:path'
import { fileURLToPath } from 'node:url'

const root = path.resolve(path.dirname(fileURLToPath(import.meta.url)), '..')
const frontend = path.join(root, 'frontend')
const temp = path.join(frontend, 'phase1g.local')
await mkdir(temp, { recursive: true })
// Keep generated input/output out of production source and version control.
await writeFile(path.join(temp, 'driver.ts'), `import ${JSON.stringify(path.join(root, 'integration/browser-driver.ts'))}`)
await writeFile(path.join(temp, 'index.html'), '<!doctype html><html lang="en"><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width, initial-scale=1.0"><title>eLabTrack disposable integration</title></head><body><div id="root"></div><script type="module" src="./driver.ts"></script></body></html>')
const { build } = await import(path.join(frontend, 'node_modules/vite/dist/node/index.js'))
await build({ root: frontend, configFile: path.join(frontend, 'vite.config.ts'),
  define: { 'import.meta.env.VITE_API_URL': JSON.stringify('/api/v1') },
  build: { outDir: path.join(temp, 'dist'), emptyOutDir: true, rollupOptions: { input: path.join(temp, 'index.html') } },
})
