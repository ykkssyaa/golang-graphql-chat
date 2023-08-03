import asyncio
import threading

from gql import Client, gql
from gql.transport.websockets import WebsocketsTransport
from gql.transport.aiohttp import AIOHTTPTransport
from gql.transport import exceptions as transportexp
import logging

logging.basicConfig(level=logging.INFO)

tasks = set()

class ClientChat:

    allMessages = []

    def __init__(self, URL):

        self.url = URL

        transport = AIOHTTPTransport(url="http://" + URL)
        self.client = Client(transport=transport, fetch_schema_from_transport=True, )

        self._authComplete = False
        self.currentID = 0
        self._currentName = ""

    async def auth(self, userID: str):
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

        transport = WebsocketsTransport(url="ws://" + self.url)

        # TODO: Обработка сообщений, добавление в список чатов
        async with Client(transport=transport, fetch_schema_from_transport=True) as session:
            try:

                self._authComplete = True
                self.currentID = userID

                async for message in session.subscribe(subscription, variable_values=args):
                    await self._message_processing(message)

            except transportexp.TransportQueryError as exp:
                self._authComplete = False
                self.currentID = 0

                raise Exception("Error with subscription" + exp.data)

    async def createUser(self, name: str):
        mutation = gql("""
            mutation CreateUser ($name: String!) {
                createUser(input: {name: $name}) {
                    id
                    name
                }
            }
        """)

        args = {"name": name}

        result = await self.client.execute_async(mutation, variable_values=args)

        return result["createUser"]

    async def allUsers(self):

        query = gql("""
            query Users {
                users {
                    id
                    name
                }
            }
        """)

        result = await self.client.execute_async(query)

        return result["users"]

    async def getName(self) -> str:

        if self.currentID == 0 or not self._authComplete:
            return ""
        if self._currentName != "":
            return self._currentName

        for user in await self.allUsers():
            if user["id"] == self.currentID:
                self._currentName = user["name"]
                return self._currentName

        return "User not found"

    async def _message_processing(self, message):
        print(message)

