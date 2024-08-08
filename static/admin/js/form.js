// Toolbar functions for the rich text editor
function wrapText(tag) {
  const editor = document.getElementById("editor");
  const start = editor.selectionStart;
  const end = editor.selectionEnd;
  const selectedText = editor.value.substring(start, end);
  const wrappedText = `<${tag}>${selectedText}</${tag}>`;
  editor.setRangeText(wrappedText, start, end, "end");
}
function addLink() {
  const url = prompt("Enter the URL", "https://example.com");
  if (url) {
    wrapTextWithCustomTag(`<a href="${url}">`, "</a>");
  }
}
function addCode() {
  wrapTextWithCustomTag("<pre><code>", "</code></pre>");
}
function addTable() {
  const rows = prompt("Enter the number of rows", "2");
  const cols = prompt("Enter the number of columns", "2");
  if (rows && cols) {
    let table = '<table border="1">\n';
    for (let i = 0; i < rows; i++) {
      table += "  <tr>\n";
      for (let j = 0; j < cols; j++) {
        table += "    <td>&nbsp;</td>\n";
      }
      table += "  </tr>\n";
    }
    table += "</table>\n";
    insertTextAtCursor(table);
  }
}
function addBulletList() {
  wrapTextWithCustomTag("<ul><li>", "</li></ul>");
}
function addNumberedList() {
  wrapTextWithCustomTag("<ol><li>", "</li></ol>");
}

// Helper functions
function wrapTextWithCustomTag(openTag, closeTag) {
  const editor = document.getElementById("editor");
  const start = editor.selectionStart;
  const end = editor.selectionEnd;
  const selectedText = editor.value.substring(start, end);
  const wrappedText = `${openTag}${selectedText}${closeTag}`;
  editor.setRangeText(wrappedText, start, end, "end");
}
function insertTextAtCursor(text) {
  const editor = document.getElementById("editor");
  const start = editor.selectionStart;
  editor.setRangeText(text, start, start, "end");
}
