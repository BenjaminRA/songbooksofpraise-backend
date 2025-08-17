#!/usr/bin/env python3
"""
Database Migration Script
Migrates from himnos_coros.sqlite to songbooks_of_praise.sqlite
"""

import sqlite3
import hashlib
import logging
import os
import sys
import boto3
from botocore.exceptions import ClientError, NoCredentialsError
from datetime import datetime
from typing import Dict, List, Tuple, Optional
from dotenv import load_dotenv

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(levelname)s - %(message)s',
    handlers=[
        logging.FileHandler('migration.log'),
        logging.StreamHandler(sys.stdout)
    ]
)

logger = logging.getLogger(__name__)

# Load environment variables from parent directory
load_dotenv(dotenv_path='../.env')

class DatabaseMigrator:
    def __init__(self, source_db_path: str, target_db_path: str):
        self.source_db_path = source_db_path
        self.target_db_path = target_db_path
        self.source_conn = None
        self.target_conn = None
        self.s3_client = None
        
        # S3 configuration
        self.aws_region = os.getenv("AWS_S3_REGION")
        self.aws_access_key = os.getenv("AWS_S3_ACCESS_KEY")
        self.aws_secret_key = os.getenv("AWS_S3_SECRET_ACCESS_KEY")
        self.music_sheet_bucket = os.getenv("AWS_S3_MUSIC_SHEET_BUCKET")
        self.voices_bucket = os.getenv("AWS_S3_VOICES_BUCKET")
        
        # File paths
        self.music_sheets_path = "/Users/benjaminrodriguez/Documents/Algoritmos/server_nico/himnario/music_sheets/hymns"
        self.voices_path = "/Users/benjaminrodriguez/Documents/Algoritmos/server_nico/himnario/voices/voices"
        
    def get_s3_client(self):
        """Initialize and return S3 client."""
        if self.s3_client is None:
            try:
                self.s3_client = boto3.client(
                    's3',
                    region_name=self.aws_region,
                    aws_access_key_id=self.aws_access_key,
                    aws_secret_access_key=self.aws_secret_key
                )
                logger.info("S3 client initialized successfully")
            except Exception as e:
                logger.error(f"Failed to initialize S3 client: {e}")
                raise
        return self.s3_client
    
    def upload_file_to_s3(self, local_file_path: str, bucket: str, s3_key: str) -> Optional[str]:
        """Upload a file to S3 and return the URL."""
        try:
            if not os.path.exists(local_file_path):
                logger.warning(f"File not found: {local_file_path}")
                return None
                
            s3_client = self.get_s3_client()
            s3_client.upload_file(local_file_path, bucket, s3_key)
            
            # Generate the S3 URL
            url = f"https://{bucket}.s3.{self.aws_region}.amazonaws.com/{s3_key}"
            logger.debug(f"Uploaded {local_file_path} to {url}")
            return url
            
        except ClientError as e:
            logger.error(f"Error uploading {local_file_path} to S3: {e}")
            return None
        except Exception as e:
            logger.error(f"Unexpected error uploading {local_file_path}: {e}")
            return None
    
    def upload_music_sheet(self, song_id: int) -> Optional[str]:
        """Upload music sheet for a song and return the S3 URL."""
        music_sheet_file = os.path.join(self.music_sheets_path, f"{song_id}.jpg")
        if os.path.exists(music_sheet_file):
            s3_key = f"{song_id}.jpg"
            return self.upload_file_to_s3(music_sheet_file, self.music_sheet_bucket, s3_key)
        return None
    
    def upload_voices(self, song_id: int) -> Dict[str, Optional[str]]:
        """Upload voice files for a song and return the S3 URLs."""
        voices_folder = os.path.join(self.voices_path, str(song_id))
        voice_urls = {
            'voices_all': None,
            'voices_soprano': None,
            'voices_contralto': None,
            'voices_tenor': None,
            'voices_bass': None
        }
        
        if not os.path.exists(voices_folder):
            return voice_urls
        
        # Voice file mappings
        voice_files = {
            'voices_all': 'Todos.mp3',
            'voices_soprano': 'Soprano.mp3',
            'voices_contralto': 'ContraAlto.mp3',
            'voices_tenor': 'Tenor.mp3',
            'voices_bass': 'Bajo.mp3'
        }
        
        # S3 key mappings
        s3_key_mappings = {
            'voices_all': f"{song_id}_all.mp3",
            'voices_soprano': f"{song_id}_soprano.mp3",
            'voices_contralto': f"{song_id}_contralto.mp3",
            'voices_tenor': f"{song_id}_tenor.mp3",
            'voices_bass': f"{song_id}_bass.mp3"
        }
        
        for voice_type, filename in voice_files.items():
            local_file = os.path.join(voices_folder, filename)
            if os.path.exists(local_file):
                s3_key = s3_key_mappings[voice_type]
                voice_urls[voice_type] = self.upload_file_to_s3(local_file, self.voices_bucket, s3_key)
        
        return voice_urls
        
    def connect_databases(self):
        """Connect to source and target databases."""
        try:
            # Check if source database exists
            if not os.path.exists(self.source_db_path):
                raise FileNotFoundError(f"Source database not found: {self.source_db_path}")
            
            self.source_conn = sqlite3.connect(self.source_db_path)
            self.source_conn.row_factory = sqlite3.Row
            
            # Remove target database if it exists
            if os.path.exists(self.target_db_path):
                os.remove(self.target_db_path)
                logger.info(f"Removed existing target database: {self.target_db_path}")
            
            self.target_conn = sqlite3.connect(self.target_db_path)
            self.target_conn.row_factory = sqlite3.Row
            
            logger.info("Successfully connected to databases")
            
        except Exception as e:
            logger.error(f"Error connecting to databases: {e}")
            raise
    
    def create_new_schema(self):
        """Create the new database schema."""
        try:
            cursor = self.target_conn.cursor()
            
            # Create users table
            cursor.execute("""
                CREATE TABLE users (
                    id INTEGER PRIMARY KEY AUTOINCREMENT,
                    first_name VARCHAR(255) NOT NULL,
                    last_name VARCHAR(255) NOT NULL,
                    email VARCHAR(255) UNIQUE NOT NULL,
                    password VARCHAR(255) NOT NULL,
                    admin BOOLEAN DEFAULT FALSE,
                    editor BOOLEAN DEFAULT FALSE,
                    moderator BOOLEAN DEFAULT FALSE,
                    verified BOOLEAN DEFAULT FALSE,
                    created_at DATETIME NOT NULL,
                    updated_at DATETIME NOT NULL
                )
            """)
            
            # Create songbooks table
            cursor.execute("""
                CREATE TABLE songbooks (
                    id INTEGER PRIMARY KEY AUTOINCREMENT,
                    title VARCHAR(255) NOT NULL,
                    verified BOOLEAN DEFAULT FALSE,
                    in_verification BOOLEAN DEFAULT FALSE,
                    rejected BOOLEAN DEFAULT FALSE,
                    owner_id INTEGER NOT NULL,
                    created_at DATETIME NOT NULL,
                    updated_at DATETIME NOT NULL,
                    FOREIGN KEY (owner_id) REFERENCES users(id)
                )
            """)
            
            # Create songbook_editors table
            cursor.execute("""
                CREATE TABLE songbook_editors (
                    id INTEGER PRIMARY KEY AUTOINCREMENT,
                    songbook_id INTEGER NOT NULL,
                    user_id INTEGER NOT NULL,
                    created_at DATETIME NOT NULL,
                    updated_at DATETIME NOT NULL,
                    FOREIGN KEY (songbook_id) REFERENCES songbooks(id),
                    FOREIGN KEY (user_id) REFERENCES users(id)
                )
            """)
            
            # Create categories table (merging temas and sub_temas)
            cursor.execute("""
                CREATE TABLE categories (
                    id INTEGER PRIMARY KEY AUTOINCREMENT,
                    name VARCHAR(255) NOT NULL,
                    parent_category_id INTEGER,
                    songbook_id INTEGER,
                    created_at DATETIME NOT NULL,
                    updated_at DATETIME NOT NULL,
                    FOREIGN KEY (parent_category_id) REFERENCES categories(id),
                    FOREIGN KEY (songbook_id) REFERENCES songbooks(id)
                )
            """)
            
            # Create songs table (renamed from himnos with additional columns)
            cursor.execute("""
                CREATE TABLE songs (
                    id INTEGER PRIMARY KEY AUTOINCREMENT,
                    title VARCHAR(255) NOT NULL,
                    lyrics TEXT,
                    music_sheet TEXT,
                    music TEXT,
                    music_only TEXT,
                    youtube_url TEXT,
                    description TEXT,
                    number INTEGER,
                    voices_all TEXT,
                    voices_soprano TEXT,
                    voices_contralto TEXT,
                    voices_tenor TEXT,
                    voices_bass TEXT,
                    transpose INTEGER,
                    scroll_speed INTEGER,
                    songbook_id INTEGER NOT NULL,
                    created_at DATETIME NOT NULL,
                    updated_at DATETIME NOT NULL,
                    FOREIGN KEY (songbook_id) REFERENCES songbooks(id)
                )
            """)
            
            # Create song_categories table (merging tema_himnos and sub_tema_himnos)
            cursor.execute("""
                CREATE TABLE song_categories (
                    id INTEGER PRIMARY KEY AUTOINCREMENT,
                    song_id INTEGER NOT NULL,
                    category_id INTEGER NOT NULL,
                    created_at DATETIME NOT NULL,
                    updated_at DATETIME NOT NULL,
                    FOREIGN KEY (song_id) REFERENCES songs(id),
                    FOREIGN KEY (category_id) REFERENCES categories(id)
                )
            """)
            
            # Create session_tokens table for authentication
            cursor.execute("""
                CREATE TABLE session_tokens (
                    id INTEGER PRIMARY KEY AUTOINCREMENT,
                    access_uuid VARCHAR(255) UNIQUE NOT NULL,
                    refresh_uuid VARCHAR(255) UNIQUE NOT NULL,
                    user_id INTEGER NOT NULL,
                    at_exp INTEGER NOT NULL,
                    rt_exp INTEGER NOT NULL,
                    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
                    FOREIGN KEY (user_id) REFERENCES users(id)
                )
            """)
            
            # Create verification_tokens table for email verification
            cursor.execute("""
                CREATE TABLE verification_tokens (
                    id INTEGER PRIMARY KEY AUTOINCREMENT,
                    key VARCHAR(255) UNIQUE NOT NULL,
                    token TEXT NOT NULL,
                    user_id INTEGER NOT NULL,
                    expiration INTEGER NOT NULL,
                    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
                    FOREIGN KEY (user_id) REFERENCES users(id)
                )
            """)
            
            # Create reset_tokens table for password reset
            cursor.execute("""
                CREATE TABLE reset_tokens (
                    id INTEGER PRIMARY KEY AUTOINCREMENT,
                    key VARCHAR(255) UNIQUE NOT NULL,
                    token TEXT NOT NULL,
                    user_id INTEGER NOT NULL,
                    expiration INTEGER NOT NULL,
                    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
                    FOREIGN KEY (user_id) REFERENCES users(id)
                )
            """)
            
            self.target_conn.commit()
            logger.info("Successfully created new database schema")
            
        except Exception as e:
            logger.error(f"Error creating new schema: {e}")
            raise
    
    def get_current_timestamp(self) -> str:
        """Get current timestamp in the format used by the database."""
        return datetime.now().strftime('%Y-%m-%d %H:%M:%S')
    
    def create_initial_data(self):
        """Create initial users and songbooks."""
        try:
            cursor = self.target_conn.cursor()
            current_time = self.get_current_timestamp()
            
            # Create admin user
            hashed_password = "99adc231b045331e514a516b4b7680f588e3823213abe901738bc3ad67b2f6fcb3c64efb93d18002588d3ccc1a49efbae1ce20cb43df36b38651f11fa75678e8"  # Hash for "root"
            cursor.execute("""
                INSERT INTO users (first_name, last_name, email, password, admin, editor, moderator, verified, created_at, updated_at)
                VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
            """, ("Benjamín", "Rodríguez", "benjamin.gra720@gmail.com", hashed_password, True, True, True, True, current_time, current_time))
            
            admin_user_id = cursor.lastrowid
            logger.info(f"Created admin user with ID: {admin_user_id}")
            
            # Create songbooks
            cursor.execute("""
                INSERT INTO songbooks (title, verified, in_verification, rejected, owner_id, created_at, updated_at)
                VALUES (?, ?, ?, ?, ?, ?, ?)
            """, ("Himnos y Cánticos del Evangelio", True, False, False, admin_user_id, current_time, current_time))
            
            himnos_songbook_id = cursor.lastrowid
            logger.info(f"Created 'Himnos y Cánticos del Evangelio' songbook with ID: {himnos_songbook_id}")
            
            cursor.execute("""
                INSERT INTO songbooks (title, verified, in_verification, rejected, owner_id, created_at, updated_at)
                VALUES (?, ?, ?, ?, ?, ?, ?)
            """, ("Coros", True, False, False, admin_user_id, current_time, current_time))
            
            coros_songbook_id = cursor.lastrowid
            logger.info(f"Created 'Coros' songbook with ID: {coros_songbook_id}")
            
            self.target_conn.commit()
            
            return {
                'admin_user_id': admin_user_id,
                'himnos_songbook_id': himnos_songbook_id,
                'coros_songbook_id': coros_songbook_id
            }
            
        except Exception as e:
            logger.error(f"Error creating initial data: {e}")
            raise
    
    def migrate_categories(self, initial_data):
        """Migrate temas and sub_temas to categories table."""
        try:
            source_cursor = self.source_conn.cursor()
            target_cursor = self.target_conn.cursor()
            
            # Migrate top-level categories (temas)
            source_cursor.execute("SELECT * FROM temas ORDER BY id")
            temas = source_cursor.fetchall()
            
            tema_id_mapping = {}
            
            for tema in temas:
                target_cursor.execute("""
                    INSERT INTO categories (name, parent_category_id, songbook_id, created_at, updated_at)
                    VALUES (?, ?, ?, ?, ?)
                """, (tema['tema'], None, initial_data['himnos_songbook_id'], tema['createdAt'], tema['updatedAt']))
                
                new_categoria_id = target_cursor.lastrowid
                tema_id_mapping[tema['id']] = new_categoria_id
            
            logger.info(f"Migrated {len(temas)} top-level categories")
            
            # Migrate subcategories (sub_temas)
            source_cursor.execute("SELECT * FROM sub_temas ORDER BY id")
            sub_temas = source_cursor.fetchall()
            
            sub_tema_id_mapping = {}
            
            for sub_tema in sub_temas:
                parent_id = tema_id_mapping.get(sub_tema['tema_id'])
                if parent_id is None:
                    logger.warning(f"Parent tema not found for sub_tema {sub_tema['id']}")
                    continue
                
                target_cursor.execute("""
                    INSERT INTO categories (name, parent_category_id, songbook_id, created_at, updated_at)
                    VALUES (?, ?, ?, ?, ?)
                """, (sub_tema['sub_tema'], parent_id, initial_data['himnos_songbook_id'], sub_tema['createdAt'], sub_tema['updatedAt']))
                
                new_categoria_id = target_cursor.lastrowid
                sub_tema_id_mapping[sub_tema['id']] = new_categoria_id
            
            logger.info(f"Migrated {len(sub_temas)} subcategories")
            
            self.target_conn.commit()
            
            return {
                'tema_id_mapping': tema_id_mapping,
                'sub_tema_id_mapping': sub_tema_id_mapping
            }
            
        except Exception as e:
            logger.error(f"Error migrating categories: {e}")
            raise
    
    def get_song_lyrics(self, himno_id: int) -> str:
        """Get concatenated lyrics for a song from parrafos table."""
        try:
            cursor = self.source_conn.cursor()
            cursor.execute("""
                SELECT parrafo, coro FROM parrafos 
                WHERE himno_id = ? 
                ORDER BY id
            """, (himno_id,))
            
            parrafos = cursor.fetchall()
            
            if not parrafos:
                return ""
            
            lyrics_parts = []
            verse_count = 0
            chorus_count = 0
            
            for parrafo in parrafos:
                if parrafo['coro']:
                    chorus_count += 1
                    lyrics_parts.append(f"Chorus{' ' + str(chorus_count) if chorus_count > 1 else ''}")
                else:
                    verse_count += 1
                    lyrics_parts.append(f"Verse {verse_count}")
                
                lyrics_parts.append(parrafo['parrafo'])
                lyrics_parts.append("")  # Add empty line between sections
            
            # Remove the last empty line
            if lyrics_parts and lyrics_parts[-1] == "":
                lyrics_parts.pop()
            
            return "\n".join(lyrics_parts)
            
        except Exception as e:
            logger.error(f"Error getting lyrics for song {himno_id}: {e}")
            return ""
    
    def migrate_songs(self, songbook_ids: Dict[str, int]):
        """Migrate himnos to songs table."""
        try:
            source_cursor = self.source_conn.cursor()
            target_cursor = self.target_conn.cursor()
            
            # Get all himnos
            source_cursor.execute("SELECT * FROM himnos ORDER BY id")
            himnos = source_cursor.fetchall()
            
            himno_id_mapping = {}
            
            for himno in himnos:
                # Determine which songbook this song belongs to
                songbook_id = (songbook_ids['himnos_songbook_id'] 
                             if himno['id'] <= 517 
                             else songbook_ids['coros_songbook_id'])
                
                number = himno['id'] if himno['id'] <= 517 else None
                
                # Get lyrics for this song
                lyrics = self.get_song_lyrics(himno['id'])
                
                # Get transpose and scroll_speed safely
                transpose = himno['transpose'] if 'transpose' in himno.keys() else None
                scroll_speed = himno['scroll_speed'] if 'scroll_speed' in himno.keys() else None
                
                # Initialize S3 URLs
                music_sheet_url = None
                voices_urls = {
                    'voices_all': None,
                    'voices_soprano': None,
                    'voices_contralto': None,
                    'voices_tenor': None,
                    'voices_bass': None
                }
                
                # Upload files to S3 only for "Himnos y Cánticos del Evangelio" songbook
                if songbook_id == songbook_ids['himnos_songbook_id']:
                    logger.info(f"Processing files for hymn {himno['id']}: {himno['titulo']}")
                    
                    # Upload music sheet
                    music_sheet_url = self.upload_music_sheet(himno['id'])
                    if music_sheet_url:
                        logger.info(f"Uploaded music sheet for song {himno['id']}")
                    
                    # Upload voice files
                    voices_urls = self.upload_voices(himno['id'])
                    uploaded_voices = [k for k, v in voices_urls.items() if v is not None]
                    if uploaded_voices:
                        logger.info(f"Uploaded voices for song {himno['id']}: {uploaded_voices}")
                
                # Insert into songs table
                target_cursor.execute("""
                    INSERT INTO songs (
                        title, lyrics, music_sheet, music, music_only, youtube_url, description, 
                        number, voices_all, voices_soprano, voices_contralto, voices_tenor, voices_bass,
                        transpose, scroll_speed, songbook_id, created_at, updated_at
                    ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
                """, (
                    himno['titulo'], lyrics, music_sheet_url, None, None, None, None,
                    number, voices_urls['voices_all'], voices_urls['voices_soprano'], 
                    voices_urls['voices_contralto'], voices_urls['voices_tenor'], voices_urls['voices_bass'],
                    transpose, scroll_speed, songbook_id,
                    himno['createdAt'], himno['updatedAt']
                ))
                
                new_song_id = target_cursor.lastrowid
                himno_id_mapping[himno['id']] = new_song_id
            
            logger.info(f"Migrated {len(himnos)} songs")
            
            self.target_conn.commit()
            
            return himno_id_mapping
            
        except Exception as e:
            logger.error(f"Error migrating songs: {e}")
            raise
    
    def migrate_song_categories(self, himno_id_mapping: Dict[int, int], category_mappings: Dict[str, Dict[int, int]]):
        """Migrate tema_himnos and sub_tema_himnos to song_categories table, only for leaf categories."""
        try:
            source_cursor = self.source_conn.cursor()
            target_cursor = self.target_conn.cursor()
            
            # Get all categories that have children (parent categories)
            target_cursor.execute("SELECT DISTINCT parent_category_id FROM categories WHERE parent_category_id IS NOT NULL")
            parent_category_ids = set(row[0] for row in target_cursor.fetchall())
            
            # Migrate tema_himnos - only if the tema doesn't have subcategories (is a leaf)
            source_cursor.execute("SELECT * FROM tema_himnos")
            tema_himnos = source_cursor.fetchall()
            
            migrated_count = 0
            skipped_count = 0
            
            for tema_himno in tema_himnos:
                himno_id = tema_himno['himno_id']
                tema_id = tema_himno['tema_id']
                
                new_song_id = himno_id_mapping.get(himno_id)
                new_category_id = category_mappings['tema_id_mapping'].get(tema_id)
                
                # Skip if song or category doesn't exist
                if not new_song_id or not new_category_id:
                    logger.warning(f"Could not migrate tema_himno: himno_id={himno_id}, tema_id={tema_id}")
                    continue
                
                # Skip if this category has children (is not a leaf)
                if new_category_id in parent_category_ids:
                    logger.debug(f"Skipping tema_himno for parent category: tema_id={tema_id}, new_category_id={new_category_id}")
                    skipped_count += 1
                    continue
                
                # Only add relationship if category is a leaf (has no children)
                target_cursor.execute("""
                    INSERT INTO song_categories (song_id, category_id, created_at, updated_at)
                    VALUES (?, ?, ?, ?)
                """, (new_song_id, new_category_id, tema_himno['createdAt'], tema_himno['updatedAt']))
                migrated_count += 1
            
            logger.info(f"Migrated {migrated_count} tema-song relationships (skipped {skipped_count} parent categories)")
            
            # Migrate sub_tema_himnos - these should all be leaf categories
            source_cursor.execute("SELECT * FROM sub_tema_himnos")
            sub_tema_himnos = source_cursor.fetchall()
            
            migrated_count = 0
            
            for sub_tema_himno in sub_tema_himnos:
                himno_id = sub_tema_himno['himno_id']
                sub_tema_id = sub_tema_himno['sub_tema_id']
                
                new_song_id = himno_id_mapping.get(himno_id)
                new_category_id = category_mappings['sub_tema_id_mapping'].get(sub_tema_id)
                
                if new_song_id and new_category_id:
                    target_cursor.execute("""
                        INSERT INTO song_categories (song_id, category_id, created_at, updated_at)
                        VALUES (?, ?, ?, ?)
                    """, (new_song_id, new_category_id, sub_tema_himno['createdAt'], sub_tema_himno['updatedAt']))
                    migrated_count += 1
                else:
                    logger.warning(f"Could not migrate sub_tema_himno: himno_id={himno_id}, sub_tema_id={sub_tema_id}")
            
            logger.info(f"Migrated {migrated_count} sub-tema-song relationships")
            
            self.target_conn.commit()
            
        except Exception as e:
            logger.error(f"Error migrating song categories: {e}")
            raise
    
    def verify_migration(self):
        """Verify the migration was successful."""
        try:
            cursor = self.target_conn.cursor()
            
            # Check table counts
            tables_to_check = ['users', 'songbooks', 'categories', 'songs', 'song_categories']
            
            for table in tables_to_check:
                cursor.execute(f"SELECT COUNT(*) FROM {table}")
                count = cursor.fetchone()[0]
                logger.info(f"Table {table}: {count} records")
            
            # Check specific data
            cursor.execute("SELECT * FROM users WHERE admin = 1")
            admin_user = cursor.fetchone()
            if admin_user:
                logger.info(f"Admin user created: {admin_user['first_name']} {admin_user['last_name']}")
            
            cursor.execute("SELECT title FROM songbooks")
            songbooks = cursor.fetchall()
            logger.info(f"Songbooks created: {[sb['title'] for sb in songbooks]}")
            
            # Check song distribution
            cursor.execute("""
                SELECT s.title as songbook, COUNT(so.id) as song_count 
                FROM songbooks s 
                LEFT JOIN songs so ON s.id = so.songbook_id 
                GROUP BY s.id, s.title
            """)
            distribution = cursor.fetchall()
            for dist in distribution:
                logger.info(f"Songbook '{dist['songbook']}': {dist['song_count']} songs")
            
            logger.info("Migration verification completed successfully")
            
        except Exception as e:
            logger.error(f"Error during migration verification: {e}")
            raise
    
    def close_connections(self):
        """Close database connections."""
        if self.source_conn:
            self.source_conn.close()
        if self.target_conn:
            self.target_conn.close()
        logger.info("Database connections closed")
    
    def migrate(self):
        """Perform the complete migration."""
        try:
            logger.info("Starting database migration...")
            
            # Connect to databases
            self.connect_databases()
            
            # Create new schema
            self.create_new_schema()
            
            # Create initial data (users and songbooks)
            initial_data = self.create_initial_data()
            
            # Migrate categories
            category_mappings = self.migrate_categories(initial_data)
            
            # Migrate songs
            himno_id_mapping = self.migrate_songs(initial_data)

            # Migrate song-category relationships
            self.migrate_song_categories(himno_id_mapping, category_mappings)
            
            # Verify migration
            self.verify_migration()
            
            logger.info("Database migration completed successfully!")
            
        except Exception as e:
            logger.error(f"Migration failed: {e}")
            raise
        finally:
            self.close_connections()

def main():
    """Main function to run the migration."""
    # Get the directory of this script
    script_dir = os.path.dirname(os.path.abspath(__file__))
    
    # Database paths
    source_db = os.path.join(script_dir, '..', 'himnos_coros.sqlite')
    target_db = os.path.join(script_dir, '..', 'songbooks_of_praise.sqlite')
    
    # Convert to absolute paths
    source_db = os.path.abspath(source_db)
    target_db = os.path.abspath(target_db)
    
    logger.info(f"Source database: {source_db}")
    logger.info(f"Target database: {target_db}")
    
    # Create and run migrator
    migrator = DatabaseMigrator(source_db, target_db)
    
    try:
        migrator.migrate()
        print("\n✅ Migration completed successfully!")
        print(f"New database created at: {target_db}")
    except Exception as e:
        print(f"\n❌ Migration failed: {e}")
        sys.exit(1)

if __name__ == "__main__":
    main()
