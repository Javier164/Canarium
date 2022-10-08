import pygame
import random
import os

directory = 'music'
mp3 = [f for f in os.listdir(f"{os.getcwd()}/assets/{directory}") if f.endswith('.mp3')]
pygame.mixer.pre_init(frequency=48000, size=-16, channels=2)
pygame.init()

def music():
    current = []
    if not pygame.mixer.music.get_busy():
        if not current:
           current = mp3[:]
           random.shuffle(current)
        song = current.pop(0)
        pygame.mixer.music.load(os.path.join(f"{os.getcwd()}/assets/{directory}/", song))
        for item in range(len(current)):
            pygame.mixer.music.queue(os.path.join(f"{os.getcwd()}/assets/{directory}/", current[item]))
 
        pygame.mixer.music.play()
        
run = True
while run:
    pygame.time.Clock().tick(100)
    for event in pygame.event.get():
        if event.type == pygame.QUIT:
            run = False
            break
        
    music()