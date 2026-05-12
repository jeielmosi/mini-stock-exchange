import time
import requests
from . import API_URL, newUUID


def run():
    print("Testing No Match Flow...")
    symbol = newUUID()
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
