import asyncio
import time

from gql import Client, gql
from gql.transport.websockets import WebsocketsTransport
import logging

#session = None
logging.basicConfig(level=logging.INFO)


async def auth():
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

    logging.log(level=logging.INFO, msg="in auth function")

    args = {"UserID": "1"}
    transport = WebsocketsTransport(url="ws://" + "localhost:8080/query")

    async with Client(transport=transport, fetch_schema_from_transport=True) as session:
        logging.log(level=logging.INFO, msg="in session")
        async for message in session.subscribe(subscription, variable_values=args):
            logging.log(level=logging.INFO, msg="in in async for loop")
            print(message)
        logging.log(level=logging.INFO, msg="after async for loop")

    logging.log(level=logging.INFO, msg="after session")


async def main():
    print("Before auth")

    coro = auth()

    task = asyncio.create_task(auth())
    await asyncio.sleep(0)

    print("After auth")

    await asyncio.sleep(1)

    while True:
        await asyncio.sleep(1)
        inp = input(">>>")



if __name__ == '__main__':
    asyncio.run(main())
