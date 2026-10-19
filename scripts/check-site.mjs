// Validate the generated Pages artifact without requiring a local HTTP server.
import assert from 'node:assert/strict';
import { readFile, access } from 'node:fs/promises';
import { BoobaTerminal, installActionBar, installRendererHud, parseRendererFromURL } from '../web/booba/booba.js';
// Link the actual browser runtime exports: syntax checks alone miss missing APIs.
for (const api of [BoobaTerminal, installActionBar, installRendererHud, parseRendererFromURL]) {
  assert.equal(typeof api, 'function', 'required browser runtime API');
}
for (const file of [
  'index.html', 'site.css', 'site.js', 'wasm_exec.js', 'booba/booba.js',
  'ghostty-web/ghostty-web.js', 'ghostty-web/ghostty-vt.wasm',
  'LICENSE', 'THIRD_PARTY_NOTICES.md', 'demos/demo.wasm',
]) await access(`web/${file}`);
const bytes = await readFile('web/demos/demo.wasm');
assert.deepEqual([...bytes.subarray(0, 4)], [0, 97, 115, 109], 'WASM header');
assert.ok(WebAssembly.validate(bytes), 'valid WASM');
const html = await readFile('web/index.html', 'utf8');
const layouts = [...html.matchAll(/data-layout="([^"]+)"/g)].map(m => m[1]);
assert.deepEqual(layouts, ['places', 'columns', 'list'], 'layout cards');
const js = await readFile('web/site.js', 'utf8');
for (const layout of layouts) assert.ok(js.includes(`${layout}:`), `${layout}: site.js entry`);
for (const [, path] of html.matchAll(/(?:src|href)="([^"?#]+)"/g)) {
  if (/^(https?:|\.\/)/.test(path)) continue;
  await access(`web/${path}`);
}
console.log(`Site links, runtime assets, and ${layouts.length} layouts checked.`);
