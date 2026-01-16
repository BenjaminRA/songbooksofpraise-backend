import json
import re

# Load guitar chords to get the chord names
with open('assets/chords/guitar/chords.json', 'r') as f:
    guitar_chords = json.load(f)

# Ukulele standard tuning: G C E A (strings 0-3, from low to high)
# Note values: C=0, C#/Db=1, D=2, D#/Eb=3, E=4, F=5, F#/Gb=6, G=7, G#/Ab=8, A=9, A#/Bb=10, B=11
UKULELE_TUNING = [7, 0, 4, 9]  # G, C, E, A

NOTE_MAP = {
    'C': 0, 'C#': 1, 'Db': 1, 'D': 2, 'D#': 3, 'Eb': 3, 'E': 4, 'F': 5,
    'F#': 6, 'Gb': 6, 'G': 7, 'G#': 8, 'Ab': 8, 'A': 9, 'A#': 10, 'Bb': 10, 'B': 11
}

def parse_chord_name(chord_name):
    """Extract root note and chord type from chord name."""
    # Match root note (can be one or two characters)
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

def generate_ukulele_voicings(root_note, chord_type, max_voicings=8):
    """Generate ukulele chord voicings for a given chord."""
    root_value = NOTE_MAP.get(root_note)
    if root_value is None:
        return []
    
    intervals = get_chord_intervals(chord_type)
    chord_notes = [(root_value + interval) % 12 for interval in intervals]
    
    voicings = []
    
    # Generate voicings by trying different fret combinations
    # We'll check frets 0-12 for each string
    for fret0 in range(13):
        for fret1 in range(13):
            for fret2 in range(13):
                for fret3 in range(13):
                    positions = [fret0, fret1, fret2, fret3]
                    
                    # Calculate the notes played
                    notes_played = [
                        (UKULELE_TUNING[i] + positions[i]) % 12 
                        for i in range(4)
                    ]
                    
                    # Check if all notes are in the chord
                    if all(note in chord_notes for note in notes_played):
                        # Check if the root note is present
                        if root_value % 12 in notes_played:
                            # Calculate finger span (max fret - min fret, excluding open strings)
                            non_zero_frets = [f for f in positions if f > 0]
                            if non_zero_frets:
                                span = max(non_zero_frets) - min(non_zero_frets)
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
                                        
                                        if len(voicings) >= max_voicings:
                                            return voicings
    
    return voicings

# Generate ukulele chords
ukulele_chords = {}
processed = 0

print("Generating ukulele chords...")
for chord_name in guitar_chords.keys():
    root, chord_type = parse_chord_name(chord_name)
    
    if root and root in NOTE_MAP:
        voicings = generate_ukulele_voicings(root, chord_type, max_voicings=8)
        
        if voicings:
            ukulele_chords[chord_name] = voicings
            processed += 1
            
            if processed % 500 == 0:
                print(f"Processed {processed} chords...")

# Save to file
with open('assets/chords/ukulele/chords.json', 'w') as f:
    json.dump(ukulele_chords, f, indent=2)

print(f"\nGenerated {len(ukulele_chords)} ukulele chords")
print(f"Total variations: {sum(len(v) for v in ukulele_chords.values())}")
print("\nSample chords generated:")
for i, chord in enumerate(list(ukulele_chords.keys())[:15]):
    print(f"  {chord}: {len(ukulele_chords[chord])} variations")
    if i < 3:
        print(f"    First voicing: {ukulele_chords[chord][0]['positions']}")
