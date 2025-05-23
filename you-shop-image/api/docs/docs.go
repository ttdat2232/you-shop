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
        "/images/banner": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "images"
                ],
                "summary": "Get Banners Data",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/model.ApiResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/definitions/image.ImageResponse"
                                            }
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/model.ApiResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/err.AppError"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            },
            "delete": {
                "tags": [
                    "images"
                ],
                "summary": "Delete Banners",
                "responses": {
                    "202": {
                        "description": "Accepted"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/model.ApiResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/err.AppError"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/images/upload": {
            "post": {
                "consumes": [
                    "multipart/form-data"
                ],
                "tags": [
                    "images"
                ],
                "summary": "Upload Image",
                "parameters": [
                    {
                        "type": "file",
                        "description": "Image File",
                        "name": "image_file",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Image Alt",
                        "name": "alt",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Owner ID",
                        "name": "owner_id",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/model.ApiResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/err.AppError"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "default": {
                        "description": ""
                    }
                }
            }
        },
        "/images/upload/banner": {
            "post": {
                "consumes": [
                    "multipart/form-data"
                ],
                "tags": [
                    "images"
                ],
                "summary": "Upload Banners",
                "parameters": [
                    {
                        "type": "file",
                        "description": "Banner Files",
                        "name": "banner_files",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/model.ApiResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/err.AppError"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/images/{id}": {
            "get": {
                "produces": [
                    "image/jpeg",
                    "image/png"
                ],
                "tags": [
                    "images"
                ],
                "summary": "Serve Image",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Image ID in UUID format (e.g: 123e4567-e89b-12d3-a456-426614174000)",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {}
            },
            "delete": {
                "tags": [
                    "images"
                ],
                "summary": "Delete Image",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Image ID in UUID format (e.g: 123e4567-e89b-12d3-a456-426614174000)",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "202": {
                        "description": "Accepted"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/model.ApiResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/err.AppError"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "err.AppError": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer"
                },
                "data": {},
                "message": {
                    "type": "string"
                },
                "title": {
                    "type": "string"
                }
            }
        },
        "image.ImageResponse": {
            "type": "object",
            "properties": {
                "alt": {
                    "type": "string"
                },
                "content_type": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "image_url": {
                    "type": "string"
                }
            }
        },
        "model.ApiResponse": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer"
                },
                "data": {},
                "is_success": {
                    "type": "boolean"
                },
                "message": {
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
