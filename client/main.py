from ClientChat import ClientChat
import menus
import asyncio


PORT = "8080"
URL = f"localhost:{PORT}/query"


async def main():

    inpt = menus.auth_menu()

    clientChat = ClientChat(URL)

    match inpt:
        case "1":
            clientChat.createUser("Egor")
        case "2":
            await clientChat.auth("1")
        case "0":
            exit(0)


if __name__ == '__main__':
    asyncio.run(main())
