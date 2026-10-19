// Exercise the real browser binary through the booba bridge under Node.
// This checks the chooser's behavior, not browser rendering.
import assert from 'node:assert/strict';
import { readFile } from 'node:fs/promises';
import vm from 'node:vm';
import { setTimeout as delay } from 'node:timers/promises';
vm.runInThisContext(await readFile('web/wasm_exec.js', 'utf8'));
const go = new Go();
go.argv = ['demo', '-layout', 'list'];
let exitCode;
go.exit = code => { exitCode = code; };
const { instance } = await WebAssembly.instantiate(await readFile('web/demos/demo.wasm'), go.importObject);
const timeout = setTimeout(() => { console.error('demo timed out'); process.exit(1); }, 20000);
const running = go.run(instance);
// Waits for text to appear in what the program writes from now on.
async function expectOutput(text) {
  let output = '';
  for (let n = 0; n < 160 && !output.includes(text); n++) {
    await delay(25);
    output += bubbletea_read();
  }
  assert.ok(output.includes(text), `expected terminal output: ${text}`);
  return output;
}
try {
  for (let n = 0; typeof globalThis.bubbletea_read !== 'function' && n < 100; n++) await delay(20);
  assert.equal(typeof globalThis.bubbletea_read, 'function');
  bubbletea_resize(100, 30);
  await expectOutput('Documents');
  // Negotiate the same program-status support as the browser action bar.
  bubbletea_write('\x1b]7501;?\x1b\\');
  const advertised = await expectOutput(':msg=');
  const encoded = advertised.match(/:msg=([A-Za-z0-9+/=]+)/)?.[1];
  assert.ok(encoded, 'action status payload');
  const description = Buffer.from(encoded, 'base64').toString('utf8');
  for (const label of ['Layout', 'Places', 'File types', 'Up a folder']) {
    assert.ok(description.includes(label), `advertised control: ${label}`);
  }
  bubbletea_write('Docu');
  await expectOutput('Docu');
  bubbletea_write('\r');
  await expectOutput('Recipes');
  bubbletea_write('2026');
  await expectOutput('budget');
  bubbletea_write('\r');
  await expectOutput('Picked');
  bubbletea_resize(30, 8);
  await expectOutput('Picked');
  bubbletea_write('\x03');
  await running;
  assert.equal(exitCode, 0);
  console.log('Demo: startup, actions, filter, open a folder, pick a file, resize, and quit passed');
} finally {
  clearTimeout(timeout);
}
process.exit(0);
