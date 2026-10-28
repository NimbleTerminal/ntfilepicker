import { BoobaTerminal, installActionBar, installRendererHud, parseRendererFromURL } from './booba/booba.js';

const layouts = { places: 'Places', columns: 'Columns', list: 'List' };
const params = new URLSearchParams(location.search);
const pathPicker = params.get('widget') === 'path';
const requested = params.get('layout');
const layout = Object.hasOwn(layouts, requested) ? requested : 'places';
const status = document.querySelector('#status');
document.querySelector('#title').textContent = pathPicker ? 'Complete a path · Path picker' : `Pick a file · ${layouts[layout]}`;
const restart = new URL(location.href);
if (pathPicker) {
  restart.searchParams.set('widget', 'path');
  restart.searchParams.delete('layout');
  document.querySelector('[data-widget="path"]').setAttribute('aria-current', 'page');
  document.querySelector('#help').textContent = 'Type a path to fuzzy-filter · Tab enters a folder · Shift+Tab goes up · Enter chooses · Esc cancels · Ctrl+U starts over. Try Documents/budget to see quoting for spaces. The files are pretend. Click the terminal to focus.';
} else {
  restart.searchParams.set('layout', layout);
  restart.searchParams.delete('widget');
}
document.querySelector('#restart').href = restart.href;
for (const link of document.querySelectorAll('[data-layout]')) {
  if (!pathPicker && link.dataset.layout === layout) link.setAttribute('aria-current', 'page');
}
function failed(error) {
  status.dataset.error = '';
  status.textContent = `Could not run the demo: ${error.message}. Try Restart.`;
  console.error(error);
}
async function main() {
  const go = new Go();
  go.argv = pathPicker ? ['demo', '-widget', 'path'] : ['demo', '-layout', layout];
  const response = await fetch('demos/demo.wasm');
  if (!response.ok) throw new Error(`Download failed (${response.status})`);
  const { instance } = await WebAssembly.instantiateStreaming(response, go.importObject);
  const running = go.run(instance);
  running.then(() => { status.textContent = 'Session ended. Choose Restart to run again.'; }, failed);
  const terminal = new BoobaTerminal('terminal', { renderer: parseRendererFromURL(), allowOSC7501Actions: true });
  await terminal.init();
  terminal.setActionKeyEncoder({
    canEncode: gesture => terminal.encodeKey(gesture) !== null,
    encode: gesture => terminal.encodeKey(gesture),
  });
  installActionBar(terminal.actions, {
    parent: document.querySelector('#controls'),
    openKeyboard: () => terminal.focus(),
    onLayout: () => terminal.fitAddon?.fit(),
  });
  terminal.onActionError = failed;
  window.addEventListener('pagehide', () => terminal.dispose(), { once: true });
  terminal.connectWasm();
  terminal.focus();
  installRendererHud(terminal.term, {
    parent: document.querySelector('.terminal-frame'),
    position: 'absolute',
    bindToggleHotkey: false,
  });
  status.textContent = 'Running locally in your browser. Nothing is read from your computer.';
}
main().catch(failed);
