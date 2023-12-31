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
        "/health": {
            "get": {
                "description": "to check http server health",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "health"
                ],
                "summary": "show http server health",
                "responses": {}
            }
        },
        "/user/words": {
            "get": {
                "description": "GetUserWords",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "word"
                ],
                "summary": "GetUserWords",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.UserWords"
                        }
                    }
                }
            }
        },
        "/words/add": {
            "post": {
                "description": "Add word",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "word"
                ],
                "summary": "Add word",
                "parameters": [
                    {
                        "description": "AddWordRequest body",
                        "name": "AddWordRequest",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.AddWordRequest"
                        }
                    }
                ],
                "responses": {}
            }
        }
    },
    "definitions": {
        "models.AddWordRequest": {
            "type": "object",
            "properties": {
                "definitions": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.Definition"
                    }
                },
                "word": {
                    "$ref": "#/definitions/models.Word"
                }
            }
        },
        "models.Definition": {
            "type": "object",
            "properties": {
                "definition": {
                    "type": "string"
                },
                "lang": {
                    "type": "string"
                }
            }
        },
        "models.UserWord": {
            "type": "object",
            "properties": {
                "definitions": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.Definition"
                    }
                },
                "word": {
                    "$ref": "#/definitions/models.Word"
                }
            }
        },
        "models.UserWords": {
            "type": "object",
            "properties": {
                "words": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.UserWord"
                    }
                }
            }
        },
        "models.Word": {
            "type": "object",
            "properties": {
                "example": {
                    "type": "string"
                },
                "image_url": {
                    "type": "string"
                },
                "lang": {
                    "type": "string"
                },
                "link": {
                    "type": "string"
                },
                "word": {
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
