// static/scripts/marked.esm.js
function z() {
  return { async: false, breaks: false, extensions: null, gfm: true, hooks: null, pedantic: false, renderer: null, silent: false, tokenizer: null, walkTokens: null };
}
var T = z();
function N(l3) {
  T = l3;
}
var _ = { exec: () => null };
function E(l3) {
  let e = [];
  return (t) => {
    let n = Math.max(0, Math.min(3, t - 1)), s = e[n];
    return s || (s = l3(n), e[n] = s), s;
  };
}
function d(l3, e = "") {
  let t = typeof l3 == "string" ? l3 : l3.source, n = { replace: (s, r) => {
    let i = typeof r == "string" ? r : r.source;
    return i = i.replace(m.caret, "$1"), t = t.replace(s, i), n;
  }, getRegex: () => new RegExp(t, e) };
  return n;
}
var Te = ((l3 = "") => {
  try {
    return !!new RegExp("(?<=1)(?<!1)" + l3);
  } catch {
    return false;
  }
})();
var m = { codeRemoveIndent: /^(?: {1,4}| {0,3}\t)/gm, outputLinkReplace: /\\([\[\]])/g, indentCodeCompensation: /^(\s+)(?:```)/, beginningSpace: /^\s+/, endingHash: /#$/, startingSpaceChar: /^ /, endingSpaceChar: / $/, nonSpaceChar: /[^ ]/, newLineCharGlobal: /\n/g, tabCharGlobal: /\t/g, multipleSpaceGlobal: /\s+/g, blankLine: /^[ \t]*$/, doubleBlankLine: /\n[ \t]*\n[ \t]*$/, blockquoteStart: /^ {0,3}>/, blockquoteSetextReplace: /\n {0,3}((?:=+|-+) *)(?=\n|$)/g, blockquoteSetextReplace2: /^ {0,3}>[ \t]?/gm, listReplaceNesting: /^ {1,4}(?=( {4})*[^ ])/g, listIsTask: /^\[[ xX]\] +\S/, listReplaceTask: /^\[[ xX]\] +/, listTaskCheckbox: /\[[ xX]\]/, anyLine: /\n.*\n/, hrefBrackets: /^<(.*)>$/, tableDelimiter: /[:|]/, tableAlignChars: /^\||\| *$/g, tableRowBlankLine: /\n[ \t]*$/, tableAlignRight: /^ *-+: *$/, tableAlignCenter: /^ *:-+: *$/, tableAlignLeft: /^ *:-+ *$/, startATag: /^<a /i, endATag: /^<\/a>/i, startPreScriptTag: /^<(pre|code|kbd|script)(\s|>)/i, endPreScriptTag: /^<\/(pre|code|kbd|script)(\s|>)/i, startAngleBracket: /^</, endAngleBracket: />$/, pedanticHrefTitle: /^([^'"]*[^\s])\s+(['"])(.*)\2/, unicodeAlphaNumeric: /[\p{L}\p{N}]/u, escapeTest: /[&<>"']/, escapeReplace: /[&<>"']/g, escapeTestNoEncode: /[<>"']|&(?!(#\d{1,7}|#[Xx][a-fA-F0-9]{1,6}|\w+);)/, escapeReplaceNoEncode: /[<>"']|&(?!(#\d{1,7}|#[Xx][a-fA-F0-9]{1,6}|\w+);)/g, caret: /(^|[^\[])\^/g, percentDecode: /%25/g, findPipe: /\|/g, splitPipe: / \|/, slashPipe: /\\\|/g, carriageReturn: /\r\n|\r/g, spaceLine: /^ +$/gm, notSpaceStart: /^\S*/, endingNewline: /\n$/, listItemRegex: (l3) => new RegExp(`^( {0,3}${l3})((?:[	 ][^\\n]*)?(?:\\n|$))`), nextBulletRegex: E((l3) => new RegExp(`^ {0,${l3}}(?:[*+-]|\\d{1,9}[.)])((?:[ 	][^\\n]*)?(?:\\n|$))`)), hrRegex: E((l3) => new RegExp(`^ {0,${l3}}((?:- *){3,}|(?:_ *){3,}|(?:\\* *){3,})(?:\\n+|$)`)), fencesBeginRegex: E((l3) => new RegExp(`^ {0,${l3}}(?:\`\`\`|~~~)`)), headingBeginRegex: E((l3) => new RegExp(`^ {0,${l3}}#`)), htmlBeginRegex: E((l3) => new RegExp(`^ {0,${l3}}<(?:[a-z].*>|!--)`, "i")), blockquoteBeginRegex: E((l3) => new RegExp(`^ {0,${l3}}>`)) };
var Oe = /^(?:[ \t]*(?:\n|$))+/;
var we = /^((?: {4}| {0,3}\t)[^\n]+(?:\n(?:[ \t]*(?:\n|$))*)?)+/;
var ye = /^ {0,3}(`{3,}(?=[^`\n]*(?:\n|$))|~{3,})([^\n]*)(?:\n|$)(?:|([\s\S]*?)(?:\n|$))(?: {0,3}\1[~`]* *(?=\n|$)|$)/;
var B = /^ {0,3}((?:-[\t ]*){3,}|(?:_[ \t]*){3,}|(?:\*[ \t]*){3,})(?:\n+|$)/;
var Pe = /^ {0,3}(#{1,6})(?=\s|$)(.*)(?:\n+|$)/;
var j = / {0,3}(?:[*+-]|\d{1,9}[.)])/;
var oe = /^(?!bull |blockCode|fences|blockquote|heading|html|table)((?:.|\n(?!\s*?\n|bull |blockCode|fences|blockquote|heading|html|table))+?)\n {0,3}(=+|-+) *(?:\n+|$)/;
var ae = d(oe).replace(/bull/g, j).replace(/blockCode/g, /(?: {4}| {0,3}\t)/).replace(/fences/g, / {0,3}(?:`{3,}|~{3,})/).replace(/blockquote/g, / {0,3}>/).replace(/heading/g, / {0,3}#{1,6}(?:\s|$)/).replace(/html/g, / {0,3}<[^\n>]+>\n/).replace(/\|table/g, "").getRegex();
var Se = d(oe).replace(/bull/g, j).replace(/blockCode/g, /(?: {4}| {0,3}\t)/).replace(/fences/g, / {0,3}(?:`{3,}|~{3,})/).replace(/blockquote/g, / {0,3}>/).replace(/heading/g, / {0,3}#{1,6}(?:\s|$)/).replace(/html/g, / {0,3}<[^\n>]+>\n/).replace(/table/g, / {0,3}\|?(?:[:\- ]*\|)+[\:\- ]*\n/).getRegex();
var F = /^([^\n]+(?:\n(?!hr|heading|lheading|blockquote|fences|list|html|table|[ \t]+\n)[^\n]+)*)/;
var $e = /^[^\n]+/;
var U = /(?!\s*\])(?:\\[\s\S]|[^\[\]\\])+/;
var Le = d(/^ {0,3}\[(label)\]: *(?:\n[ \t]*)?([^<\s][^\s]*|<.*?>)(?:(?: +(?:\n[ \t]*)?| *\n[ \t]*)(title))? *(?:\n+|$)/).replace("label", U).replace("title", /(?:"(?:\\"?|[^"\\])*"|'[^'\n]*(?:\n[^'\n]+)*\n?'|\([^()]*\))/).getRegex();
var _e = d(/^(bull)([ \t][^\n]*?)?(?:\n|$)/).replace(/bull/g, j).getRegex();
var H = "address|article|aside|base|basefont|blockquote|body|caption|center|col|colgroup|dd|details|dialog|dir|div|dl|dt|fieldset|figcaption|figure|footer|form|frame|frameset|h[1-6]|head|header|hr|html|iframe|legend|li|link|main|menu|menuitem|meta|nav|noframes|ol|optgroup|option|p|param|search|section|summary|table|tbody|td|tfoot|th|thead|title|tr|track|ul";
var K = /<!--(?:-?>|[\s\S]*?(?:-->|$))/;
var Me = d("^ {0,3}(?:<(script|pre|style|textarea)[\\s>][\\s\\S]*?(?:</\\1>[^\\n]*\\n*|$)|comment[^\\n]*(\\n+|$)|<\\?[\\s\\S]*?(?:\\?>[^\\n]*\\n*|$)|<![A-Z][\\s\\S]*?(?:>[^\\n]*\\n*|$)|<!\\[CDATA\\[[\\s\\S]*?(?:\\]\\]>[^\\n]*\\n*|$)|</?(tag)(?: +|\\n|/?>)[\\s\\S]*?(?:(?:\\n[ 	]*)+\\n|$)|<(?!script|pre|style|textarea)([a-z][\\w-]*)(?:attribute)*? */?>(?=[ \\t]*(?:\\n|$))[\\s\\S]*?(?:(?:\\n[ 	]*)+\\n|$)|</(?!script|pre|style|textarea)[a-z][\\w-]*\\s*>(?=[ \\t]*(?:\\n|$))[\\s\\S]*?(?:(?:\\n[ 	]*)+\\n|$))", "i").replace("comment", K).replace("tag", H).replace("attribute", / +[a-zA-Z:_][\w.:-]*(?: *= *"[^"\n]*"| *= *'[^'\n]*'| *= *[^\s"'=<>`]+)?/).getRegex();
var le = (l3) => d(F).replace("hr", B).replace("heading", " {0,3}#{1,6}(?:\\s|$)").replace("|lheading", "").replace("|table", "").replace("blockquote", " {0,3}>").replace("fences", " {0,3}(?:`{3,}(?=[^`\\n]*\\n)|~~~)[^\\n]*\\n").replace("list", l3).replace("html", "</?(?:tag)(?: +|\\n|/?>)|<(?:script|pre|style|textarea|!--)").replace("tag", H).getRegex();
var ze = le(/ {0,3}(?:[*+-]|1[.)])[ \t]+[^ \t\n]/);
var Ee = le(/ {0,3}(?:[*+-]|\d{1,9}[.)])(?:[ \t]|\n|$)/);
var Ce = d(/^( {0,3}> ?(paragraph|[^\n]*)(?:\n|$))+/).replace("paragraph", Ee).getRegex();
var W = { blockquote: Ce, code: we, def: Le, fences: ye, heading: Pe, hr: B, html: Me, lheading: ae, list: _e, newline: Oe, paragraph: ze, table: _, text: $e };
var se = d("^ *([^\\n ].*)\\n {0,3}((?:\\| *)?:?-+:? *(?:\\| *:?-+:? *)*(?:\\| *)?)(?:\\n((?:(?! *\\n|hr|heading|blockquote|code|fences|list|html).*(?:\\n|$))*)\\n*|$)").replace("hr", B).replace("heading", " {0,3}#{1,6}(?:\\s|$)").replace("blockquote", " {0,3}>").replace("code", "(?: {4}| {0,3}	)[^\\n]").replace("fences", " {0,3}(?:`{3,}(?=[^`\\n]*\\n)|~~~)[^\\n]*\\n").replace("list", " {0,3}(?:[*+-]|1[.)])[ \\t]").replace("html", "</?(?:tag)(?: +|\\n|/?>)|<(?:script|pre|style|textarea|!--)").replace("tag", H).getRegex();
var Ae = { ...W, lheading: Se, table: se, paragraph: d(F).replace("hr", B).replace("heading", " {0,3}#{1,6}(?:\\s|$)").replace("|lheading", "").replace("table", se).replace("blockquote", " {0,3}>").replace("fences", " {0,3}(?:`{3,}(?=[^`\\n]*\\n)|~~~)[^\\n]*\\n").replace("list", " {0,3}(?:[*+-]|1[.)])[ \\t]+[^ \\t\\n]").replace("html", "</?(?:tag)(?: +|\\n|/?>)|<(?:script|pre|style|textarea|!--)").replace("tag", H).getRegex() };
var Ie = { ...W, html: d(`^ *(?:comment *(?:\\n|\\s*$)|<(tag)[\\s\\S]+?</\\1> *(?:\\n{2,}|\\s*$)|<tag(?:"[^"]*"|'[^']*'|\\s[^'"/>\\s]*)*?/?> *(?:\\n{2,}|\\s*$))`).replace("comment", K).replace(/tag/g, "(?!(?:a|em|strong|small|s|cite|q|dfn|abbr|data|time|code|var|samp|kbd|sub|sup|i|b|u|mark|ruby|rt|rp|bdi|bdo|span|br|wbr|ins|del|img)\\b)\\w+(?!:|[^\\w\\s@]*@)\\b").getRegex(), def: /^ *\[([^\]]+)\]: *<?([^\s>]+)>?(?: +(["(][^\n]+[")]))? *(?:\n+|$)/, heading: /^(#{1,6})(.*)(?:\n+|$)/, fences: _, lheading: /^(.+?)\n {0,3}(=+|-+) *(?:\n+|$)/, paragraph: d(F).replace("hr", B).replace("heading", ` *#{1,6} *[^
]`).replace("lheading", ae).replace("|table", "").replace("blockquote", " {0,3}>").replace("|fences", "").replace("|list", "").replace("|html", "").replace("|tag", "").getRegex() };
var Be = /^\\([!"#$%&'()*+,\-./:;<=>?@\[\]\\^_`{|}~])/;
var qe = /^(`+)([^`]|[^`][\s\S]*?[^`])\1(?!`)/;
var ue = /^( {2,}|\\)\n(?!\s*$)/;
var De = /^(`+|[^`])(?:(?= {2,}\n)|[\s\S]*?(?:(?=[\\<!\[`*_]|\b_|$)|[^ ](?= {2,}\n)))/;
var C = /[\p{P}\p{S}]/u;
var Z = /[\s\p{P}\p{S}]/u;
var X = /[^\s\p{P}\p{S}]/u;
var ve = d(/^((?![*_])punctSpace)/, "u").replace(/punctSpace/g, Z).getRegex();
var pe = /(?!~)[\p{P}\p{S}]/u;
var He = /(?!~)[\s\p{P}\p{S}]/u;
var Ze = /(?:[^\s\p{P}\p{S}]|~)/u;
var Ge = d(/link|precode-code|html/, "g").replace("link", /\[(?:[^\[\]`]|(?<a>`+)[^`]+\k<a>(?!`))*?\]\((?:\\[\s\S]|[^\\\(\)]|\((?:\\[\s\S]|[^\\\(\)])*\))*\)/).replace("precode-", Te ? "(?<!`)()" : "(^^|[^`])").replace("code", /(?<b>`+)[^`]+\k<b>(?!`)/).replace("html", /<(?! )[^<>]*?>/).getRegex();
var ce = /^(?:\*+(?:((?!\*)punct)|([^\s*]))?)|^_+(?:((?!_)punct)|([^\s_]))?/;
var Ne = d(ce, "u").replace(/punct/g, C).getRegex();
var Qe = d(ce, "u").replace(/punct/g, pe).getRegex();
var he = "^[^_*]*?__[^_*]*?\\*[^_*]*?(?=__)|[^*]+(?=[^*])|(?!\\*)punct(\\*+)(?=[\\s]|$)|notPunctSpace(\\*+)(?!\\*)(?=punctSpace|$)|(?!\\*)punctSpace(\\*+)(?=notPunctSpace)|[\\s](\\*+)(?!\\*)(?=punct)|(?!\\*)punct(\\*+)(?!\\*)(?=punct)|notPunctSpace(\\*+)(?=notPunctSpace)";
var je = d(he, "gu").replace(/notPunctSpace/g, X).replace(/punctSpace/g, Z).replace(/punct/g, C).getRegex();
var Fe = d(he, "gu").replace(/notPunctSpace/g, Ze).replace(/punctSpace/g, He).replace(/punct/g, pe).getRegex();
var Ue = d("^[^_*]*?\\*\\*[^_*]*?_[^_*]*?(?=\\*\\*)|[^_]+(?=[^_])|(?!_)punct(_+)(?=[\\s]|$)|notPunctSpace(_+)(?!_)(?=punctSpace|$)|(?!_)punctSpace(_+)(?=notPunctSpace)|[\\s](_+)(?!_)(?=punct)|(?!_)punct(_+)(?!_)(?=punct)", "gu").replace(/notPunctSpace/g, X).replace(/punctSpace/g, Z).replace(/punct/g, C).getRegex();
var Ke = d(/^~~?(?:((?!~)punct)|[^\s~])/, "u").replace(/punct/g, C).getRegex();
var We = "^[^~]+(?=[^~])|(?!~)punct(~~?)(?=[\\s]|$)|notPunctSpace(~~?)(?!~)(?=punctSpace|$)|(?!~)punctSpace(~~?)(?=notPunctSpace)|[\\s](~~?)(?!~)(?=punct)|(?!~)punct(~~?)(?!~)(?=punct)|notPunctSpace(~~?)(?=notPunctSpace)";
var Xe = d(We, "gu").replace(/notPunctSpace/g, X).replace(/punctSpace/g, Z).replace(/punct/g, C).getRegex();
var Je = d(/\\(punct)/, "gu").replace(/punct/g, C).getRegex();
var Ve = d(/^<(scheme:[^\s\x00-\x1f<>]*|email)>/).replace("scheme", /[a-zA-Z][a-zA-Z0-9+.-]{1,31}/).replace("email", /[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+(@)[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)+(?![-_])/).getRegex();
var Ye = d(K).replace("(?:-->|$)", "-->").getRegex();
var et = d("^comment|^</[a-zA-Z][\\w:-]*\\s*>|^<[a-zA-Z][\\w-]*(?:attribute)*?\\s*/?>|^<\\?[\\s\\S]*?\\?>|^<![a-zA-Z]+\\s[\\s\\S]*?>|^<!\\[CDATA\\[[\\s\\S]*?\\]\\]>").replace("comment", Ye).replace("attribute", /\s+[a-zA-Z:_][\w.:-]*(?:\s*=\s*"[^"]*"|\s*=\s*'[^']*'|\s*=\s*[^\s"'=<>`]+)?/).getRegex();
var v = /(?:\[(?:\\[\s\S]|[^\[\]\\])*\]|\\[\s\S]|`+(?!`)[^`]*?`+(?!`)|``+(?=\])|[^\[\]\\`])*?/;
var tt = d(/^!?\[(label)\]\(\s*(href)(?:(?:[ \t]+(?:\n[ \t]*)?|\n[ \t]*)(title))?\s*\)/).replace("label", v).replace("href", /<(?:\\.|[^\n<>\\])+>|[^ \t\n\x00-\x1f]+|(?=\))/).replace("title", /"(?:\\"?|[^"\\])*"|'(?:\\'?|[^'\\])*'|\((?:\\\)?|[^)\\])*\)/).getRegex();
var ke = d(/^!?\[(label)\]\[(ref)\]/).replace("label", v).replace("ref", U).getRegex();
var de = d(/^!?\[(ref)\](?:\[\])?/).replace("ref", U).getRegex();
var nt = d("reflink|nolink(?!\\()", "g").replace("reflink", ke).replace("nolink", de).getRegex();
var ie = /[hH][tT][tT][pP][sS]?|[fF][tT][pP]/;
var J = { _backpedal: _, anyPunctuation: Je, autolink: Ve, blockSkip: Ge, br: ue, code: qe, del: _, delLDelim: _, delRDelim: _, emStrongLDelim: Ne, emStrongRDelimAst: je, emStrongRDelimUnd: Ue, escape: Be, link: tt, nolink: de, punctuation: ve, reflink: ke, reflinkSearch: nt, tag: et, text: De, url: _ };
var rt = { ...J, link: d(/^!?\[(label)\]\((.*?)\)/).replace("label", v).getRegex(), reflink: d(/^!?\[(label)\]\s*\[([^\]]*)\]/).replace("label", v).getRegex() };
var Q = { ...J, emStrongRDelimAst: Fe, emStrongLDelim: Qe, delLDelim: Ke, delRDelim: Xe, url: d(/^((?:protocol):\/\/|www\.)(?:[a-zA-Z0-9\-]+\.?)+[^\s<]*|^email/).replace("protocol", ie).replace("email", /[A-Za-z0-9._+-]+(@)[a-zA-Z0-9-_]+(?:\.[a-zA-Z0-9-_]*[a-zA-Z0-9])+(?![-_])/).getRegex(), _backpedal: /(?:[^?!.,:;*_'"~()&]+|\([^)]*\)|&(?![a-zA-Z0-9]+;$)|[?!.,:;*_'"~)]+(?!$))+/, del: /^(~~?)(?=[^\s~])((?:\\[\s\S]|[^\\])*?(?:\\[\s\S]|[^\s~\\]))\1(?=[^~]|$)/, text: d(/^(`+|~+|[^`~])(?:(?=[`~])|(?= {2,}\n)|(?=[a-zA-Z0-9.!#$%&'*+\/=?_`{\|}~-]+@)|[\s\S]*?(?:(?=[\\<!\[`*~_]|\b_|protocol:\/\/|www\.|$)|[^ ](?= {2,}\n)|[^a-zA-Z0-9.!#$%&'*+\/=?_`{\|}~-](?=[a-zA-Z0-9.!#$%&'*+\/=?_`{\|}~-]+@)))/).replace("protocol", ie).getRegex() };
var st = { ...Q, br: d(ue).replace("{2,}", "*").getRegex(), text: d(Q.text).replace("\\b_", "\\b_| {2,}\\n").replace(/\{2,\}/g, "*").getRegex() };
var q = { normal: W, gfm: Ae, pedantic: Ie };
var A = { normal: J, gfm: Q, breaks: st, pedantic: rt };
var it = { "&": "&amp;", "<": "&lt;", ">": "&gt;", '"': "&quot;", "'": "&#39;" };
var ge = (l3) => it[l3];
function O(l3, e) {
  if (e) {
    if (m.escapeTest.test(l3)) return l3.replace(m.escapeReplace, ge);
  } else if (m.escapeTestNoEncode.test(l3)) return l3.replace(m.escapeReplaceNoEncode, ge);
  return l3;
}
function V(l3) {
  try {
    l3 = encodeURI(l3).replace(m.percentDecode, "%");
  } catch {
    return null;
  }
  return l3;
}
function Y(l3, e) {
  let t = l3.replace(m.findPipe, (r, i, o) => {
    let u = false, a = i;
    for (; --a >= 0 && o[a] === "\\"; ) u = !u;
    return u ? "|" : " |";
  }), n = t.split(m.splitPipe), s = 0;
  if (n[0].trim() || n.shift(), n.length > 0 && !n.at(-1)?.trim() && n.pop(), e) if (n.length > e) n.splice(e);
  else for (; n.length < e; ) n.push("");
  for (; s < n.length; s++) n[s] = n[s].trim().replace(m.slashPipe, "|");
  return n;
}
function $(l3, e, t) {
  let n = l3.length;
  if (n === 0) return "";
  let s = 0;
  for (; s < n; ) {
    let r = l3.charAt(n - s - 1);
    if (r === e && !t) s++;
    else if (r !== e && t) s++;
    else break;
  }
  return l3.slice(0, n - s);
}
function ee(l3) {
  let e = l3.split(`
`), t = e.length - 1;
  for (; t >= 0 && m.blankLine.test(e[t]); ) t--;
  return e.length - t <= 2 ? l3 : e.slice(0, t + 1).join(`
`);
}
function fe(l3, e) {
  if (l3.indexOf(e[1]) === -1) return -1;
  let t = 0;
  for (let n = 0; n < l3.length; n++) if (l3[n] === "\\") n++;
  else if (l3[n] === e[0]) t++;
  else if (l3[n] === e[1] && (t--, t < 0)) return n;
  return t > 0 ? -2 : -1;
}
function me(l3, e = 0) {
  let t = e, n = "";
  for (let s of l3) if (s === "	") {
    let r = 4 - t % 4;
    n += " ".repeat(r), t += r;
  } else n += s, t++;
  return n;
}
function xe(l3, e, t, n, s) {
  let r = e.href, i = e.title || null, o = l3[1].replace(s.other.outputLinkReplace, "$1");
  n.state.inLink = true;
  let u = { type: l3[0].charAt(0) === "!" ? "image" : "link", raw: t, href: r, title: i, text: o, tokens: n.inlineTokens(o) };
  return n.state.inLink = false, u;
}
function ot(l3, e, t) {
  let n = l3.match(t.other.indentCodeCompensation);
  if (n === null) return e;
  let s = n[1];
  return e.split(`
`).map((r) => {
    let i = r.match(t.other.beginningSpace);
    if (i === null) return r;
    let [o] = i;
    return o.length >= s.length ? r.slice(s.length) : r;
  }).join(`
`);
}
var w = class {
  options;
  rules;
  lexer;
  constructor(e) {
    this.options = e || T;
  }
  space(e) {
    let t = this.rules.block.newline.exec(e);
    if (t && t[0].length > 0) return { type: "space", raw: t[0] };
  }
  code(e) {
    let t = this.rules.block.code.exec(e);
    if (t) {
      let n = this.options.pedantic ? t[0] : ee(t[0]), s = n.replace(this.rules.other.codeRemoveIndent, "");
      return { type: "code", raw: n, codeBlockStyle: "indented", text: s };
    }
  }
  fences(e) {
    let t = this.rules.block.fences.exec(e);
    if (t) {
      let n = t[0], s = ot(n, t[3] || "", this.rules);
      return { type: "code", raw: n, lang: t[2] ? t[2].trim().replace(this.rules.inline.anyPunctuation, "$1") : t[2], text: s };
    }
  }
  heading(e) {
    let t = this.rules.block.heading.exec(e);
    if (t) {
      let n = t[2].trim();
      if (this.rules.other.endingHash.test(n)) {
        let s = $(n, "#");
        (this.options.pedantic || !s || this.rules.other.endingSpaceChar.test(s)) && (n = s.trim());
      }
      return { type: "heading", raw: $(t[0], `
`), depth: t[1].length, text: n, tokens: this.lexer.inline(n) };
    }
  }
  hr(e) {
    let t = this.rules.block.hr.exec(e);
    if (t) return { type: "hr", raw: $(t[0], `
`) };
  }
  blockquote(e) {
    let t = this.rules.block.blockquote.exec(e);
    if (t) {
      let n = $(t[0], `
`).split(`
`), s = "", r = "", i = [];
      for (; n.length > 0; ) {
        let o = false, u = [], a;
        for (a = 0; a < n.length; a++) if (this.rules.other.blockquoteStart.test(n[a])) u.push(n[a]), o = true;
        else if (!o) u.push(n[a]);
        else break;
        n = n.slice(a);
        let p = u.join(`
`), c = p.replace(this.rules.other.blockquoteSetextReplace, `
    $1`).replace(this.rules.other.blockquoteSetextReplace2, "");
        s = s ? `${s}
${p}` : p, r = r ? `${r}
${c}` : c;
        let h = this.lexer.state.top;
        if (this.lexer.state.top = true, this.lexer.blockTokens(c, i, true), this.lexer.state.top = h, n.length === 0) break;
        let k = i.at(-1);
        if (k?.type === "code") break;
        if (k?.type === "blockquote") {
          let R = k, f = R.raw + `
` + n.join(`
`), S = this.blockquote(f);
          i[i.length - 1] = S, s = s.substring(0, s.length - R.raw.length) + S.raw, r = r.substring(0, r.length - R.text.length) + S.text;
          break;
        } else if (k?.type === "list") {
          let R = k, f = R.raw + `
` + n.join(`
`), S = this.list(f);
          i[i.length - 1] = S, s = s.substring(0, s.length - k.raw.length) + S.raw, r = r.substring(0, r.length - R.raw.length) + S.raw, n = f.substring(i.at(-1).raw.length).split(`
`);
          continue;
        }
      }
      return { type: "blockquote", raw: s, tokens: i, text: r };
    }
  }
  list(e) {
    let t = this.rules.block.list.exec(e);
    if (t) {
      let n = t[1].trim(), s = n.length > 1, r = { type: "list", raw: "", ordered: s, start: s ? +n.slice(0, -1) : "", loose: false, items: [] };
      n = s ? `\\d{1,9}\\${n.slice(-1)}` : `\\${n}`, this.options.pedantic && (n = s ? n : "[*+-]");
      let i = this.rules.other.listItemRegex(n), o = false;
      for (; e; ) {
        let a = false, p = "", c = "";
        if (!(t = i.exec(e)) || this.rules.block.hr.test(e)) break;
        p = t[0], e = e.substring(p.length);
        let h = me(t[2].split(`
`, 1)[0], t[1].length), k = e.split(`
`, 1)[0], R = !h.trim(), f = 0;
        if (this.options.pedantic ? (f = 2, c = h.trimStart()) : R ? f = t[1].length + 1 : (f = h.search(this.rules.other.nonSpaceChar), f = f > 4 ? 1 : f, c = h.slice(f), f += t[1].length), R && this.rules.other.blankLine.test(k) && (p += k + `
`, e = e.substring(k.length + 1), a = true), !a) {
          let S = this.rules.other.nextBulletRegex(f), te = this.rules.other.hrRegex(f), ne = this.rules.other.fencesBeginRegex(f), re = this.rules.other.headingBeginRegex(f), be = this.rules.other.htmlBeginRegex(f), Re = this.rules.other.blockquoteBeginRegex(f);
          for (; e; ) {
            let G = e.split(`
`, 1)[0], I;
            if (k = G, this.options.pedantic ? (k = k.replace(this.rules.other.listReplaceNesting, "  "), I = k) : I = k.replace(this.rules.other.tabCharGlobal, "    "), ne.test(k) || re.test(k) || be.test(k) || Re.test(k) || S.test(k) || te.test(k)) break;
            if (I.search(this.rules.other.nonSpaceChar) >= f || !k.trim()) c += `
` + I.slice(f);
            else {
              if (R || h.replace(this.rules.other.tabCharGlobal, "    ").search(this.rules.other.nonSpaceChar) >= 4 || ne.test(h) || re.test(h) || te.test(h)) break;
              c += `
` + k;
            }
            R = !k.trim(), p += G + `
`, e = e.substring(G.length + 1), h = I.slice(f);
          }
        }
        r.loose || (o ? r.loose = true : this.rules.other.doubleBlankLine.test(p) && (o = true)), r.items.push({ type: "list_item", raw: p, task: !!this.options.gfm && this.rules.other.listIsTask.test(c), loose: false, text: c, tokens: [] }), r.raw += p;
      }
      let u = r.items.at(-1);
      if (u) u.raw = u.raw.trimEnd(), u.text = u.text.trimEnd();
      else return;
      r.raw = r.raw.trimEnd();
      for (let a of r.items) {
        this.lexer.state.top = false, a.tokens = this.lexer.blockTokens(a.text, []);
        let p = a.tokens[0];
        if (a.task && (p?.type === "text" || p?.type === "paragraph")) {
          a.text = a.text.replace(this.rules.other.listReplaceTask, ""), p.raw = p.raw.replace(this.rules.other.listReplaceTask, ""), p.text = p.text.replace(this.rules.other.listReplaceTask, "");
          for (let h = this.lexer.inlineQueue.length - 1; h >= 0; h--) if (this.rules.other.listIsTask.test(this.lexer.inlineQueue[h].src)) {
            this.lexer.inlineQueue[h].src = this.lexer.inlineQueue[h].src.replace(this.rules.other.listReplaceTask, "");
            break;
          }
          let c = this.rules.other.listTaskCheckbox.exec(a.raw);
          if (c) {
            let h = { type: "checkbox", raw: c[0] + " ", checked: c[0] !== "[ ]" };
            a.checked = h.checked, r.loose ? a.tokens[0] && ["paragraph", "text"].includes(a.tokens[0].type) && "tokens" in a.tokens[0] && a.tokens[0].tokens ? (a.tokens[0].raw = h.raw + a.tokens[0].raw, a.tokens[0].text = h.raw + a.tokens[0].text, a.tokens[0].tokens.unshift(h)) : a.tokens.unshift({ type: "paragraph", raw: h.raw, text: h.raw, tokens: [h] }) : a.tokens.unshift(h);
          }
        } else a.task && (a.task = false);
        if (!r.loose) {
          let c = a.tokens.filter((k) => k.type === "space"), h = c.length > 0 && c.some((k) => this.rules.other.anyLine.test(k.raw));
          r.loose = h;
        }
      }
      if (r.loose) for (let a of r.items) {
        a.loose = true;
        for (let p of a.tokens) p.type === "text" && (p.type = "paragraph");
      }
      return r;
    }
  }
  html(e) {
    let t = this.rules.block.html.exec(e);
    if (t) {
      let n = ee(t[0]);
      return { type: "html", block: true, raw: n, pre: t[1] === "pre" || t[1] === "script" || t[1] === "style", text: n };
    }
  }
  def(e) {
    let t = this.rules.block.def.exec(e);
    if (t) {
      let n = t[1].toLowerCase().replace(this.rules.other.multipleSpaceGlobal, " "), s = t[2] ? t[2].replace(this.rules.other.hrefBrackets, "$1").replace(this.rules.inline.anyPunctuation, "$1") : "", r = t[3] ? t[3].substring(1, t[3].length - 1).replace(this.rules.inline.anyPunctuation, "$1") : t[3];
      return { type: "def", tag: n, raw: $(t[0], `
`), href: s, title: r };
    }
  }
  table(e) {
    let t = this.rules.block.table.exec(e);
    if (!t || !this.rules.other.tableDelimiter.test(t[2])) return;
    let n = Y(t[1]), s = t[2].replace(this.rules.other.tableAlignChars, "").split("|"), r = t[3]?.trim() ? t[3].replace(this.rules.other.tableRowBlankLine, "").split(`
`) : [], i = { type: "table", raw: $(t[0], `
`), header: [], align: [], rows: [] };
    if (n.length === s.length) {
      for (let o of s) this.rules.other.tableAlignRight.test(o) ? i.align.push("right") : this.rules.other.tableAlignCenter.test(o) ? i.align.push("center") : this.rules.other.tableAlignLeft.test(o) ? i.align.push("left") : i.align.push(null);
      for (let o = 0; o < n.length; o++) i.header.push({ text: n[o], tokens: this.lexer.inline(n[o]), header: true, align: i.align[o] });
      for (let o of r) i.rows.push(Y(o, i.header.length).map((u, a) => ({ text: u, tokens: this.lexer.inline(u), header: false, align: i.align[a] })));
      return i;
    }
  }
  lheading(e) {
    let t = this.rules.block.lheading.exec(e);
    if (t) {
      let n = t[1].trim();
      return { type: "heading", raw: $(t[0], `
`), depth: t[2].charAt(0) === "=" ? 1 : 2, text: n, tokens: this.lexer.inline(n) };
    }
  }
  paragraph(e) {
    let t = this.rules.block.paragraph.exec(e);
    if (t) {
      let n = t[1].charAt(t[1].length - 1) === `
` ? t[1].slice(0, -1) : t[1];
      return { type: "paragraph", raw: t[0], text: n, tokens: this.lexer.inline(n) };
    }
  }
  text(e) {
    let t = this.rules.block.text.exec(e);
    if (t) return { type: "text", raw: t[0], text: t[0], tokens: this.lexer.inline(t[0]) };
  }
  escape(e) {
    let t = this.rules.inline.escape.exec(e);
    if (t) return { type: "escape", raw: t[0], text: t[1] };
  }
  tag(e) {
    let t = this.rules.inline.tag.exec(e);
    if (t) return !this.lexer.state.inLink && this.rules.other.startATag.test(t[0]) ? this.lexer.state.inLink = true : this.lexer.state.inLink && this.rules.other.endATag.test(t[0]) && (this.lexer.state.inLink = false), !this.lexer.state.inRawBlock && this.rules.other.startPreScriptTag.test(t[0]) ? this.lexer.state.inRawBlock = true : this.lexer.state.inRawBlock && this.rules.other.endPreScriptTag.test(t[0]) && (this.lexer.state.inRawBlock = false), { type: "html", raw: t[0], inLink: this.lexer.state.inLink, inRawBlock: this.lexer.state.inRawBlock, block: false, text: t[0] };
  }
  link(e) {
    let t = this.rules.inline.link.exec(e);
    if (t) {
      let n = t[2].trim();
      if (!this.options.pedantic && this.rules.other.startAngleBracket.test(n)) {
        if (!this.rules.other.endAngleBracket.test(n)) return;
        let i = $(n.slice(0, -1), "\\");
        if ((n.length - i.length) % 2 === 0) return;
      } else {
        let i = fe(t[2], "()");
        if (i === -2) return;
        if (i > -1) {
          let u = (t[0].indexOf("!") === 0 ? 5 : 4) + t[1].length + i;
          t[2] = t[2].substring(0, i), t[0] = t[0].substring(0, u).trim(), t[3] = "";
        }
      }
      let s = t[2], r = "";
      if (this.options.pedantic) {
        let i = this.rules.other.pedanticHrefTitle.exec(s);
        i && (s = i[1], r = i[3]);
      } else r = t[3] ? t[3].slice(1, -1) : "";
      return s = s.trim(), this.rules.other.startAngleBracket.test(s) && (this.options.pedantic && !this.rules.other.endAngleBracket.test(n) ? s = s.slice(1) : s = s.slice(1, -1)), xe(t, { href: s && s.replace(this.rules.inline.anyPunctuation, "$1"), title: r && r.replace(this.rules.inline.anyPunctuation, "$1") }, t[0], this.lexer, this.rules);
    }
  }
  reflink(e, t) {
    let n;
    if ((n = this.rules.inline.reflink.exec(e)) || (n = this.rules.inline.nolink.exec(e))) {
      let s = (n[2] || n[1]).replace(this.rules.other.multipleSpaceGlobal, " "), r = t[s.toLowerCase()];
      if (!r) {
        let i = n[0].charAt(0);
        return { type: "text", raw: i, text: i };
      }
      return xe(n, r, n[0], this.lexer, this.rules);
    }
  }
  emStrong(e, t, n = "") {
    let s = this.rules.inline.emStrongLDelim.exec(e);
    if (!s || !s[1] && !s[2] && !s[3] && !s[4] || s[4] && n.match(this.rules.other.unicodeAlphaNumeric)) return;
    if (!(s[1] || s[3] || "") || !n || this.rules.inline.punctuation.exec(n)) {
      let i = [...s[0]].length - 1, o, u, a = i, p = 0, c = s[0][0] === "*" ? this.rules.inline.emStrongRDelimAst : this.rules.inline.emStrongRDelimUnd;
      for (c.lastIndex = 0, t = t.slice(-1 * e.length + i); (s = c.exec(t)) !== null; ) {
        if (o = s[1] || s[2] || s[3] || s[4] || s[5] || s[6], !o) continue;
        if (u = [...o].length, s[3] || s[4]) {
          a += u;
          continue;
        } else if ((s[5] || s[6]) && i % 3 && !((i + u) % 3)) {
          p += u;
          continue;
        }
        if (a -= u, a > 0) continue;
        u = Math.min(u, u + a + p);
        let h = [...s[0]][0].length, k = e.slice(0, i + s.index + h + u);
        if (Math.min(i, u) % 2) {
          let f = k.slice(1, -1);
          return { type: "em", raw: k, text: f, tokens: this.lexer.inlineTokens(f) };
        }
        let R = k.slice(2, -2);
        return { type: "strong", raw: k, text: R, tokens: this.lexer.inlineTokens(R) };
      }
    }
  }
  codespan(e) {
    let t = this.rules.inline.code.exec(e);
    if (t) {
      let n = t[2].replace(this.rules.other.newLineCharGlobal, " "), s = this.rules.other.nonSpaceChar.test(n), r = this.rules.other.startingSpaceChar.test(n) && this.rules.other.endingSpaceChar.test(n);
      return s && r && (n = n.substring(1, n.length - 1)), { type: "codespan", raw: t[0], text: n };
    }
  }
  br(e) {
    let t = this.rules.inline.br.exec(e);
    if (t) return { type: "br", raw: t[0] };
  }
  del(e, t, n = "") {
    let s = this.rules.inline.delLDelim.exec(e);
    if (!s) return;
    if (!(s[1] || "") || !n || this.rules.inline.punctuation.exec(n)) {
      let i = [...s[0]].length - 1, o, u, a = i, p = this.rules.inline.delRDelim;
      for (p.lastIndex = 0, t = t.slice(-1 * e.length + i); (s = p.exec(t)) !== null; ) {
        if (o = s[1] || s[2] || s[3] || s[4] || s[5] || s[6], !o || (u = [...o].length, u !== i)) continue;
        if (s[3] || s[4]) {
          a += u;
          continue;
        }
        if (a -= u, a > 0) continue;
        u = Math.min(u, u + a);
        let c = [...s[0]][0].length, h = e.slice(0, i + s.index + c + u), k = h.slice(i, -i);
        return { type: "del", raw: h, text: k, tokens: this.lexer.inlineTokens(k) };
      }
    }
  }
  autolink(e) {
    let t = this.rules.inline.autolink.exec(e);
    if (t) {
      let n, s;
      return t[2] === "@" ? (n = t[1], s = "mailto:" + n) : (n = t[1], s = n), { type: "link", raw: t[0], text: n, href: s, tokens: [{ type: "text", raw: n, text: n }] };
    }
  }
  url(e) {
    let t;
    if (t = this.rules.inline.url.exec(e)) {
      let n, s;
      if (t[2] === "@") n = t[0], s = "mailto:" + n;
      else {
        let r;
        do
          r = t[0], t[0] = this.rules.inline._backpedal.exec(t[0])?.[0] ?? "";
        while (r !== t[0]);
        n = t[0], t[1] === "www." ? s = "http://" + t[0] : s = t[0];
      }
      return { type: "link", raw: t[0], text: n, href: s, tokens: [{ type: "text", raw: n, text: n }] };
    }
  }
  inlineText(e) {
    let t = this.rules.inline.text.exec(e);
    if (t) {
      let n = this.lexer.state.inRawBlock;
      return { type: "text", raw: t[0], text: t[0], escaped: n };
    }
  }
};
var x = class l {
  tokens;
  options;
  state;
  inlineQueue;
  tokenizer;
  constructor(e) {
    this.tokens = [], this.tokens.links = /* @__PURE__ */ Object.create(null), this.options = e || T, this.options.tokenizer = this.options.tokenizer || new w(), this.tokenizer = this.options.tokenizer, this.tokenizer.options = this.options, this.tokenizer.lexer = this, this.inlineQueue = [], this.state = { inLink: false, inRawBlock: false, top: true };
    let t = { other: m, block: q.normal, inline: A.normal };
    this.options.pedantic ? (t.block = q.pedantic, t.inline = A.pedantic) : this.options.gfm && (t.block = q.gfm, this.options.breaks ? t.inline = A.breaks : t.inline = A.gfm), this.tokenizer.rules = t;
  }
  static get rules() {
    return { block: q, inline: A };
  }
  static lex(e, t) {
    return new l(t).lex(e);
  }
  static lexInline(e, t) {
    return new l(t).inlineTokens(e);
  }
  lex(e) {
    e = e.replace(m.carriageReturn, `
`), this.blockTokens(e, this.tokens);
    for (let t = 0; t < this.inlineQueue.length; t++) {
      let n = this.inlineQueue[t];
      this.inlineTokens(n.src, n.tokens);
    }
    return this.inlineQueue = [], this.tokens;
  }
  blockTokens(e, t = [], n = false) {
    this.tokenizer.lexer = this, this.options.pedantic && (e = e.replace(m.tabCharGlobal, "    ").replace(m.spaceLine, ""));
    let s = 1 / 0;
    for (; e; ) {
      if (e.length < s) s = e.length;
      else {
        this.infiniteLoopError(e.charCodeAt(0));
        break;
      }
      let r;
      if (this.options.extensions?.block?.some((o) => (r = o.call({ lexer: this }, e, t)) ? (e = e.substring(r.raw.length), t.push(r), true) : false)) continue;
      if (r = this.tokenizer.space(e)) {
        e = e.substring(r.raw.length);
        let o = t.at(-1);
        r.raw.length === 1 && o !== void 0 ? o.raw += `
` : t.push(r);
        continue;
      }
      if (r = this.tokenizer.code(e)) {
        e = e.substring(r.raw.length);
        let o = t.at(-1);
        o?.type === "paragraph" || o?.type === "text" ? (o.raw += (o.raw.endsWith(`
`) ? "" : `
`) + r.raw, o.text += `
` + r.text, this.inlineQueue.at(-1).src = o.text) : t.push(r);
        continue;
      }
      if (r = this.tokenizer.fences(e)) {
        e = e.substring(r.raw.length), t.push(r);
        continue;
      }
      if (r = this.tokenizer.heading(e)) {
        e = e.substring(r.raw.length), t.push(r);
        continue;
      }
      if (r = this.tokenizer.hr(e)) {
        e = e.substring(r.raw.length), t.push(r);
        continue;
      }
      if (r = this.tokenizer.blockquote(e)) {
        e = e.substring(r.raw.length), t.push(r);
        continue;
      }
      if (r = this.tokenizer.list(e)) {
        e = e.substring(r.raw.length), t.push(r);
        continue;
      }
      if (r = this.tokenizer.html(e)) {
        e = e.substring(r.raw.length), t.push(r);
        continue;
      }
      if (r = this.tokenizer.def(e)) {
        e = e.substring(r.raw.length);
        let o = t.at(-1);
        o?.type === "paragraph" || o?.type === "text" ? (o.raw += (o.raw.endsWith(`
`) ? "" : `
`) + r.raw, o.text += `
` + r.raw, this.inlineQueue.at(-1).src = o.text) : this.tokens.links[r.tag] || (this.tokens.links[r.tag] = { href: r.href, title: r.title }, t.push(r));
        continue;
      }
      if (r = this.tokenizer.table(e)) {
        e = e.substring(r.raw.length), t.push(r);
        continue;
      }
      if (r = this.tokenizer.lheading(e)) {
        e = e.substring(r.raw.length), t.push(r);
        continue;
      }
      let i = e;
      if (this.options.extensions?.startBlock) {
        let o = 1 / 0, u = e.slice(1), a;
        this.options.extensions.startBlock.forEach((p) => {
          a = p.call({ lexer: this }, u), typeof a == "number" && a >= 0 && (o = Math.min(o, a));
        }), o < 1 / 0 && o >= 0 && (i = e.substring(0, o + 1));
      }
      if (this.state.top && (r = this.tokenizer.paragraph(i))) {
        let o = t.at(-1);
        n && o?.type === "paragraph" ? (o.raw += (o.raw.endsWith(`
`) ? "" : `
`) + r.raw, o.text += `
` + r.text, this.inlineQueue.pop(), this.inlineQueue.at(-1).src = o.text) : t.push(r), n = i.length !== e.length, e = e.substring(r.raw.length);
        continue;
      }
      if (r = this.tokenizer.text(e)) {
        e = e.substring(r.raw.length);
        let o = t.at(-1);
        o?.type === "text" ? (o.raw += (o.raw.endsWith(`
`) ? "" : `
`) + r.raw, o.text += `
` + r.text, this.inlineQueue.pop(), this.inlineQueue.at(-1).src = o.text) : t.push(r);
        continue;
      }
      if (e) {
        this.infiniteLoopError(e.charCodeAt(0));
        break;
      }
    }
    return this.state.top = true, t;
  }
  inline(e, t = []) {
    return this.inlineQueue.push({ src: e, tokens: t }), t;
  }
  inlineTokens(e, t = []) {
    this.tokenizer.lexer = this;
    let n = e;
    if (this.tokens.links) {
      let o = Object.keys(this.tokens.links);
      o.length > 0 && (n = n.replace(this.tokenizer.rules.inline.reflinkSearch, (u) => o.includes(u.slice(u.lastIndexOf("[") + 1, -1)) ? "[" + "a".repeat(u.length - 2) + "]" : u));
    }
    n = n.replace(this.tokenizer.rules.inline.anyPunctuation, "++"), n = n.replace(this.tokenizer.rules.inline.blockSkip, (o, u, a) => {
      let p = a ? a.length : 0;
      return o.slice(0, p) + "[" + "a".repeat(o.length - p - 2) + "]";
    }), n = this.options.hooks?.emStrongMask?.call({ lexer: this }, n) ?? n;
    let s = false, r = "", i = 1 / 0;
    for (; e; ) {
      if (e.length < i) i = e.length;
      else {
        this.infiniteLoopError(e.charCodeAt(0));
        break;
      }
      s || (r = ""), s = false;
      let o;
      if (this.options.extensions?.inline?.some((a) => (o = a.call({ lexer: this }, e, t)) ? (e = e.substring(o.raw.length), t.push(o), true) : false)) continue;
      if (o = this.tokenizer.escape(e)) {
        e = e.substring(o.raw.length), t.push(o);
        continue;
      }
      if (o = this.tokenizer.tag(e)) {
        e = e.substring(o.raw.length), t.push(o);
        continue;
      }
      if (o = this.tokenizer.link(e)) {
        e = e.substring(o.raw.length), t.push(o);
        continue;
      }
      if (o = this.tokenizer.reflink(e, this.tokens.links)) {
        e = e.substring(o.raw.length);
        let a = t.at(-1);
        o.type === "text" && a?.type === "text" ? (a.raw += o.raw, a.text += o.text) : t.push(o);
        continue;
      }
      if (o = this.tokenizer.emStrong(e, n, r)) {
        e = e.substring(o.raw.length), t.push(o);
        continue;
      }
      if (o = this.tokenizer.codespan(e)) {
        e = e.substring(o.raw.length), t.push(o);
        continue;
      }
      if (o = this.tokenizer.br(e)) {
        e = e.substring(o.raw.length), t.push(o);
        continue;
      }
      if (o = this.tokenizer.del(e, n, r)) {
        e = e.substring(o.raw.length), t.push(o);
        continue;
      }
      if (o = this.tokenizer.autolink(e)) {
        e = e.substring(o.raw.length), t.push(o);
        continue;
      }
      if (!this.state.inLink && (o = this.tokenizer.url(e))) {
        e = e.substring(o.raw.length), t.push(o);
        continue;
      }
      let u = e;
      if (this.options.extensions?.startInline) {
        let a = 1 / 0, p = e.slice(1), c;
        this.options.extensions.startInline.forEach((h) => {
          c = h.call({ lexer: this }, p), typeof c == "number" && c >= 0 && (a = Math.min(a, c));
        }), a < 1 / 0 && a >= 0 && (u = e.substring(0, a + 1));
      }
      if (o = this.tokenizer.inlineText(u)) {
        e = e.substring(o.raw.length), o.raw.slice(-1) !== "_" && (r = o.raw.slice(-1)), s = true;
        let a = t.at(-1);
        a?.type === "text" ? (a.raw += o.raw, a.text += o.text) : t.push(o);
        continue;
      }
      if (e) {
        this.infiniteLoopError(e.charCodeAt(0));
        break;
      }
    }
    return t;
  }
  infiniteLoopError(e) {
    let t = "Infinite loop on byte: " + e;
    if (this.options.silent) console.error(t);
    else throw new Error(t);
  }
};
var y = class {
  options;
  parser;
  constructor(e) {
    this.options = e || T;
  }
  space(e) {
    return "";
  }
  code({ text: e, lang: t, escaped: n }) {
    let s = (t || "").match(m.notSpaceStart)?.[0], r = e.replace(m.endingNewline, "") + `
`;
    return s ? '<pre><code class="language-' + O(s) + '">' + (n ? r : O(r, true)) + `</code></pre>
` : "<pre><code>" + (n ? r : O(r, true)) + `</code></pre>
`;
  }
  blockquote({ tokens: e }) {
    return `<blockquote>
${this.parser.parse(e)}</blockquote>
`;
  }
  html({ text: e }) {
    return e;
  }
  def(e) {
    return "";
  }
  heading({ tokens: e, depth: t }) {
    return `<h${t}>${this.parser.parseInline(e)}</h${t}>
`;
  }
  hr(e) {
    return `<hr>
`;
  }
  list(e) {
    let t = e.ordered, n = e.start, s = "";
    for (let o = 0; o < e.items.length; o++) {
      let u = e.items[o];
      s += this.listitem(u);
    }
    let r = t ? "ol" : "ul", i = t && n !== 1 ? ' start="' + n + '"' : "";
    return "<" + r + i + `>
` + s + "</" + r + `>
`;
  }
  listitem(e) {
    return `<li>${this.parser.parse(e.tokens)}</li>
`;
  }
  checkbox({ checked: e }) {
    return "<input " + (e ? 'checked="" ' : "") + 'disabled="" type="checkbox"> ';
  }
  paragraph({ tokens: e }) {
    return `<p>${this.parser.parseInline(e)}</p>
`;
  }
  table(e) {
    let t = "", n = "";
    for (let r = 0; r < e.header.length; r++) n += this.tablecell(e.header[r]);
    t += this.tablerow({ text: n });
    let s = "";
    for (let r = 0; r < e.rows.length; r++) {
      let i = e.rows[r];
      n = "";
      for (let o = 0; o < i.length; o++) n += this.tablecell(i[o]);
      s += this.tablerow({ text: n });
    }
    return s && (s = `<tbody>${s}</tbody>`), `<table>
<thead>
` + t + `</thead>
` + s + `</table>
`;
  }
  tablerow({ text: e }) {
    return `<tr>
${e}</tr>
`;
  }
  tablecell(e) {
    let t = this.parser.parseInline(e.tokens), n = e.header ? "th" : "td";
    return (e.align ? `<${n} align="${e.align}">` : `<${n}>`) + t + `</${n}>
`;
  }
  strong({ tokens: e }) {
    return `<strong>${this.parser.parseInline(e)}</strong>`;
  }
  em({ tokens: e }) {
    return `<em>${this.parser.parseInline(e)}</em>`;
  }
  codespan({ text: e }) {
    return `<code>${O(e, true)}</code>`;
  }
  br(e) {
    return "<br>";
  }
  del({ tokens: e }) {
    return `<del>${this.parser.parseInline(e)}</del>`;
  }
  link({ href: e, title: t, tokens: n }) {
    let s = this.parser.parseInline(n), r = V(e);
    if (r === null) return s;
    e = r;
    let i = '<a href="' + e + '"';
    return t && (i += ' title="' + O(t) + '"'), i += ">" + s + "</a>", i;
  }
  image({ href: e, title: t, text: n, tokens: s }) {
    s && (n = this.parser.parseInline(s, this.parser.textRenderer));
    let r = V(e);
    if (r === null) return O(n);
    e = r;
    let i = `<img src="${e}" alt="${O(n)}"`;
    return t && (i += ` title="${O(t)}"`), i += ">", i;
  }
  text(e) {
    return "tokens" in e && e.tokens ? this.parser.parseInline(e.tokens) : "escaped" in e && e.escaped ? e.text : O(e.text);
  }
};
var L = class {
  strong({ text: e }) {
    return e;
  }
  em({ text: e }) {
    return e;
  }
  codespan({ text: e }) {
    return e;
  }
  del({ text: e }) {
    return e;
  }
  html({ text: e }) {
    return e;
  }
  text({ text: e }) {
    return e;
  }
  link({ text: e }) {
    return "" + e;
  }
  image({ text: e }) {
    return "" + e;
  }
  br() {
    return "";
  }
  checkbox({ raw: e }) {
    return e;
  }
};
var b = class l2 {
  options;
  renderer;
  textRenderer;
  constructor(e) {
    this.options = e || T, this.options.renderer = this.options.renderer || new y(), this.renderer = this.options.renderer, this.renderer.options = this.options, this.renderer.parser = this, this.textRenderer = new L();
  }
  static parse(e, t) {
    return new l2(t).parse(e);
  }
  static parseInline(e, t) {
    return new l2(t).parseInline(e);
  }
  parse(e) {
    this.renderer.parser = this;
    let t = "";
    for (let n = 0; n < e.length; n++) {
      let s = e[n];
      if (this.options.extensions?.renderers?.[s.type]) {
        let i = s, o = this.options.extensions.renderers[i.type].call({ parser: this }, i);
        if (o !== false || !["space", "hr", "heading", "code", "table", "blockquote", "list", "html", "def", "paragraph", "text"].includes(i.type)) {
          t += o || "";
          continue;
        }
      }
      let r = s;
      switch (r.type) {
        case "space": {
          t += this.renderer.space(r);
          break;
        }
        case "hr": {
          t += this.renderer.hr(r);
          break;
        }
        case "heading": {
          t += this.renderer.heading(r);
          break;
        }
        case "code": {
          t += this.renderer.code(r);
          break;
        }
        case "table": {
          t += this.renderer.table(r);
          break;
        }
        case "blockquote": {
          t += this.renderer.blockquote(r);
          break;
        }
        case "list": {
          t += this.renderer.list(r);
          break;
        }
        case "checkbox": {
          t += this.renderer.checkbox(r);
          break;
        }
        case "html": {
          t += this.renderer.html(r);
          break;
        }
        case "def": {
          t += this.renderer.def(r);
          break;
        }
        case "paragraph": {
          t += this.renderer.paragraph(r);
          break;
        }
        case "text": {
          t += this.renderer.text(r);
          break;
        }
        default: {
          let i = 'Token with "' + r.type + '" type was not found.';
          if (this.options.silent) return console.error(i), "";
          throw new Error(i);
        }
      }
    }
    return t;
  }
  parseInline(e, t = this.renderer) {
    this.renderer.parser = this;
    let n = "";
    for (let s = 0; s < e.length; s++) {
      let r = e[s];
      if (this.options.extensions?.renderers?.[r.type]) {
        let o = this.options.extensions.renderers[r.type].call({ parser: this }, r);
        if (o !== false || !["escape", "html", "link", "image", "strong", "em", "codespan", "br", "del", "text"].includes(r.type)) {
          n += o || "";
          continue;
        }
      }
      let i = r;
      switch (i.type) {
        case "escape": {
          n += t.text(i);
          break;
        }
        case "html": {
          n += t.html(i);
          break;
        }
        case "link": {
          n += t.link(i);
          break;
        }
        case "image": {
          n += t.image(i);
          break;
        }
        case "checkbox": {
          n += t.checkbox(i);
          break;
        }
        case "strong": {
          n += t.strong(i);
          break;
        }
        case "em": {
          n += t.em(i);
          break;
        }
        case "codespan": {
          n += t.codespan(i);
          break;
        }
        case "br": {
          n += t.br(i);
          break;
        }
        case "del": {
          n += t.del(i);
          break;
        }
        case "text": {
          n += t.text(i);
          break;
        }
        default: {
          let o = 'Token with "' + i.type + '" type was not found.';
          if (this.options.silent) return console.error(o), "";
          throw new Error(o);
        }
      }
    }
    return n;
  }
};
var P = class {
  options;
  block;
  constructor(e) {
    this.options = e || T;
  }
  static passThroughHooks = /* @__PURE__ */ new Set(["preprocess", "postprocess", "processAllTokens", "emStrongMask"]);
  static passThroughHooksRespectAsync = /* @__PURE__ */ new Set(["preprocess", "postprocess", "processAllTokens"]);
  preprocess(e) {
    return e;
  }
  postprocess(e) {
    return e;
  }
  processAllTokens(e) {
    return e;
  }
  emStrongMask(e) {
    return e;
  }
  provideLexer(e = this.block) {
    return e ? x.lex : x.lexInline;
  }
  provideParser(e = this.block) {
    return e ? b.parse : b.parseInline;
  }
};
var D = class {
  defaults = z();
  options = this.setOptions;
  parse = this.parseMarkdown(true);
  parseInline = this.parseMarkdown(false);
  Parser = b;
  Renderer = y;
  TextRenderer = L;
  Lexer = x;
  Tokenizer = w;
  Hooks = P;
  constructor(...e) {
    this.use(...e);
  }
  walkTokens(e, t) {
    let n = [];
    for (let s of e) switch (n = n.concat(t.call(this, s)), s.type) {
      case "table": {
        let r = s;
        for (let i of r.header) n = n.concat(this.walkTokens(i.tokens, t));
        for (let i of r.rows) for (let o of i) n = n.concat(this.walkTokens(o.tokens, t));
        break;
      }
      case "list": {
        let r = s;
        n = n.concat(this.walkTokens(r.items, t));
        break;
      }
      default: {
        let r = s;
        this.defaults.extensions?.childTokens?.[r.type] ? this.defaults.extensions.childTokens[r.type].forEach((i) => {
          let o = r[i].flat(1 / 0);
          n = n.concat(this.walkTokens(o, t));
        }) : r.tokens && (n = n.concat(this.walkTokens(r.tokens, t)));
      }
    }
    return n;
  }
  use(...e) {
    let t = this.defaults.extensions || { renderers: {}, childTokens: {} };
    return e.forEach((n) => {
      let s = { ...n };
      if (s.async = this.defaults.async || s.async || false, n.extensions && (n.extensions.forEach((r) => {
        if (!r.name) throw new Error("extension name required");
        if ("renderer" in r) {
          let i = t.renderers[r.name];
          i ? t.renderers[r.name] = function(...o) {
            let u = r.renderer.apply(this, o);
            return u === false && (u = i.apply(this, o)), u;
          } : t.renderers[r.name] = r.renderer;
        }
        if ("tokenizer" in r) {
          if (!r.level || r.level !== "block" && r.level !== "inline") throw new Error("extension level must be 'block' or 'inline'");
          let i = t[r.level];
          i ? i.unshift(r.tokenizer) : t[r.level] = [r.tokenizer], r.start && (r.level === "block" ? t.startBlock ? t.startBlock.push(r.start) : t.startBlock = [r.start] : r.level === "inline" && (t.startInline ? t.startInline.push(r.start) : t.startInline = [r.start]));
        }
        "childTokens" in r && r.childTokens && (t.childTokens[r.name] = r.childTokens);
      }), s.extensions = t), n.renderer) {
        let r = this.defaults.renderer || new y(this.defaults);
        for (let i in n.renderer) {
          if (!(i in r)) throw new Error(`renderer '${i}' does not exist`);
          if (["options", "parser"].includes(i)) continue;
          let o = i, u = n.renderer[o], a = r[o];
          r[o] = (...p) => {
            let c = u.apply(r, p);
            return c === false && (c = a.apply(r, p)), c || "";
          };
        }
        s.renderer = r;
      }
      if (n.tokenizer) {
        let r = this.defaults.tokenizer || new w(this.defaults);
        for (let i in n.tokenizer) {
          if (!(i in r)) throw new Error(`tokenizer '${i}' does not exist`);
          if (["options", "rules", "lexer"].includes(i)) continue;
          let o = i, u = n.tokenizer[o], a = r[o];
          r[o] = (...p) => {
            let c = u.apply(r, p);
            return c === false && (c = a.apply(r, p)), c;
          };
        }
        s.tokenizer = r;
      }
      if (n.hooks) {
        let r = this.defaults.hooks || new P();
        for (let i in n.hooks) {
          if (!(i in r)) throw new Error(`hook '${i}' does not exist`);
          if (["options", "block"].includes(i)) continue;
          let o = i, u = n.hooks[o], a = r[o];
          P.passThroughHooks.has(i) ? r[o] = (p) => {
            if (this.defaults.async && P.passThroughHooksRespectAsync.has(i)) return (async () => {
              let h = await u.call(r, p);
              return a.call(r, h);
            })();
            let c = u.call(r, p);
            return a.call(r, c);
          } : r[o] = (...p) => {
            if (this.defaults.async) return (async () => {
              let h = await u.apply(r, p);
              return h === false && (h = await a.apply(r, p)), h;
            })();
            let c = u.apply(r, p);
            return c === false && (c = a.apply(r, p)), c;
          };
        }
        s.hooks = r;
      }
      if (n.walkTokens) {
        let r = this.defaults.walkTokens, i = n.walkTokens;
        s.walkTokens = function(o) {
          let u = [];
          return u.push(i.call(this, o)), r && (u = u.concat(r.call(this, o))), u;
        };
      }
      this.defaults = { ...this.defaults, ...s };
    }), this;
  }
  setOptions(e) {
    return this.defaults = { ...this.defaults, ...e }, this;
  }
  lexer(e, t) {
    return x.lex(e, t ?? this.defaults);
  }
  parser(e, t) {
    return b.parse(e, t ?? this.defaults);
  }
  parseMarkdown(e) {
    return (n, s) => {
      let r = { ...s }, i = { ...this.defaults, ...r }, o = this.onError(!!i.silent, !!i.async);
      if (this.defaults.async === true && r.async === false) return o(new Error("marked(): The async option was set to true by an extension. Remove async: false from the parse options object to return a Promise."));
      if (typeof n > "u" || n === null) return o(new Error("marked(): input parameter is undefined or null"));
      if (typeof n != "string") return o(new Error("marked(): input parameter is of type " + Object.prototype.toString.call(n) + ", string expected"));
      if (i.hooks && (i.hooks.options = i, i.hooks.block = e), i.async) return (async () => {
        let u = i.hooks ? await i.hooks.preprocess(n) : n, p = await (i.hooks ? await i.hooks.provideLexer(e) : e ? x.lex : x.lexInline)(u, i), c = i.hooks ? await i.hooks.processAllTokens(p) : p;
        i.walkTokens && await Promise.all(this.walkTokens(c, i.walkTokens));
        let k = await (i.hooks ? await i.hooks.provideParser(e) : e ? b.parse : b.parseInline)(c, i);
        return i.hooks ? await i.hooks.postprocess(k) : k;
      })().catch(o);
      try {
        i.hooks && (n = i.hooks.preprocess(n));
        let a = (i.hooks ? i.hooks.provideLexer(e) : e ? x.lex : x.lexInline)(n, i);
        i.hooks && (a = i.hooks.processAllTokens(a)), i.walkTokens && this.walkTokens(a, i.walkTokens);
        let c = (i.hooks ? i.hooks.provideParser(e) : e ? b.parse : b.parseInline)(a, i);
        return i.hooks && (c = i.hooks.postprocess(c)), c;
      } catch (u) {
        return o(u);
      }
    };
  }
  onError(e, t) {
    return (n) => {
      if (n.message += `
Please report this to https://github.com/markedjs/marked.`, e) {
        let s = "<p>An error occurred:</p><pre>" + O(n.message + "", true) + "</pre>";
        return t ? Promise.resolve(s) : s;
      }
      if (t) return Promise.reject(n);
      throw n;
    };
  }
};
var M = new D();
function g(l3, e) {
  return M.parse(l3, e);
}
g.options = g.setOptions = function(l3) {
  return M.setOptions(l3), g.defaults = M.defaults, N(g.defaults), g;
};
g.getDefaults = z;
g.defaults = T;
g.use = function(...l3) {
  return M.use(...l3), g.defaults = M.defaults, N(g.defaults), g;
};
g.walkTokens = function(l3, e) {
  return M.walkTokens(l3, e);
};
g.parseInline = M.parseInline;
g.Parser = b;
g.parser = b.parse;
g.Renderer = y;
g.TextRenderer = L;
g.Lexer = x;
g.lexer = x.lex;
g.Tokenizer = w;
g.Hooks = P;
g.parse = g;
var Kt = g.options;
var Wt = g.setOptions;
var Xt = g.use;
var Jt = g.walkTokens;
var Vt = g.parseInline;
var en = b.parse;
var tn = x.lex;

// static/scripts/guides.js
var guideWindow = document.querySelector(".guide-overlay");
var markdownContainer = document.querySelector(".markdown-body");
var closeBtn = document.getElementById("guide-close");
var guidesCont = document.getElementById("guideList");
var docHeaders = [];
var checkForFirstGuide = async () => {
  const current = localStorage.getItem("new");
  if (current == null) {
    openGuides();
    localStorage.setItem("new", "1");
  }
};
var openGuides = async () => {
  guideWindow.style.display = "flex";
  initButtons();
  await fetchAndFillGuides();
};
var initButtons = () => {
  closeBtn.addEventListener("click", () => {
    guideWindow.style.display = "none";
  });
};
var fetchAndFillGuides = async () => {
  const response = await fetch(`./docs`);
  const guides = await response.json();
  guides.forEach((element) => {
    const guide = document.createElement("div");
    guide.innerText = element;
    guide.classList.add("guide-item");
    guide.id = element;
    docHeaders.push(element);
    if (element == "Introduction.md") {
      guide.classList.add("active");
      fetchAndFill(element);
    }
    guide.addEventListener("click", () => {
      fetchAndFill(element);
    });
    guidesCont.appendChild(guide);
  });
};
var makeActive = (name) => {
  for (const child of guidesCont.children) {
    if (child.classList.contains("active")) {
      child.classList.remove("active");
    }
  }
  const item = document.getElementById(name);
  if (item) {
    item.classList.add("active");
  }
};
var fetchAndFill = async (name) => {
  makeActive(name);
  const response = await fetch(`./docs/${name}`);
  const md = await response.text();
  markdownContainer.innerHTML = g.parse(md);
  hijackLinks(g.parse(md));
  markdownContainer.scrollTo({ top: 0, behavior: "smooth" });
};
var hijackLinks = (content) => {
  const links = document.querySelectorAll("a");
  const start2 = window.location.href.length;
  links.forEach((element) => {
    const dest = element.href.slice(start2);
    if (dest.includes("docs")) {
      const doc = dest.slice(5);
      console.log(doc, docHeaders);
      if (docHeaders.includes(doc)) {
        console.log("included hiajck");
        element.addEventListener("click", (e) => {
          e.preventDefault();
          fetchAndFill(doc);
        });
      }
    }
  });
};

// static/scripts/rewind.js
var track = document.getElementById("items-container");
var preview_img = document.getElementById("preview-img");
var rewind_overlay = document.getElementById("rewind-overlay");
var timeline = document.getElementById("rewind-timeline");
var highligh_view = document.getElementById("rewind-highlight");
var back = document.getElementById("rewind-back");
var next = document.getElementById("rewind-next");
var start = document.getElementById("rewind-start");
var end = document.getElementById("rewind-end");
var load = document.getElementById("review-load");
var cart = document.getElementById("review-cart");
var saveCap = document.getElementById("rewind-cap");
var saveArea = document.getElementById("rewind-area");
var slot2 = document.getElementById("rewind-slot");
var activeID = 0;
var clicker = new IntersectionObserver(
  (entries) => {
    entries.forEach((entry) => {
      if (entry.isIntersecting) {
        clickSound();
        makeActive2(entry.target);
        updateHighlight();
      }
    });
  },
  {
    root: track,
    rootMargin: "0px -45% 0px -45%",
    threshold: 0.1
  }
);
var ctx = new AudioContext();
var createTilesFromSnapshots = async (snapshots) => {
  console.log("creating snapshot tiles");
  openCap();
  let lastTile = null;
  rewind_overlay.style.display = "flex";
  setupRewindButtons(snapshots.length);
  setupTimeline(snapshots.length);
  for (let i = 0; i < snapshots.length; i++) {
    let snapshot = snapshots[i];
    let tile = document.createElement("div");
    tile.id = "tile-" + i;
    clicker.observe(tile);
    tile.classList.add("tape-item");
    let img = document.createElement("canvas");
    img.height = snapshot.image.height;
    img.width = snapshot.image.width;
    tile.addEventListener("mousedown", () => {
      makeActive2(tile);
    });
    img.getContext("2d").drawImage(snapshot.image, 0, 0);
    lastTile = tile;
    tile.append(img);
    track.append(tile);
  }
  lastTile.scrollIntoView({
    behavior: "smooth",
    inline: "nearest",
    block: "center"
  });
};
var setupRewindButtons = (length) => {
  next.addEventListener("click", () => {
    let target = parseInt(activeID, 10) + 1;
    if (target > length - 1) {
      target = length - 1;
    }
    let targetEl = document.getElementById("tile-" + target);
    makeActive2(document.getElementById("tile-" + target));
    scrollReelToValue(target, convertValToPerc().center);
  });
  back.addEventListener("click", () => {
    let target = parseInt(activeID, 10) - 1;
    if (target < 0) {
      target = 0;
    }
    let targetEl = document.getElementById("tile-" + target);
    makeActive2(targetEl);
    scrollReelToValue(target, convertValToPerc().center);
  });
  end.addEventListener("click", async () => {
    let target = length - 1;
    let targetEl = document.getElementById("tile-" + target);
    scrollReelToValue(target, false);
    await wait(200);
    makeActive2(targetEl);
  });
  start.addEventListener("click", async () => {
    let target = 0;
    let targetEl = document.getElementById("tile-" + target);
    scrollReelToValue(target, false);
    await wait(200);
    makeActive2(targetEl);
  });
  load.addEventListener("click", async () => {
    console.log("loading snapshot:", activeID);
    await startSaveLoad();
    await loadSnap(parseInt(activeID, 10));
  });
};
var startSaveLoad = async () => {
  let clone = cart.cloneNode(true);
  clone.classList.add("rewind-cart-clone");
  let originalRect = cart.getBoundingClientRect();
  let slotRect = slot2.getBoundingClientRect();
  let areaRect = saveArea.getBoundingClientRect();
  let areaHeight = areaRect.top - areaRect.bottom;
  let sourceCanvas = cart.querySelector("canvas");
  let width = originalRect.right - originalRect.left;
  let targetWidth = slotRect.right - slotRect.left;
  let scaleFactor2 = targetWidth / width;
  clone.style.position = "fixed";
  clone.style.left = originalRect.left + "px";
  clone.style.top = originalRect.top + "px";
  clone.style.width = width + "px";
  clone.style.transformOrigin = "top left";
  clone.style.transition = "transform 0.8s cubic-bezier(0.34, 1.56, 0.67 , 1)";
  clone.style.zIndex = 0;
  clone.querySelector("canvas").getContext("2d").drawImage(sourceCanvas, 0, 0);
  document.body.appendChild(clone);
  await wait(10);
  rewind_overlay.style.display = "none";
  let dx = slotRect.left - originalRect.left;
  let dy = slotRect.top - originalRect.top;
  clone.style.transform = `translate(${dx}px, ${dy}px) scale(${scaleFactor2})`;
  await wait(1e3);
  let movedown = slotRect.bottom - slotRect.top;
  clone.style.transition = "transform 0.8s ease";
  clone.style.transform = `translate(${dx}px, ${dy - areaHeight}px) scale(${scaleFactor2})`;
  closeCap();
};
var openCap = () => {
  console.log(slot2);
  let Rect = saveArea.getBoundingClientRect();
  let dy = Rect.top - Rect.bottom;
  console.log(dy);
  let capRect = saveCap.getBoundingClientRect();
  let capheight = capRect.bottom - capRect.top;
  saveArea.style.transform = `translateY(${dy}px)`;
  saveCap.style.transform = `translateY(${dy - capheight}px)`;
};
var closeCap = () => {
  saveArea.style.transform = `translateY(0)`;
  saveCap.style.transform = `translateY(0)`;
};
var lastactive = null;
var makeActive2 = (el) => {
  if (lastactive != null) {
    lastactive.classList.remove("tape-active");
  }
  let source = el.querySelector("canvas");
  preview_img.height = 240;
  preview_img.width = 256;
  preview_img.getContext("2d").drawImage(source, 0, 0);
  el.classList.add("tape-active");
  let index = el.id.slice(5);
  timeline.value = index;
  activeID = index;
  lastactive = el;
};
var setupTimeline = (val) => {
  timeline.min = 0;
  timeline.max = val - 1;
};
timeline.addEventListener("input", (e) => {
  scrollReelToValue(timeline.value, convertValToPerc().center);
});
var updateHighlight = () => {
  let res = convertValToPerc();
  highligh_view.style.left = res.perc - 6 + "%";
  return res.center;
};
var scrollReelToValue = (value, center) => {
  let elemtentID = "tile-" + Math.round(value);
  let tile = document.getElementById(elemtentID);
  tile.scrollIntoView({
    behavior: "instant",
    inline: center ? "center" : "nearest",
    block: "center"
  });
};
var convertValToPerc = () => {
  let max = timeline.max;
  let min = timeline.min;
  let val = timeline.value;
  let perc = (val - min) / (max - min) * 100;
  let center = true;
  if (perc < 6) {
    perc = 6;
    center = false;
  }
  if (perc > 94) {
    perc = 94;
    center = false;
  }
  return { perc, center };
};
var clickSound = () => {
  if (ctx == null) return;
  const osc = ctx.createOscillator();
  const gain2 = ctx.createGain();
  osc.connect(gain2);
  gain2.connect(ctx.destination);
  osc.frequency.value = 900;
  gain2.gain.setValueAtTime(0.4, ctx.currentTime);
  gain2.gain.exponentialRampToValueAtTime(1e-3, ctx.currentTime + 0.02);
  osc.start();
  osc.stop(ctx.currentTime + 0.02);
};

// static/scripts/driver.js
var worker = new Worker(new URL("./emuWorker.js", import.meta.url));
var fBytes = null;
var inputBuf = new SharedArrayBuffer(4);
var inputState = new Int32Array(inputBuf);
var cmdMap = {
  "STOP": 1,
  "RESET": 2
};
var canvas = document.getElementById("screen");
var ctx2 = canvas.getContext("2d");
var imageData = ctx2.createImageData(256, 240);
var speedBuf = new SharedArrayBuffer(4);
var speedNum = new Int32Array(speedBuf);
var NES_FPS = 60.0948;
var FRAME_MS = 1e3 / NES_FPS;
var frameSig = null;
var rafMode = false;
var rafId = null;
var lastT = 0;
var acc = 0;
var audioBufS = null;
var AudioControl = null;
var gameControl = null;
Atomics.store(speedNum, 0, 1e3);
var state = {
  romRunning: false,
  runMode: 0
  //0-aduio run 1 - raf run
};
var audioCtx = null;
var gain = null;
window.addEventListener("keydown", async (e) => {
  if (e.code == "KeyR") {
    await getSnapList();
  } else if (e.code == "KeyT") {
    await PauseGame();
    worker.postMessage({ type: "reset" });
  }
});
var setUpAudio = async (audioBufS2, SIZE) => {
  if (audioCtx !== null && gain !== null) {
    return;
  }
  audioCtx = new AudioContext({ sampleRate: 44100 });
  await audioCtx.audioWorklet.addModule(new URL("./driverWorklet.js", import.meta.url));
  const node = new AudioWorkletNode(audioCtx, "apu-proc", {
    outputChannelCount: [1]
  });
  gain = audioCtx.createGain();
  node.port.postMessage({ audioBufS: audioBufS2, SIZE });
  node.connect(gain);
  gain.connect(audioCtx.destination);
};
var PauseGame = async () => {
  if (audioCtx == null) {
    return;
  }
  await audioCtx.suspend();
  Atomics.store(gameControl, 0, 2);
  Atomics.notify(gameControl, 0);
};
var ResumeGame = async () => {
  if (audioCtx == null) {
    return;
  }
  await audioCtx.resume();
  Atomics.store(gameControl, 0, cmdMap.STOP);
  Atomics.notify(gameControl, 0);
  worker.postMessage({ "type": "pump" });
};
var getSnapList = async () => {
  await PauseGame();
  worker.postMessage({ type: "getsnap" });
};
var loadSnap = async (index) => {
  worker.postMessage({ type: "loadSnapshot", index });
  await wait(100);
  await ResumeGame();
};
var updateVolume = (intensity) => {
  if (gain !== null) {
    gain.gain.setValueAtTime(gain.gain.value, audioCtx.currentTime);
    gain.gain.linearRampToValueAtTime(intensity, audioCtx.currentTime + 1);
  }
};
var initConsole = () => {
  worker.postMessage({ type: "init", speedBuf, inputBuf });
};
worker.onmessage = async ({ data }) => {
  switch (data.type) {
    case "init":
      fBytes = new Uint8Array(data.FBuf);
      audioBufS = data.audioBufS;
      AudioControl = new Int32Array(audioBufS, data.SIZE * 4, 3);
      gameControl = new Int32Array(data.gameBuf);
      frameSig = new Int32Array(data.frameSigBuf);
      await setUpAudio(data.audioBufS, data.SIZE);
      break;
    case "wasm":
      createModal("Console Ready!!", "Press the power button to start the console");
      worker.postMessage({ type: "init", inputBuf, speedBuf });
      break;
    case "snaps":
      await createTilesFromSnapshots(data.snaps);
      break;
    case "start":
      console.log("starting/restarting");
      await startConsole();
      break;
    case "frameUp":
      imageData.data.set(fBytes);
      ctx2.putImageData(imageData, 0, 0);
      break;
  }
};
var startConsole = async () => {
  await ResumeGame();
  if (state.romRunning) {
    return;
  }
  state.romRunning = true;
  console.log("starting the fucking game");
  if (state.runMode == 0) {
    await startAudioBuf();
  } else {
    await startRaf();
  }
};
var ControlMap = {
  "joypad-A": 0,
  "joypad-B": 1,
  "joypad-select": 2,
  "joypad-start": 3,
  "dpad-up": 4,
  "dpad-down": 5,
  "dpad-left": 6,
  "dpad-right": 7
};
var switchMode = (mode) => {
  if (mode == 0) {
    state.runMode = 0;
    startAudioBuf();
  } else {
    state.runMode = 1;
    startRaf();
  }
};
var loadRom = async (game) => {
  await PauseGame();
  console.log("loading game");
  worker.postMessage({ type: "loadRom", rom: game });
};
var UpdateSpeed = (speed) => {
  Atomics.store(speedNum, 0, speed);
};
var UpdatePress = (btn) => {
  Atomics.or(inputState, 0, 1 << ControlMap[btn]);
};
var UpdateRelease = (btn) => {
  Atomics.and(inputState, 0, ~(1 << ControlMap[btn]));
};
var startAudioBuf = async () => {
  if (!state.romRunning) return;
  stopRaF();
  if (audioCtx) {
    await audioCtx.resume();
  }
  await ResumeGame();
  worker.postMessage({ type: "pump" });
};
var startRaf = () => {
  if (!state.romRunning) return;
  if (!frameSig) return;
  if (gameControl) {
    Atomics.store(gameControl, 0, cmdMap.STOP);
    Atomics.notify(gameControl, 0);
  }
  if (audioCtx) audioCtx.suspend();
  Atomics.store(frameSig, 2, 0);
  Atomics.store(frameSig, 1, Atomics.load(frameSig, 0));
  rafMode = true;
  lastT = performance.now();
  acc = 0;
  worker.postMessage({ type: "startRaF" });
  rafId = requestAnimationFrame(rafLoop);
};
var rafLoop = (now) => {
  if (!rafMode) return;
  acc += now - lastT;
  lastT = now;
  if (acc > 250) acc = 250;
  while (acc > FRAME_MS) {
    Atomics.add(frameSig, 0, 1);
    Atomics.notify(frameSig, 0);
    acc -= FRAME_MS;
  }
  rafId = requestAnimationFrame(rafLoop);
};
var stopRaF = () => {
  if (rafId) cancelAnimationFrame(rafId);
  rafId = null;
  rafMode = false;
  if (frameSig) {
    Atomics.store(frameSig, 2, 1);
    Atomics.add(frameSig, 0, 1);
    Atomics.notify(frameSig, 0);
  }
};

// static/scripts/joypad.js
var updateBtn = document.getElementById("control-update");
var buttons = document.querySelectorAll(".joypad-button");
var UpdatingKey = "";
var updatingControls = false;
var keyMap = {
  "KeyZ": "joypad-A",
  "KeyX": "joypad-B",
  "ShiftLeft": "joypad-select",
  "Enter": "joypad-start",
  "ArrowUp": "dpad-up",
  "ArrowDown": "dpad-down",
  "ArrowLeft": "dpad-left",
  "ArrowRight": "dpad-right"
};
var neededControls = ["joypad-A", "joypad-B", "joypad-select", "joypad-start", "dpad-up", "dpad-down", "dpad-left", "dpad-right"];
document.addEventListener("DOMContentLoaded", () => {
  updateBtn.addEventListener("click", () => {
    if (updatingControls) {
      closeControlPanel();
      removeUpdateListeners();
    } else {
      openControlPanel();
      handleUpdateListeners();
    }
    updatingControls = !updatingControls;
    UpdatingKey = "";
  });
});
var controller = null;
var handleUpdateListeners = () => {
  controller = new AbortController();
  const { signal } = controller;
  buttons.forEach((element) => {
    element.addEventListener("mousedown", async (e) => {
      console.log(element.id, getKeyFromAction(element.id));
      let x2 = getKeyFromAction(element.id);
      console.log(x2);
      addUpdatePanel(element.id, getKeyFromAction(element.id));
      UpdatingKey = element.id;
    }, { signal });
  });
};
var getKeyFromAction = (action) => {
  let keys = Object.keys(keyMap);
  let res = "";
  keys.forEach((key) => {
    if (keyMap[key] == action) {
      res = key;
    }
  });
  return res;
};
var removeUpdateListeners = () => {
  if (controller) {
    controller.abort();
  }
};
var updateBinding = (action, newKey) => {
  let old = structuredClone(keyMap);
  let Keys = Object.keys(keyMap);
  Keys.forEach((a) => {
    if (keyMap[a] == action) {
      delete keyMap[a];
    }
  });
  keyMap[newKey] = action;
  checkForMissingKeys();
};
var checkForMissingKeys = () => {
  let keys = Object.keys(keyMap);
  let needed = structuredClone(neededControls);
  let maped = [];
  keys.forEach((key) => {
    maped.push(keyMap[key]);
  });
  needed.forEach((action) => {
    let x2 = document.getElementById(action);
    if (maped.includes(action)) {
      if (x2.classList.contains("key-missing")) {
        x2.classList.remove("key-missing");
      }
    } else {
      x2.classList.add("key-missing");
    }
  });
};
window.addEventListener("keydown", (e) => {
  if (UpdatingKey !== "" && updatingControls) {
    updateKey(e.code);
    updateBinding(UpdatingKey, e.code);
  }
  if (keyMap[e.code] !== void 0 && state.romRunning) {
    UpdatePress(keyMap[e.code]);
    let btn = document.getElementById(keyMap[e.code]);
    PressBtn(btn);
  }
});
window.addEventListener("keyup", (e) => {
  if (keyMap[e.code] !== void 0 && state.romRunning) {
    UpdateRelease(keyMap[e.code]);
    let btn = document.getElementById(keyMap[e.code]);
    ReleaseBtn(btn);
  }
});
var PressBtn = async (button) => {
  button.classList.add("active");
};
var ReleaseBtn = async (button) => {
  button.classList.remove("active");
};
function wait(ms) {
  return new Promise((r) => setTimeout(r, ms));
}

// static/scripts/modal.js
var createModal = async (head, body2, error = false, fleeting = true) => {
  let modalList = document.querySelector(".moodle-bar");
  const modal = document.createElement("div");
  modal.classList.add("modal", "modal-corners");
  const title = document.createElement("h1");
  if (error) {
    title.classList.add("error");
  } else {
    title.classList.add("info");
  }
  title.textContent = head;
  modal.appendChild(title);
  const para = document.createElement("p");
  para.classList.add("modal-text");
  para.textContent = body2;
  modal.appendChild(para);
  modalList.appendChild(modal);
  if (fleeting) {
    setUpForFailure(modal);
  }
  activateModal(modal);
  return modal;
};
var activateModal = (modal) => {
  modal.classList.add("active");
};
var deactivateModal = async (modal) => {
  modal.classList.remove("active");
  modal.classList.add("inactive");
  await wait(400);
  modal.remove();
};
var setUpForFailure = async (modal) => {
  await wait(3500);
  if (modal.classList.contains("active")) {
    await deactivateModal(modal);
  } else {
    console.log("bro alr failed");
  }
};

// static/scripts/tooltips.js
var activateTip = (tip) => {
  if (TipsState[tip]["state"]) {
    TipsState[tip]["function"]();
    TipsState[tip]["state"] = false;
  }
};
var startRandomTipEngine = async () => {
  await wait(1e3);
  scheduleTip();
};
var scheduleTip = () => {
  let delay = Math.random() * (6e4 - 1e4) + 1e4;
  setTimeout(() => {
    let tip = tips[Math.floor(Math.random() * tips.length)];
    TipsState[tip]["function"]();
    scheduleTip();
  }, delay);
};
var tip_fullscreen = () => {
  createModal("Go big or go home", "click on the tv screen to go fullscreen mode (resolution not garunteed)");
};
var tip_knobs = () => {
  createModal("Play with the dials!!", "Adjust the dials on the television to control the sound and speed levels");
};
var tip_swap = () => {
  createModal("You dont need to reload(prolly)", "Click the cartridge slot to change games while the console is running");
};
var tip_feedback = () => {
  createModal("Have suggestions or found errors?", "Dm emptiedfull on slack please");
};
var TipsState = {
  "fullscreen": { "state": true, "function": tip_fullscreen },
  "knobs": { "state": true, "function": tip_knobs },
  "hotswap": { "state": true, "function": tip_swap },
  "feedback": { "state": true, "function": tip_feedback },
  "controls": {}
};
var tips = Object.keys(TipsState);

// static/scripts/graphics.js
var body = document.querySelector("body");
var middle = document.getElementById("middle");
var romled = document.getElementById("rom");
var cap = document.getElementById("cartridge");
var screen = document.getElementById("screen");
var control_cont = document.getElementById("joypad-cont");
var panel = document.getElementById("update-panel");
var joypad = document.getElementById("joypad");
var updateBtn2 = document.getElementById("control-update");
var keyDisplay = document.getElementById("key-display");
var keyText = document.getElementById("show-text");
var knobSettings = { "speed": { "angle": 0, "setting": 0 }, "sound": { "angle": 0, "setting": 0 } };
var scaleFactor = 2;
function wait2(ms) {
  return new Promise((r) => setTimeout(r, ms));
}
window.addEventListener("DOMContentLoaded", async () => {
  initCables();
  drawNav();
  initCanvas();
  drawKnob("sound", 0);
  drawKnob("speed", 0);
});
var addUpdatePanel = async (action, key) => {
  keyText.innerText = action + ":  ";
  await updateKey(key);
};
var flashAnim = keyDisplay.animate([
  { opacity: 1 },
  { opacity: 0.2 },
  { opacity: 1 }
], {
  duration: 1e3,
  iterations: Infinity
});
flashAnim.pause();
var updateKey = async (key) => {
  if (key == void 0 || key == "") {
    keyDisplay.style.backgroundImage = "none";
    keyDisplay.innerText = "NA";
    keyDisplay.removeAttribute("style");
    flashAnim.play();
    return;
  }
  const res = await fetch(`./keys/${key.toLocaleLowerCase()}.png`);
  if (res.ok) {
    flashAnim.cancel();
    keyDisplay.innerText = "";
    const blob = await res.blob();
    const qwn = await createImageBitmap(blob);
    let h = qwn.height * scaleFactor;
    let wFull = (qwn.width - 2) * scaleFactor;
    let w2 = wFull / 2;
    keyDisplay.style.height = h + "px";
    keyDisplay.style.width = w2 + "px";
    keyDisplay.style.backgroundSize = `${wFull}px ${h}px `;
    keyDisplay.style.backgroundImage = `url('./keys/${key.toLocaleLowerCase()}.png')`;
    keyDisplay.animate([
      { backgroundPosition: "0px 0px" },
      { backgroundPosition: `-${wFull}px 0px` }
    ], {
      duration: 700,
      easing: "steps(2)",
      iterations: Infinity
    });
  }
};
var initCanvas = () => {
  screen.addEventListener("mouseenter", async () => {
    activateTip("fullscreen");
  });
  screen.addEventListener("click", async () => {
    try {
      await screen.requestFullscreen();
    } catch (e) {
      await createModal("Unable to go fullscreen", e, true);
    }
  });
};
var initCables = () => {
  let cable1 = document.getElementById("cable-1");
  let port1 = document.getElementById("port-1");
  let wire1 = document.getElementById("wire-1");
  let cable2 = document.getElementById("cable-2");
  let port2 = document.getElementById("port-2");
  let wire2 = document.getElementById("wire-2");
  startCable(cable1, port1, wire1);
  startCable(cable2, port2, wire2);
};
var alignCable = (cableID) => {
  if (cableID == 1) {
    let cable1 = document.getElementById("cable-1");
    let port1 = document.getElementById("port-1");
    let wire1 = document.getElementById("wire-1");
    alignCables(cable1, port1, wire1);
  } else if (cableID == 2) {
    let cable2 = document.getElementById("cable-2");
    let port2 = document.getElementById("port-2");
    let wire2 = document.getElementById("wire-2");
    alignCables(cable2, port2, wire2);
  }
};
var alignCables = (cable, port, wire) => {
  let portBox = port.getBoundingClientRect();
  cable.style.top = portBox.top + "px";
  cable.style.left = portBox.left + "px";
  let bodyBox = body.getBoundingClientRect();
  let dy = bodyBox.bottom - portBox.bottom + 200;
  let dx = portBox.width / 4;
  wire.style.height = dy + "px";
  wire.style.top = portBox.top + portBox.height / 2 + "px";
  wire.style.left = portBox.left + portBox.width / 4 + "px";
  wire.style.width = dx + "px";
};
var startCable = (cable, port, wire) => {
  let portBox = port.getBoundingClientRect();
  let bodyBox = document.querySelector("body").getBoundingClientRect();
  cable.style.position = "absolute";
  cable.style.width = portBox.width + "px";
  cable.style.height = portBox.height + "px";
  cable.style.left = portBox.left + "px";
  cable.style.top = bodyBox.bottom + "px";
  wire.style.position = "absolute";
  let dy = bodyBox.bottom - portBox.bottom + 200;
  let dx = portBox.width / 4;
  wire.style.height = dy + "px";
  wire.style.top = bodyBox.bottom + portBox.height / 2 + "px";
  wire.style.left = portBox.left + portBox.width / 4 + "px";
};
var activateJoypad = () => {
  let cont = document.getElementById("joypad-cont");
  cont.classList.add("active");
};
var openControlPanel = () => {
  control_cont.classList.add("updating");
  panel.classList.add("active");
  panel.style.pointerEvents = "all";
  updateBtn2.innerText = "Save Controls";
};
var closeControlPanel = () => {
  control_cont.classList.remove("updating");
  panel.classList.remove("active");
  panel.style.pointerEvents = "none";
  updateBtn2.innerText = "Update Controls";
};
var C2 = {
  //this color pallete was ai generated
  body: "#8B6F4E",
  bodyMid: "#7A5F3E",
  hi: "#C4A882",
  hiTop: "#D9C4A0",
  shadow: "#4A3522",
  shadowD: "#332614",
  rim: "#3B2A16",
  rimHi: "#6B5030",
  dot: "#1E1208",
  dotHi: "#3A2810"
};
function drawKnob(id, angleDeg) {
  knobSettings[id].angle = angleDeg;
  const canvas2 = document.getElementById(id);
  const ctx3 = canvas2.getContext("2d");
  const width = canvas2.width;
  const height = canvas2.height;
  const cx = Math.floor(width / 2);
  const cy = Math.floor(height / 2);
  const r = Math.floor(width / 2) - 1;
  const fill = (x2, y2, col) => {
    if (x2 < 0 || y2 < 0 || x2 >= width || y2 >= height) return;
    ctx3.fillStyle = col;
    ctx3.fillRect(x2, y2, 1, 1);
  };
  for (let py = -r; py <= r; py++) {
    for (let px = -r; px <= r; px++) {
      const dist = Math.sqrt(px * px + py * py);
      if (dist > r + 0.5) continue;
      const x2 = cx + px;
      const y2 = cy + py;
      if (dist > r - 1) {
        if (px < 0 && py < 0) {
          fill(x2, y2, C2.rimHi);
        } else {
          fill(x2, y2, C2.rim);
        }
        continue;
      }
      let col = C2.body;
      if (px <= -Math.floor(r * 0.1) && py <= -Math.floor(r * 0.1) && dist < r * 0.75) {
        col = C2.hi;
      }
      if (px <= -Math.floor(r * 0.3) && py <= -Math.floor(r * 0.3) && dist < r * 0.45) {
        col = C2.hiTop;
      }
      if (px >= Math.floor(r * 0.2) && py >= Math.floor(r * 0.2)) {
        col = C2.bodyMid;
      }
      if (px >= Math.floor(r * 0.45) && py >= Math.floor(r * 0.45)) {
        col = C2.shadow;
      }
      if (px >= Math.floor(r * 0.65) && py >= Math.floor(r * 0.65)) {
        col = C2.shadowD;
      }
      fill(x2, y2, col);
    }
  }
  const rad = (angleDeg - 90) * Math.PI / 180;
  const len = r - 3;
  const ex = Math.round(cx + Math.cos(rad) * len);
  const ey = Math.round(cy + Math.sin(rad) * len);
  const steps = Math.max(Math.abs(ex - cx), Math.abs(ey - cy));
  for (let i = 1; i <= steps; i++) {
    const t = i / steps;
    const px = Math.round(cx + (ex - cx) * t);
    const py = Math.round(cy + (ey - cy) * t);
    const d2 = Math.sqrt((px - cx) ** 2 + (py - cy) ** 2);
    if (d2 < r - 1) fill(px, py, i < 2 ? C2.dotHi : C2.dot);
  }
  fill(cx, cy, C2.dot);
}
async function turnKnob(id, targetAngle, durationMs = 400, steps = 20) {
  const startAngle = knobSettings[id].angle;
  const diff = targetAngle - startAngle;
  const stepTime = durationMs / steps;
  for (let i = 1; i <= steps; i++) {
    const t = i / steps;
    const eased = t < 0.5 ? 2 * t * t : -1 + (4 - 2 * t) * t;
    const angle = startAngle + diff * eased;
    drawKnob(id, angle);
    await wait2(stepTime);
  }
  drawKnob(id, targetAngle);
}
async function wiggleKnob(id, intensity = 15) {
  const current = knobSettings[id].angle;
  await turnKnob(id, current + intensity, 80, 6);
  await turnKnob(id, current - intensity, 100, 8);
  await turnKnob(id, current + intensity * 0.5, 70, 5);
  await turnKnob(id, current, 80, 6);
}
var closeCap2 = async () => {
  cap.style.transition = "all ease 1s";
  cap.classList.remove("open");
  await wait2(200);
};
function pushbtn(canvasId) {
  const c = document.getElementById(canvasId);
  c.style.filter = "brightness(2)";
  c.style.transform = "scale(0.85)";
  setTimeout(() => {
    c.style.filter = "brightness(1.3)";
    c.style.transform = "scale(1.1)";
  }, 150);
  setTimeout(() => {
    c.style.filter = "";
    c.style.transform = "scale(1)";
  }, 280);
}
var openCap2 = async () => {
  cap.style.transition = "all ease 1s";
  cap.classList.add("open");
  await await 200;
};
var slotCart = async () => {
  await wait2(800);
  overlay.style.opacity = 0;
  const spine = middle.querySelector(".rom-spine");
  const spineRect = spine.getBoundingClientRect();
  const clone = spine.cloneNode(true);
  const computedStyle = window.getComputedStyle(spine);
  clone.style.font = computedStyle.font;
  clone.style.color = computedStyle.color;
  clone.style.letterSpacing = computedStyle.letterSpacing;
  clone.style.textTransform = computedStyle.textTransform;
  clone.style.position = "fixed";
  clone.style.top = spineRect.top + "px";
  clone.style.left = spineRect.left + "px";
  clone.style.width = spineRect.width + "px";
  clone.style.height = slot.height + "px";
  clone.style.transform = "none";
  clone.style.margin = "0";
  clone.style.zIndex = "3";
  document.body.appendChild(clone);
  const slotRect = slot.getBoundingClientRect();
  const stripRect = strip.getBoundingClientRect();
  let dy = stripRect.left - slotRect.left;
  clone.style.transition = "all 1s ease";
  clone.style.top = slotRect.top + "px";
  clone.style.left = slotRect.left + "px";
  clone.style.width = dy + "px";
  await wait2(1e3);
  clone.style.transition = "all 0.25s ease-in";
  clone.style.transform = " scale(0.88)";
  clone.style.transformOrigin = "center top";
  await wait2(500);
  clone.style.transition = "all 0.15s ease-out";
  clone.style.transform = "scale(0.9)";
  flickerLed(romled);
  await wait2(500);
  await closeCap2();
  cleapUpOverlay();
};
var flickerLed = async (led) => {
  led.classList.remove("off");
  await wait2(60);
  led.classList.add("off");
  await wait2(80);
  led.classList.remove("off");
  await wait2(150);
  led.classList.add("off");
  await wait2(120);
  led.classList.remove("off");
};
var cleapUpOverlay = () => {
  overlay.style.display = "none";
  middle.classList.remove("rotated");
  middle.classList.add("active");
};
var drawNav = () => {
  var c = document.getElementById("canvas-back");
  var ctx3 = c.getContext("2d");
  var PS = 16, COLS = 14, ROWS = 14;
  var T2 = { hi: "#ffc040", md: "#b05000", dk: "#3a1400" };
  var G = [0, 0, 0, 0, 0, 0, 0, "hi", "hi", "hi", "hi", "hi", "hi", 0, 0, 0, 0, 0, 0, 0, "hi", "md", "md", "md", "md", "dk", 0, 0, 0, 0, 0, 0, 0, "hi", "md", "md", "md", "md", "dk", 0, 0, 0, 0, 0, 0, 0, "hi", "md", "md", "md", "md", "dk", 0, 0, 0, 0, 0, 0, 0, "hi", "md", "md", "md", "md", "dk", 0, 0, 0, 0, 0, 0, 0, "hi", "md", "md", "md", "md", "dk", 0, 0, 0, 0, 0, 0, 0, "hi", "md", "md", "md", "md", "md", 0, 0, 0, 0, 0, 0, 0, 0, "dk", "md", "md", "md", "md", "md", 0, 0, 0, 0, 0, 0, 0, 0, 0, "dk", "md", "md", "md", "md", "hi", 0, 0, 0, 0, 0, 0, 0, 0, 0, "dk", "md", "md", "md", "md", "hi", 0, 0, 0, 0, 0, 0, 0, 0, 0, "dk", "md", "md", "md", "md", "hi", 0, 0, 0, 0, 0, 0, 0, 0, 0, "dk", "md", "md", "md", "md", "hi", 0, 0, 0, 0, 0, 0, 0, 0, 0, "dk", "md", "md", "md", "md", "hi", 0, 0, 0, 0, 0, 0, 0, 0, 0, "dk", "dk", "dk", "dk", "dk", "hi", 0];
  for (var i = 0; i < G.length; i++) {
    if (!G[i]) continue;
    var r = Math.floor(i / COLS), col = i % COLS;
    ctx3.fillStyle = T2[G[i]];
    ctx3.fillRect(col * PS, r * PS, PS, PS);
  }
  var c = document.getElementById("canvas-next");
  var ctx3 = c.getContext("2d");
  var PS = 16, COLS = 14, ROWS = 14;
  var T2 = { hi: "#ffc040", md: "#b05000", dk: "#3a1400" };
  var G = [0, "hi", "hi", "hi", "hi", "hi", "hi", 0, 0, 0, 0, 0, 0, 0, 0, 0, "dk", "md", "md", "md", "md", "hi", 0, 0, 0, 0, 0, 0, 0, 0, 0, "dk", "md", "md", "md", "md", "hi", 0, 0, 0, 0, 0, 0, 0, 0, 0, "dk", "md", "md", "md", "md", "hi", 0, 0, 0, 0, 0, 0, 0, 0, 0, "dk", "md", "md", "md", "md", "hi", 0, 0, 0, 0, 0, 0, 0, 0, 0, "dk", "md", "md", "md", "md", "hi", 0, 0, 0, 0, 0, 0, 0, 0, 0, "md", "md", "md", "md", "md", "hi", 0, 0, 0, 0, 0, 0, 0, 0, "md", "md", "md", "md", "md", "dk", 0, 0, 0, 0, 0, 0, 0, "hi", "md", "md", "md", "md", "dk", 0, 0, 0, 0, 0, 0, 0, "hi", "md", "md", "md", "md", "dk", 0, 0, 0, 0, 0, 0, 0, "hi", "md", "md", "md", "md", "dk", 0, 0, 0, 0, 0, 0, 0, "hi", "md", "md", "md", "md", "dk", 0, 0, 0, 0, 0, 0, 0, "hi", "md", "md", "md", "md", "dk", 0, 0, 0, 0, 0, 0, 0, "hi", "dk", "dk", "dk", "dk", "dk", 0, 0, 0, 0, 0, 0, 0];
  for (var i = 0; i < G.length; i++) {
    if (!G[i]) continue;
    var r = Math.floor(i / COLS), col = i % COLS;
    ctx3.fillStyle = T2[G[i]];
    ctx3.fillRect(col * PS, r * PS, PS, PS);
  }
};
var setUpMainLinks = () => {
  const docLink = document.getElementById("docs");
  const hubLink = document.getElementById("github");
  docLink.addEventListener("click", () => {
    openGuides();
  });
  hubLink.addEventListener("click", () => {
    window.location.href = "https://github.com/Emptiedfull/ExNES/";
  });
};

// static/scripts/nes.js
var currentIndex = 3;
var middle2 = document.getElementById("middle");
var left = document.getElementById("left");
var right = document.getElementById("right");
var slot3 = document.getElementById("slot");
var overlay2 = document.getElementById("overlay");
var strip2 = document.getElementById("strip");
var cap2 = document.getElementById("cartridge");
var GamesArray = [];
var powerLed = document.getElementById("power");
var powerBtn = document.getElementById("start");
var pauseBtn = document.getElementById("pause");
var rocker = document.getElementById("rocker");
var paddle = document.getElementById("paddle");
var romLoaded = "";
var power = false;
var roms = [left, middle2, right];
document.addEventListener("DOMContentLoaded", async () => {
  checkForFirstGuide();
  setUpMainLinks();
  await startRandomTipEngine();
  await setUpRocker();
  await setUpKnobs();
  await setUpButtons();
  await initGames();
  await wait(500);
  updateRom();
});
var BeginKnobs = async () => {
  knobSettings["sound"].setting = 4;
  knobSettings["angle"] = 180;
  await turnKnob("sound", 180);
  await turnKnob("speed", 180);
};
var setUpRocker = () => {
  paddle.addEventListener("click", () => {
    const m2 = rocker.dataset.mode == "1" ? "0" : "1";
    rocker.dataset.mode = m2;
    paddle.innerHTML = m2 === "1" ? "rAF" : "Audio";
    switchMode(Number(m2));
  });
};
var setUpKnobs = async () => {
  const soundKnob = document.getElementById("sound");
  const soundIncrements = [0, 45, 90, 135, 180, 225, 270, 315, 360];
  soundKnob.addEventListener("click", async () => {
    if (romLoaded == "" || power == false) {
      await wiggleKnob("sound", 45);
      await turnKnob("sound", 0);
      return;
    }
    knobSettings["sound"].setting = knobSettings["sound"].setting + 1;
    if (knobSettings["sound"].setting >= soundIncrements.length - 1) {
      knobSettings["sound"].setting = 0;
    }
    let targetAngle = soundIncrements[knobSettings["sound"].setting];
    if (targetAngle == 0) {
      updateVolume(0);
    } else {
      updateVolume(targetAngle / 180);
    }
    await turnKnob("sound", targetAngle);
  });
  const speedKnob = document.getElementById("speed");
  const speedIncrements = [45, 90, 135, 180, 225, 270, 315, 360];
  speedKnob.addEventListener("click", async () => {
    if (romLoaded == "" || power == false) {
      await wiggleKnob("speed", 45);
      await turnKnob("speed", 0);
      return;
    }
    knobSettings["speed"].setting = knobSettings["speed"].setting + 1;
    if (knobSettings["speed"].setting >= speedIncrements.length - 1) {
      knobSettings["speed"].setting = 0;
    }
    let targetAngle = speedIncrements[knobSettings["speed"].setting];
    UpdateSpeed(1e3 * targetAngle / 180);
    await turnKnob("speed", targetAngle);
  });
};
var initGames = async () => {
  let response = await fetch("./games");
  if (response.ok) {
    let data = await response.json();
    data.forEach((element) => {
      GamesArray.push(element);
    });
  }
};
function updateRom() {
  let romSelection = GamesArray.slice(currentIndex - 1, currentIndex + 2);
  if (currentIndex == 0) {
    romSelection = [GamesArray.at(-1), GamesArray[0], GamesArray[1]];
  }
  if (currentIndex == GamesArray.length - 1) {
    romSelection = [GamesArray.at(-2), GamesArray.at(-1), GamesArray[0]];
  }
  if (romSelection.length == 3) {
    for (let i = 0; i < 3; i++) {
      const element = romSelection[i];
      const rom = roms[i];
      rom.id = element.ID;
      const rom_title = rom.querySelector("span");
      const img = rom.querySelector("img");
      img.src = "/rom_images/" + element.ID + ".webp";
      const spine = rom.querySelector(".rom-spine");
      spine.innerText = element.name;
      rom_title.innerText = element.name;
    }
  }
}
var setUpButtons = async () => {
  pauseBtn.addEventListener("click", async () => {
    if (pauseBtn.classList.contains("paused")) {
      pauseBtn.classList.remove("paused");
      pauseBtn.classList.add("unpaused");
      pauseBtn.innerHTML = `   <svg viewBox="0 0 24 24" fill="#5a3a1a" stroke="#5a3a1a" stroke-width="1.5" stroke-linecap="square" xmlns="http://www.w3.org/2000/svg">
                                     <rect x="4" y="3" width="4" height="18"/>
                                     <rect x="16" y="3" width="4" height="18"/></svg>`;
      await ResumeGame();
    } else {
      pauseBtn.classList.remove("unpaused");
      pauseBtn.classList.add("paused");
      pauseBtn.innerHTML = ` <svg viewBox="0 0 24 24" fill="#5a3a1a" stroke="#5a3a1a" stroke-width="1.5"
                                    stroke-linecap="square">
                                    <polygon points="6,3 20,12 6,21" />
                                </svg>`;
      await PauseGame();
    }
  });
  powerBtn.addEventListener("click", async () => {
    power = true;
    if (romLoaded != "") {
      begin(romLoaded);
      return;
    }
    powerLed.classList.remove("off");
    await openCap2(cap2);
    await wait(200);
    overlay2.style.display = "flex";
    overlay2.style.transition = "all ease 1";
    overlay2.style.opacity = 1;
  });
  const nav_back = document.getElementById("nav-back");
  nav_back.addEventListener("click", async () => {
    move(-1);
    pushbtn("canvas-back");
  });
  const nav_front = document.getElementById("nav-front");
  nav_front.addEventListener("click", async () => {
    move(1);
    pushbtn("canvas-next");
  });
  middle2.addEventListener("click", async () => {
    middle2.classList.remove("active");
    await wait(50);
    let span = middle2.querySelector("span");
    span.style.display = "none";
    middle2.classList.add("rotated");
    romLoaded = middle2.id;
    await slotCart();
    if (power) {
      activateJoypad();
      await begin(romLoaded);
    }
  });
  cap2.addEventListener("click", async () => {
    await PauseGame();
    await openCap2(cap2);
    await wait(200);
    overlay2.style.display = "flex";
    overlay2.style.transition = "all ease 1";
    overlay2.style.opacity = 1;
  });
};
async function begin(game) {
  initConsole();
  await loadRom(game);
  await BeginKnobs();
  alignCable(1);
  state.romRunning = true;
}
function move(direction) {
  let newIdx = currentIndex + direction;
  if (newIdx < 0) {
    currentIndex = GamesArray.length - 1;
    updateRom();
    return;
  }
  if (newIdx >= GamesArray.length) {
    currentIndex = newIdx - GamesArray.length;
    updateRom();
    return;
  }
  currentIndex = newIdx;
  updateRom();
}
