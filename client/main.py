from ClientChat import ClientChat
import menus, auth
import asyncio
from time import sleep
from chatFunctionality import *

PORT = "8080"
URL = f"localhost:{PORT}/query"

loop = asyncio.get_event_loop()


def main():

    clientChat = ClientChat(URL)

    # Входим в меню авторизации
    while True:
        inpt = menus.auth_menu()
        print("\n\n")
        try:
            match inpt:
                case "1":
                    thr = auth.CreateUser(clientChat)

                    if thr is None:
                        continue

                case "2":
                    thr = auth.SignIn(clientChat)

                case "0":
                    exit(0)
        except Exception as e:
            print("Проблема в авторизации " + str(e))

        print("\n\n")

        if clientChat.currentID == "0":
            print("Error with auth, try again!\n\n\n")
        else:
            break

    sleep(0.5)
    print(f"{ clientChat.getName()}, Добро пожаловать! Ваш ID - {clientChat.currentID}")

    # Входим в меню для взаимодействия с чатом
    while True:
        inpt = menus.main_menu()

        print("\n\n")
        match inpt:
            case "1":
                allChats(clientChat)
            case "2":
                openChat(clientChat)
            case "3":
                createChat(clientChat)
            case "4":
                deleteChat(clientChat)
            case "5":
                allUsers(clientChat)
            case "0":
                exit(0)
        print("\n\n")


if __name__ == '__main__':
    main()

