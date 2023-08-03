from ClientChat import ClientChat
import asyncio


async def CreateUser(client: ClientChat):
    newName = input("Введите имя пользователя: ")
    result = await client.createUser(newName)
    newid = result["id"]

    asyncio.create_task(client.auth(newid))


def SignIn(client):

    singInId = input("Введите ID пользователя: ")
    asyncio.run(client.auth(singInId))


