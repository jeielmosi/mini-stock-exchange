import os
import shutil
import sys
from utils import run_command

def main():
    volume_path = ".volume/postgresql"
    
    print("Erasing database...")
    
    if not run_command("docker compose down"):
        print("Failed to stop containers")
        sys.exit(1)
        
    if os.path.exists(volume_path):
        try:
            abs_volume_parent = os.path.abspath(".volume")
            if not run_command(f"docker run --rm -v {abs_volume_parent}:/mnt alpine rm -rf /mnt/postgresql"):
                raise Exception("Docker removal failed")
            print(f"Removed volume at {volume_path}")
        except Exception as e:
            print(f"Error removing volume: {e}")
            sys.exit(1)
    else:
        print("Volume directory not found, skipping removal")

    print("Restarting database and applying migrations...")
    if not run_command("docker compose up -d db"):
        print("Failed to start database")
        sys.exit(1)
        
    if not run_command("make migrate"):
        print("Failed to apply migrations")
        sys.exit(1)
        
    print("Database erased and reset successfully.")

if __name__ == "__main__":
    main()
