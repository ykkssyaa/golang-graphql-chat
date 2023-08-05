from ClientChat import ClientChat


def CreateUser(client: ClientChat):
    newName = input("Введите имя пользователя: ")

    if len(newName) < 3:
        print("Имя слишком короткое! (менее 3х символов)")
        return None

    result = client.createUser(newName)
    newId = result["id"]

    return client.auth(newId)


def SignIn(client):

    singInId = input("Введите ID пользователя: ")
    return client.auth(singInId)


