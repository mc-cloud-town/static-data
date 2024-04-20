import yaml
import aiohttp
import asyncio


async def fetch(s: aiohttp.ClientSession, url: str) -> dict:
    try:
        async with s.get(url) as response:
            return await response.json()
    except aiohttp.ClientError as e:
        print(e)
        return await fetch(s, url)


async def main():

    async with aiohttp.ClientSession() as s:
        print(yaml.safe_load(""))
        # url = "https://api.coinmarketcap.com/v1/ticker/"
        # data = await fetch(s, url)
        # print(json.dumps(data, indent=4))


asyncio.run(main())
