package openapi

// OpenAPI 3.0 定義
openapi: "3.0.3"

info: {
	title:       "Todo API"
	version:     "0.1.0"
	description: "A simple Todo API built with Go and PostgreSQL"
}

servers: [{
	url:         "http://localhost:4000"
	description: "Local development server"
}]

// スキーマ定義
#Todo: {
	type: "object"
	properties: {
		id: {
			type:   "integer"
			format: "int64"
		}
		title: type:       "string"
		description: type: "string"
		completed: type:   "boolean"
		created_at: {
			type:   "string"
			format: "date-time"
		}
		updated_at: {
			type:   "string"
			format: "date-time"
		}
		user_id: {
			type:   "integer"
			format: "int64"
		}
	}
	required: ["id", "title", "completed", "user_id", "created_at", "updated_at"]
}

#CreateTodoRequest: {
	type: "object"
	properties: {
		title: type:       "string"
		description: type: "string"
	}
	required: ["title"]
}

#UpdateTodoRequest: {
	type: "object"
	properties: {
		title: type:       "string"
		description: type: "string"
		completed: type:   "boolean"
	}
}

// Todoのバッチ処理関連
#BatchTodoRequest: {
	type: "object"
	properties: {
		ids: {
			type: "array"
			items: {
				type:   "integer"
				format: "int64"
			}
			minItems: 1
			maxItems: 100
		}
	}
	required: ["ids"]
}

#BatchCompleteResponse: {
	type: "object"
	properties: {
		succeeded: {
			type: "array"
			items: "$ref": "#/components/schemas/Todo"
		}
		failed: {
			type: "array"
			items: "$ref": "#/components/schemas/BatchFailedItem"
		}
	}
	required: ["succeeded", "failed"]
}

#BatchDeleteResponse: {
	type: "object"
	properties: {
		succeeded: {
			type: "array"
			items: {
				type:   "integer"
				format: "int64"
			}
		}
		failed: {
			type: "array"
			items: "$ref": "#/components/schemas/BatchFailedItem"
		}
	}
	required: ["succeeded", "failed"]
}

#BatchFailedItem: {
	type: "object"
	properties: {
		id: {
			type:   "integer"
			format: "int64"
		}
		error: type: "string"
	}
	required: ["id", "error"]
}

#ErrorResponse: {
	type: "object"
	properties: message: type: "string"
	required: ["message"]
}

#HealthResponse: {
	type: "object"
	properties: status: type: "string"
	required: ["status"]
}

#InfoResponse: {
	type: "object"
	properties: {
		name: type:    "string"
		version: type: "string"
	}
	required: ["name", "version"]
}

// パス定義
paths: {
	"/": get: {
		summary:     "API information"
		description: "Get API information including name and version"
		operationId: "getInfo"
		tags: ["general"]
		responses: "200": {
			description: "OK"
			content: "application/json": schema: "$ref": "#/components/schemas/InfoResponse"
		}
	}
	"/health": get: {
		summary:     "Health check"
		description: "Check if the API is running"
		operationId: "getHealth"
		tags: ["general"]
		responses: "200": {
			description: "OK"
			content: "application/json": schema: "$ref": "#/components/schemas/HealthResponse"
		}
	}
	"/todos": {
		get: {
			summary:     "List all todos"
			description: "Get all todos for the authenticated user"
			operationId: "listTodos"
			tags: ["todos"]
			security: [{cookieAuth: []}]
			responses: {
				"200": {
					description: "OK"
					content: "application/json": schema: {
						type: "array"
						items: "$ref": "#/components/schemas/Todo"
					}
				}
				"401": {
					description: "Unauthorized"
					content: "application/json": schema: "$ref": "#/components/schemas/ErrorResponse"
				}
				"500": {
					description: "Internal server error"
					content: "application/json": schema: "$ref": "#/components/schemas/ErrorResponse"
				}
			}
		}
		post: {
			summary:     "Create a new todo"
			description: "Create a new todo with the provided information"
			operationId: "createTodo"
			tags: ["todos"]
			security: [{cookieAuth: []}]
			requestBody: {
				required: true
				content: "application/json": schema: "$ref": "#/components/schemas/CreateTodoRequest"
			}
			responses: {
				"201": {
					description: "Created"
					content: "application/json": schema: "$ref": "#/components/schemas/Todo"
				}
				"400": {
					description: "Bad request"
					content: "application/json": schema: "$ref": "#/components/schemas/ErrorResponse"
				}
				"401": {
					description: "Unauthorized"
					content: "application/json": schema: "$ref": "#/components/schemas/ErrorResponse"
				}
				"500": {
					description: "Internal server error"
					content: "application/json": schema: "$ref": "#/components/schemas/ErrorResponse"
				}
			}
		}
	}
	"/todos/{id}": {
		get: {
			summary:     "Get a todo by ID"
			description: "Get a single todo by its ID"
			operationId: "getTodo"
			tags: ["todos"]
			security: [{cookieAuth: []}]
			parameters: [{
				name:        "id"
				in:          "path"
				required:    true
				description: "Todo ID"
				schema: type: "integer", format: "int64"
			}]
			responses: {
				"200": {
					description: "OK"
					content: "application/json": schema: "$ref": "#/components/schemas/Todo"
				}
				"400": {
					description: "Invalid ID"
					content: "application/json": schema: "$ref": "#/components/schemas/ErrorResponse"
				}
				"401": {
					description: "Unauthorized"
					content: "application/json": schema: "$ref": "#/components/schemas/ErrorResponse"
				}
				"404": {
					description: "Todo not found"
					content: "application/json": schema: "$ref": "#/components/schemas/ErrorResponse"
				}
				"500": {
					description: "Internal server error"
					content: "application/json": schema: "$ref": "#/components/schemas/ErrorResponse"
				}
			}
		}
		put: {
			summary:     "Update a todo"
			description: "Update an existing todo by ID"
			operationId: "updateTodo"
			tags: ["todos"]
			security: [{cookieAuth: []}]
			parameters: [{
				name:        "id"
				in:          "path"
				required:    true
				description: "Todo ID"
				schema: type: "integer", format: "int64"
			}]
			requestBody: {
				required: true
				content: "application/json": schema: "$ref": "#/components/schemas/UpdateTodoRequest"
			}
			responses: {
				"200": {
					description: "OK"
					content: "application/json": schema: "$ref": "#/components/schemas/Todo"
				}
				"400": {
					description: "Invalid ID or request body"
					content: "application/json": schema: "$ref": "#/components/schemas/ErrorResponse"
				}
				"401": {
					description: "Unauthorized"
					content: "application/json": schema: "$ref": "#/components/schemas/ErrorResponse"
				}
				"404": {
					description: "Todo not found"
					content: "application/json": schema: "$ref": "#/components/schemas/ErrorResponse"
				}
				"500": {
					description: "Internal server error"
					content: "application/json": schema: "$ref": "#/components/schemas/ErrorResponse"
				}
			}
		}
		delete: {
			summary:     "Delete a todo"
			description: "Delete a todo by ID"
			operationId: "deleteTodo"
			tags: ["todos"]
			security: [{cookieAuth: []}]
			parameters: [{
				name:        "id"
				in:          "path"
				required:    true
				description: "Todo ID"
				schema: type: "integer", format: "int64"
			}]
			responses: {
				"204": description: "No Content"
				"400": {
					description: "Invalid ID"
					content: "application/json": schema: "$ref": "#/components/schemas/ErrorResponse"
				}
				"401": {
					description: "Unauthorized"
					content: "application/json": schema: "$ref": "#/components/schemas/ErrorResponse"
				}
				"404": {
					description: "Todo not found"
					content: "application/json": schema: "$ref": "#/components/schemas/ErrorResponse"
				}
				"500": {
					description: "Internal server error"
					content: "application/json": schema: "$ref": "#/components/schemas/ErrorResponse"
				}
			}
		}
	}
	"/todos/batch/complete": post: {
		summary:     "Batch complete todos"
		description: "Mark multiple todos as completed"
		operationId: "batchCompleteTodos"
		tags: ["todos"]
		security: [{cookieAuth: []}]
		requestBody: {
			required: true
			content: "application/json": schema: "$ref": "#/components/schemas/BatchTodoRequest"
		}
		responses: {
			"200": {
				description: "Batch operation completed"
				content: "application/json": schema: "$ref": "#/components/schemas/BatchCompleteResponse"
			}
			"400": {
				description: "Bad request"
				content: "application/json": schema: "$ref": "#/components/schemas/ErrorResponse"
			}
			"401": {
				description: "Unauthorized"
				content: "application/json": schema: "$ref": "#/components/schemas/ErrorResponse"
			}
			"500": {
				description: "Internal server error"
				content: "application/json": schema: "$ref": "#/components/schemas/ErrorResponse"
			}
		}
	}
	"/todos/batch/delete": post: {
		summary:     "Batch delete todos"
		description: "Soft delete multiple todos"
		operationId: "batchDeleteTodos"
		tags: ["todos"]
		security: [{cookieAuth: []}]
		requestBody: {
			required: true
			content: "application/json": schema: "$ref": "#/components/schemas/BatchTodoRequest"
		}
		responses: {
			"200": {
				description: "Batch operation completed"
				content: "application/json": schema: "$ref": "#/components/schemas/BatchDeleteResponse"
			}
			"400": {
				description: "Bad request"
				content: "application/json": schema: "$ref": "#/components/schemas/ErrorResponse"
			}
			"401": {
				description: "Unauthorized"
				content: "application/json": schema: "$ref": "#/components/schemas/ErrorResponse"
			}
			"500": {
				description: "Internal server error"
				content: "application/json": schema: "$ref": "#/components/schemas/ErrorResponse"
			}
		}
	}
}

components: {
	schemas: {
		Todo:                  #Todo
		CreateTodoRequest:     #CreateTodoRequest
		UpdateTodoRequest:     #UpdateTodoRequest
		BatchTodoRequest:      #BatchTodoRequest
		BatchCompleteResponse: #BatchCompleteResponse
		BatchDeleteResponse:   #BatchDeleteResponse
		BatchFailedItem:       #BatchFailedItem
		ErrorResponse:         #ErrorResponse
		HealthResponse:        #HealthResponse
		InfoResponse:          #InfoResponse
	}
	securitySchemes: cookieAuth: {
		type: "apiKey"
		in:   "cookie"
		name: "session"
	}
}

tags: [
	{name: "general", description: "General endpoints"},
	{name: "todos", description: "Todo management endpoints"},
]
