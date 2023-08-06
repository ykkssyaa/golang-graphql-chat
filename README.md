# golang-graphql-chat

## Описание

Проект представляет собой чат, использующий GraphQL. 
Общение пользователей в реальном времени проходит по протоколу WebSocket

[Информация о клиенте](client/README.md)

### Реализованный функционал

- Добавление полтзователей
- Создание чата
- Удаление чата
- Вывод списка чатов (всех и для конкретного пользователя)
- Вывод всех пользователей
- Отправка сообщения
- Удаление сообщения
- Вывод всех сообщений от пользователя с пагинацией

## Запуск

1. `go mod download`
2. `docker compose up -d`
3. `go run main.go`

## Query запросы

### users

```
query Users {
    users {
        id
        name
    }
}
```

### chats

```
query Chats {
    chats {
        id
        user_1 {
            id
            name
        }
        user_2 {
            id
            name
        }
    }
}
```

### messagesFromUser

```
query MessagesFromUser {
    messagesFromUser(input: {receiver: "2", sender: "1"}, first: 5) {
        totalCount
        edges {
            cursor
            node {
                id
                payload
                chatID
                time
                sender {
                    id
                    name
                }
                receiver {
                    id
                    name
                }
            }
        }
        pageInfo {
            startCursor
            endCursor
            hasNextPage
        }
    }
}
```

## Mutation

### createUser

```
mutation CreateUser {
    createUser(input: {name: "World"}) {
        id
        name
    }
}
```

### postMessage

```
mutation PostMessage {
    postMessage(input: {payload: "Hello World", sender: "1", receiver: "2"}) {
        id
        payload
        chatID
        time
        sender {
            id
            name
        }
        receiver {
            id
            name
        }
    }
}
```

### CreateChat

```
mutation CreateChat {
    createChat(input: {user1: "2", user2: "1"}) {
        id
    }
}
```

### deleteChat

```
mutation DeleteChat {
    deleteChat(id: "1")
}
```

### deleteMessage

```
mutation DeleteMessage {
    deleteMessage(id: "1")
}
```

## Subscription

### userJoined

```
subscription{
  userJoined(user:1){
    id
    chatID
    
    payload
    
    receiver{
      id
      name
    }
    sender{
      id
      name
    }
  time
  
  }
}
```