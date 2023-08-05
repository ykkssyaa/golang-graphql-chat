from ClientChat import ClientChat
import asyncio


def CreateUser(client: ClientChat):
    newName = input("Введите имя пользователя: ")
    result = client.createUser(newName)
    newId = result["id"]

    return client.auth(newId)


def SignIn(client):

    singInId = input("Введите ID пользователя: ")
    return client.auth(singInId)


