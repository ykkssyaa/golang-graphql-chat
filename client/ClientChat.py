import threading
import logging

from gql import Client, gql
from gql.transport.websockets import WebsocketsTransport
from gql.transport.aiohttp import AIOHTTPTransport
from gql.transport import exceptions as transportexp

logging.basicConfig(level=logging.CRITICAL)


def print_message(payload, senderName, mesId, time):
    print(f"{senderName:<8}: {payload:<15}   ({time[11:19]}) |ID:{mesId}")


class ClientChat:

    def __init__(self, URL):

        self.url = URL

        transport = AIOHTTPTransport(url="http://" + URL)
        self.client = Client(transport=transport, fetch_schema_from_transport=True, )

        self._authComplete = False
        self.currentID = "0"
        self.currentChatID = "0"
        self._currentName = ""
        self.currentInterlocutor = "0"

        self.chats = []

    def auth(self, userID: str):
        subscription = gql("""
        subscription ($UserID:ID!){
            userJoined(user: $UserID){
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
        """)

        args = {"UserID": userID}

        self.currentID = userID
        self._authComplete = True
        transport = WebsocketsTransport(url="ws://" + self.url)

        session = Client(transport=transport, fetch_schema_from_transport=True)

        def message_handler():
            # TODO: Обработка сообщений, добавление в список чатов
            # TODO: Обработка ошибок внутри потока
            # TODO: При удалении отправляется нулевое сообщение, из-за чего вылезает ошибка
            for message in session.subscribe(subscription, variable_values=args):
                self._message_processing(message["userJoined"])

        try:
            thr = threading.Thread(target=message_handler)
            thr.daemon = True
            thr.start()

            return thr

        except transportexp.TransportQueryError as exp:
            self._authComplete = False
            self.currentID = 0
            raise Exception("Error with subscription")

    def _message_processing(self, message):
        logging.log(level=logging.INFO, msg=message)

        if message['chatID'] == self.currentChatID:
            print_message(payload=message['payload'],
                          senderName=message['sender']['name'],
                          mesId=message['id'],
                          time=message['time']
                          )

    def setCurrentChat(self, currentChatID: str):
        if not self._checkChat(currentChatID):
            raise Exception("no permissions to set this chat as current")

        self.currentChatID = currentChatID
        self.currentInterlocutor = self.getInterlocutor()

    def exitChat(self):
        self.currentChatID = "0"
        self.currentInterlocutor = "0"

    def getInterlocutor(self):
        if self.currentChatID == "0" or not self._authComplete:
            return "0"

        if self.currentInterlocutor != "0":
            return self.currentInterlocutor

        for chat in self.chats:
            if chat["id"] == self.currentChatID:
                if chat['user_1']['id'] != self.currentID:
                    return chat['user_1']['id']
                else:
                    return chat['user_2']['id']

        for chat in self.chats_of_user():
            if chat["id"] == self.currentChatID:
                if chat['user_1']['id'] != self.currentID:
                    return chat['user_1']['id']
                else:
                    return chat['user_2']['id']

        return "0"

    def createUser(self, name: str):
        mutation = gql("""
            mutation CreateUser ($name: String!) {
                createUser(input: {name: $name}) {
                    id
                    name
                }
            }
        """)

        args = {"name": name}

        result = self.client.execute(mutation, variable_values=args)

        return result["createUser"]

    def allUsers(self):

        query = gql("""
            query Users {
                users {
                    id
                    name
                }
            }
        """)

        result = self.client.execute(query)
        return result["users"]

    def getName(self, userID="") -> str:

        if len(userID) == 0:
            userID = self.currentID

        if (self.currentID == "0" or not self._authComplete) and userID == self.currentID:
            return ""
        if self._currentName != "" and userID == self.currentID:
            return self._currentName

        for user in self.allUsers():
            if user["id"] == userID:

                if userID == self.currentID:
                    self._currentName = user["name"]

                return user["name"]

        return "User not found"

    def chats_of_user(self):

        query = gql("""
        query ($user: ID) {
            chats (user: $user) {
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
        """)

        args = {"user": self.currentID}

        result = self.client.execute(query, variable_values=args)["chats"]
        self.chats = result

        return result

    def createChat(self, otherUserID: str):
        mutation = gql("""
        
        mutation ($user1:ID!, $user2:ID!) {
            createChat(input: {user1: $user1, user2: $user2}) {
                id
                user_1
                user_2
            }
        }
        """)

        args = {"user1": self.currentID, "user2": otherUserID}

        result = self.client.execute(mutation, variable_values=args)
        return result["createChat"]

    def _checkChat(self, chatID: str) -> bool:

        chats = self.chats_of_user()

        for chat in chats:
            if chat["id"] == chatID:
                return True

        return False

    def deleteChat(self, chatID: str):

        if not self._checkChat(chatID):
            raise Exception("no permissions to delete this chat")

        mutation = gql("""
        mutation ($chat:ID!) {
            deleteChat(id: $chat)
        }
        """)

        args = {"chat": chatID}

        result = self.client.execute(mutation, variable_values=args)

        return result["deleteChat"]

    def postMessage(self, payload: str):

        if payload.replace(" ", "") == "":
            raise Exception("empty message")

        mutation = gql("""
        mutation ($payload: String!, $sender: ID!, $receiver: ID!) {
            postMessage(input: {payload: $payload, sender: $sender, receiver: $receiver}) {
                id
                payload
                chatID
                time
            }
        }
        """)

        args = {"payload": payload, "sender": self.currentID, "receiver": self.getInterlocutor()}

        resultMessage = self.client.execute(mutation, variable_values=args)["postMessage"]

        return resultMessage

    messagesCountOnPage = 15

    def _loadMessages(self, sender, receiver):
        query = gql("""
        query ($receiver: ID!, $sender: ID!, $count: Int!, $after: ID) {
            messagesFromUser(input: 
            {receiver: $receiver, sender: $sender},
             first: $count, 
            after: $after) 
            {
                edges {
                    node {
                        id
                        payload
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
        """)

        args = {"count": self.messagesCountOnPage, "receiver": receiver,
                "sender": sender, "after": None}

        messages = []

        while True:
            result = self.client.execute(query, variable_values=args)["messagesFromUser"]

            messages += result["edges"]

            if not result["pageInfo"]["hasNextPage"]:
                break

            args["after"] = result["pageInfo"]["endCursor"]

        return messages

    def loadMessages(self):

        messages = self._loadMessages(self.currentID, self.getInterlocutor()) + \
                   self._loadMessages(self.getInterlocutor(), self.currentID)

        messages = list(map(lambda d: d['node'], messages))

        return sorted(messages, key=lambda d: d['id'])
