import json
import re

# Load guitar chords to get the chord names
with open('assets/chords/guitar/chords.json', 'r') as f:
    guitar_chords = json.load(f)

# Banjo standard tuning (Open G): gDGBD
# String 5 (short): G (high, typically not fretted or fretted at 5th fret)
# String 4: D
# String 3: G
# String 2: B
# String 1: D
# Note values: C=0, C#/Db=1, D=2, D#/Eb=3, E=4, F=5, F#/Gb=6, G=7, G#/Ab=8, A=9, A#/Bb=10, B=11
BANJO_TUNING = [7, 2, 7, 11, 2]  # G, D, G, B, D (strings 0-4, from 5th to 1st)

NOTE_MAP = {
    'C': 0, 'C#': 1, 'Db': 1, 'D': 2, 'D#': 3, 'Eb': 3, 'E': 4, 'F': 5,
    'F#': 6, 'Gb': 6, 'G': 7, 'G#': 8, 'Ab': 8, 'A': 9, 'A#': 10, 'Bb': 10, 'B': 11
}

def parse_chord_name(chord_name):
    """Extract root note and chord type from chord name."""
    match = re.match(r'^([A-G][#b]?)', chord_name)
    if match:
        root = match.group(1)
        chord_type = chord_name[len(root):]
        return root, chord_type
    return None, None

def get_chord_intervals(chord_type):
    """Return intervals for a chord type (semitones from root)."""
    chord_intervals = {
        '': [0, 4, 7],  # Major
        'm': [0, 3, 7],  # Minor
        '7': [0, 4, 7, 10],  # Dominant 7th
        'maj7': [0, 4, 7, 11],  # Major 7th
        'M7': [0, 4, 7, 11],  # Major 7th (alternative)
        'm7': [0, 3, 7, 10],  # Minor 7th
        'dim': [0, 3, 6],  # Diminished
        'dim7': [0, 3, 6, 9],  # Diminished 7th
        'aug': [0, 4, 8],  # Augmented
        'sus2': [0, 2, 7],  # Suspended 2nd
        'sus4': [0, 5, 7],  # Suspended 4th
        '6': [0, 4, 7, 9],  # Major 6th
        'm6': [0, 3, 7, 9],  # Minor 6th
        '9': [0, 4, 7, 10, 14],  # Dominant 9th
        'maj9': [0, 4, 7, 11, 14],  # Major 9th
        'm9': [0, 3, 7, 10, 14],  # Minor 9th
        'add9': [0, 4, 7, 14],  # Add 9
        '7sus4': [0, 5, 7, 10],  # 7sus4
    }
    
    # Handle slash chords
    if '/' in chord_type:
        chord_type = chord_type.split('/')[0]
    
    return chord_intervals.get(chord_type, [0, 4, 7])  # Default to major

def generate_banjo_voicings(root_note, chord_type, max_voicings=6):
    """Generate banjo chord voicings for a given chord."""
    root_value = NOTE_MAP.get(root_note)
    if root_value is None:
        return []
    
    intervals = get_chord_intervals(chord_type)
    chord_notes = [(root_value + interval) % 12 for interval in intervals]
    
    voicings = []
    
    # For banjo, the 5th string can be: open (0), at 5th fret (5), fretted higher, or muted (x)
    # Try common positions: muted, 0, 2, 5, 7
    for fret5_option in ['x', 0, 2, 5, 7]:
        if len(voicings) >= max_voicings:
            break
            
        # Try different fret combinations for strings 1-4
        # Search up to fret 12 
        for fret4 in range(13):
            if len(voicings) >= max_voicings:
                break
            for fret3 in range(13):
                if len(voicings) >= max_voicings:
                    break
                for fret2 in range(13):
                    if len(voicings) >= max_voicings:
                        break
                    for fret1 in range(13):
                        if len(voicings) >= max_voicings:
                            break
                            
                        positions = [fret5_option, fret4, fret3, fret2, fret1]
                        
                        # Calculate the notes played (skip muted strings)
                        notes_played = []
                        for i in range(5):
                            if positions[i] != 'x':
                                note = (BANJO_TUNING[i] + positions[i]) % 12
                                notes_played.append(note)
                        
                        # Check if all notes are in the chord
                        if all(note in chord_notes for note in notes_played):
                            # Check if the root note is present
                            if root_value % 12 in notes_played:
                                # Make sure we're playing at least 3 strings
                                if len(notes_played) >= 3:
                                    # Calculate finger span for strings 1-4 (excluding 5th string)
                                    non_zero_frets = [positions[i] for i in range(1, 5) if positions[i] not in ['x', 0]]
                                    if non_zero_frets:
                                        span = max(non_zero_frets) - min(non_zero_frets)
                                    else:
                                        span = 0
                                    
                                    # Only accept reasonable spans (0-4 frets)
                                    if span <= 4:
                                        # Generate fingering (simplified)
                                        fingerings = [str(p) for p in positions]
                                        
                                        voicing = {
                                            'positions': [str(p) for p in positions],
                                            'fingerings': [fingerings]
                                        }
                                        
                                        # Avoid duplicates
                                        if voicing not in voicings:
                                            voicings.append(voicing)
    
    return voicings

# Generate banjo chords
banjo_chords = {}
processed = 0
total_chords = len(guitar_chords)

print(f"Generating banjo chords for {total_chords} chord types...")
for chord_name in guitar_chords.keys():
    root, chord_type = parse_chord_name(chord_name)
    
    if root and root in NOTE_MAP:
        voicings = generate_banjo_voicings(root, chord_type, max_voicings=6)
        
        if voicings:
            banjo_chords[chord_name] = voicings
            processed += 1
            
            if processed % 100 == 0:
                print(f"Processed {processed}/{total_chords} chords... ({100*processed//total_chords}%)")
        else:
            # Still count it even if no voicings found
            processed += 1
            if processed % 100 == 0:
                print(f"Processed {processed}/{total_chords} chords... ({100*processed//total_chords}%)")
    else:
        processed += 1

# Save to file
with open('assets/chords/banjo/chords.json', 'w') as f:
    json.dump(banjo_chords, f, indent=2)

print(f"\nGenerated {len(banjo_chords)} banjo chords")
print(f"Total variations: {sum(len(v) for v in banjo_chords.values())}")
print("\nSample chords generated:")
for i, chord in enumerate(list(banjo_chords.keys())[:15]):
    print(f"  {chord}: {len(banjo_chords[chord])} variations")
    if i < 3:
        print(f"    First voicing: {banjo_chords[chord][0]['positions']}")
