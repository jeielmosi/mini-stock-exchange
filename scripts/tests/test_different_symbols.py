import time
import requests
from . import API_URL, newUUID


def run():
    print("Testing Different Symbols Flow...")
    symbol1 = newUUID()
    symbol2 = newUUID()

    # 1. Create Ask for symbol1
    ask_payload = {
        "broker_id": "b1",
        "owner_doc": "d1",
        "type": "ASK",
        "symbol": symbol1,
        "price": 100.0,
        "quantity": 10,
        "valid_until": time.strftime("%Y-%m-%d")
    }
    resp = requests.post(f"{API_URL}/orders", json=ask_payload)
    ask_id = resp.json().get("id")

    # 2. Create Bid for symbol2 (Matching price)
    bid_payload = {
        "broker_id": "b2",
        "owner_doc": "d2",
        "type": "BID",
        "symbol": symbol2,
        "price": 110.0,
        "quantity": 10,
        "valid_until": time.strftime("%Y-%m-%d")
    }
    resp = requests.post(f"{API_URL}/orders", json=bid_payload)
    bid_id = resp.json().get("id")

    time.sleep(1)

    # 3. Verify both pending
    assert requests.get(f"{API_URL}/orders/{ask_id}").json().get("status") == "PENDING"
    assert requests.get(f"{API_URL}/orders/{bid_id}").json().get("status") == "PENDING"
    print("Different symbols did not match as expected")
