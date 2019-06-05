(function(){
	var session = window.sessionStorage;
	window._sysVersion = 'V2.3.95';
	window._configPath = hierarchy+'hq/configuration.js?v='+Math.random();
	// 加载iscroll4所需css文件，优化在使用iscroll4时才加载样式文件，导致上下拉刷新有片刻无样式的支持
	addCss(["css/hsea.min.css", "css/iscroll4-custom.css"]);
	addJS(["js/zepto.js", "js/hsea.min.js"]);
	var request = new UrlSearch();
	if(!(request["openid"] === undefined)) request["clientid"] = request["openid"]; // 微信的标识key为openid
	if(!(request["clientid"] === undefined)) { // 处理刷新页面因没这个参数则会清理掉原有的同步标识
		var clientid = request["clientid"] || null;
		if(!clientid || clientid == "null" || clientid == "undefined") clientid = "";
		sessionStorage.setItem("clientid", clientid); // 获取同步自选股标识 一般为手机号 外部模块可能通过url链接带过来
	}
	sessionStorage.setItem("moduleName", request["moduleName"]); // 获取从外部模块到行情的模块名 用于绑定手机号去他们的页面 可不传
	var cssType = request["cssType"] || sessionStorage.getItem("cssType");
	sessionStorage.setItem("cssType", cssType || "");
	var modelName = "";
	try{
		modelName = model_name ? model_name + "_" : "";
	}catch( e ){
		
	}
	switch(cssType) { // 样式类型 用于换肤
		case "black": // 黑版行情
			addCss(["css/common.css", "css/animate.css", "hq/css/"+ modelName +"style_black.css"]);
			break;
		case "red": // 红版行情
			addCss(["css/common.css", "css/animate.css", "hq/css/"+ modelName +"style_red.css"]);
			break;
		case "blue": // 蓝版行情
			addCss(["css/common.css", "css/animate.css", "hq/css/"+ modelName +"style.css"]);
			break;
		// ... 可添加其它换肤样式
		default: // 默认蓝版行情
			addCss(["css/common.css", "css/animate.css", "hq/css/"+ modelName +"style_black.css"]);
			break;
	}
	
	function addCss(a){
		for(var i = 0; i < a.length; i++){
			var link = document.createElement("link");
			link.rel = "stylesheet";
			link.type = "text/css";
			link.charset= "utf-8";
			link.href = hierarchy + a[i] + "?v=" + _sysVersion;
			document.getElementsByTagName("head")[0].appendChild(link);
		}
	}
	
	function addJS(a){
		for(var i = 0; i < a.length; i++){
			var url = hierarchy + a[i] + "?v=" + _sysVersion;
			document.write(unescape("%3Cscript src='"+url+"' type='text/javascript'%3E%3C/script%3E"));
		}
	}
	
	function UrlSearch() { // 截取当前浏览器的url参数，并根据指定key获取参数值value
		var name, value;
		var str = location.href; //取得整个地址栏
		var num = str.indexOf("?")
		str = str.substr(num + 1); //取得所有参数   stringvar.substr(start [, length ]
		var arr = str.split("&"); //各个参数放到数组里
		for (var i = 0; i < arr.length; i++) {
			num = arr[i].indexOf("=");
			if (num > 0) {
				name = arr[i].substring(0, num);
				value = arr[i].substr(num + 1);
				this[name] = value;
			}
		}
	}
})()