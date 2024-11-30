# GoatRobotics
This project is a dynamic chat application built using Go, leveraging goroutines, channels, and RESTful APIs to allow multiple clients to join, send messages, and leave a chat room concurrently. The chat room ensures thread-safe operations and efficient message broadcasting.

# GoatRobotics API

Welcome to the GoatRobotics API! This API allows you to interact with the chat system, where you can join, leave, send messages, and retrieve all chat logs.


## Installation
Clone the Project from Github

```bash
  git clone https://github.com/ROHITHSAKTHIVEL/GoatRobotics.git
  cd GoatRobotics
```

#### Pre requisite
```bash
  # If you want to Change Something in Code then you need to have :
   Go lang Installed on your Computer
```
## Run
Using go run main.go

```bash
 go run main.go
```

Using Docker-compose

```bash
docker compose up -d
```


#### To Change Server Configuration
Change this file to Change the Server Configuration
```bash
./config.json 
```
Feel free to Change the Server Configuration as Per your needs


    
## API Reference

#### 1. Join the Chat Room

```http
  GET http://localhost:8080/rpc/GOATROBOTICS/join?id=12345
```
#### Curl 

```http
 curl --location 'http://localhost:9000/join?id=12345'
```

| Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `id` | `string` | **Required**. To Join the ID is Required |



#### Response 
```json 
{ 
    "message": "User Chat Successfully", 
}
```
| Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `message` | `string` |  Status of the Response  |


---
---
#### 2. Leave the Chat Room

```http
  GET http://localhost:9000/leave?id=12345
```
#### Curl 

```http
 curl --location 'http://localhost:9000/leave?id=12345'
```

| Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `id` | `string` | **Required**. To Left the ID is Required |



#### Response 
```json 
{
    "message": "user Left Chat Successfully",
}
```
| Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `message` | `string` |  Status of the Response  |

---
---
#### 3. To Send Message to the Chat Room

```http
  GET http://localhost:9000/send?id=12345&message=Hello!
```
#### Curl 

```http
curl --location 'http://localhost:9000/send?id=12345&message=Hello!'
```

| Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `id` | `string` | **Required**. To Message the ID is Required |
| `messgae` | `string` | **Required**. Message to Publish is Required |



#### Response 
```json 
{
    "message": "Message Sent Successfully",
}
```
| Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `message` | `string` |  Status of the Response  |

---
---
#### 4. To Get Messages form Chat Room

```http
  GET http://localhost:9000/messages?id=12345&message=Hello!
```
#### Curl 

```http
curl --location 'http://localhost:9000/messages?id=12345&message=Hello!'
```

| Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `id` | `string` | **Required**. To Message the ID is Required |
| `messgae` | `string` | **Required**. Message to Publish is Required |



#### Response 
```json 
{
    "messages": [
        {
            "clientID": "client-1",
            "message": "hello!",
            "sentTime": "2024-11-30T22:36:38.9496343+05:30"
        }
    ],
    "ReponseTime": "2024-11-30T22:36:49.0027989+05:30",
    "userId": "client-1"
}

```
---
---
#### 5. Ping to get Server Status

```http
  GET http://localhost:9000/ping
```
#### Curl 

```http
curl --location 'http://localhost:9000/ping'
```




#### Response 
```json 
{
    "message": "Pinged Successfully",
}
```
| Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `message` | `string` |  Status of the Response  |1

---
---
#### 6. Get Audit or Logs

```http
  http://localhost:8080/viewlogs
```
#### Curl 

```http
ccurl --location 'http://localhost:9000/viewlogs' 
```




#### Response 
```json 
[
    {
        "clientID": "4",
        "startTime": "2024-11-30T22:36:26.5243771+05:30",
        "endTime": "2024-11-30T22:36:26.526367+05:30",
        "methodName": "GET",
        "requestURL": "/send?id=4&message=hello!",
        "requestHeaders": "User-Agent: PostmanRuntime/7.43.0\nAccept: */*\nPostman-Token: e4ac0a0b-02c6-44e2-8b78-458229044ca6\nAccept-Encoding: gzip, deflate, br\nConnection: keep-alive\n",
        "queryParameters": "id=4&message=hello%21",
        "responseDuration": 1989900,
        "responseBody": "client 4 does not exist\n",
        "responseHeaders": "Vary: Origin\nContent-Type: text/plain; charset=utf-8\nX-Content-Type-Options: nosniff\n",
        "statusCode": 404
    },
    {
        "clientID": "client-1",
        "startTime": "2024-11-30T22:36:38.9496343+05:30",
        "endTime": "2024-11-30T22:36:38.9503725+05:30",
        "methodName": "GET",
        "requestURL": "/send?id=client-1&message=hello!",
        "requestHeaders": "User-Agent: PostmanRuntime/7.43.0\nAccept: */*\nPostman-Token: d10e380b-6ac6-4e5e-b7f5-737657f763d9\nAccept-Encoding: gzip, deflate, br\nConnection: keep-alive\n",
        "queryParameters": "id=client-1&message=hello%21",
        "responseDuration": 738200,
        "responseBody": "{\"message\":\"Message sent successfully\"}\n",
        "responseHeaders": "Vary: Origin\n",
        "statusCode": 200
    },
    {
        "clientID": "client-1",
        "startTime": "2024-11-30T22:36:49.0027989+05:30",
        "endTime": "2024-11-30T22:36:49.0037856+05:30",
        "methodName": "GET",
        "requestURL": "/messages?id=client-1",
        "requestHeaders": "User-Agent: PostmanRuntime/7.43.0\nAccept: */*\nPostman-Token: 3cd31ac3-a9b9-45f4-9f2f-5f000f0583b9\nAccept-Encoding: gzip, deflate, br\nConnection: keep-alive\n",
        "queryParameters": "id=client-1",
        "responseDuration": 986700,
        "responseBody": "{\"messages\":[{\"clientID\":\"client-1\",\"message\":\"hello!\",\"sentTime\":\"2024-11-30T22:36:38.9496343+05:30\"}],\"ReponseTime\":\"2024-11-30T22:36:49.0027989+05:30\",\"userId\":\"client-1\"}\n",
        "responseHeaders": "Vary: Origin\nContent-Type: application/json\n",
        "statusCode": 200
    },
    {
        "clientID": "client-1",
        "startTime": "2024-11-30T22:36:56.9825099+05:30",
        "endTime": "2024-11-30T22:36:56.9825099+05:30",
        "methodName": "GET",
        "requestURL": "/leave?id=client-1",
        "requestHeaders": "Connection: keep-alive\nUser-Agent: PostmanRuntime/7.43.0\nAccept: */*\nPostman-Token: 62da19a5-664f-4277-b3ab-8184701fe8a9\nAccept-Encoding: gzip, deflate, br\n",
        "queryParameters": "id=client-1",
        "responseBody": "{\"message\":\"User left successfully\"}\n",
        "responseHeaders": "Vary: Origin\n",
        "statusCode": 200
    },
    {
        "clientID": "client-1",
        "startTime": "2024-11-30T22:36:20.0931171+05:30",
        "endTime": "2024-11-30T22:36:20.0931171+05:30",
        "methodName": "GET",
        "requestURL": "/join?id=client-1",
        "requestHeaders": "User-Agent: PostmanRuntime/7.43.0\nAccept: */*\nPostman-Token: 675a57b0-8b21-4b7b-8caa-a0eeabafe160\nAccept-Encoding: gzip, deflate, br\nConnection: keep-alive\n",
        "queryParameters": "id=client-1",
        "responseBody": "{\"message\":\"User joined successfully\"}\n",
        "responseHeaders": "Vary: Origin\n",
        "statusCode": 200
    }
]
```


---

### Lessons Learned

A summary of key concepts and best practices implemented throughout my recent projects. These lessons have contributed to the development of robust, scalable, and efficient systems.

#### 1. Thread-Safe Operations & Efficient Message Broadcasting
Ensuring thread-safe operations while maintaining high performance in concurrent environments. Efficient message broadcasting ensures that messages are transmitted seamlessly to multiple clients.
> Optimized for high concurrency 🏎️💨

#### 2. Middleware for Request/Response Interception
Implemented middleware to intercept and log all incoming and outgoing requests and responses. Logs are saved in `./logs/Audit.audit` for traceability.
> Enhanced security and traceability 🔄

#### 3. Docker for Dependency Management
Used Docker to create a consistent environment, eliminating dependency issues across different environments and ensuring smooth application execution.
> Dockerized for consistency 🐋

#### 4. Custom Error Types for Scalability
Implemented custom error types to support scalable error handling, ensuring clean, manageable, and maintainable code.
> Structured for growth ⚙️

#### 5. Constants Instead of Hardcoding
Replaced hardcoded values with constants to improve maintainability, reduce errors, and make the code more adaptable to changes.
> Improved code quality 🔧

#### 6. Everything Configurable with `config.json`
All key configurations are easily manageable through a `config.json` file, allowing for flexible and centralized configuration management across environments.  
> Simplifies deployment and environment setup 🛠️

---

These lessons have shaped my development process, focusing on best practices that ensure scalable, reliable, and efficient software solutions. I continue to refine these skills, aiming to deliver high-quality code in every project.



## Support

For support, email rohithsakthivel2002@gmail.com 

