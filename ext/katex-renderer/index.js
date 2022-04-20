const fs = require('fs');
const katex = require('katex');

let content;
try {
    content = fs.readFileSync(process.stdin.fd, 'utf-8');
} catch (e) {
    console.error('reading from stdin failed:', e);
    process.exit(1);
}

if (!content) {
    console.error('unexpected empty input');
    process.exit(1);
}

let inputEntries;
try {
    inputEntries = JSON.parse(content);
} catch (e) {
    console.error('failed to parse input:', e);
    process.exit(1);
}

if (!Array.isArray(inputEntries)) {
    console.error('input must be an array');
    process.exit(1);
}

const outputEntries = [];

for (const entry of inputEntries) {
    let htmlOut;
    try {
        htmlOut = katex.renderToString(entry.input, {
            displayMode: !!entry.display,
            throwOnError: true,
            strict: true,
            macros: entry.macros,
        });
    } catch (e) {
        console.error('failed to render input:', entry)
        console.error('error:', e);
        process.exit(1);
    }
    outputEntries.push(htmlOut);
}

console.log(JSON.stringify(outputEntries));
