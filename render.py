import pandas
import matplotlib.pyplot as plt
import numpy as np
from time import sleep
import imageio
import os

points = []
with open("output.csv", "r") as file:
    for line in file:
        if len(line) == 0:
            break
        dat = np.reshape(np.array(line.split(",")), (-1,2)).T
        points.append(dat.astype(float))

i=0
while i<len(points):
    plt.scatter(points[i][0], points[i][1], c="black")
    plt.xlim(-10, 10)
    plt.ylim(-10,10)
    plt.savefig("output/fig" + str(i), dpi=400)
    plt.clf()
    i+=1

writer = imageio.get_writer('output/out.mp4', mode="I", fps=24)
for i in range(0, len(points)):
    writer.append_data(imageio.imread("output/fig"+str(i)+".png"))
    os.remove("output/fig"+str(i)+".png")


