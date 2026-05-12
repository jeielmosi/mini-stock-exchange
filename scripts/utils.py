import subprocess
import time
import requests

API_URL = "http://localhost:8080"

def run_command(command):
    print(f"Running: {command}")
    result = subprocess.run(command, shell=True, capture_output=True, text=True)
    if result.returncode != 0:
        print(f"Error running command: {command}")
        print(f"Stderr: {result.stderr}")
        return False
    return True

def wait_for_app(timeout=30):
    start_time = time.time()
    while time.time() - start_time < timeout:
        try:
            requests.get(API_URL + "/health", timeout=1)
            print("Application is up and running")
            return True
        except requests.exceptions.ConnectionError:
            time.sleep(1)
    return False
