import os
import requests
import urllib.parse
import json
import time
import subprocess

import textwrap
import tkinter
from tkinter import *
from geopy.geocoders import Nominatim
from sys import platform
from threading import Thread

# JSON Background initial values
"""
"infobg": "#00007D",
"marqueebg": "#000055"
"""

class Marquee(tkinter.Canvas):
    def __init__(self, parent, text, margin=0, borderwidth=1, fps=30):
        super().__init__(parent, borderwidth=borderwidth)

        self.fps = fps

        text = self.create_text(0, -1000, text=text, font=('STAR JR', 20, "bold"), fill="white", anchor="w", tags=("text"))
        (x0, y0, x1, y1) = self.bbox("text")
        width = (x1 - x0) + (2 * margin) + (2 * borderwidth)
        height = (y1 - y0) + (2 * margin) + (2 * borderwidth)
        self.configure(width=width, height=height, bg=data["marqueebg"], highlightthickness=0)

        self.animate()

    def animate(self):
        (x0, y0, x1, y1) = self.bbox("text")
        if x1 < 0 or y0 < 0:
            x0 = self.winfo_width()
            y0 = int(self.winfo_height() / 2)
            self.coords("text", x0, y0)
            
        else:
            self.move("text", -1, 0)

        self.after_id = self.after(int(1000 / self.fps), self.animate)

def ExecMusic():
    if platform == "linux" or platform == "linux2":
        subprocess.Popen("python3 music.py", shell=True)
    else:
        subprocess.Popen("python music.py", shell=True)

mp3 = Thread(target=ExecMusic())
running = True
day = 0

with open(f"{os.getcwd()}/data.json","r") as file:
    data = json.load(file)
        
def clock():
    current = time.strftime("%I:%M:%S %p")
    digital.configure(text=current)
    root.after(1000, clock)

response = requests.get(f'https://api.weather.com/v1/location/{data["zip"]}:4:US/observations/current.json?language=en-US&units=e&apiKey=21d8a80b3d6b444998a80b3d6b1449d3').json()

temperature = response["observation"]["imperial"]["temp"]
wind = response["observation"]["imperial"]["wspd"]
dew = response["observation"]["imperial"]["dewpt"]
visibility = response["observation"]["imperial"]["vis"]
index = response["observation"]["uv_index"]
high = response["observation"]["imperial"]['temp_max_24hour']
low = response["observation"]["imperial"]['temp_min_24hour']
description = response["observation"]["phrase_32char"]

root = tkinter.Tk()

canvas = tkinter.Canvas(root, bg=data["infobg"], height=720, width=480, highlightthickness=0)
location = Nominatim(user_agent="CanariumApp")
locdata = location.geocode(f'{data["city"]}, {data["state"]}')
wxforecast = requests.get(f"https://api.weather.com/v3/aggcommon/v3-wx-forecast-daily-5day?geocodes={locdata.latitude},{locdata.longitude}&language=en-US&units=e&format=json&apiKey=e1f10a1e78da46f5b10a1e78da96f525").json()
foretext = textwrap.fill(wxforecast[0]["v3-wx-forecast-daily-5day"]["narrative"][0], width=25)

def wxupdate():
    response = requests.get(f'https://api.weather.com/v1/location/{data["zip"]}:4:US/observations/current.json?language=en-US&units=e&apiKey=21d8a80b3d6b444998a80b3d6b1449d3').json()
    temperature = response["observation"]["imperial"]["temp"]
    wind = response["observation"]["imperial"]["wspd"]
    dew = response["observation"]["imperial"]["dewpt"]
    visibility = response["observation"]["imperial"]["vis"]
    index = response["observation"]["uv_index"]
    foretext = textwrap.fill(wxforecast[0]["v3-wx-forecast-daily-5day"]["narrative"][0], width=25)
    current = time.strftime("%I:%M:%S %p")
    
    canvas.itemconfigure(temp, text=f'Temperature: {temperature}\u00b0F')
    canvas.itemconfigure(wspd, text=f'Wind Speed: {wind}mph')
    canvas.itemconfigure(dwpt, text=f'Dew Point: {dew}\u00b0')
    canvas.itemconfigure(vis, text=f'Visibility: {int(visibility)} miles')
    canvas.itemconfigure(uv, text=f'UV Index: {index}')
    print(f"Forecast has been updated at {current}.")
    root.after(600000, wxupdate)


if visibility != 1:
    if 10 > visibility > 1:
        vis = canvas.create_text(384, 200, text=f'Visibility: {int(visibility)} miles', font=("STAR JR", 40), fill="white")
    else:
        vis = canvas.create_text(395, 200, text=f'Visibility: {int(visibility)} miles', font=("STAR JR", 40), fill="white")
else:
    vis = canvas.create_text(395, 200, text=f'Visibility: {int(visibility)} mile', font=("STAR JR", 40), fill="white")

if wind >= 10 or wind < 100:
    if wind >= 11:
        wspd = canvas.create_text(360, 150, text=f'Wind Speed: {wind} mph', font=("STAR JR", 40), fill="white")
    else:    
        wspd = canvas.create_text(345, 150, text=f'Wind Speed: {wind} mph', font=("STAR JR", 40), fill="white")

temp = canvas.create_text(350, 50, text=f'Temperature: {temperature}\u00b0F', font=("STAR JR", 40), fill="white")
dwpt = canvas.create_text(320, 100, text=f'Dew Point: {dew}\u00b0', font=("STAR JR", 40), fill="white")
uv = canvas.create_text(860, 100, text=f'UV Index: {index}', font=("STAR JR", 40), fill="white")

channels = canvas.create_text(390, 320, text=f'LOCAL CHANNELS', font=("VCR OSD Mono", 45), fill="white")

channel1 = canvas.create_text(300, 380, text=f'{data["channels"][0]["id"]} {data["channels"][0]["name"]}', font=("STAR JR", 40), fill="white")
channel2 = canvas.create_text(620, 380, text=f'{data["channels"][1]["id"]} {data["channels"][1]["name"]}', font=("STAR JR", 40), fill="white")
channel3 = canvas.create_text(300, 440, text=f'{data["channels"][2]["id"]} {data["channels"][2]["name"]}', font=("STAR JR", 40), fill="white")
channel4 = canvas.create_text(640, 440, text=f'{data["channels"][3]["id"]} {data["channels"][3]["name"]}', font=("STAR JR", 40), fill="white")

title = canvas.create_text(325, 505, text="TODAY'S FORECAST", font=("STAR JR", 35), fill="white")

desc = canvas.create_text(435, 595, text=foretext, font=("STAR JR", 30), fill="white")

digital = Label(root, text="", font=("STAR JR", 40), fg="white", bg=data["infobg"])
digital.place(x=720, y=18)
clock()

root.title('Canarium')
root.geometry("1280x720")
root.resizable(False, False)
root.config(cursor="none")
root.attributes('-fullscreen', True)

if platform == "linux" or platform == "linux2":
    root.attributes('-zoomed', True)

wxupdate()
marquee = Marquee(root, text="You are now watching america's #1 uninterrupted weather forecast channel, Canarium.", margin=2, borderwidth=2, fps=60)
marquee.pack(side="bottom", fill="both")

canvas.pack(fill="both", expand=True)

mp3.start()
root.mainloop()