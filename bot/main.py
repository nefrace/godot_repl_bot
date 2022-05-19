from urllib import response
from aiogram import types, Bot, Dispatcher
import requests
from dotenv import load_dotenv
import os

load_dotenv()
TOKEN = os.getenv("TOKEN")
HOST = os.getenv("HOST", "http://127.0.0.1:8080")
if not TOKEN:
    print("Token not specified")
    exit(1)

bot: Bot = Bot(TOKEN)

dp: Dispatcher = Dispatcher()

@dp.message(commands=['run'])
async def message_handler(message: types.Message):
    if not message.reply_to_message: return
    msg = message.reply_to_message
    if not msg.entities: return
    print(msg.text)

    for entity in msg.entities:
        if entity.type != "pre" and entity.type != "code": continue
        code = entity.extract(msg.text)
        resp = requests.post(HOST+"/run", data={'code': code})
        json = resp.json()
        if resp.status_code != 200:
            return await msg.answer("Произошла ошибка:\n" + json["stderr"])
        await msg.reply(f"Резульат\n```gd\n{json['result']}\n```", parse_mode="MarkdownV2")

def main():
    dp.run_polling(bot)

if __name__ == "__main__":
    main()