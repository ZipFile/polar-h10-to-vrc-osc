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
    harmonics=None,
    compression=1.0,
):
    t = np.linspace(0, duration, int(sample_rate * duration), endpoint=False)

    # Base sine thump with subtle harmonics for natural sound
    if harmonics is None:
        # Default harmonic structure - more gentle, natural falloff
        harmonics = [(1.0, 1.0), (0.3, 2.0), (0.1, 3.0)]
    
    # Generate the base wave with subtle harmonics
    wave = np.zeros_like(t)
    for amp, harmonic_mult in harmonics:
        # Use smoother sine waves with subtle phase shifts for more organic sound
        phase = np.random.uniform(0, 0.1)  # Tiny random phase for natural variation
        wave += amp * np.sin(2 * np.pi * freq * harmonic_mult * t + phase)
    
    # Normalize the combined wave
    wave = wave / max(abs(wave)) if max(abs(wave)) > 0 else wave

    # Add pink noise instead of white noise for more natural texture
    # Pink noise has more energy in lower frequencies (1/f spectrum)
    noise_length = len(t)
    noise = np.random.normal(0, 1, noise_length)
    # Create pink-ish noise by applying a 1/f filter in frequency domain
    fft_noise = np.fft.rfft(noise)
    frequencies = np.fft.rfftfreq(noise_length, 1/sample_rate)
    # Avoid division by zero
    frequencies[0] = 1
    pink_filter = 1 / np.sqrt(frequencies)
    # Apply filter
    pink_noise = np.fft.irfft(fft_noise * pink_filter, noise_length)
    # Normalize
    pink_noise = pink_noise / np.max(np.abs(pink_noise))
    
    # Apply gentle lowpass to the noise
    noise = lowpass(pink_noise, cutoff=noise_cutoff, sample_rate=sample_rate)
    
    # Add subtle mid-frequency component for audibility on small speakers
    # Use a gentler approach with a softer sine
    mid_boost = np.sin(2 * np.pi * (freq * 2) * t) * 0.15
    noise = noise + mid_boost
    noise = noise / max(abs(noise)) if max(abs(noise)) > 0 else noise
    
    # Mix the harmonic wave with the noise
    wave = wave * (1 - noise_mix) + noise * noise_mix

    # Apply natural envelope - smoother attack and decay
    # Smoother envelope models the natural sound of a heart more accurately
    # Based on the Hann window but with adjustable attack/decay
    attack = int(len(wave) * 0.15)  # 15% attack
    decay = len(wave) - attack
    # Create smoother envelope curve using cosine
    env = np.zeros_like(wave)
    env[:attack] = 0.5 * (1 - np.cos(np.pi * np.linspace(0, 1, attack)))
    env[attack:] = 0.5 * (1 + np.cos(np.pi * np.linspace(0, 1, decay)))
    wave *= env
    
    # Apply very gentle compression to increase audibility without harshness
    if compression > 1.0:
        # Use a gentler compression curve that doesn't distort the sound
        # This is more like a soft-knee compressor that maintains the timbre
        scale = 0.5 * compression  # Reduce the effect of compression
        wave = np.sign(wave) * (1 - np.exp(-scale * np.abs(wave)))
    
    wave *= amplitude

    # Convert to 16-bit audio
    wave_int16 = (wave * 32767).astype(np.int16)
    write(filename, sample_rate, wave_int16)


# Generate more natural lub and dub with improved audibility
# "Lub" sound - first heart sound
lub_harmonics = [(1.0, 1.0), (0.25, 2.0), (0.07, 3.0)]
generate_thump_with_noise(
    "lub.wav", 
    freq=60,  # Slightly higher frequency for better audibility 
    duration=0.12,  # Slightly longer for more body
    amplitude=0.75,  
    noise_mix=0.25,  # Less noise for cleaner sound
    noise_cutoff=130,  # Gentle lowpass
    harmonics=lub_harmonics,
    compression=1.2,  # Very gentle compression
)

# "Dub" sound - second heart sound
dub_harmonics = [(1.0, 1.0), (0.15, 1.5), (0.08, 2.0)]
generate_thump_with_noise(
    "dub.wav", 
    freq=70, 
    duration=0.09,  # Shorter than lub
    amplitude=0.6,
    noise_mix=0.2,
    noise_cutoff=150,
    harmonics=dub_harmonics,
    compression=1.15,  # Even gentler compression
)
