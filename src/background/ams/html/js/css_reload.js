!function () {
    function e(e) {
        t(e), window.PrefixFree && StyleFix.process()
    }

    function t(e) {
        var t = n(), a = document.createElement("style");
        a.type = "text/css", a.className = "cp-pen-styles", a.styleSheet ? a.styleSheet.cssText = e : a.appendChild(document.createTextNode(e)), c.appendChild(a), t && t.parentNode.removeChild(t)
    }

    function n() {
        for (var e = document.getElementsByTagName("style"), t = e.length - 1; t >= 0; t--)if ("cp-pen-styles" === e[t].className)return e[t];
        return !1
    }

    function a(e) {
        window.addEventListener ? window.addEventListener("message", e, !1) : window.attachEvent("onmessage", e)
    }

    function s(e, t) {
        try {
            if (!/codepen/.test(e.origin))return null;
            if ("object" != typeof e.data)return null;
            if (e.data.action === t)return e.data
        } catch (n) {
        }
        return null
    }

    var c = document.head || document.getElementsByTagName("head")[0], r = "ACTION_LIVE_VIEW_RELOAD_CSS";
    a(function (t) {
        var n = s(t, r);
        n && e(n.data.css)
    })
}();