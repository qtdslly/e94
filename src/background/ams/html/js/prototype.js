String.prototype.isEmpty=function(){return this==undefined||this.trim()==""||this.trim()=="undefined"?true:false};String.prototype.isNotEmpty=function(){return !this.isEmpty()};String.prototype.trim=function(){return $.trim(this)};String.prototype.equals=function(a){return this==a};String.prototype.isEmail=function(){return/^([a-zA-Z0-9]+[_|\_|\.]?)*[a-zA-z0-9]+@([a-zA-Z0-9]+[_|\_|\.]?)*[a-zA-z0-9]+\.[a-zA-Z]{2,3}$/.test(this)};String.prototype.isMobile=function(){return/^1[34578]{1}\d{9}$/.test(this)};String.prototype.isPhone=function(){return/^(([0\+]\d{2,3}-)?(0\d{2,3})-)(\d{7,8})(-(\d{3,}))?$/.test(this)};String.prototype.isNumber=function(){return/^\d+$/.test(this)};String.prototype.isDecimal=function(){return/^\d+(\.\d+)?$/.test(this)};String.prototype.isLetter=function(){return/^[a-zA-Z]+$/.test(this)};String.prototype.toNumber=function(){return parseInt(this)};String.prototype.toFloat=function(){return parseFloat(this)};String.prototype.decimalFormat=function(b,c){var d=$(b).val().trim();if(!d.isDecimal()){return d}var a=c+1;if(d.indexOf(".")>0&&d.substring(d.indexOf(".")).length>a){d=d.substring(0,d.indexOf(".")+a);$(b).val(d)}return d};String.prototype.matchers=function(){var a="string";if(this.trim().isEmail()){a="email"}if(this.trim().isMobile()){a="mobile"}if(this.trim().isPhone()){a="phone"}if(this.trim().isNumber()){a="number"}if(this.trim().isLetter()){a="letter"}if(this.trim().isDecimal()){a="decimal"}return a};Number.prototype.compare=function(b,a){return this>=b&&this<=a?true:false};var BwNumber={};BwNumber.add=function(a,g){var e=a.indexOf(".");var d=e>0?a.length-e:0;var c=g.indexOf(".");var b=c>0?g.length-c:0;var f=(a*Math.pow(10,d+b)).toFixed(0)+(g*Math.pow(10,d+b)).toFixed(0);return f/(Math.pow(10,d+b))};BwNumber.minus=function(a,g){var e=a.indexOf(".");var d=e>0?a.length-e:0;var c=g.indexOf(".");var b=c>0?g.length-c:0;var f=(a*Math.pow(10,d+b)).toFixed(0)-(g*Math.pow(10,d+b)).toFixed(0);return f/(Math.pow(10,d+b))};BwNumber.multiply=function(a,g){var e=a.indexOf(".");var d=e>0?a.length-e:0;var c=g.indexOf(".");var b=c>0?g.length-c:0;var f=(a*Math.pow(10,d)).toFixed(0)*(g*Math.pow(10,b)).toFixed(0);return f/(Math.pow(10,d+b))};BwNumber.divider=function(a,g){var e=a.indexOf(".");var d=e>0?a.length-e:0;var c=g.indexOf(".");var b=c>0?g.length-c:0;var f=(a*Math.pow(10,d+b)).toFixed(0)/(g*Math.pow(10,b)).toFixed(0);return f/(Math.pow(10,d))};