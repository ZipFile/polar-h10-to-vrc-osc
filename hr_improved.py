#!/usr/bin/env python

import numpy as np
from scipy.io.wavfile import write, read
from scipy.signal import butter, lfilter, convolve

SAMPLE_RATE = 44100


def generate_ambient_hum(
    filename, 
    base_cycles=8, 
    amplitude=0.15, 
    noise_mix=0.6, 
    noise_cutoff=80, 
    sample_rate=SAMPLE_RATE,
):
    """
    Generate an ambient background hum that simulates the continuous sounds in the body.
    Uses exact cycle counts rather than time for perfect looping.
    
    Parameters:
    -----------
    filename : str
        Output filename
    base_cycles : int
        Number of cycles of the lowest frequency
    amplitude : float
        Overall amplitude of the output
    noise_mix : float
        Amount of noise to mix with the tonal elements (0.0 to 1.0)
    noise_cutoff : float
        Cutoff frequency for the lowpass filter on noise
    sample_rate : int
        Sample rate of the output
    """
    # Calculate a good length based on the number of cycles
    # Using a lower base frequency (16 Hz) for deep ambient sound
    base_freq = 16
    
    # Calculate duration based on cycles for a perfectly loopable sound
    duration = base_cycles / base_freq
    
    # Generate sample positions
    num_samples = int(sample_rate * duration)
    
    # Use a precise approach: create cycle-exact angles from 0 to 2π * cycles
    # This ensures the waveform starts and ends at exactly the same point
    angles = np.linspace(0, 2 * np.pi * base_cycles, num_samples, endpoint=False)
    
    # Create the base sine wave with exact cycles
    base_wave = np.sin(angles)
    
    # Create the first harmonic (exact multiple of base cycles)
    second_harmonic = np.sin(angles * 1.5)  # 1.5 times the frequency = 1.5 * base_cycles total cycles
    
    # Create the second harmonic (exact multiple of base cycles)
    third_harmonic = np.sin(angles * 2.0)  # 2 times the frequency = 2 * base_cycles total cycles
    
    # Mix the harmonics with decreasing amplitude
    wave = base_wave + 0.33 * second_harmonic + 0.2 * third_harmonic
    
    # Normalize the combined wave
    wave = wave / np.max(np.abs(wave))
    
    # Generate pink noise for texture - with improved loopability
    # We'll create overlapping noise segments and crossfade them for seamless looping
    
    # Create a longer buffer to allow for crossfading at the edges
    crossfade_duration = 0.05  # 50ms crossfade on each end
    crossfade_samples = int(crossfade_duration * sample_rate)
    
    # Generate a longer noise buffer to allow for proper crossfading
    extended_length = num_samples + 2 * crossfade_samples
    noise = np.random.normal(0, 1, extended_length)
    
    # Create pink-ish noise by applying a 1/f filter in frequency domain
    fft_noise = np.fft.rfft(noise)
    frequencies = np.fft.rfftfreq(extended_length, 1/sample_rate)
    # Avoid division by zero
    frequencies[0] = 1
    pink_filter = 1 / np.sqrt(frequencies)
    # Apply filter
    pink_noise = np.fft.irfft(fft_noise * pink_filter, extended_length)
    
    # Apply lowpass to the noise
    filtered_noise = lowpass(pink_noise, cutoff=noise_cutoff, sample_rate=sample_rate)
    
    # Extract the main portion with extra for crossfading
    main_portion = filtered_noise[crossfade_samples:crossfade_samples+num_samples]
    
    # Create crossfade windows (squared sine for smoother transition)
    fade_in = np.sin(np.linspace(0, np.pi/2, crossfade_samples))**2
    fade_out = np.sin(np.linspace(np.pi/2, 0, crossfade_samples))**2
    
    # Apply special treatment to the beginning and end for perfect looping
    # Crossfade the beginning of the buffer with the end of the original noise
    # and crossfade the end of the buffer with the beginning of the original noise
    start_portion = filtered_noise[:crossfade_samples]
    end_portion = filtered_noise[crossfade_samples+num_samples:extended_length]
    
    # Apply crossfades at both ends
    main_portion[:crossfade_samples] = (main_portion[:crossfade_samples] * fade_in + 
                                     end_portion * fade_out)
    
    main_portion[num_samples-crossfade_samples:] = (main_portion[num_samples-crossfade_samples:] * fade_out + 
                                                 start_portion * fade_in)
    
    # Normalize the noise
    main_portion = main_portion / np.max(np.abs(main_portion))
    
    # Mix the tonal elements with the crossfaded noise
    mixed_wave = (1 - noise_mix) * wave + noise_mix * main_portion
    
    # Apply a slow amplitude modulation with whole cycles
    # Use 2 complete cycles for the modulation
    mod_angles = np.linspace(0, 2 * np.pi * 2, num_samples, endpoint=False)
    amplitude_mod = 1.0 + 0.1 * np.sin(mod_angles)
    
    # Apply the modulation
    mixed_wave *= amplitude_mod
    
    # Normalize again
    mixed_wave = mixed_wave / np.max(np.abs(mixed_wave))

    # Apply final amplitude
    mixed_wave *= amplitude

    # Check loop quality by measuring discontinuity at the loop point
    loop_jump = np.abs(mixed_wave[-1] - mixed_wave[0])
    derivative_jump = np.abs((mixed_wave[1] - mixed_wave[0]) - (mixed_wave[0] - mixed_wave[-1]))
    print(f"Hum loop quality - amplitude jump: {loop_jump:.8f}, derivative jump: {derivative_jump:.8f}")
    
    # Convert to 16-bit audio
    wave_int16 = (mixed_wave * 32767).astype(np.int16)
    write(filename, sample_rate, wave_int16)
    
    return duration, loop_jump, derivative_jump


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
lub_harmonics = [(1.0, 1.0), (0.05, 2.0), (0.0025, 3.0)]
generate_thump_with_noise(
    "lub.wav", 
    freq=58,  # Slightly higher frequency for better audibility 
    duration=0.12,  # Slightly longer for more body
    amplitude=0.75,  
    noise_mix=0.25,  # Less noise for cleaner sound
    noise_cutoff=130,  # Gentle lowpass
    harmonics=lub_harmonics,
    compression=1.2,  # Very gentle compression
)

# "Dub" sound - second heart sound
dub_harmonics = [(1.0, 1.0), (0.025, 1.5), (0.0125, 2.0)]
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

min_hum = 1
min_idx = 0

from shutil import copyfile

for i in range(9000):
    hum, l, d = generate_ambient_hum(
        f"/tmp/hum.wav",
        base_cycles=128,
        amplitude=0.18,
        noise_mix=0.65,
        noise_cutoff=75,
    )

    if l < min_hum:
        min_hum = l
        min_idx = i
        copyfile(
            "/tmp/hum.wav",
            "/tmp/hum_best.wav",
        )

print(f"min hum: {min_idx}/{min_hum:.8f}")

