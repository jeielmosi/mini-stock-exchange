import time
import requests
from . import API_URL, newUUID

def run():
    print("Testing Order Flow...")
    symbol = newUUID()
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
