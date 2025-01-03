// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/auth/categories": {
            "post": {
                "description": "创建分类",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "category"
                ],
                "summary": "创建分类",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Authorization",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "分类信息",
                        "name": "category",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handles.CategoryReq"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "创建成功",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/auth/upload": {
            "post": {
                "description": "上传图片",
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "上传图片",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Token",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "file",
                        "description": "图片",
                        "name": "image",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "自定义短链",
                        "name": "short_link",
                        "in": "formData"
                    },
                    {
                        "type": "string",
                        "description": "分类",
                        "name": "category",
                        "in": "formData"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "图片上传成功",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/private/setting": {
            "get": {
                "description": "获取设置",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "setting"
                ],
                "summary": "获取设置",
                "parameters": [
                    {
                        "type": "string",
                        "description": "key",
                        "name": "key",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "keys",
                        "name": "keys",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "设置",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "post": {
                "description": "保存设置",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "setting"
                ],
                "summary": "保存设置",
                "parameters": [
                    {
                        "description": "设置",
                        "name": "settings",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/model.SettingItem"
                            }
                        }
                    }
                ],
                "responses": {}
            },
            "delete": {
                "description": "删除设置",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "setting"
                ],
                "summary": "删除设置",
                "parameters": [
                    {
                        "type": "string",
                        "description": "key",
                        "name": "key",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {}
            }
        },
        "/api/private/setting/token": {
            "post": {
                "description": "重置token",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "setting"
                ],
                "summary": "重置token",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/private/settings": {
            "get": {
                "description": "列出设置",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "setting"
                ],
                "summary": "列出设置",
                "parameters": [
                    {
                        "type": "string",
                        "description": "group",
                        "name": "group",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "groups",
                        "name": "groups",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "设置列表",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/model.SettingItem"
                            }
                        }
                    }
                }
            }
        },
        "/api/public/categories": {
            "get": {
                "description": "获取图片分类",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "category"
                ],
                "summary": "获取图片分类",
                "responses": {
                    "200": {
                        "description": "图片分类",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/public/categories/{name}": {
            "get": {
                "description": "根据分类名获取图片",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "category"
                ],
                "summary": "根据分类名获取图片",
                "parameters": [
                    {
                        "type": "string",
                        "description": "分类名",
                        "name": "name",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "图片内容",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/public/images": {
            "get": {
                "description": "分页列出图片",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "image"
                ],
                "summary": "分页列出图片",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "页码",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "每页数量",
                        "name": "page_size",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "图片列表",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/public/images/{name}": {
            "get": {
                "description": "根据文件名获取图片",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "image"
                ],
                "summary": "根据文件名获取图片",
                "parameters": [
                    {
                        "type": "string",
                        "description": "文件名",
                        "name": "name",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "图片内容",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/public/login": {
            "post": {
                "description": "登录",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "登录",
                "parameters": [
                    {
                        "description": "用户信息",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handles.LoginReq"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "登录成功",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/public/random": {
            "get": {
                "description": "随机获取一个图片 支持分类",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "image"
                ],
                "summary": "随机获取一个图片 支持分类",
                "parameters": [
                    {
                        "type": "string",
                        "description": "分类",
                        "name": "category",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "图片内容",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/public/settings": {
            "get": {
                "description": "获取公共设置",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "setting"
                ],
                "summary": "获取公共设置",
                "responses": {
                    "200": {
                        "description": "公共设置",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/api/public/short/{short_link}": {
            "get": {
                "description": "根据短链获取图片",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "image"
                ],
                "summary": "根据短链获取图片",
                "parameters": [
                    {
                        "type": "string",
                        "description": "短链",
                        "name": "short_link",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "图片内容",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "handles.CategoryReq": {
            "type": "object",
            "properties": {
                "is_public": {
                    "type": "boolean"
                },
                "is_random": {
                    "type": "boolean"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "handles.LoginReq": {
            "type": "object",
            "properties": {
                "password": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "model.SettingItem": {
            "type": "object",
            "required": [
                "key"
            ],
            "properties": {
                "flag": {
                    "description": "0 = public, 1 = private, 2 = readonly, 3 = deprecated, etc.",
                    "type": "integer"
                },
                "group": {
                    "description": "use to group setting in frontend",
                    "type": "integer"
                },
                "help": {
                    "description": "help message",
                    "type": "string"
                },
                "index": {
                    "type": "integer"
                },
                "key": {
                    "description": "unique key",
                    "type": "string"
                },
                "options": {
                    "description": "values for select",
                    "type": "string"
                },
                "type": {
                    "description": "string, number, bool, select",
                    "type": "string"
                },
                "value": {
                    "description": "value",
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
