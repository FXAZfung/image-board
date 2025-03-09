package common

// 原泛型定义
// type Resp[T any] struct {
//     Code    int    `json:"code"`
//     Message string `json:"message"`
//     Data    T      `json:"data"`
// }

// swagger:model BaseResponse
type Resp struct {
	// 状态码
	// Example: 200
	Code int `json:"code"`
	// 消息描述
	// Example: success
	Message string `json:"message"`
	// 返回数据
	Data interface{} `json:"data"`
}

// swagger:model PageResponse
type PageResp struct {
	// 分页内容
	Content interface{} `json:"content"`
	// 总记录数
	// Example: 100
	Total int64 `json:"total"`
}
