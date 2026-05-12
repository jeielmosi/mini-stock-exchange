import time
import requests
from . import API_URL, newUUID

def run():
    print("Testing Partial Fill Flow...")
    symbol = newUUID()
    ask_price = 140.0
    bid_price = 150.0

    # 1. Create Ask Order (Small quantity)
    ask_payload = {
        "broker_id": "broker1",
        "owner_doc": "doc1",
        "type": "ASK",
        "symbol": symbol,
        "price": ask_price,
        "quantity": 5,
        "valid_until": time.strftime("%Y-%m-%d")
    }
    resp = requests.post(f"{API_URL}/orders", json=ask_payload)
    assert resp.status_code == 201
    ask_id = resp.json().get("id")

    # 2. Create Matching Bid Order (Larger quantity)
    bid_payload = {
        "broker_id": "broker2",
        "owner_doc": "doc2",
        "type": "BID",
        "symbol": symbol,
        "price": bid_price,
        "quantity": 15,
        "valid_until": time.strftime("%Y-%m-%d")
    }
    resp = requests.post(f"{API_URL}/orders", json=bid_payload)
    assert resp.status_code == 201
    bid_id = resp.json().get("id")

    time.sleep(1)

    # 3. Verify
    ask_order = requests.get(f"{API_URL}/orders/{ask_id}").json()
    assert ask_order.get("status") == "FILLED", f"Ask order should be FILLED: {ask_order}"
    
    bid_order = requests.get(f"{API_URL}/orders/{bid_id}").json()
    print("BID status:", bid_order.get("status"))
    assert bid_order.get("status") != "FILLED", f"Bid order should not be fully filled: {bid_order}"
    print("Partial fill successful")
