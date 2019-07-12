/**
 * Created by ccgo on 2019/6/12.
 */
function getParam(paramName) {
  paramValue = "", isFound = !1;
  if (this.location.search.indexOf("?") == 0 && this.location.search.indexOf("=") > 1) {
    arrSource = unescape(this.location.search).substring(1, this.location.search.length).split("&"), i = 0;
    while (i < arrSource.length && !isFound) arrSource[i].indexOf("=") > 0 && arrSource[i].split("=")[0].toLowerCase() == paramName.toLowerCase() && (paramValue = arrSource[i].split("=")[1], isFound = !0), i++
  }
  return paramValue == "" && (paramValue = null), paramValue
}


angular.module('stock', ['ionic'])
  .controller( 'listCtrl',['$scope','$timeout' ,'$http',function($scope,$timeout,$http){
    $scope.detail = function(code) {
      window.location.href="http://localhost/html/detail.html?code="+code;
    }

    $http({
      method: 'GET',
      url: 'http://localhost/stock/notice/list',
    }).then(function successCallback(response) {
      $scope.notices = response.data.stocks;
    }, function errorCallback(response) {
      // 请求失败执行代码
    }).finally(function() {
      $scope.$broadcast('scroll.refreshComplete');
    });
  }])
  .controller( 'detailCtrl',['$scope','$timeout' ,'$ionicPopup','$http',function($scope,$timeout,$http){
    var code = getParam("code");

    $http({
      method: 'GET',
      url: 'http://localhost/stock/notice/detail?code='+code,
    }).then(function successCallback(response) {
      $scope.stock = response.data.notice;
      if($scope.stock.state == 1){
        document.getElementById("state").checked = true;
      }else{
        document.getElementById("state").checked = false;
      }
      if($scope.stock.frequency == 1){
        document.getElementById("frequency").checked = true;
      }else{
        document.getElementById("frequency").checked = false;
      }
    }, function errorCallback(response) {
      // 请求失败执行代码
    }).finally(function() {
      $scope.$broadcast('scroll.refreshComplete');
    });

    $scope.cancel = function() {
      window.history.go(-1);
    }

    $scope.update = function() {
      var stat = "0";
      if(document.getElementById("state").checked == true){
        stat = "1";
      }

      var fre = "0";
      if(document.getElementById("frequency").checked == true){
        fre = "1";
      }

      $http({
        method: 'POST',
        url: 'http://localhost/stock/notice/update',
        data:{code:$scope.stock.code,name:$scope.stock.name,buy_price:Number($scope.stock.buy_price),
          sell_price:Number($scope.stock.sell_price),buy_count:Number($scope.stock.buy_count),
          sell_count:Number($scope.stock.sell_count),state:stat,frequency:fre}
      }).then(function successCallback(response) {
        $scope.showAlert = function() {
          var alertPopup = $ionicPopup.alert({
            title: '提示',
            template: '修改成功'
          });
          alertPopup.then(function(res) {
            console.log('修改成功');
          });
        };
      }, function errorCallback(response) {
        // 请求失败执行代码
      }).finally(function() {
        $scope.$broadcast('scroll.refreshComplete');
      });
    }


  }])
