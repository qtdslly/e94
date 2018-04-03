package constant

/*
	Warning: make sure the type defined here consistent with OperationLog type
*/
const (
	OperationTypeUnknown = 0 //未知类型
	OperationTypeCreate  = 1 //创建
	OperationTypeUpdate  = 2 //更新
	OperationTypeList    = 3 //列表查询
	OperationTypeDetail  = 4 //详情
	OperationTypeDelete  = 5 //删除
	OperationTypeLogin   = 6 //登录
	OperationTypeLogout  = 7 //登出

)
