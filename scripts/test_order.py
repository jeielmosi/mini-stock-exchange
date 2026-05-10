import subprocess
import time
import requests
import sys

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
            # Try to hit a simple endpoint or just check if port is open
            # Since we don't have a health check endpoint, we'll try to hit /orders
            # and ignore the 405/404 as long as the server responds.
            requests.get(API_URL + "/health", timeout=1)
            print("Application is up and running")
            return True
        except requests.exceptions.ConnectionError:
            time.sleep(1)
    return False

def test_order_flow():
    print("Testing Order Flow...")
    symbol = "AAPL"
    ask_price = 140.0
    bid_price = 150.0

    # 1. Create Ask Order
    ask_payload = {
        "broker_id": "broker1",
        "owner_doc": "doc1",
        "type": "ASK",
        "symbol": symbol,
        "price": ask_price,
        "quantity": 10,
        "valid_until": time.strftime("%Y-%m-%d")
    }
    resp = requests.post(f"{API_URL}/orders", json=ask_payload)
    assert resp.status_code == 201, f"Failed to create Ask order: {resp.text}"
    ask_id = resp.json().get("id")
    print(f"Created Ask order: {ask_id}")

    # 2. Create Matching Bid Order
    bid_payload = {
        "broker_id": "broker2",
        "owner_doc": "doc2",
        "type": "BID",
        "symbol": symbol,
        "price": bid_price,
        "quantity": 10,
        "valid_until": time.strftime("%Y-%m-%d")
    }
    resp = requests.post(f"{API_URL}/orders", json=bid_payload)
    assert resp.status_code == 201, f"Failed to create Bid order: {resp.text}"
    bid_id = resp.json().get("id")
    print(f"Created matching Bid order: {bid_id}")

    # Give time for matching
    time.sleep(1)

    # 3. Verify filled
    ask_order = requests.get(f"{API_URL}/orders/{ask_id}").json()
    print("ASK:", ask_order)
    assert ask_order.get("status") == "FILLED", f"Ask order not filled: {ask_order}"
    
    bid_order = requests.get(f"{API_URL}/orders/{bid_id}").json()
    print("BID:", bid_order)
    assert bid_order.get("status") == "FILLED", f"Bid order not filled: {bid_order}"
    print("Order flow matched and filled successfully")

def test_no_match():
    print("Testing No Match Flow...")
    symbol = "MSFT"
    ask_price = 200.0
    bid_price = 100.0

    # 1. Create Ask Order
    ask_payload = {
        "broker_id": "broker1",
        "owner_doc": "doc1",
        "type": "ASK",
        "symbol": symbol,
        "price": ask_price,
        "quantity": 10,
        "valid_until": time.strftime("%Y-%m-%d")
    }
    resp = requests.post(f"{API_URL}/orders", json=ask_payload)
    assert resp.status_code == 201
    ask_id = resp.json().get("id")

    # 2. Create non-matching Bid Order
    bid_payload = {
        "broker_id": "broker2",
        "owner_doc": "doc2",
        "type": "BID",
        "symbol": symbol,
        "price": bid_price,
        "quantity": 10,
        "valid_until": time.strftime("%Y-%m-%d")
    }
    resp = requests.post(f"{API_URL}/orders", json=bid_payload)
    assert resp.status_code == 201
    bid_id = resp.json().get("id")

    time.sleep(1)

    # 3. Verify pending
    ask_order = requests.get(f"{API_URL}/orders/{ask_id}").json()
    assert ask_order.get("status") == "PENDING", f"Ask order should be pending: {ask_order}"
    print("ASK:", ask_order)
    
    bid_order = requests.get(f"{API_URL}/orders/{bid_id}").json()
    assert bid_order.get("status") == "PENDING", f"Bid order should be pending: {bid_order}"
    print("BID:", ask_order)
    print("Non-matching orders remained pending as expected")

def main():
    try:
        if not run_command("make up"):
            sys.exit(1)
        
        if not run_command("make migrate"):
            sys.exit(1)
        
        if not wait_for_app():
            print("Application failed to start in time")
            sys.exit(1)
        
        test_order_flow()
        test_no_match()
        
        print("\nAll tests passed!")
    except Exception as e:
        print(f"\nTest failed: {e}")
        sys.exit(1)
    finally:
        input("Press Enter to shut down the application...")
        run_command("make down")

if __name__ == "__main__":
    main()
