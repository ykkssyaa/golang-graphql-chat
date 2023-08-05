from ClientChat import ClientChat


def allUsers(client: ClientChat):
    users = client.allUsers()
    userID = client.currentID

    for user in users:
        print(f"ID:{user['id']}| {user['name'] + ('(Вы)' if user['id'] == userID else '')}")


def allChats(client: ClientChat):
    userID = client.currentID
    chats = client.chats_of_user()

    for chat in chats:
        otherUser = {}
        if chat["user_1"]["id"] != userID:
            otherUser = chat["user_1"]
        else:
            otherUser = chat["user_2"]

        print(f"ID: {chat['id']} - User {otherUser['name']} | ID:{otherUser['id']}")


def deleteChat(client: ClientChat):

    chatID = input("ID чата: ")

    res = client.deleteChat(chatID=chatID)

    print(res)


def openChatWithID(client: ClientChat, chatID: str):

    try:
        client.setCurrentChat(chatID)
    except Exception:
        print("Нет доступа к этому чату или его не существует")


def createChat(client: ClientChat):
    otherUserID = input("ID пользователя: ")

    res = client.createChat(otherUserID)

    print(f"Создан чат с пользователем {client.getName(otherUserID)} (ID:{otherUserID}),"
          f" chatID: {res['id']}")

    openChatWithID(client, res['id'])


def openChat(client: ClientChat):
    chatID = input("ID чата: ")

    openChatWithID(client, chatID)
