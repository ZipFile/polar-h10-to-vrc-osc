#!/usr/bin/env python

min_hr = 32
max_hr = 192
delta_hr = max_hr - min_hr
fps = 60
lub_frames = fps // 6
dub_frames = fps - lub_frames
lub_time = lub_frames / fps
dub_time = dub_frames / fps


def target_duration(x):
    return fps / (min_hr + delta_hr * x)


def sA(x):
    return 1 + 0.25 * x


def sB(x):
    return dub_time / (target_duration(x) - lub_time / sA(x))


print("sA(0) =", sA(0))
print("sA(1) =", sA(1))
print("sB(0) =", sB(0))
print("sB(1) =", sB(1))