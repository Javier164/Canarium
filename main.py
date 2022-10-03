import os, sys
import requests
import urllib.parse
import json
import time
import textwrap

import tkinter
import pygame
from tkinter import *
from pygame import mixer
from datetime import datetime
from geopy.geocoders import Nominatim

pygame.init()

running = True
day = 0

with open(f"{os.getcwd()}/data.json","r") as file:
    data = json.load(file)

def wxupdate():
    response = requests.get(f'https://api.weather.com/v1/location/{data["zip"]}:4:US/observations/current.json?language=en-US&units=e&apiKey=21d8a80b3d6b444998a80b3d6b1449d3').json()
    temperature = response["observation"]["imperial"]["temp"]
    wind = response["observation"]["imperial"]["wspd"]
    dew = response["observation"]["imperial"]["dewpt"]
    visibility = response["observation"]["imperial"]["vis"]
    index = response["observation"]["uv_index"]
    current = time.strftime("%I:%M:%S %p")
    
    canvas.itemconfigure(temp, text=f'Temperature: {temperature}\u00b0F')
    canvas.itemconfigure(wspd, text=f'Wind Speed: {wind}mph')
    canvas.itemconfigure(dwpt, text=f'Dew Point: {dew}\u00b0')
    canvas.itemconfigure(vis, text=f'Visibility: {int(visibility)} miles')
    canvas.itemconfigure(uv, text=f'UV Index: {index}')
    print(f"Forecast has been updated at {current}.")
    root.after(600000, wxupdate)

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

canvas = tkinter.Canvas(root, bg="#00007D", height=720, width=480, highlightthickness=0)
location = Nominatim(user_agent="CanariumApp")
locdata = location.geocode(f'{data["city"]}, {data["state"]}')
wxforecast = requests.get(f"https://api.weather.com/v3/aggcommon/v3-wx-forecast-daily-5day?geocodes={locdata.latitude},{locdata.longitude}&language=en-US&units=e&format=json&apiKey=e1f10a1e78da46f5b10a1e78da96f525").json()
foretext = textwrap.fill(wxforecast[0]["v3-wx-forecast-daily-5day"]["narrative"][0], width=25)

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
else:
    wspd = canvas.create_text(350, 150, text=f'Wind Speed: {wind} mph', font=("STAR JR", 40), fill="white")

temp = canvas.create_text(350, 50, text=f'Temperature: {temperature}\u00b0F', font=("STAR JR", 40), fill="white")
dwpt = canvas.create_text(320, 100, text=f'Dew Point: {dew}\u00b0', font=("STAR JR", 40), fill="white")
uv = canvas.create_text(860, 100, text=f'UV Index: {index}', font=("STAR JR", 40), fill="white")

channels = canvas.create_text(385, 320, text=f'LOCAL CHANNELS', font=("VCR OSD Mono", 45), fill="white")

channel1 = canvas.create_text(300, 380, text=f'{data["channels"][0]["id"]} {data["channels"][0]["name"]}', font=("STAR JR", 40), fill="white")
channel2 = canvas.create_text(620, 380, text=f'{data["channels"][1]["id"]} {data["channels"][1]["name"]}', font=("STAR JR", 40), fill="white")
channel3 = canvas.create_text(300, 440, text=f'{data["channels"][2]["id"]} {data["channels"][2]["name"]}', font=("STAR JR", 40), fill="white")
channel4 = canvas.create_text(640, 440, text=f'{data["channels"][3]["id"]} {data["channels"][3]["name"]}', font=("STAR JR", 40), fill="white")

title = canvas.create_text(325, 535, text="TODAY'S FORECAST", font=("STAR JR", 35), fill="white")
    
desc = canvas.create_text(455, 635, text=foretext, font=("STAR JR", 35), fill="white")

digital = Label(root, text="", font=("STAR JR", 40), fg="white", bg="#00007D")
digital.place(x=720, y=18)
clock()

root.title('Canarium')
root.geometry("720x480")
root.resizable(False, False)
root.config(cursor="none")
root.attributes('-fullscreen', True)

wxupdate()

canvas.pack(fill="both", expand=True)

"""
( ͡° ͜ʖ ͡°)

mixer.music.load(f'{os.getcwd()}/assets/music/cafe.mp3')
mixer.music.play(loops=1)
"""

root.mainloop()
