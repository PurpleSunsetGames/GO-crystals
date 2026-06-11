import pandas
import matplotlib.pyplot as plt
import numpy as np
from time import sleep

points = []
with open("output.csv", "r") as file:
    for line in file:
        if len(line) == 0:
            break
        dat = np.reshape(np.array(line.split(",")), (2,-1))
        points.append(dat.astype(float))

i=0
while i<len(points):
    plt.scatter(points[i][0], points[i][1])
    plt.show()
    sleep(1)
    i+=1
