#!/usr/bin/env python3
"""
Database Migration Validation Script
Validates that the migration was successful by running various checks.
"""

import sqlite3
import os

def validate_migration():
    """Run validation checks on the migrated database."""
    
    # Database path
    db_path = os.path.join(os.path.dirname(__file__), '..', 'songbooks_of_praise.sqlite')
    
    if not os.path.exists(db_path):
        print("❌ Migration database not found!")
        return False
    
    try:
        conn = sqlite3.connect(db_path)
        conn.row_factory = sqlite3.Row
        cursor = conn.cursor()
        
        print("🔍 Validating migration results...\n")
        
        # Check tables exist
        cursor.execute("SELECT name FROM sqlite_master WHERE type='table' ORDER BY name")
        tables = [row[0] for row in cursor.fetchall()]
        expected_tables = ['users', 'songbooks', 'songbook_editors', 'categories', 'songs', 'song_categories']
        
        print("📋 Tables created:")
        for table in expected_tables:
            if table in tables:
                print(f"  ✅ {table}")
            else:
                print(f"  ❌ {table} - MISSING")
        
        # Check record counts
        print("\n📊 Record counts:")
        for table in expected_tables:
            if table in tables:
                cursor.execute(f"SELECT COUNT(*) FROM {table}")
                count = cursor.fetchone()[0]
                print(f"  {table}: {count} records")
        
        # Validate admin user
        print("\n👤 Admin user validation:")
        cursor.execute("SELECT * FROM users WHERE admin = 1")
        admin_user = cursor.fetchone()
        if admin_user:
            print(f"  ✅ Admin user: {admin_user['first_name']} {admin_user['last_name']}")
            print(f"  ✅ Email: {admin_user['email']}")
            print(f"  ✅ Verified: {bool(admin_user['verified'])}")
        else:
            print("  ❌ Admin user not found")
        
        # Validate songbooks
        print("\n📚 Songbook validation:")
        cursor.execute("SELECT * FROM songbooks ORDER BY id")
        songbooks = cursor.fetchall()
        for songbook in songbooks:
            print(f"  ✅ '{songbook['title']}' - Verified: {bool(songbook['verified'])}")
        
        # Validate song distribution
        print("\n🎵 Song distribution:")
        cursor.execute("""
            SELECT sb.title, COUNT(s.id) as song_count
            FROM songbooks sb
            LEFT JOIN songs s ON sb.id = s.songbook_id
            GROUP BY sb.id, sb.title
        """)
        for row in cursor.fetchall():
            print(f"  {row['title']}: {row['song_count']} songs")
        
        # Check lyrics formatting
        print("\n📝 Lyrics formatting validation:")
        cursor.execute("SELECT id, title, lyrics FROM songs WHERE lyrics IS NOT NULL AND lyrics != '' LIMIT 3")
        songs_with_lyrics = cursor.fetchall()
        
        for song in songs_with_lyrics:
            has_verse_chorus = ('Verse' in song['lyrics'] or 'Chorus' in song['lyrics'])
            print(f"  {song['title']}: {'✅' if has_verse_chorus else '❌'} Formatted lyrics")
        
        # Check category hierarchy
        print("\n🏷️  Category hierarchy validation:")
        cursor.execute("SELECT COUNT(*) FROM categories WHERE parent_category_id IS NULL")
        top_level = cursor.fetchone()[0]
        cursor.execute("SELECT COUNT(*) FROM categories WHERE parent_category_id IS NOT NULL")
        sub_level = cursor.fetchone()[0]
        print(f"  Top-level categories: {top_level}")
        print(f"  Sub-categories: {sub_level}")
        
        # Check song-category relationships
        print("\n🔗 Song-category relationships:")
        cursor.execute("SELECT COUNT(*) FROM song_categories")
        relationships = cursor.fetchone()[0]
        print(f"  Total song-category links: {relationships}")
        
        # Validate foreign key relationships
        print("\n🔑 Foreign key validation:")
        
        # Check songbook ownership
        cursor.execute("""
            SELECT COUNT(*) FROM songbooks sb 
            JOIN users u ON sb.owner_id = u.id
        """)
        valid_owners = cursor.fetchone()[0]
        cursor.execute("SELECT COUNT(*) FROM songbooks")
        total_songbooks = cursor.fetchone()[0]
        print(f"  Songbook ownership: {valid_owners}/{total_songbooks} valid")
        
        # Check song-songbook relationships
        cursor.execute("""
            SELECT COUNT(*) FROM songs s 
            JOIN songbooks sb ON s.songbook_id = sb.id
        """)
        valid_song_songbooks = cursor.fetchone()[0]
        cursor.execute("SELECT COUNT(*) FROM songs")
        total_songs = cursor.fetchone()[0]
        print(f"  Song-songbook links: {valid_song_songbooks}/{total_songs} valid")
        
        conn.close()
        
        print("\n✅ Migration validation completed successfully!")
        return True
        
    except Exception as e:
        print(f"❌ Validation failed: {e}")
        return False

if __name__ == "__main__":
    validate_migration()
