import logging

from ClientChat import ClientChat, print_message
from menus import chat_menu


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

        print(f"chatID: {chat['id']} - User {otherUser['name']}|ID:{otherUser['id']}")


def deleteChat(client: ClientChat):
    chatID = input("ID чата: ")

    try:
        res = client.deleteChat(chatID=chatID)
        print(res)

    except Exception as e:
        print("Ошибка удаления чата " + str(e))


def deleteMessage(client: ClientChat):
    mesID = input("ID сообщения: ")

    try:
        res = client.deleteChat(mesID)
        print(res)
    except Exception as e:
        print("Ошибка удаления сообщения " + str(e))


def Messaging(client: ClientChat, chatID: str):
    client.setCurrentChat(chatID)

    messages = client.loadMessages()
    for mes in messages:
        print_message(payload=mes['payload'],
                      senderName=mes['sender']['name'],
                      mesId=mes['id'],
                      time=mes['time']
                      )

    while True:
        inpt = input()
        if inpt == "!q":
            break

        try:
            message = client.postMessage(inpt)

            print_message(payload=message['payload'],
                          mesId=message['id'],
                          senderName=client.getName(),
                          time=message['time'])

        except Exception as e:
            logging.log(level=logging.INFO, msg=str(e))

    client.exitChat()


def openChatWithID(client: ClientChat, chatID: str):
    try:
        client.setCurrentChat(chatID)
    except Exception as exp:
        print("\n\nНет доступа к этому чату или его не существует: " + str(exp))
        client.exitChat()
        return

    client.exitChat()

    while True:
        inpt = chat_menu()

        match inpt:
            case "1":
                Messaging(client, chatID)
            case "2":
                deleteMessage(client)
            case "0":
                client.exitChat()
                return


def createChat(client: ClientChat):
    otherUserID = input("ID пользователя: ")

    res = client.createChat(otherUserID)

    print(f"Создан чат с пользователем {client.getName(otherUserID)} (ID:{otherUserID}),"
          f" chatID: {res['id']}")

    openChatWithID(client, res['id'])


def openChat(client: ClientChat):
    chatID = input("ID чата: ")

    openChatWithID(client, chatID)
