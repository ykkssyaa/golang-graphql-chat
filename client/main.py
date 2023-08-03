from ClientChat import ClientChat
import menus, auth
import asyncio
from time import sleep

PORT = "8080"
URL = f"localhost:{PORT}/query"

loop = asyncio.get_event_loop()


async def main():

    clientChat = ClientChat(URL)

    # Входим в меню авторизации
    while True:
        inpt = menus.auth_menu()
        match inpt:
            case "1":
                # TODO: Настроить создание пользователя
                auth.CreateUser(clientChat)
            case "2":
                singInId = input("Введите ID пользователя: ")
                task = asyncio.create_task(clientChat.auth(singInId))
                await asyncio.sleep(0)

            case "0":
                exit(0)

        await asyncio.sleep(1)
        if clientChat.currentID == "0":
            print("Error with auth, try again!\n\n\n")
        else:
            break

    await asyncio.sleep(0)
    print(f"{await clientChat.getName()}, Добро пожаловать! Ваш ID - {clientChat.currentID}")

    # Входим в меню для взаимодействия с чатом
    while True:
        await asyncio.sleep(0)
        inpt = menus.main_menu()

        match inpt:
            case "1":
                pass
            case "2":
                pass
            case "3":
                pass
            case "4":
                pass
            case "0":
                # TODO: при выходе вызывается ошибка, что task не закрыто
                exit(0)


if __name__ == '__main__':
    asyncio.run(main())

    exit(0)

