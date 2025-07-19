#!/usr/bin/env python

import numpy as np
from scipy.io.wavfile import write, read
from scipy.signal import butter, lfilter, convolve

SAMPLE_RATE = 44100


def apply_reverb(signal, reverb_level=0.2, decay=0.5, sample_rate=SAMPLE_RATE):
    """
    Apply a simple room reverb effect to an audio signal.
    
    Parameters:
    -----------
    signal : numpy array
        Input audio signal
    reverb_level : float
        Level of reverb effect (0.0 to 1.0)
    decay : float
        Decay time in seconds
    sample_rate : int
        Sample rate of the audio
        
    Returns:
    --------
    numpy array
        Signal with reverb applied
    """
    # Skip if no reverb requested
    if reverb_level <= 0:
        return signal
    
    # Create impulse response for simple room reverb
    decay_samples = int(decay * sample_rate)
    
    # Exponential decay
    impulse_response = np.exp(-np.linspace(0, 4, decay_samples))
    
    # Add some early reflections for a more natural sound
    for i in range(3, 10):
        idx = int(i * 0.01 * sample_rate)
        if idx < decay_samples:
            impulse_response[idx] += 0.5 * np.exp(-i * 0.5)
    
    # Normalize impulse response
    impulse_response = impulse_response / np.sum(impulse_response)
    
    # Apply the reverb via convolution
    reverb_signal = convolve(signal, impulse_response, mode='full')[:len(signal)]
    
    # Mix dry and wet signals
    return signal * (1 - reverb_level) + reverb_signal * reverb_level


def generate_ambient_hum(
    filename, 
    base_cycles=8, 
    amplitude=0.15, 
    noise_mix=0.6, 
    noise_cutoff=80, 
    sample_rate=SAMPLE_RATE,
    reverb_amount=0.2,
    reverb_decay=0.1
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
    reverb_amount : float
        Amount of reverb to add (0.0 to 1.0)
    reverb_decay : float
        Decay time of the reverb in seconds
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
    
    # Generate pink noise for texture
    noise = np.random.normal(0, 1, num_samples)
    
    # Create pink-ish noise by applying a 1/f filter in frequency domain
    fft_noise = np.fft.rfft(noise)
    frequencies = np.fft.rfftfreq(num_samples, 1/sample_rate)
    # Avoid division by zero
    frequencies[0] = 1
    pink_filter = 1 / np.sqrt(frequencies)
    # Apply filter
    pink_noise = np.fft.irfft(fft_noise * pink_filter, num_samples)
    
    # Apply lowpass to the noise
    filtered_noise = lowpass(pink_noise, cutoff=noise_cutoff, sample_rate=sample_rate)
    
    # Normalize the noise
    filtered_noise = filtered_noise / np.max(np.abs(filtered_noise))
    
    # Mix the tonal elements with the noise
    mixed_wave = (1 - noise_mix) * wave + noise_mix * filtered_noise
    
    # Apply a slow amplitude modulation with whole cycles
    # Use 2 complete cycles for the modulation
    mod_angles = np.linspace(0, 2 * np.pi * 2, num_samples, endpoint=False)
    amplitude_mod = 1.0 + 0.1 * np.sin(mod_angles)
    
    # Apply the modulation
    mixed_wave *= amplitude_mod
    
    # Normalize again
    mixed_wave = mixed_wave / np.max(np.abs(mixed_wave))
    
    # Apply reverb if requested - reverb can affect looping if not careful
    if reverb_amount > 0:
        # Use a very short decay to avoid affecting the loop points
        safe_decay = min(reverb_decay, duration * 0.05)
        mixed_wave = apply_reverb(mixed_wave, reverb_level=reverb_amount, decay=safe_decay, sample_rate=sample_rate)
    
    # Apply final amplitude
    mixed_wave *= amplitude
    
    # Check loop quality by measuring discontinuity at the loop point
    loop_jump = np.abs(mixed_wave[-1] - mixed_wave[0])
    derivative_jump = np.abs((mixed_wave[1] - mixed_wave[0]) - (mixed_wave[0] - mixed_wave[-1]))
    print(f"Hum loop quality - amplitude jump: {loop_jump:.8f}, derivative jump: {derivative_jump:.8f}")
    
    # Convert to 16-bit audio
    wave_int16 = (mixed_wave * 32767).astype(np.int16)
    write(filename, sample_rate, wave_int16)
    
    return duration  # Return duration for reference when combining sounds

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
    reverb_amount=0.0,
    reverb_decay=0.1,
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
    
    # Apply reverb if requested
    if reverb_amount > 0:
        wave = apply_reverb(wave, reverb_level=reverb_amount, decay=reverb_decay, sample_rate=sample_rate)

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
    reverb_amount=0.15,  # Subtle reverb for spatial depth
    reverb_decay=0.08,  # Short decay for "room" feel
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
    reverb_amount=0.12,  # Slightly less reverb than lub
    reverb_decay=0.06,  # Shorter decay
)

# Generate ambient background hum with perfect looping
hum_duration = generate_ambient_hum(
    "hum.wav",
    base_cycles=12,  # Use exactly 12 cycles for perfect looping
    amplitude=0.2,  # Quieter than the heart sounds
    noise_mix=0.6,  # More noise content for "air" quality
    noise_cutoff=80,  # Very gentle lowpass
)

print(f"Generated hum with duration: {hum_duration:.2f} seconds")

# Generate a longer perfectly loopable hum for extended use
long_hum_duration = generate_ambient_hum(
    "hum_loop.wav",
    base_cycles=48,  # Use more cycles for a longer loop
    amplitude=0.18,  # Slightly quieter for extended listening
    noise_mix=0.65,
    noise_cutoff=75,
)

print(f"Generated longer hum with duration: {long_hum_duration:.2f} seconds")


def combine_heartbeat_with_hum(
    output_filename="heartbeat.wav",
    lub_file="lub.wav",
    dub_file="dub.wav",
    hum_file="hum.wav",
    bpm=60,
    sample_rate=SAMPLE_RATE,
    duration=5.0
):
    """
    Combine lub, dub, and ambient hum into a realistic heartbeat sequence.
    
    Parameters:
    -----------
    output_filename : str
        Output wav file name
    lub_file : str
        File containing the "lub" sound
    dub_file : str
        File containing the "dub" sound
    hum_file : str
        File containing the ambient hum sound
    bpm : float
        Beats per minute for the heart rate
    sample_rate : int
        Sample rate of the output
    duration : float
        Duration of the output file in seconds
    """
    # Read the wav files
    lub_rate, lub = read_wav(lub_file)
    dub_rate, dub = read_wav(dub_file)
    hum_rate, hum = read_wav(hum_file)
    
    # Create output buffer
    total_samples = int(duration * sample_rate)
    output = np.zeros(total_samples)
    
    # Calculate timing
    beat_interval = 60.0 / bpm  # Time for one full beat in seconds
    lub_dub_gap = beat_interval * 0.25  # Gap between lub and dub (systole)
    
    # Calculate sample positions
    beat_samples = int(beat_interval * sample_rate)
    lub_dub_samples = int(lub_dub_gap * sample_rate)
    
    # Add heartbeats
    current_sample = 0
    while current_sample < total_samples - len(lub) - len(dub):
        # Add lub at current position
        end_idx = min(current_sample + len(lub), total_samples)
        output[current_sample:end_idx] += lub[:end_idx-current_sample]
        
        # Add dub after lub
        dub_start = current_sample + lub_dub_samples
        dub_end = min(dub_start + len(dub), total_samples)
        if dub_end > dub_start:
            output[dub_start:dub_end] += dub[:dub_end-dub_start]
        
        # Move to next beat
        current_sample += beat_samples
    
    # Loop the hum as needed to fill the entire duration
    hum_loops_needed = int(np.ceil(duration / (len(hum) / hum_rate)))
    hum_extended = np.tile(hum, hum_loops_needed)[:total_samples]
    
    # Combine with hum - apply slight ducking to the hum during heartbeats
    hum_envelope = np.ones(total_samples)
    
    # Create ducking at each heartbeat
    current_sample = 0
    while current_sample < total_samples - len(lub) - len(dub):
        # Create a ducking region around each heartbeat
        duck_start = max(0, current_sample - int(0.03 * sample_rate))
        duck_end = min(total_samples, current_sample + len(lub) + lub_dub_samples + len(dub) + int(0.05 * sample_rate))
        duck_length = duck_end - duck_start
        
        # Create a gentle ducking curve
        duck_curve = 0.8 + 0.2 * np.cos(np.linspace(0, np.pi, duck_length))
        hum_envelope[duck_start:duck_end] = np.minimum(hum_envelope[duck_start:duck_end], duck_curve)
        
        # Move to next beat
        current_sample += beat_samples
    
    # Apply envelope to the hum
    hum_extended *= hum_envelope
    
    # Mix the signals
    output = output + hum_extended
    
    # Normalize to prevent clipping
    max_amplitude = np.max(np.abs(output))
    if max_amplitude > 0:
        output = output / max_amplitude * 0.95  # Leave some headroom
    
    # Apply a short fade-in and fade-out
    fade_samples = int(0.03 * sample_rate)
    if fade_samples > 0:
        fade_in = np.linspace(0, 1, fade_samples)
        fade_out = np.linspace(1, 0, fade_samples)
        
        output[:fade_samples] *= fade_in
        output[-fade_samples:] *= fade_out
    
    # Convert to 16-bit
    output_int16 = (output * 32767).astype(np.int16)
    write(output_filename, sample_rate, output_int16)
    
    print(f"Generated {duration}-second heartbeat at {bpm} BPM")


def read_wav(filename):
    """Helper function to read WAV files"""
    from scipy.io.wavfile import read
    rate, data = read(filename)
    # Convert to float in the -1.0 to 1.0 range
    if data.dtype == np.int16:
        data = data.astype(np.float64) / 32767.0
    elif data.dtype == np.int32:
        data = data.astype(np.float64) / 2147483647.0
    elif data.dtype == np.uint8:
        data = (data.astype(np.float64) - 128) / 128.0
    return rate, data


# Generate the complete heartbeat
combine_heartbeat_with_hum(
    output_filename="heartbeat.wav",
    bpm=70,  # Normal resting heart rate
    duration=10.0  # Generate a 10-second sample
)

print("Generated all sound files: lub.wav, dub.wav, hum.wav, hum_loop.wav, and heartbeat.wav")
