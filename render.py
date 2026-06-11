import pandas
import matplotlib.pyplot as plt
import numpy as np
from time import sleep
import imageio

points = []
with open("output.csv", "r") as file:
    for line in file:
        if len(line) == 0:
            break
        dat = np.reshape(np.array(line.split(",")), (2,-1))
        points.append(dat.astype(float))

i=0
while i<len(points):
    plt.scatter(points[i][0], points[i][1], c="black")
    plt.xlim(-50, 50)
    plt.ylim(-50,50)
    plt.savefig("fig" + str(i))
    plt.clf()
    i+=1

writer = imageio.get_writer('out.mp4', mode="I", fps=2)
for i in range(0, len(points)):
    writer.append_data(imageio.imread("fig"+str(i)+".png"))


