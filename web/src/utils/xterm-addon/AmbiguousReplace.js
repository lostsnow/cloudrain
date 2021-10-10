
let charsA1 = '─┌┎╓└┖╙' + '┬┭┰┱├┞┟┠┴┵┸┹' + '╥╟╨' + '┼╁╀╂┽╃╅╉' + '╫╭╰←';
let charsA2 = '━┏┍┗┕' + '┳┲┯┮┣┢┡┝┻┺┷┶' + '╋╇╈┿╊╆╄┾';
let charsA3 = '═╔╒╚╘' + '╦╠╧╞╤╩' + '╬╪';
let charsA0 = '│┃║' + '┐┤┘┑┥┙┒┨┚┓┫┛┩┪' + '╗╝╣╕╛╢╕╡╛╮╯' + '╱╲↑↓→≈≌◎①②③④⑤⑥⑦⑧⑨★◆';
let charsAA = '█▇▆▅▄▃▂▁▀▔┄┅┈┉';
let chars = charsA1 + charsA2 + charsA3 + charsA0 + charsAA;
let re = new RegExp('[' + chars + ']', 'g');

let replaceMap = function () {
  let m = {};
  for (var i = 0; i < chars.length; i++) {
    let c = chars[i]
    if (charsA1.indexOf(c) !== -1) {
      m[c] = c + '─';
    } else if (charsA2.indexOf(c) !== -1) {
      m[c] = c + '━';
    } else if (charsA3.indexOf(c) !== -1) {
      m[c] = c + '═';
    } else if (charsA0.indexOf(c) !== -1) {
      m[c] = c + ' ';
    } else if (charsAA.indexOf(c) !== -1) {
      m[c] = c + c;
    } else {
      m[c] = c;
    }
  }
  return m
}

let rMap = replaceMap();

export function AmbiguousReplace(msg) {
  msg = msg.replace(re, m => rMap[m]);
  return msg;
}