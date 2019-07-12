/**
 * Created by ccgo on 2019/6/28.
 */


function GetGoodsList(page,category){
  if($("#hidCategory").val() != category){
    $("#hidPage").val("1");
    $("#goods_list").html("");
  }
  $("#hidCategory").val(category);

  $.getJSON("/cms/goods/list?page="+page+"&size=20&category="+category, function (json) {
    var list = json.tbk_dg_item_coupon_get_response.results.tbk_coupon;

    var htmlStr = "";
    for (var i = 0; i < list.length; i++) {
      htmlStr += '<a href="javascript:;" onclick="getTpwd(\'' + list[i].coupon_click_url + '\',\'' + list[i].title + '\',\'' + list[i].pict_url + '\')" class="aui-flex">';
      htmlStr += '<div class="aui-pd-img">';
      htmlStr += '<img src="' + list[i].pict_url + '" alt="">';
      htmlStr += '</div>';
      htmlStr += '<div class="aui-flex-box">';
      htmlStr += '<h2><i class="icon icon-mall"></i>' + list[i].title + '</h2>';
      var s1 = list[i].coupon_info;
      var s2 = s1.substring(s1.indexOf("减") + 1).replace("元","");
      htmlStr += '<div style="width: 100%;"><span>' + s2 + '元隐藏劵</span></div></br>';
      htmlStr += "<div style='display: inline'>"
      htmlStr += '<div class="aui-flex-box" style="display:inline;width:50%;float:left;">';
      var finalPrice = (Number(list[i].zk_final_price) - Number(s2)).toFixed(2);
      htmlStr += '<h1>劵后￥' + finalPrice + '</h1>';
      htmlStr += '<h2>天猫价￥' + list[i].zk_final_price + '</h2>';
      htmlStr += '</div>';
      htmlStr += '<div class="aui-chang-box" style="width:50%;float: right;inline;float: right;">';
      htmlStr += '<div style="width:100%;"><span style="float: right;display:inline;"><em>领券</em>￥' + s2 + '</span></div>';
      htmlStr += '<div style="width:100%;display:inline;float: right;"><p>月销2398</p></div>';
      htmlStr += '</div>';
      htmlStr += '</div>';
      htmlStr += '</div>';
      htmlStr += '</a>';
    }
    $("#goods_list").append(htmlStr);
  });
}



function getTpwd(url,title,logo) {
  $.getJSON("/cms/tpwd?url="+url+"&title="+title+"&logo="+logo, function (json) {
    alert(json.tbk_tpwd_create_response.data.model)

    new $.flavr({ content : '赋值到淘宝领取优惠券吧', dialog : 'prompt',

      prompt : { placeholder: 'Enter something' }, onConfirm : function( $container, $prompt

      ){ alert( $prompt.val() ); return false; } });

//      $.zclip({
//        path: "/js/ZeroClipboard.swf",
//        copy: function () {//复制内容
//          return "lly";
//        },
//        afterCopy: function () {//复制成功
//          alert("已复制到剪贴板");
//        }
//      });
  });
}


// $(document).ready(function() {
//  $(window).scroll(function() {
//    if ($(document).scrollTop() >= $(document).height() - $(window).height()) {
//      var page = Number($("#hidPage").val()) + 1;
//      var category = $("#hidCategory").val();
//      GetGoodsList(page,category);
//      $("#hidPage").val(page + "");
//    }
//  });
// });
