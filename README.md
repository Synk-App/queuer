# Queuer
Service to manage queue from system.

# Getting Started

First step is to create a `.env` file in project root and change example values to your config. You can use `example.env` file from `_setup` folder as template.

And then, run `docker compose up -d` into project root to start project.

## Tests

The easy way to run tests is just run `docker compose up -d` command to start project with variables. So, enter in `synk_queuer` with `docker exec` and run `go test ./tests -v`.

## Certificates

This app must run in HTTPS to authentication works properly. So, to install it, just setup `[mkcert](https://github.com/FiloSottile/mkcert)` into your machine and then run command below into root directory of this project.

```
mkcert -key-file ./.cert/key.pem -cert-file ./.cert/cert.pem localhost synk_queuer
```

## Network

You can use a custom network for this services, using then `synk_network` you must create before run it. So, to create on just run command below once during initial setup.

```
docker network create synk_network
```

# Routes

## Get info about app

> `GET` /about

### Response

```json
{
	"ok": true,
	"error": "",
	"info": {
		"server_port": "8080",
		"app_port": "8083",
		"db_working": true
	},
	"list": null
}
```

## Publish content from posts

> `POST` /send

### Request

```json
{
	"posts": [1]
}
```

### Response

```json
{
	"resource": {
		"ok": true,
		"error": ""
	},
	"posts": {
		"1": {
			"2": {
				"resource": {
					"ok": true,
					"error": ""
				},
				"http_code": 200,
				"raw": "{\"ok\":true,\"result\":{\"message_id\":9,\"from\":{\"id\":12345678910,\"is_bot\":true,\"first_name\":\"Bot\",\"username\":\"bot\"},\"chat\":{\"id\":12345678910,\"first_name\":\"Cristian\",\"last_name\":\"Prochnow\",\"type\":\"private\"},\"date\":1764218621,\"text\":\"Nova publica\\u00e7\\u00e3o show\"}}"
			},
			"3": {
				"resource": {
					"ok": true,
					"error": ""
				},
				"http_code": 200,
				"raw": "{\"type\":0,\"content\":\"Nova publica\\u00e7\\u00e3o show\",\"mentions\":[],\"mention_roles\":[],\"attachments\":[],\"embeds\":[],\"timestamp\":\"2025-11-27T04:43:41.712000+00:00\",\"edited_timestamp\":null,\"flags\":0,\"components\":[],\"id\":\"12345678910\",\"channel_id\":\"144300515123456789104281394400\",\"author\":{\"id\":\"12345678910\",\"username\":\"Captain Hook\",\"avatar\":null,\"discriminator\":\"0000\",\"public_flags\":0,\"flags\":0,\"bot\":true,\"global_name\":null,\"clan\":null,\"primary_guild\":null},\"pinned\":false,\"mention_everyone\":false,\"tts\":false,\"webhook_id\":\"12345678910\"}\n"
			}
		}
	}
}
```