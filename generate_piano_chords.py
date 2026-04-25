import json
import re

# Load guitar chords to get the chord names
with open('assets/chords/guitar/chords.json', 'r') as f:
    guitar_chords = json.load(f)

# Note names
NOTE_NAMES = ['C', 'C#', 'D', 'D#', 'E', 'F', 'F#', 'G', 'G#', 'A', 'A#', 'B']
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
        '11': [0, 4, 7, 10, 14, 17],  # Dominant 11th
        '13': [0, 4, 7, 10, 14, 21],  # Dominant 13th
        'sus': [0, 5, 7],  # Sus (same as sus4)
    }
    
    # Handle slash chords
    if '/' in chord_type:
        chord_type = chord_type.split('/')[0]
    
    return chord_intervals.get(chord_type, [0, 4, 7])  # Default to major

def semitone_to_note(semitone, octave=4):
    """Convert semitone number to note name with octave."""
    note_index = semitone % 12
    note_octave = octave + (semitone // 12)
    note_name = NOTE_NAMES[note_index]
    return f"{note_name}{note_octave}"

def generate_piano_voicings(root_note, chord_type):
    """Generate piano chord voicings across 2 octaves."""
    root_value = NOTE_MAP.get(root_note)
    if root_value is None:
        return []
    
    intervals = get_chord_intervals(chord_type)
    
    voicings = []
    
    # Root position - starting from octave 3 and 4
    for start_octave in [3, 4]:
        notes = []
        for interval in intervals:
            semitone = root_value + interval
            note = semitone_to_note(semitone, start_octave)
            notes.append(note)
        voicings.append({"notes": notes})
    
    # First inversion (second note becomes bass)
    if len(intervals) >= 3:
        notes = []
        # Start with second interval note in lower octave
        second_interval = intervals[1]
        notes.append(semitone_to_note(root_value + second_interval, 3))
        # Add root note in higher position
        notes.append(semitone_to_note(root_value + 12, 4))
        # Add remaining notes
        for interval in intervals[2:]:
            notes.append(semitone_to_note(root_value + interval, 4))
        voicings.append({"notes": notes})
    
    # Second inversion (third note becomes bass)
    if len(intervals) >= 3:
        notes = []
        # Start with third interval note in lower octave
        third_interval = intervals[2]
        notes.append(semitone_to_note(root_value + third_interval, 3))
        # Add root note in higher position
        notes.append(semitone_to_note(root_value + 12, 4))
        # Add second note
        notes.append(semitone_to_note(root_value + intervals[1], 4))
        # Add remaining notes
        for interval in intervals[3:]:
            notes.append(semitone_to_note(root_value + interval, 4))
        voicings.append({"notes": notes})
    
    # Spread voicing (wider spacing)
    notes = []
    notes.append(semitone_to_note(root_value, 3))  # Root in bass
    for i, interval in enumerate(intervals[1:]):
        octave_offset = 4 if i < 2 else 5
        notes.append(semitone_to_note(root_value + interval, octave_offset))
    voicings.append({"notes": notes})
    
    return voicings

# Generate piano chords
piano_chords = {}
processed = 0

print("Generating piano chords...")
for chord_name in guitar_chords.keys():
    root, chord_type = parse_chord_name(chord_name)
    
    if root and root in NOTE_MAP:
        voicings = generate_piano_voicings(root, chord_type)
        
        if voicings:
            piano_chords[chord_name] = voicings
            processed += 1
            
            if processed % 500 == 0:
                print(f"Processed {processed} chords...")

# Save to file
with open('assets/chords/piano/chords.json', 'w') as f:
    json.dump(piano_chords, f, indent=2)

print(f"\nGenerated {len(piano_chords)} piano chords")
print(f"Total variations: {sum(len(v) for v in piano_chords.values())}")
print("\nSample chords generated:")
for i, chord in enumerate(list(piano_chords.keys())[:10]):
    print(f"  {chord}: {piano_chords[chord][0]['notes']}")
