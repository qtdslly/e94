function SlideBox(e, j, v) {
    var v = v || {};
    var e = e;
    var j = j;
    var f = v.swiftCss ? v.swiftCss : "";
    var n = v.swiftText ? v.swiftText : "";
    var i = v.autoLeft ? v.autoLeft : false;
    var c = v.autoTop ? v.autoTop : false;
    var d = v.trigger ? v.trigger : "mouseover";
    var k = v.hoverDelay ? v.hoverDelay : 200;
    var g = v.outDelay ? v.outDelay : 0;
    var a = v.showEffect ? v.showEffect : "slideDown";
    var u = v.hideEffect ? v.hideEffect : "fadeOut";
    var r = v.showTime ? v.showTime : 100;
    var q = v.hideTime ? v.hideTime : 100;
    if (d == "mouseover") {
        $(e).hoverDelay({
            hoverDuring: k, outDuring: g, hoverEvent: function () {
                if (i) {
                    var x = ($(e).outerWidth() - $(j).outerWidth()) / 2;
                    $(j).css("left", x)
                }
                if (c) {
                    var y = ($(e).outerHeight() - $(j).outerHeight()) / 2;
                    $(j).css("top", y)
                }
                if (a == "slideDown") {
                    $(j).slideDown(r)
                }
                if (a == "show") {
                    $(j).show(r)
                }
                if (a == "fadeIn") {
                    $(j).fadeIn(r)
                }
            }, outEvent: function () {
                if (u == "slideUp") {
                    $(j).slideUp(q)
                }
                if (u == "hide") {
                    $(j).hide(q)
                }
                if (u == "fadeOut") {
                    $(j).fadeOut(q)
                }
            }
        })
    }
    if (d == "click") {
        $(e).bind("click", (function () {
            if (a == "slideDown") {
                $(j).slideToggle(r)
            }
            if (a == "show") {
                $(j).toggle(r)
            }
            if (a == "fadeIn") {
                $(j).fadeToggle(r)
            }
            if (f) {
                $(e).toggleClass(f)
            }
            if (n.indexOf("$$") != -1) {
                var x = n.split("$$");
                if ($(e).text() == x[0]) {
                    $(e).text(x[1])
                } else {
                    $(e).text(x[0])
                }
            }
        }))
    }
}
function IsMobile() {
    var d = navigator.userAgent.toLowerCase();
    var i = d.match(/ipad/i) == "ipad";
    var j = d.match(/iphone os/i) == "iphone os";
    var g = d.match(/midp/i) == "midp";
    var e = d.match(/rv:1.2.3.4/i) == "rv:1.2.3.4";
    var f = d.match(/ucweb/i) == "ucweb";
    var a = d.match(/android/i) == "android";
    var c = d.match(/windows ce/i) == "windows ce";
    var k = d.match(/windows mobile/i) == "windows mobile";
    if (i || j || g || e || f || a || c || k) {
        return true
    } else {
        return false
    }
}
function bindScroll(a) {
    a.jscroll({
        W: "15px",
        BgUrl: "url(/statics/img/ui/s_bg.gif)",
        Bg: "right 0 repeat-y",
        Bar: {
            Bd: {Out: "#a3c3d5", Hover: "#b7d5e6"},
            Bg: {Out: "-45px 0 repeat-y", Hover: "-58px 0 repeat-y", Focus: "-71px 0 repeat-y"}
        },
        Btn: {
            btn: true,
            uBg: {Out: "0 0", Hover: "-15px 0", Focus: "-30px 0"},
            dBg: {Out: "0 -15px", Hover: "-15px -15px", Focus: "-30px -15px"}
        },
        Fn: function () {
        }
    })
}
function thisScroll(d, e) {
    var f = $(d);
    var a = false;
    var c = f.innerHeight();
    if (c > e) {
        a = true
    }
    if (IsMobile()) {
        a = false
    }
    if (a) {
        bindScroll(f)
    }
}
function getScrollWidth() {
    var c, d, a = document.createElement("DIV");
    a.style.cssText = "position:absolute; top:-1000px; width:100px; height:100px; overflow:hidden;";
    c = document.body.appendChild(a).clientWidth;
    a.style.overflowY = "scroll";
    d = a.clientWidth;
    document.body.removeChild(a);
    return c - d
}
function JuaInputOn(a) {
    $(a).removeClass("hold");
    if (a.type == "password") {
        if ($(a).prev(".onlyshowpass")) {
            $(a).prev(".onlyshowpass").hide()
        }
    } else {
        if ($(a).val() == $(a).attr("defaultV")) {
            $(a).val("")
        }
    }
}
function JuaInputOut(a) {
    $(a).addClass("hold");
    if (a.type == "password") {
        if ($(a).prev(".onlyshowpass") && $(a).val() == "") {
            $(a).hide();
            $(a).prev(".onlyshowpass").show().blur()
        }
    } else {
        if ($(a).val() == "") {
            $(a).val($(a).attr("defaultV"))
        }
    }
}
function JuaInputSwt(a) {
    if ($(a).next(":password")) {
        $(a).hide();
        $(a).next(":password").show().focus()
    } else {
        return
    }
}
function reLoading(c, a) {
    var c = arguments[0] ? arguments[0] : "#reloadbar";
    var a = arguments[1] ? arguments[1] : 5;
    if (a != 5) {
        $(".fullwidth .expand").css("animation", "fullexpand " + a + "s ease");
        $(".fullwidth .expand").css("-moz-animation", "fullexpand " + a + "s ease");
        $(".fullwidth .expand").css("-webkit-animation", "fullexpand " + a + "s ease")
    }
    $(c).removeClass("fullwidth").delay(a).queue(function (d) {
        $(this).addClass("fullwidth");
        d()
    })
}
var JuaSec;
var JuaTimeOut = false;
var JuaBoxNum = 0;
var JuaBoxIndex;
var JuaCallBack = new Array();
var JuaIsFrame = (window.top == window.self) ? false : true;
var JuaBox = {
    info: function (g, e) {
        var e = e || {};
        var g = g ? g : "";
        var C = e.title ? e.title : chbtc.L("币网提示您");
        var q = e.btnNum ? e.btnNum : 1;
        var c = e.btnFuc1 ? e.btnFuc1 : "JuaBox.close();";
        var a = e.btnFuc2 ? e.btnFuc2 : "JuaBox.close();";
        var v = e.btnName1 ? e.btnName1 : chbtc.L("确定");
        var u = e.btnName2 ? e.btnName2 : chbtc.L("取消");
        var d = e.numSec ? e.numSec : 5000;
        var z = (e.isKill === false) ? false : true;
        if (z) {
            JuaBox.close()
        }
        JuaTimeOut = e.isSec ? e.isSec : JuaTimeOut;
        JuaCallBack[JuaBoxNum] = "";
        JuaCallBack[JuaBoxNum] += e.callback ? e.callback : "";
        JuaBoxNum += 1;
        JuaBoxIndex = JuaBoxNum - 1;
        var i = "btn-red";
        var f = "btn-gray";
        if (v.indexOf("$$") != -1) {
            i = v.split("$$")[1];
            v = v.split("$$")[0]
        }
        if (u.indexOf("$$") != -1) {
            f = u.split("$$")[1];
            u = u.split("$$")[0]
        }
        var k = "";
        k += '<div class="ch_mask hide" onClick="JuaBox.close()"></div>';
        k += '<div class="ch_info hide">';
        if (C != "none") {
            k += '  <div class="userBox-header"><h2>' + C + "</h2></div>"
        }
        k += '  <div class="userBox-body">';
        k += '  <p class="yzinfo">' + g + "</p>";
        if (q != 0) {
            k += '  <div class="form-button">';
            if (q > 1) {
                k += '<a href="javascript:Void()" onClick="' + a + '" id="b2" class="btn ' + f + '">' + u + "</a>"
            }
            if (q > 0) {
                k += '<a href="javascript:Void()" onClick="' + c + '" id="b1" class="btn ' + i + '">' + v + "</a>"
            }
            k += "  </div>"
        }
        k += "  </div>";
        k += '  <div class="close" onClick="JuaBox.close()"></div>';
        k += "</div>";
        if ($("body").length != 0) {
            $("body").append(k)
        }
        var r = $(document).width();
        var A = $(window).height();
        var B = e.width ? e.width : $(".ch_info").eq(JuaBoxIndex).outerWidth();
        var j = e.height ? e.height : $(".ch_info").eq(JuaBoxIndex).outerHeight();
        var n = e.height ? e.height : "auto";
        var y = A / 2 - j / 2 + JuaBox.getScrollTop() - 30;
        var x = r / 2 - B / 2 - 30;
        if (y < 50) {
            y = 50;
            A = j + 100;
            if (A > parent.$(".JuaFrame").eq(-1).height()) {
                parent.$(".JuaFrame").eq(-1).height(A)
            }
        }
        $(".ch_mask").eq(JuaBoxIndex).css("height", $(document).height());
        $(".ch_mask").eq(JuaBoxIndex).css("width", $(document).width());
        $(".ch_info").eq(JuaBoxIndex).css("width", B);
        $(".ch_info").eq(JuaBoxIndex).css("height", "auto");
        $(".ch_info").eq(JuaBoxIndex).find(".info").css("height", n);
        $(".ch_info").eq(JuaBoxIndex).css("left", x);
        $(".ch_info").eq(JuaBoxIndex).css("top", y);
        $(".ch_info").eq(JuaBoxIndex).css("z-index", parseInt($(".ch_info").eq(JuaBoxIndex).css("z-index")) + JuaBoxNum);
        $(".ch_mask").eq(JuaBoxIndex).css("z-index", parseInt($(".ch_mask").eq(JuaBoxIndex).css("z-index")) + JuaBoxNum);
        $(".ch_mask").eq(JuaBoxIndex).removeClass("hide");
        $(".ch_info").eq(JuaBoxIndex).fadeIn(400);
        if (JuaTimeOut) {
            JuaSec = setTimeout("JuaBox.close();", d)
        }
    }, frame: function (e, f) {
        if (JuaIsFrame) {
            parent.JuaBox.frame(e, f);
            return
        }
        var f = f || {};
        var e = e ? e : "";
        if (e.indexOf("?") < 0) {
            e += "?a=1"
        }
        if (e.indexOf("iframe=1") < 0) {
            e += "&iframe=1"
        }
        var E = f.title ? f.title : "none";
        var r = f.btnNum ? f.btnNum : 0;
        var c = f.btnFuc1 ? f.btnFuc1 : "JuaBox.close();";
        var a = f.btnFuc2 ? f.btnFuc2 : "JuaBox.close();";
        var x = f.btnName1 ? f.btnName1 : chbtc.L("确定");
        var v = f.btnName2 ? f.btnName2 : chbtc.L("取消");
        var d = f.numSec ? f.numSec : 5000;
        var A = f.needLogin || false;
        var B = f.isKill ? f.isKill : false;
        if (B) {
            JuaBox.close()
        }
        JuaTimeOut = f.isSec ? f.isSec : JuaTimeOut;
        JuaCallBack[JuaBoxNum] = "";
        JuaCallBack[JuaBoxNum] += f.callback ? f.callback : "";
        JuaBoxNum += 1;
        JuaBoxIndex = JuaBoxNum - 1;
        var i = "orange";
        var g = "green";
        if (x.indexOf("$$") != -1) {
            i = x.split("$$")[1];
            x = x.split("$$")[0]
        }
        if (v.indexOf("$$") != -1) {
            g = v.split("$$")[1];
            v = v.split("$$")[0]
        }
        if (A && !chbtc.user.isLogin()) {
            e = "/user/log?iframe=1"
        }
        var k = "";
        k += '<div class="ch_mask hide" onClick="JuaBox.close()"></div>';
        k += '<div class="ch_info hide">';
        if (E != "none") {
            k += '  <div class="head"><h3>' + E + "</h3></div>"
        }
        k += '  <div class="frame JuaPage">';
        k += "  <iframe src='" + e + "' id='JuaFrame' class='JuaFrame' name='JuaFrame' allowtransparency='true' frameBorder='0' width='100%' hspace='0' scrolling='no' onLoad='SetWinHeight(this)'></iframe>";
        k += "  </div>";
        if (r != 0) {
            k += '  <div class="done">';
            if (r > 1) {
                k += '  <a href="javascript:Void()" onClick="' + a + '" id="b2" class="btn ' + g + '">' + v + "</a>"
            }
            if (r > 0) {
                k += '  <a href="javascript:Void()" onClick="' + c + '" id="b1" class="btn ' + i + '">' + x + "</a>"
            }
            k += "  </div>"
        }
        k += '  <div class="close" onClick="JuaBox.close()"></div>';
        k += "</div>";
        if ($("body").length != 0) {
            $("body").append(k)
        }
        var u = $(document).width();
        var C = $(window).height();
        var n = parseInt($(".ch_info").eq(JuaBoxIndex).outerHeight()) + 415;
        var D = f.width ? f.width : 500;
        var j = f.height ? f.height : n;
        var q = f.height ? parseInt(parseInt(f.height) + 5) : "auto";
        var z = C / 2 - j / 2 + JuaBox.getScrollTop() - 50;
        if (z < 0) {
            z = 50
        }
        var y = u / 2 - D / 2 - 50;
        $(".ch_mask").eq(JuaBoxIndex).css("height", $(document).height());
        $(".ch_mask").eq(JuaBoxIndex).css("width", $(document).width());
        $(".ch_info").eq(JuaBoxIndex).css("width", D);
        $(".ch_info").eq(JuaBoxIndex).css("height", "auto");
        $(".ch_info").eq(JuaBoxIndex).find(".info").css("height", q);
        $(".ch_info").eq(JuaBoxIndex).css("left", y);
        $(".ch_info").eq(JuaBoxIndex).css("top", z);
        $(".ch_info").eq(JuaBoxIndex).css("z-index", parseInt($(".ch_info").eq(JuaBoxIndex).css("z-index")) + JuaBoxNum);
        $(".ch_mask").eq(JuaBoxIndex).css("z-index", parseInt($(".ch_mask").eq(JuaBoxIndex).css("z-index")) + JuaBoxNum);
        $(".ch_mask").eq(JuaBoxIndex).removeClass("hide");
        $(".ch_info").eq(JuaBoxIndex).fadeIn(400);
        if (JuaTimeOut) {
            JuaSec = setTimeout("JuaBox.close();", d)
        }
    }, close: function (ops) {
        if (JuaBoxNum != 0) {
            $(".ch_mask").eq(JuaBoxIndex).hide(0, function () {
                $(".ch_mask").eq(JuaBoxIndex).remove();
                $(".ch_info").eq(JuaBoxIndex).remove();
                if (JuaTimeOut) {
                    clearTimeout(JuaSec);
                    JuaTimeOut = false
                }
                JuaBoxNum -= 1;
                JuaBoxIndex = JuaBoxNum - 1;
                if (JuaCallBack[JuaBoxNum] != "") {
                    eval(JuaCallBack[JuaBoxNum]);
                    JuaCallBack[JuaBoxNum] = ""
                }
            })
        }
    }, closeAll: function () {
        if (JuaBoxNum != 0) {
            $(".ch_info").animate({top: "-=50px", opacity: "0"}, 300, function () {
                $(".ch_mask").fadeOut(10, function () {
                    $(".ch_mask").remove();
                    $(".ch_info").remove();
                    if (JuaTimeOut) {
                        clearTimeout(JuaSec);
                        JuaTimeOut = false
                    }
                    JuaBoxNum = 0;
                    JuaBoxIndex = JuaBoxNum - 1;
                    if (JuaCallBack[JuaBoxNum] != "") {
                        eval(JuaCallBack[JuaBoxNum]);
                        JuaCallBack[JuaBoxNum] = ""
                    }
                })
            })
        }
    }, reSetHeight: function (a) {
        var c = $("body").height() + 10;
        if (a) {
            parent.$(".JuaFrame").eq(-1).height(c);
            parent.$(".JuaPage.frame").eq(-1).height(c)
        } else {
            parent.$(".JuaFrame").eq(-1).height(c);
            parent.$(".JuaPage").eq(-1).height(c)
        }
    }, getScrollTop: function () {
        var a = 0;
        if (document.documentElement && document.documentElement.scrollTop) {
            a = document.documentElement.scrollTop
        } else {
            if (document.body) {
                a = document.body.scrollTop
            }
        }
        return a
    }
};
function SetWinHeight(c) {
    var a = c;
    if (document.getElementById) {
        if (a && !window.opera) {
            if (a.contentDocument && a.contentDocument.body.offsetHeight) {
                a.height = a.contentDocument.body.offsetHeight + 5;
                $(".ch_info").eq(JuaBoxIndex).find(".info").height(a.height);
                a.contentDocument.body.style.backgroundImage = "none";
                a.contentDocument.body.style.backgroundColor = "transparent"
            } else {
                if (a.Document && a.Document.body.scrollHeight) {
                    a.height = a.Document.body.scrollHeight + 5;
                    $(".ch_info").eq(JuaBoxIndex).find(".info").height(a.height);
                    a.Document.body.style.backgroundImage = "none";
                    a.Document.body.style.backgroundColor = "transparent"
                }
            }
        }
    }
}
function showLoad(a) {
    var a = arguments[0] ? arguments[0] : "提交请求中，请稍候...";
    showTips("<h3>" + a + "</h3><div class='ch_load'></div>", "1", "", "", "确定$$blue", "", "false")
}
function closeTips() {
    if ($(".ch_mask").length != 0) {
        $(".ch_info,.ch_mask").hide(0, function () {
            $(".ch_info,.ch_mask").remove()
        });
        clearTimeout(d_sec)
    }
}
function getScrollTop() {
    var a = 0;
    if (document.documentElement && document.documentElement.scrollTop) {
        a = document.documentElement.scrollTop
    } else {
        if (document.body) {
            a = document.body.scrollTop
        }
    }
    return a
}
function set_menu_hover(a) {
    $(".d_menu1 li").removeClass("on");
    $(".d_menu1 li.m" + a + "").addClass("on")
}
function voids(a) {
}
String.prototype.replaceAll = function (c, a) {
    var d = new RegExp(c.replace(/([\(\)\[\]\{\}\^\$\+\-\*\?\.\"\'\|\/\\])/g, "\\$1"), "ig");
    return this.replace(d, a)
};
Object.extend = function (a, d) {
    for (var c in d) {
        a[c] = d[c]
    }
    return a
};
Object.extend(Object, {
    inspect: function (a) {
        try {
            if (a == undefined) {
                return "undefined"
            }
            if (a == null) {
                return "null"
            }
            return a.inspect ? a.inspect() : a.toString()
        } catch (c) {
            if (c instanceof RangeError) {
                return "..."
            }
            throw c
        }
    }, keys: function (a) {
        var c = [];
        for (var d in a) {
            c.push(d)
        }
        return c
    }, values: function (c) {
        var a = [];
        for (var d in c) {
            a.push(c[d])
        }
        return a
    }, clone: function (a) {
        return Object.extend({}, a)
    }
});
function getX(a) {
    return a.offsetLeft + (a.offsetParent ? getX(a.offsetParent) : a.x ? a.x : 0)
}
function getY(a) {
    return (a.offsetParent ? a.offsetTop + getY(a.offsetParent) : a.y ? a.y : 0)
}
var Style = {
    getFinalStyle: function (c, a) {
        if (window.getComputedStyle) {
            return window.getComputedStyle(c, null)[a]
        } else {
            if (c.currentStyle) {
                return c.currentStyle[a]
            } else {
                return c.style[a]
            }
        }
    }, height: function (c) {
        if (this.getFinalStyle(c, "display") !== "none") {
            return c.offsetHeight || c.clientHeight
        } else {
            c.style.display = "block";
            var a = c.offsetHeight || c.clientHeight;
            c.style.display = "none";
            return a
        }
    }, width: function (c) {
        if (this.getFinalStyle(c, "display") !== "none") {
            return c.offsetWidth || c.clientWidth
        } else {
            c.style.display = "block";
            var a = c.offsetWidth || c.clientWidth;
            c.style.display = "none";
            return a
        }
    }
};
function numberID() {
    return Math.round(Math.random() * 10000) * Math.round(Math.random() * 10000)
}
function T$(a) {
    return document.getElementById(a)
}
var sysCloseCount = 0;
NetBox = function () {
    var u, z, D, i, y, d, k, x, g, C = 0;
    var F, c, E, e, r, A;
    var j = false;
    var q = false;
    var n = "";
    var v;
    return {
        show: function (L, O, J, I, f, N) {
            v = N;
            g = f;
            if (!C) {
                u = document.createElement("div");
                u.id = "tinybox";
                z = document.createElement("div");
                z.id = "tinymask";
                D = document.createElement("div");
                D.id = "tinycontent";
                document.body.appendChild(z);
                document.body.appendChild(u);
                u.appendChild(D);
                C = 1
            }
            n = L;
            var P = (NetPage.height() / 2) - (J / 2);
            P = P < 10 ? 10 : P;
            var H = 0;
            if (J < 100) {
                H = (P + NetPage.top() - 60)
            } else {
                H = (P + NetPage.top())
            }
            var K = (NetPage.width() / 2) - (O / 2);
            if (!j) {
                y = L;
                d = O;
                k = J;
                if (f) {
                    r = F = getY(f);
                    A = c = getX(f);
                    E = Style.width(f) + 12;
                    e = Style.height(f) + 12
                }
                u.style.backgroundImage = "none";
                u.innerHTML = "";
                if (f) {
                    u.style.width = (E - 12) + "px";
                    u.style.height = (e - 12) + "px";
                    u.style.top = (F - 6) + "px";
                    u.style.left = (c - 6) + "px"
                } else {
                    u.style.width = O + "px";
                    if (J > 99) {
                        u.style.height = J + "px";
                        u.style.top = (H - 6) + "px";
                        u.style.left = (K - 6) + "px"
                    } else {
                        u.style.height = "auto";
                        u.style.top = (H + 37) + "px";
                        u.style.left = (K) + "px"
                    }
                }
                this.mask();
                if (f) {
                    u.style.display = "block";
                    $("#tinybox").animate({left: K, top: H, width: O, height: J}, 150, "", function () {
                        u.innerHTML = n;
                        u.style.height = "auto"
                    })
                } else {
                    u.innerHTML = n;
                    $(u).show()
                }
            } else {
                u.style.backgroundImage = "none";
                u.style.display = "block";
                if (J < 100) {
                    H = H - 50
                }
                u.style.height = $(u).height() + "px";
                u.innerHTML = n;
                var G = O;
                var M = 0;
                if (J == 99) {
                    if ($(u).find("iframe:first").length > 0) {
                        M = $(u).find("iframe:first").height()
                    } else {
                        M = $(u).find("div:first").height()
                    }
                    this.resize(G, M, null, 99)
                } else {
                    M = J;
                    this.resize(G, M)
                }
                this.mask()
            }
            j = true;
            if (I) {
                sysCloseCount = 0;
                setTimeout(function () {
                    if (sysCloseCount == 0) {
                        $("#tinybox").animate({top: -300}, 300);
                        NetBox.hide()
                    }
                }, 1000 * I)
            }
        }, hide: function (f) {
            sysCloseCount = -1;
            j = false;
            $(z).fadeOut(50);
            if (g) {
                u.style.height = $(u).height() + "px";
                u.innerHTML = "";
                $("#tinybox").animate({
                    left: (c - 6),
                    top: (F) - 6,
                    width: (E - 12),
                    height: (e - 12)
                }, 200, "", function () {
                    u.style.display = "none"
                })
            } else {
                $(u).fadeOut(50)
            }
        }, resize: function (f, K, L, H) {
            var I = (NetPage.height() / 2) - (K / 2);
            I = I < 10 ? 10 : I;
            var J = (I + NetPage.top());
            var G = (NetPage.width() / 2) - (f / 2);
            if (d == f) {
                f = 0
            }
            if (k == K) {
                K = 0
            }
            if (f > 0 && K > 0) {
                $("#tinybox").animate({left: G, top: J, width: f, height: K}, 150, "", function () {
                    if (H == 99) {
                        u.style.height = "auto"
                    }
                    $("#fancybox-frame").css("height", "100%");
                    if (L != null) {
                        L
                    }
                })
            } else {
                if (f > 0) {
                    $("#tinybox").animate({left: G, width: f}, 150, "", function () {
                        if (H == 99) {
                            u.style.height = "auto"
                        }
                        $("#fancybox-frame").css("height", "100%");
                        if (L != null) {
                            L
                        }
                    })
                } else {
                    if (K > 0) {
                        $("#tinybox").animate({top: J, height: K}, 150, "", function () {
                            if (H == 99) {
                                u.style.height = "auto"
                            }
                            $("#fancybox-frame").css("height", K - 44);
                            if (L != null) {
                                L
                            }
                        })
                    } else {
                        if (H == 99) {
                            u.style.height = "auto"
                        }
                        if (L != null) {
                            L
                        }
                    }
                }
            }
        }, mask: function () {
            z.style.display = "block";
            z.style.height = NetPage.theight() + "px";
            z.style.width = NetPage.twidth() + "px";
            z.style.opacity = 0.4;
            z.style.filter = "alpha(opacity=40)";
            if (v) {
                z.style.cursor = "pointer";
                z.title = "双击关闭";
                z.ondblclick = function () {
                    NetBox.hide()
                }
            }
        }, pos: function () {
        }, size: function (J, L, I, O) {
            J = typeof J == "object" ? J : T$(J);
            clearInterval(J.si);
            var f = J.offsetWidth, M = J.offsetHeight, K = f - parseInt(J.style.width), N = M - parseInt(J.style.height);
            var G = f - K > L ? -1 : 1, H = (M - N > I) ? -1 : 1;
            J.si = setInterval(function () {
                NetBox.twsize(J, L, K, G, I, N, H, O)
            }, 10)
        }, twsize: function (J, L, K, G, I, N, H, O) {
            var f = J.offsetWidth - K, M = J.offsetHeight - N;
            if (f == L && M == I) {
                clearInterval(J.si);
                u.style.backgroundImage = "none";
                D.style.display = "block";
                u.innerHTML = n
            } else {
                if (f != L) {
                    J.style.width = f + (Math.ceil(Math.abs(L - f) / O) * G) + "px"
                }
                if (M != I) {
                    J.style.height = M + (Math.ceil(Math.abs(I - M) / O) * H) + "px"
                }
                this.pos();
                if (a == f && B == M) {
                    clearInterval(J.si)
                }
                a = f;
                B = M
            }
        }
    };
    var a = 0;
    var B = 0
}();
NetPage = function () {
    return {
        top: function () {
            return document.body.scrollTop || document.documentElement.scrollTop
        }, width: function () {
            return self.innerWidth || document.documentElement.clientWidth
        }, height: function () {
            return self.innerHeight || document.documentElement.clientHeight
        }, theight: function () {
            var f = document, a = f.body, c = f.documentElement;
            return Math.max(Math.max(a.scrollHeight, c.scrollHeight), Math.max(a.clientHeight, c.clientHeight))
        }, twidth: function () {
            var f = document, a = f.body, c = f.documentElement;
            return Math.max(Math.max(a.scrollWidth, c.scrollWidth), Math.max(a.clientWidth, c.clientWidth))
        }
    }
}();
$.fn.extend({
    jscroll: function (a) {
        return this.each(function () {
            a = a || {};
            a.Bar = a.Bar || {};
            a.Btn = a.Btn || {};
            a.Bar.Bg = a.Bar.Bg || {};
            a.Bar.Bd = a.Bar.Bd || {};
            a.Btn.uBg = a.Btn.uBg || {};
            a.Btn.dBg = a.Btn.dBg || {};
            var e = {
                W: "15px",
                BgUrl: "",
                Bg: "#efefef",
                Bar: {
                    Pos: "up",
                    Bd: {Out: "#b5b5b5", Hover: "#ccc"},
                    Bg: {Out: "#fff", Hover: "#fff", Focus: "orange"}
                },
                Btn: {
                    btn: true,
                    uBg: {Out: "#ccc", Hover: "#fff", Focus: "orange"},
                    dBg: {Out: "#ccc", Hover: "#fff", Focus: "orange"}
                },
                Fn: function () {
                }
            };
            a.W = a.W || e.W;
            a.BgUrl = a.BgUrl || e.BgUrl;
            a.Bg = a.Bg || e.Bg;
            a.Bar.Pos = a.Bar.Pos || e.Bar.Pos;
            a.Bar.Bd.Out = a.Bar.Bd.Out || e.Bar.Bd.Out;
            a.Bar.Bd.Hover = a.Bar.Bd.Hover || e.Bar.Bd.Hover;
            a.Bar.Bg.Out = a.Bar.Bg.Out || e.Bar.Bg.Out;
            a.Bar.Bg.Hover = a.Bar.Bg.Hover || e.Bar.Bg.Hover;
            a.Bar.Bg.Focus = a.Bar.Bg.Focus || e.Bar.Bg.Focus;
            a.Btn.btn = a.Btn.btn != undefined ? a.Btn.btn : e.Btn.btn;
            a.Btn.uBg.Out = a.Btn.uBg.Out || e.Btn.uBg.Out;
            a.Btn.uBg.Hover = a.Btn.uBg.Hover || e.Btn.uBg.Hover;
            a.Btn.uBg.Focus = a.Btn.uBg.Focus || e.Btn.uBg.Focus;
            a.Btn.dBg.Out = a.Btn.dBg.Out || e.Btn.dBg.Out;
            a.Btn.dBg.Hover = a.Btn.dBg.Hover || e.Btn.dBg.Hover;
            a.Btn.dBg.Focus = a.Btn.dBg.Focus || e.Btn.dBg.Focus;
            a.Fn = a.Fn || e.Fn;
            var f = this;
            var u, D = 0, F = 0;
            $(f).css({overflow: "hidden", position: "relative", padding: "0px"});
            var r = $(f).width(), B = $(f).height() - 1;
            var x = a.W ? parseInt(a.W) : 21;
            var C = r - x;
            var A = a.Btn.btn == true ? x : 0;
            if ($(f).children(".jscroll-c").height() == null) {
                $(f).wrapInner("<div class='jscroll-c' style='top:0px;z-index:8000;zoom:1;position:relative'></div>");
                $(f).children(".jscroll-c").prepend("<div style='height:0px;overflow:hidden'></div>");
                $(f).append("<div class='jscroll-e' unselectable='on' style=' height:100%;top:0px;right:1px;-moz-user-select:none;position:absolute;overflow:hidden;z-index:8002;'><div class='jscroll-u' style='position:absolute;top:0px;width:100%;left:0;background:blue;overflow:hidden'></div><div class='jscroll-h'  unselectable='on' style='background:green;position:absolute;left:0;-moz-user-select:none;border:1px solid'></div><div class='jscroll-d' style='position:absolute;bottom:0px;width:100%;left:0;background:blue;overflow:hidden'></div></div>")
            }
            var n = $(f).children(".jscroll-c");
            var j = $(f).children(".jscroll-e");
            var i = j.children(".jscroll-h");
            var c = j.children(".jscroll-u");
            var k = j.children(".jscroll-d");
            if ($.browser.msie) {
                document.execCommand("BackgroundImageCache", false, true)
            }
            n.css({"padding-right": x});
            j.css({width: x, background: a.Bg, "background-image": a.BgUrl});
            i.css({
                top: A,
                background: a.Bar.Bg.Out,
                "background-image": a.BgUrl,
                "border-color": a.Bar.Bd.Out,
                width: x - 2
            });
            c.css({height: A, background: a.Btn.uBg.Out, "background-image": a.BgUrl});
            k.css({height: A, background: a.Btn.dBg.Out, "background-image": a.BgUrl});
            i.hover(function () {
                if (F == 0) {
                    $(this).css({
                        background: a.Bar.Bg.Hover,
                        "background-image": a.BgUrl,
                        "border-color": a.Bar.Bd.Hover
                    })
                }
            }, function () {
                if (F == 0) {
                    $(this).css({background: a.Bar.Bg.Out, "background-image": a.BgUrl, "border-color": a.Bar.Bd.Out})
                }
            });
            c.hover(function () {
                if (F == 0) {
                    $(this).css({background: a.Btn.uBg.Hover, "background-image": a.BgUrl})
                }
            }, function () {
                if (F == 0) {
                    $(this).css({background: a.Btn.uBg.Out, "background-image": a.BgUrl})
                }
            });
            k.hover(function () {
                if (F == 0) {
                    $(this).css({background: a.Btn.dBg.Hover, "background-image": a.BgUrl})
                }
            }, function () {
                if (F == 0) {
                    $(this).css({background: a.Btn.dBg.Out, "background-image": a.BgUrl})
                }
            });
            var d = n.height();
            var E = (B - 2 * A) * B / d;
            if (E < 10) {
                E = 10
            }
            var g = E / 6;
            var q = 0, z = false;
            i.height(E);
            if (d <= B) {
                n.css({padding: 0});
                j.css({display: "none"})
            } else {
                z = true
            }
            if (a.Bar.Pos != "up") {
                q = B - E - A;
                y()
            }
            i.bind("mousedown", function (I) {
                a.Fn && a.Fn.call(f);
                F = 1;
                i.css({background: a.Bar.Bg.Focus, "background-image": a.BgUrl});
                var H = I.pageY, G = parseInt($(this).css("top"));
                $(document).mousemove(function (J) {
                    q = G + J.pageY - H;
                    y()
                });
                $(document).mouseup(function () {
                    F = 0;
                    i.css({background: a.Bar.Bg.Out, "background-image": a.BgUrl, "border-color": a.Bar.Bd.Out});
                    $(document).unbind()
                });
                return false
            });
            c.bind("mousedown", function (G) {
                a.Fn && a.Fn.call(f);
                F = 1;
                c.css({background: a.Btn.uBg.Focus, "background-image": a.BgUrl});
                f.timeSetT("u");
                $(document).mouseup(function () {
                    F = 0;
                    c.css({background: a.Btn.uBg.Out, "background-image": a.BgUrl});
                    $(document).unbind();
                    clearTimeout(u);
                    D = 0
                });
                return false
            });
            k.bind("mousedown", function (G) {
                a.Fn && a.Fn.call(f);
                F = 1;
                k.css({background: a.Btn.dBg.Focus, "background-image": a.BgUrl});
                f.timeSetT("d");
                $(document).mouseup(function () {
                    F = 0;
                    k.css({background: a.Btn.dBg.Out, "background-image": a.BgUrl});
                    $(document).unbind();
                    clearTimeout(u);
                    D = 0
                });
                return false
            });
            f.timeSetT = function (I) {
                var G = this;
                if (I == "u") {
                    q -= g
                } else {
                    q += g
                }
                y();
                D += 2;
                var H = 500 - D * 50;
                if (H <= 0) {
                    H = 0
                }
                u = setTimeout(function () {
                    G.timeSetT(I)
                }, H)
            };
            j.bind("mousedown", function (G) {
                a.Fn && a.Fn.call(f);
                q = q + G.pageY - i.offset().top - E / 2;
                v();
                return false
            });
            function v() {
                if (q < A) {
                    q = A
                }
                if (q > B - E - A) {
                    q = B - E - A
                }
                i.stop().animate({top: q}, 100);
                var G = -((q - A) * d / (B - 2 * A));
                n.stop().animate({top: G}, 1000)
            }

            function y() {
                if (q < A) {
                    q = A
                }
                if (q > B - E - A) {
                    q = B - E - A
                }
                i.css({top: q});
                var G = -((q - A) * d / (B - 2 * A));
                n.css({top: G})
            }

            $(f).mousewheel(function () {
                if (z != true) {
                    return
                }
                a.Fn && a.Fn.call(f);
                if (this.D > 0) {
                    q -= g
                } else {
                    q += g
                }
                y()
            })
        })
    }
});
$.fn.extend({
    mousewheel: function (a) {
        return this.each(function () {
            var c = this;
            c.D = 0;
            if ($.browser.msie || $.browser.safari) {
                c.onmousewheel = function () {
                    c.D = event.wheelDelta;
                    event.returnValue = false;
                    a && a.call(c)
                }
            } else {
                c.addEventListener("DOMMouseScroll", function (d) {
                    c.D = d.detail > 0 ? -1 : 1;
                    d.preventDefault();
                    a && a.call(c)
                }, false)
            }
        })
    }
});
function Close() {
    NetBox.hide()
}
var jqTransformGetLabel = function (d) {
    var a = $(d.get(0).form);
    var e = d.next();
    if (!e.is("label")) {
        e = d.prev();
        if (e.is("label")) {
            var c = d.attr("id");
            if (c) {
                e = a.find('label[for="' + c + '"]')
            }
        }
    }
    if (e.is("label")) {
        return e.css("cursor", "pointer")
    }
    return false
};
function Ask(a) {
    var d = {CloseTime: 0, Msg: "？", Title: "？", Height: 99, callback: "Close()", callback2: "Close()", fromObj: null};
    a = Object.extend(d, a);
    if (a.Title != "？") {
        a.Msg = a.Title
    }
    var c = '<div class="AlertMessage"><div class="MessageTitle"><a class="xx" title="' + chbtc.L("关闭") + '" onclick="Close();" href="javascript:Void();">' + chbtc.L("关闭") + "</a>" + chbtc.L("BW提示您") + '</div><div class="MessageHelp"><div class="Message">' + a.Msg + '</div></div><div class="MessageControl"><div class="MessageControl2"><a onfocus="this.blur()" id="nobackButton_1"  class="noback" onclick="' + a.callback2 + '" href="javascript:Void();"><h4>' + chbtc.L("取消") + '</h4></a><a onfocus="this.blur()" id="okButton_1" class="ok" onclick="' + a.callback + '" href="javascript:Void();"><h4>' + chbtc.L("确定") + "</h4></a></div></div></div>";
    NetBox.show(c, 550, a.Height, a.CloseTime, a.fromObj);
    $("body").bind("keyup", function (e) {
        if (e.keyCode == "13") {
            $("#okButton_1").trigger("click")
        }
    })
}
function Ask2(a) {
    var d = {
        call: function (e) {
        }, data: "", CloseTime: 0, Msg: "？", Height: 99, Title: "？", callback: "", callback2: "Close()", fromObj: null
    };
    a = $.extend(d, a);
    if (a.Title != "？") {
        a.Msg = a.Title
    }
    var c = '<div class="AlertMessage"><div class="MessageTitle"><a class="xx" title="关闭" onclick="Close();" href="javascript:Void();">' + chbtc.L("关闭") + "</a>" + chbtc.L("BW提示您") + '</div><div class="MessageHelp"><div class="Message">' + a.Msg + '</div></div><div class="MessageControl"><div class="MessageControl2"><a onfocus="this.blur()" id="nobackButton_2"  class="noback" onclick="' + a.callback2 + '" href="javascript:Void();"><h4>' + chbtc.L("取消") + '</h4></a><a onfocus="this.blur()" id="okButton_2" class="ok" onclick="' + a.callback + '" href="javascript:Void();"><h4>' + chbtc.L("确定") + "</h4></a></div></div></div>";
    NetBox.show(c, 550, a.Height, a.CloseTime, a.fromObj);
    $("#okButton_2").bind("click", function () {
        a.call(a.data);
        return false
    });
    $("body").bind("keyup", function (e) {
        if (e.keyCode == "13") {
            $("#okButton_2").trigger("click")
        }
    })
}
$.fn.Ask = function (a) {
    $(this).click(function () {
        Ask2(a);
        return false
    })
};
$.fn.LineOne = function () {
    $(this).bind("focus", function () {
        if (this.blur) {
            this.blur()
        }
    })
};
function Message(c, g, d, a) {
    var f = {
        CloseTime: 0,
        Msg: "",
        Style: "Alert",
        Tit: "Title",
        Height: 99,
        callback: "Close()",
        fromObj: null,
        call: null
    };
    a = Object.extend(f, a);
    if (g.length > 0) {
        a.Msg = g
    }
    if (c != "Title") {
        a.Tit = c
    }
    if (d) {
        a.Style = d
    }
    var e = '<div class="AlertMessage">';
    if (a.CloseTime == 0) {
        e += '<div class="MessageTitle"><a class="xx" title="' + chbtc.L("关闭") + '" onclick="Close();" href="javascript:Void();">' + chbtc.L("关闭") + "</a>" + a.Tit + "</div>"
    }
    e += '<div class="Message' + a.Style + '"><div class="Message">' + a.Msg + "</div></div>";
    if (a.CloseTime == 0) {
        e += '<div class="MessageControl"><div class="MessageControl2"><a onfocus="this.blur()" id="okButton_2013"  class="ok" onclick="' + a.callback + '" href="javascript:Void();"><h4>' + chbtc.L("确定") + "</h4></a></div></div>"
    }
    e += "</div>";
    NetBox.show(e, 550, a.Height, a.CloseTime, a.fromObj);
    if (typeof a.call == "function") {
        $("#okButton_2013").attr("href", "javascript:;");
        $("#okButton_2013").click(function () {
            a.call();
            return false
        })
    }
    $("body").bind("keyup", function (i) {
        if (i.keyCode == "13") {
            $("#okButton_2013").trigger("click")
        }
    })
}
function Alert(c, a) {
    Message(chbtc.L("BW提示您"), c + "", "Alert", a)
}
function Info(c, a) {
    Message(chbtc.L("BW提示您"), c + "", "Info", a)
}
function msg(c, a) {
    Message(chbtc.L("BW提示您"), c + "", "Info", a)
}
function Msg(c, a) {
    Message(chbtc.L("BW提示您"), c + "", "Info", a)
}
function Wrong(c, a) {
    Message(chbtc.L("BW提示您"), c + "", "Wrong", a)
}
function Right(c, a) {
    Message(chbtc.L("BW提示您"), c + "", "Right", a)
}
function Help(c, a) {
    Message(chbtc.L("BW提示您"), c + "", "Help", a)
}
function Iframe(c) {
    var g = {
        Url: "",
        zoomSpeedIn: 100,
        zoomSpeedOut: 100,
        Width: 540,
        Height: 190,
        Title: "",
        overlayShow: false,
        modal: true,
        isShowIframeTitle: true,
        isShowClose: true,
        isIframeAutoHeight: false,
        scrolling: "auto",
        overlayOpacity: 0.5,
        overlayColor: "#000000",
        padding: 0,
        IsShow: false,
        IsCloseOnModal: false,
        fromObj: null
    };
    c = Object.extend(g, c);
    var e = "";
    if (c.Url.indexOf("?") < 0) {
        e += "?a=1"
    }
    if (c.Url.indexOf("iframe=1") < 0) {
        e += "&iframe=1"
    }
    c.Url += e;
    var f = (c.height - 44) + "px";
    var d = "";
    if (c.isShowIframeTitle) {
        var a = "";
        if (!c.isShowClose) {
            a = " noClose"
        }
        if (c.isIframeAutoHeight) {
            d = '<div class="popIframeTitle' + a + '"><div class="popIframeTitle">' + c.Title + '</div><div class="popIframeCloseC"><a class="popIframeClose" onfocus="this.blur()" href="javascript:Close()" title="关闭">Close</a></div></div><iframe id="fancybox-frame" name="fancybox-frame' + new Date().getTime() + '" frameborder="0"   hspace="0" ' + ($.browser.msie ? 'allowtransparency="true""' : "") + ' scrolling="' + c.scrolling + '" onload="$(this).height($(this).contents().height());NetBox.resize(0,($(this).contents().height()+44));" src="' + c.Url + '"></iframe>'
        } else {
            d = '<div class="popIframeTitle' + a + '"><div class="popIframeTitle">' + c.Title + '</div><div class="popIframeCloseC"><a class="popIframeClose" onfocus="this.blur()" href="javascript:Close()" title="关闭">Close</a></div></div><iframe id="fancybox-frame" name="fancybox-frame' + new Date().getTime() + '" frameborder="0" hspace="0"  scrolling="' + c.scrolling + '" style="height:' + (c.Height - 44) + 'px"  src="' + c.Url + '"></iframe>'
        }
    } else {
        if (c.isIframeAutoHeight) {
            d = '<iframe id="fancybox-frame" name="fancybox-frame' + new Date().getTime() + '" frameborder="0" hspace="0" ' + ($.browser.msie ? 'allowtransparency="true""' : "") + '  onload="$(this).height($(this).contents().height());NetBox.resize(0,($(this).contents().height()));" scrolling="' + c.scrolling + '" src="' + c.Url + '"></iframe>'
        } else {
            d = '<iframe id="fancybox-frame" name="fancybox-frame' + new Date().getTime() + '" frameborder="0" hspace="0" ' + ($.browser.msie ? 'allowtransparency="true""' : "") + ' scrolling="' + c.scrolling + '" src="' + c.Url + '"></iframe>'
        }
    }
    NetBox.show(d, c.Width, c.Height, 0, c.fromObj, c.IsCloseOnModal)
}
$.fn.Iframe = function (a) {
    $(this).click(function () {
        Iframe(a);
        return false
    })
};
$.extend({
    Iframe: function (a) {
        Iframe(a)
    }
});
jQuery.cookie = function (c, n, u) {
    if (typeof n != "undefined") {
        u = u || {};
        if (n === null) {
            n = "";
            u.expires = -1
        }
        var f = "";
        if (u.expires && (typeof u.expires == "number" || u.expires.toUTCString)) {
            var g;
            if (typeof u.expires == "number") {
                g = new Date();
                g.setTime(g.getTime() + (u.expires * 24 * 60 * 60 * 1000))
            } else {
                g = u.expires
            }
            f = "; expires=" + g.toUTCString()
        }
        var r = u.path ? "; path=" + (u.path) : "";
        var j = u.domain ? "; domain=" + (u.domain) : "";
        var a = u.secure ? "; secure" : "";
        document.cookie = [c, "=", encodeURIComponent(n), f, r, j, a].join("")
    } else {
        var e = null;
        if (document.cookie && document.cookie != "") {
            var q = document.cookie.split(";");
            for (var k = 0; k < q.length; k++) {
                var d = jQuery.trim(q[k]);
                if (d.substring(0, c.length + 1) == (c + "=")) {
                    e = decodeURIComponent(d.substring(c.length + 1));
                    break
                }
            }
        }
        return e
    }
};
$.fn.DropSelecter = function (k) {
    var d = {
        StyleHover: "hover",
        StyleSelect: "",
        ControlSeleter: "",
        Left: 0,
        Top: 0,
        IsShow: false,
        Auto: true,
        call: function () {
        }
    };
    k = $.extend(d, k);
    if (k.StyleSelect == "") {
        k.StyleSelect = k.StyleHover
    }
    var g = $(this);
    var j = $(k.ControlSeleter);
    var c;
    var i = 100;
    var e = g.position();
    if (!e) {
        return
    }
    j.css({left: e.left + k.Left, top: e.top + g.height() + k.Top});
    var f = true;
    var a = "";
    if (!$(this).data("DropSelecterId")) {
        a = numberID();
        g.data("DropSelecterId", a);
        if (!k.Auto) {
            $("body").click(function () {
                if (g.data("opened")) {
                    j.hide();
                    g.removeClass(k.StyleSelect);
                    g.removeData("opened")
                }
            });
            g.click(function (n) {
                if (g.data("opened")) {
                    j.hide();
                    g.removeClass(k.StyleSelect);
                    g.removeData("opened")
                } else {
                    $("body").trigger("click");
                    g.data("opened", true);
                    j.show();
                    g.addClass(k.StyleSelect)
                }
                n.stopPropagation()
            });
            j.click(function (n) {
                n.stopPropagation()
            })
        } else {
            j.mouseenter(function () {
                clearTimeout(c);
                c = setTimeout(function () {
                    j.show();
                    g.addClass(k.StyleHover)
                }, i)
            }).mouseleave(function () {
                clearTimeout(c);
                c = setTimeout(function () {
                    j.hide();
                    g.removeClass(k.StyleHover)
                }, i)
            })
        }
        g.focus(function () {
            $(this).blur()
        }).mouseover(function () {
            clearTimeout(c);
            c = setTimeout(function () {
                e = g.position();
                j.css({left: e.left + k.Left, top: e.top + g.height() + k.Top});
                if (k.Auto) {
                    j.show();
                    if (f) {
                        (function (n) {
                            return k.call(n)
                        })(1)
                    }
                }
                g.addClass(k.StyleHover);
                f = false
            }, i)
        }).mouseleave(function () {
            clearTimeout(c);
            c = setTimeout(function () {
                if (k.Auto) {
                    j.hide()
                }
                g.removeClass(k.StyleHover)
            }, i)
        });
        $(window).resize(function () {
            e = g.position();
            j.css({left: e.left + k.Left, top: e.top + g.height() + k.Top})
        })
    }
};
$.fn.UiTitle = function () {
    $(this).MyTitle({defaultCss: "tipsr", title: $(this).attr("mytitle"), html: true})
};
$.fn.UIButton = function () {
    if ($(this).attr("NewStyle")) {
        return
    }
    $(this).attr("NewStyle", "true");
    $(this).wrap("<a  href='javascript:void(0)'></a>");
    if (!$(this).attr("StyleName")) {
        $(this).attr("StyleName", "buttonCommon")
    }
    if ($(this).attr("disabled")) {
        $(this).parent().addClass($(this).attr("StyleName") + "Disabled")
    } else {
        $(this).parent().addClass($(this).attr("StyleName"))
    }
    $(this).bind("focus", function () {
        if (this.blur) {
            this.blur()
        }
    })
};
$.fn.UiText = function () {
    if (!$(this).attr("NewStyle")) {
        $(this).attr("NewStyle");
        var a = $(this);
        if (a.attr("mytitle")) {
            a.MyTitle({defaultCss: "tipsr", title: a.attr("mytitle"), trigger: "manual"})
        }
        if ($(this).attr("valueDemo")) {
            if ($(this).val() == "") {
                $(this).val($(this).attr("valueDemo"));
                $(this).css({color: "#C7C7C7"})
            }
        }
        $(this).bind("focus", function () {
            $(this).removeAttr("errorStyle");
            if ($(this).attr("valueDemo")) {
                if ($(this).attr("valueDemo") == $(this).val()) {
                    $(this).val("");
                    $(this).css({color: "#333333"})
                }
            }
            a.addClass("inputFocue");
            if (a.data("tipsy")) {
                a.MyTitle("show")
            }
        }).bind("blur", function () {
            if ($(this).attr("valueDemo")) {
                if ($(this).val() == "") {
                    $(this).val($(this).attr("valueDemo"));
                    $(this).css({color: "#C7C7C7"})
                }
            }
            a.removeClass("inputFocue");
            if (a.data("tipsy")) {
                a.MyTitle("hide")
            }
            if ($(this).val().length > 0 && $(this).attr("valueDemo") != $(this).val() && $(this).attr("defaultV") != $(this).val()) {
                CheckTextBox($(this))
            }
        });
        if (!$(this).attr("noHide") && a.data("tipsy")) {
            a.keydown(function () {
                a.MyTitle("hide")
            })
        }
    }
};
$.fn.UiCheckbox = function (a) {
    $(":checkbox+label", this).each(function () {
        $(this).addClass("checkbox");
        if ($(this).prev().is(":disabled") == false) {
            if ($(this).prev().is(":checked")) {
                $(this).addClass("checked")
            }
        } else {
            $(this).addClass("disabled")
        }
    }).click(function (c) {
        if (!$(this).prev().is(":checked")) {
            $(this).addClass("checked");
            $(this).prev()[0].checked = true
        } else {
            $(this).removeClass("checked");
            $(this).prev()[0].checked = false
        }
        c.stopPropagation()
    }).prev().hide()
};
function EnableRadio(c) {
    var d = $(c);
    var a = d.parent().find("a:first");
    a.removeClass("jqTransformRadioEnabled").css({cursor: "pointer"});
    d.removeAttr("disabled");
    d.change(function () {
        d[0].checked && a.addClass("jqTransformRadioChecked") || a.removeClass("jqTransformRadioChecked");
        return true
    });
    a.click(function () {
        if (d.attr("disabled")) {
            return false
        }
        d.trigger("click").trigger("change");
        $('input[name="' + d.attr("name") + '"]').not(d).each(function () {
            $(this).attr("type") == "radio" && $(this).trigger("change")
        });
        return false
    })
}
function DisenableRadio(c) {
    var d = $(c);
    var a = d.parent().find("a:first");
    d.attr("disabled", "disabled");
    a.addClass("jqTransformRadioEnabled").css({cursor: "default"});
    a.unbind("click");
    d.unbind("change")
}
$.fn.UiRadio = function () {
    return this.each(function () {
        if (!$(this).attr("NewStyle")) {
            if ($(this).hasClass("jqTransformHidden")) {
                return
            }
            var d = $(this);
            var c = this;
            d.addClass("jqTransformHidden");
            var a;
            if ($(this).parent().hasClass("jqTransformRadioWrapper")) {
                a = d.parent().find("a:first")
            } else {
                a = $('<a  class="jqTransformRadio" style="cursor:pointer;"></a>');
                d.wrap('<span class="jqTransformRadioWrapper"></span>').parent().prepend(a)
            }
            if (d.attr("disabled")) {
                a.addClass("jqTransformRadioEnabled").css({cursor: "default"});
                return
            } else {
                d.change(function () {
                    c.checked && a.addClass("jqTransformRadioChecked") || a.removeClass("jqTransformRadioChecked");
                    return true
                });
                a.click(function () {
                    if (d.attr("disabled")) {
                        return false
                    }
                    d.trigger("click").trigger("change");
                    $('input[name="' + d.attr("name") + '"]', c.form).not(d).each(function () {
                        $(this).attr("type") == "radio" && $(this).trigger("change")
                    });
                    return false
                });
                c.checked && a.addClass("jqTransformRadioChecked")
            }
            d.hide()
        }
    })
};
$.fn.UiSelect = function (d) {
    var g = {
        StyleNormal: "SelectGray",
        StyleHover: "SelectBlue",
        StyleDropDown: "SelectDropDown",
        itemHeight: 24,
        InnerWidth: 0,
        InnerWidthOffset: -15,
        OuterWidth: 0,
        OuterWidthOffset: 12,
        MaxShow: 10,
        Top: 0,
        Left: 0,
        ControlSeleter: "",
        Auto: false,
        IsShow: true
    };
    if ($(this).attr("NewStyle")) {
        return
    }
    d = $.extend(g, d);
    var f = $(this);
    var c = $(this);
    var i = "0";
    var e = false;
    if (f.attr("SelectId")) {
        i = f.attr("SelectId").split("_")[1];
        c = $("#select_" + i);
        c.removeData("dropDownSelected");
        c.find("span i").text(f.find("option:selected").text());
        e = true
    } else {
        i = numberID();
        f.attr("SelectId", "select_" + i);
        c = $("<div id='select_" + i + "' class='" + d.StyleNormal + "'><span><i><i><span><div>").insertAfter(f);
        c.find("span i").text(f.find("option:selected").text());
        c.mouseover(function () {
            c.addClass(d.StyleHover)
        }).mouseout(function () {
            c.removeClass(d.StyleHover)
        })
    }
    if ($.browser.safari) {
        d.InnerWidthOffset = d.InnerWidthOffset + 26;
        d.OuterWidthOffset = d.OuterWidthOffset + 26
    }
    var a = f.selectedIndex;
    if (d.InnerWidth == 0) {
        d.InnerWidth = f.width() + d.InnerWidthOffset
    }
    if (d.OuterWidth == 0) {
        d.OuterWidth = f.width() + d.OuterWidthOffset
    }
    c.find("span i").css("width", d.InnerWidth);
    c.mouseover(function () {
        d.InnerWidth = f.width() + d.InnerWidthOffset;
        d.OuterWidth = f.width() + d.OuterWidthOffset;
        e = $("#down_" + i).length > 0;
        if (!$(this).data("dropDownSelected")) {
            $(this).data("dropDownSelected", true);
            var r = "";
            if (e) {
                r = "<span id='downDiv_" + i + "' >"
            } else {
                r = "<p id='down_" + i + "' ><span id='downDiv_" + i + "' >"
            }
            var q = 0;
            $("option", f).each(function (u) {
                if (a == u) {
                    r += '<a href="javascript:void(0)"  class="selected" index="' + u + '">' + $(this).text() + "</a>"
                } else {
                    r += '<a href="javascript:void(0)"  index="' + u + '">' + $(this).text() + "</a>"
                }
                q++
            });
            if (e) {
                r += "</span>"
            } else {
                r += "</span></p>"
            }
            var k = false;
            var n = q * d.itemHeight;
            if (d.MaxShow < q) {
                k = true;
                n = d.MaxShow * d.itemHeight
            }
            var j = $(r);
            if (e) {
                j = $("#down_" + i).html(r)
            }
            j.css({width: d.OuterWidth, height: n});
            j.find("span").css({width: d.OuterWidth, height: n});
            j.find("a").click(function () {
                if (c.data("tipsy")) {
                    c.MyTitle("hide")
                }
                f[0].selectedIndex = $(this).attr("index");
                j.find("a").removeClass("selected");
                $(this).addClass("selected");
                c.find("span i").text($(this).text());
                f.trigger("change");
                $(".TitleErrorRight:visible").hide("fast");
                $(".TitleErrorTop:visible").hide("fast");
                $("body").trigger("click");
                return false
            });
            j.find("a").focus(function () {
                $(this).blur()
            });
            if (!e) {
                j.appendTo(c)
            }
            c.DropSelecter({
                ControlSeleter: "#down_" + i,
                StyleHover: d.StyleHover,
                StyleSelect: d.StyleDropDown,
                Left: 0,
                Top: -13,
                Auto: false,
                IsShow: true
            });
            if (k) {
                j.css("padding-bottom", "3px");
                $("#down_" + i).show();
                $("#downDiv_" + i).jscroll({
                    W: "17px",
                    BgUrl: "url(/body/images/form/s_bg.png)",
                    Bg: "right 0 repeat-y",
                    Bar: {
                        Pos: "up",
                        Bd: {Out: "#a3c3d5", Hover: "#b7d5e6"},
                        Bg: {Out: "-51px 0 repeat-y", Hover: "-66px 0 repeat-y", Focus: "-81px 0 repeat-y"}
                    },
                    Btn: {
                        btn: true,
                        uBg: {Out: "0 0", Hover: "-17px 0", Focus: "-34px 0"},
                        dBg: {Out: "0 -21px", Hover: "-17px -21px", Focus: "-34px -21px"}
                    },
                    Fn: function () {
                    }
                })
            }
            j.hide()
        }
    });
    $(this).hide()
};
function changeCheckBox(d) {
    var c = T$("ck_" + d);
    var a = T$(d);
    if (a.checked) {
        c.className = "checkbox";
        a.checked = false
    } else {
        c.className = "checkbox checked";
        a.checked = true
    }
    $(a).trigger("change")
}
function UICheckbox(c) {
    if ($(c).attr("newStyle")) {
        return
    }
    var a = c.getAttribute("id");
    if (a) {
        if ($("#" + a).attr("checked")) {
            $(c).after("<label onclick=\"changeCheckBox('" + a + "')\" class='checkbox checked' id='ck_" + a + "'></label>")
        } else {
            $(c).after("<label onclick=\"changeCheckBox('" + a + "')\" class='checkbox' id='ck_" + a + "'></label>")
        }
    } else {
        var d = numberID();
        c.setAttribute("id", d);
        if ($("#" + d).attr("checked")) {
            $(c).after("<label onclick=\"changeCheckBox('" + d + "')\" class='checkbox' id='ck_" + d + "'></label>")
        } else {
            $(c).after("<label onclick=\"changeCheckBox('" + d + "')\" class='checkbox checked' id='ck_" + d + "'></label>")
        }
    }
}
function matchof(c, d, e) {
    var a = $("#" + d).val();
    a = a.replaceAll("(", "（");
    a = a.replaceAll(")", "）");
    a = a.replaceAll("\\", "|");
    a = a.replace(/["'\n\r\t]/g, " ");
    if (a == c) {
        return true
    } else {
        return false
    }
}
function num(c, d) {
    var a = /^[0-9]{1,20}$id/;
    if (c.search("^-?\\d+(\\.\\d+)?$") == 0) {
        return true
    } else {
        if (c.toString().length == 0) {
            return true
        } else {
            return false
        }
    }
}
function cnChar(c, d) {
    var a = new RegExp("[^\x00-\xff]");
    return a.exec(c)
}
function strDateTime(f, g) {
    if (f == null) {
        return true
    }
    if (f == "") {
        return true
    }
    var a = /^(\d{1,4})(-|\/)(\d{1,2})\2(\d{1,2}) (\d{1,2}):(\d{1,2}):(\d{1,2})$/;
    var c = f.match(a);
    if (c == null) {
        return false
    }
    var e = new Date(c[1], c[3] - 1, c[4], c[5], c[6], c[7]);
    return (e.getFullYear() == c[1] && (e.getMonth() + 1) == c[3] && e.getDate() == c[4] && e.getHours() == c[5] && e.getMinutes() == c[6] && e.getSeconds() == c[7])
}
function limit(f, e, d, g) {
    var c = f;
    var a = c.match(/[^ -~]/g) == null ? c.length : c.length + c.match(/[^ -~]/g).length;
    if (g && f == g) {
        return true
    }
    if (a >= e && a <= d) {
        return true
    } else {
        return false
    }
}
function notmatch(c, a, d) {
    if (c == a) {
        return false
    } else {
        return true
    }
}
function email(a, c) {
    if (a.length == 0) {
        return false
    }
    if (a.indexOf(".") > 0 && a.indexOf("@") > 0) {
        return true
    } else {
        return false
    }
}
function telphone(c, d) {
    var a = new RegExp("[0-9]{2})+-([0-9]{4})+-([0-9]{4}");
    if (!a.exec(c)) {
        return false
    } else {
        return true
    }
}
function isMobile(a) {
    var c = /^1[3|4|5|6|7|8][0-9]\d{8}$/;
    if (c.test(a)) {
        return true
    } else {
        return false
    }
}
function isEmail(a) {
    var c = /^([a-zA-Z0-9_\\.-])+\@([a-z0-9A-Z]+(-[a-z0-9A-Z]+)?\.)+[a-zA-Z]{2,}$/;
    if (c.test(a)) {
        return true
    } else {
        return false
    }
}
function isPhone(a) {
    var c = /((\d{11})|^((\d{7,8})|(\d{4}|\d{3})-(\d{7,8})|(\d{4}|\d{3})-(\d{7,8})-(\d{4}|\d{3}|\d{2}|\d{1})|(\d{7,8})-(\d{4}|\d{3}|\d{2}|\d{1}))$)/;
    if (c.test(a)) {
        return true
    } else {
        return false
    }
}
function Ui(f) {
    var g = "";
    if (!$(f)) {
        return null
    }
    var a = T$(f).getElementsByTagName("*");
    for (var d = 0; d < a.length; d++) {
        if (a[d].tagName != undefined) {
            var c = a[d].tagName.toLowerCase();
            var e = a[d].getAttribute("id");
            if (c == "input" || c == "textarea") {
                if (a[d].type.toLowerCase() == "text" || a[d].type.toLowerCase() == "hidden" || a[d].type.toLowerCase() == "password" || a[d].tagName.toLowerCase() == "textarea") {
                    $(a[d]).UiText()
                } else {
                    if (a[d].type.toLowerCase() == "radio") {
                        $(a[d]).UiRadio()
                    } else {
                        if (a[d].type.toLowerCase() == "checkbox") {
                            UICheckbox(a[d])
                        } else {
                            if (a[d].type.toLowerCase() == "button") {
                                $(a[d]).UIButton()
                            }
                        }
                    }
                }
            } else {
                if (c == "select") {
                    $(a[d]).UiSelect()
                }
            }
        }
    }
}
$.fn.Ui = function () {
    var a = [];
    if ($(this)[0]) {
        a = $(this)[0].getElementsByTagName("*")
    }
    for (var d = 0; d < a.length; d++) {
        if (a[d].tagName != undefined) {
            var c = a[d].tagName.toLowerCase();
            var e = a[d].getAttribute("id");
            if (c == "input" || c == "textarea") {
                if (a[d].type.toLowerCase() == "text" || a[d].type.toLowerCase() == "hidden" || a[d].type.toLowerCase() == "password" || a[d].tagName.toLowerCase() == "textarea") {
                    $(a[d]).UiText()
                } else {
                    if (a[d].type.toLowerCase() == "radio") {
                        $(a[d]).UiRadio()
                    } else {
                        if (a[d].type.toLowerCase() == "checkbox") {
                            UICheckbox(a[d])
                        } else {
                            if (a[d].type.toLowerCase() == "button") {
                                $(a[d]).UIButton()
                            }
                        }
                    }
                }
            } else {
                if (c == "select") {
                    $(a[d]).UiSelect()
                }
            }
        }
    }
};
$.fn.AStop = function (a) {
    $(this).click(function (d) {
        var d = d ? d : window.event, c;
        d.preventDefault();
        if (navigator.appName == "Microsoft Internet Explorer") {
            c = d.srcElement
        } else {
            c = d.target
        }
        if (c.getAttribute("disabled")) {
            return false
        }
    })
};
jQuery.fn.fixedIn = function (c) {
    var d = {offsetX: 0, offsetY: 0, selecter: ""};
    c = $.extend(d, c);
    var f = $(this);
    var e = numberID();
    var a = $(c.selecter);
    $(window).scroll(function () {
        var k = $(window).scrollTop();
        var n = f.offset().top;
        var i = n + f.height() - 2 * c.offsetY - a.height();
        var j = k - n + c.offsetY;
        var g = k - n + c.offsetY;
        if (k > n && k <= i) {
            if (jQuery.browser.msie && jQuery.browser.version == "6.0") {
                a.css({position: "absolute", top: j + "px", left: c.offsetX + "px"})
            } else {
                a.css({position: "fixed", top: c.offsetY + "px", left: (f.offset().left + c.offsetX) + "px"})
            }
        } else {
            if (k > i) {
                a.css({position: "absolute", top: (f.height() - c.offsetY - a.height()) + "px", left: c.offsetX + "px"})
            } else {
                a.css({position: "absolute", top: c.offsetY + "px", left: c.offsetX + "px"})
            }
        }
    });
    $(window).trigger("scroll")
};
$.fn.Loadding = function (a) {
    var e = {
        StyleName: "Loadding",
        Str: chbtc.L("加载中") + "…",
        Width: 0,
        Height: 0,
        Position: "on",
        IsShow: true,
        OffsetY: 0,
        OffsetX: 0,
        OffsetYGIF: 0,
        OffsetXGIF: 0
    };
    a = $.extend(e, a);
    if (!a.IsShow) {
        var d = $(this).attr("LoddingId");
        $(".Loadding").remove();
        $(this).removeAttr("LoddingId");
        return
    }
    if ($(this).offset() == null) {
        return
    }
    var f = numberID();
    var c = ";position:absolute;";
    if (a.Position == "top") {
        c += "margin-top:-" + $(this).outerHeight() + "px;"
    } else {
        if (a.Position == "bottom") {
            c += "margin-top:0;"
        } else {
            if (a.Position == "left") {
                c += "margin-left:-" + $(this).outerWidth() + "px;"
            } else {
                if (a.Position == "right") {
                    c += "margin-right:" + $(this).outerWidth() + "px;"
                } else {
                    c += "top:" + $(this).offset().top + "px;left:" + $(this).offset().left + "px;"
                }
            }
        }
    }
    if (a.Width == 0) {
        c += "width:" + $(this).outerWidth() + "px;"
    } else {
        if (a.Width > 0) {
            c += "width:" + a.Width + "px;"
        }
    }
    if (a.OffsetX == 0) {
        a.OffsetXGIF = $(this).outerWidth() / 2 - 110
    } else {
        a.OffsetXGIF = $(this).outerWidth() / 2 - 110 + a.OffsetXGIF
    }
    if (a.Height == 0) {
        c += "height:" + $(this).outerHeight() + "px;"
    } else {
        if (a.Height > 0) {
            c += "height:" + a.Height + "px;"
        }
    }
    if ($(".Loadding").length < 1) {
        $(this).after('<div id="Lodding_' + f + '" class="' + a.StyleName + '" style="' + c + '" ><div class="GIF"><div class="img"></div><div class="str">' + a.Str + "</div></div></div>");
        $(this).attr("LoddingId", f);
        $("#Lodding_" + f).fadeIn()
    }
};
function stopDefault(a) {
    if (a && a.preventDefault) {
        a.preventDefault()
    } else {
        window.event.returnValue = false
    }
    return false
}
function checkPageInput(c, a) {
    c.keyCode == 13 && $("#JumpButton").trigger("click")
}
jQuery.fn.ScrollTo = function (a) {
    var c = {Speed: 300, Top: 0};
    a = $.extend(c, a);
    o = jQuery.speed(a.Speed);
    return this.each(function () {
        new jQuery.fx.ScrollTo(this, o, a.Top)
    })
};
jQuery.fx.ScrollTo = function (c, f, a) {
    var d = this;
    d.o = f;
    d.e = c;
    d.p = jQuery.getPos(c);
    d.s = jQuery.getScroll();
    d.clear = function () {
        clearInterval(d.timer);
        d.timer = null
    };
    d.t = (new Date).getTime();
    d.step = function () {
        var e = (new Date).getTime();
        var g = (e - d.t) / d.o.duration;
        if (e >= d.o.duration + d.t) {
            d.clear();
            setTimeout(function () {
                d.scroll(d.p.y, d.p.x)
            }, 13)
        } else {
            st = ((-Math.cos(g * Math.PI) / 2) + 0.5) * (d.p.y - d.s.t) + d.s.t;
            sl = ((-Math.cos(g * Math.PI) / 2) + 0.5) * (d.p.x - d.s.l) + d.s.l;
            d.scroll(st, sl)
        }
    };
    d.scroll = function (g, e) {
        window.scrollTo(e - a, g - a)
    };
    d.timer = setInterval(function () {
        d.step()
    }, 13)
};
jQuery.intval = function (a) {
    a = parseInt(a);
    return isNaN(a) ? 0 : a
};
jQuery.getClient = function (a) {
    if (a) {
        w = a.clientWidth;
        h = a.clientHeight
    } else {
        w = (window.innerWidth) ? window.innerWidth : (document.documentElement && document.documentElement.clientWidth) ? document.documentElement.clientWidth : document.body.offsetWidth;
        h = (window.innerHeight) ? window.innerHeight : (document.documentElement && document.documentElement.clientHeight) ? document.documentElement.clientHeight : document.body.offsetHeight
    }
    return {w: w, h: h}
};
jQuery.getScroll = function (a) {
    if (a) {
        t = a.scrollTop;
        l = a.scrollLeft;
        w = a.scrollWidth;
        h = a.scrollHeight
    } else {
        if (document.documentElement && document.documentElement.scrollTop) {
            t = document.documentElement.scrollTop;
            l = document.documentElement.scrollLeft;
            w = document.documentElement.scrollWidth;
            h = document.documentElement.scrollHeight
        } else {
            if (document.body) {
                t = document.body.scrollTop;
                l = document.body.scrollLeft;
                w = document.body.scrollWidth;
                h = document.body.scrollHeight
            }
        }
    }
    return {t: t, l: l, w: w, h: h}
};
jQuery.getPos = function (i) {
    var c = 0;
    var d = 0;
    var a = jQuery.intval(jQuery.css(i, "width"));
    var f = jQuery.intval(jQuery.css(i, "height"));
    var g = i.offsetWidth;
    var j = i.offsetHeight;
    while (i.offsetParent) {
        c += i.offsetLeft + (i.currentStyle ? jQuery.intval(i.currentStyle.borderLeftWidth) : 0);
        d += i.offsetTop + (i.currentStyle ? jQuery.intval(i.currentStyle.borderTopWidth) : 0);
        i = i.offsetParent
    }
    c += i.offsetLeft + (i.currentStyle ? jQuery.intval(i.currentStyle.borderLeftWidth) : 0);
    d += i.offsetTop + (i.currentStyle ? jQuery.intval(i.currentStyle.borderTopWidth) : 0);
    return {x: c, y: d, w: a, h: f, wb: g, hb: j}
};
jQuery.fx.ScrollTo = function (c, f, a) {
    var d = this;
    d.o = f;
    d.e = c;
    d.p = jQuery.getPos(c);
    d.s = jQuery.getScroll();
    d.clear = function () {
        clearInterval(d.timer);
        d.timer = null
    };
    d.t = (new Date).getTime();
    d.step = function () {
        var e = (new Date).getTime();
        var g = (e - d.t) / d.o.duration;
        if (e >= d.o.duration + d.t) {
            d.clear();
            setTimeout(function () {
                d.scroll(d.p.y, d.p.x)
            }, 13)
        } else {
            st = ((-Math.cos(g * Math.PI) / 2) + 0.5) * (d.p.y - d.s.t) + d.s.t;
            sl = ((-Math.cos(g * Math.PI) / 2) + 0.5) * (d.p.x - d.s.l) + d.s.l;
            d.scroll(st, sl)
        }
    };
    d.scroll = function (g, e) {
        window.scrollTo(e - a, g - a)
    };
    d.timer = setInterval(function () {
        d.step()
    }, 13)
};
function ajaxNextPage(c) {
    var a = $(this).attr("href").replace("ForBuy", "ForBuy/ajaxList");
    ajaxUrl(a, true, true)
}
Title = function () {
    var d, c, i, g, j, a, e = 0;
    return {
        Show: function (f) {
            var k = {CloseTime: 0, Height: 99, callback: "Close()", fromObj: null};
            f = Object.extend(k, f);
            if (!e) {
                d = document.createElement("div");
                p.id = "tinybox";
                c = document.createElement("div");
                m.id = "tinymask";
                b = document.createElement("div");
                b.id = "tinycontent";
                document.body.appendChild(m);
                document.body.appendChild(p);
                p.appendChild(b);
                m.onclick = NetBox.hide;
                window.onresize = NetBox.resize;
                e = 1
            }
        }, Close: function () {
        }
    }
}();
$.fn.Title = function (a) {
    var d = {StyleName: "TitleError", Title: "Message", Position: "w", Offset: 10, OffsetY: 5, IsShow: true};
    a = $.extend(d, a);
    if ($(this).attr("TitleStyleName")) {
        a.StyleName = $(this).attr("TitleStyleName")
    }
    if ($(this).attr("title")) {
        a.Title = $(this).attr("title")
    }
    if (a.Title.indexOf("#") == 0) {
        a.Title = $(a.Title).html()
    }
    if ($(this).attr("Offset")) {
        a.Offset = $(this).attr("Offset")
    }
    if ($(this).attr("TitlePosition")) {
        a.Position = $(this).attr("TitlePosition")
    }
    var e = $(this).attr("id") + a.StyleName;
    var c = $(this).offset();
    if (a.IsShow) {
        $(this).MyTitle({defaultCss: "tipsr", title: chbtc.L("BW提示您"), trigger: "manual"});
        if (!$(this).attr("errmsg")) {
            $(this).attr("errmsg", a.Title)
        }
        $(this).attr("errorStyle", "tipsr");
        $(this).MyTitle("show")
    } else {
        $(this).MyTitle("hide");
        $(this).removeAttr("errorStyle")
    }
};
function CheckTextBox(obj) {
    var pattern = obj.attr("pattern");
    var defaultV = obj.attr("defaultV");
    if (obj.attr("webpattern")) {
        pattern = obj.attr("webpattern")
    }
    if (pattern != null && pattern != undefined && pattern != "") {
        if ($(obj).attr("valueDemo") || defaultV && $(obj).is(":visible")) {
            if ($(obj).attr("valueDemo") == $(obj).val() || $(obj).val() == defaultV) {
                if (!$(obj).data("tipsy")) {
                    $(obj).MyTitle({defaultCss: "tipsr", title: "An error occurs here!", trigger: "manual"})
                }
                $(obj).attr("errorStyle", "tipsr");
                $(obj).MyTitle("show");
                return false
            }
        }
        var subs = pattern.split(";");
        for (var j = 0; j < subs.length; j++) {
            var cValue = obj.val();
            cValue = cValue.replaceAll("(", "（");
            cValue = cValue.replaceAll(")", "）");
            cValue = cValue.replaceAll("\\", "|");
            var paraValue = cValue.replace(/["'\n\r\t]/g, " ");
            var execStr = subs[j].replace("(", '("' + paraValue + '",');
            if (execStr.indexOf(",)") > 0) {
                execStr = execStr.replace(")", "'" + obj.attr("defaultV") + "')")
            } else {
                execStr = execStr.replace(")", ",'" + obj.attr("defaultV") + "')")
            }
            if (!this.eval(execStr)) {
                var srcEm = $(obj).attr("errormsg");
                if (chbtc.form.error.length > 0) {
                    $(obj).attr("errormsg", chbtc.form.error)
                }
                if (!$(obj).data("tipsy")) {
                    $(obj).MyTitle({defaultCss: "tipsr", title: "请重新填写本信息。", trigger: "manual"})
                }
                $(obj).attr("errorStyle", "tipsr");
                $(obj).MyTitle("show");
                if (chbtc.form.error.length > 0) {
                    $(obj).attr("errormsg", srcEm);
                    chbtc.form.error = ""
                }
                return false
            } else {
                $(obj).removeAttr("errorStyle");
                showRightIcon(obj)
            }
        }
        if ($(obj).data("tipsy")) {
            $(obj).MyTitle("hide")
        }
        return true
    } else {
        return true
    }
}
function showRightIcon(c) {
    var d = $(c).offset();
    var a = $("<div class='rightIcon' style='left:" + (d.left + $(c).width() + 20) + "px;top:" + (d.top + 5) + "px'></div>").appendTo("body");
    $(c).bind("focus", function () {
        a.remove()
    });
    $(c).bind("keydown", function () {
        a.remove()
    })
}
function CheckSelect(a) {
    var c = a.selectedIndex;
    var e = a.getAttribute("pattern");
    if (a.getAttribute("webpattern")) {
        e = a.getAttribute("webpattern")
    }
    if (e != null && e != undefined && c <= 0) {
        var d = $("#" + $(a).attr("selectid"));
        if (!d.data("tipsy")) {
            d.MyTitle({defaultCss: "tipsr", title: "请选择该选项。", trigger: "manual"})
        }
        if ($(a).attr("errormsg")) {
            d.attr("errormsg", $(a).attr("errormsg"))
        }
        d.attr("errorStyle", "tipsr");
        d.MyTitle("show");
        return false
    } else {
        return true
    }
}
function GoOn(a) {
    var c = window.eventList[a];
    window.eventList[a] = null;
    if (c.NextStep) {
        c.NextStep()
    } else {
        c()
    }
}
function Pause(d, e) {
    if (window.eventList == null) {
        window.eventList = new Array()
    }
    var c = -1;
    for (var a = 0; a < window.eventList.length; a++) {
        if (window.eventList[a] == null) {
            window.eventList[a] = d;
            c = a;
            break
        }
    }
    if (c == -1) {
        c = window.eventList.length;
        window.eventList[c] = d
    }
    setTimeout("GoOn(" + c + ")", e)
}
function FormToStr(a) {
    return FormToStrFun(a, false)
}
function errCalback(a) {
    Close()
}
function errSelectCalback(a) {
    Close()
}
var currentErrorObj = null;
function FormToStrFun(g, f) {
    var j = "";
    if (!T$(g)) {
        return null
    }
    var a = T$(g).getElementsByTagName("*");
    for (var e = 0; e < a.length; e++) {
        if (a[e].tagName != undefined) {
            if (a[e].tagName.toLowerCase() == "input" || a[e].tagName.toLowerCase() == "textarea") {
                if (a[e].type.toLowerCase() == "text" || a[e].type.toLowerCase() == "hidden" || a[e].type.toLowerCase() == "password" || a[e].tagName.toLowerCase() == "textarea") {
                    if (a[e].name.length > 0) {
                        var d = $(a[e]).attr("defaultV");
                        if ($(a[e]).attr("valueDemo") != $(a[e]).val() && d != $(a[e]).val() || a[e].type.toLowerCase() == "hidden" || $(a[e]).attr("readonly")) {
                            j += "&" + encodeURIComponent(a[e].name) + "=" + encodeURIComponent($.trim(a[e].value))
                        } else {
                            j += "&" + encodeURIComponent(a[e].name) + "="
                        }
                    }
                    if (!CheckTextBox($(a[e]))) {
                        if ($(a[e]).attr("errorName")) {
                            currentErrorObj = a[e]
                        } else {
                            if (f) {
                                $(a[e]).ScrollTo({Top: 100})
                            }
                        }
                        return null
                    }
                } else {
                    if (a[e].type.toLowerCase() == "checkbox") {
                        if (a[e].checked) {
                            if (j.indexOf("&" + encodeURIComponent(a[e].name) + "=") > -1) {
                                j = j.replace("&" + encodeURIComponent(a[e].name) + "=", "&" + encodeURIComponent(a[e].name) + "=" + encodeURIComponent(a[e].value) + encodeURIComponent("\u2229"))
                            } else {
                                j += "&" + encodeURIComponent(a[e].name) + "=" + encodeURIComponent(a[e].value)
                            }
                        }
                    } else {
                        if (a[e].type.toLowerCase() == "radio" || a[e].type.toLowerCase() == "checkbox") {
                            if (a[e].checked) {
                                j += "&" + encodeURIComponent(a[e].name) + "=" + encodeURIComponent(a[e].value)
                            }
                        }
                    }
                }
            } else {
                if (a[e].tagName.toLowerCase() == "select") {
                    var c = a[e].selectedIndex;
                    if (a[e].getAttribute("pattern") != null && a[e].getAttribute("pattern") != undefined && c <= 0) {
                        if (!CheckSelect(a[e])) {
                            if ($(a[e]).attr("errorName")) {
                                currentErrorObj = a[e];
                                Wrong("请选择" + $(a[e]).attr("errorName") + "！", {callback: "errSelectCalback(" + f + ")"})
                            } else {
                                if (f) {
                                    $(a[e]).parent().ScrollTo({Top: 100})
                                }
                            }
                            alert(20);
                            return null
                        }
                    } else {
                    }
                    if (c >= 0) {
                        if (a[e].options) {
                            if (a[e].name.length > 0) {
                                j += "&" + encodeURIComponent(a[e].name) + "=" + encodeURIComponent(a[e].options[c].value)
                            }
                        }
                    }
                }
            }
        }
    }
    return j.substring(1, j.length)
}
function Redirect(a) {
    self.location = a
}
(function (d) {
    function a(e) {
        if (e.attr("title") || typeof(e.attr("mytitle")) != "string") {
            e.attr("mytitle", e.attr("title") || "").removeAttr("title")
        }
    }

    function c(f, e) {
        this.$element = d(f);
        this.options = e;
        this.enabled = true;
        a(this.$element)
    }

    c.prototype = {
        show: function () {
            var u = this.getTitle();
            if (u && this.enabled) {
                var i = this.tip();
                i.addClass("tipsr");
                var x = this.$element.attr("errorStyle");
                if (!x) {
                    if (this.$element.attr("mytitle")) {
                        i.find(".inner-inner")[this.options.html ? "html" : "text"](this.$element.attr("mytitle"))
                    } else {
                        i.find(".inner-inner")[this.options.html ? "html" : "text"](u)
                    }
                    i[0].className = "tipsy"
                } else {
                    var n = this.$element.attr("errormsg");
                    if (this.$element.attr("errmsg")) {
                        n = this.$element.attr("errmsg")
                    }
                    i.find(".inner-inner")[this.options.html ? "html" : "text"](n);
                    i[0].className = "tipsy " + x
                }
                i.remove().css({top: 0, left: 0, visibility: "hidden", display: "block"}).appendTo(document.body);
                var v = i.width();
                i.find(".tipsy-inner").css({width: ""});
                v = 236;
                var q = d.extend({}, this.$element.offset(), {
                    width: this.$element[0].offsetWidth,
                    height: this.$element[0].offsetHeight
                });
                var e = i[0].offsetWidth, j = i[0].offsetHeight;
                var y = (typeof this.options.gravity == "function") ? this.options.gravity.call(this.$element[0]) : this.options.gravity;
                var g = this.$element.attr("position");
                if (g) {
                    y = g
                }
                var k = d(window).width();
                var z = d(window).height();
                if ((k < (q.left + q.width + v + 10)) && (y.charAt(0) == "w")) {
                    y = "n"
                }
                if (((z < (q.top + q.height + j + 10)) || k < (q.left + q.width / 2 + v / 2 + 10)) && (y.charAt(0) == "n")) {
                    y = "e"
                }
                if ((q.left < (v + 10)) && (y.charAt(0) == "e")) {
                    y = "s"
                }
                if ((q.top < (j + 10)) && (y.charAt(0) == "s")) {
                    y = "n"
                }
                var r;
                switch (y.charAt(0)) {
                    case"n":
                        r = {top: q.top + q.height + this.options.offset, left: q.left + q.width / 2 - e / 2};
                        break;
                    case"s":
                        var f = i.find(".tipsy-arrow");
                        f.css({top: j - f.height()});
                        r = {top: q.top - j - this.options.offset, left: q.left + q.width / 2 - e / 2};
                        break;
                    case"e":
                        r = {top: q.top + q.height / 2 - j / 2, left: q.left - e - this.options.offset};
                        break;
                    case"w":
                        r = {top: q.top + q.height / 2 - j / 2, left: q.left + q.width + this.options.offset};
                        break
                }
                if (y.length == 2) {
                    if (y.charAt(1) == "w") {
                        r.left = q.left + q.width / 2 - 15
                    } else {
                        r.left = q.left + q.width / 2 - e + 15
                    }
                }
                i.css(r).addClass("tipsy-" + y);
                if (this.options.fade) {
                    i.stop().css({
                        opacity: 0,
                        display: "block",
                        visibility: "visible"
                    }).animate({opacity: this.options.opacity})
                } else {
                    i.css({visibility: "visible", opacity: this.options.opacity})
                }
            }
        }, hide: function () {
            if (this.options.fade) {
                this.tip().stop().fadeOut(function () {
                    d(this).remove()
                })
            } else {
                this.tip().remove()
            }
        }, getTitle: function () {
            var g, e = this.$element, f = this.options;
            a(e);
            var g, f = this.options;
            if (typeof f.title == "string") {
                g = f.title == "title" ? e.attr("mytitle") : f.title
            } else {
                if (typeof f.title == "function") {
                    g = f.title.call(e[0])
                }
            }
            g = ("" + g).replace(/(^\s*|\s*$)/, "");
            return g || f.fallback
        }, tip: function () {
            if (!this.$tip) {
                this.$tip = d('<div class="tipsy"></div>').html('<div class="tipsy-arrow"></div><div class="tipsy-inner"><div class="inner-inner"></div></div></div>')
            }
            return this.$tip
        }, validate: function () {
            if (!this.$element[0].parentNode) {
                this.hide();
                this.$element = null;
                this.options = null
            }
        }, enable: function () {
            this.enabled = true
        }, disable: function () {
            this.enabled = false
        }, toggleEnabled: function () {
            this.enabled = !this.enabled
        }
    };
    d.fn.MyTitle = function (j) {
        if (j === true) {
            return this.data("tipsy")
        } else {
            if (typeof j == "string") {
                if (this.data("tipsy")) {
                    return this.data("tipsy")[j]()
                } else {
                    return null
                }
            }
        }
        j = d.extend({}, d.fn.MyTitle.defaults, j);
        function i(q) {
            var r = d.data(q, "tipsy");
            if (!r) {
                r = new c(q, d.fn.MyTitle.elementOptions(q, j));
                d.data(q, "tipsy", r)
            }
            return r
        }

        function n() {
            var q = i(this);
            q.hoverState = "in";
            if (j.delayIn == 0) {
                q.show()
            } else {
                setTimeout(function () {
                    if (q.hoverState == "in") {
                        q.show()
                    }
                }, j.delayIn)
            }
        }

        function g() {
            var q = i(this);
            q.hoverState = "out";
            if (j.delayOut == 0) {
                q.hide()
            } else {
                setTimeout(function () {
                    if (q.hoverState == "out") {
                        q.hide()
                    }
                }, j.delayOut)
            }
        }

        if (!j.live) {
            this.each(function () {
                i(this)
            })
        }
        if (j.trigger != "manual") {
            var e = j.live ? "live" : "bind", k = j.trigger == "hover" ? "mouseenter" : "focus", f = j.trigger == "hover" ? "mouseleave" : "blur";
            this[e](k, n)[e](f, g)
        }
        return this
    };
    d.fn.MyTitle.defaults = {
        defaultCss: " ",
        delayIn: 0,
        delayOut: 0,
        fade: false,
        fallback: "",
        gravity: "w",
        html: true,
        live: false,
        offset: 0,
        opacity: 1,
        title: "title",
        trigger: "hover"
    };
    d.fn.MyTitle.elementOptions = function (f, e) {
        return d.metadata ? d.extend({}, e, d(f).metadata()) : e
    };
    d.fn.MyTitle.autoNS = function () {
        return d(this).offset().top > (d(document).scrollTop() + d(window).height() / 2) ? "s" : "n"
    };
    d.fn.MyTitle.autoWE = function () {
        return d(this).offset().left > (d(document).scrollLeft() + d(window).width() / 2) ? "e" : "w"
    }
})(jQuery);
function changeTitle(c, a) {
    if (a.length > 0) {
        $(c).attr("errorStyle", a)
    } else {
        $(c).removeAttr("errorStyle")
    }
}
function redirecToWithReferer(c) {
    var a = document.createElement("a");
    a.href = c;
    document.body.appendChild(a);
    a.click()
}
jQuery.fn.FixedHeader = function (a) {
    var d = jQuery.extend({headerrowsize: 1, highlightrow: false, highlightclass: "highlight"}, a);
    this.each(function (j) {
        var q = $(this);
        var k = $(this).parent();
        var f = q.find("tr:lt(" + d.headerrowsize + ")");
        var g = "th";
        if (f.find(g).length == 0) {
            g = "td"
        }
        if (f.find(g).length > 0) {
            f.find(g).each(function () {
                $(this).css("width", $(this).width())
            });
            var n = q.clone().empty();
            var e = c(q);
            n.attr("id", "fixedtableheader" + j).css({
                position: "fixed",
                top: "0",
                left: q.offset().left
            }).append(f.clone()).width(e).hide().appendTo(k);
            if (d.highlightrow) {
                $("tr:gt(" + (d.headerrowsize - 1) + ")", q).hover(function () {
                    $(this).addClass(d.highlightclass)
                }, function () {
                    $(this).removeClass(d.highlightclass)
                })
            }
            $(window).scroll(function () {
                if (jQuery.browser.msie && jQuery.browser.version == "6.0") {
                    n.css({position: "absolute", top: $(window).scrollTop(), left: q.offset().left})
                } else {
                    n.css({position: "fixed", top: "0", left: q.offset().left - $(window).scrollLeft()})
                }
                var i = $(window).scrollTop();
                var r = f.offset().top;
                if (i > r && i <= (r + q.height() - f.height())) {
                    n.show()
                } else {
                    n.hide()
                }
            });
            $(window).resize(function () {
                if (n.outerWidth() != q.outerWidth()) {
                    f.find(g).each(function (r) {
                        var i = $(this).width();
                        $(this).css("width", i);
                        n.find(g).eq(r).css("width", i)
                    });
                    n.width(q.outerWidth())
                }
                n.css("left", q.offset().left)
            })
        }
    });
    function c(f) {
        var e = f.outerWidth();
        return e
    }
};
function encodeURI(a) {
    a = a + "";
    a = a.replaceAll("-", "__");
    a = a.replaceAll(".", "___");
    a = a.replaceAll("[.]", "___");
    return encodeURIComponent(a)
}
$.extend({
    loadScript: function (c) {
        var k = new Date().getTime();
        var g = typeof c == "string" ? [c] : c;
        var n = $("head");
        for (var a = 0; a < g.length; a++) {
            var e = document.getElementsByTagName("head")[0];
            var j = document.createElement("script");
            j.type = "text/javascript";
            j.id = k + "-" + a;
            j.src = g[a] + "?" + k;
            e.appendChild(j)
        }
    }, getServerDate: function (a) {
        return new Date(a.getResponseHeader("Date"))
    }, getServerTime: function (a) {
        return this.getServerDate(a).getTime()
    }
});
function Void() {
}
var urlreplace = "/ajaxList-";
function newUrl(d, a) {
    var c = $(a).attr("href");
    if (c && c.indexOf("javascript") != 0 && c.indexOf("-") > 0) {
        c = c.replace("-", urlreplace);
        ajaxUrl(c, true, false, "text")
    }
    if (d && d.preventDefault) {
        d.preventDefault()
    } else {
        window.event.returnValue = false
    }
    return false
}
function checkPageInput(c, a) {
    c.keyCode == 13 && $("#JumpButton").trigger("click")
}
function jumpPage(e) {
    var d = T$("PagerInput").value;
    if (d > e) {
        Wrong("page number is too large", {CloseTime: 1});
        return
    }
    if (d < 1) {
        if (d == "") {
            Info("please input page number", {CloseTime: 1})
        } else {
            Wrong("page number is too small", {CloseTime: 1})
        }
        return
    } else {
        var c = /^[0-9]+.?[0-9]*$/;
        if (c.test(d)) {
            var a = getPageUrl();
            if (a.length > 0) {
                ajaxUrl(a, true, true, "text")
            }
        } else {
            Wrong("Page number you entered is not valid, please re-enter.", {CloseTime: 1})
        }
    }
}
var chbtc = {
    cookiKeys: {
        uon: JsCommon.uon,
        uname: JsCommon.uname,
        uid: JsCommon.uid,
        aid: JsCommon.aid,
        rid: JsCommon.rid,
        aname: JsCommon.aname,
        note: JsCommon.note,
        lan: JsCommon.lan
    },
    mainDomain: JsCommon.mainDomain,
    vipDomain: JsCommon.vipDomain,
    p2pDomain: JsCommon.p2pDomain,
    transDomain: JsCommon.transDomain,
    staticDomain: JsCommon.staticDomain,
    ltcDomain: JsCommon.ltcDomain,
    first: 0,
    urlsAjax: {},
    ajax: function (d) {
        var c = true;
        if (d.needLoading == undefined) {
            c = true
        } else {
            c = d.needLoading == false ? false : true
        }
        var f = this;
        var a = d.url.indexOf("?") > -1 ? d.url.substring(0, d.url.indexOf("?")) : d.url;
        var e = f.urlsAjax[a];
        if (e == null) {
            f.urlsAjax[a] = {loading: false, needLoading: c, url: a};
            e = f.urlsAjax[a]
        }
        new this.ajaxDeal(d, e)
    },
    ajaxDeal: function (e, j) {
        if (!e.url) {
            Wrong("url参数必须传递！")
        }
        if (j.loading) {
            return
        }
        var i = this, a = "", c = e.div || null, d = e.needLogin || false, g = e.timeout || 60000;
        var f = true;
        if (e.async == undefined) {
            f = true
        } else {
            f = e.async == false ? false : true
        }
        i.dataType = e.dataType || "xml";
        i.dataType = i.dataType.toLowerCase();
        i.type = e.type || "post";
        if (j.needLogin && !chbtc.user.checkLogin()) {
            return
        }
        if (e.formId) {
            a = FormToStr(e.formId);
            if (a == null) {
                JuaBox.close();
                return
            }
        }
        j.loading = true;
        if (c && j.needLoading) {
            $("#" + c).Loadding({OffsetXGIF: 0, OffsetYGIF: 0})
        }
        $.ajax({
            async: f, cache: false, type: i.type, dataType: i.dataType, url: e.url, data: a, error: function (k) {
                j.loading = false;
                if (c) {
                    $("#" + c).Loadding({IsShow: false})
                }
            }, timeout: g, contentType: "application/x-www-form-urlencoded; charset=UTF-8", success: function (k) {
                j.loading = false;
                if (c) {
                    $("#" + c).Loadding({IsShow: false})
                }
                if (i.dataType == "xml") {
                    if ($(k).find("State").text() == "true") {
                        if (typeof e.suc == "function") {
                            (function (n) {
                                return e.suc(n)
                            })(k)
                        } else {
                            JuaBox.info($(k).find("Des").text())
                        }
                    } else {
                        if (typeof e.err == "function") {
                            (function (n) {
                                return e.err(n)
                            })(k)
                        } else {
                            JuaBox.info($(k).find("Des").text())
                        }
                    }
                } else {
                    if (i.dataType.indexOf("json") == 0) {
                        if (k.isSuc) {
                            if (typeof e.suc == "function") {
                                (function (n) {
                                    return e.suc(n)
                                })(k)
                            } else {
                                JuaBox.info(k.des)
                            }
                        } else {
                            if (typeof e.err == "function") {
                                (function (n) {
                                    return e.err(n)
                                })(k)
                            } else {
                                JuaBox.info(k.des)
                            }
                        }
                    } else {
                        if (typeof e.suc == "function") {
                            (function (n) {
                                return e.suc(n)
                            })(k)
                        }
                        if (typeof e.err == "function") {
                            (function (n) {
                                return e.err(n)
                            })(k)
                        }
                    }
                }
            }
        })
    },
    top: {
        init: function () {
            chbtc.user.init()
        }, close_city: function () {
            $("#selectcity_div").hide()
        }, show_city: function () {
            $("#selectcity_div").show()
        }, changeCity: function (a) {
            $(".tuangouflsy_pf_chengshi").each(function () {
                $(this).hide()
            });
            $("#city_span_" + a).show()
        }
    },
    tool: {
        param: function (a) {
            var c = window.location.search.substr(1).match(new RegExp("(^|&)" + a + "=([^&]*)(&|$)", "i"));
            if (c != null) {
                return unescape(c[2])
            }
            return null
        }, addBookmark: function (c) {
            var a = location.href;
            if (window.sidebar) {
                window.sidebar.addPanel(c, a, "")
            } else {
                if (document.all) {
                    window.external.AddFavorite(a, c)
                } else {
                    if (window.opera && window.print) {
                        return true
                    }
                }
            }
        }, initBackBtn: function () {
            var a = $.cookie(chbtc.cookiKeys.rid);
            if (a > 0) {
                var c = "/admin/Module/list.js";
                $.loadScript(c)
            }
        }, getTimeShowByMillSeconds: function (a) {
            var d = parseInt(a / (24 * 60 * 60 * 1000));
            var c = parseInt(a / (60 * 60 * 1000) - d * 24);
            c = c >= 10 ? c : "0" + c;
            var e = parseInt((a / (60 * 1000)) - d * 24 * 60 - c * 60);
            e = e >= 10 ? e : "0" + e;
            var f = parseInt(a / 1000 - d * 24 * 60 * 60 - c * 60 * 60 - e * 60);
            f = f >= 10 ? f : "0" + f;
            return "" + d + "天" + c + "小时" + e + "分" + f + "秒"
        }, isFloat: function (a) {
            if (a) {
                var c = /^[0-9]*\.?[0-9]*$/;
                if (!c.test(a)) {
                    return false
                }
            } else {
                return false
            }
            return true
        }, getHostName: function () {
            var a = location.href;
            return this.getHostNameByHref(a)
        }, getHostNameByHref: function (a) {
            try {
                var f = a.match(new RegExp("//(.*)[.]"))[1];
                var c = f.split(".");
                if (c.length > 0) {
                    return c[0]
                } else {
                    return null
                }
            } catch (d) {
            }
            return null
        }
    },
    user: {
        cookieInit: false, loginStatus: false, inAjaxing: false, lastPrice: 0, tickJson: null, isLogin: function () {
            if (!this.cookieInit) {
                this.init()
            }
            return this.loginStatus
        }, init: function () {
            var c = chbtc.getLan();
            $("#lanSelectA").attr("class", "da " + c);
            $("#languageSelectA").DropSelecter({
                ControlSeleter: "#languageSelectCont",
                Top: 0,
                Left: -5,
                StyleHover: "download_hover",
                call: function () {
                }
            });
            this.cookieInit = true;
            this.loginStatus = $.cookie(chbtc.cookiKeys.uon) == "1";
            var e = this;
            if (e.firstTitle.length == 0) {
                e.firstTitle = document.title
            }
            var d = document.location.pathname;
            if (e.loginStatus) {
            }
            setInterval(function () {
            }, 32000);
            if (e.loginStatus) {
                var a = $.cookie(chbtc.cookiKeys.uname);
                $(".logined .a1").html(a);
                $(".nologin").hide();
                $(".logined").show();
                $("#bwbank").show();
                $(".showSubUser").show()
            } else {
                $(".nologin").show();
                $(".logined").hide();
                $("#bwbank").hide();
                $(".showSubUser").hide()
            }
        }, ticker: function () {
            var a = this;
            $.getJSON(chbtc.mainDomain + "/api/ticker/topall?jsoncallback=?", function (c) {
                $(".realtime li span").each(function (d) {
                    $(this).text(c[d])
                })
            })
        }, showTime: function () {
            var c = new Date();
            var a = c.getHours();
            return a >= 0 && a < 6 ? chbtc.L("早上好") : a >= 6 && a < 12 ? chbtc.L("上午好") : a >= 12 && a < 13 ? chbtc.L("中午好") : a >= 13 && a < 18 ? chbtc.L("下午好") : chbtc.L("晚上好")
        }, firstTitle: "", lastP: [0, 0], upOrDown: function (a, e, d, c) {
            if (a != 0) {
                return
            }
            if (this.lastP[e] == 0) {
                this.lastP[e] = d[0]
            } else {
                if (d[0] > this.lastP[e]) {
                    $("#statisticsDiv li." + c).addClass("up").removeClass("down")
                } else {
                    if (d[0] < this.lastP[e]) {
                        $("#statisticsDiv li." + c).addClass("down").removeClass("up")
                    }
                }
                this.lastP[e] = d[0]
            }
        }, uticker: function () {
            if (this.isLogin()) {
                var a = this;
                $.getJSON(chbtc.mainDomain + "/user/userticker?callback=?", function (d) {
                    var c = $("#indexUserFunds p span");
                    if (c.length == 4) {
                        c.eq(0).text(d.total)
                    }
                    $("#finaPanelDown1 .rmb .d1 b").html(d.total)
                })
            }
        }, exeTickerJson: function () {
            if (this.tickJson != null) {
                var a = this.tickJson.ticker;
                if (($("#userPayAccountInfo").length > 0) && a != null) {
                    accountInfoShow(a.rmb, a.btcs, a.frmb, a.fbtc, a.buyBtc, a.sellRmb, a.total, a.buy, a.sell, a.btq, a.noBtc, a.nextBtc, a.ltc, a.fltc, a.buyLtc, a.sellLtc, a.nextLtcs)
                }
            }
        }, zcticker: function () {
            $.getJSON(chbtc.mainDomain + "/user/zcticker?callback=?", function (c) {
                var a = c.assets;
                $("#finaPanelDown1 .rmb b").each(function (d) {
                    $(this).text(a[d])
                });
                if ($("#myTotalBalance").length > 0) {
                    $("#myTotalBalance").text(a[0])
                }
            })
        }, resetUserFunds: function () {
            var a = this;
            $.getJSON(chbtc.mainDomain + "/u/resetUserFunds?callback=?", function (c) {
                if (c.isSuc) {
                    Right("资金刷新成功", {callback: "Close()"});
                    a.balance()
                } else {
                    Info(c.des, {callback: "Close()"})
                }
            })
        }, balance: function () {
            if (this.isLogin()) {
                $.getJSON(chbtc.mainDomain + "/u/getBalance?callback=?", function (c) {
                    var f = c.datas;
                    var e = Number(f[0]);
                    var a = Number(f[1]);
                    var d = Number(f[2]);
                    $("#myTotalBalance").text("฿ " + (e + a + d).toFixed(3));
                    $("#myRmbBalance").text("฿ " + e.toFixed(3));
                    $("#myBtcBalance").text("฿ " + a.toFixed(3));
                    $("#myLtcBalance").text("฿ " + d.toFixed(3))
                })
            }
        }, userInfo: {}, login: function (c, f, e) {
            var a = f || true;
            var d = c || "Close";
            Redirect("/user/log")
        }, reg: function (a, d, c) {
            Redirect("/user/reg")
        }, checkLogin: function (d, f, c) {
            this.loginStatus = $.cookie(chbtc.cookiKeys.uon) == "1";
            var a = f ? f : true;
            var e = c;
            if (!this.loginStatus) {
                this.login(d, a, e);
                return false
            }
            return true
        }, loginScuess: function (c, a) {
            if (a) {
                Close()
            } else {
                Right(chbtc.L("登陆成功！"))
            }
            $("#userNameCookie").text(c);
            this.loginStatus = true;
            this.init()
        }, getJuaBao: function () {
            if (this.checkLogin()) {
                JuaBox.frame(chbtc.mainDomain + "/u/juaBao?iframe=1", {width: 500})
            }
        }
    },
    form: {error: ""},
    getLan: function () {
        var a = $.cookie(chbtc.cookiKeys.lan);
        if (a == "en") {
            return "en"
        } else {
            if (a == "hk") {
                return "hk"
            } else {
                if (a == "ew") {
                    return "ew"
                } else {
                    if (a == "es") {
                        return "es"
                    }
                }
            }
        }
        return "cn"
    },
    L: function (c, a) {
        try {
            var g = "";
            if (this.getLan() == "cn") {
                g = this.tips[c][0]
            } else {
                if (this.getLan() == "en") {
                    g = this.tips[c][1]
                } else {
                    g = this.tips[c][2]
                }
            }
            if (a && a.length > 0) {
                for (var f in a) {
                    g = g.replace(a[f].k, a[f].v)
                }
            }
            return g
        } catch (d) {
            return c
        }
    },
    addTips: function (a, c) {
        if (this.getLan() == "cn") {
            chbtc.tips[a] = [c, ""]
        } else {
            if (this.getLan() == "en") {
                chbtc.tips[a] = ["", c]
            } else {
                chbtc.tips[a] = ["", "", c]
            }
        }
    },
    tips: {
        "登陆成功！": ["登陆成功！", "Login Success!", ""],
        "您确定注销登录么？": ["您确定注销登录么？", "Are you sure you want to log out!", ""],
        "请检查网络，可能是网络过慢导致超时或者远程服务出现故障!": ["请检查网络，可能是网络过慢导致超时或者远程服务出现故障!", "net errors!", ""],
        "确定": ["确定", "Ok", ""],
        "取消": ["取消", "Cancel", ""],
        "确定删除": ["您确定要删除该项吗？删除后无法恢复！", "Are you sure you want to delete that? Deleted can not be recovered!", ""],
        "确定要取消本次充值吗？": ["确定要取消本次充值吗？", "Sure you want to cancel the recharge?", ""],
        "确定删除本项吗？": ["确定删除本项吗？", "make sure delete this items?", ""],
        "确定取消吗？": ["确定取消吗？", "make sure cancle?", ""],
        "用户登录": ["用户登录", "User Login", ""],
        "登录": ["登录", "Login", ""],
        "注册": ["注册", "Sign up", ""],
        "收件箱": ["收件箱", "Inbox", ""],
        "未读邮件": ["未读邮件", "", ""],
        "请选择一项": ["请选择一项", "Please select a letter to be deleted!", ""],
        "确定要删除选中的项吗？": ["确定要删除选中的项吗？", "Are you sure you want to delete the selected messages?", ""],
        "开启": ["开启", "Open", ""],
        "关闭": ["关闭", "Close", ""],
        "系统忙碌，请稍候！": ["系统忙碌，请稍候！", "System is busy, please wait!", ""],
        "待成交": ["待成交", "Wating", ""],
        "已成交": ["已成交", "has", ""],
        "已取消": ["已取消", "canceled", ""],
        "查看原因": ["查看原因", "To see why.", ""],
        "您确定要开启安全密码吗？": ["您确定要开启安全密码吗？", "Are you sure you want to open the security code?", ""],
        "请正确输入您的安全密码,6-16位。": ["请正确输入您的安全密码,6-16位。", "Please enter your security code correctly, 6 to 16. ", ""],
        "暂时没有符合要求的记录": ["暂时没有符合要求的记录", "No records", ""],
        "您的X1余额：": ["您的X1余额：", "Your X1 balance:", ""],
        "地址已被复制到剪贴板，请核对：": ["地址已被复制到剪贴板，请核对：", "Address has been copied to the clipboard, please check it: ", ""],
        "BW提示您": ["提示", "JUA prompts you", ""],
        "请填写手机号码。": ["请填写手机号码。", "Please enter the phone number.", ""],
        "请填写您的安全密码。": ["请填写您的安全密码。", "Please enter your security code.", ""],
        "获取验证码": ["获取验证码", "Get verification code", ""],
        "x1秒后重新获取": ["x1秒后重新获取", "To obtain the x1 seconds", ""],
        "重新获取语音验证码": ["重新获取语音验证码", "Retrieve audio captcha", ""],
        "x1秒后重新获取语音验证码": ["x1秒后重新获取语音验证码", "x1 seconds to get voice verification code", ""],
        "提交成功": ["提交成功", "Submitted successfully", ""],
        "请选择要删除的比特币地址": ["请选择要删除的比特币地址", "Please select the currency you want to remove bits of address", ""],
        "删除比特币地址": ["删除比特币地址", "Delete bitcoin address", ""],
        "详细记录": ["详细记录", "All", ""],
        "": ["", "", ""],
        "": ["", "", ""],
        "": ["", "", ""]
    },
    setLan: function (a) {
        $.getJSON(chbtc.mainDomain + "/setlan?callback=?&lan=" + a, function (c) {
            if (c.isSuc) {
                location.reload()
            }
        })
    }
};
chbtc.list = {
    basePath: "", funcName: "", isInit: false, ui: function (c) {
        if (!c) {
            c = {}
        }
        var d = c.formId || "searchContaint";
        var a = $("#" + d);
        if (a.length > 0) {
            a.Ui()
        }
        this.pageInit();
        this.tabInit()
    }, tabInit: function () {
        var a = chbtc.tool.param("tab");
        if (a != null) {
            this.dealTab(a)
        }
    }, aoru: function (f) {
        if (!f) {
            f = {}
        }
        if (!f.id) {
            Wrong("ID参数必须传递！")
        }
        var g = f.id, c = f.width || 560, e = f.height || 360, d = f.title || "添加/编辑" + this.funcName, a = f.scroll || "auto";
        otherParam = f.otherParam || "", _this = this, url = f.url || _this.basePath + "aoru?id=" + g + otherParam;
        Iframe({Url: url, Width: c, Height: e, scrolling: a, Title: d})
    }, del: function (a) {
        if (!a) {
            a = {}
        }
        if (!a.id) {
            Wrong("ID参数必须传递！")
        }
        var e = a.id, c = a.otherParam || "", d = this;
        Ask2({
            Title: chbtc.L("确定删除本项吗？"), call: function () {
                chbtc.ajax({
                    url: d.basePath + "doDel?id=" + e + c, suc: function (f) {
                        d.reload();
                        Right($(f).find("Des").text())
                    }
                })
            }
        })
    }, dealTab: function (a) {
        $("#" + a).parent("li").parent("ul").find("li").removeClass("on");
        $("#" + a).parent("li").addClass("on")
    }, search: function (d) {
        if (!d) {
            d = {}
        }
        var r = d.formId || "searchContaint", c = d.div || this.defaultDiv, i = d.page || 1, f = "", j = d.special || false, q = d.validate || false, e = d.tab, a = d.url, g = this;
        var u = true;
        if (d.needLoading == undefined) {
            u = true
        } else {
            u = d.needLoading == false ? false : true
        }
        if ($("#" + r).length > 0) {
            if (e) {
                var n = $("#" + r).find("input[name='tab']");
                if (n.length > 0) {
                    n.val(e)
                } else {
                    $("<input name='tab' value='" + e + "' type='hidden'>").appendTo($("#" + r))
                }
                this.dealTab(e)
            }
            f = FormToStr(r);
            if (q) {
                if (f == null) {
                    return
                }
            } else {
                f = f == null ? "" : f
            }
            f = "?" + f + "&page=" + i
        }
        var k = "ajax" + f;
        chbtc.ajax({
            url: g.basePath + k, needLoading: u, suc: function (v) {
                if (j) {
                    $("#" + c).html(v)
                } else {
                    $("#" + c).html(v)
                }
                g.pageInit()
            }, dataType: "text", div: c
        })
    }, defaultDiv: "shopslist", defaultDivs: {}, ajaxPage: function (d, i) {
        if (d.needInit && !this.isInit) {
            return
        }
        if (!d) {
            d = {}
        }
        if (!d.url) {
            Wrong("url参数必须传递！")
        }
        var c = d.url, a = d.div, k = d.timeout || 60000, g = this;
        var f = c.indexOf("?") > -1 ? c.substring(0, c.indexOf("?")) : d.url;
        var j = g.defaultDivs[f];
        if (j == null && a) {
            g.defaultDivs[f] = {divId: a, url: f};
            j = g.defaultDivs[f]
        } else {
            if (j) {
                a = j.divId
            } else {
                a = g.defaultDiv
            }
        }
        var n = true;
        if (d.needLoading == undefined) {
            n = true
        } else {
            n = d.needLoading == false ? false : true
        }
        chbtc.ajax({
            url: c, needLoading: n, suc: function (q) {
                var e = q != $("#" + a).html();
                if (e) {
                    $("#" + a).html(q)
                }
                if (d.needInit) {
                    g.pageInit()
                }
                if (typeof d.suc == "function") {
                    (function () {
                        d.suc(e)
                    }())
                }
            }, dataType: "text", div: a, timeout: k
        });
        if (i) {
            if (i && i.preventDefault) {
                i.preventDefault()
            } else {
                window.event.returnValue = false
            }
            return false
        }
    }, jumpPage: function (d) {
        var c = $("#PagerInput").val();
        if (c > d) {
            JuaBox.info("您输入的页码过大，请重新输入。", {numSec: 1000, isSec: true});
            return
        }
        if (c < 1) {
            if (c == "") {
                JuaBox.info("请填写页码", {numSec: 1000, isSec: true})
            } else {
                JuaBox.info("您输入的页码太小了，请重新输入。", {numSec: 1000, isSec: true})
            }
            return
        } else {
            var a = /^[0-9]+.?[0-9]*$/;
            if (a.test(c)) {
                this.search({page: c})
            } else {
                JuaBox.info("您输入不是有效页码，请重新输入。", {numSec: 1000, isSec: true})
            }
        }
    }, resetForm: function (a) {
        if (!a) {
            a = {}
        }
        var c = "", d = this;
        c = a.formId || "searchContaint";
        $("#" + c).each(function () {
            this.reset();
            d.search()
        })
    }, config: function () {
        $(".item_list_bd").each(function (a) {
            $(this).mouseover(function () {
                $(this).css("background", "#fff8e1")
            }).mouseout(function () {
                $(this).css("background", "#ffffff")
            })
        })
    }, reload: function (a) {
        if (!a) {
            a = {}
        }
        a.page = $("#PagerInput").val() || 1;
        this.search(a)
    }, look: function (f) {
        if (!f) {
            f = {}
        }
        if (!f.url) {
            Wrong("url参数必须传递！")
        }
        var d = f.url, a = f.width || 820, c = f.height || 716, g = f.isShowIframeTile == undefined ? true : f.isShowIframeTile, i = f.scrolling || "auto", e = f.needLogin || false, j = this;
        if (e && !chbtc.user.checkLogin()) {
            return
        }
        Iframe({Url: d, Width: a, Height: c, scrolling: i, isShowIframeTitle: g, Title: "查看" + j.funcName})
    }, noPass: function (a) {
        Iframe({
            Url: "/admin/user/reason?id=" + a,
            zoomSpeedIn: 200,
            zoomSpeedOut: 200,
            Width: 600,
            Height: 460,
            Title: "填写不通过的原因"
        })
    }, url: function (c) {
        if (!c) {
            c = {}
        }
        if (!c.url) {
            Wrong("url参数必须传递！")
        }
        var a = c.url;
        chbtc.ajax({url: a})
    }, pageInit: function () {
        $("#ListTable").FixedHeader();
        $("#pagin a,#page_navA a").each(function () {
            $(this).AStop()
        });
        $("#PagerInput").UiText();
        this.isInit = true
    }, reloadAsk: function (a) {
        if (!a) {
            a = {}
        }
        if (!a.title) {
            Wrong("请传递您的询问标题！");
            return
        }
        if (!a.url) {
            Wrong("请传递您的后台url！");
            return
        }
        var c = this;
        Ask2({
            Msg: a.title, call: function () {
                chbtc.ajax({
                    url: a.url, suc: function (d) {
                        Right($(d).find("Des").text(), {
                            call: function () {
                                c.reload(a);
                                Close()
                            }
                        })
                    }, err: function (d) {
                        Wrong($(d).find("Des").text())
                    }
                })
            }
        })
    }
};
function CurrencyFormatted(c) {
    var a = parseFloat(c);
    if (isNaN(a)) {
        a = 0
    }
    var d = "";
    if (a < 0) {
        d = "-"
    }
    a = Math.abs(a);
    a = parseInt((a + 0.005) * 100);
    a = a / 100;
    s = new String(a);
    if (s.indexOf(".") < 0) {
        s += ".00"
    }
    if (s.indexOf(".") == (s.length - 2)) {
        s += "0"
    }
    s = d + s;
    return s
}
(function (a) {
    a.fn.hoverDelay = function (d) {
        var g = {
            hoverDuring: 200, outDuring: 200, hoverEvent: function () {
                a.noop()
            }, outEvent: function () {
                a.noop()
            }
        };
        var f = a.extend(g, d || {});
        var c, e;
        return a(this).each(function () {
            a(this).hover(function () {
                clearTimeout(e);
                c = setTimeout(f.hoverEvent, f.hoverDuring)
            }, function () {
                clearTimeout(c);
                e = setTimeout(f.outEvent, f.outDuring)
            })
        })
    }
})(jQuery);