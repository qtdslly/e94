define('hq/scripts/market/mainPanel.js',function(require, exports, module) { // 市场页面
	var gconfig = $.config;
	var global = gconfig.global;
	var VIscroll = require("vIscroll");
	var SHIscroll = require("shIscroll");
	var HIscroll = require("hIscroll");
	var layerUtils = require("layerUtils");
	var common = require("common");
	var service = require("hqService").getInstance();
	var informService = require("informService").getInstance();
	var interval = -1; // 一个以指定周期刷新返回的标识
	var newShareInterval = null; // 新股日历定时器
	var myVIscroll = null; // 垂直滑动组件
	var mySHIscroll = null; // 水平滑动组件 非分页式
	var myHIscroll = null; // 水平滑动组件 分页式
	var myHSHIscroll = null; // 水平沪深指数滑动组件 分页式
	var mainPanelExpResults = []; // 保存主面板指数查询结果集
	var mainPanelListResults = []; // 保存主面板列表查询结果集  mainPanelListResults[id] id:0-热门板块，1-领涨股，2-领跌股，3-基金，4-债券，5-资金流入，6-资金流出，7-换手率，8-港股通，9-沪股通，10-AH股，11-认购，12-认沽
	var currentCheckIndex = 0; // 当前选中市场的下标  与该页面的li一一对应   用于用户触发垂直滑动组件的同时也响应了分页式水平滑动组件的滑动完毕后的处理函数
	var _pageId = "#market_mainPanel ";

	function init() { // 初始化方法
		initMainPanel(); // 初始化面板样式
		initSHIScroll(); // 初始化非分页式左右滑动组件
		common.mainFooterFunc(); // 行情页面一级页面底部
//		initHIScroll(); // 初始化分页式左右滑动组件 初始化后可能影响到跳更多页面后返回不回来的问题
		initHSHIScroll(); // 初始化沪深指数水平滑组 
		initVIScroll(); // 初始化垂直滑动组件 防止在没有网络的情况下无法执行初始化操作
	}
	
	function load() {
		updatePanel(); // 加载数据和初始化变量 
//		getNewShareCal(); // 新股日历信息查询
		initHQUpdate(); // 行情定时刷新程序
		unbindPageEvent();
	}
	
	function unbindPageEvent() { // 事件用完后需要清理的方法
		$.bindEvent($(_pageId + ".fund_box li"), function(e) { // 主面板主导部分的点击 沪深指数 港股指数
			var id = $(this).parents(".fund_box").attr("id");
			var index = $(this).index();
			var param = {};
			if(mainPanelExpResults[id] && mainPanelExpResults[id].length > 0) {
				param.stockInfo = JSON.stringify(mainPanelExpResults[id][index]);
				if(id == 0) { // 指数
					$.pageInit("market/mainPanel", "market/expStockInfo", param);
				} else if(id == 1) { // 港股
					$.pageInit("market/mainPanel", "market/hkExpStockInfo", param);
				}
			}
			$(this).unbind(); // 删除对该事件的绑定 避免触发两次的情况
		});
	}

	function bindPageEvent() { // 页面元素事件绑定方法
		$.bindEvent($(_pageId + ".header_inner .button-group a:eq(0)"), function(e) { // 自选股与市场切换
			$.pageInit("market/mainPanel", "hqZxg/zxgList", null, "");
		});

		$.bindEvent($(_pageId + ".header_inner>a"), function(e) { // 刷新 搜索 菜单
			var cla = $(this).attr("class");
			if(cla == "icon_refresh") { // 刷新
				updatePanel(); // 更新数据
			} else if(cla == "icon_search") { // 搜索
				$.pageInit("market/mainPanel", "hqZxg/search");	
			}else if(cla == "icon_nav") { // 设置中心
				$.pageInit("market/mainPanel", "hqSet/set");	
			}
		});
		
		$.bindEvent($(_pageId + "footer #trade"), function() { // 切至交易模块
			common.hqExternalJump("3");
		});
		
		$.bindEvent($(_pageId + ".new_tab li"), function() { // 市场之间的切换
			if(!$(this).hasClass("current")) {
				$(this).addClass("current").siblings("li").removeClass("current");
				currentCheckIndex = $(this).index();
				var $current_fund_list = $(_pageId + ".tab-content:eq(" + currentCheckIndex + ")"); // 获取当前显示的列表模块
				$(_pageId + "#container_VIscroll").height($(_pageId + "article").height() - ($current_fund_list.find(".other_main").length ? 0 : 10));
				// 将未显示市场的模块赋予当前显示模块的高度并删除本身模块的高度，避免市场切换时显示模块高度被之前显示过的模块高度撑开导致上下滑动时高度错误引发的问题，且因上下滑动组件的原因需移除当前显示模块的高度
				$current_fund_list.height("").siblings(".tab-content").height($current_fund_list.height());
//				setTimeout(function() { // 处理tab切换时会偶尔scrollToElement执行失败的问题
//					mySHIscroll ? mySHIscroll.scrollToElement($(_pageId + ".new_tab li").eq(currentCheckIndex ? (currentCheckIndex - 1) : 0), 0) : "";
//				}, 80);
				mySHIscroll ? mySHIscroll.scrollToElement($(_pageId + ".new_tab li").eq(currentCheckIndex ? (currentCheckIndex - 1) : 0), 500) : "";
				myHIscroll ? myHIscroll.scrollToPage(currentCheckIndex) : $current_fund_list.show().siblings(".tab-content").hide(); // 跳至指定页
				/**
				 * scrollTo方法 将滑动到指定的地方
				 * 参数1 定位的横坐标
				 * 参数2 定位的纵坐标
				 * 参数3 定位的时间
				 * 参数4 是否以当前为原点，false表示以实际原点开始
				 * */
				myVIscroll.getWrapperObj() ? myVIscroll.getWrapperObj().scrollTo(0, -40, 0, false) : "";
				updatePanel(); // 更新数据
			}
		});

		$.bindEvent($(_pageId + ".home_box .title"), function(e) { // 点击分类最右边的更多和左边的箭头
			var childrenName = e.target.localName; // 获取事件的调用对象的名称
			if($(this).hasClass("calendarText")) return; // 新股日历点击
			if(childrenName == "h3") { // 左侧的箭头
				$(this).toggleClass("hide");
				$(this).next().slideToggle("500", initVIScroll);
			} else if(childrenName == "span" || childrenName == "a") { // 右侧的更多
				if($(this).next().find(".part").length > 0) { // 存在更多数据
					var id = $(this).parents(".home_box").attr("id");
					var title = $(this).children("h3").html();
					var param = {
							"id" : id,
							"title" : title
					};
					if(id == 0 || id == 90) { // 热门行业
						$.pageInit("market/mainPanel", "market/hqSortList", param);
					} else {
						$.pageInit("market/mainPanel", "market/hqList", param);
					}
				}
			}
		});

		$.bindEvent($(_pageId + ".calendarText .more"), function() { // 点击新股日历进入新股日历列表
			if($(_pageId + ".calendarText div").hasClass("active")) {
				$.pageInit("market/mainPanel", "market/newCalendarList");
			}
		});
		
		$.bindEvent($(_pageId + ".calendarText span"), function() { // 点击新股日历进入新股日历列表
			if($(_pageId + ".calendarText div").hasClass("active")) {
				$.pageInit("market/mainPanel", "market/newCalendarList");
			}
		});
		
		$.bindEvent($(_pageId + "#newCalendar a"), function() { // 沪深新股列点击
			if($(this).index() == 0 || $(this).hasClass("more")) $.pageInit("market/mainPanel", "market/newCalendarList"); // 新股日历
			else {	// 龙虎榜
				var url = global.znxgUrl + "/lhb/request.html";
				var paramArr = [];
				paramArr.push("isH5Hq=true");
				var cssType = sessionStorage.getItem("cssType")||"black";
				cssType ? cssType : "";
				paramArr.push("skin=" + cssType);
				window.location.href = url + "?" + paramArr.join("&");
			} 
		});

		$.preBindEvent($(_pageId + ".home_plate"), " .part", function(e) { // 点击某个分类的板块进入对应页面查询对应板块的成分股查询
			e.stopPropagation();
			var param = {
					"id" : $(this).parents(".home_plate").attr("id"),
					"bk_code" : $(this).attr("bk_code"),
					"title" : $(this).find("h5").text()
			};
			$.pageInit("market/mainPanel", "market/hqList", param);
		});

		$.preBindEvent($(_pageId + ".home_rank"), " .part", function() { // 点击某只股票进入股票详情页面
			var id = $(this).parents(".home_rank").attr("id");
			if(common.indexOfFunc(id, [8, 9, 22, 23, 24, 25, 26, 27, 28, 29], true)) { // 港股分类
				$.pageInit("market/mainPanel", "market/hkStockInfo", {"stockInfo" : JSON.stringify(mainPanelListResults[id][$(this).index()])});
			} else if(common.indexOfFunc(id, [60], true)){  // 指数模块
				$.pageInit("market/mainPanel", "market/expStockInfo", {"stockInfo" : JSON.stringify(mainPanelListResults[id][$(this).index()])});
			} else if(common.indexOfFunc(id, [91,92,93,94,95], true)){  //全球指数 (不跳转)
				
			}else if(id == 21) { //新股（新股多了一个title头，所以index需要减一）
				$.pageInit("market/mainPanel", "market/ggStockInfo", {"stockInfo" : JSON.stringify(mainPanelListResults[id][$(this).index() - 1])});
			}else {
				$.pageInit("market/mainPanel", "market/ggStockInfo", {"stockInfo" : JSON.stringify(mainPanelListResults[id][$(this).index()])});
			}
		});
		
		$.bindEvent($(_pageId + " .more_market .more_list a"),function() { // 其它按钮点击
			var stockType = $(this).attr("stockType"); // 股票类型
			var bkType = $(this).attr("bkType"); // 板块类型
			var otctype = $(this).attr("otctype"); // 新板块类型
			var id = $(this).attr("id"); // 个股期权id 11-认购 12-认沽 13-个股期权(新版)
			var param = {
					"stockType" : stockType,
					"bkType" : bkType,
					"id" : id,
					"title" : $(this).text(),
					"otctype" : otctype
			}
			if(id == 13) $.pageInit("market/mainPanel", "contract/contractIndex");
			else if(id == 114) $.pageInit("market/mainPanel", "market/reverseList"); // 国债逆回购
//			else if(id) $.pageInit("market/mainPanel", "contract/contractIndex", param);
			else if(id || stockType || otctype) $.pageInit("market/mainPanel", "market/hqList", param);
			else if(bkType == 0) $.pageInit("market/mainPanel", "market/allPlate");
			else if(bkType) $.pageInit("market/mainPanel", "market/hqSortList", param);
			else layerUtils.dialog(0, "此功能暂未开放");
		});
		
		// 手指触发滑动组件的动作记录判断
//		$.bindEvent($(_pageId + ".tab-content"), function(e) {
//			var startX = e.originalEvent.changedTouches[0].clientX;
//			var startY = e.originalEvent.changedTouches[0].clientY;
//			var hIscrollCss = $(_pageId + "#scroller_HIscroll").css("transform"); // 保存原来样式
//			console.dir(hIscrollCss);
//			$.bindEvent($(this), function(e) {
//				e.preventDefault();
//				var moveToX = Math.abs(e.originalEvent.changedTouches[0].clientX - startX);
//				var moveToY = Math.abs(e.originalEvent.changedTouches[0].clientY - startY);
//				if(moveToX < moveToY && moveToY > 50) { // 满足条件则禁掉水平滑动
//					console.dir("隐藏其它市场");
//					$(this).show().siblings(".tab-content").hide();
//					myHIscroll.getWrapperObj().hScroll = false;
//				} 
//				/*else {
//					console.dir("显示所有市场");
//					myHIscroll.getWrapperObj().hScroll = true;
//					$(_pageId + "#scroller_HIscroll").css("transform", hIscrollCss); // 还原样式
//					$(_pageId + ".tab-content").show();
//				}*/
//			}, "touchmove"); // 微信x5内核 touchmove 不持续触发的问题 可通过定时器在touchmove里面延迟touchend代码的触发
//			$.bindEvent($(this), function(e) { 
//				e.preventDefault();
//				console.dir("显示所有市场，还原水平滑动");
//				$(_pageId + ".tab-content").show();
//				$(_pageId + "#scroller_HIscroll").css("transform", hIscrollCss); // 还原样式
//				myHIscroll.getWrapperObj().hScroll = true; // 还原正常滑动
//			}, "touchend");
//		}, "touchstart");
	}
	
	/**************************** 函数方法开始 *******************************/
	/**
	 * 初始化页面样式的方法
	 * 锁定页面内容区高度 减去头部、底部
	 * 根据个数确定模块宽度
	 * 确定内容区的最小高度
	 * */
	function initMainPanel() {
		var liLength = $(_pageId + ".new_tab li").length; // 获取页面模块个数
//		if(liLength == 5 || liLength == 6) { // 当有5、6个模块时页面显示3个模块
//			liLength = 3;
//		} else if(liLength > 6) { // 超过7个模块时显示4个模块
//			liLength = 4;
//		}
		liLength = 4; // 项目要求显示4个
		$(_pageId + ".tab-nav").show();
		$(_pageId + ".new_tab li").width(Math.floor($(window).width() / liLength));
//		$(_pageId + ".new_tab li").css("min-width", $(window).width() / liLength);
	}

	function initSHIScroll() { // 初始化水平滑动组件 非分页式
		if(!mySHIscroll) {
			var config = {
					wrapper : $(_pageId + "#wrapper_SHIscroll"), // wrapper对象
					container : $(_pageId + "#container_SHIscroll") // container对象
			};
			mySHIscroll = new SHIscroll(config);
			// 非分页式滑动组件会给外层加上position：absolute属性，使得页面其它节点忽略它的高度和存在
			$(_pageId + ".new_tab").css("position", "initial"); // 组件加的绝对定位属性，建议由组件去除
		}
	}

	function initHIScroll() { // 初始化水平滑动组件 分页式
		if(!myHIscroll) {
			$(_pageId + ".tab-content").css("min-height", $(_pageId + "article").height() - 10); // 10为内补丁高度
			var config = {
					wrapper : $(_pageId + "#wrapper_HIscroll"), // wrapper对象
					scroller : $(_pageId + "#scroller_HIscroll"), // scroller对象
					perCount : 1, // 每个可视区域显示的子元素，就是每个滑块区域显示几个子元素
					showTab : true, //是否有导航点
					tabDiv : $(_pageId + "#container_SHIscroll"), // 导航点集合对象
					auto : false, //是否自动播放
					onScrollEnd : function() { // 滑动完毕后的自定义处理函数
//						$(_pageId + ".new_tab li").removeClass("current"); // 去掉由该滑动组件给导航点点增加的样式 建议由组件去除
						var index = $(_pageId + ".new_tab li.current").index(); // 获取滑动后导航点的index 也可通过myHIscroll.getActivePage()获取
						if(index != currentCheckIndex) { // 表示用户通过水平滑动组件切换了市场，因为用户在触发垂直滑动组件的同时也会触发响应这个函数可能导致没有切换市场而请求了没必要的接口
							currentCheckIndex = index;
							var $current_fund_list = $(_pageId + ".tab-content:eq(" + currentCheckIndex + ")"); // 获取当前显示的列表模块
							$(_pageId + "#container_VIscroll").height($(_pageId + "article").height() - ($current_fund_list.find(".other_main").length ? 0 : 10));
							// 将未显示市场的模块赋予当前显示模块的高度并删除本身模块的高度，避免市场切换时显示模块高度被之前显示过的模块高度撑开导致上下滑动时高度错误引发的问题，且因上下滑动组件的原因需移除当前显示模块的高度
							$current_fund_list.height("").siblings(".tab-content").height($current_fund_list.height()); 
							mySHIscroll ? mySHIscroll.scrollToElement($(_pageId + ".new_tab li").eq(currentCheckIndex ? (currentCheckIndex - 1) : 0), 300) : "";
							myVIscroll.getWrapperObj() ? myVIscroll.getWrapperObj().scrollTo(0, -40, 0, false) : "";
							updatePanel(); // 更新数据
						}
					}
			};
			myHIscroll = new HIscroll(config);
//			$(_pageId + ".new_tab li").removeClass("current"); // 去掉由该滑动组件给导航点点增加的样式 建议由组件去除
			$(_pageId + ".tab-content").show();
			$(_pageId + "#scroller_HIscroll").height(""); // 去掉由水平滑动组件给div加的高度 建议由组件去除
			$(_pageId + ".tab-content").height(""); // 去掉由水平滑动组件给div加的高度 建议由组件去除
		}
	}
	
	function initHSHIScroll() { // 初始化沪深指数水平滑动组件 分页式
		if(!myHSHIscroll) {
			var config = {
					wrapper : $(_pageId + "#wrapper_HS_HIscroll"), // wrapper对象
					scroller : $(_pageId + "#scroller_HS_HIscroll"), // scroller对象
					perCount : 3, // 每个可视区域显示的子元素，就是每个滑块区域显示几个子元素
					showTab : true, //是否有导航点
					tabDiv : $(_pageId + "#tag_HS_HIscroll"), // 导航点集合对象
					auto : false, //是否自动播放
					onScrollEnd : function() { // 滑动完毕后的自定义处理函数
						
					}
			};
			myHSHIscroll = new HIscroll(config);
		}
	}

	function initVIScroll() { // 垂直滑动组件
		if(!myVIscroll) {
			var config = {
					"isPagingType" : false,
//					"isvScrollbar":false,
					"visibleHeight" : $(_pageId + "article").height() - 10, 
					"container" : $(_pageId + "#container_VIscroll"), // wrapper对象
					"wrapper" : $(_pageId + "#wrapper_VIscroll"), // scroller对象
					"downHandle" : function() {
						updatePanel();
					},
					"oPullDownTips" : {
						"still" : "下拉刷新", // 静止时的文本
						"flip" : "释放刷新", // 上拉过程中的文本
						"loading" : "正在加载..." // 加载中的文本
					}
			};
			myVIscroll = new VIscroll(config);
			$(_pageId + ".visc_pullUp").css("display", "none");
		} else {
			 // 每次刷新时都有可能使滑动内容区高度改变，所以需要将滑动区非显示的其它市场高度改变，不然就算现实隐藏模块高度也无法改变
			var $current_fund_list = $(_pageId + ".tab-content:eq(" + currentCheckIndex + ")"); 
			$current_fund_list.height("").siblings(".tab-content").height($current_fund_list.height());
			myVIscroll.getWrapperObj() ? myVIscroll.refresh() : "";
		}
	}

	function initHQUpdate() { // 行情定时刷新函数
		var hqUpdateTime = global.hqUpdateTime ? global.hqUpdateTime : 10; // 行情定时更新频率
		interval = setInterval(function() {
			var id = $(_pageId + ".new_tab li.current").attr("id");
			var isTrade = null;
			if(id == 1 || id == 10) { // 港股通、港股行情
				isTrade = common.nativePluginServiceFunction(2, 1, "isHKTrade", null);
			} else if(id == 0) { // 沪深A股行情
				isTrade = common.nativePluginServiceFunction(2, 1, "isTrade", null);
			} else if(id == 4) { // 全球行情
				isTrade = "true";
			}
			if(isTrade == "true") {
				updatePanel();
			}
		}, hqUpdateTime * 1000);
//		common.screenOrientationListener(function() { // 监听手机屏幕旋转时触发的方法
//			myVIscroll ? myVIscroll.destroy() : null; // 旋转后需清除滑动组件并重新定制页面
//			myHIscroll ? myHIscroll.destroy() : null; 
////			mySHIscroll ? mySHIscroll.destroy() : null; // 非分页式滑动组件经测试不需要重新初始化
//			myVIscroll = null;
//			myHIscroll = null;
////			mySHIscroll = null;
//			initMainPanel(); // 初始化面板样式
//			initSHIScroll(); // 初始化分页式左右滑动组件
//			global.isMainPanelHIscroll && initHIScroll(); // 初始化分页式左右滑动组件
//			initVIScroll(); // 初始化垂直滑动组件 防止在没有网络的情况下无法执行初始化操作
////			updatePanel(); // 加载数据和初始化变量 
//			
//		});
	}

	function updatePanel() { // 加载数据
		initVIScroll(); // 处理因网络等原因导致接口不通时无法回调执行刷新滑动组件的问题导致页面切换tab时高度不变的情况
		var id = $(_pageId + ".new_tab li.current").attr("id");
		switch(id) {
			case "0": // 沪深市场
				mainPanelLeadPartFunc(global.HSExpList, 0); // 查询沪深指数数据 0表示沪深市场
				getHotPlateResults(false); // 查询热门板块结果集
				getRankList(1, 0, 1); // 根据涨跌幅降序排列 领涨股查询
				getRankList(1, 1, 2); // 根据涨跌幅升序排序 领跌股查询
				getMoneyFlowsList(1, 0); // 资金流入查询
				getMoneyFlowsList(0, 1); // 资金流出查询
				getRankList(8, 0, 7); // 根据换手率降序排序 换手率榜查询
				break;
			case "1": // 港股通市场
				getGgtFundLimit(); //查询港股通可用额度和资金流入
				getHKTList(); // 查询港股通列表
				getRankList(1, 0, 9 ,98); // 查询深港通列表
				break;
			case "3": // 美股市场
	//			mainPanelLeadPartFunc("", 3); // 查询美股指数数据 3表示美股市场
				initVIScroll();
				break;
			case "4": // 全球市场
				getGlobalList("110", 91); //重要指数
				getGlobalList("111", 92); //欧洲
				getGlobalList("112", 93); //亚洲
				getGlobalList("113", 94); //美洲
				getGlobalList("115", 100); //非洲
				getGlobalList("114", 101); //澳洲
//				getGlobalList("114:115", 95); //其他
				getGlobalList("117", 96); //热门汇率
				getGlobalList("119", 97); //人民币中间价
				break;
			case "6": // 新股
				getNewStockList(); // 新股行情查询
				break;
			case "8": // 债券
				getFundOrBondSortList(4, "21:25"); // 国债
				getFundOrBondSortList(41, "22:26"); // 企债
				getFundOrBondSortList(42, "23:27"); // 可转债
				getFundOrBondSortList(43, "24:30"); // 回购
				break;
			case "10": // 港股
				mainPanelLeadPartFunc(global.HKExpList, 1); // 查询港股指数数据 1表示港股市场
				getHotPlateResults(true); // 查询港股热门板块结果集
				getRankList(1, 0, 22, 102, true); // 港股主板涨幅榜查询
				getRankList(1, 1, 23, 102, true); // 港股主板跌幅榜查询
				getRankList(14, 0, 24, 102, true); // 港股主板成交额榜查询
				getRankList(1, 0, 25, 104, true); // 港股创业板涨幅榜查询
				getRankList(1, 1, 26, 104, true); // 港股创业板跌幅榜查询
				getRankList(14, 0, 27, 104, true); // 港股创业板成交额榜查询
				getRankList(14, 0, 28, 101, true); // 港股认股证成交额榜查询
				getRankList(14, 0, 29, 105, true); // 港股牛熊证成交额榜查询
				break;
			case "12": // 指数
				getFundOrBondSortList(60, "7:15"); // 沪深指数
//				getFundOrBondSortList(62, "42"); // 三板指数
//				getFundOrBondSortList(61, 15); // 深圳指数
				break;
			case "13": // 其他
				initVIScroll();
				break;
		}
	}

	function mainPanelLeadPartFunc(stock_list, id) { // 查询主面板主导部分的数据 包括沪深指数、港股指数、期权标的
		var param = {
				"stock_list" : stock_list,
				"field" : "1:2:3:21:22:23:24"
		};
		var getStockListCallBack = function(data) {
			mainPanelLeadPartFuncCallBack(data, id);
		};
		if(id == 0){ // 沪深指数
			service.getStockList(param, getStockListCallBack);
		}else if(id == 1){ // 港股指数
			service.getHKZSList(param, getStockListCallBack);
		}
	}
	
	function getGgtFundLimit(){ //查询港股通可用额度和资金流入
		if($(_pageId + ".fund_box_hk").length && $(_pageId + ".fund_box_hk").css("display") != "none") {
			service.getGgtFundLimit(null , function(data){
				if(data){
					if(data.errorNo == 0 ){
						var results = data.results[0];
						
						var flowFund1 = results[0] - results[1];//沪港通初始额度-沪港通日中剩余额度
						var flowFund2 = results[4] - results[5];//深港通初始额度-深港通日中剩余额度
						flowFund1 >= 0 ? $(_pageId + ".fund_box_hk .row-1").eq(0).find("p").eq(1).html("资金流入"): $(_pageId + ".fund_box_hk .row-1").eq(0).find("p").eq(1).html("资金流出");
						flowFund2 >= 0 ? $(_pageId + ".fund_box_hk .row-1").eq(1).find("p").eq(1).html("资金流入"): $(_pageId + ".fund_box_hk .row-1").eq(1).find("p").eq(1).html("资金流出");
						
						var flowFund1 = common.judgeColorValue(3, flowFund1 > 0 ? flowFund1 : -flowFund1);
						var flowFund2 = common.judgeColorValue(3, flowFund2 > 0 ? flowFund2 : -flowFund2);
						
						$(_pageId + ".fund_box_hk .row-1").eq(0).find("strong").eq(0).html(common.judgeColorValue(3, results[1]) ? common.judgeColorValue(3, results[1]) : "--");
						$(_pageId + ".fund_box_hk .row-1").eq(1).find("strong").eq(0).html(common.judgeColorValue(3, results[5]) ? common.judgeColorValue(3, results[5]) : "--");
						
						$(_pageId + ".fund_box_hk .row-1").eq(0).find("strong").eq(1).html(flowFund1 ? flowFund1 : "--");
						$(_pageId + ".fund_box_hk .row-1").eq(1).find("strong").eq(1).html(flowFund2 ? flowFund2 : "--");
					}
				}
			});
		}
	}

	function getFundOrBondSortList(_id, type) { // 基金、债券、股转分类、指数查询
		var param = {
				"type" : type, // 查询类型
				"sort" : 1, // 排序类型
				"order" : 0, // 排序方向 默认降序
				"curPage" : 1, // 当前页 默认第1页
				"rowOfPage" : global.defaultCount, // 查询条数
				"field" : "1:2:3:21:22:23:24"
		};
		var getFundOrBondSortListCallBack = function(data) {
			callBackFunction(data, _id);
		};
		service.getRankList(param, getFundOrBondSortListCallBack);
	}
	
	function getNewShareCal() { // 查询新股日历信息
		informService.getNewShareCal({}, getNewShareCalCallBack);
	}

	function getHotPlateResults(isHG) { // 查询热门板块结果集
		var param = {
				"bkType" : 1, // 默认传1，显示行业板块
				"sort" : 1, // 默认传1，按照板块的涨跌幅排序
				"order" : 0, // 默认传0，按照降序查询
				"curPage" : 1, // 默认传1，从第一页开始查询
				"rowOfPage" : 6, // 查询条数，沪深面板显示6条数据
				"field" : "1:3:22:24:38:39:40:41:42" // 板块涨跌幅:板块涨跌:板块名称:板块代码:领涨股票名称:领涨股票涨跌幅:领涨股票现价:领涨股票涨跌:领涨股票类型
		};
		var callback = function(data) {
			getHotPlateResultsCallBack(data, isHG);
		}
		if(!isHG) {
			service.getHotPlateResults(param, callback);
		} else {
			service.getHKHotPlateResults(param, callback);
		}
		
	}

	function getRankList(sort, order, _id, type, isHG) { // 查询领涨股榜、领跌股榜、换手率榜
		var param = {
				"sort" : sort, // 排序字段
				"order" : order, // 排序方向
				"curPage" : 1, // 默认传1，从第一页开始查询
				"rowOfPage" : global.defaultCount, // 查询条数
				"type" : type || global.HSListType, // 查询列表的类型
				"field" : "1:2:3:21:22:23:24:31:8:14"
		}
		var getRankListCallBack = function(data) {
			callBackFunction(data, _id);
		};
		if(isHG) { // 港股板块查询
			service.getHKRankList(param, getRankListCallBack);
		} else {
			service.getRankList(param, getRankListCallBack);
		}
	}

	function getGlobalList(type, id){ //全球
		var param = {
				"sort" : 1, // 排序字段
				"order" : 0, // 排序方向
				"rowOfPage" : global.defaultCount, // 查询条数
				"curPage" : 1, // 默认传1，从第一页开始查询
				"field" : "1:2:3:21:22:24",
				"type" : type // 查询列表的类型
		}
		
		var getGlobalListCallBack = function(data) {
			getGlobalListCallBackFunction(data, id);
		};
		service.getGlobalList(param, getGlobalListCallBack);
	}
	
	function getMoneyFlowsList(ioType, sortType) { // 资金流入、流出列表查询
		var param = {
//				"ioType" : 2, // 2表示净流向排序
//				"sortType" : sortType, // 排序类型
				"ioType" : ioType, // 流入、流出类型
				"sortType" : 0, // 排序类型 当用总流向排序时默认降序
				"trading_day" : 0, // 交易日期 默认为0，表示当天
				"curpage" : 1, // 当前页 默认第1页
				"rowofpage" : global.defaultCount // 每页数量
		};
		var callback = function(data) {
			getMoneyFlowListCallBack(data, ioType);
		};
		service.getMoneyFlowsList(param, callback);
	}
	
	function getHKTList() { // 查询港股通列表数据
		var param = {
				"sort" : 1, // 排序类型
				"order" : 0, // 排序方向 默认降序
				"curPage" : 1, // 当前页 默认第1页
				"rowOfPage" : global.defaultCount, // 每页数量
				"field" : "1:2:3:21:22:23:24"
		};
		var getHKTListCallBack = function(data) {
			callBackFunction(data, 8);
		}
		service.getHKTList(param, getHKTListCallBack);
	}
	
	function getNewStockList() { // 新股行情查询
		var param = {
				"sort" : 62, // 排序类型
				"order" : 0, // 排序方向 默认降序
				"curPage" : 1, // 当前页 默认第1页
				"rowOfPage" : global.defaultCount, // 每页数量
				"field" : "1:2:3:21:22:23:24:62"
		};
		var getNewStockListCallBack = function(data) {
			callBackFunction(data, 21);
		}
		service.getNewStockList(param, getNewStockListCallBack);
	}

	function listIsDataFunc(id) { // 根据id获取变量时候有数据并展示在页面上
		if(mainPanelListResults[id] && mainPanelListResults[id].length > 0) {
			var stockArray = common.addListData(id, mainPanelListResults[id]); 
			if(id == 21){ //新股显示title
				$(_pageId + ".home_box[id='" + id + "']").find(".fund_table").html("<dl class=\"title\"><dt>名称代码</dt><dd>最新</dd><dd>累计涨幅</dd></dl>" + stockArray.join(""));
			}else{
				$(_pageId + ".home_box[id='" + id + "']").find(".fund_table").html(stockArray.join(""));
			}
			
		} else {
			$(_pageId + ".home_box[id='" + id + "']").find(".fund_table").html("<div style=\"height: 0.45rem;line-height: 0.45rem;text-align: center;\">暂无数据</div>");
		}
	}
	/**************************** 函数方法结束 *******************************/

	/**************************** 回调函数开始 *******************************/
	function mainPanelLeadPartFuncCallBack(data, id) { // 主面板主导部分的的回调函数 包括沪深指数、港股指数 的回调
		if(data) {
			if(data.errorNo == 0) {
				mainPanelExpResults[id] = []; // 清空之前保存的查询值
				if(data.results && data.results.length > 0) {
					for (var i = 0; i < data.results.length; i++) {
						var stock = data.results[i];
						var expLi = $(_pageId + ".fund_box[id='" + id + "']").find("li:eq(" + i + ")");
						var claColor = common.judgeColorValue(2, stock[2]);
						var point = common.typePoint(stock[3]); // 根据股票类型判断价格位数
						var zdf = (stock[0] * 100).toFixed(2);
						var zde = common.judgeColorValue(1, stock[2].toFixed(point));
						zdf = common.judgeColorValue(1, zdf) + "%";
						expLi.find("h5").html(stock[4]); // 股票名称  期权标的页面没有名称 需要显示
						expLi.find("strong").removeClass().html(stock[1].toFixed(point)).addClass(claColor); // 现价
//						expLi.find("span").html(zde +　"&nbsp;&nbsp;" + zdf); // 涨跌额
						expLi.find("p").removeClass().addClass(claColor).find(".lt").html(zde); // 涨跌额
						expLi.find(".rt").html(zdf); // 涨跌幅
						var expStock = common.unifyPackParam(i, {}, stock); // 统一封装固定参数
						mainPanelExpResults[id].push(expStock); // 存值
					}
				}
			} else {
				layerUtils.dialog(0, data.errorInfo);
			}
		}
	}

	function getNewShareCalCallBack(data) { // 新股日历信息查询的回调函数
		if(data) {
			if(data.error_no == 0) {
				var newShareArray = [];
				if(data.results && data.results.length > 0) {
					var newShareInfo = data.results[0];
					if(newShareInfo.total_subNewStockList != "0") { // 今日申购股票个数
						newShareArray.push("<div>今日" + newShareInfo.total_subNewStockList + "只新股可申购</div>");
					}
					if(newShareInfo.total_list != "0") { // 今日上市股票个数
						newShareArray.push("<div>今日" + newShareInfo.total_list + "只新股上市</div>");
					}
					if(newShareInfo.total_refundment != "0") { // 今日中签股票个数
						newShareArray.push("<div>今日" + newShareInfo.total_refundment + "只新股中签公布</div>");
					}
					if(newShareInfo.total_noList != "0") { // 即将上市股票个数
						newShareArray.push("<div>" + newShareInfo.total_noList + "只新股即将上市</div>");
					}
					if(newShareInfo.total_issueList != "0") { // 即将发行股票个数
						newShareArray.push("<div>" + newShareInfo.total_issueList + "只新股即将发行</div>");
					}
					if(newShareArray.length > 0) {
						var index = $(_pageId + ".calendarText span div[class=active]").index(); // 获取上一次接口响应时被选中的下标，第一次则为0
						index = index >= 0 ? index : 0; // 第一次index为-1，赋值为0
						$(_pageId + ".calendarText span").html(newShareArray.join(""));
						$(_pageId + ".calendarText span div:eq(" + index + ")").show().siblings().hide();
						$(_pageId + ".calendarText span div:eq(" + index + ")").addClass("active").siblings().removeClass("active");
						var length = $(_pageId + ".calendarText span div").length; 
						if(length > 1) { // 大于两个，需要轮播
							clearInterval(newShareInterval);
							newShareInterval = null;
							newShareInterval = setInterval(function () {
								if(index == length){
									index = 0;
								}
								$(_pageId + ".calendarText span div:eq(" + index + ")").show().siblings().hide();
								$(_pageId + ".calendarText span div:eq(" + index+ ")").addClass("active").siblings().removeClass("active");
								index++;
							}, 3000);
						}
					} else {
						$(_pageId + ".calendarText span").html("<div>无新股发行或上市</div>");
					}
					
					if(newShareInfo.total_noList != "0"){
						$(_pageId + "#newCalendar a em").html(newShareInfo.total_subNewStockList).show();
					}else{
						$(_pageId + "#newCalendar a em").html(" ").hide();
					}
					//$(_pageId + "#newCalendar a em").html("今日可申购"+newShareInfo.total_subNewStockList+"只");
				} else {
					$(_pageId + ".calendarText span").html("<div>无新股发行或上市</div>");
				}
			} else {
				layerUtils.dialog(0, data.error_info);
			}
		}
	}

	function getHotPlateResultsCallBack(data, isHG) { // 查询热门板块结果集
		if(data) {
			if(data.errorNo == 0) {
				if(data.results && data.results.length > 0) {
					var _results = [];
					for (var i = 0; i < data.results.length; i++) {
						var stock = data.results[i];
						var bk_upercent = (stock[0] * 100).toFixed(2); // 板块涨跌幅
						var point = common.typePoint(stock[8]); // 根据领涨股类型判断位数
						var lzg_now = stock[6].toFixed(point); // 领涨幅现价
						var lzg_upercent = (stock[5] * 100).toFixed(2); // 领涨股票涨跌幅
						var spanFontSize = common.judgeStrGetLenthFun(stock[2]) > 12 ? 0.12 : ""; // 小字体 超过10个字符
						var bkName = common.nameMecismFunc(stock[2], spanFontSize, null, 0, $(window).width() * 1 / 3);
						var hotPlateLi = "<li class=\"part\" bk_code=" + stock[3] + ">" +
						"<a><h5 style=\"" + (spanFontSize ? ("font-size:" + spanFontSize + "rem;") : "") + "\">" + bkName + "</h5>" +
						"<strong class=" + common.judgeColorValue(2, stock[1]) + ">" + common.judgeColorValue(1, bk_upercent) + "%</strong>" +
						"<p>" + stock[4].substring(0, 4) + "<em class=" + common.judgeColorValue(2, stock[7]) + ">" + common.judgeColorValue(1, lzg_upercent) + "%</em></p>" +
						"</a></li>";
						_results.push(hotPlateLi);
					}
					$(_pageId + ".home_plate[id='" + (isHG ?　90 : 0) + "'] .hot").html(_results.join(""));
				} else {
					$(_pageId + ".home_plate[id='" + (isHG ?　90 : 0) + "'] .hot").html("<div style=\"height: 0.45rem;line-height: 0.45rem;text-align: center;\">暂无数据</div>");
				}
			} else {
				layerUtils.dialog(0, data.errorInfo);
			}
			initVIScroll();
		}
	}

	function getMoneyFlowListCallBack(data, ioType) { // 资金流向的回调函数
		if(data) {
			if(data.errorNo == 0) {
				var id = ioType ? 5 : 6;
				mainPanelListResults[id] = [];
				if(data.results && data.results.length > 0) {
					for (var i = 0, len = data.results.length; i < len; i++) {
						var stock = data.results[i];
						var newStock = {
								"id" : i,
								"uppercent" : stock.uppercent, // 涨跌幅
								"now" : stock.now, // 现价
								"up" : stock.up, // 涨跌额
								"stockType" : stock.stktype, // 股票类型
								"stockCode" : stock.stockcode, // 股票代码
								"stockName" : stock.stockname, // 股票名称
								"market" : stock.market, // 股票市场
								"net_money_flow" : Number(stock.net_money_flow) * 10000 // 流入、流出资金  因资金流向返回的金额是以万为单位的，所以需转换
						};
						mainPanelListResults[id].push(newStock);
					}
				}
				listIsDataFunc(id);
			} else {
				layerUtils.dialog(0, data.errorInfo);
			}
			initVIScroll();
		}
	}
	
	function callBackFunction(data, id) { // 查询领涨股榜、领跌股榜、基金、债券列表、换手率榜、港股通列表、期权标的合约的列表的回调函数
		if(data) {
			if(data.errorNo == 0) {
				mainPanelListResults[id] = [];
				if(data.results && data.results.length > 0) {
					for (var i = 0, len = data.results.length; i < len; i++) {
						var stock = data.results[i];
						var newStock = {
								"ljzdf": stock[7] ||　"", //累积涨跌幅 
								"hsl" : stock[8] || "", // 换手率
								"cje" : stock[9] || "" //成交额
						};
						newStock = common.unifyPackParam(i, newStock, stock); // 统一封装固定参数
						mainPanelListResults[id].push(newStock);
					}
				} 
				listIsDataFunc(id);
			} else {
				layerUtils.dialog(0, data.errorInfo);
			}
			initVIScroll();
		}
	}
	
	function getGlobalListCallBackFunction(data, id){
		if(data) {
			if(data.errorNo == 0) {
				mainPanelListResults[id] = [];
				if(data.results && data.results.length > 0) {
					for (var i = 0, len = data.results.length; i < len; i++) {
						var stock = data.results[i];
						var globalStock = {
								"id" : i,
								"uppercent" : stock[0],
								"now" : stock[1],
								"up" : stock[2],
								"stockType" : stock[3],
								"stockName" : stock[4],
								"stockCode" : stock[5]
						};
						mainPanelListResults[id].push(globalStock);
					}
				} 
				listIsDataFunc(id);
			} else {
				layerUtils.dialog(0, data.errorInfo);
			}
			initVIScroll();
		}
	}
	/**************************** 回调函数结束 *******************************/

	function destroy() { // 销毁方法
		window.clearInterval(interval); // 销毁掉当前页面的定时刷新机制
		interval = -1;
		$(_pageId + ".draver").css("left","-"+$(_pageId + ".draver").css("width")).hide().removeClass("show");
		$(_pageId + ".menuUl").hide();
	}

	var mainPanel = {
			"init" : init,
			"load" : load,
			"bindPageEvent": bindPageEvent,
			"destroy": destroy
	};
	module.exports = mainPanel;
});