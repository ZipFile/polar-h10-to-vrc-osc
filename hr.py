#!/usr/bin/env python

import numpy as np
from scipy.io.wavfile import write
from scipy.signal import butter, lfilter

SAMPLE_RATE = 44100


def lowpass(data, cutoff=150, order=4, sample_rate=SAMPLE_RATE):
    nyq = 0.5 * sample_rate
    norm_cutoff = cutoff / nyq
    b, a = butter(order, norm_cutoff, btype="low", analog=False)
    return lfilter(b, a, data)


def generate_thump_with_noise(
    filename,
    freq=60,
    duration=0.1,
    amplitude=0.5,
    noise_mix=0.3,
    noise_cutoff=200,
    sample_rate=SAMPLE_RATE,
):
    t = np.linspace(0, duration, int(sample_rate * duration), endpoint=False)

    # Base sine thump
    sine_wave = np.sin(2 * np.pi * freq * t)

    # Add white noise and lowpass it
    noise = np.random.normal(0, 1, len(t))
    noise = lowpass(noise, cutoff=noise_cutoff, sample_rate=sample_rate)
    wave = sine_wave * (1 - noise_mix) + noise * noise_mix

    # Apply envelope
    envelope = np.hanning(len(wave))
    wave *= envelope
    wave *= amplitude

    wave_int16 = (wave * 32767).astype(np.int16)
    write(filename, sample_rate, wave_int16)


MODES = {
    "neutral": {
        "lub": {
            "freq": 60,
            "amplitude": 0.6,
            "noise_mix": 0.3,
        },
        "dub": {
            "freq": 80,
            "amplitude": 0.4,
            "noise_mix": 0.4,
        },
    },
    "calm": {
        "lub": {
            "freq": 58,
            "amplitude": 0.6,
            "noise_mix": 0.45, # 0.25 meh
            "noise_cutoff": 120,
        },
        "dub": {
            "freq": 75,
            "amplitude": 0.4,
            "noise_mix": 0.6, # 0.2 meh, 0.6 good
            "noise_cutoff": 150,
        },
    },
    "excited": {
        "lub": {
            "freq": 68,
            "amplitude": 0.7,
            "noise_mix": 0.4,
            "noise_cutoff": 150,
        },
        "dub": {
            "freq": 85,
            "amplitude": 0.5,
            "noise_mix": 0.35,
            "noise_cutoff": 200,
        },
    },
    "workout": {
        "lub": {
            "freq": 75,
            "amplitude": 0.8,
            "noise_mix": 0.4,
            "noise_cutoff": 150,
        },
        "dub": {
            "freq": 91,
            "amplitude": 0.6,
            "noise_mix": 0.5,
            "noise_cutoff": 200,
        },
    },
}


for mode_name, sounds in MODES.items():
    for sound_name, params in sounds.items():
        generate_thump_with_noise(f"{mode_name}-{sound_name}.wav", **params)

# Generate more organic lub and dub
#generate_thump_with_noise("lub.wav", freq=60, amplitude=0.6, noise_mix=0.3)
#generate_thump_with_noise("dub.wav", freq=80, amplitude=0.4, noise_mix=0.4)


