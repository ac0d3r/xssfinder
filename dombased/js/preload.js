/*!
 * https://github.com/AsaiKen/dom-based-xss-finder
 * License: MIT
 */

const __dombasedxssfinder_String = function (str, parent) {
    this.str = '' + str;
    this.sources = []; // 传播记录
    parent.sources.forEach(e => this.sources.push(e));

    this.valueOf = function () {
        return this;
    };

    this.toString = function () {
        return this.str;
    };

    // str.length
    Object.defineProperty(this, 'length', {
        set: () => null,
        get: () => this.str.length
    });

    // str[0]
    for (let i = 0; i < this.str.length; i++) {
        Object.defineProperty(this, i, {
            set: () => null,
            get: () => new __dombasedxssfinder_String(this.str[i], this)
        });
    }
    // hook String标记: __dombasedxssfinder_string
    Object.defineProperty(this, '__dombasedxssfinder_string', {
        set: () => null,
        get: () => true
    });
};
__dombasedxssfinder_String.prototype = String.prototype;

function __is_dombasedxssfinder_string(o) {
    return o && o.__dombasedxssfinder_string;
}

function __is_dombasedxssfinder_string_html(o) {
    // <svg/onload=alert()>
    o = __convert_to_dombasedxssfinder_string_if_location(o);
    return __is_dombasedxssfinder_string(o);
}

function __is_dombasedxssfinder_string_data_html(o) {
    // data:text/html,<script>alert(1)</script>
    o = __convert_to_dombasedxssfinder_string_if_location(o);
    return __is_dombasedxssfinder_string(o);
}

function __is_dombasedxssfinder_string_script(o) {
    // alert()
    // javascript:alert()
    o = __convert_to_dombasedxssfinder_string_if_location(o);
    return __is_dombasedxssfinder_string(o);
}

function __is_dombasedxssfinder_string_url(o) {
    // //14.rs
    o = __convert_to_dombasedxssfinder_string_if_location(o);
    return __is_dombasedxssfinder_string(o);
}

function __convert_to_dombasedxssfinder_string_if_location(o) {
    if (o === window.location) {
        o = new __dombasedxssfinder_String(o.toString(), {
            sources: [__dombasedxssfinder_get_source('window.location')],
        });
    }
    return o;
}

function __dombasedxssfinder_get_source(label) {
    return { label, stacktrace: __dombasedxssfinder_get_stacktrace() };
}

function __dombasedxssfinder_get_sink(label) {
    return { label, stacktrace: __dombasedxssfinder_get_stacktrace() };
}

// 堆栈信息
function __dombasedxssfinder_get_stacktrace() {
    const o = {};
    Error.captureStackTrace(o);
    // console.debug(o.stack.replace(/^Error\n/, '').replace(/^\s+at\s+/mg, ''));
    const regExp = /(https?:\/\/\S+):(\d+):(\d+)/;
    return o.stack.replace(/^Error\n/, '').replace(/^\s+at\s+/mg, '').split('\n')
        .filter(e => regExp.test(e))
        .map(e => {
            const m = e.match(regExp);
            const url = m[1];
            const line = m[2]; // start from 1
            const column = m[3]; // start from 1
            return { url, line, column, code: null };
        });
}


(function () {
    ///////////////////////////////////////////////
    // String.prototype
    ///////////////////////////////////////////////
    const stringPrototypeAnchor = String.prototype.anchor;
    String.prototype.anchor = function () {
        if (__is_dombasedxssfinder_string(this)) {
            const str = stringPrototypeAnchor.apply(this.toString(), arguments);
            return new __dombasedxssfinder_String(str, this);
        }
        return stringPrototypeAnchor.apply(this, arguments);
    };

    const stringPrototypeBig = String.prototype.big;
    String.prototype.big = function () {
        if (__is_dombasedxssfinder_string(this)) {
            const str = stringPrototypeBig.apply(this.toString(), arguments);
            return new __dombasedxssfinder_String(str, this);
        }
        return stringPrototypeBig.apply(this, arguments);
    };

    const stringPrototypeBlink = String.prototype.blink;
    String.prototype.blink = function () {
        if (__is_dombasedxssfinder_string(this)) {
            const str = stringPrototypeBlink.apply(this.toString(), arguments);
            return new __dombasedxssfinder_String(str, this);
        }
        return stringPrototypeBlink.apply(this, arguments);
    };

    const stringPrototypeBold = String.prototype.bold;
    String.prototype.bold = function () {
        if (__is_dombasedxssfinder_string(this)) {
            const str = stringPrototypeBold.apply(this.toString(), arguments);
            return new __dombasedxssfinder_String(str, this);
        }
        return stringPrototypeBold.apply(this, arguments);
    };

    const stringPrototypeCharAt = String.prototype.charAt;
    String.prototype.charAt = function () {
        if (__is_dombasedxssfinder_string(this)) {
            const str = stringPrototypeCharAt.apply(this.toString(), arguments);
            return new __dombasedxssfinder_String(str, this);
        }
        return stringPrototypeCharAt.apply(this, arguments);
    };

    const stringPrototypeCharCodeAt = String.prototype.charCodeAt;
    String.prototype.charCodeAt = function () {
        return stringPrototypeCharCodeAt.apply(this.toString(), arguments);
    };

    const stringPrototypeCodePointAt = String.prototype.codePointAt;
    String.prototype.codePointAt = function () {
        return stringPrototypeCodePointAt.apply(this.toString(), arguments);
    };

    const stringPrototypeConcat = String.prototype.concat;
    String.prototype.concat = function () {
        const sources = [];
        for (let i = 0; i < arguments.length; i++) {
            arguments[i] = __convert_to_dombasedxssfinder_string_if_location(arguments[i]);
            if (__is_dombasedxssfinder_string(arguments[i])) {
                arguments[i].sources.forEach(e => sources.push(e));
            }
        }
        if (__is_dombasedxssfinder_string(this)) {
            this.sources.forEach(e => sources.push(e));
        }
        if (sources.size > 0) {
            const str = stringPrototypeConcat.apply(this.toString(), arguments);
            return new __dombasedxssfinder_String(str, { sources });
        }
        return stringPrototypeConcat.apply(this, arguments);
    };

    const stringPrototypeEndsWith = String.prototype.endsWith;
    String.prototype.endsWith = function () {
        return stringPrototypeEndsWith.apply(this.toString(), arguments);
    };

    const stringPrototypeFixed = String.prototype.fixed;
    String.prototype.fixed = function () {
        if (__is_dombasedxssfinder_string(this)) {
            const str = stringPrototypeFixed.apply(this.toString(), arguments);
            return new __dombasedxssfinder_String(str, this);
        }
        return stringPrototypeFixed.apply(this, arguments);
    };

    const stringPrototypeFontcolor = String.prototype.fontcolor;
    String.prototype.fontcolor = function () {
        if (__is_dombasedxssfinder_string(this)) {
            const str = stringPrototypeFontcolor.apply(this.toString(), arguments);
            return new __dombasedxssfinder_String(str, this);
        }
        return stringPrototypeFontcolor.apply(this, arguments);
    };

    const stringPrototypeFontsize = String.prototype.fontsize;
    String.prototype.fontsize = function () {
        if (__is_dombasedxssfinder_string(this)) {
            const str = stringPrototypeFontsize.apply(this.toString(), arguments);
            return new __dombasedxssfinder_String(str, this);
        }
        return stringPrototypeFontsize.apply(this, arguments);
    };

    const stringPrototypeIncludes = String.prototype.includes;
    String.prototype.includes = function () {
        return stringPrototypeIncludes.apply(this.toString(), arguments);
    };

    const stringPrototypeIndexOf = String.prototype.indexOf;
    String.prototype.indexOf = function () {
        return stringPrototypeIndexOf.apply(this.toString(), arguments);
    };

    const stringPrototypeItalics = String.prototype.italics;
    String.prototype.italics = function () {
        if (__is_dombasedxssfinder_string(this)) {
            const str = stringPrototypeItalics.apply(this.toString(), arguments);
            return new __dombasedxssfinder_String(str, this);
        }
        return stringPrototypeItalics.apply(this, arguments);
    };

    const stringPrototypeLastIndexOf = String.prototype.lastIndexOf;
    String.prototype.lastIndexOf = function () {
        return stringPrototypeLastIndexOf.apply(this.toString(), arguments);
    };

    const stringPrototypeLink = String.prototype.link;
    String.prototype.link = function () {
        if (__is_dombasedxssfinder_string(this)) {
            const str = stringPrototypeLink.apply(this.toString(), arguments);
            return new __dombasedxssfinder_String(str, this);
        }
        return stringPrototypeLink.apply(this, arguments);
    };

    const stringPrototypeLocaleCompare = String.prototype.localeCompare;
    String.prototype.localeCompare = function () {
        return stringPrototypeLocaleCompare.apply(this.toString(), arguments);
    };

    const stringPrototypeMatch = String.prototype.match;
    // TODO propagate taints of the regexp argument
    String.prototype.match = function () {
        if (__is_dombasedxssfinder_string(this)) {
            const array = stringPrototypeMatch.apply(this.toString(), arguments);
            if (array === null) {
                return null;
            }
            for (let i = 0; i < array.length; i++) {
                array[i] = new __dombasedxssfinder_String(array[i], this);
            }
            return array;
        }
        return stringPrototypeMatch.apply(this, arguments);
    };

    const stringPrototypeMatchAll = String.prototype.matchAll;
    // TODO propagate taints of the regexp argument
    String.prototype.matchAll = function () {
        if (__is_dombasedxssfinder_string(this)) {
            const iterator = stringPrototypeMatchAll.apply(this.toString(), arguments);
            return function* () {
                for (const array of iterator) {
                    for (let i = 0; i < array.length; i++) {
                        array[i] = new __dombasedxssfinder_String(array[i], this);
                    }
                    yield array;
                }
            };
        }
        return stringPrototypeMatchAll.apply(this, arguments);
    };

    const stringPrototypeNormalize = String.prototype.normalize;
    String.prototype.normalize = function () {
        if (__is_dombasedxssfinder_string(this)) {
            const str = stringPrototypeNormalize.apply(this.toString(), arguments);
            return new __dombasedxssfinder_String(str, this);
        }
        return stringPrototypeNormalize.apply(this, arguments);
    };

    const stringPrototypePadEnd = String.prototype.padEnd;
    String.prototype.padEnd = function () {
        const sources = [];
        arguments[1] = __convert_to_dombasedxssfinder_string_if_location(arguments[1]);
        if (__is_dombasedxssfinder_string(arguments[1])) {
            arguments[1].sources.forEach(e => sources.push(e));
        }
        if (__is_dombasedxssfinder_string(this)) {
            this.sources.forEach(e => sources.push(e));
        }
        if (sources.size > 0) {
            const str = stringPrototypePadEnd.apply(this.toString(), arguments);
            return new __dombasedxssfinder_String(str, { sources });
        }
        return stringPrototypePadEnd.apply(this, arguments);
    };

    const stringPrototypePadStart = String.prototype.padStart;
    String.prototype.padStart = function () {
        const sources = [];
        arguments[1] = __convert_to_dombasedxssfinder_string_if_location(arguments[1]);
        if (__is_dombasedxssfinder_string(arguments[1])) {
            arguments[1].sources.forEach(e => sources.push(e));
        }
        if (__is_dombasedxssfinder_string(this)) {
            this.sources.forEach(e => sources.push(e));
        }
        if (sources.size > 0) {
            const str = stringPrototypePadStart.apply(this.toString(), arguments);
            return new __dombasedxssfinder_String(str, { sources });
        }
        return stringPrototypePadStart.apply(this, arguments);
    };

    const stringPrototypeRepeat = String.prototype.repeat;
    String.prototype.repeat = function () {
        if (__is_dombasedxssfinder_string(this)) {
            const str = stringPrototypeRepeat.apply(this.toString(), arguments);
            return new __dombasedxssfinder_String(str, this);
        }
        return stringPrototypeRepeat.apply(this, arguments);
    };

    const stringPrototypeReplace = String.prototype.replace;
    String.prototype.replace = function () {
        const sources = [];
        arguments[1] = __convert_to_dombasedxssfinder_string_if_location(arguments[1]);
        if (__is_dombasedxssfinder_string(arguments[1])) {
            arguments[1].sources.forEach(e => sources.push(e));
        }
        if (__is_dombasedxssfinder_string(this)) {
            this.sources.forEach(e => sources.push(e));
        }
        if (sources.size > 0) {
            const str = stringPrototypeReplace.apply(this.toString(), arguments);
            return new __dombasedxssfinder_String(str, { sources });
        }
        return stringPrototypeReplace.apply(this, arguments);
    };

    const stringPrototypeSearch = String.prototype.search;
    String.prototype.search = function () {
        return stringPrototypeSearch.apply(this.toString(), arguments);
    };

    const stringPrototypeSlice = String.prototype.slice;
    String.prototype.slice = function () {
        if (__is_dombasedxssfinder_string(this)) {
            const str = stringPrototypeSlice.apply(this.toString(), arguments);
            return new __dombasedxssfinder_String(str, this);
        }
        return stringPrototypeSlice.apply(this, arguments);
    };

    const stringPrototypeSmall = String.prototype.small;
    String.prototype.small = function () {
        if (__is_dombasedxssfinder_string(this)) {
            const str = stringPrototypeSmall.apply(this.toString(), arguments);
            return new __dombasedxssfinder_String(str, this);
        }
        return stringPrototypeSlice.apply(this, arguments);
    };

    const stringPrototypeSplit = String.prototype.split;
    String.prototype.split = function () {
        if (__is_dombasedxssfinder_string(this)) {
            const array = stringPrototypeSplit.apply(this.toString(), arguments);
            for (let i = 0; i < array.length; i++) {
                array[i] = new __dombasedxssfinder_String(array[i], this);
            }
            return array;
        }
        return stringPrototypeSplit.apply(this, arguments);
    };

    const stringPrototypeStartsWith = String.prototype.startsWith;
    String.prototype.startsWith = function () {
        return stringPrototypeStartsWith.apply(this.toString(), arguments);
    };

    const stringPrototypeStrike = String.prototype.strike;
    String.prototype.strike = function () {
        if (__is_dombasedxssfinder_string(this)) {
            const str = stringPrototypeStrike.apply(this.toString(), arguments);
            return new __dombasedxssfinder_String(str, this);
        }
        return stringPrototypeStrike.apply(this, arguments);
    };

    const stringPrototypeSub = String.prototype.sub;
    String.prototype.sub = function () {
        if (__is_dombasedxssfinder_string(this)) {
            const str = stringPrototypeSub.apply(this.toString(), arguments);
            return new __dombasedxssfinder_String(str, this);
        }
        return stringPrototypeSub.apply(this, arguments);
    };

    const stringPrototypeSubstr = String.prototype.substr;
    String.prototype.substr = function () {
        if (__is_dombasedxssfinder_string(this)) {
            const str = stringPrototypeSubstr.apply(this.toString(), arguments);
            return new __dombasedxssfinder_String(str, this);
        }
        return stringPrototypeSubstr.apply(this, arguments);
    };

    const stringPrototypeSubstring = String.prototype.substring;
    String.prototype.substring = function () {
        if (__is_dombasedxssfinder_string(this)) {
            const str = stringPrototypeSubstring.apply(this.toString(), arguments);
            return new __dombasedxssfinder_String(str, this);
        }
        return stringPrototypeSubstring.apply(this, arguments);
    };

    const stringPrototypeSup = String.prototype.sup;
    String.prototype.sup = function () {
        if (__is_dombasedxssfinder_string(this)) {
            const str = stringPrototypeSup.apply(this.toString(), arguments);
            return new __dombasedxssfinder_String(str, this);
        }
        return stringPrototypeSup.apply(this, arguments);
    };

    const stringPrototypeToLocaleLowerCase = String.prototype.toLocaleLowerCase;
    String.prototype.toLocaleLowerCase = function () {
        if (__is_dombasedxssfinder_string(this)) {
            const str = stringPrototypeToLocaleLowerCase.apply(this.toString(), arguments);
            return new __dombasedxssfinder_String(str, this);
        }
        return stringPrototypeToLocaleLowerCase.apply(this, arguments);
    };

    const stringPrototypeToLocaleUpperCase = String.prototype.toLocaleUpperCase;
    String.prototype.toLocaleUpperCase = function () {
        if (__is_dombasedxssfinder_string(this)) {
            const str = stringPrototypeToLocaleUpperCase.apply(this.toString(), arguments);
            return new __dombasedxssfinder_String(str, this);
        }
        return stringPrototypeToLocaleUpperCase.apply(this, arguments);
    };

    const stringPrototypeToLowerCase = String.prototype.toLowerCase;
    String.prototype.toLowerCase = function () {
        if (__is_dombasedxssfinder_string(this)) {
            const str = stringPrototypeToLowerCase.apply(this.toString(), arguments);
            return new __dombasedxssfinder_String(str, this);
        }
        return stringPrototypeToLowerCase.apply(this, arguments);
    };

    // skip String.prototype.toString, which is overwritten in __dombasedxssfinder_String

    const stringPrototypeToUpperCase = String.prototype.toUpperCase;
    String.prototype.toUpperCase = function () {
        if (__is_dombasedxssfinder_string(this)) {
            const str = stringPrototypeToUpperCase.apply(this.toString(), arguments);
            return new __dombasedxssfinder_String(str, this);
        }
        return stringPrototypeToUpperCase.apply(this, arguments);
    };

    const stringPrototypeTrim = String.prototype.trim;
    String.prototype.trim = function () {
        if (__is_dombasedxssfinder_string(this)) {
            const str = stringPrototypeTrim.apply(this.toString(), arguments);
            return new __dombasedxssfinder_String(str, this);
        }
        return stringPrototypeTrim.apply(this, arguments);
    };

    const stringPrototypeTrimEnd = String.prototype.trimEnd;
    String.prototype.trimEnd = function () {
        if (__is_dombasedxssfinder_string(this)) {
            const str = stringPrototypeTrimEnd.apply(this.toString(), arguments);
            return new __dombasedxssfinder_String(str, this);
        }
        return stringPrototypeTrimEnd.apply(this, arguments);
    };

    const stringPrototypeTrimStart = String.prototype.trimStart;
    String.prototype.trimStart = function () {
        if (__is_dombasedxssfinder_string(this)) {
            const str = stringPrototypeTrimStart.apply(this.toString(), arguments);
            return new __dombasedxssfinder_String(str, this);
        }
        return stringPrototypeTrimStart.apply(this, arguments);
    };

    // skip String.prototype.valueOf, which is overwritten in __dombasedxssfinder_String

    ///////////////////////////////////////////////
    // RegExp.prototype
    ///////////////////////////////////////////////

    const regExpPrototypeExec = RegExp.prototype.exec;
    RegExp.prototype.exec = function () {
        const array = regExpPrototypeExec.apply(this, arguments);
        if (array !== null && __is_dombasedxssfinder_string(arguments[0])) {
            for (let i = 0; i < array.length; i++) {
                array[i] = new __dombasedxssfinder_String(array[i], arguments[0]);
            }
        }
        return array;
    };

    ///////////////////////////////////////////////
    // global functions
    ///////////////////////////////////////////////

    const _decodeURI = decodeURI;
    decodeURI = function (encodedURI) {
        encodedURI = __convert_to_dombasedxssfinder_string_if_location(encodedURI);
        if (__is_dombasedxssfinder_string(encodedURI)) {
            const str = _decodeURI.apply(this, [encodedURI.toString()]);
            const newStr = new __dombasedxssfinder_String(str, encodedURI);
            return newStr;
        }
        return _decodeURI.apply(this, arguments);
    };

    const _encodeURI = encodeURI;
    encodeURI = function (URI) {
        URI = __convert_to_dombasedxssfinder_string_if_location(URI);
        if (__is_dombasedxssfinder_string(URI)) {
            const str = _encodeURI.apply(this, [URI.toString()]);
            const newStr = new __dombasedxssfinder_String(str, URI);
            return newStr;
        }
        return _encodeURI.apply(this, arguments);
    };

    const _decodeURIComponent = decodeURIComponent;
    decodeURIComponent = function (encodedURI) {
        encodedURI = __convert_to_dombasedxssfinder_string_if_location(encodedURI);
        if (__is_dombasedxssfinder_string(encodedURI)) {
            const str = _decodeURIComponent.apply(this, [encodedURI.toString()]);
            const newStr = new __dombasedxssfinder_String(str, encodedURI);
            return newStr;
        }
        return _decodeURIComponent.apply(this, arguments);
    };

    const _encodeURIComponent = encodeURIComponent;
    encodeURIComponent = function (URI) {
        URI = __convert_to_dombasedxssfinder_string_if_location(URI);
        if (__is_dombasedxssfinder_string(URI)) {
            const str = _encodeURIComponent.apply(this, [URI.toString()]);
            const newStr = new __dombasedxssfinder_String(str, URI);
            return newStr;
        }
        return _encodeURIComponent.apply(this, arguments);
    };

    const _unescape = unescape;
    unescape = function (escapedString) {
        escapedString = __convert_to_dombasedxssfinder_string_if_location(escapedString);
        if (__is_dombasedxssfinder_string(escapedString)) {
            const str = _unescape.apply(this, [escapedString.toString()]);
            const newStr = new __dombasedxssfinder_String(str, escapedString);
            return newStr;
        }
        return _unescape.apply(this, arguments);
    };

    const _escape = escape;
    escape = function (string) {
        string = __convert_to_dombasedxssfinder_string_if_location(string);
        if (__is_dombasedxssfinder_string(string)) {
            const str = _escape.apply(this, [string.toString()]);
            const newStr = new __dombasedxssfinder_String(str, string);
            return newStr;
        }
        return _escape.apply(this, arguments);
    };

    const _postMessage = postMessage;
    postMessage = function (message) {
        if (__is_dombasedxssfinder_string(message)) {
            arguments[0] = message.toString();
        }
        return _postMessage.apply(this, arguments);
    };

})();


/*********************************************/
// sinks
/*********************************************/


(function () {
    ///////////////////////////////////////////////
    // Range.prototype
    ///////////////////////////////////////////////

    const rangeCreateContextualFragment = Range.prototype.createContextualFragment;
    Range.prototype.createContextualFragment = function (fragment) {
        if (__is_dombasedxssfinder_string_html(fragment)) {
            __dombasedxssfinder_vulns_push(fragment.sources, 'Range.prototype.createContextualFragment()');
        }
        return rangeCreateContextualFragment.apply(this, arguments);
    };

    ///////////////////////////////////////////////
    // document
    ///////////////////////////////////////////////

    const documentWrite = document.write;
    document.write = function (...text) {
        for (let i = 0; i < text.length; i++) {
            if (__is_dombasedxssfinder_string_html(text[i])) {
                __dombasedxssfinder_vulns_push(text[i].sources, 'document.write()');
            }
        }
        return documentWrite.apply(this, arguments);
    };

    const documentWriteln = document.writeln;
    document.writeln = function (...text) {
        for (let i = 0; i < text.length; i++) {
            if (__is_dombasedxssfinder_string_html(text[i])) {
                __dombasedxssfinder_vulns_push(text[i].sources, 'document.writeln()');
            }
        }
        return documentWriteln.apply(this, arguments);
    };

    ///////////////////////////////////////////////
    // global functions
    ///////////////////////////////////////////////

    const _eval = eval;
    eval = function (x) {
        if (__is_dombasedxssfinder_string_script(x)) {
            __dombasedxssfinder_vulns_push(x.sources, 'eval()');
            // eval requires toString()
            return _eval.apply(this, [x.toString()]);
        }
        return _eval.apply(this, arguments);
    };

    const _setInterval = setInterval;
    setInterval = function (handler) {
        if (__is_dombasedxssfinder_string_script(handler)) {
            __dombasedxssfinder_vulns_push(handler.sources, 'setTimeout()');
        }
        return _setInterval.apply(this, arguments);
    };

    const _setTimeout = setTimeout;
    setTimeout = function (handler) {
        if (__is_dombasedxssfinder_string_script(handler)) {
            __dombasedxssfinder_vulns_push(handler.sources, 'setTimeout()');
        }
        return _setTimeout.apply(this, arguments);
    };

})();


// @asthook: +,+=
function __dombasedxssfinder_plus(left, right) {
    left = __convert_to_dombasedxssfinder_string_if_location(left);
    right = __convert_to_dombasedxssfinder_string_if_location(right);
    if (__is_dombasedxssfinder_string(left) || __is_dombasedxssfinder_string(right)) {
        const sources = [];
        if (__is_dombasedxssfinder_string(left)) {
            left.sources.forEach(e => sources.push(e));
        }
        if (__is_dombasedxssfinder_string(right)) {
            right.sources.forEach(e => sources.push(e));
        }
        return new __dombasedxssfinder_String('' + left + right, { sources });
    }
    try {
        return left + right;
    } catch (e) {
        return left.toString() + right.toString();
    }
}

// @asthook: object.key || object[key]
function __dombasedxssfinder_get(object, key) {
    if (object === window.location) {
        if (key === 'hash') {
            return new __dombasedxssfinder_String(object[key], {
                sources: [__dombasedxssfinder_get_source('window.location.hash')],
            });
        } else if (key === 'href') {
            return new __dombasedxssfinder_String(object[key], {
                sources: [__dombasedxssfinder_get_source('window.location.href')],
            });
        } else if (key === 'pathname') {
            return new __dombasedxssfinder_String(object[key], {
                sources: [__dombasedxssfinder_get_source('window.location.pathname')],
            });
        } else if (key === 'search') {
            return new __dombasedxssfinder_String(object[key], {
                sources: [__dombasedxssfinder_get_source('window.location.search')],
            });
        }
    } else if (object === document) {
        if (key === 'documentURI') {
            return new __dombasedxssfinder_String(object[key], {
                sources: [__dombasedxssfinder_get_source('document.documentURI')],
            });
        } else if (key === 'baseURI') {
            return new __dombasedxssfinder_String(object[key], {
                sources: [__dombasedxssfinder_get_source('document.baseURI')],
            });
        } else if (key === 'URL') {
            return new __dombasedxssfinder_String(object[key], {
                sources: [__dombasedxssfinder_get_source('document.URL')],
            });
        } else if (key === 'referrer' && object[key]) {
            return new __dombasedxssfinder_String(object[key], {
                sources: [__dombasedxssfinder_get_source('document.referrer')],
            });
        }
    }
    return object[key];
}

// @asthook: object.key = value || object[key] = value
function __dombasedxssfinder_put(object, key, value) {
    if (object[key] === window.location && __is_dombasedxssfinder_string_script(value)) {
        // __dombasedxssfinder_vulns_push(value.sources, 'window.location');
        // kill navigation
        return;
    } else if (object === window.location && key === 'href' && __is_dombasedxssfinder_string_script(value) && value.toString() !== object[key]) {
        // __dombasedxssfinder_vulns_push(value.sources, 'window.location.href');
        // kill navigation
        return;
    } else if (object instanceof Element && key === 'innerHTML' && __is_dombasedxssfinder_string_html(value)) {
        __dombasedxssfinder_vulns_push(value.sources, 'Element.innerHTML');
    } else if (object instanceof Element && key === 'outerHTML' && __is_dombasedxssfinder_string_html(value)) {
        __dombasedxssfinder_vulns_push(value.sources, 'Element.outerHTML');
    } else if (object instanceof HTMLScriptElement && key === 'src' && __is_dombasedxssfinder_string_url(value)) {
        __dombasedxssfinder_vulns_push(value.sources, 'HTMLScriptElement.src');
    } else if (object instanceof HTMLEmbedElement && key === 'src' && __is_dombasedxssfinder_string_url(value)) {
        __dombasedxssfinder_vulns_push(value.sources, 'HTMLEmbedElement.src');
    } else if (object instanceof HTMLIFrameElement && key === 'src' && __is_dombasedxssfinder_string_script(value)) {
        __dombasedxssfinder_vulns_push(value.sources, 'HTMLIFrameElement.src');
    } else if (object instanceof HTMLAnchorElement && key === 'href' && __is_dombasedxssfinder_string_script(value)) {
        __dombasedxssfinder_vulns_push(value.sources, 'HTMLAnchorElement.href');
    } else if (object instanceof HTMLFormElement && key === 'action' && __is_dombasedxssfinder_string_script(value)) {
        __dombasedxssfinder_vulns_push(value.sources, 'HTMLFormElement.action');
    } else if (object instanceof HTMLInputElement && key === 'formAction' && __is_dombasedxssfinder_string_script(value)) {
        __dombasedxssfinder_vulns_push(value.sources, 'HTMLInputElement.formAction');
    } else if (object instanceof HTMLButtonElement && key === 'formAction' && __is_dombasedxssfinder_string_script(value)) {
        __dombasedxssfinder_vulns_push(value.sources, 'HTMLButtonElement.formAction');
    } else if (object instanceof HTMLObjectElement && key === 'data' && __is_dombasedxssfinder_string_data_html(value)) {
        __dombasedxssfinder_vulns_push(value.sources, 'HTMLObjectElement.data');
    } else if (object instanceof HTMLScriptElement && key === 'text' && __is_dombasedxssfinder_string_script(value)) {
        __dombasedxssfinder_vulns_push(value.sources, 'HTMLScriptElement.text');
    } else if (object instanceof HTMLScriptElement && key === 'textContent' && __is_dombasedxssfinder_string_script(value)) {
        __dombasedxssfinder_vulns_push(value.sources, 'HTMLScriptElement.textContent');
    } else if (object instanceof HTMLScriptElement && key === 'innerText' && __is_dombasedxssfinder_string_script(value)) {
        __dombasedxssfinder_vulns_push(value.sources, 'HTMLScriptElement.innerText');
    }
    return object[key] = value;
}

// @asthook: function
function __dombasedxssfinder_new_Function() {
    const f = new Function(...arguments);
    if (__is_dombasedxssfinder_string_script(arguments[arguments.length - 1])) {
        __dombasedxssfinder_vulns_push(arguments[arguments.length - 1].sources, 'new Function()');
        f.__dombasedxssfinder_str = arguments[arguments.length - 1];
    }
    return f;
}

// @asthook: ==
function __dombasedxssfinder_equal(left, right) {
    if (__is_dombasedxssfinder_string(left)) {
        left = left.toString();
    }
    if (__is_dombasedxssfinder_string(right)) {
        right = right.toString();
    }
    return left == right;
}

// @asthook: !=
function __dombasedxssfinder_notEqual(left, right) {
    if (__is_dombasedxssfinder_string(left)) {
        left = left.toString();
    }
    if (__is_dombasedxssfinder_string(right)) {
        right = right.toString();
    }
    return left != right;
}

// @asthook: ===
function __dombasedxssfinder_strictEqual(left, right) {
    if (__is_dombasedxssfinder_string(left)) {
        left = left.toString();
    }
    if (__is_dombasedxssfinder_string(right)) {
        right = right.toString();
    }
    return left === right;
}

// @asthook: !==
function __dombasedxssfinder_strictNotEqual(left, right) {
    if (__is_dombasedxssfinder_string(left)) {
        left = left.toString();
    }
    if (__is_dombasedxssfinder_string(right)) {
        right = right.toString();
    }
    return left !== right;
}

// @asthook: typeof
function __dombasedxssfinder_typeof(o) {
    if (__is_dombasedxssfinder_string(o)) {
        return 'string';
    }
    return typeof o;
}

// @asthook: object.key(...arguments) || object[key](...arguments)
function __dombasedxssfinder_property_call(object, key, ...arguments) {
    if (object[key] === window.location.assign) {
        // cannot overwrite, replace it when called.
        return (function (url) {
            if (__is_dombasedxssfinder_string_script(url)) {
                // __dombasedxssfinder_vulns_push(url.sources, 'window.location.assign()');
                // kill navigation
                return;
            }
        }).apply(object, arguments);
    } else if (object[key] === window.location.replace) {
        // cannot overwrite, replace it when called.
        return (function (url) {
            if (__is_dombasedxssfinder_string_script(url)) {
                // __dombasedxssfinder_vulns_push(url.sources, 'window.location.replace()');
                // kill navigation
                return;
            }
        }).apply(object, arguments);
    } else if (object instanceof Element && key === 'setAttribute') {
        const elementSetAttribute = object[key];
        return (function (qualifiedName, value) {
            if (qualifiedName.startsWith('on') && __is_dombasedxssfinder_string_script(value)) {
                __dombasedxssfinder_vulns_push(value.sources, `Element.setAttribute('${qualifiedName}')`);
            } else if (this instanceof HTMLScriptElement && qualifiedName === 'src' && __is_dombasedxssfinder_string_url(value)) {
                __dombasedxssfinder_vulns_push(value.sources, 'HTMLScriptElement.setAttribute(\'src\')');
            } else if (this instanceof HTMLEmbedElement && qualifiedName === 'src' && __is_dombasedxssfinder_string_url(value)) {
                __dombasedxssfinder_vulns_push(value.sources, 'HTMLEmbedElement.setAttribute(\'src\')');
            } else if (this instanceof HTMLIFrameElement && qualifiedName === 'src' && __is_dombasedxssfinder_string_script(value)) {
                __dombasedxssfinder_vulns_push(value.sources, 'HTMLIFrameElement.setAttribute(\'src\')');
            } else if (this instanceof HTMLAnchorElement && qualifiedName === 'href' && __is_dombasedxssfinder_string_script(value)) {
                __dombasedxssfinder_vulns_push(value.sources, 'HTMLAnchorElement.setAttribute(\'href\')');
            } else if (this instanceof HTMLFormElement && qualifiedName === 'action' && __is_dombasedxssfinder_string_script(value)) {
                __dombasedxssfinder_vulns_push(value.sources, 'HTMLFormElement.setAttribute(\'action\')');
            } else if (this instanceof HTMLInputElement && qualifiedName === 'formaction' && __is_dombasedxssfinder_string_script(value)) {
                __dombasedxssfinder_vulns_push(value.sources, 'HTMLInputElement.setAttribute(\'formaction\')');
            } else if (this instanceof HTMLButtonElement && qualifiedName === 'formaction' && __is_dombasedxssfinder_string_script(value)) {
                __dombasedxssfinder_vulns_push(value.sources, 'HTMLButtonElement.setAttribute(\'formaction\')');
            } else if (this instanceof HTMLObjectElement && qualifiedName === 'data' && __is_dombasedxssfinder_string_data_html(value)) {
                __dombasedxssfinder_vulns_push(value.sources, 'HTMLObjectElement.setAttribute(\'data\')');
            }
            elementSetAttribute.apply(this, arguments);
        }).apply(object, arguments);
    } else if (object instanceof Element && key === 'addEventListener') {
        const elementAddEventListener = object[key];
        return (function (type, listener) {
            if (type === 'click' && listener && listener.__dombasedxssfinder_str && __is_dombasedxssfinder_string_script(listener.__dombasedxssfinder_str)) {
                __dombasedxssfinder_vulns_push(listener.__dombasedxssfinder_str.sources, 'Element.addEventListener(\'click\')');
            }
            elementAddEventListener.apply(this, arguments);
        }).apply(object, arguments);
    }

    return object[key](...arguments);
}

// @asthook: func(...arguments)
function __dombasedxssfinder_call(func, ...arguments) {
    if (func === window.location.assign) {
        // cannot overwrite, replace it when called.
        func = function (url) {
            if (__is_dombasedxssfinder_string_script(url)) {
                // __dombasedxssfinder_vulns_push(url.sources, 'window.location.assign()');
                // kill navigation
                return;
            }
        };
    } else if (func === window.location.replace) {
        // cannot overwrite, replace it when called.
        func = function (url) {
            if (__is_dombasedxssfinder_string_script(url)) {
                // __dombasedxssfinder_vulns_push(url.sources, 'window.location.replace()');
                // kill navigation
                return;
            }
        };
    }

    return func(...arguments);
}