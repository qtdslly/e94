/**
 * Created by ccgo on 2019/4/30.
 */
function getParam(paramName) {
    // if(this.location.search.indexOf(paramName) < 1){
    //     return ""
    // }
    paramValue = "", isFound = !1;
    if (this.location.search.indexOf("?") == 0 && this.location.search.indexOf("=") > 1) {
        arrSource = unescape(this.location.search).substring(1, this.location.search.length).split("&"), i = 0;
        while (i < arrSource.length && !isFound) arrSource[i].indexOf("=") > 0 && arrSource[i].split("=")[0].toLowerCase() == paramName.toLowerCase() && (paramValue = arrSource[i].split("=")[1], isFound = !0), i++
    }
    return paramValue == "" && (paramValue = null), paramValue
}

function jump(obj,id){
    var token = getParam("Authorization");
    alert(token)

    if(obj == 0){
        location.href = "http://" + window.location.host + "/html/video/list.html?" + "Authorization=" + token;
    }else if(obj == 1){
        location.href = "http://" + window.location.host + "/html/user1.html?" + "Authorization=" + token;
    }else if(obj == 2){
        location.href = "http://" + window.location.host + "/html/userwallet.html?" + "Authorization=" + token + "&user_id="+id;
    }else if(obj == 3){
        location.href = "http://" + window.location.host + "/html/sonuser.html?" + "Authorization=" + token + "&parent_id="+id;
    }else if(obj == 4){
        location.href = "http://" + window.location.host + "/html/usermachine.html?" + "Authorization=" + token + "&user_id="+id;
    }else if(obj == 5){
        location.href = "http://" + window.location.host + "/html/contract.html?" + "Authorization=" + token;
    }else if(obj == 6){
        location.href = "http://" + window.location.host + "/html/paramset.html?" + "Authorization=" + token;
    }else if(obj == 7){
        location.href = "http://" + window.location.host + "/html/notice.html?" + "Authorization=" + token;
    }else if(obj == 8){
        location.href = "http://" + window.location.host + "/html/question.html?" + "Authorization=" + token;
    }else if(obj == 9){
        location.href = "http://" + window.location.host + "/html/msgboard.html?" + "Authorization=" + token;
    }else if(obj == 10){
        location.href = "http://" + window.location.host + "/html/worklist.html?" + "Authorization=" + token;
    }
}