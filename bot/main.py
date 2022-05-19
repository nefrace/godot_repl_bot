from urllib import response
from aiogram import types, Bot, Dispatcher
from aiogram.utils.text_decorations import markdown_decoration

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


def escape_md(text):
    return markdown_decoration.quote(text)


@dp.message(commands=['run'])
async def message_handler(message: types.Message):
    if not message.reply_to_message:
        return
    msg = message.reply_to_message
    if not msg.entities:
        return
    print(msg.text)

    for entity in msg.entities:
        if entity.type != "pre" and entity.type != "code":
            continue
        code = entity.extract(msg.text)
        resp = requests.post(HOST+"/run", data={'code': code})
        json = resp.json()
        if resp.status_code != 200:
            if json['stderr'].strip() == "":
                return await msg.reply("Таймаут скрипта")
            return await msg.reply("Произошла ошибка:\n" + json["stderr"])
        result = json['result'][:4095]
        answer = f"{result}"
        repr(result)
        if result.strip() == "":
            answer = "Скрипт выполнен, вывод пуст."
        await msg.reply(answer)


def main():
    dp.run_polling(bot)


if __name__ == "__main__":
    main()
