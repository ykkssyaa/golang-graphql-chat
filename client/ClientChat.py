from gql import Client, gql
from gql.transport.websockets import WebsocketsTransport
from gql.transport.aiohttp import AIOHTTPTransport

class ClientChat:

    allMessages = []

    def __init__(self, URL):

        self.url = URL

        transport = AIOHTTPTransport(url="http://" + URL)
        self.client = Client(transport=transport, fetch_schema_from_transport=True, )

        self._authComplete = False

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

        async with Client(transport=transport, fetch_schema_from_transport=True) as session:
            async for message in session.subscribe(subscription, variable_values=args):
                print(message)
                self.allMessages.append(message)

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
        print(result)
